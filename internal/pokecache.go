package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createAt time.Time
	val      []byte
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]cacheEntry
}

// reapLoop läuft kontinuierlich und entfernt Einträge, die älter als das angegebene Intervall sind.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.data {
			if time.Since(v.createAt) > interval {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}

// Erzeugt ein neues Cache-Objekt mit einem angegebenen Intervall.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}