package provider

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Retry Logic
// ---------------------------------------------------------------------------

// RetryConfig controls retry behavior for API calls.
type RetryConfig struct {
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
	JitterFraction    float64
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
		JitterFraction:    0.1,
	}
}

// RetryClassification describes how an error should be handled.
type RetryClassification int

const (
	RetryClassNoRetry RetryClassification = iota
	RetryClassRetry
	RetryClassRateLimit
	RetryClassOverloaded
	RetryClassPromptTooLong
)

// ClassifyError determines the retry classification for an error.
func ClassifyError(errMsg string) RetryClassification {
	lower := strings.ToLower(errMsg)

	switch {
	case strings.Contains(lower, "prompt is too long") ||
		strings.Contains(lower, "prompt_too_long") ||
		strings.Contains(lower, "context_length_exceeded"):
		return RetryClassPromptTooLong

	case strings.Contains(lower, "overloaded") || strings.Contains(lower, "529"):
		return RetryClassOverloaded

	case strings.Contains(lower, "rate_limit") ||
		strings.Contains(lower, "rate limit") ||
		strings.Contains(lower, "429"):
		return RetryClassRateLimit

	case strings.Contains(lower, "503") ||
		strings.Contains(lower, "500") ||
		strings.Contains(lower, "internal server error") ||
		strings.Contains(lower, "bad gateway") ||
		strings.Contains(lower, "service unavailable"):
		return RetryClassRetry

	default:
		return RetryClassNoRetry
	}
}

// ParseRetryAfter extracts a Retry-After duration from an error message or header value.
func ParseRetryAfter(s string) time.Duration {
	if secs, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := time.Parse(time.RFC1123, strings.TrimSpace(s)); err == nil {
		d := time.Until(t)
		if d > 0 {
			return d
		}
	}
	return 0
}

// ComputeBackoff calculates the backoff duration for a given attempt.
func ComputeBackoff(config RetryConfig, attempt int) time.Duration {
	if attempt <= 0 {
		return config.InitialBackoff
	}
	backoff := float64(config.InitialBackoff) * math.Pow(config.BackoffMultiplier, float64(attempt))
	if backoff > float64(config.MaxBackoff) {
		backoff = float64(config.MaxBackoff)
	}
	if config.JitterFraction > 0 {
		jitter := backoff * config.JitterFraction * (rand.Float64()*2 - 1)
		backoff += jitter
	}
	if backoff < 0 {
		backoff = float64(config.InitialBackoff)
	}
	return time.Duration(backoff)
}

// WaitForRetry blocks until the backoff period elapses or context is cancelled.
func WaitForRetry(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// RetryableError wraps an error with retry classification metadata.
type RetryableError struct {
	Err            error
	Classification RetryClassification
	RetryAfter     time.Duration
	Attempt        int
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error (class=%d, attempt=%d): %v", e.Classification, e.Attempt, e.Err)
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// ShouldRetry returns true if the error classification allows retry.
func (e *RetryableError) ShouldRetry() bool {
	switch e.Classification {
	case RetryClassRetry, RetryClassRateLimit, RetryClassOverloaded:
		return true
	default:
		return false
	}
}
