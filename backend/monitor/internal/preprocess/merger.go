package preprocess

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/pool"
)

// MultiSourceMerger is a Stage that reads multiple Redis keys and merges their data.
type MultiSourceMerger struct {
	rdb  cache.Client
	keys []string
}

// NewMultiSourceMerger creates a merger stage.
func NewMultiSourceMerger(rdb cache.Client, keys []string) *MultiSourceMerger {
	return &MultiSourceMerger{rdb: rdb, keys: keys}
}

func (m *MultiSourceMerger) Name() string { return "multi-source-merger" }

func (m *MultiSourceMerger) Process(ctx context.Context, _ any) (any, error) {
	type fetchResult struct {
		key  string
		data json.RawMessage
	}

	results := pool.BoundedPool(ctx, 5, m.keys, func(ctx context.Context, key string) (fetchResult, error) {
		raw, err := m.rdb.Get(ctx, key)
		if err != nil {
			return fetchResult{}, fmt.Errorf("merger read %s: %w", key, err)
		}
		return fetchResult{key: key, data: raw}, nil
	})

	merged := make(map[string]json.RawMessage, len(m.keys))
	for _, r := range results {
		if r.Err == nil && r.Value.data != nil {
			merged[r.Value.key] = r.Value.data
		}
	}

	return merged, nil
}
