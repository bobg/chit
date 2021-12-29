package chit

import "context"

// Filter filters the elements of an iterator according to a predicate function.
// Only the elements in the input iterator producing a true value appear in the output iterator.
func Filter[T any](ctx context.Context, inp *Iter[T], f func(T) (bool, error)) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		for {
			x, ok, err := inp.Read()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			ok, err = f(x)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			err = chwrite(ctx, ch, x)
			if err != nil {
				return err
			}
		}
	})
}
