package cmap

type UpsertFunc func(exists bool, oldValue interface{}) (newValue interface{})

type RemoveIfFunc func(exists bool, value interface{}) bool

type CMap interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
	Remove(string) (interface{}, bool)
	Len() int
	Keys() []string

	Upsert(string, UpsertFunc) interface{}
	SetIfAbsent(string, interface{}) bool
	RemoveIf(string, RemoveIfFunc) bool
}

// compile check
var (
	_ CMap = (*defaultCMap)(nil)
)

type defaultCMap struct {
	s *slab
}

func New(funcs ...cmapOptionFunc) CMap {
	opt := newDefaultOption()
	for _, fn := range funcs {
		fn(opt)
	}
	return &defaultCMap{
		s: newSlab(opt),
	}
}

func (c *defaultCMap) Set(key string, value interface{}) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	m.Set(key, value)
}

func (c *defaultCMap) Get(key string) (interface{}, bool) {
	m := c.s.GetShard(key)
	m.RLock()
	defer m.RUnlock()

	return m.Get(key)
}

func (c *defaultCMap) Remove(key string) (interface{}, bool) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	return m.Remove(key)
}

func (c *defaultCMap) Len() int {
	count := 0
	for _, m := range c.s.Shards() {
		m.RLock()
		count += m.Len()
		m.RUnlock()
	}
	return count
}

func (c *defaultCMap) Keys() []string {
	shards := c.s.Shards()
	keys := make([]string, 0, len(shards))
	for _, m := range shards {
		m.RLock()
		keys = append(keys, m.Keys()...)
		m.RUnlock()
	}
	return keys
}

func (c *defaultCMap) Upsert(key string, fn UpsertFunc) (newValue interface{}) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	oldValue, ok := m.Get(key)
	newValue = fn(ok, oldValue)
	m.Set(key, newValue)
	return
}

func (c *defaultCMap) SetIfAbsent(key string, value interface{}) (updated bool) {
	m := c.s.GetShard(key)
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Get(key); ok != true {
		m.Set(key, value)
		return true
	}
	return false
}

func (c *defaultCMap) RemoveIf(key string, fn RemoveIfFunc) (removed bool) {
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
