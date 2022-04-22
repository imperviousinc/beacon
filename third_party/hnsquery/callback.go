package hnsquery

/*
   #include <stdio.h>
   #include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/imperviousinc/hnsquery/resource"
	"github.com/miekg/dns"
	"log"
	"net"
	"sync"
	"unsafe"
)

const HandshakeTTL uint32 = 21600

var base32Hex = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding)

type Callback struct {
	Func func(rrs []dns.RR, err error)
}

type CallbackFunc func(rrs []dns.RR, err error)

type cgoHSKAccess struct {
	callbacks map[string][]*CallbackFunc
	sync.RWMutex
}

type ctxTable struct {
	contexts map[uint64]*cgoHSKAccess
	sync.RWMutex
}

// used to fetch back callbacks
// coming from C
var ctxMap ctxTable

func init() {
	ctxMap.contexts = make(map[uint64]*cgoHSKAccess)
}

func newCGOHSK() *cgoHSKAccess {
	return &cgoHSKAccess{
		callbacks: make(map[string][]*CallbackFunc),
	}
}

func (c *cgoHSKAccess) getCallbacks(name string) ([]*CallbackFunc, bool) {
	c.RLock()
	cbs, ok := c.callbacks[name]
	c.RUnlock()

	return cbs, ok
}

func (c *cgoHSKAccess) addCallback(name string, cb *CallbackFunc) bool {
	c.Lock()
	defer c.Unlock()

	nameFuncs, ok := c.callbacks[name]

	if !ok {
		nameFuncs = []*CallbackFunc{cb}
		c.callbacks[name] = nameFuncs
		return true
	}

	c.callbacks[name] = append(nameFuncs, cb)
	return false
}

func (c *cgoHSKAccess) removeCallback(name string, cb *CallbackFunc) {
	c.Lock()
	defer c.Unlock()

	nameFuncs, ok := c.callbacks[name]
	if !ok {
		return
	}

	if len(nameFuncs) == 1 {
		delete(c.callbacks, name)
		return
	}

	for i, f := range nameFuncs {
		if f != cb {
			continue
		}

		if i < len(nameFuncs)-1 {
			copy(nameFuncs[i:], nameFuncs[i+1:])
		}
		nameFuncs[len(nameFuncs)-1] = nil
		c.callbacks[name] = nameFuncs[:len(nameFuncs)-1]
		return
	}

	return
}

//export cgoAfterResolve
func cgoAfterResolve(name *C.char, status C.int, exists C.int, data unsafe.Pointer, dataLen C.size_t, v unsafe.Pointer) {
	// hns_ctx is passed to v
	// we need ctx->id to fetch callbacks for this ctx
	ctxId := getContextId(v)

	// get hns cgo callbacks
	ctxMap.RLock()
	hnsCgo, ok := ctxMap.contexts[ctxId]
	ctxMap.RUnlock()

	if !ok || hnsCgo == nil {
		log.Printf("cgoAfterResolve go not callbacks for ctx with id %d", ctxId)
		return
	}

	var goName string = C.GoString(name)
	size := C.int(dataLen)

	callbacks, ok := hnsCgo.getCallbacks(goName)
	if !ok || len(callbacks) == 0 {
		return
	}

	if status != 0 {
		for _, cb := range callbacks {
			if cb == nil {
				continue
			}

			(*cb)(nil, fmt.Errorf("after resolve %s: %w", goName, hskCodeToError(status)))
		}
		return
	}

	if exists != 1 || size == 0 {
		for _, cb := range callbacks {
			if cb == nil {
				continue
			}

			(*cb)(nil, nil)
		}
		return
	}

	if size < 0 {
		err := fmt.Errorf("invalid response len = %d", int(size))
		for _, cb := range callbacks {
			if cb == nil {
				continue
			}

			(*cb)(nil, err)
		}
		return
	}

	buf := C.GoBytes(data, size)
	r := &resource.Resource{}
	err := r.Decode(bytes.NewReader(buf))

	if err != nil {
		err = fmt.Errorf("failed decoding resource: %v", err)
		for _, cb := range callbacks {
			if cb == nil {
				continue
			}

			(*cb)(nil, err)
		}
		return
	}

	for _, cb := range callbacks {
		if cb == nil {
			continue
		}
		// for now allocate new set of rrs
		// for every callback
		rrs := resourceToDNS(goName, r)
		(*cb)(rrs, nil)
	}
}

func ipToSynth(ip net.IP) string {
	if len(ip) == 0 {
		ip = net.ParseIP("0.0.0.0")
	}
	if ip.To4() != nil {
		ip = ip[len(ip)-net.IPv4len:]
	}

	synth := base32Hex.EncodeToString(ip)
	return fmt.Sprintf("_%s._synth.", synth)
}

func resourceToDNS(name string, res *resource.Resource) []dns.RR {
	var rrs []dns.RR
	name = dns.CanonicalName(name)

	for _, hr := range res.Records {
		hType := hr.Type()

		switch hType {
		case resource.RecordTypeDS:
			hnsRR := hr.(*resource.DSRecord)
			dnsRR := &dns.DS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeDS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Digest:     hex.EncodeToString(hnsRR.Digest),
				DigestType: hnsRR.DigestType,
				Algorithm:  hnsRR.Algorithm,
				KeyTag:     hnsRR.KeyTag,
			}
			rrs = append(rrs, dnsRR)
			continue
		case resource.RecordTypeNS:
			hnsRR := hr.(*resource.NSRecord)

			dnsRR := &dns.NS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Ns: dns.CanonicalName(hnsRR.NS),
			}
			rrs = append(rrs, dnsRR)
			continue
		case resource.RecordTypeGlue4:
			hnsRR := hr.(*resource.Glue4Record)
			dnsRR := &dns.NS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Ns: hnsRR.NS,
			}
			dnsRR2 := &dns.A{
				Hdr: dns.RR_Header{
					Name:   hnsRR.NS,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				A: hnsRR.Address,
			}
			rrs = append(rrs, dnsRR)
			rrs = append(rrs, dnsRR2)
			continue
		case resource.RecordTypeGlue6:
			hnsRR := hr.(*resource.Glue6Record)
			dnsRR := &dns.NS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Ns: hnsRR.NS,
			}
			dnsRR2 := &dns.A{
				Hdr: dns.RR_Header{
					Name:   hnsRR.NS,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				A: hnsRR.Address,
			}
			rrs = append(rrs, dnsRR)
			rrs = append(rrs, dnsRR2)
			continue
		case resource.RecordTypeSynth4:
			hnsRR := hr.(*resource.Synth4Record)
			synthOwner := ipToSynth(hnsRR.Address)
			dnsRR := &dns.NS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Ns: synthOwner,
			}
			dnsRR2 := &dns.A{
				Hdr: dns.RR_Header{
					Name:   synthOwner,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				A: hnsRR.Address,
			}
			rrs = append(rrs, dnsRR)
			rrs = append(rrs, dnsRR2)
			continue
		case resource.RecordTypeSynth6:
			hnsRR := hr.(*resource.Synth6Record)
			synthOwner := ipToSynth(hnsRR.Address)
			dnsRR := &dns.NS{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Ns: synthOwner,
			}
			dnsRR2 := &dns.A{
				Hdr: dns.RR_Header{
					Name:   synthOwner,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				A: hnsRR.Address,
			}
			rrs = append(rrs, dnsRR)
			rrs = append(rrs, dnsRR2)
			continue
		case resource.RecordTypeTXT:
			hnsRR := hr.(*resource.TXTRecord)
			dnsRR := &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    HandshakeTTL,
				},
				Txt: hnsRR.Entries,
			}
			rrs = append(rrs, dnsRR)
			continue
		}
	}

	return rrs
}
