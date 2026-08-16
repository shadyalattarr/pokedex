package pokecache

import (
	"sync"
	"time"
)

func NewCache(interval time.Duration) Cache {
	cacheToReturn := Cache{
		myCache:  map[string]CacheEntry{}, // or make(map[string]Entry)? both same
		interval: interval,
		mu:       &sync.RWMutex{},
	}
	go cacheToReturn.reapLoop() //? ?? ?? ? ?
	return cacheToReturn
}

func (c *Cache) Add(key string, val []byte) { // writing
	c.mu.Lock()
	defer c.mu.Unlock()
	c.myCache[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) { // reading
	c.mu.RLock()
	defer c.mu.RUnlock()
	cEntry, ok := c.myCache[key]
	return cEntry.val, ok
}

func (c *Cache) reapLoop() { // writing!
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for { // an interval has passed
		time := <-ticker.C
		c.mu.Lock()
		for key, cEntry := range c.myCache { // reading
			if cEntry.createdAt.Before(time) { // reading
				delete(c.myCache, key) // writing / editing
			}
		}
		c.mu.Unlock() // ?
	}

}
