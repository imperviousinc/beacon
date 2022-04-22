package hip5

import (
	"sync"
	"time"
)

type entry struct {
	msg interface{}
	ttl time.Time
}

type cache struct {
	m    map[string]*entry
	maxN int

	sync.RWMutex
}

func newCache(maxN int) (m *cache) {
	return &cache{m: make(map[string]*entry), maxN: maxN}
}

func (c *cache) set(key string, item *entry) {
	c.Lock()
	defer c.Unlock()

	if c.maxN == len(c.m) {
		for k := range c.m {
			delete(c.m, k)
			break
		}
	}

	c.m[key] = item
}

func (c *cache) get(key string) (*entry, bool) {
	c.RLock()
	defer c.RUnlock()

	i, ok := c.m[key]
	return i, ok
}

func (c *cache) remove(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, key)
}

func (c *cache) len() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.m)
}
