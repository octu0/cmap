package cmap

type slab struct {
	shards []Cache
	size   uint64
	hash   CMapHashFunc
}

func (s *slab) GetShard(key string) Cache {
	idx := int(s.hash.Hash64(key) % s.size)
	return s.shards[idx]
}

func (s *slab) Shards() []Cache {
	return s.shards
}

func newSlab(opt *cmapOption) *slab {
	shards := make([]Cache, opt.slabSize)
	size64 := uint64(opt.slabSize)
	for i := 0; i < opt.slabSize; i += 1 {
		shards[i] = newDefaultCache(opt.cacheCapacity)
	}
	return &slab{
		shards: shards,
		size:   size64,
		hash:   opt.hashFunc,
	}
}
