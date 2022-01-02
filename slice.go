package chit

import "context"

// FromSlice creates a channel iterator over a slice.
func FromSlice[T any](ctx context.Context, inp []T) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		for _, x := range inp {
			err := Send(ctx, ch, x)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ToSlice consumes all of an iterator's elements and returns them as a slice.
// Be sure your input isn't infinite, or very large!
// (Consider using FirstN to ensure the input has a reasonable size.)
func ToSlice[T any](ctx context.Context, inp *Iter[T]) ([]T, error) {
	var result []T
	for {
		x, ok, err := inp.Next()
		if err != nil {
			return nil, err
		}
		if !ok {
			return result, nil
		}
		result = append(result, x)
	}
}
