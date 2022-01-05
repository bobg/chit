package chit

import "context"

// Group partitions the elements of an input iterator into separate groups according to a given partitioning function.
// The output is an iterator of (K,iterator) pairs,
// where each K is a value produced by the partitioning function,
// and each iterator contains the elements from the input belonging to that group.
//
// The logic that consumes the input iterator is single-threaded.
// This means that a caller consuming the top-level output of Group
// (i.e., the iterator of pairs)
// should launch goroutines to consume the nested iterator in each pair.
// Otherwise a reader waiting for a value on the K1 sub-iterator
// may deadlock waiting for someone to read the value
// that Group is trying to supply to the K2 sub-iterator.
//
// Illustration:
//
//   outer := Group(ctx, inp, partitionFunc)
//   for {
//     pair, ok, err := outer.Next()
//     if err != nil { ... }
//     if !ok { break }
//     k, inner := pair.X, pair.Y
//     go func() {
//       for {
//         x, ok, err := inner.Next()
//         if err != nil { ... }
//         if !ok { break }
//         ...handle value x in the k group...
//       }
//     }()
//   }
func Group[T any, K comparable](ctx context.Context, inp *Iter[T], f func(T) (K, error)) *Iter[Pair[K, *Iter[T]]] {
	// When we discover a new partition (a new K value),
	// we create a channel to feed the corresponding Iter[T].
	m := make(map[K]chan<- T)

	return New(ctx, func(outerSend func(Pair[K, *Iter[T]]) error) error {
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

				err = chsend(ctx, ch, x)
				if err != nil {
					return err
				}
				continue
			}

			// This is a new partition.

			ch := make(chan T, 32)
			m[k] = ch
			iter := Chan(ctx, ch)
			err = outerSend(Pair[K, *Iter[T]]{X: k, Y: iter})
			if err != nil {
				return err
			}

			err = chsend(ctx, ch, x)
			if err != nil {
				return err
			}
		}
	})
}
