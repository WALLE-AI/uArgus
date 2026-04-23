package research

import (
	"context"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/parser"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/ratelimit"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

var arxivCategories = []string{"cs.AI", "cs.CR", "cs.LG"}

// ArxivSource fetches recent papers from arXiv Atom API.
type ArxivSource struct {
	spec       registry.SourceSpec
	httpClient *fetcher.Client
	limiter    ratelimit.RateLimiter
}

// NewArxivSource creates an ArxivSource.
func NewArxivSource(spec registry.SourceSpec, httpClient *fetcher.Client) *ArxivSource {
	return &ArxivSource{
		spec:       spec,
		httpClient: httpClient,
		limiter:    ratelimit.NewFixedInterval(3 * time.Second),
	}
}

func (s *ArxivSource) Name() string                { return s.spec.CanonicalKey }
func (s *ArxivSource) Spec() registry.SourceSpec    { return s.spec }
func (s *ArxivSource) Dependencies() []string       { return nil }

func (s *ArxivSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()
	var allPapers []parser.AtomEntry

	for _, cat := range arxivCategories {
		if err := s.limiter.Wait(ctx); err != nil {
			return nil, err
		}
		url := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=cat:%s&sortBy=submittedDate&sortOrder=descending&max_results=50", cat)
		body, status, err := s.httpClient.Get(ctx, url)
		if err != nil || status >= 400 {
			continue
		}
		entries, err := parser.ParseAtomXML(body)
		if err != nil {
			continue
		}
		allPapers = append(allPapers, entries...)
	}

	return &registry.FetchResult{
		Data: allPapers,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(allPapers),
		},
	}, nil
}
