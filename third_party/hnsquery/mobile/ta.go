package mobile

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/imperviousinc/hnsquery"
	"github.com/imperviousinc/hnsquery/dnssec"
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

type tldCacheEntry struct {
	expire time.Time
	rrs    []dns.RR
	hip5   bool
}

func queryTLDWithCache(ctx context.Context, h *HNS, name string) ([]dns.RR, time.Duration, error) {
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

func queryTLD(ctx context.Context, h *HNS, name string) (rrs []dns.RR, ttl time.Duration, err error) {
	name = dns.CanonicalName(name)
	// remove dot
	name = name[:len(name)-1]
	if dns.CountLabel(name) != 1 {
		return nil, 0, fmt.Errorf("not a tld")
	}

	var ok bool
	var cached []byte

	if cached, ttl, ok = h.tldDiskCache.get(name); ok {
		// stale item
		if ttl.Seconds() == 0 {
			if rrs, err = h.client.GetZone(ctx, name); err != nil {
				// serving stale on error with zero ttl
				log.Printf("disk stale cache hit %s", name)
				rrs, err = bytesToRecords(cached)
				return
			}
			// found good record update cache and return
			go h.tldDiskCache.set(name, recordsToBytes(rrs))
			return
		}

		// cached item
		log.Printf("disk cache hit %s", name)
		rrs, err = bytesToRecords(cached)
		return
	}

	if rrs, err = h.client.GetZone(ctx, name); err != nil {
		return
	}

	ttl = 6 * time.Hour
	// no disk caching for zero records
	if len(rrs) == 0 {
		return
	}

	go func() {
		err := h.tldDiskCache.set(name, recordsToBytes(rrs))
		if err != nil {
			log.Printf("failed storing key %s: %v", name, err)
		}
	}()

	return

}

func getPowTrustAnchor(h *HNS) hnsquery.TrustAnchorCallbackFunc {
	return func(ctx context.Context, cut string) (*dnssec.Zone, bool, error) {
		if cut == "." {
			return RootAnchor(ctx, h, cut)
		}

		// cut should be a TLD
		if dns.CountLabel(cut) != 1 {
			return nil, false, nil
		}

		cut = strings.TrimSuffix(strings.ToLower(cut), ".")
		rrs, ttl, err := queryTLDWithCache(ctx, h, cut)
		if err != nil {
			return nil, false, err
		}
		if len(rrs) == 0 {
			return nil, false, nil
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

		return zone, true, err
	}
}

func rootVerify(ctx context.Context, h *HNS, msg *dns.Msg) (bool, error) {
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

	if h.secureChannel {
		return true, nil
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

var RootAnchor = func(ctx context.Context, v *HNS, cut string) (*dnssec.Zone, bool, error) {
	zone, err := dnssec.NewZone(cut, nil)
	if err != nil {
		return nil, false, err
	}
	zone.VerifyCallback = func(ctx context.Context, msg *dns.Msg) (bool, error) {
		return rootVerify(ctx, v, msg)
	}
	return zone, true, nil
}

func bytesToRecords(buf []byte) (rrs []dns.RR, err error) {
	sc := bufio.NewScanner(bytes.NewReader(buf))

	for sc.Scan() {
		var rr dns.RR
		if rr, err = dns.NewRR(sc.Text()); err != nil {
			return
		}
		rrs = append(rrs, rr)
	}
	return
}

func recordsToBytes(rrs []dns.RR) []byte {
	var b bytes.Buffer

	for _, rr := range rrs {
		b.WriteString(rr.String())
		b.WriteRune('\n')
	}

	return b.Bytes()
}
