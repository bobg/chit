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
