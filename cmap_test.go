package cmap

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/octu0/chanque"

	ccmap "github.com/orcaman/concurrent-map"
)

func BenchmarkCmapNoParallel(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	sequenceSet := func(c *CMap, start, end int) {
		for i := start; i < end; i += 1 {
			key := strconv.Itoa(i)
			c.Set(key, key)
		}
	}
	sequenceGet := func(c *CMap, start, end int) {
		for i := start; i < end; i += 1 {
			key := strconv.Itoa(i)
			c.Get(key)
		}
	}
	randomSet := func(c *CMap, count int) {
		for i := 0; i < count; i += 1 {
			key := strconv.Itoa(rand.Int())
			c.Set(key, key)
		}
	}
	randomGet := func(c *CMap, count int) {
		for i := 0; i < count; i += 1 {
			key := strconv.Itoa(rand.Int())
			c.Get(key)
		}
	}

	b.Run("sequence/rw/32", func(tb *testing.B) {
		c := New(WithSlabSize(32))
		for i := 0; i < tb.N; i += 1 {
			sequenceSet(c, 0, 5000)
			sequenceGet(c, 0, 5000)
			sequenceSet(c, 5000, 10000)
			sequenceGet(c, 5000, 10000)
		}
	})
	b.Run("sequence/rw/1024", func(tb *testing.B) {
		c := New(WithSlabSize(1024))
		for i := 0; i < tb.N; i += 1 {
			sequenceSet(c, 0, 5000)
			sequenceGet(c, 0, 5000)
			sequenceSet(c, 5000, 10000)
			sequenceGet(c, 5000, 10000)
		}
	})
	b.Run("random/rw/32", func(tb *testing.B) {
		c := New(WithSlabSize(32))
		for i := 0; i < tb.N; i += 1 {
			randomSet(c, 5000)
			randomSet(c, 5000)
			randomGet(c, 5000)
			randomGet(c, 5000)
		}
	})
	b.Run("random/rw/1024", func(tb *testing.B) {
		c := New(WithSlabSize(1024))
		for i := 0; i < tb.N; i += 1 {
			randomSet(c, 5000)
			randomSet(c, 5000)
			randomGet(c, 5000)
			randomGet(c, 5000)
		}
	})
}

func BenchmarkCmapParallel(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	sequenceSetWG := func(wg *sync.WaitGroup, c *CMap, start, end int) {
		defer wg.Done()

		for i := start; i < end; i += 1 {
			key := strconv.Itoa(i)
			c.Set(key, key)
		}
	}
	sequenceGetWG := func(wg *sync.WaitGroup, c *CMap, start, end int) {
		defer wg.Done()

		for i := start; i < end; i += 1 {
			key := strconv.Itoa(i)
			c.Get(key)
		}
	}

	randomSetWG := func(wg *sync.WaitGroup, c *CMap, count int) {
		defer wg.Done()

		for i := 0; i < count; i += 1 {
			key := strconv.Itoa(rand.Int())
			c.Set(key, key)
		}
	}
	randomGetWG := func(wg *sync.WaitGroup, c *CMap, count int) {
		defer wg.Done()

		for i := 0; i < count; i += 1 {
			key := strconv.Itoa(rand.Int())
			c.Get(key)
		}
	}

	b.Run("sequence/rw/32", func(tb *testing.B) {
		c := New(WithSlabSize(32))
		wg := new(sync.WaitGroup)
		for i := 0; i < tb.N; i += 1 {
			wg.Add(4)
			go sequenceSetWG(wg, c, 0, 5000)
			go sequenceGetWG(wg, c, 0, 5000)
			go sequenceSetWG(wg, c, 5000, 10000)
			go sequenceGetWG(wg, c, 5000, 10000)
		}
		wg.Wait()
	})
	b.Run("sequence/rw/1024", func(tb *testing.B) {
		c := New(WithSlabSize(1024))
		wg := new(sync.WaitGroup)
		for i := 0; i < tb.N; i += 1 {
			wg.Add(4)
			go sequenceSetWG(wg, c, 0, 5000)
			go sequenceGetWG(wg, c, 0, 5000)
			go sequenceSetWG(wg, c, 5000, 10000)
			go sequenceGetWG(wg, c, 5000, 10000)
		}
		wg.Wait()
	})
	b.Run("random/rw/32", func(tb *testing.B) {
		c := New(WithSlabSize(32))
		wg := new(sync.WaitGroup)
		for i := 0; i < tb.N; i += 1 {
			wg.Add(4)
			go randomSetWG(wg, c, 5000)
			go randomSetWG(wg, c, 5000)
			go randomGetWG(wg, c, 5000)
			go randomGetWG(wg, c, 5000)
		}
		wg.Wait()
	})
	b.Run("random/rw/1024", func(tb *testing.B) {
		c := New(WithSlabSize(1024))
		wg := new(sync.WaitGroup)
		for i := 0; i < tb.N; i += 1 {
			wg.Add(4)
			go randomSetWG(wg, c, 5000)
			go randomSetWG(wg, c, 5000)
			go randomGetWG(wg, c, 5000)
			go randomGetWG(wg, c, 5000)
		}
		wg.Wait()
	})
}

func BenchmarkCompare(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	gen := func(tb *testing.B, count int) []string {
		tb.StopTimer()
		defer tb.StartTimer()

		keys := make([]string, count)
		for i := 0; i < count; i += 1 {
			keys[i] = strconv.Itoa(rand.Int())
		}
		return keys
	}

	randomSetSyncMap := func(sm *sync.Map, keys []string) {
		for _, key := range keys {
			sm.Store(key, key)
		}
	}
	randomGetSyncMap := func(sm *sync.Map, keys []string) {
		for _, key := range keys {
			sm.Load(key)
		}
	}

	randomSetConcurrentMap := func(ccm ccmap.ConcurrentMap, keys []string) {
		for _, key := range keys {
			ccm.Set(key, key)
		}
	}
	randomGetConcurrentMap := func(ccm ccmap.ConcurrentMap, keys []string) {
		for _, key := range keys {
			ccm.Get(key)
		}
	}

	randomSetCMap := func(c *CMap, keys []string) {
		for _, key := range keys {
			c.Set(key, key)
		}
	}
	randomGetCMap := func(c *CMap, keys []string) {
		for _, key := range keys {
			c.Get(key)
		}
	}

	b.Run("ConcurrentMap", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		tb.Cleanup(func() {
			e.Release()
		})

		ccm := ccmap.New()
		for i := 0; i < tb.N; i += 1 {
			keys1 := gen(tb, 5000)
			keys2 := gen(tb, 5000)

			se := e.SubExecutor()
			se.Submit(func() {
				randomSetConcurrentMap(ccm, keys1)
			})
			se.Submit(func() {
				randomGetConcurrentMap(ccm, keys2)
			})
			se.Wait()
		}
	})
	b.Run("sync.Map", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		tb.Cleanup(func() {
			e.Release()
		})

		sm := new(sync.Map)
		for i := 0; i < tb.N; i += 1 {
			keys1 := gen(tb, 5000)
			keys2 := gen(tb, 5000)

			se := e.SubExecutor()
			se.Submit(func() {
				randomSetSyncMap(sm, keys1)
			})
			se.Submit(func() {
				randomGetSyncMap(sm, keys2)
			})
			se.Wait()
		}
	})
	b.Run("cmap", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		tb.Cleanup(func() {
			e.Release()
		})

		c := New()
		for i := 0; i < tb.N; i += 1 {
			keys1 := gen(tb, 5000)
			keys2 := gen(tb, 5000)

			se := e.SubExecutor()
			se.Submit(func() {
				randomSetCMap(c, keys1)
			})
			se.Submit(func() {
				randomGetCMap(c, keys2)
			})
			se.Wait()
		}
	})
}

func TestCmapSetGetRemove(t *testing.T) {
	c := New()
	if _, ok := c.Get("foobar"); ok {
		t.Errorf("key foobar not exists")
	}
	if _, ok := c.Remove("foobar"); ok {
		t.Errorf("key foobar not exists")
	}

	c.Set("foobar", "123456")
	c.Set("hello", "world")

	if v, ok := c.Get("foobar"); ok != true {
		t.Errorf("foobar exists")
	} else {
		if s, ok := v.(string); ok != true {
			t.Errorf("foobar is string")
		} else {
			if s != "123456" {
				t.Errorf("value is 123456")
			}
		}
	}

	if v := c.GetRLocked("foobar", func(ok bool, v interface{}) interface{} {
		return len(v.(string))
	}); v == nil {
		t.Errorf("foobar exists")
	} else {
		if s, ok := v.(int); ok != true {
			t.Errorf("foobar is int")
		} else {
			if s != 6 {
				t.Errorf("value is len(123456)")
			}
		}
	}

	if old, ok := c.Remove("foobar"); ok != true {
		t.Errorf("foobar exists")
	} else {
		if s, ok := old.(string); ok != true {
			t.Errorf("old value is string")
		} else {
			if s != "123456" {
				t.Errorf("old value is 123456")
			}
		}
	}

	if _, ok := c.Get("foobar"); ok {
		t.Errorf("foobar removed")
	}
	if ok := c.GetRLocked("foobar", func(ok bool, v interface{}) interface{} {
		return ok
	}); ok.(bool) {
		t.Errorf("foobar removed")
	}
	if _, ok := c.Remove("foobar"); ok {
		t.Errorf("foobar removed")
	}
	if _, ok := c.Get("hello"); ok != true {
		t.Errorf("key hello not removed")
	}
}

func TestCmapLenKeys(t *testing.T) {
	c := New()

	if c.Len() != 0 {
		t.Errorf("no keys")
	}
	if len(c.Keys()) != 0 {
		t.Errorf("no keys")
	}

	size := 1000
	for i := 0; i < size; i += 1 {
		key := strconv.Itoa(i)
		c.Set(key, key)
	}

	if c.Len() != size {
		t.Errorf("%d key set", size)
	}

	if len(c.Keys()) != size {
		t.Errorf("%d key set", size)
	}

	removeSize := size / 2
	for i := 0; i < removeSize; i += 1 {
		key := strconv.Itoa(i)
		c.Remove(key)
	}

	if c.Len() != removeSize {
		t.Errorf("%d key set", removeSize)
	}

	if len(c.Keys()) != removeSize {
		t.Errorf("%d key set", removeSize)
	}
}

func TestCmapUpsert(t *testing.T) {
	t.Run("exists", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		v := c.Upsert("foo", func(ok bool, oldValue interface{}) interface{} {
			if ok != true {
				tt.Errorf("foo is exists")
			}
			if s, ok := oldValue.(string); ok != true {
				tt.Errorf("old value is string")
			} else {
				if s != "bar" {
					tt.Errorf("old value is bar")
				}
			}
			return oldValue.(string) + "qwerty"
		})

		if s, ok := v.(string); ok != true {
			tt.Errorf("upsert value is string")
		} else {
			if s != "barqwerty" {
				tt.Errorf("new value is barqwerty")
			}
		}

		if v, ok := c.Get("foo"); ok != true {
			tt.Errorf("foo exists")
		} else {
			if v.(string) != "barqwerty" {
				tt.Errorf("new value updated")
			}
		}
	})
	t.Run("notexists", func(tt *testing.T) {
		c := New()

		v := c.Upsert("foobar", func(ok bool, oldValue interface{}) interface{} {
			if ok {
				tt.Errorf("foobar not exists")
			}
			if oldValue != nil {
				tt.Errorf("old value not exists")
			}
			return "quux"
		})

		if s, ok := v.(string); ok != true {
			tt.Errorf("upsert value is string")
		} else {
			if s != "quux" {
				tt.Errorf("new value is quux")
			}
		}

		if v, ok := c.Get("foobar"); ok != true {
			tt.Errorf("foobar exists")
		} else {
			if v.(string) != "quux" {
				tt.Errorf("new value updated")
			}
		}
	})
}

func TestCmapSetIfAbsent(t *testing.T) {
	t.Run("exists", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		if c.SetIfAbsent("foo", "12345") {
			tt.Errorf("foo exists")
		}

		if v, ok := c.Get("foo"); ok != true {
			tt.Errorf("foo exists")
		} else {
			if v.(string) != "bar" {
				tt.Errorf("not updated")
			}
		}
	})
	t.Run("notexists", func(tt *testing.T) {
		c := New()

		if c.SetIfAbsent("foobar", "12345") != true {
			tt.Errorf("no key = updated")
		}

		if v, ok := c.Get("foobar"); ok != true {
			tt.Errorf("foobar exists")
		} else {
			if v.(string) != "12345" {
				tt.Errorf("updated value 12345")
			}
		}
	})
}

func TestCmapSetIf(t *testing.T) {
	t.Run("notexists/noset", func(tt *testing.T) {
		c := New()

		c.SetIf("foo", func(exists bool, oldValue interface{}) (interface{}, bool) {
			if exists {
				tt.Errorf("key foo must not exists")
			}
			return "noset-testdata", false
		})
		if _, ok := c.Get("foo"); ok {
			tt.Errorf("must no set foo")
		}
	})
	t.Run("notexists/set", func(tt *testing.T) {
		c := New()

		c.SetIf("foo", func(exists bool, oldValue interface{}) (interface{}, bool) {
			if exists {
				tt.Errorf("key foo must not exists")
			}
			return "new-value", true
		})
		if v, ok := c.Get("foo"); ok != true {
			tt.Errorf("must foo exists")
		} else {
			if v.(string) != "new-value" {
				tt.Errorf("updated new-value")
			}
		}
	})
	t.Run("exists/noset", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		c.SetIf("foo", func(exists bool, oldValue interface{}) (interface{}, bool) {
			if exists != true {
				tt.Errorf("foo is exists")
			}
			if oldValue.(string) != "bar" {
				tt.Errorf("old value is bar")
			}
			return "noset-testdata", false
		})
		if v, ok := c.Get("foo"); ok != true {
			tt.Errorf("foo is exists")
		} else {
			if v.(string) != "bar" {
				tt.Errorf("not updated")
			}
		}
	})
	t.Run("exists/set", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		c.SetIf("foo", func(exists bool, oldValue interface{}) (interface{}, bool) {
			if exists != true {
				tt.Errorf("foo is exists")
			}
			if oldValue.(string) != "bar" {
				tt.Errorf("old value is bar")
			}
			return "new-value", true
		})
		if v, ok := c.Get("foo"); ok != true {
			tt.Errorf("foo is exists")
		} else {
			if v.(string) != "new-value" {
				tt.Errorf("updated new-value")
			}
		}
	})
}

func TestCmapRemoveIf(t *testing.T) {
	t.Run("exists/remove", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		removed := c.RemoveIf("foo", func(exists bool, oldValue interface{}) bool {
			if exists != true {
				tt.Errorf("key foo exists")
			}
			return true
		})
		if removed != true {
			tt.Errorf("key removed")
		}
		if _, ok := c.Get("foo"); ok {
			tt.Errorf("foo removed!")
		}
	})
	t.Run("exists/noremove", func(tt *testing.T) {
		c := New()
		c.Set("foo", "bar")

		removed := c.RemoveIf("foo", func(exists bool, oldValue interface{}) bool {
			if exists != true {
				tt.Errorf("key foo exists")
			}
			return false
		})
		if removed != false {
			tt.Errorf("key removed")
		}
		if _, ok := c.Get("foo"); ok != true {
			tt.Errorf("foo not removed")
		}
	})
	t.Run("notexists/remove", func(tt *testing.T) {
		c := New()

		removed := c.RemoveIf("foobar", func(exists bool, oldValue interface{}) bool {
			if exists {
				tt.Errorf("key foobar notexists")
			}
			return true
		})
		if removed {
			tt.Errorf("key not exists")
		}
		if _, ok := c.Get("foobar"); ok {
			tt.Errorf("foobar not exists")
		}
	})
	t.Run("notexists/noremove", func(tt *testing.T) {
		c := New()

		removed := c.RemoveIf("foobar", func(exists bool, oldValue interface{}) bool {
			if exists {
				tt.Errorf("key foobar notexists")
			}
			return false
		})
		if removed {
			tt.Errorf("key not exists")
		}
		if _, ok := c.Get("foobar"); ok {
			tt.Errorf("foobar not exists")
		}
	})
}
