package registry

import (
	"sync"
	"time"
)

// HealthState describes the current health of a source.
type HealthState string

const (
	HealthStateHealthy  HealthState = "healthy"
	HealthStateDegraded HealthState = "degraded"
	HealthStateFailing  HealthState = "failing"
)

// HealthStatus holds the live health info for a single source.
type HealthStatus struct {
	LastSuccessAt    time.Time     `json:"lastSuccessAt"`
	LastAttemptAt    time.Time     `json:"lastAttemptAt"`
	ConsecutiveFails int           `json:"consecutiveFails"`
	LastRecordCount  int           `json:"lastRecordCount"`
	AvgDuration      time.Duration `json:"avgDurationMs"`
	State            HealthState   `json:"state"`

	// thresholds (from SourceSpec)
	MaxStaleDuration time.Duration `json:"maxStaleDuration"`
	MinRecordCount   int           `json:"minRecordCount"`
}

// HealthTracker records attempts/successes/failures for all registered sources.
type HealthTracker struct {
	mu      sync.RWMutex
	entries map[string]*HealthStatus
}

// NewHealthTracker creates a new tracker.
func NewHealthTracker() *HealthTracker {
	return &HealthTracker{entries: make(map[string]*HealthStatus)}
}

// Register initialises tracking for a source using its SourceSpec thresholds.
func (ht *HealthTracker) Register(s Source) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	spec := s.Spec()
	ht.entries[s.Name()] = &HealthStatus{
		State:            HealthStateHealthy,
		MaxStaleDuration: spec.MaxStaleDuration,
		MinRecordCount:   spec.MinRecordCount,
	}
}

// RecordAttempt marks the start of a run.
func (ht *HealthTracker) RecordAttempt(name string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	if e, ok := ht.entries[name]; ok {
		e.LastAttemptAt = time.Now()
	}
}

// RecordSuccess records a successful run.
func (ht *HealthTracker) RecordSuccess(name string, m FetchMetrics) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	e, ok := ht.entries[name]
	if !ok {
		return
	}
	now := time.Now()
	e.LastSuccessAt = now
	e.ConsecutiveFails = 0
	e.LastRecordCount = m.RecordCount
	// exponential moving average for duration
	if e.AvgDuration == 0 {
		e.AvgDuration = m.Duration
	} else {
		e.AvgDuration = (e.AvgDuration*7 + m.Duration*3) / 10
	}
	e.State = ht.deriveState(e, now)
}

// RecordFailure records a failed run.
func (ht *HealthTracker) RecordFailure(name string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	e, ok := ht.entries[name]
	if !ok {
		return
	}
	e.ConsecutiveFails++
	e.State = ht.deriveState(e, time.Now())
}

// Snapshot returns a copy of all health statuses.
func (ht *HealthTracker) Snapshot() map[string]HealthStatus {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	out := make(map[string]HealthStatus, len(ht.entries))
	for k, v := range ht.entries {
		out[k] = *v
	}
	return out
}

func (ht *HealthTracker) deriveState(e *HealthStatus, now time.Time) HealthState {
	if e.ConsecutiveFails >= 3 {
		return HealthStateFailing
	}
	if e.MaxStaleDuration > 0 && !e.LastSuccessAt.IsZero() &&
		now.Sub(e.LastSuccessAt) > e.MaxStaleDuration {
		return HealthStateDegraded
	}
	if e.MinRecordCount > 0 && e.LastRecordCount > 0 &&
		e.LastRecordCount < e.MinRecordCount {
		return HealthStateDegraded
	}
	return HealthStateHealthy
}
