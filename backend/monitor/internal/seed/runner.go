package seed

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// Runner orchestrates the full seed lifecycle:
// lock → source.Run() → publishWithEnvelope → writeMeta → verify → unlock.
// On failure: extendTTL to preserve existing data (v1 resilience model).
type Runner struct {
	client  cache.Client
	metrics *cache.Metrics
	lock    *Lock
	pub     *Publisher
	meta    *MetaWriter
	ttl     *TTLManager
}

// NewRunner creates a Runner with all sub-components.
func NewRunner(client cache.Client, metrics *cache.Metrics) *Runner {
	return &Runner{
		client:  client,
		metrics: metrics,
		lock:    NewLock(client),
		pub:     NewPublisher(client),
		meta:    NewMetaWriter(client),
		ttl:     NewTTLManager(client),
	}
}

// Run executes the full seed cycle for a Source.
func (r *Runner) Run(ctx context.Context, src registry.Source) error {
	spec := src.Spec()
	name := src.Name()
	runID := fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
	dataTTL := ResolveTTL(spec)
	metaTTL := ResolveMetaTTL(spec)

	logger := slog.With("source", name, "runID", runID)

	// 1. acquire lock
	lockTTL := spec.LockTTL
	if lockTTL == 0 {
		lockTTL = 5 * time.Minute
	}
	acquired, err := r.lock.Acquire(ctx, spec.Domain, spec.Resource, runID, lockTTL)
	if err != nil {
		logger.Warn("lock acquire error", "err", err)
		return fmt.Errorf("seed run %s: lock: %w", name, err)
	}
	if !acquired {
		logger.Info("lock not acquired, skipping")
		return nil
	}

	// ensure lock release
	defer func() {
		if relErr := r.lock.Release(ctx, spec.Domain, spec.Resource, runID); relErr != nil {
			logger.Warn("lock release error", "err", relErr)
		}
	}()

	// 2. execute source
	start := time.Now()
	result, err := src.Run(ctx)
	duration := time.Since(start)

	if err != nil {
		logger.Error("source run failed, extending existing TTL", "err", err, "duration", duration)
		// resilience: extend existing data TTL
		keys := append([]string{spec.CanonicalKey}, spec.ExtraKeys...)
		if extErr := r.ttl.ExtendExisting(ctx, keys, dataTTL); extErr != nil {
			logger.Warn("extend TTL failed", "err", extErr)
		}
		return fmt.Errorf("seed run %s: %w", name, err)
	}

	// 3. publish with envelope
	meta := SeedMeta{
		FetchedAt:     time.Now().UnixMilli(),
		RecordCount:   result.Metrics.RecordCount,
		SourceVersion: "v2",
		SchemaVersion: 1,
		State:         "OK",
	}
	if result.Metrics.RecordCount == 0 {
		meta.State = "OK_ZERO"
	}

	if pubErr := r.pub.PublishWithEnvelope(ctx, spec.CanonicalKey, meta, result.Data, dataTTL); pubErr != nil {
		logger.Error("publish failed, extending existing TTL", "err", pubErr)
		keys := append([]string{spec.CanonicalKey}, spec.ExtraKeys...)
		_ = r.ttl.ExtendExisting(ctx, keys, dataTTL)
		return fmt.Errorf("seed run %s: publish: %w", name, pubErr)
	}

	// 4. write extra keys
	for _, ek := range result.ExtraKeys {
		ekTTL := ek.TTL
		if ekTTL == 0 {
			ekTTL = dataTTL
		}
		ekMeta := meta
		if ekErr := r.pub.PublishWithEnvelope(ctx, ek.Key, ekMeta, ek.Data, ekTTL); ekErr != nil {
			logger.Warn("extra key publish failed", "key", ek.Key, "err", ekErr)
		}
	}

	// 5. write freshness metadata
	if metaErr := r.meta.WriteFreshness(ctx, spec.Domain, spec.Resource, result.Metrics.RecordCount, name, metaTTL); metaErr != nil {
		logger.Warn("write freshness failed", "err", metaErr)
	}

	logger.Info("seed run complete",
		"records", result.Metrics.RecordCount,
		"duration", duration,
		"bytes", result.Metrics.BytesRead)

	return nil
}
