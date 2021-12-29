package chit

import "context"

// Zip takes two iterators and produces a new iterator containing pairs of corresponding elements.
// If one input iterator ends before the other,
// Zip produces zero values of the appropriate type in constructing pairs.
func Zip[T, U any](ctx context.Context, t *Iter[T], u *Iter[U]) *Iter[Pair[T, U]] {
	return New(ctx, func(ctx context.Context, ch chan<- Pair[T, U]) error {
		okx, oky := true, true

		for {
			var (
				x T
				y U
			)

			if okx {
				xx, ok, err := t.Read()
				if err != nil {
					return err
				}
				if ok {
					x = xx
				}
				okx = ok
			}
			if oky {
				yy, ok, err := u.Read()
				if err != nil {
					return err
				}
				if ok {
					y = yy
				}
				oky = ok
			}

			if !okx && !oky {
				return nil
			}

			err := chwrite(ctx, ch, Pair[T, U]{X: x, Y: y})
			if err != nil {
				return err
			}
		}
	})
}
