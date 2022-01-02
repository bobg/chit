package chit

import (
	"context"
	"sync"
)

// Dup[T] duplicates the contents of an iterator,
// producing n new iterators,
// each containing the members of the original.
//
// An internal buffer grows to roughly the size
// of the difference between the output iterator that is farthest ahead in the stream,
// and the one that is farthest behind.
func Dup[T any](ctx context.Context, inp *Iter[T], n int) []*Iter[T] {
	var (
		mu        sync.Mutex
		buf       []T
		bufoffset int
		offsets   = make([]int, n)
		iters     []*Iter[T]
	)

	for idx := 0; idx < n; idx++ {
		idx := idx // Go loop-var pitfall
		var iter *Iter[T]
		iter = New(ctx, func(ctx context.Context, ch chan<- T) error {
			for {
				x, ok, err := func() (T, bool, error) {
					mu.Lock()
					defer mu.Unlock()

					if iter.Err != nil {
						var x T
						return x, false, iter.Err
					}

					for offsets[idx] >= bufoffset+len(buf) {
						x, ok, err := inp.Next()
						if err != nil {
							// xxx cancel other iters?
							// for j := 0; j < n; j++ {
							// 	if iters[j].Err == nil {
							// 		iters[j].Err = err
							// 	}
							// }
							return x, false, err
						}
						if !ok {
							var x T
							return x, false, nil
						}
						buf = append(buf, x)
					}

					x := buf[offsets[idx]-bufoffset]
					offsets[idx]++

					minoffset := offsets[0]
					for j := 1; j < n; j++ {
						if offsets[j] < minoffset {
							minoffset = offsets[j]
						}
					}
					if minoffset > bufoffset {
						buf = buf[minoffset-bufoffset:]
						bufoffset = minoffset
					}

					return x, true, nil
				}()
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}
				err = Send(ctx, ch, x)
				if err != nil {
					return err
				}
			}
		})

		iters = append(iters, iter)
	}

	return iters
}
