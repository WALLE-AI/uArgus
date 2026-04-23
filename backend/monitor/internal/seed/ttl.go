package seed

import (
	"context"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// TTLManager handles TTL extension and resolution from SourceSpec.
type TTLManager struct {
	client cache.Client
}

// NewTTLManager creates a TTLManager.
func NewTTLManager(client cache.Client) *TTLManager {
	return &TTLManager{client: client}
}

// ExtendExisting extends the TTL on existing keys (resilience: don't lose data on failure).
func (t *TTLManager) ExtendExisting(ctx context.Context, keys []string, ttl time.Duration) error {
	return cache.PipelineExpire(ctx, t.client, keys, ttl)
}

// ResolveTTL reads DataTTL from SourceSpec.
func ResolveTTL(spec registry.SourceSpec) time.Duration {
	if spec.DataTTL > 0 {
		return spec.DataTTL
	}
	return 24 * time.Hour // fallback
}

// ResolveMetaTTL returns max(7d, DataTTL) for seed-meta keys.
func ResolveMetaTTL(spec registry.SourceSpec) time.Duration {
	sevenDays := 7 * 24 * time.Hour
	if spec.DataTTL > sevenDays {
		return spec.DataTTL
	}
	return sevenDays
}
