package pokecache

import (
	"time"
)

type Cache struct {
	cache map[string]cache_entry
}

type cache_entry struct {
	data       []byte
	created_at time.Time
}

func NewCache() Cache {
	return Cache{
		cache: make(map[string]cache_entry),
	}
}

func (c *Cache) Add(key string, data []byte) {
	c.cache[key] = cache_entry{data: data, created_at: time.Now().UTC()}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return entry.data, true
}

