package news

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/pool"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/classify"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/scoring"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/tiers"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/tracking"
)

// DigestSource fetches RSS feeds, classifies, scores, tracks, and produces the news digest.
type DigestSource struct {
	variant    string
	lang       string
	spec       registry.SourceSpec
	httpClient *fetcher.Client
	classifier *classify.Classifier
	scorer     *scoring.WeightedScorer
	tracker    *tracking.StoryTracker
	rdb        cache.Client
}

// DigestOpts configures a DigestSource.
type DigestOpts struct {
	Variant    string
	Lang       string
	Spec       registry.SourceSpec
	HTTPClient *fetcher.Client
	RDB        cache.Client
}

// NewDigestSource creates a DigestSource.
func NewDigestSource(opts DigestOpts) *DigestSource {
	return &DigestSource{
		variant:    opts.Variant,
		lang:       opts.Lang,
		spec:       opts.Spec,
		httpClient: opts.HTTPClient,
		classifier: classify.NewNewsClassifier(),
		scorer:     scoring.NewImportanceScorer(),
		tracker:    tracking.NewStoryTracker(opts.RDB),
		rdb:        opts.RDB,
	}
}

func (d *DigestSource) Name() string                { return d.spec.CanonicalKey }
func (d *DigestSource) Spec() registry.SourceSpec    { return d.spec }
func (d *DigestSource) Dependencies() []string       { return nil }

func (d *DigestSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()
	feeds := GetFeeds(d.variant, d.lang)
	logger := slog.With("source", d.Name(), "feeds", len(feeds))

	// 1. fetch all feeds in bounded pool (concurrency 20)
	type fetchResult struct {
		items []ParsedItem
	}
	results := pool.BoundedPool(ctx, 20, feeds, func(ctx context.Context, feed FeedEntry) (fetchResult, error) {
		body, status, err := d.httpClient.Get(ctx, feed.URL)
		if err != nil || status >= 400 {
			return fetchResult{}, err
		}
		items := ParseFeedResponse(body, feed)
		return fetchResult{items: items}, nil
	})

	// 2. collect all items
	var allItems []ParsedItem
	for _, r := range results {
		if r.Err == nil {
			allItems = append(allItems, r.Value.items...)
		}
	}
	logger.Info("feeds fetched", "totalItems", len(allItems))

	// 3. dedup
	seen := make(map[string]bool)
	var deduped []ParsedItem
	for _, item := range allItems {
		if !seen[item.Hash] {
			seen[item.Hash] = true
			deduped = append(deduped, item)
		}
	}

	// 4. classify + enrich tier
	for i := range deduped {
		cr := d.classifier.Classify(deduped[i].Title)
		deduped[i].Categories = cr.Categories
		deduped[i].Severity = cr.Severity
		deduped[i].Tier = tiers.GetTier(deduped[i].Source)
	}

	// 5. corroboration (from tracker)
	trackables := make([]tracking.Trackable, len(deduped))
	for i := range deduped {
		trackables[i] = deduped[i]
	}
	corrMap := d.tracker.ComputeCorroboration(trackables)

	// 6. score
	for i := range deduped {
		tierNorm := 1.0 - float64(deduped[i].Tier-1)/3.0 // 1→1.0, 4→0.0
		corrNorm := float64(corrMap[deduped[i].Hash]) / 10.0
		if corrNorm > 1.0 {
			corrNorm = 1.0
		}
		sevNorm := float64(deduped[i].Severity) / 7.0
		deduped[i].Score = d.scorer.Score(scoring.ScoringInput{
			"severity":      sevNorm,
			"tier":          tierNorm,
			"corroboration": corrNorm,
			"recency":       1.0, // TODO: compute from pubDate
		})
		deduped[i].Corroboration = corrMap[deduped[i].Hash]
	}

	// 7. track (write state)
	if err := d.tracker.Write(ctx, trackables); err != nil {
		logger.Warn("tracker write error", "err", err)
	}

	// 8. sort by score desc + truncate
	sort.Slice(deduped, func(i, j int) bool { return deduped[i].Score > deduped[j].Score })
	maxItems := 200
	if len(deduped) > maxItems {
		deduped = deduped[:maxItems]
	}

	return &registry.FetchResult{
		Data: deduped,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(deduped),
		},
	}, nil
}
