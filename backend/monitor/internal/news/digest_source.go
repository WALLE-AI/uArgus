package news

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/pool"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/proxy"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/classify"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/scoring"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/tiers"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/tracking"
)

const digestCacheTTL = 900 // 15 minutes — matches v1 cachedFetchJson TTL

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
	cache      *cache.FetchThrough[[]ParsedItem]
	strategy   proxy.Strategy // nil = direct only
	relayHTTP  *fetcher.Client
}

// DigestOpts configures a DigestSource.
type DigestOpts struct {
	Variant      string
	Lang         string
	Spec         registry.SourceSpec
	HTTPClient   *fetcher.Client
	RDB          cache.Client
	RelayBaseURL string // optional; enables DirectFirst proxy fallback
}

// NewDigestSource creates a DigestSource.
func NewDigestSource(opts DigestOpts) *DigestSource {
	ds := &DigestSource{
		variant:    opts.Variant,
		lang:       opts.Lang,
		spec:       opts.Spec,
		httpClient: opts.HTTPClient,
		classifier: classify.NewNewsClassifier(),
		scorer:     scoring.NewImportanceScorer(),
		tracker:    tracking.NewStoryTracker(opts.RDB),
		rdb:        opts.RDB,
		cache:      cache.NewFetchThrough[[]ParsedItem](opts.RDB, nil),
	}
	if opts.RelayBaseURL != "" {
		ds.strategy = proxy.DirectFirst{}
		ds.relayHTTP = fetcher.NewClient(fetcher.ClientOpts{
			Timeout:   10 * time.Second,
			UserAgent: "uArgus-Monitor/2.0 relay",
		})
	}
	return ds
}

func (d *DigestSource) Name() string              { return d.spec.CanonicalKey }
func (d *DigestSource) Spec() registry.SourceSpec { return d.spec }
func (d *DigestSource) Dependencies() []string    { return nil }

func (d *DigestSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	cacheKey := fmt.Sprintf("news:digest-cache:v2:%s:%s", d.variant, d.lang)
	result, err := d.cache.Fetch(ctx, cache.FetchOpts[[]ParsedItem]{
		Key: cacheKey,
		TTL: digestCacheTTL * time.Second,
		Fetcher: func(ctx context.Context) (*[]ParsedItem, error) {
			items, err := d.buildDigest(ctx)
			if err != nil {
				return nil, err
			}
			return &items, nil
		},
	})
	if err != nil {
		return nil, err
	}
	var items []ParsedItem
	if result != nil {
		items = *result
	}
	return &registry.FetchResult{
		Data: items,
		Metrics: registry.FetchMetrics{
			RecordCount: len(items),
		},
	}, nil
}

// buildDigest performs the full fetch → parse → classify → score → track pipeline.
func (d *DigestSource) buildDigest(ctx context.Context) ([]ParsedItem, error) {
	start := time.Now()
	feeds := GetFeeds(d.variant, d.lang)
	logger := slog.With("source", d.Name(), "feeds", len(feeds))

	// 1. fetch all feeds in bounded pool (concurrency 20)
	type fetchResult struct {
		items []ParsedItem
	}
	results := pool.BoundedPool(ctx, 20, feeds, func(ctx context.Context, feed FeedEntry) (fetchResult, error) {
		var body []byte
		var status int
		var err error

		if d.strategy != nil {
			body, status, err = d.strategy.Execute(ctx, feed.URL,
				d.httpClient.Get,
				d.relayHTTP.Get,
			)
		} else {
			body, status, err = d.httpClient.Get(ctx, feed.URL)
		}

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

	logger.Info("digest built", "items", len(deduped), "duration", time.Since(start))
	return deduped, nil
}
