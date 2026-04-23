package research

import (
	"context"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/parser"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

const techmemeICSURL = "https://www.techmeme.com/events.ics"

// TechEventsSource fetches tech events from Techmeme ICS.
type TechEventsSource struct {
	spec       registry.SourceSpec
	httpClient *fetcher.Client
}

// NewTechEventsSource creates a TechEventsSource.
func NewTechEventsSource(spec registry.SourceSpec, httpClient *fetcher.Client) *TechEventsSource {
	return &TechEventsSource{spec: spec, httpClient: httpClient}
}

func (s *TechEventsSource) Name() string             { return s.spec.CanonicalKey }
func (s *TechEventsSource) Spec() registry.SourceSpec { return s.spec }
func (s *TechEventsSource) Dependencies() []string    { return nil }

func (s *TechEventsSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()

	body, _, err := s.httpClient.Get(ctx, techmemeICSURL)
	if err != nil {
		return nil, err
	}

	events := parser.ParseICS(body)

	return &registry.FetchResult{
		Data: events,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(events),
		},
	}, nil
}
