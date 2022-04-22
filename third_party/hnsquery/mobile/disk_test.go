package mobile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDiskCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping disk cache test")
		return
	}

	dc := &diskCache{}
	var err error
	if dc.dir, err = ioutil.TempDir("", "disk_cache_test"); err != nil {
		t.Fatalf("failed making cache dir: %v", err)
	}
	defer os.RemoveAll(dc.dir)

	dc.ttl = time.Second * 5
	dc.maxItems = 100

	// write
	for i := 0; i < 100; i++ {
		err := dc.set("foo"+strconv.Itoa(i), []byte("bar"+strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
	}

	// read
	for i := 0; i < 100; i++ {
		k := "foo" + strconv.Itoa(i)
		val, ttl, ok := dc.get(k)
		if !ok {
			t.Fatalf("want key %s", k)
			return
		}

		if ttl.Seconds() == 0 {
			t.Fatalf("got zero ttl, want ttl = %f", ttl.Seconds())
			return
		}

		exp := "bar" + strconv.Itoa(i)
		if exp != string(val) {
			t.Fatalf("got cache value = %s, want %s", string(val), exp)
			return
		}
	}

	// read not exit
	if _, _, ok := dc.get("notexist"); ok {
		t.Fatalf("got an item that shouldn't exist")
	}

	// adding an item to a full cache
	dc.set("item", []byte("foo"))

	time.Sleep(time.Second)

	if dc.count != 0 {
		t.Fatalf("cache should be empty")
	}
}

func TestDiskCacheTTL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping disk cache test")
		return
	}

	dc := &diskCache{}
	var err error
	if dc.dir, err = ioutil.TempDir("", "disk_cache_test"); err != nil {
		t.Fatalf("failed making cache dir: %v", err)
	}
	defer os.RemoveAll(dc.dir)

	dc.ttl = time.Second * 2
	dc.maxItems = 100
	dc.set("foo", []byte("bar"))

	_, ttl, ok := dc.get("foo")
	if !ok {
		t.Fatalf("should be ok")
	}

	if ttl.Seconds() == 0 {
		t.Fatalf("bad ttl")
	}

	time.Sleep(2 * time.Second)

	_, ttl, ok = dc.get("foo")
	if !ok {
		fmt.Println("should be ok")
	}

	if ttl.Seconds() != 0 {
		t.Fatalf("got ttl = %f, stale item should have zero ttl", ttl.Seconds())
	}
}
