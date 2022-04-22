package mobile

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type diskCache struct {
	dir      string
	ttl      time.Duration
	maxItems uint32

	// items written during this session
	count uint32

	nextCleanup time.Time
	sync.RWMutex
}

func newDiskCache(dir string) (*diskCache, error) {
	c := &diskCache{
		dir:      dir,
		ttl:      time.Hour * 6,
		maxItems: 100,
	}

	return c, nil
}

func (d *diskCache) set(key string, data []byte) error {
	defer d.maybeCleanUp()
	atomic.AddUint32(&d.count, 1)
	curr := atomic.LoadUint32(&d.count)
	if curr > d.maxItems {
		log.Printf("disk cache: cache is full")
		return nil
	}

	p := d.toPath(key)
	err := ioutil.WriteFile(p, data, 0644)
	if err != nil {
		return fmt.Errorf("failed writing cache: %v", err)
	}
	return nil
}

func (d *diskCache) get(key string) ([]byte, time.Duration, bool) {
	defer d.maybeCleanUp()

	p := d.toPath(key)
	f, err := os.Open(p)
	if err != nil {
		return nil, 0, false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, 0, false
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, 0, false
	}

	mod := stat.ModTime()
	expire := mod.Add(d.ttl)
	if time.Now().After(expire) {
		return data, 0, true
	}

	return data, expire.Sub(time.Now()), true
}

func (d *diskCache) maybeCleanUp() {
	d.RLock()
	due := d.nextCleanup
	d.RUnlock()

	curr := atomic.LoadUint32(&d.count)
	// skip if cache isn't full and clean up time isn't due yet
	if curr < d.maxItems && !due.IsZero() && time.Now().Before(due) {
		return
	}

	d.Lock()
	d.nextCleanup = time.Now().Add(3 * time.Hour)
	d.Unlock()

	atomic.StoreUint32(&d.count, 0)

	go d.cleanUp(false)
}

func (d *diskCache) cleanUp(force bool) {
	infos, err := ioutil.ReadDir(d.dir)
	if err != nil {
		log.Printf("disk cache: clean up failed: %v", err)
		return
	}

	// allow stale records up to 24 hours
	cutoff := time.Hour * 24
	if force {
		cutoff = 0
	}

	if len(infos) > int(d.maxItems) {
		cutoff = time.Nanosecond
	}

	deleted := 0
	now := time.Now()
	for _, info := range infos {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			deleted++
			log.Printf("disk cache: deleting old file %s", info.Name())
			err := os.Remove(path.Join(d.dir, info.Name()))
			if err != nil {
				log.Printf("disk cache: failed deleting old file %s: %v", info.Name(), err)
			}
		}
	}

	log.Printf("disk cache: clean up completed")
}

func (d *diskCache) toPath(key string) string {
	return path.Join(d.dir, base64.RawStdEncoding.EncodeToString([]byte(key)))
}
