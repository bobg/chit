package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestLast(t *testing.T) {
	ctx := context.Background()
	ints := Ints(ctx, 0, 1)
	first100 := FirstN(ctx, ints, 100)
	nineties := LastN(ctx, first100, 10)
	got, err := ToSlice(ctx, nineties)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{90, 91, 92, 93, 94, 95, 96, 97, 98, 99}) {
		t.Errorf("got %v, want [90 91 92 93 94 95 96 97 98 99]", got)
	}
}
