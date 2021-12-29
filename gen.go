package chit

import "context"

// Gen produces an iterator whose members are generated by successive calls to a given function.
func Gen[T any](ctx context.Context, f func() (T, bool, error)) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		for {
			x, ok, err := f()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			err = chwrite(ctx, ch, x)
			if err != nil {
				return err
			}
		}
	})
}

// Ints produces an infinite iterator of integers,
// starting and start and incrementing by delta.
func Ints(ctx context.Context, start, delta int) *Iter[int] {
	n := start
	return Gen(ctx, func() (int, bool, error) {
		res := n
		n += delta
		return res, true, nil
	})
}

// Repeat produces an infinite iterator of the given element, over and over.
// Useful in combination with FirstN when you want a certain number of the same item.
func Repeat[T any](ctx context.Context, val T) *Iter[T] {
	return Gen(ctx, func() (T, bool, error) {
		return val, true, nil
	})
}
