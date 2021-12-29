package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestDup(t *testing.T) {
	ctx := context.Background()
	inp := FromSlice(ctx, []int{1, 2, 3})
	dups := Dup(ctx, inp, 2)
	s1, err := ToSlice(ctx, dups[0])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s1, []int{1, 2, 3}) {
		t.Errorf("got %v, want [1 2 3]", s1)
	}
	s2, err := ToSlice(ctx, dups[1])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s2, []int{1, 2, 3}) {
		t.Errorf("got %v, want [1 2 3]", s2)
	}
}
