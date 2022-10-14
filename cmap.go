package cmap

import (
	"runtime"
	"sync"
)

type UpsertFunc func(exists bool, oldValue interface{}) (newValue interface{})

type RemoveIfFunc func(exists bool, value interface{}) bool

type CMap struct {
	opt *cmapOption
	s   *slab
	p   *gopool
}

func (c *CMap) Set(key string, value interface{}) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	m.Set(key, value)
}

func (c *CMap) Get(key string) (interface{}, bool) {
	m := c.s.GetShard(key)
	m.RLock()
	defer m.RUnlock()

	return m.Get(key)
}

func (c *CMap) Remove(key string) (interface{}, bool) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	return m.Remove(key)
}

func (c *CMap) Len() int {
	count := 0
	for _, m := range c.s.Shards() {
		m.RLock()
		count += m.Len()
		m.RUnlock()
	}
	return count
}

func (c *CMap) Keys() []string {
	shards := c.s.Shards()
	keys := make([]string, 0, len(shards))
	for _, m := range shards {
		m.RLock()
		keys = append(keys, m.Keys()...)
		m.RUnlock()
	}
	return keys
}

func (c *CMap) KeysParallel() []string {
	chKeys := make(chan []string)
	keys := make(chan []string)

	// read chan keys
	c.p.Submit(func(ch chan []string, result chan []string) taskFunc {
		return func() {
			buf := make([]string, 0, c.opt.slabSize)
			for {
				select {
				case keys, ok := <-ch:
					if ok != true {
						result <- buf
						close(result)
						return
					}
					buf = append(buf, keys...)
				}
			}
		}
	}(chKeys, keys))

	wg := new(sync.WaitGroup)
	for _, shard := range c.s.Shards() {
		wg.Add(1)

		// write chan keys
		c.p.Submit(func(w *sync.WaitGroup, m Cache, ch chan []string) taskFunc {
			return func() {
				defer w.Done()

				m.RLock()
				defer m.RUnlock()

				ch <- m.Keys()
			}
		}(wg, shard, chKeys))
	}
	wg.Wait()
	close(chKeys)

	return <-keys
}

func (c *CMap) Upsert(key string, fn UpsertFunc) (newValue interface{}) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	oldValue, ok := m.Get(key)
	newValue = fn(ok, oldValue)
	m.Set(key, newValue)
	return
}

func (c *CMap) SetIfAbsent(key string, value interface{}) (updated bool) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Get(key); ok != true {
		m.Set(key, value)
		return true
	}
	return false
}

func (c *CMap) RemoveIf(key string, fn RemoveIfFunc) (removed bool) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	v, ok := m.Get(key)
	remove := fn(ok, v)
	if remove && ok {
		m.Remove(key)
		return true
	}
	return false
}

func New(funcs ...cmapOptionFunc) *CMap {
	opt := newDefaultOption()
	for _, fn := range funcs {
		fn(opt)
	}
	c := &CMap{
		opt: opt,
		s:   newSlab(opt),
		p:   newGopool(opt.gopoolSize),
	}
	runtime.SetFinalizer(c, func(me *CMap) {
		me.p.Close()
	})
	return c
}
