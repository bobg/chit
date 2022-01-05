package chit

import "context"

// FromChan creates an iterator reading from a channel.
func FromChan[T any](ctx context.Context, inp <-chan T) *Iter[T] {
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

// ToChan copies the members of an input iterator to a new Go channel.
// If errptr is non-nil,
// any error encountered reading from the iterator or writing to the channel is placed there
// and can be read after the end of the channel is reached.
func ToChan[T any](ctx context.Context, iter *Iter[T], errptr *error) <-chan T {
	ch := make(chan T, 32)
	go func() {
		// Note: this must appear before the other defer,
		// so that this close executes later.
		// Otherwise there is a race condition
		// where the caller might try to read from *errptr before it is set.
		defer close(ch)

		var err error
		if errptr != nil {
			defer func() {
				*errptr = err
			}()
		}

		for {
			x, ok, err := iter.Next()
			if err != nil {
				return
			}
			if !ok {
				return
			}
			err = chsend(ctx, ch, x)
			if err != nil {
				return
			}
		}
	}()
	return ch
}
