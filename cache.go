package cmap

import (
	"sync"
)

type Cache interface {
	Lock()
	Unlock()

	RLock()
	RUnlock()

	Set(string, interface{})
	Get(string) (interface{}, bool)
	Remove(string) (interface{}, bool)

	Len() int
	Keys() []string
}

// compile check
var (
	_ Cache = (*defaultCache)(nil)
)

type defaultCache struct {
	mutex  *sync.RWMutex
	values map[string]interface{}
}

func (c *defaultCache) Lock() {
	c.mutex.Lock()
}

func (c *defaultCache) RLock() {
	c.mutex.RLock()
}

func (c *defaultCache) Unlock() {
	c.mutex.Unlock()
}

func (c *defaultCache) RUnlock() {
	c.mutex.RUnlock()
}

func (c *defaultCache) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *defaultCache) Get(key string) (interface{}, bool) {
	v, ok := c.values[key]
	return v, ok
}

func (c *defaultCache) Remove(key string) (interface{}, bool) {
	v, ok := c.values[key]
	delete(c.values, key)
	return v, ok
}

func (c *defaultCache) Len() int {
	return len(c.values)
}

func (c *defaultCache) Keys() []string {
	keys := make([]string, 0, len(c.values))
	for k, _ := range c.values {
		keys = append(keys, k)
	}
	return keys
}

func newDefaultCache(size int) *defaultCache {
	return &defaultCache{
		mutex:  new(sync.RWMutex),
		values: make(map[string]interface{}, size),
	}
}
