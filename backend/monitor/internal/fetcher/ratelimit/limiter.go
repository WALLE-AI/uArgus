package ratelimit

import (
	"context"
	"time"
)

// RateLimiter controls the pace of outgoing requests.
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// fixedInterval enforces a minimum interval between requests.
type fixedInterval struct {
	interval time.Duration
	last     time.Time
}

// NewFixedInterval creates a RateLimiter with a fixed delay between requests.
func NewFixedInterval(d time.Duration) RateLimiter {
	return &fixedInterval{interval: d}
}

func (f *fixedInterval) Wait(ctx context.Context) error {
	if !f.last.IsZero() {
		elapsed := time.Since(f.last)
		if elapsed < f.interval {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(f.interval - elapsed):
			}
		}
	}
	f.last = time.Now()
	return nil
}
