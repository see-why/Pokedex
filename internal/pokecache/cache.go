package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	store map[string]cacheEntry
	mux   *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		store: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, ok := c.store[key]
	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reapOld(time.Now().UTC(), interval)
	}
}

func (c *Cache) reapOld(now time.Time, last time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for k, v := range c.store {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.store, k)
		}
	}
}
