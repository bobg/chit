package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	ctx := context.Background()
	c := Concat(
		ctx,
		FromSlice(ctx, []int{1, 2, 3}),
		FromSlice(ctx, []int{4, 5, 6}),
	)
	got, err := ToSlice(ctx, c)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("got %v, want [1 2 3 4 5 6]", got)
	}
}
