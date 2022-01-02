package chit

import "context"

// Accum accumulates the result of repeatedly applying a function to the elements of an iterator.
// If inp[i] is the ith element of the input
// and out[i] is the ith element of the output,
// then:
//   out[0] == inp[0]
// and
//   out[i+1] == f(out[i], inp[i+1])
func Accum[T any](ctx context.Context, inp *Iter[T], f func(T, T) (T, error)) *Iter[T] {
	return New(ctx, func(send func(T) error) error {
		var (
			last  T
			first = true
		)
		for {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			if first {
				last = x
				first = false
			} else {
				last, err = f(last, x)
				if err != nil {
					return err
				}
			}
			err = send(last)
			if err != nil {
				return err
			}
		}
	})
}
