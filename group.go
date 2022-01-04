package chit

func Group[T any, K comparable](ctx context.Context, inp *Iter[T], f func(T) (K, error)) *Iter[Pair[K, *Iter[T]]] {
	// When we discover a new partition (a new K value),
	// we create a channel to feed the corresponding Iter[T].
	m := make(map[K]chan<- T)

	return New(ctx, func(outerSend func(Pair[K, *Iter[T]]) error {
		defer func() {
			for _, ch := range m {
				close(ch)
			}
		}()

		for {
			x, ok, err := inp.Next()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			k, err := f(x)
			if err != nil {
				return err
			}

			if ch, ok := m[k]; ok {
				// This is an existing partition.
				// Supply the current value to its iterator.

				select {
				case ch <- x:
					// ok
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}

			// This is a new partition.
			// 

			ch := make(chan T, 32)
			m[k] = ch
			iter := Chan(ctx, ch)
			err = outerSend(Pair[K, *Iter[T]]{X: k, Y: iter})
			if err != nil {
				return err
			}

			select {
			case ch <- x:
				// ok
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}))
}
