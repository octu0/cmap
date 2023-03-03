package cmap

type GetFunc func(exists bool, value interface{}) (processedValue interface{})

type UpsertFunc func(exists bool, oldValue interface{}) (newValue interface{})

type SetIfFunc func(exists bool, value interface{}) (newValue interface{}, isSetValue bool)

type RemoveIfFunc func(exists bool, value interface{}) bool

type CMap struct {
	s *slab
}

func New(funcs ...cmapOptionFunc) *CMap {
	opt := newDefaultOption()
	for _, fn := range funcs {
		fn(opt)
	}
	return &CMap{
		s: newSlab(opt),
	}
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

func (c *CMap) GetRLocked(key string, fn GetFunc) interface{} {
	m := c.s.GetShard(key)
	m.RLock()
	defer m.RUnlock()

	v, ok := m.Get(key)
	return fn(ok, v)
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

func (c *CMap) SetIf(key string, fn SetIfFunc) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	v, ok := m.Get(key)
	setValue, isSet := fn(ok, v)
	if isSet {
		m.Set(key, setValue)
	}
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
