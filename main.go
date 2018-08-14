package interviewcache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	m  map[string]item

	getter getter
}

type getter interface {
	Get(string) ([]byte, error)
}

type item struct {
	data       []byte
	expiration time.Time
}

func NewCache(getter getter) *Cache {
	return &Cache{m: make(map[string]item), getter: getter}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.m[key]
	if !ok {
		data, err := c.getter.Get(key)
		if err != nil {
			// log err
			return nil, false
		}
		go c.Set(key, data)
		return data, true
	}
	b := make([]byte, len(item.data))
	copy(b, item.data)

	if item.expiration.Before(time.Now()) {
		defer c.lazyLoad(key)
	}

	return b, true
}

func (c *Cache) lazyLoad(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := c.getter.Get(key)
	if err != nil {
		// log error
		return
	}
	go c.Set(key, data)
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	b := make([]byte, len(value))
	copy(b, value)
	c.m[key] = item{data: b, expiration: time.Now().Add(1 * time.Minute)}
}
