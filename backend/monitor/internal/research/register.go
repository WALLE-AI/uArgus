package research

import (
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// RegisterAll registers all research Sources (arxiv, hackernews, tech-events, trending).
func RegisterAll(reg *registry.Registry, rdb cache.Client) {
	httpClient := fetcher.NewClient(fetcher.ClientOpts{})

	reg.Register(NewArxivSource(registry.SourceSpec{
		Domain:           "research",
		Resource:         "arxiv",
		Version:          1,
		CanonicalKey:     "research:arxiv:v1",
		DataTTL:          24 * time.Hour,
		MaxStaleDuration: 6 * time.Hour,
		Schedule:         registry.CronSchedule{Expr: "0 0 */6 * * *"}, // every 6h
		LockTTL:          10 * time.Minute,
	}, httpClient))

	reg.Register(NewHackerNewsSource(registry.SourceSpec{
		Domain:           "research",
		Resource:         "hackernews",
		Version:          1,
		CanonicalKey:     "research:hackernews:v1",
		DataTTL:          1 * time.Hour,
		MaxStaleDuration: 30 * time.Minute,
		Schedule:         registry.CronSchedule{Expr: "0 */15 * * * *"}, // every 15 min
		LockTTL:          5 * time.Minute,
	}, httpClient))

	reg.Register(NewTechEventsSource(registry.SourceSpec{
		Domain:           "research",
		Resource:         "tech-events",
		Version:          1,
		CanonicalKey:     "research:tech-events:v1",
		DataTTL:          24 * time.Hour,
		MaxStaleDuration: 12 * time.Hour,
		Schedule:         registry.CronSchedule{Expr: "0 0 */12 * * *"}, // every 12h
		LockTTL:          5 * time.Minute,
	}, httpClient))

	reg.Register(NewTrendingSource(registry.SourceSpec{
		Domain:           "research",
		Resource:         "trending",
		Version:          1,
		CanonicalKey:     "research:trending:v1",
		DataTTL:          12 * time.Hour,
		MaxStaleDuration: 6 * time.Hour,
		Schedule:         registry.CronSchedule{Expr: "0 0 */6 * * *"}, // every 6h
		LockTTL:          5 * time.Minute,
	}, httpClient))
}
