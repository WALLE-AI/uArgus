package registry

import "time"

// SourceSpec is the declarative descriptor for a data source.
// Each Source fills one at construction time; scheduling, health thresholds,
// TTL and canonical keys are all derived from this single struct.
type SourceSpec struct {
	// ── identity ────────────────────────────────────────────
	Domain   string // "news", "cyber", "resilience", "research" …
	Resource string // "digest", "threats", "scores", "arxiv" …
	Version  int    // 1, 2 …

	// ── scheduling ──────────────────────────────────────────
	Schedule Schedule      // Cron | Interval | OnDemand
	LockTTL  time.Duration // distributed lock TTL for seed runs

	// ── data keys ───────────────────────────────────────────
	CanonicalKey string   // e.g. "news:digest:v1:full:en"
	ExtraKeys    []string // bootstrap mirror keys
	DataTTL      time.Duration

	// ── health ──────────────────────────────────────────────
	MaxStaleDuration time.Duration // degrade after this
	MinRecordCount   int           // min expected records
}

// ── Schedule types ──────────────────────────────────────────

// Schedule is a sealed interface for the three scheduling modes.
type Schedule interface{ scheduleTag() }

// CronSchedule wraps a cron expression (Railway Cron / Bundle style).
type CronSchedule struct{ Expr string }

func (CronSchedule) scheduleTag() {}

// IntervalSchedule triggers at a fixed interval (relay setInterval style).
type IntervalSchedule struct{ Every time.Duration }

func (IntervalSchedule) scheduleTag() {}

// OnDemandSchedule is triggered externally (RPC warm-ping style).
type OnDemandSchedule struct{}

func (OnDemandSchedule) scheduleTag() {}
