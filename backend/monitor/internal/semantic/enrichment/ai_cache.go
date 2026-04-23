package enrichment

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// ClassifiedItem is the interface for items that can be enriched with AI classification data.
type ClassifiedItem interface {
	EnrichHash() string
	SetAIClassification(data json.RawMessage)
}

// Enricher performs batch AI cache enrichment via MGET.
type Enricher struct {
	rdb cache.Client
}

// New creates an Enricher.
func New(rdb cache.Client) *Enricher {
	return &Enricher{rdb: rdb}
}

// EnrichBatch attempts to fill AI classification fields from cached MGET results.
// Keys follow the pattern: classify:sebuf:v3:{hash}
func (e *Enricher) EnrichBatch(ctx context.Context, items []ClassifiedItem) error {
	if len(items) == 0 {
		return nil
	}

	cmds := make([]cache.Cmd, len(items))
	for i, item := range items {
		cmds[i] = cache.Cmd{
			Op:   "GET",
			Args: []any{"classify:sebuf:v3:" + item.EnrichHash()},
		}
	}

	results, err := e.rdb.Pipeline(ctx, cmds)
	if err != nil {
		slog.Warn("enrichment: pipeline error", "err", err)
		return nil // non-fatal
	}

	for i, r := range results {
		if r.Err != nil || r.Value == nil || string(r.Value) == "null" {
			continue
		}
		items[i].SetAIClassification(r.Value)
	}
	return nil
}
