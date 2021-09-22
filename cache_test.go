package cmap

import (
	"sort"
	"strconv"
	"testing"
	"time"
)

type testLatch struct {
	ch chan struct{}
}

func (l *testLatch) Release() {
	close(l.ch)
}

func (l *testLatch) Wait() {
	<-l.ch
}

func newTestLatch() *testLatch {
	return &testLatch{make(chan struct{}, 0)}
}

func TestCacheLockUnlock(t *testing.T) {
	c := newDefaultCache(16)
	locker := func(boot, blockLatch *testLatch, dur time.Duration) {
		boot.Wait()

		c.Lock()
		blockLatch.Release()

		// blocking time
		time.Sleep(dur)

		c.Unlock()
	}
	nolocker := func(boot, blockLatch *testLatch, dur time.Duration) {
		boot.Wait()
		blockLatch.Release()
		time.Sleep(dur)
	}

	blocker := func(blockLatch *testLatch, ch chan time.Duration) {
		blockLatch.Wait()

		e := time.Now()
		c.Lock()
		c.Unlock()
		ch <- time.Since(e)
	}

	testBlock := func(tt *testing.T, blockTime time.Duration) {
		d := make(chan time.Duration)
		startup := newTestLatch()
		blockLatch := newTestLatch()
		go locker(startup, blockLatch, blockTime)
		go blocker(blockLatch, d)

		startup.Release()
		dur := <-d
		if dur < blockTime {
			tt.Errorf("no block %s", dur)
		}
		tt.Logf("block time = %s", dur)
	}

	t.Run("30ms", func(tt *testing.T) {
		testBlock(tt, 30*time.Millisecond)
	})
	t.Run("100ms", func(tt *testing.T) {
		testBlock(tt, 100*time.Millisecond)
	})
	t.Run("nolock", func(tt *testing.T) {
		blockTime := 150 * time.Millisecond

		d := make(chan time.Duration)
		startup := newTestLatch()
		blockLatch := newTestLatch()
		go nolocker(startup, blockLatch, blockTime)
		go blocker(blockLatch, d)

		startup.Release()
		dur := <-d
		if dur < blockTime {
			tt.Logf("no block %s", dur)
		} else {
			tt.Errorf("blocked %s", dur)
		}
	})
}

func TestCacheRlockRUnlock(t *testing.T) {
	c := newDefaultCache(16)
	locker := func(boot, blockLatch *testLatch, dur time.Duration) {
		boot.Wait()

		c.Lock()
		blockLatch.Release()

		// blocking time
		time.Sleep(dur)

		c.Unlock()
	}
	lockerR := func(boot, blockLatch *testLatch, dur time.Duration) {
		boot.Wait()

		c.RLock()
		blockLatch.Release()

		// blocking time
		time.Sleep(dur)

		c.RUnlock()
	}

	blocker := func(blockLatch *testLatch, ch chan time.Duration) {
		blockLatch.Wait()

		e := time.Now()
		c.Lock()
		c.Unlock()
		ch <- time.Since(e)
	}

	blockerR := func(blockLatch *testLatch, ch chan time.Duration) {
		blockLatch.Wait()

		e := time.Now()
		c.RLock()
		c.RUnlock()
		ch <- time.Since(e)
	}

	t.Run("lock+blockR", func(tt *testing.T) {
		blockTime := 50 * time.Millisecond

		d := make(chan time.Duration)
		startup := newTestLatch()
		blockLatch := newTestLatch()
		go locker(startup, blockLatch, blockTime)
		go blockerR(blockLatch, d)

		startup.Release()
		dur := <-d
		if dur < blockTime {
			tt.Errorf("must block %s", dur)
		}
	})
	t.Run("lockR+block", func(tt *testing.T) {
		blockTime := 50 * time.Millisecond

		d := make(chan time.Duration)
		startup := newTestLatch()
		blockLatch := newTestLatch()
		go lockerR(startup, blockLatch, blockTime)
		go blocker(blockLatch, d)

		startup.Release()
		dur := <-d
		if dur < blockTime {
			tt.Errorf("must block %s", dur)
		}
	})
	t.Run("lockR+blockR", func(tt *testing.T) {
		blockTime := 50 * time.Millisecond

		d := make(chan time.Duration)
		startup := newTestLatch()
		blockLatch := newTestLatch()
		go lockerR(startup, blockLatch, blockTime)
		go blockerR(blockLatch, d)

		startup.Release()
		dur := <-d
		if blockTime <= dur {
			tt.Errorf("must no block %s", dur)
		}
	})
}

func TestCacheSetGetRemove(t *testing.T) {
	c := newDefaultCache(0)

	if _, ok := c.Get("foo"); ok {
		t.Errorf("foo not set")
	}
	if _, ok := c.Remove("foo"); ok {
		t.Errorf("foo not exists")
	}

	c.Set("foo", "bar")
	if v, ok := c.Get("foo"); ok != true {
		t.Errorf("key foo not exists")
	} else {
		if s, ok := v.(string); ok != true {
			t.Errorf("value bar is string")
		} else {
			if s != "bar" {
				t.Errorf("value is bar")
			}
			t.Logf("foo = %s", s)
		}
	}

	c.Set("12345", 456)
	if v, ok := c.Get("12345"); ok != true {
		t.Errorf("key 12345 not exists")
	} else {
		if i, ok := v.(int); ok != true {
			t.Errorf("value 456 is int")
		} else {
			if i != 456 {
				t.Errorf("value is 456")
			}
			t.Logf("12345 = %d", i)
		}
	}

	if old, ok := c.Remove("foo"); ok != true {
		t.Errorf("foo exists")
	} else {
		if s, ok := old.(string); ok != true {
			t.Errorf("old value is string")
		} else {
			if s != "bar" {
				t.Errorf("old value is bar")
			}
			t.Logf("old value = %s", s)
		}
	}

	if _, ok := c.Remove("foo"); ok {
		t.Errorf("foo already removed")
	}
}

func TestCacheLen(t *testing.T) {
	c1 := newDefaultCache(16)
	c2 := newDefaultCache(0)

	if c1.Len() != 0 {
		t.Errorf("capacity != len")
	}
	if c2.Len() != 0 {
		t.Errorf("no keys")
	}

	size := 100
	for i := 0; i < size; i += 1 {
		c1.Set(strconv.Itoa(i), i)
		c2.Set(strconv.Itoa(i), i)
	}

	if c1.Len() != size {
		t.Errorf("%d keys setup", size)
	}
	if c2.Len() != size {
		t.Errorf("%d keys setup", size)
	}
	t.Logf("c1 = %d c2 = %d", c1.Len(), c2.Len())

	for i := 0; i < size; i += 1 {
		c1.Remove(strconv.Itoa(i))
		c2.Remove(strconv.Itoa(i))
	}

	if c1.Len() != 0 {
		t.Errorf("all keys removed")
	}
	if c2.Len() != 0 {
		t.Errorf("all keys removed")
	}
	t.Logf("c1 = %d c2 = %d", c1.Len(), c2.Len())
}

func TestCacheKeys(t *testing.T) {
	c := newDefaultCache(8)

	if ks := c.Keys(); 0 < len(ks) {
		t.Errorf("no key")
	}

	keys := []string{
		"foo1",
		"foo2",
		"foo3",
	}
	sort.Strings(keys)

	for _, k := range keys {
		c.Set(k, k)
	}

	if ks := c.Keys(); len(ks) < 1 {
		t.Errorf("keys setup")
	} else {
		sort.Strings(ks)
		if len(ks) != len(keys) {
			t.Errorf("all keys")
		}
		if ks[0] != keys[0] {
			t.Errorf("key[0] = %s actual = %s", keys[0], ks[0])
		}
		if ks[1] != keys[1] {
			t.Errorf("key[1] = %s actual = %s", keys[1], ks[1])
		}
		if ks[2] != keys[2] {
			t.Errorf("key[2] = %s actual = %s", keys[2], ks[2])
		}
	}
}
