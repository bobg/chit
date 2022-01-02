package chit

import "context"

// Filter filters the elements of an iterator according to a predicate function.
// Only the elements in the input iterator producing a true value appear in the output iterator.
func Filter[T any](ctx context.Context, inp *Iter[T], f func(T) (bool, error)) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		for {
			x, ok, err := inp.Next()
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
			err = Send(ctx, ch, x)
			if err != nil {
				return err
			}
		}
	})
}

// SkipUntil skips the inital elements of the input
// until calling the given predicate on an element returns true;
// then it copies that and the remaining elements to the output.
// The predicate is not called again after the first time it returns true.
func SkipUntil[T any](ctx context.Context, inp *Iter[T], f func(T) (bool, error)) *Iter[T] {
	skipping := true
	return Filter(ctx, inp, func(t T) (bool, error) {
		if !skipping {
			return true, nil
		}
		ok, err := f(t)
		if err != nil {
			return false, err
		}
		if ok {
			skipping = false
		}
		return !skipping, nil
	})
}
