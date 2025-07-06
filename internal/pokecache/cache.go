package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data            map[string]cacheEntry
	invalidInterval time.Duration

	mu sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	nc := Cache{data: make(map[string]cacheEntry), invalidInterval: interval, mu: sync.RWMutex{}}
	go nc.reapLoop()
	return &nc
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return val.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.invalidInterval)
	for t := range ticker.C {
		for k, v := range c.data {
			if t.Sub(v.createdAt) > c.invalidInterval {
				c.mu.Lock()
				delete(c.data, k)
				c.mu.Unlock()
			}
		}
	}
}
