package fallback

import (
	"context"
	"fmt"
	"log/slog"
)

// Provider is a named data provider that may fail.
type Provider[T any] struct {
	Name string
	Fn   func(ctx context.Context) (T, error)
}

// Chain executes providers in order, returning the first success.
type Chain[T any] struct {
	providers []Provider[T]
}

// NewChain creates a fallback chain from providers.
func NewChain[T any](providers ...Provider[T]) *Chain[T] {
	return &Chain[T]{providers: providers}
}

// Execute tries each provider in order. Returns the first success.
func (c *Chain[T]) Execute(ctx context.Context) (T, error) {
	var lastErr error
	var zero T

	for _, p := range c.providers {
		val, err := p.Fn(ctx)
		if err == nil {
			return val, nil
		}
		slog.Warn("fallback: provider failed", "name", p.Name, "err", err)
		lastErr = err
	}

	return zero, fmt.Errorf("fallback: all %d providers failed, last: %w", len(c.providers), lastErr)
}
