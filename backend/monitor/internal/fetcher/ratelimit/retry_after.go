package ratelimit

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// HandleRetryAfter parses a 429 response's Retry-After header and waits.
// Returns true if it waited, false if no Retry-After was present.
func HandleRetryAfter(ctx context.Context, resp *http.Response) (bool, error) {
	if resp.StatusCode != http.StatusTooManyRequests {
		return false, nil
	}

	val := resp.Header.Get("Retry-After")
	if val == "" {
		return false, nil
	}

	// try as seconds
	if secs, err := strconv.Atoi(val); err == nil {
		return true, waitDuration(ctx, time.Duration(secs)*time.Second)
	}

	// try as HTTP-date
	if t, err := http.ParseTime(val); err == nil {
		d := time.Until(t)
		if d > 0 {
			return true, waitDuration(ctx, d)
		}
		return true, nil
	}

	return false, nil
}

func waitDuration(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
