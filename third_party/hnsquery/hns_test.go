package hnsquery

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
		return
	}

	c, err := NewClient(&Config{
		DataDir: os.TempDir(),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Destroy()

	c.Run()
}

func TestClient_GetZone(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping get zone integration test")
		return
	}

	c, err := NewClient(&Config{
		DataDir: os.TempDir(),
	})

	if err != nil {
		t.Fatal(err)
	}
	defer c.Destroy()

	ready := make(chan error)
	c.Start(ready)

	<-ready

	names := []string{"proofofconcept",
		"3b",
		"3b",
		"schematic",
		"nb",
		"tlsa",
		"letsdane",
		"forever",
	}

	// duplicate
	names = append(names, names...)

	var wg sync.WaitGroup
	wg.Add(len(names))
	for _, name := range names {
		go func(zone string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			rrs, err := c.GetZone(ctx, zone)
			if err != nil {
				t.Error(err)
				return
			}

			if len(rrs) == 0 {
				t.Fatal("got no records")
			}

			for _, rr := range rrs {
				fmt.Println(rr)
			}

		}(name)
	}

	wg.Wait()
	fmt.Println("done")
}

func TestClient_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping client start integration test")
		return
	}

	c, err := NewClient(&Config{
		DataDir: os.TempDir(),
	})

	if err != nil {
		t.Fatal(err)
	}
	defer c.Destroy()

	ready := make(chan error)
	c.Start(ready)

	<-ready

	if !c.Ready() {
		t.Fatal("chain should be ready")
	}

	if c.ActivePeerCount() == 0 {
		t.Fatal("got no active peers")
	}

	fmt.Println("Progress: ", c.Progress())
	fmt.Println("Peers: ", c.PeerCount())
	fmt.Println("Active Peers: ", c.ActivePeerCount())
	fmt.Println("Block height: ", c.Height())

	fmt.Println("proofofconcept zone:")

	zone, err := c.GetZone(context.Background(), "proofofconcept")

	if err != nil {
		t.Fatal(err)
	}

	for _, rr := range zone {
		fmt.Println(rr)
	}
}
