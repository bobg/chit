package chit

import "context"

// FirstN produces an iterator containing the first n elements of the input
// (or all of the input, if there are fewer than n elements).
// Excess elements in the input are discarded by calling inp.Cancel.
func FirstN[T any](ctx context.Context, inp *Iter[T], n int) *Iter[T] {
	return New(ctx, func(send func(T) error) error {
		defer inp.Cancel()

		for i := 0; i < n; i++ {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			err = send(x)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// SkipN produces an iterator containing all but the first n elements of the input.
// If the input contains n or fewer elements,
// the output iterator will be empty.
func SkipN[T any](ctx context.Context, inp *Iter[T], n int) *Iter[T] {
	return New(ctx, func(send func(T) error) error {
		for i := 0; i < n; i++ {
			_, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
		}
		for {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			err = send(x)
			if err != nil {
				return err
			}
		}
	})
}

// LastN produces an iterator containing the last n elements of the input
// (or all of the input, if there are fewer than n elements).
// This requires buffering up to n elements.
// There is no guarantee that any elements will ever be produced:
// the input iterator may be infinite!
func LastN[T any](ctx context.Context, inp *Iter[T], n int) *Iter[T] {
	return New(ctx, func(send func(T) error) error {
		var (
			// Circular buffer.
			buf   = make([]T, 0, n)
			start = 0
		)
		for {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				for i := start; i < len(buf); i++ {
					err = send(buf[i])
					if err != nil {
						return err
					}
				}
				for i := 0; i < start; i++ {
					err = send(buf[i])
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
