package chit

import "context"

// Concat[T] takes a sequence of iterators and produces an iterator over all the elements of the input iterators, in sequence.
func Concat[T any](ctx context.Context, inps ...*Iter[T]) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		for _, inp := range inps {
			for {
				x, ok, err := inp.Read()
				if err != nil {
					return err
				}
				if !ok {
					break
				}
				err = chwrite(ctx, ch, x)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
