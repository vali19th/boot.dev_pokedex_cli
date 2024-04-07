package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cache_entry
	mux   *sync.Mutex
}

type cache_entry struct {
	data       []byte
	created_at time.Time
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cache_entry),
		mux:   &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, data []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cache_entry{data: data, created_at: time.Now().UTC()}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, ok := c.cache[key]
	if ok {
		return entry.data, true
	} else {
		return nil, false
	}

}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

func (c *Cache) reap(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	t := time.Now().UTC().Add(-interval)
	for k, v := range c.cache {
		if v.created_at.Before(t) {
			delete(c.cache, k)
		}
	}
}

