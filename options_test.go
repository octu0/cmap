package cmap

import (
	"testing"
)

func TestDefaultOption(t *testing.T) {
	d := newDefaultOption()

	if d.slabSize != defaultSlabSize {
		t.Errorf("default slab size = %d", defaultSlabSize)
	}
	if d.cacheCapacity != defaultCacheCapacity {
		t.Errorf("default cache capacity size = %d", defaultCacheCapacity)
	}
	if d.hashFunc == nil {
		t.Errorf("default hash func not nil")
	}
}
