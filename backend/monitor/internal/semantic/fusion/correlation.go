package fusion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

// FuseInput describes one domain's data for correlation.
type FuseInput struct {
	Domain string
	Key    string
}

// CorrelationCard is a cross-domain correlation signal.
type CorrelationCard struct {
	ID       string   `json:"id"`
	Domains  []string `json:"domains"`
	Signal   string   `json:"signal"`
	Strength float64  `json:"strength"`
}

// CorrelationFuser reads multiple domain seed keys and produces cross-domain correlation cards.
type CorrelationFuser struct {
	rdb cache.Client
}

// NewCorrelationFuser creates a CorrelationFuser.
func NewCorrelationFuser(rdb cache.Client) *CorrelationFuser {
	return &CorrelationFuser{rdb: rdb}
}

// Fuse reads data from multiple domains and computes correlation cards.
func (f *CorrelationFuser) Fuse(ctx context.Context, sources []FuseInput) ([]CorrelationCard, error) {
	domainData := make(map[string]json.RawMessage, len(sources))

	for _, src := range sources {
		raw, err := f.rdb.Get(ctx, src.Key)
		if err != nil {
			return nil, fmt.Errorf("fusion: read %s: %w", src.Key, err)
		}
		if raw != nil {
			domainData[src.Domain] = raw
		}
	}

	if len(domainData) < 2 {
		return nil, nil // need at least 2 domains to correlate
	}

	// placeholder: real implementation would compute cross-domain signals
	var cards []CorrelationCard
	var domains []string
	for d := range domainData {
		domains = append(domains, d)
	}
	cards = append(cards, CorrelationCard{
		ID:       fmt.Sprintf("corr-%d", len(domainData)),
		Domains:  domains,
		Signal:   "cross-domain-detected",
		Strength: 0.5,
	})
	return cards, nil
}
