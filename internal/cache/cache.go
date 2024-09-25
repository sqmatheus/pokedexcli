package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries   map[string]cacheEntry
	entriesMu sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{entries: make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.entriesMu.Lock()
	defer c.entriesMu.Unlock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.entriesMu.Lock()
	defer c.entriesMu.Unlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.entriesMu.Lock()
		for key, entry := range c.entries {
			if t.After(entry.createdAt.Add(interval)) {
				delete(c.entries, key)
			}
		}
		c.entriesMu.Unlock()
	}
}
