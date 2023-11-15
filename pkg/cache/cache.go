package cache

import (
	"context"
	"log"
	"sync"
	"testGRPC/pkg/api"
	"time"
)

type Cache struct {
	m sync.RWMutex

	expirationDelta time.Duration
	cleanupInterval time.Duration
	cache           map[string]*Item
}

type Item struct {
	timeAccessed time.Time
	data         api.Description
}

func NewCache(exInterval, cleanupInterval time.Duration) *Cache {
	return &Cache{
		expirationDelta: exInterval,
		cache:           make(map[string]*Item),
	}
}

func (c *Cache) Get(key string) (api.Description, bool) {
	c.m.RLock()
	defer c.m.RUnlock()

	data, ok := c.cache[key]
	if ok {
		c.cache[key].timeAccessed = time.Now()
		return data.data, ok
	}

	return api.Description{}, ok
}

func (c *Cache) Set(key string, value api.Description) {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = &Item{
		timeAccessed: time.Now(),
		data:         value,
	}
}

func (c *Cache) CleaningUp(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context closed. Exiting..")
			return
		case <-time.After(c.cleanupInterval):
			if c.cache == nil {
				return
			}

			if keys := c.expiredKeys(); len(keys) != 0 {
				c.clearItems(keys)
			}
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {

	c.m.RLock()

	defer c.m.RUnlock()

	for k, i := range c.cache {
		if time.Now().After(i.timeAccessed.Add(c.expirationDelta)) {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearItems(keys []string) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, key := range keys {
		delete(c.cache, key)
	}
}
