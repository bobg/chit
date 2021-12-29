package chit

import "context"

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

func Ints(ctx context.Context, start, delta int) *Iter[int] {
	n := start
	return Gen(ctx, func() (int, bool, error) {
		res := n
		n += delta
		return res, true, nil
	})
}

func Repeat[T any](ctx context.Context, val T) *Iter[T] {
	return Gen(ctx, func() (T, bool, error) {
		return val, true, nil
	})
}
