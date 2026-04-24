package research

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/fallback"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// TrendingRepo represents a trending GitHub repository.
type TrendingRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	URL         string `json:"html_url"`
}

// TrendingSource fetches trending repos via OSSInsight → GitHub API fallback.
type TrendingSource struct {
	spec       registry.SourceSpec
	httpClient *fetcher.Client
}

// NewTrendingSource creates a TrendingSource.
func NewTrendingSource(spec registry.SourceSpec, httpClient *fetcher.Client) *TrendingSource {
	return &TrendingSource{spec: spec, httpClient: httpClient}
}

func (s *TrendingSource) Name() string              { return s.spec.CanonicalKey }
func (s *TrendingSource) Spec() registry.SourceSpec { return s.spec }
func (s *TrendingSource) Dependencies() []string    { return nil }

func (s *TrendingSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()

	chain := fallback.NewChain(
		fallback.Provider[[]TrendingRepo]{
			Name: "ossinsight",
			Fn:   func(ctx context.Context) ([]TrendingRepo, error) { return s.fetchOSSInsight(ctx) },
		},
		fallback.Provider[[]TrendingRepo]{
			Name: "github-api",
			Fn:   func(ctx context.Context) ([]TrendingRepo, error) { return s.fetchGitHub(ctx) },
		},
	)

	repos, err := chain.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return &registry.FetchResult{
		Data: repos,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(repos),
		},
	}, nil
}

type ossInsightRepo struct {
	RepoName    string `json:"repo_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Language    string `json:"language"`
	RepoID      int    `json:"repo_id"`
}

func (s *TrendingSource) fetchOSSInsight(ctx context.Context) ([]TrendingRepo, error) {
	body, status, err := s.httpClient.Get(ctx, "https://api.ossinsight.io/v1/trends/repos")
	if err != nil || status >= 400 {
		return nil, fmt.Errorf("ossinsight: status %d: %w", status, err)
	}
	// OSSInsight returns { data: { rows: [...] } }
	var resp struct {
		Data struct {
			Rows []ossInsightRepo `json:"rows"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	repos := make([]TrendingRepo, 0, len(resp.Data.Rows))
	for _, r := range resp.Data.Rows {
		repos = append(repos, TrendingRepo{
			Name:        r.RepoName,
			FullName:    r.RepoName,
			Description: r.Description,
			Stars:       r.Stars,
			Language:    r.Language,
			URL:         fmt.Sprintf("https://github.com/%s", r.RepoName),
		})
	}
	if len(repos) == 0 {
		return nil, fmt.Errorf("ossinsight: 0 repos returned")
	}
	return repos, nil
}

func (s *TrendingSource) fetchGitHub(ctx context.Context) ([]TrendingRepo, error) {
	cutoff := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:%%3E1000+pushed:%%3E%s&sort=stars&order=desc&per_page=30", cutoff)
	body, status, err := s.httpClient.Get(ctx, url)
	if err != nil || status >= 400 {
		return nil, fmt.Errorf("github: status %d: %w", status, err)
	}
	var resp struct {
		Items []TrendingRepo `json:"items"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}
