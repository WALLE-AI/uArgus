package registry

import (
	"context"
	"time"
)

// Source is the interface every data source must implement.
type Source interface {
	// Name returns a unique identifier, e.g. "news:digest:full:en".
	Name() string

	// Spec returns the declarative metadata for this source.
	Spec() SourceSpec

	// Dependencies returns Names of other Sources that must run first.
	Dependencies() []string

	// Run executes the data fetch/transform/publish cycle.
	Run(ctx context.Context) (*FetchResult, error)
}

// FetchResult is the output of Source.Run().
type FetchResult struct {
	Data        any
	ExtraKeys   []ExtraKey
	Metrics     FetchMetrics
	NeedsEnrich bool
}

// ExtraKey describes an additional Redis key to write alongside the canonical key.
type ExtraKey struct {
	Key  string
	Data any
	TTL  time.Duration
}

// FetchMetrics captures telemetry from a single Run invocation.
type FetchMetrics struct {
	Duration       time.Duration
	RecordCount    int
	UpstreamStatus int
	BytesRead      int64
}
