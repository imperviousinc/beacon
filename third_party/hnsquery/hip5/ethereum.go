package hip5

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/miekg/dns"
	"strings"
	"time"
)

// hardcoded .eth NS rrset pointing to their registry
var ethNS = []*dns.NS{
	{
		Hdr: dns.RR_Header{
			Name:   "eth.",
			Rrtype: dns.TypeNS,
			Class:  1,
			Ttl:    86400,
		},
		Ns: "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e._eth.",
	},
}

type Ethereum struct {
	client *ethclient.Client
	// resolver cache
	rCache *cache
	// query cache
	qCache map[uint16]*cache
}

type queryCacheData struct {
	registry string
	rrs      []dns.RR
}

func NewEthereum(rawurl string) (*Ethereum, error) {
	conn, err := ethclient.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	e := &Ethereum{
		client: conn,
		rCache: newCache(200),
		qCache: make(map[uint16]*cache),
	}

	// caching lower level lookups only
	// other types should be cached by users
	// of this client
	e.qCache[dns.TypeCNAME] = newCache(50)
	e.qCache[dns.TypeNS] = newCache(100)
	e.qCache[dns.TypeDS] = newCache(100)
	return e, nil
}

func (e *Ethereum) GetResolverAddress(node, registryAddress string) (common.Address, error) {
	key := node + ";" + registryAddress
	r, ok := e.rCache.get(key)
	if ok {
		if time.Now().Before(r.ttl) {
			return r.msg.(common.Address), nil
		}
		e.rCache.remove(key)
	}

	registry, err := NewENSRegistry(common.HexToAddress(registryAddress), e.client)
	if err != nil {
		return common.Address{}, err
	}

	addr, err := registry.Resolver(nil, EnsNode(node))
	if err != nil {
		return common.Address{}, err
	}

	e.rCache.set(key, &entry{
		msg: addr,
		ttl: time.Now().Add(6 * time.Hour),
	})

	return addr, nil
}

func isZero(addr common.Address) bool {
	for _, b := range addr {
		if b != 0 {
			return false
		}
	}

	return true
}

func (e *Ethereum) Resolve(registry string, ra common.Address, qname string, qtype uint16, single bool) ([]dns.RR, error) {
	if isZero(ra) {
		return nil, nil
	}

	r, err := NewDNSResolver(ra, e.client)
	if err != nil {
		return nil, err
	}

	qname = dns.CanonicalName(qname)
	node := toNode(qname)
	nodeHash, err := NameHash(node)
	if err != nil {
		return nil, err
	}

	res, err := e.queryWithResolver(registry, r, nodeHash, qname, qtype, single)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Ethereum) checkQueryCache(registry string, qname string, qtype uint16) ([]dns.RR, bool) {
	c, ok := e.qCache[qtype]
	if !ok {
		return nil, false
	}
	entry, ok := c.get(qname)
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.ttl) {
		c.remove(qname)
		return nil, false
	}

	m := entry.msg.(*queryCacheData)
	if !strings.EqualFold(m.registry, registry) {
		c.remove(qname)
		return nil, false
	}

	return m.rrs, true
}

func (e *Ethereum) dnsRecord(registry string, r *DNSResolver, node [32]byte, qname string, qtype uint16) ([]dns.RR, error) {
	if rrs, ok := e.checkQueryCache(registry, qname, qtype); ok {
		return rrs, nil
	}

	qnameHash, err := hashDnsName(qname)
	if err != nil {
		return nil, err
	}

	raw, err := r.DnsRecord(nil, node, qnameHash, qtype)
	if err != nil {
		return nil, err
	}

	rrs := unpackRRSet(raw)

	if qtype == dns.TypeCNAME || qtype == dns.TypeNS || qtype == dns.TypeDS {
		e.qCache[qtype].set(qname, &entry{
			msg: &queryCacheData{
				registry: registry,
				rrs:      rrs,
			},
			ttl: time.Now().Add(getTTL(rrs)),
		})
	}

	return rrs, nil
}

func (e *Ethereum) queryWithResolver(registry string, r *DNSResolver,
	nodeHash [32]byte, qname string, qtype uint16, single bool) ([]dns.RR, error) {

	rawRecords, err := e.dnsRecord(registry, r, nodeHash, qname, qtype)
	if err != nil {
		return nil, err
	}
	if single {
		return rawRecords, nil
	}

	maxLabels := dns.CountLabel(qname)
	if maxLabels > 3 {
		maxLabels = 3
	}

	// Look for NS records up to maxLabels
	if len(rawRecords) == 0 {
		labels := 2
		for {
			if labels > maxLabels {
				break
			}

			name := dns.Fqdn(LastNLabels(qname, labels))
			labels++

			if rawRecords, err = e.dnsRecord(registry, r, nodeHash, name, dns.TypeNS); err != nil {
				return nil, err
			}

			// a delegation exists check if it's signed
			if len(rawRecords) > 0 {
				var dsSet []dns.RR
				if dsSet, err = e.dnsRecord(registry, r, nodeHash, name, dns.TypeDS); err != nil {
					return nil, err
				}

				if len(dsSet) > 0 {
					rawRecords = append(rawRecords, dsSet...)
				}

				return rawRecords, nil
			}
		}
	}

	if len(rawRecords) == 0 {
		// no records for original qname and no delegations
		// check if a CNAME exists
		if rawRecords, err = e.dnsRecord(registry, r, nodeHash, qname, dns.TypeCNAME); err != nil {
			return nil, err
		}
	}

	return rawRecords, nil
}

func (e *Ethereum) Handler(ctx context.Context, qname string, qtype uint16, ns *dns.NS, exact bool) ([]dns.RR, error) {
	registryAddress := FirstNLabels(ns.Ns, 1)
	node := toNode(qname)

	var resolverAddr common.Address
	var err error

	resolverAddr, err = e.GetResolverAddress(node, registryAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to get resolver address from registry %s: %v", registryAddress, err)
	}

	return e.Resolve(registryAddress, resolverAddr, qname, qtype, exact)
}
