package chit

import (
	"context"
	"testing"
)

func BenchmarkHardcoded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sum int
		for j := 1; j <= 1000; j++ {
			if j^2 == 1 {
				sum += j * j
			}
		}
	}
}

func BenchmarkChit(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		var (
			ints      = Ints(ctx, 1, 1)
			first1000 = FirstN(ctx, ints, 1000)
			odds      = Filter(ctx, first1000, func(x int) (bool, error) { return x%2 == 1, nil })
			squares   = Map(ctx, odds, func(x int) (int, error) { return x * x, nil })
			sums      = Accum(ctx, squares, func(x, y int) (int, error) { return x + y, nil })
			sum       = LastN(ctx, sums, 1)
		)
		_, ok, err := sum.Read()
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal("no value in sum iterator")
		}
	}
}

type iterIF[T any] interface {
	Next() bool
	Val() T
}

// generates ints
type intIter struct {
	val, next, delta int
}

func (it *intIter) Next() bool {
	it.val = it.next
	it.next += it.delta
	return true
}

func (it *intIter) Val() int {
	return it.val
}

// produces the first N items
type firstNIter[T any] struct {
	inp iterIF[T]
	n   int
}

func (it *firstNIter[T]) Next() bool {
	if it.n <= 0 {
		return false
	}
	it.n--
	return it.inp.Next()
}

func (it *firstNIter[T]) Val() T {
	return it.inp.Val()
}

// filters according to a predicate
type filterIter[T any] struct {
	inp   iterIF[T]
	f     func(T) bool
	latch *T
}

func (it *filterIter[T]) Next() bool {
	for it.inp.Next() {
		val := it.inp.Val()
		if it.f(val) {
			it.latch = &val
			return true
		}
	}
	return false
}

func (it *filterIter[T]) Val() T {
	return *it.latch
}

// transforms
type mapIter[T, U any] struct {
	inp iterIF[T]
	f   func(T) U
}

func (it *mapIter[T, U]) Next() bool {
	return it.inp.Next()
}

func (it *mapIter[T, U]) Val() U {
	return it.f(it.inp.Val())
}

// accumulates
type accumIter[T any] struct {
	inp    iterIF[T]
	f      func(T, T) T
	latest *T
}

func (it *accumIter[T]) Next() bool {
	return it.inp.Next()
}

func (it *accumIter[T]) Val() T {
	var val T
	if it.latest == nil {
		val = it.inp.Val()
	} else {
		val = it.f(*it.latest, it.inp.Val())
	}
	it.latest = &val
	return val
}

// produces the last N items
type lastNIter[T any] struct {
	inp    iterIF[T]
	n      int
	filled bool
	buf    []T
	val    T
}

func (it *lastNIter[T]) Next() bool {
	if !it.filled {
		var (
			start int
			buf   []T
		)
		for it.inp.Next() {
			val := it.inp.Val()
			if len(buf) < it.n {
				buf = append(buf, val)
				continue
			}
			buf[start] = val
			start = (start + 1) % it.n
		}
		// straighten out the circular buffer
		it.buf = buf[start:]
		if start > 0 {
			it.buf = append(it.buf, buf[:start]...)
		}
	}
	if len(it.buf) == 0 {
		return false
	}
	it.val = it.buf[0]
	it.buf = it.buf[1:]
	return true
}

func (it *lastNIter[T]) Val() T {
	return it.val
}

func BenchmarkIF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			ints      = &intIter{next: 1, delta: 1}
			first1000 = &firstNIter[int]{inp: ints, n: 1000}
			odds      = &filterIter[int]{inp: first1000, f: func(x int) bool { return x%2 == 1 }}
			squares   = &mapIter[int, int]{inp: odds, f: func(x int) int { return x * x }}
			sums      = &accumIter[int]{inp: squares, f: func(x, y int) int { return x + y }}
			sum       = &lastNIter[int]{inp: sums, n: 1}
		)
		if !sum.Next() {
			b.Fatal("no value in sum iterator")
		}
	}
}
