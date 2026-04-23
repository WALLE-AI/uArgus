package tracking

import (
	"context"
	"fmt"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// Accumulator manages Redis sorted sets (ZADD/ZRANGEBYSCORE) for accumulating items.
type Accumulator struct {
	rdb cache.Client
}

// NewAccumulator creates an Accumulator.
func NewAccumulator(rdb cache.Client) *Accumulator {
	return &Accumulator{rdb: rdb}
}

// Add adds a member with a score to a sorted set via ZADD.
func (a *Accumulator) Add(ctx context.Context, key, member string, score float64) error {
	cmds := []cache.Cmd{{
		Op:   "ZADD",
		Args: []any{key, score, member},
	}}
	results, err := a.rdb.Pipeline(ctx, cmds)
	if err != nil {
		return fmt.Errorf("accumulator add: %w", err)
	}
	if results[0].Err != nil {
		return results[0].Err
	}
	return nil
}

// Trim keeps only the top N members by score (removes lower-scored ones).
func (a *Accumulator) Trim(ctx context.Context, key string, keep int) error {
	// ZREMRANGEBYRANK key 0 -(keep+1) removes all but top `keep`
	cmds := []cache.Cmd{{
		Op:   "ZREMRANGEBYRANK",
		Args: []any{key, 0, -(keep + 1)},
	}}
	results, err := a.rdb.Pipeline(ctx, cmds)
	if err != nil {
		return fmt.Errorf("accumulator trim: %w", err)
	}
	if results[0].Err != nil {
		return results[0].Err
	}
	return nil
}
