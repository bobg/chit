# Chit - Channel-based generic iterators

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/chit.svg)](https://pkg.go.dev/github.com/bobg/chit)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/chit)](https://goreportcard.com/report/github.com/bobg/chit)
[![Tests](https://github.com/bobg/chit/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/chit/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/chit/badge.svg?branch=master)](https://coveralls.io/github/bobg/chit?branch=master)

This is chit,
an implementation of channel-based generic iterators for Go 1.18 and later.

These iterators are implemented in terms of channels and goroutines,
details that are mostly invisible to the caller.
This is as opposed to
(what might at first seem to be)
the more obvious approach:
defining an abstract iterator interface,
then writing a bunch of concrete types implementing that interface.

The abstract-iterator approach is fine for things like a filter:

```go
type Iter[T any] interface {
	Next() (T, bool)
}

type FilterIter[T any] struct {
	inp Iter[T]
	f   func(T) bool
}

func (f *FilterIter[T]) Next() (T, bool) {
	for {
		x, ok := f.inp.Next()
		if !ok {
			return x, false
		}
		if f.f(x) {
			return x, true
		}
	}
}
```

but less-trivial iterators quickly run into concurrency issues,
code repetition,
and complexity arising from splitting state information between the concrete data structure and the `Next` function.

All of this is improved with a channel-and-goroutine-based approach,
at the cost of some performance.
Channels are concurrency-safe;
and with a single callback function
(running in a goroutine)
producing all of an iterator’s elements,
there is no need to store intermediate state in a data object.
