// Chit defines functions for operating on channels as generic iterators.
package chit

import "context"

// Iter[T] is an iterator over items of type T.
// It contains an underlying channel of type <-chan T.
// Create an Iter[T] with New[T].
type Iter[T any] struct {
	// Err contains any error that might have closed the channel prematurely.
	// Callers should read it only after a call to Iter.Next returns a false "ok" value.
	Err error

	ch     <-chan T
	ctx    context.Context
	cancel context.CancelFunc
}

// New[T] creates a new Iter[T].
// The writer function is invoked once in a goroutine,
// and must supply all of the iterator's elements on the given channel.
// The writer function must not close the channel;
// this will happen automatically when the function exits.
func New[T any](ctx context.Context, writer func(context.Context, chan<- T) error) *Iter[T] {
	ctx, cancel := context.WithCancel(ctx)

	// Benchmark results:
	//   No channel buffering: 1.00 (baseline)
	//   Buffer size 1:        0.81
	//   Buffer size 32:       0.36
	//   Abstract iterators:   0.08 (vs. chit-style channels and goroutines)

	ch := make(chan T, 32)
	iter := &Iter[T]{
		ch:     ch,
		ctx:    ctx,
		cancel: cancel,
	}
	go func() {
		iter.Err = writer(ctx, ch)
		close(ch)
	}()
	return iter
}

// Next reads the next item from the iterator.
func (it *Iter[T]) Next() (T, bool, error) {
	select {
	case x, ok := <-it.ch:
		// xxx call it.cancel()
		return x, ok, nil
	case <-it.ctx.Done():
		var x T
		it.Err = it.ctx.Err()
		return x, false, it.Err
	}
}

// Cancel cancels the context in the iterator.
// This normally causes the iterator's "writer" function to terminate early,
// closing the iterator's underlying channel and causing Next calls to return context.Canceled.
func (it *Iter[T]) Cancel() {
	it.cancel()
}

// Send sends a value on a channel but returns early (with an error) if the context is canceled before the value can be sent.
func Send[T any](ctx context.Context, ch chan<- T, x T) error {
	select {
	case ch <- x:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
