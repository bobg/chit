package chit

import "context"

// Map transforms a sequence of T-type elements into a sequence of U-type elements
// by applying a function to each one.
func Map[T, U any](ctx context.Context, inp *Iter[T], f func(T) (U, error)) *Iter[U] {
	return New(ctx, func(send func(U) error) error {
		for {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			y, err := f(x)
			if err != nil {
				return err
			}
			err = send(y)
			if err != nil {
				return err
			}
		}
	})
}
