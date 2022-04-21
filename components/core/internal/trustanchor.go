package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/imperviousinc/hnsquery"
	"github.com/imperviousinc/hnsquery/dnssec"
	"github.com/imperviousinc/hnsquery/hip5"
	"github.com/miekg/dns"
)

type tldCacheEntry struct {
	expire time.Time
	rrs    []dns.RR
	hip5   bool
}

type ZoneQuery interface {
	GetZone(ctx context.Context, name string) (rrs []dns.RR, err error)
}

type RootZoneConfig struct {
	client      ZoneQuery
	eth         *hip5.Ethereum
	tldMemCache *lru.Cache
}

func queryTLDWithCache(ctx context.Context, h *RootZoneConfig, name string) ([]dns.RR, time.Duration, error) {
	name = dns.CanonicalName(name)
	// remove dot
	name = name[:len(name)-1]
	if dns.CountLabel(name) != 1 {
		return nil, 0, fmt.Errorf("not a tld")
	}

	if res, ok := h.tldMemCache.Get(name); ok {
		entry := res.(tldCacheEntry)
		if time.Now().Before(entry.expire) {
			log.Printf("mem cache hit for name %s", name)
			return entry.rrs, time.Now().Sub(entry.expire), nil
		}
		h.tldMemCache.Remove(name)
	}

	rrs, ttl, err := queryTLD(ctx, h, name)
	if err != nil {
		return nil, 0, err
	}

	h.tldMemCache.Add(name, tldCacheEntry{
		expire: time.Now().Add(ttl),
		rrs:    rrs,
	})

	return rrs, ttl, nil
}

func queryTLD(ctx context.Context, h *RootZoneConfig, name string) (rrs []dns.RR, ttl time.Duration, err error) {
	name = dns.CanonicalName(name)
	// remove dot
	name = name[:len(name)-1]
	if dns.CountLabel(name) != 1 {
		return nil, 0, fmt.Errorf("not a tld")
	}

	if rrs, err = h.client.GetZone(ctx, name); err != nil {
		return
	}

	return
}

func getPowTrustAnchor(h *RootZoneConfig) hnsquery.TrustAnchorPointFunc {
	return func(ctx context.Context, cut string) (*dnssec.Zone, error) {
		if cut == "." {
			return RootAnchor(ctx, h, cut)
		}

		// cut should be a TLD
		if dns.CountLabel(cut) != 1 {
			return nil, nil
		}

		cut = strings.TrimSuffix(strings.ToLower(cut), ".")
		rrs, ttl, err := queryTLDWithCache(ctx, h, cut)
		if err != nil {
			return nil, err
		}
		if len(rrs) == 0 {
			return nil, nil
		}

		var dsSet []dns.RR
		for _, rr := range rrs {
			if rr.Header().Rrtype == dns.TypeDS {
				dsSet = append(dsSet, rr)
			}
		}

		zone, err := dnssec.NewZone(cut, dsSet)
		if err == nil {
			zone.Expire = time.Now().Add(ttl)
		}

		return zone, err
	}
}

func rootVerify(ctx context.Context, h *RootZoneConfig, msg *dns.Msg) (bool, error) {
	qname := msg.Question[0].Name
	qtype := msg.Question[0].Qtype

	labelOffs := dns.Split(qname)
	if len(labelOffs) == 0 {
		return false, fmt.Errorf("root can't verify itself")
	}

	tld := qname[labelOffs[len(labelOffs)-1]:]

	// if qname = tld
	if len(labelOffs) == 1 {
		if qtype != dns.TypeDS && qtype != dns.TypeTXT {
			return false, fmt.Errorf("can't verify type %d from root", qtype)
		}

		rrs, _, err := queryTLDWithCache(ctx, h, tld)
		if err != nil {
			return false, err
		}

		if len(rrs) == 0 {
			// name doesn't exist on handshake
			// icann fallback
			return false, fmt.Errorf("no trust anchors for this zone")
		}

		msg.Rcode = dns.RcodeSuccess
		msg.AuthenticatedData = true
		msg.Authoritative = false
		msg.CheckingDisabled = false
		msg.Truncated = false
		msg.Answer = []dns.RR{}
		msg.Ns = []dns.RR{}
		msg.Extra = []dns.RR{}

		for _, rr := range rrs {
			if rr.Header().Rrtype == qtype {
				msg.Answer = append(msg.Answer, rr)
			}
		}

		return true, nil
	}

	// could be a HIP-5 SLD
	rrs, _, err := queryTLDWithCache(ctx, h, tld)
	if err != nil {
		return false, err
	}

	var hip5NS []*dns.NS
	for _, rr := range rrs {
		if rr.Header().Rrtype == dns.TypeDS {
			return false, fmt.Errorf("bad zone cut for name %s", tld)
		}
		if rr.Header().Rrtype == dns.TypeNS {
			ns := rr.(*dns.NS)
			if strings.HasSuffix(ns.Ns, "._eth.") {
				hip5NS = append(hip5NS, ns)
			}
		}
	}
	if len(hip5NS) == 0 {
		return false, fmt.Errorf("cannot verify %s", tld)
	}

	// No HIP-5 Handler, trust upstream resolver
	// this is equivalent to using a gateway
	// if the upstream resolver is over a secure channel
	// and the DoH has the same level of trust
	// as the gateway
	if h.eth == nil {
		return msg.AuthenticatedData, nil
	}

	// HIP-5 record look for the type we should request
	if len(msg.Answer) != 0 {
		var t uint16
		for _, rr := range msg.Answer {
			if !strings.EqualFold(rr.Header().Name, qname) {
				continue
			}
			rtype := rr.Header().Rrtype

			if rtype == qtype || rtype == dns.TypeCNAME {
				t = rr.Header().Rrtype
				break
			}
		}

		rrs, err := h.eth.Handler(ctx, qname, t, hip5NS[0], true)
		if err != nil {
			return false, fmt.Errorf("hip-5: %v", err)
		}

		msg.Rcode = dns.RcodeSuccess
		msg.CheckingDisabled = false
		msg.AuthenticatedData = true
		msg.Truncated = false
		msg.Answer = rrs
		msg.Ns = []dns.RR{}
		msg.Extra = []dns.RR{}
		return true, nil
	}

	if len(msg.Ns) == 0 {
		return false, fmt.Errorf("nothing to verify empty message")
	}
	rrs, err = h.eth.Handler(ctx, qname, qtype, hip5NS[0], true)
	if err != nil {
		return false, fmt.Errorf("hip-5: %v", err)
	}

	if len(rrs) == 0 {
		msg.Ns = []dns.RR{}
		msg.Extra = []dns.RR{}
		msg.Rcode = dns.RcodeSuccess
		return true, nil
	}

	return false, fmt.Errorf("hip-5: record exists")
}

var RootAnchor = func(ctx context.Context, v *RootZoneConfig, cut string) (*dnssec.Zone, error) {
	zone, err := dnssec.NewZone(cut, nil)
	if err != nil {
		return nil, err
	}
	zone.VerifyCallback = func(ctx context.Context, msg *dns.Msg) (bool, error) {
		return rootVerify(ctx, v, msg)
	}
	return zone, nil
}
