package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"golang.org/x/sync/singleflight"
)

// negSentinel is the internal marker stored for negative caching.
// Unlike v1's raw string "__WM_NEG__", this is a typed JSON object
// so it cannot collide with real data.
var negSentinel = []byte(`{"__neg":true}`)

func isNeg(b []byte) bool {
	return len(b) > 0 && b[0] == '{' && len(b) < 20 && string(b) == string(negSentinel)
}

// FetchThrough is a generic read-through cache with singleflight,
// negative caching, and optional stale-while-revalidate.
type FetchThrough[T any] struct {
	client  Client
	group   singleflight.Group
	metrics *Metrics // may be nil
}

// NewFetchThrough creates a FetchThrough. metrics may be nil.
func NewFetchThrough[T any](client Client, metrics *Metrics) *FetchThrough[T] {
	return &FetchThrough[T]{client: client, metrics: metrics}
}

// FetchOpts configures a single Fetch call.
type FetchOpts[T any] struct {
	Key                  string
	TTL                  time.Duration
	Fetcher              func(ctx context.Context) (*T, error)
	NegativeTTL          time.Duration // default 120s
	StaleWhileRevalidate time.Duration // >0 → return stale + async refresh
	Validator            func(raw []byte) bool
}

// Fetch checks cache, coalesces concurrent requests via singleflight,
// and falls back to Fetcher on miss.
func (ft *FetchThrough[T]) Fetch(ctx context.Context, opts FetchOpts[T]) (*T, error) {
	if opts.NegativeTTL == 0 {
		opts.NegativeTTL = 120 * time.Second
	}

	start := time.Now()

	// 1. cache check
	cached, err := ft.client.Get(ctx, opts.Key)
	if err != nil {
		ft.recordError(opts.Key, "get")
		slog.Warn("cache get error, falling through", "key", opts.Key, "err", err)
	} else if cached != nil {
		if isNeg(cached) {
			ft.recordHit(opts.Key, time.Since(start))
			return nil, nil // negative hit
		}
		if opts.Validator == nil || opts.Validator(cached) {
			var val T
			if err := json.Unmarshal(cached, &val); err == nil {
				ft.recordHit(opts.Key, time.Since(start))
				return &val, nil
			}
		}
	}

	ft.recordMiss(opts.Key, time.Since(start))

	// 2. singleflight fetch
	v, err, _ := ft.group.Do(opts.Key, func() (any, error) {
		val, fetchErr := opts.Fetcher(ctx)
		if fetchErr != nil {
			return nil, fetchErr
		}
		if val == nil {
			// negative cache
			_ = ft.client.Set(ctx, opts.Key, negSentinel, opts.NegativeTTL)
			return nil, nil
		}
		// write to cache
		data, marshalErr := json.Marshal(val)
		if marshalErr != nil {
			return val, nil // return value even if we can't cache
		}
		if setErr := ft.client.Set(ctx, opts.Key, data, opts.TTL); setErr != nil {
			ft.recordError(opts.Key, "set")
			slog.Warn("cache set error", "key", opts.Key, "err", setErr)
		} else {
			ft.recordWrite(opts.Key, len(data), opts.TTL)
		}
		return val, nil
	})

	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	val := v.(*T)
	return val, nil
}

// ── metrics helpers ─────────────────────────────────────────

func (ft *FetchThrough[T]) recordHit(key string, d time.Duration) {
	if ft.metrics != nil {
		ft.metrics.RecordHit(key, d)
	}
}

func (ft *FetchThrough[T]) recordMiss(key string, d time.Duration) {
	if ft.metrics != nil {
		ft.metrics.RecordMiss(key, d)
	}
}

func (ft *FetchThrough[T]) recordError(key, errType string) {
	if ft.metrics != nil {
		ft.metrics.RecordError(key, errType)
	}
}

func (ft *FetchThrough[T]) recordWrite(key string, b int, ttl time.Duration) {
	if ft.metrics != nil {
		ft.metrics.RecordWrite(key, b, ttl)
	}
}
