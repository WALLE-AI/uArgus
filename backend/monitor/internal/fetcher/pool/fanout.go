package pool

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// FanOut executes multiple functions concurrently and returns the first error.
func FanOut(ctx context.Context, fns ...func(ctx context.Context) error) error {
	g, gCtx := errgroup.WithContext(ctx)
	for _, fn := range fns {
		fn := fn
		g.Go(func() error { return fn(gCtx) })
	}
	return g.Wait()
}
