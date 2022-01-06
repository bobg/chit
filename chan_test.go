package chit

import (
	"context"
	"testing"
)

func TestToChan(t *testing.T) {
	ctx := context.Background()
	ints := Ints(ctx, 0, 1)
	ch := ToChan(ctx, ints, nil)
	want := 0
	for got := range ch {
		if got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		want++
		if want > 10 {
			break
		}
	}
}
