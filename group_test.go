package chit

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestGroup(t *testing.T) {
	ctx := context.Background()
	inp := FirstN(ctx, Ints(ctx, 1, 1), 10)
	groups := Group(ctx, inp, func(x int) (int, error) { return x % 3, nil })
	m := map[int][]int{
		0: nil,
		1: nil,
		2: nil,
	}
	var wg sync.WaitGroup
	for {
		pair, ok, err := groups.Next()
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			break
		}
		wg.Add(1)
		go func() {
			s, err := ToSlice(ctx, pair.Y)
			if err != nil {
				panic(err)
			}
			m[pair.X] = s
			wg.Done()
		}()
	}
	wg.Wait()

	if !reflect.DeepEqual(m[0], []int{3, 6, 9}) {
		t.Errorf("got %v, want [3 6 9]", m[0])
	}
	if !reflect.DeepEqual(m[1], []int{1, 4, 7, 10}) {
		t.Errorf("got %v, want [1 4 7 10]", m[1])
	}
	if !reflect.DeepEqual(m[2], []int{2, 5, 8}) {
		t.Errorf("got %v, want [2 5 8]", m[2])
	}
}
