package in_memory_cache

import (
	"sync"
	"time"
)

type Writingtothecache struct {
	data         []byte
	creationtime time.Time
}
type InMemoryCache struct {
	mu    sync.RWMutex
	cache map[int64]Writingtothecache
	ttl   time.Duration
}

func NewInMemoryCache(expiration time.Duration) *InMemoryCache {
	cache := &InMemoryCache{
		cache: make(map[int64]Writingtothecache),
		ttl:   expiration,
	}
	go func() {
		t := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-t.C:

				for key, data := range cache.cache {
					if time.Since(data.creationtime) > cache.ttl {
						delete(cache.cache, key)
					}
				}
			default:

			}
		}
	}()
	return cache
}

func (c *InMemoryCache) Get(key int64) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return data.data, true
}

func (c *InMemoryCache) Set(key int64, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = Writingtothecache{
		data:         data,
		creationtime: time.Now(),
	}
}

func (c *InMemoryCache) Delete(key int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
