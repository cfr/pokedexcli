package pokecache

import (
    "time"
    "sync"
)

type Cache struct {
    mu sync.Mutex
    es map[string]cacheEntry
    interval time.Duration
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}

func NewCache(interval time.Duration) Cache {
    cache := Cache{es: make(map[string]cacheEntry), interval: interval}
    go reapLoop(&cache)
    return cache
}

func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.es[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    e, ok := c.es[key]
    if ok {
        return e.val, true
    } else {
        return nil, false
    }
}

func reapLoop(c *Cache) {
    ticker := time.NewTicker(c.interval)
    for range ticker.C {
        c.mu.Lock()
        defer c.mu.Unlock()
        for k, e := range c.es {
            if time.Since(e.createdAt) > c.interval {
                delete(c.es, k)
            }
        }
    }

}
