package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestAccum(t *testing.T) {
	ctx := context.Background()
	inp := FromSlice(ctx, []int{1, 2, 3, 4})
	a := Accum(ctx, inp, func(a, b int) (int, error) { return a + b, nil })
	got, err := ToSlice(ctx, a)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{1, 3, 6, 10}) {
		t.Errorf("got %v, want [1 3 6 10]", got)
	}
}
