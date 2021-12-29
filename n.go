package chit

import "context"

func FirstN[T any](ctx context.Context, inp *Iter[T], n int) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		defer inp.Cancel()

		for i := 0; i < n; i++ {
			x, ok, err := inp.Read()
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

		return nil
	})
}

func LastN[T any](ctx context.Context, inp *Iter[T], n int) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		var (
			// Circular buffer.
			buf   = make([]T, 0, n)
			start = 0
		)
		for {
			x, ok, err := inp.Read()
			if err != nil {
				return err
			}
			if !ok {
				for i := start; i < len(buf); i++ {
					err = chwrite(ctx, ch, buf[i])
					if err != nil {
						return err
					}
				}
				for i := 0; i < start; i++ {
					err = chwrite(ctx, ch, buf[i])
					if err != nil {
						return err
					}
				}
				return nil
			}
			if len(buf) < n {
				buf = append(buf, x)
				continue
			}
			buf[start] = x
			start = (start + 1) % n
		}
	})
}
