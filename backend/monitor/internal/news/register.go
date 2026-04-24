package news

import (
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/agents"
)

// RegisterAll registers all news Sources (5 digest variants + 1 insights).
// relayBaseURL enables dual-path proxy fallback when non-empty.
func RegisterAll(reg *registry.Registry, rdb cache.Client, agentsClient agents.AgentsClient, relayBaseURL ...string) {
	httpClient := fetcher.NewClient(fetcher.ClientOpts{})
	relay := ""
	if len(relayBaseURL) > 0 {
		relay = relayBaseURL[0]
	}

	for _, v := range AllVariants() {
		canonicalKey := fmt.Sprintf("news:digest:v1:%s:en", v)
		spec := registry.SourceSpec{
			Domain:           "news",
			Resource:         fmt.Sprintf("digest-%s", v),
			Version:          1,
			CanonicalKey:     canonicalKey,
			DataTTL:          24 * time.Hour,
			MaxStaleDuration: 30 * time.Minute,
			MinRecordCount:   10,
			Schedule:         registry.CronSchedule{Expr: "0 */5 * * * *"}, // every 5 min
			LockTTL:          5 * time.Minute,
		}
		src := NewDigestSource(DigestOpts{
			Variant:      v,
			Lang:         "en",
			Spec:         spec,
			HTTPClient:   httpClient,
			RDB:          rdb,
			RelayBaseURL: relay,
		})
		reg.Register(src)
	}

	// insights depends on the full digest
	insightsSpec := registry.SourceSpec{
		Domain:           "news",
		Resource:         "insights",
		Version:          1,
		CanonicalKey:     "news:insights:v1",
		DataTTL:          24 * time.Hour,
		MaxStaleDuration: 60 * time.Minute,
		Schedule:         registry.CronSchedule{Expr: "0 */15 * * * *"}, // every 15 min
		LockTTL:          5 * time.Minute,
	}
	insightsSrc := NewInsightsSource(InsightsOpts{
		Spec:      insightsSpec,
		DigestKey: "news:digest:v1:full:en",
		RDB:       rdb,
		Agents:    agentsClient,
	})
	reg.Register(insightsSrc)
}
