package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// Publisher handles atomic writes of seed data to Redis.
type Publisher struct {
	client cache.Client
}

// NewPublisher creates a Publisher.
func NewPublisher(client cache.Client) *Publisher {
	return &Publisher{client: client}
}

// AtomicPublish writes data to the canonical key using a pipeline (MULTI/EXEC semantics).
// Replaces v1's staging→canonical→del three-step approach.
func (p *Publisher) AtomicPublish(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	cmds := []cache.Cmd{
		{Op: "SET", Args: []any{key, string(data), "EX", int(ttl.Seconds())}},
	}
	results, err := p.client.Pipeline(ctx, cmds)
	if err != nil {
		return fmt.Errorf("seed publish: %w", err)
	}
	for _, r := range results {
		if r.Err != nil {
			return fmt.Errorf("seed publish pipeline: %w", r.Err)
		}
	}
	return nil
}

// PublishWithEnvelope builds an envelope and atomically publishes canonical + seed-meta.
func (p *Publisher) PublishWithEnvelope(ctx context.Context, key string, meta SeedMeta, data any, ttl time.Duration) error {
	envelope, err := Build(meta, data)
	if err != nil {
		return err
	}

	metaKey := seedMetaKeyFromCanonical(key)
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("seed publish: marshal meta: %w", err)
	}

	// Atomic: SET canonical + SET seed-meta in one pipeline
	metaTTL := resolveMetaTTL(ttl)
	cmds := []cache.Cmd{
		{Op: "SET", Args: []any{key, string(envelope), "EX", int(ttl.Seconds())}},
		{Op: "SET", Args: []any{metaKey, string(metaBytes), "EX", int(metaTTL.Seconds())}},
	}

	results, err := p.client.Pipeline(ctx, cmds)
	if err != nil {
		return fmt.Errorf("seed publish envelope: %w", err)
	}
	for _, r := range results {
		if r.Err != nil {
			return fmt.Errorf("seed publish envelope pipeline: %w", r.Err)
		}
	}
	return nil
}

func seedMetaKeyFromCanonical(key string) string {
	return "seed-meta:" + key
}

func resolveMetaTTL(dataTTL time.Duration) time.Duration {
	sevenDays := 7 * 24 * time.Hour
	if dataTTL > sevenDays {
		return dataTTL
	}
	return sevenDays
}
