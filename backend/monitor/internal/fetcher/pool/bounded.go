package pool

import (
	"context"
	"sync"
)

// Result holds the output or error of a single pooled task.
type Result[T any] struct {
	Value T
	Err   error
	Index int
}

// BoundedPool executes fn on each item with at most n concurrent goroutines.
// All items are processed (allSettled semantics — failures don't cancel others).
func BoundedPool[T any, R any](ctx context.Context, n int, items []T, fn func(ctx context.Context, item T) (R, error)) []Result[R] {
	if n <= 0 {
		n = 1
	}
	results := make([]Result[R], len(items))
	sem := make(chan struct{}, n)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it T) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			val, err := fn(ctx, it)
			results[idx] = Result[R]{Value: val, Err: err, Index: idx}
		}(i, item)
	}

	wg.Wait()
	return results
}
