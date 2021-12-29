package chit

import "context"

// Group partitions the elements of an iterator into multiple separate iterator streams
// based on a given partitioning function.
// Each item in the input is fed to the function to see which partition it belongs in.
// The output is an iterator of X,Y pairs where X is the partition key and Y is an iterator over the elements in that partition.
//
// Callers reading the top-level output iterator should launch goroutines to consume the nested iterators.
// Otherwise the process trying to consume items in partition P1 is likely to block
// while the Group iterator waits for something to consume an item in partition P2.
//
// Example:
//
//   var (
//     groups = Group(ctx, input, partitionFunc)
//     wg     sync.WaitGroup
//   )
//   for {
//     pair, ok, err := groups.Read()
//     // ...check err...
//     if !ok {
//       break
//     }
//     wg.Add(1)
//     go func() {
//       defer wg.Done()
//       partitionKey, partitionItems := pair.X, pair.Y
//       for {
//         item, ok, err := partitionItems.Read()
//         // ...check err...
//         if !ok {
//           break
//         }
//         // ...handle item...
//       }
//     }()
//   }
//   wg.Wait()
func Group[T any, U comparable](ctx context.Context, inp *Iter[T], partition func(T) (U, error)) *Iter[Pair[U, *Iter[T]]] {
	m := make(map[U]chan<- T)

	return New(ctx, func(ctx context.Context, ch chan<- Pair[U, *Iter[T]]) error {
		defer func() {
			for _, c := range m {
				close(c)
			}
		}()

		for {
			x, ok, err := inp.Read()
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			p, err := partition(x)
			if err != nil {
				return err
			}

			c, ok := m[p]
			if !ok {
				pipe := make(chan T)
				m[p] = pipe

				iter := New(ctx, func(ctx context.Context, c chan<- T) error {
					for {
						select {
						case x, ok := <-pipe:
							if !ok {
								return nil
							}
							err = chwrite(ctx, c, x)
							if err != nil {
								return err
							}

						case <-ctx.Done():
							return ctx.Err()
						}
					}
				})
				err = chwrite(ctx, ch, Pair[U, *Iter[T]]{X: p, Y: iter})
				if err != nil {
					return err
				}
				c = pipe
			}

			err = chwrite(ctx, c, x)
			if err != nil {
				return err
			}
		}
	})
}
