package tracking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

const (
	intervalKeyPrefix = "score-interval:"
	intervalTTL       = 604800 * time.Second // 7d
)

// ScoreInterval holds a time-series of scores for trend detection.
type ScoreInterval struct {
	CountryCode string    `json:"countryCode"`
	Scores      []float64 `json:"scores"`
	Trend       string    `json:"trend"` // "stable" | "improving" | "declining"
}

// IntervalTracker tracks resilience scores over time.
type IntervalTracker struct {
	rdb cache.Client
}

// NewIntervalTracker creates an IntervalTracker.
func NewIntervalTracker(rdb cache.Client) *IntervalTracker {
	return &IntervalTracker{rdb: rdb}
}

// Read retrieves the score interval for a country.
func (t *IntervalTracker) Read(ctx context.Context, countryCode string) (*ScoreInterval, error) {
	key := intervalKeyPrefix + countryCode
	data, err := t.rdb.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var si ScoreInterval
	if err := json.Unmarshal(data, &si); err != nil {
		return nil, fmt.Errorf("interval tracker read: %w", err)
	}
	return &si, nil
}

// Write appends a score and recomputes the trend.
func (t *IntervalTracker) Write(ctx context.Context, countryCode string, score float64) error {
	key := intervalKeyPrefix + countryCode
	si := &ScoreInterval{CountryCode: countryCode}

	data, err := t.rdb.Get(ctx, key)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, si)
	}

	si.Scores = append(si.Scores, score)
	// keep last 30 data points
	if len(si.Scores) > 30 {
		si.Scores = si.Scores[len(si.Scores)-30:]
	}
	si.Trend = deriveTrend(si.Scores)

	out, err := json.Marshal(si)
	if err != nil {
		return fmt.Errorf("interval tracker write: %w", err)
	}
	return t.rdb.Set(ctx, key, out, intervalTTL)
}

func deriveTrend(scores []float64) string {
	n := len(scores)
	if n < 3 {
		return "stable"
	}
	recent := scores[n-3:]
	if recent[2] > recent[1] && recent[1] > recent[0] {
		return "improving"
	}
	if recent[2] < recent[1] && recent[1] < recent[0] {
		return "declining"
	}
	return "stable"
}
