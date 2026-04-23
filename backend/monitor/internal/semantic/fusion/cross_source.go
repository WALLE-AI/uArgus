package fusion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// CrossSignal represents a fused signal from multiple source keys.
type CrossSignal struct {
	ID      string   `json:"id"`
	Sources []string `json:"sources"`
	Type    string   `json:"type"`
	Score   float64  `json:"score"`
}

// CrossSourceFuser reads multiple seed keys and fuses them into signals.
type CrossSourceFuser struct {
	rdb cache.Client
}

// NewCrossSourceFuser creates a CrossSourceFuser.
func NewCrossSourceFuser(rdb cache.Client) *CrossSourceFuser {
	return &CrossSourceFuser{rdb: rdb}
}

// Fuse reads multiple keys and produces fused signals.
func (f *CrossSourceFuser) Fuse(ctx context.Context, keys []string) ([]CrossSignal, error) {
	data := make(map[string]json.RawMessage)
	for _, k := range keys {
		raw, err := f.rdb.Get(ctx, k)
		if err != nil {
			return nil, fmt.Errorf("cross-source: read %s: %w", k, err)
		}
		if raw != nil {
			data[k] = raw
		}
	}

	if len(data) < 2 {
		return nil, nil
	}

	// placeholder: real implementation would compute cross-source fusion
	var sources []string
	for k := range data {
		sources = append(sources, k)
	}
	return []CrossSignal{{
		ID:      fmt.Sprintf("xs-%d", len(data)),
		Sources: sources,
		Type:    "cross-source-signal",
		Score:   0.5,
	}}, nil
}
