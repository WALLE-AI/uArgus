package ratelimit

import (
	"math"
	"time"
)

// ExponentialBackoff returns the wait duration for attempt n (0-indexed).
// formula: base * 2^n, capped at maxWait.
func ExponentialBackoff(base time.Duration, n int, maxWait time.Duration) time.Duration {
	d := time.Duration(float64(base) * math.Pow(2, float64(n)))
	if d > maxWait {
		return maxWait
	}
	return d
}

// LinearBackoff returns base * (n+1), capped at maxWait.
func LinearBackoff(base time.Duration, n int, maxWait time.Duration) time.Duration {
	d := base * time.Duration(n+1)
	if d > maxWait {
		return maxWait
	}
	return d
}
