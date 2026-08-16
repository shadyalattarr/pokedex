package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	myCache  map[string]CacheEntry
	interval time.Duration
	mu       *sync.RWMutex
}
