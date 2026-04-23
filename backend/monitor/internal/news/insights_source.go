package news

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/agents"
)

// InsightsSource reads the digest, extracts top headlines, and produces AI summaries.
type InsightsSource struct {
	spec      registry.SourceSpec
	digestKey string
	rdb       cache.Client
	agents    agents.AgentsClient
}

// InsightsOpts configures an InsightsSource.
type InsightsOpts struct {
	Spec      registry.SourceSpec
	DigestKey string
	RDB       cache.Client
	Agents    agents.AgentsClient
}

// NewInsightsSource creates an InsightsSource.
func NewInsightsSource(opts InsightsOpts) *InsightsSource {
	return &InsightsSource{
		spec:      opts.Spec,
		digestKey: opts.DigestKey,
		rdb:       opts.RDB,
		agents:    opts.Agents,
	}
}

func (s *InsightsSource) Name() string             { return s.spec.CanonicalKey }
func (s *InsightsSource) Spec() registry.SourceSpec { return s.spec }
func (s *InsightsSource) Dependencies() []string    { return []string{s.digestKey} }

func (s *InsightsSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()

	// 1. read digest from cache
	raw, err := s.rdb.Get(ctx, s.digestKey)
	if err != nil || raw == nil {
		return nil, fmt.Errorf("insights: digest not available at %s", s.digestKey)
	}

	// 2. extract top headlines
	var items []ParsedItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("insights: unmarshal digest: %w", err)
	}
	headlines := extractTopHeadlines(items, 10)

	// 3. summarise via agents
	summary, err := s.agents.Summarize(ctx, headlines, agents.SummarizeOpts{
		Mode:      "brief",
		MaxTokens: 500,
	})
	if err != nil {
		return nil, fmt.Errorf("insights: summarize: %w", err)
	}

	result := map[string]any{
		"summary":   summary,
		"headlines": headlines,
		"generatedAt": time.Now().UTC().Format(time.RFC3339),
	}

	return &registry.FetchResult{
		Data: result,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(headlines),
		},
	}, nil
}

func extractTopHeadlines(items []ParsedItem, n int) []string {
	if n > len(items) {
		n = len(items)
	}
	headlines := make([]string, n)
	for i := 0; i < n; i++ {
		headlines[i] = items[i].Title
	}
	return headlines
}
