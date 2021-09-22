package cmap

import (
	"strconv"
	"testing"
)

func TestSlabShard(t *testing.T) {
	testRate := func(tt *testing.T, actual float64, expect float64, delta float64) {
		min := (expect - delta)
		max := (expect + delta)
		if min <= actual && actual <= max {
			tt.Logf("ok %3.5f <= actual(%3.5f) <= %3.5f", min, actual, max)
		} else {
			tt.Errorf("actual(%3.5f) out of range in expect(%3.5f -%3.5f +%3.5f)", actual, expect, min, max)
		}
	}

	t.Run("conflict/2", func(tt *testing.T) {
		opt := newDefaultOption()
		opt.slabSize = 2

		expectRate := 1.0 / 2.0

		slab := newSlab(opt)
		for i := 0; i < 1000; i += 1 {
			key := strconv.Itoa(i)
			slab.GetShard(key).Set(key, key)
		}
		for _, s := range slab.Shards() {
			r := float64(len(s.Keys())) / float64(1000)
			testRate(tt, r, expectRate, 0.05)
		}
	})
	t.Run("conflict/50", func(tt *testing.T) {
		opt := newDefaultOption()
		opt.slabSize = 50

		expectRate := 1.0 / 50.0

		slab := newSlab(opt)
		for i := 0; i < 1000; i += 1 {
			key := strconv.Itoa(i)
			slab.GetShard(key).Set(key, key)
		}
		for _, s := range slab.Shards() {
			r := float64(len(s.Keys())) / float64(1000)
			testRate(tt, r, expectRate, 0.01)
		}
	})
}
