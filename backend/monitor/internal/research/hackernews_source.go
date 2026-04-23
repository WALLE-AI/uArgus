package research

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/pagination"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// HNStory represents a Hacker News story.
type HNStory struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Descendants int    `json:"descendants"`
}

// HackerNewsSource fetches top stories from Hacker News API.
type HackerNewsSource struct {
	spec       registry.SourceSpec
	httpClient *fetcher.Client
}

// NewHackerNewsSource creates a HackerNewsSource.
func NewHackerNewsSource(spec registry.SourceSpec, httpClient *fetcher.Client) *HackerNewsSource {
	return &HackerNewsSource{spec: spec, httpClient: httpClient}
}

func (s *HackerNewsSource) Name() string             { return s.spec.CanonicalKey }
func (s *HackerNewsSource) Spec() registry.SourceSpec { return s.spec }
func (s *HackerNewsSource) Dependencies() []string    { return nil }

func (s *HackerNewsSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()

	// 1. fetch top story IDs
	body, _, err := s.httpClient.Get(ctx, "https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, fmt.Errorf("hn: topstories: %w", err)
	}
	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, fmt.Errorf("hn: unmarshal ids: %w", err)
	}
	if len(ids) > 500 {
		ids = ids[:500]
	}

	// 2. batch fetch individual stories
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = fmt.Sprintf("%d", id)
	}

	batcher := pagination.NewIDBatchFetcher[HNStory](50, 10)
	stories, err := batcher.FetchByIDs(ctx, strIDs, func(ctx context.Context, batch []string) ([]HNStory, error) {
		var results []HNStory
		for _, idStr := range batch {
			url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", idStr)
			b, _, fetchErr := s.httpClient.Get(ctx, url)
			if fetchErr != nil {
				continue
			}
			var story HNStory
			if err := json.Unmarshal(b, &story); err == nil && story.Title != "" {
				results = append(results, story)
			}
		}
		return results, nil
	})
	if err != nil {
		return nil, fmt.Errorf("hn: batch fetch: %w", err)
	}

	return &registry.FetchResult{
		Data: stories,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(stories),
		},
	}, nil
}
