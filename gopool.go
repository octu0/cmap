package cmap

import (
	"sync/atomic"
)

type taskFunc func()

type gopool struct {
	ch     chan taskFunc
	closed int32
}

func (p *gopool) Submit(task taskFunc) bool {
	if atomic.LoadInt32(&p.closed) == 1 {
		return false
	}

	p.ch <- task
	return true
}

func (p *gopool) Close() {
	if atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		close(p.ch)
	}
}

func runGoPool(ch chan taskFunc) {
	for {
		select {
		case task, ok := <-ch:
			if ok != true {
				return
			}
			task()
		}
	}
}

func newGopool(size int) *gopool {
	ch := make(chan taskFunc, size)
	for i := 0; i < size; i += 1 {
		go runGoPool(ch)
	}
	return &gopool{ch, 0}
}
