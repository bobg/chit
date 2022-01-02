package chit

import (
	"context"
	"testing"
)

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
