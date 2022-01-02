package chit

import "context"

// FromMap creates a channel iterator over a map.
func FromMap[K comparable, V any](ctx context.Context, inp map[K]V) *Iter[Pair[K, V]] {
	return New(ctx, func(ctx context.Context, ch chan<- Pair[K, V]) error {
		for k, v := range inp {
			k, v := k, v // Go loop-var pitfall
			err := chwrite(ctx, ch, Pair[K, V]{X: k, Y: v})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ToMap consumes all of an iterator's elements and returns them as a map of pairs.
func ToMap[K comparable, V any](ctx context.Context, inp *Iter[Pair[K, V]]) (map[K]V, error) {
	result := make(map[K]V)
	for {
		pair, ok, err := inp.Read()
		if err != nil {
			return nil, err
		}
		if !ok {
			return result, nil
		}
		result[pair.X] = pair.Y
	}
}
