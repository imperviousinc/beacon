package hnsquery

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/imperviousinc/hnsquery/dnssec"
	"github.com/miekg/dns"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TrustAnchorPointFunc Must return a valid dnssec.Zone or nil
// if the zone provably does not exist. It must return an error
// in all other cases
type TrustAnchorPointFunc func(ctx context.Context, cut string) (*dnssec.Zone, error)

var ErrDNSFatal = errors.New("unrecoverable lookup error")
var ErrDNSSECFailed = errors.New("dnssec verify failed")

type Resolver struct {
	http             http.Client
	url              *url.URL
	dns              dns.Client
	forward          string
	CheckingDisabled bool
	TrustAnchorPointHandler  TrustAnchorPointFunc
	zoneCuts         *lru.Cache

	// for testing
	exchangeTest func(ctx context.Context, msg *dns.Msg) (*dns.Msg, error)
}

type recordCache struct {
	rrs    []*dns.TLSA
	expire time.Time
}

func (r *Resolver) queryInternal(ctx context.Context, qname string, qtype uint16) (re *dns.Msg, err error) {
	msg := new(dns.Msg)
	msg.SetQuestion(qname, qtype)
	msg.SetEdns0(4096, true)
	msg.AuthenticatedData = !r.CheckingDisabled
	msg.CheckingDisabled = r.CheckingDisabled

	return r.ExchangeContext(ctx, msg)
}

func (r *Resolver) verifyMessage(ctx context.Context, msg *dns.Msg) error {
	zone, err := r.findZone(ctx, msg.Question[0].Name, msg)
	if err != nil {
		return err
	}

	secure, err := zone.Verify(ctx, msg, msg.Question[0].Name, msg.Question[0].Qtype)
	if err != nil {
		return err
	}

	msg.Authoritative = false
	msg.CheckingDisabled = false
	msg.AuthenticatedData = secure
	return nil
}

func (r *Resolver) ZoneInsecure(ctx context.Context, qname string) (bool, error) {
	// try tld first
	labels := dns.Split(qname)
	if len(labels) > 0 {
		tld := dns.CanonicalName(qname[labels[len(labels)-1]:])
		z, err := r.getTrustAnchor(ctx, tld)
		if err != nil {
			return false, err
		}

		// doesn't exist
		if z == nil {
			return true, nil
		}

		// if tld is secure, don't downgrade
		if z.Secure() {
			return false, nil
		}

		// zone == qname no need to find zone
		if strings.EqualFold(qname, z.Name) {
			return !z.Secure(), nil
		}
	}

	// check if we have an insecure cut
	zone, err := r.findZone(ctx, qname, nil)
	return err == nil && !zone.Secure(), nil
}

func (r *Resolver) Query(ctx context.Context, qname string, qtype uint16) (msg *dns.Msg, err error) {
	qname = dns.CanonicalName(qname)
	if msg, err = r.queryInternal(ctx, qname, qtype); err != nil {
		return
	}
	if len(msg.Question) != 1 {
		return nil, fmt.Errorf("bad question section")
	}
	if !strings.EqualFold(msg.Question[0].Name, qname) || msg.Question[0].Qtype != qtype {
		return nil, fmt.Errorf("question mismatch")
	}

	answerSection := msg.Answer
	if err := r.verifyMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrDNSSECFailed)
	}

	// response may have a CNAME chain that was omitted
	// during validation
	chase := msg.AuthenticatedData && len(answerSection) > len(msg.Answer)
	if !chase {
		return msg, nil
	}

	// look for cnames
	target := ""
	for _, rr := range msg.Answer {
		if rr.Header().Rrtype == dns.TypeCNAME {
			cname := rr.(*dns.CNAME)
			target = cname.Target
			break
		}
	}

	if target == "" {
		return msg, nil
	}

	attempts := 0

chaseCNAME:
	for {
		if attempts > 10 {
			return nil, fmt.Errorf("max attempts reached chasing cname chain from target: %s", target)
		}

		attempts++
		if target == "" {
			break
		}

		targetAnswer := extractSecureRRSet(target, qtype, answerSection)
		if len(targetAnswer) > 0 {
			targetMsg := new(dns.Msg)
			targetMsg.SetQuestion(target, qtype)
			targetMsg.Answer = targetAnswer
			if err := r.verifyMessage(ctx, targetMsg); err != nil {
				return nil, fmt.Errorf("faild verifying cname target: %s: %w", target, ErrDNSSECFailed)
			}

			msg.Answer = append(msg.Answer, targetMsg.Answer...)
			msg.AuthenticatedData = msg.AuthenticatedData && targetMsg.AuthenticatedData

			break
		}

		targetCNAME := extractSecureRRSet(target, dns.TypeCNAME, answerSection)
		if len(targetCNAME) > 0 {
			targetMsg := new(dns.Msg)
			targetMsg.SetQuestion(target, qtype)
			targetMsg.Answer = targetCNAME
			if err := r.verifyMessage(ctx, targetMsg); err != nil {
				return nil, fmt.Errorf("faild verifying cname target: %s: %w", target, ErrDNSSECFailed)
			}

			target = ""
			msg.Answer = append(msg.Answer, targetMsg.Answer...)
			msg.AuthenticatedData = msg.AuthenticatedData && targetMsg.AuthenticatedData
			for _, rr := range targetMsg.Answer {
				if rr.Header().Rrtype == dns.TypeCNAME {
					cname := rr.(*dns.CNAME)
					target = cname.Target
					continue chaseCNAME
				}
			}
		}

		break
	}

	return msg, nil
}

func extractSecureRRSet(sname string, stype uint16, section []dns.RR) []dns.RR {
	var rrs []dns.RR
	for _, rr := range section {
		if !strings.EqualFold(rr.Header().Name, sname) {
			continue
		}

		if stype == rr.Header().Rrtype {
			rrs = append(rrs, rr)
		}

		if dns.TypeRRSIG == rr.Header().Rrtype && rr.(*dns.RRSIG).TypeCovered == stype {
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

func findCutFromMsg(qname string, msg *dns.Msg) string {
	for _, rr := range msg.Answer {
		if rr.Header().Rrtype == dns.TypeRRSIG {
			sig := rr.(*dns.RRSIG)
			// if signer name is in bailiwick use it
			if strings.EqualFold(sig.Header().Name, qname) && dns.IsSubDomain(sig.SignerName, qname) {
				return sig.SignerName
			}
		}
	}

	if len(msg.Answer) == 0 {
		for _, rr := range msg.Ns {
			if rr.Header().Rrtype == dns.TypeSOA {
				if dns.IsSubDomain(rr.Header().Name, qname) {
					return rr.Header().Name
				}
			}
		}
	}

	return ""
}

func (r *Resolver) queryCut(ctx context.Context, qname string) (cut string, err error) {
	sname := qname
	for {
		msg := new(dns.Msg)
		msg.SetQuestion(sname, dns.TypeSOA)

		soa, soaErr := r.ExchangeContext(ctx, msg)
		if soaErr != nil {
			return "", fmt.Errorf("failed finding zone cut: %v", soaErr)
		}

		for _, rr := range soa.Answer {
			if rr.Header().Rrtype == dns.TypeSOA && strings.EqualFold(rr.Header().Name, qname) {
				return rr.Header().Name, nil
			}
		}

		for _, rr := range soa.Ns {
			if rr.Header().Rrtype == dns.TypeSOA && dns.IsSubDomain(rr.Header().Name, qname) {
				return rr.Header().Name, nil
			}
		}

		off, end := dns.NextLabel(sname, 0)
		if end {
			return ".", nil
		}
		sname = sname[off:]
	}
}

func (r *Resolver) findCut(ctx context.Context, qname string, msg *dns.Msg) (string, error) {
	if qname == "." {
		return ".", nil
	}

	if msg != nil {
		if cut := findCutFromMsg(qname, msg); cut != "" {
			return cut, nil
		}
	}

	return r.queryCut(ctx, qname)
}

func (r *Resolver) loadKeys(ctx context.Context, zone *dnssec.Zone) error {
	msg := new(dns.Msg)
	msg.SetEdns0(4096, true)
	msg.SetQuestion(zone.Name, dns.TypeDNSKEY)
	msg.AuthenticatedData = true

	re, err := r.ExchangeContext(ctx, msg)
	if err != nil {
		return err
	}

	keys, err := zone.VerifyDNSKeys(re)
	if err != nil {
		return err
	}

	zone.Keys = keys
	return nil
}

func (r *Resolver) zoneFromDS(ctx context.Context, name string, dsSet []dns.RR) (*dnssec.Zone, error) {
	zone, err := dnssec.NewZone(name, dsSet)
	if err != nil {
		return nil, err
	}

	if err = r.loadKeys(ctx, zone); err != nil {
		return nil, err
	}

	return zone, nil
}

func (r *Resolver) getTrustAnchor(ctx context.Context, cut string) (zone *dnssec.Zone, err error) {
	if zone, ok := r.zoneCuts.Get(cut); ok {
		zone := zone.(*dnssec.Zone)

		if time.Now().Before(zone.Expire) {
			log.Printf("resolver cache hit for cut: %s", cut)
			return zone, nil
		}

		r.zoneCuts.Remove(cut)
	}

	if r.TrustAnchorPointHandler == nil {
		return nil, fmt.Errorf("no trust anchor callback set")
	}
	if zone, err = r.TrustAnchorPointHandler(ctx, cut); err != nil {
		return nil, fmt.Errorf("failed getting trust anchor: %w", err)
	}
	if zone == nil {
		return nil, nil
	}
	if len(zone.TrustAnchors) > 0 && len(zone.Keys) == 0 {
		if err = r.loadKeys(ctx, zone); err != nil {
			return nil, fmt.Errorf("failed getting dnskeys for zone %s", zone.Name)
		}
	}

	log.Printf("resolver caching cut: %s", cut)
	r.zoneCuts.Add(cut, zone)
	return
}

func (r *Resolver) findZone(ctx context.Context, qname string, msg *dns.Msg) (*dnssec.Zone, error) {
	sname := qname
	var pendingValidation []string
	var baseZone *dnssec.Zone

loopFindCut:
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed finding zone for %s: context cancelled", qname)
		default:
			cut, err := r.findCut(ctx, sname, msg)
			if err != nil {
				return nil, err
			}

			baseZone, err = r.getTrustAnchor(ctx, cut)
			if err != nil {
				return nil, err
			}

			if baseZone != nil {
				break loopFindCut
			}

			pendingValidation = append(pendingValidation, cut)
			msg = nil
			off, end := dns.NextLabel(cut, 0)
			if end {
				sname = "."
				continue
			}

			sname = cut[off:]
		}
	}

	if len(pendingValidation) > 0 {
		log.Printf("find zone: pending validation: %v", pendingValidation)
		return r.verifyChain(ctx, pendingValidation, baseZone)
	}

	return baseZone, nil
}

// verifyChain verifies all zone cuts and returns the last verified zone cut
func (r *Resolver) verifyChain(ctx context.Context, cuts []string, base *dnssec.Zone) (*dnssec.Zone, error) {
	insecure := false

	for i := len(cuts) - 1; i >= 0; i-- {
		cut := cuts[i]

		// if the parent cut was insecure
		// mark the remaining cuts insecure
		if insecure {
			log.Printf("verify chain: cut %s has unsigned parent - marking insecure", cut)
			insecureCut, err := dnssec.NewZone(cut, nil)
			if err != nil {
				return nil, err
			}

			base = insecureCut
			r.zoneCuts.Add(cut, insecureCut)
			continue
		}

		msg := new(dns.Msg)
		msg.SetQuestion(cut, dns.TypeDS)
		msg.SetEdns0(4096, true)
		msg.AuthenticatedData = true

		re, err := r.ExchangeContext(ctx, msg)
		if err != nil {
			return nil, err
		}

		log.Printf("verify chain: verifying cut %s with zone %s", cut, base.Name)

		secure, err := base.Verify(ctx, re, cut, dns.TypeDS)
		if err != nil {
			return nil, err
		}

		if secure {
			if base, err = r.zoneFromDS(ctx, cut, re.Answer); err != nil {
				return nil, err
			}

			r.zoneCuts.Add(cut, base)
			continue
		}

		log.Printf("verify chain: cut %s insecure verified by %s", cut, base.Name)

		// first insecure cut found in the chain
		insecure = true
		insecureCut, err := dnssec.NewZone(cut, nil)
		if err != nil {
			return nil, err
		}

		r.zoneCuts.Add(cut, insecureCut)
		base = insecureCut
	}

	return base, nil
}

type ResolverConfig struct {
	Forward string
}

func NewResolver(config *ResolverConfig) (r *Resolver, err error) {
	r = &Resolver{}
	r.dns.SingleInflight = true
	r.dns.Net = "doh"
	r.http.Timeout = time.Second * 10

	r.url, err = url.Parse(config.Forward)

	r.zoneCuts, err = lru.New(300)
	if err != nil {
		return nil, fmt.Errorf("failed cache init: %v", err)
	}

	return
}

func (r *Resolver) ExchangeContext(ctx context.Context, msg *dns.Msg) (re *dns.Msg, err error) {
	if r.exchangeTest != nil {
		return r.exchangeTest(ctx, msg)
	}

	for i := 0; i < 3; i++ {
		if r.dns.Net == "doh" {
			re, _, err = r.dns.ExchangeWithConn(msg, &dns.Conn{Conn: &dohConn{
				endpoint: r.url,
				http:     &r.http,
				ctx:      ctx,
			}})
		} else {
			re, _, err = r.dns.ExchangeContext(ctx, msg, r.forward)
		}

		if err == nil {
			if re.Truncated {
				err = errors.New("response truncated")
				continue
			}

			return
		}
	}

	return
}
