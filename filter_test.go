package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	ctx := context.Background()
	ints := Ints(ctx, 1, 1)
	evens := Filter(ctx, ints, func(n int) (bool, error) { return n%2 == 0, nil })
	got, err := ToSlice(ctx, FirstN(ctx, evens, 3))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{2, 4, 6}) {
		t.Errorf("got %v, [2 4 6]", got)
	}
}

func TestSkipUntil(t *testing.T) {
	ctx := context.Background()
	ints := Ints(ctx, 1, 1)
	first10 := FirstN(ctx, ints, 10)
	latter := SkipUntil(ctx, first10, func(x int) (bool, error) { return x > 7, nil })
	got, err := ToSlice(ctx, latter)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{8, 9, 10}) {
		t.Errorf("got %v, want [8 9 10]", got)
	}
}
