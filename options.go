package cmap

const (
	defaultSlabSize      int = 1024
	defaultCacheCapacity int = 4096
)

type cmapOptionFunc func(*cmapOption)

type cmapOption struct {
	slabSize      int
	cacheCapacity int
	hashFunc      CMapHashFunc
}

func newDefaultOption() *cmapOption {
	return &cmapOption{
		slabSize:      defaultSlabSize,
		cacheCapacity: defaultCacheCapacity,
		hashFunc:      NewXXHashFunc(),
	}
}

func WithSlabSize(size int) cmapOptionFunc {
	return func(opt *cmapOption) {
		opt.slabSize = size
	}
}

func WithCacheCapacity(size int) cmapOptionFunc {
	return func(opt *cmapOption) {
		opt.cacheCapacity = size
	}
}

func WithHashFunc(hashFunc CMapHashFunc) cmapOptionFunc {
	return func(opt *cmapOption) {
		opt.hashFunc = hashFunc
	}
}
