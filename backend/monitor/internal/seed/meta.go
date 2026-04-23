package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// MetaWriter writes seed freshness metadata to Redis.
type MetaWriter struct {
	client cache.Client
}

// NewMetaWriter creates a MetaWriter.
func NewMetaWriter(client cache.Client) *MetaWriter {
	return &MetaWriter{client: client}
}

// WriteFreshness writes the seed-meta:{domain}:{resource} key.
func (w *MetaWriter) WriteFreshness(ctx context.Context, domain, resource string, count int, source string, ttl time.Duration) error {
	key := fmt.Sprintf("seed-meta:%s:%s", domain, resource)
	meta := map[string]any{
		"fetchedAt":   time.Now().UnixMilli(),
		"recordCount": count,
		"status":      "OK",
		"source":      source,
	}
	if count == 0 {
		meta["status"] = "OK_ZERO"
	}

	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("seed meta: marshal: %w", err)
	}
	return w.client.Set(ctx, key, data, ttl)
}
