package hnsquery

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/imperviousinc/hnsquery/dnssec"
	"github.com/miekg/dns"
	"testing"
)

func TestNewResolver(t *testing.T) {
	c, _ := lru.New(100)
	called := false
	r := &Resolver{
		TrustAnchorFunc: func(ctx context.Context, cut string) (*dnssec.Zone, bool, error) {
			if called {
				t.Fatal("should only be called once")
				return nil, false, nil
			}

			called = true
			if cut == "." {
				return nil, false, fmt.Errorf("failed")
			}

			if cut == "proofofconcept." {
				z, err := dnssec.NewZone("proofofconcept.", nil)
				return z, true, err
			}

			return nil, false, nil
		},
		zoneCuts: c,
		exchangeTest: func(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
			t.Fatal("insecure zone shouldn't call exchange")
			return nil, nil
		},
	}

	_, err := r.Query(context.Background(), "_443._tcp.proofofconcept.", dns.TypeTLSA)
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.getTrustAnchor(context.Background(), "proofofconcept.")
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.Query(context.Background(), "_443._tcp.proofofconcept.", dns.TypeTLSA)
	if err != nil {
		t.Fatal(err)
	}
}
