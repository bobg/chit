package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	ctx := context.Background()
	inp := FromSlice(ctx, []int{1, 2, 3, 4})
	m := Map(ctx, inp, func(x int) (int, error) { return x * x, nil })
	s, err := ToSlice(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, []int{1, 4, 9, 16}) {
		t.Errorf("got %v, want [1, 4, 9, 16]", s)
	}
}
