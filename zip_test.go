package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestZip(t *testing.T) {
	ctx := context.Background()
	inp1 := FromSlice(ctx, []int{1, 2, 3})
	inp2 := FromSlice(ctx, []string{"a", "b", "c", "d"})
	z := Zip(ctx, inp1, inp2)
	got, err := ToSlice(ctx, z)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []Pair[int, string]{{X: 1, Y: "a"}, {X: 2, Y: "b"}, {X: 3, Y: "c"}, {Y: "d"}}) {
		t.Errorf("got %v, want [[1 a] [2 b] [3 c] [0 d]]", got)
	}
}
