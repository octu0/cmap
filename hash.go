package cmap

import (
	"hash/fnv"
	"reflect"
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

type CMapHashFunc interface {
	Hash64(string) uint64
}

func NewXXHashFunc() CMapHashFunc {
	return &xxHashFunc{}
}

func NewFNV64HashFun() CMapHashFunc {
	return &fnv64HashFunc{}
}

type xxHashFunc struct{}

func (*xxHashFunc) Hash64(key string) uint64 {
	return xxhash.Sum64String(key)
}

type fnv64HashFunc struct{}

func (*fnv64HashFunc) Hash64(key string) uint64 {
	b := *(*[]byte)(unsafe.Pointer(&key))
	s := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	s.Len = len(key)
	s.Cap = len(key)

	f := fnv.New64()
	f.Write(b)
	return f.Sum64()
}
