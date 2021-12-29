package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestInts(t *testing.T) {
	ctx := context.Background()
	ints := Ints(ctx, 1, 2)
	got, err := ToSlice(ctx, FirstN(ctx, ints, 10))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}) {
		t.Errorf("got %v, want [1 3 5 7 9 11 13 15 17 19]", got)
	}
}

func TestRepeat(t *testing.T) {
	ctx := context.Background()
	r := Repeat(ctx, "foo")
	got, err := ToSlice(ctx, FirstN(ctx, r, 10))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo"}) {
		t.Errorf("got %v, want 10 foos", got)
	}
}
