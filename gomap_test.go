package chit

import (
	"context"
	"reflect"
	"testing"
)

func TestGomaps(t *testing.T) {
	ctx := context.Background()
	inp := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	iter := FromMap(ctx, inp)
	got, err := ToMap(ctx, iter)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, inp) {
		t.Errorf("got %v, want [one:1 two:2 three:3]", got)
	}
}
