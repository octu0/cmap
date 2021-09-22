# cmap

[![MIT License](https://img.shields.io/github/license/octu0/cmap)](https://github.com/octu0/cmap/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/cmap?status.svg)](https://godoc.org/github.com/octu0/cmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/cmap)](https://goreportcard.com/report/github.com/octu0/cmap)
[![Releases](https://img.shields.io/github/v/release/octu0/cmap)](https://github.com/octu0/cmap/releases)

`cmap` is inspired by [orcaman/concurrent-map](https://github.com/orcaman/concurrent-map), with performance improvements and some usable methods, while keeping same use cases.

## Installation

```
$ go get github.com/octu0/cmap
```

## Example

```go
import "github.com/octu0/cmap

var (
  m = cmap.New()
)

func main() {
  m.Set("foo", "bar")

  if v, ok := m.Get("foo"); ok {
    bar := v.(string)
  }

  m.Remove("foo")
}
```

## Benchmarks

```
goos: darwin
goarch: amd64
pkg: github.com/octu0/cmap
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkCompare/ConcurrentMap-8         	     496	   2806827 ns/op	  735840 B/op	    5153 allocs/op
BenchmarkCompare/sync.Map-8              	     252	   4720666 ns/op	  852128 B/op	   25158 allocs/op
BenchmarkCompare/cmap-8                  	     799	   1506810 ns/op	  468819 B/op	    5011 allocs/op
PASS
```

## License

MIT, see LICENSE file for details.
