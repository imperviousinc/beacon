package hip5

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"testing"
)

func TestEthereum_Handler(t *testing.T) {
	eth, err := NewEthereum("https://mainnet.infura.io/v3/b0933ce6026a4e1e80e89e96a5d095bc")
	if err != nil {
		t.Fatal(err)
	}

	rr, err := eth.Handler(context.Background(), "humbly.eth.", dns.TypeA, ethNS[0])
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rr)
}
