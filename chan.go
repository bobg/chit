package chit

// Chan creates an iterator reading from a channel.
func Chan[T any](ctx context.Context, inp <-chan T) *Iter[T] {
	return New(ctx, func(send func(T) error) error {
		for {
			select {
			case x, ok := <-inp:
				if !ok {
					return nil
				}
				err := send(x)
				if err != nil {
					return err
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}
