package hnsquery

import (
	"github.com/miekg/dns"
	"testing"
)

func TestCallback(t *testing.T) {
	c := newCGOHSK()

	cb1 := CallbackFunc(func(rrs []dns.RR, err error) {})
	cb2 := CallbackFunc(func(rrs []dns.RR, err error) {})
	cb3 := CallbackFunc(func(rrs []dns.RR, err error) {})

	// addCallback should return true
	// for initial
	initial := c.addCallback("test", &cb1)
	if !initial {
		t.Fatal("got initial = false, want true")
	}

	if cbs, ok := c.getCallbacks("test"); !ok || len(cbs) != 1 {
		t.Fatal("want callbacks len = 1")
	}

	// remove only callback in list
	c.removeCallback("test", &cb1)
	if _, ok := c.getCallbacks("test"); ok {
		t.Fatal("got callbacks, want none")
	}

	// add two callbacks to test2
	c.addCallback("test2", &cb1)
	if init2 := c.addCallback("test2", &cb2); init2 {
		t.Fatal("an initial callback exists")
	}

	// single callback to test3
	c.addCallback("test3", &cb3)

	if cbs, ok := c.getCallbacks("test2"); !ok || len(cbs) != 2 {
		t.Fatal("want callbacks len = 2")
	}

	// removing one callback shouldn't mess with the others
	c.removeCallback("test2", &cb1)
	cbs, ok := c.getCallbacks("test2")
	if !ok || len(cbs) != 1 {
		t.Fatal("want callbacks len = 1")
	}

	// the right callback should've been removed
	if cbs[0] != &cb2 {
		t.Fatal("removed incorrect callback")
	}

	// last callback for test2 should be removed
	c.removeCallback("test2", &cb2)
	if _, ok := c.getCallbacks("test2"); ok {
		t.Fatal("want no callbacks")
	}

	// remove only callback for test3
	c.removeCallback("test3", &cb2)
	if _, ok := c.getCallbacks("test3"); ok {
		t.Fatal("want no callbacks")
	}

	if len(c.callbacks) != 0 {
		t.Fatal("want no callbacks in map")
	}
}
