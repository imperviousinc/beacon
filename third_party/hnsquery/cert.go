package hnsquery

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

type CertVerifyInfo struct {
	Host             string
	Port             string
	Protocol         string
	RawCerts         [][]byte
	DisableNameCheck bool
}

var ErrCertVerifyFailed = errors.New("verify failed")
var ErrDNSAuthFailed = errors.New("dns authentication failed")

type CertVerifier interface {
	Verify(ctx context.Context, verifyInfo *CertVerifyInfo) (bool, error)
}

type DNSCertVerifier struct {
	Resolver  *Resolver
	tlsaCache *lru.Cache
}

func NewDNSCertVerifier(resolver *Resolver) (*DNSCertVerifier, error) {
	c, err := lru.New(100)
	if err != nil {
		return nil, err
	}

	d := &DNSCertVerifier{
		Resolver:  resolver,
		tlsaCache: c,
	}

	return d, nil
}

func (d *DNSCertVerifier) lookupTLSAWithRetry(ctx context.Context, port, proto, name string) (rrs []*dns.TLSA, err error) {
	for i := 0; i < 3; i++ {
		rrs, err = d.lookupTLSAWithDowngrade(ctx, port, proto, name)
		if err == nil {
			return
		}
		if errors.Is(err, ErrDNSFatal) || errors.Is(err, ErrDNSSECFailed) {
			return
		}
	}
	return
}

func (d *DNSCertVerifier) lookupTLSAWithDowngrade(ctx context.Context, port, proto, name string) ([]*dns.TLSA, error) {
	qname, err := dns.TLSAName(name, port, proto)
	if err != nil {
		return nil, fmt.Errorf("invalid format: %w", ErrDNSFatal)
	}
	qname = strings.ToLower(qname)

	if rrs, ok := d.tlsaCache.Get(qname); ok {
		rc := rrs.(*recordCache)
		if time.Now().Before(rc.expire) {
			log.Printf("tlsa cache hit for name: %s, port: %s, protocol: %s", name, port, proto)
			return rc.rrs, nil
		}
		d.tlsaCache.Remove(qname)
	}

	// if zone is insecure, we ignore TLSA lookup result
	// some nameservers don't behave well
	// when asked about unfamiliar record types
	// downgrade lookup should populate tld cache
	// doing it in parallel with TLSA
	var zoneStatusErr error
	var zoneStatusDowngrade bool
	done := make(chan struct{}, 1)

	findZoneStatus := func() {
		zoneStatusDowngrade, zoneStatusErr = d.Resolver.ZoneInsecure(ctx, name)
		done <- struct{}{}
	}

	go findZoneStatus()

	tlsaRecords, tlsaErr := d.lookupTLSAInternal(ctx, qname)
	<-done

	if zoneStatusErr != nil {
		return nil, zoneStatusErr
	}
	if zoneStatusDowngrade {
		return nil, nil
	}

	if tlsaErr == nil {
		ttl := time.Second * 60
		if len(tlsaRecords) > 0 {
			ttl = time.Second * time.Duration(tlsaRecords[0].Hdr.Ttl)
		}
		d.tlsaCache.Add(qname, &recordCache{
			rrs:    tlsaRecords,
			expire: time.Now().Add(ttl),
		})
	}
	return tlsaRecords, tlsaErr
}

func (d *DNSCertVerifier) lookupTLSAInternal(ctx context.Context, qname string) ([]*dns.TLSA, error) {
	msg, err := d.Resolver.Query(ctx, qname, dns.TypeTLSA)
	if err != nil {
		return nil, err
	}

	// ignore if insecure
	if !msg.AuthenticatedData {
		return nil, nil
	}

	// hard fail if lookup isn't successful
	if msg.Rcode != dns.RcodeSuccess && msg.Rcode != dns.RcodeNameError {
		return nil, fmt.Errorf("received non-success rcode: %d", msg.Rcode)
	}

	var tlsas []*dns.TLSA
	for _, rr := range msg.Answer {
		if rr.Header().Rrtype == dns.TypeTLSA {
			tlsas = append(tlsas, rr.(*dns.TLSA))
		}
	}

	if len(tlsas) == 0 && len(msg.Answer) > 0 {
		// could be a partial CNAME response
		// but we expect resolver to give a full answer
		return nil, fmt.Errorf("got a response with no tlsa records")
	}

	return tlsas, nil
}

func (d *DNSCertVerifier) Verify(ctx context.Context, verifyInfo *CertVerifyInfo) (bool, error) {
	host := strings.ToLower(verifyInfo.Host)

	if len(verifyInfo.RawCerts) == 0 {
		return false, fmt.Errorf("no certificates specified: %w", ErrCertVerifyFailed)
	}

	// currently, we only verify DANE-EE so we need
	// the leaf certificate
	leaf := verifyInfo.RawCerts[0]
	if leaf == nil {
		return false, fmt.Errorf("no leaf certificate: %w", ErrCertVerifyFailed)
	}

	if verifyInfo.Host == "" || verifyInfo.Port == "" || verifyInfo.Protocol == "" {
		return false, fmt.Errorf("missing host, port or protocol: %w", ErrCertVerifyFailed)
	}

	cert, err := x509.ParseCertificate(leaf)
	if err != nil {
		return false, fmt.Errorf("cert parse error: %w", ErrCertVerifyFailed)
	}

	// name checks
	if !verifyInfo.DisableNameCheck {
		if err := cert.VerifyHostname(host); err != nil {
			return false, fmt.Errorf("%v: %w", err, ErrCertVerifyFailed)
		}
	}

	// fetch TLSA records
	rrs, err := d.lookupTLSAWithRetry(ctx, verifyInfo.Port, verifyInfo.Protocol, dns.Fqdn(host))
	if err != nil {
		return false, err
	}

	// no records we can't
	// verify this certificate
	if len(rrs) == 0 {
		return false, nil
	}

	supportedUsage := false
	for _, rr := range rrs {
		// multiple TLSA records may exist
		// only usage 3 is supported
		// unsupported usages are ignored
		if rr.Usage == 3 {
			supportedUsage = true

			if rr.Verify(cert) == nil {
				// success
				return true, nil
			}
		}
	}

	// no TLSA records with usage 3
	// it's safe to downgrade
	if !supportedUsage {
		return false, nil
	}

	return false, fmt.Errorf("no matching tlsa record: %w", ErrDNSAuthFailed)
}
