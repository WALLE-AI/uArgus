package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm/provider"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/news"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/research"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/seed"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/agents"
)

// ── in-memory cache.Client mock ─────────────────────────────

type memClient struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemClient() *memClient { return &memClient{data: make(map[string][]byte)} }

func (m *memClient) Get(_ context.Context, k string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[k], nil
}
func (m *memClient) Set(_ context.Context, k string, v []byte, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
	return nil
}
func (m *memClient) Del(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}
func (m *memClient) Pipeline(_ context.Context, cmds []cache.Cmd) ([]cache.Result, error) {
	results := make([]cache.Result, len(cmds))
	for i, cmd := range cmds {
		switch cmd.Op {
		case "GET":
			if len(cmd.Args) > 0 {
				k, _ := cmd.Args[0].(string)
				m.mu.RLock()
				results[i] = cache.Result{Value: m.data[k]}
				m.mu.RUnlock()
			}
		case "SET":
			if len(cmd.Args) >= 2 {
				k, _ := cmd.Args[0].(string)
				v, _ := cmd.Args[1].(string)
				m.mu.Lock()
				m.data[k] = []byte(v)
				m.mu.Unlock()
			}
			results[i] = cache.Result{Value: []byte("OK")}
		default:
			results[i] = cache.Result{Value: []byte("OK")}
		}
	}
	return results, nil
}
func (m *memClient) Expire(_ context.Context, _ string, _ time.Duration) error { return nil }
func (m *memClient) Eval(_ context.Context, _ string, _ []string, _ ...any) (any, error) {
	return nil, nil
}
func (m *memClient) Info(_ context.Context) (map[string]string, error) {
	return map[string]string{"db0": fmt.Sprintf("keys=%d", len(m.data))}, nil
}

// ── helpers ─────────────────────────────────────────────────

func envFirst(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

// ── main ────────────────────────────────────────────────────

func main() {
	loadDotEnv(".env")
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	ctx := context.Background()

	// ensure output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		slog.Error("mkdir output", "err", err)
		os.Exit(1)
	}

	rdb := newMemClient()
	httpClient := fetcher.NewClient(fetcher.ClientOpts{Timeout: 30 * time.Second})
	relayURL := envFirst("PROXY_RELAY_URL")

	// ════════════════════════════════════════════════════════
	//  1. News RSS — quick variant (10 feeds, fast)
	// ════════════════════════════════════════════════════════
	fmt.Println("═══ 1. News RSS (quick variant) ════════════════")
	if relayURL != "" {
		fmt.Printf("[proxy] relay=%s\n", relayURL)
	}
	digestSrc := news.NewDigestSource(news.DigestOpts{
		Variant: "quick",
		Lang:    "en",
		Spec: registry.SourceSpec{
			Domain:       "news",
			Resource:     "digest-quick",
			Version:      1,
			CanonicalKey: "news:digest:v1:quick:en",
			DataTTL:      24 * time.Hour,
		},
		HTTPClient:   httpClient,
		RDB:          rdb,
		RelayBaseURL: relayURL,
	})

	newsStart := time.Now()
	newsResult, err := digestSrc.Run(ctx)
	newsDur := time.Since(newsStart)
	if err != nil {
		slog.Error("news digest failed", "err", err)
	} else {
		saveResult("output/news_digest.json", newsResult, newsDur, "news", "quick feeds")
	}

	// ════════════════════════════════════════════════════════
	//  2. News Insights — LLM summarisation (requires LLM_*)
	// ════════════════════════════════════════════════════════
	fmt.Println("\n═══ 2. News Insights (LLM) ═════════════════════")
	// Support both LLM_* and AGENT_ENGINE_* env var naming
	llmProvider := envFirst("LLM_PROVIDER", "AGENT_ENGINE_PROVIDER")
	llmKey := envFirst("LLM_API_KEY", "AGENT_ENGINE_API_KEY")
	llmModel := envFirst("LLM_MODEL", "AGENT_ENGINE_MODEL")
	llmBaseURL := envFirst("LLM_BASE_URL", "AGENT_ENGINE_BASE_URL")

	if llmProvider == "" || llmKey == "" {
		fmt.Println("[skip]  LLM env vars not set, skipping insights")
	} else {
		// Store digest data in memClient so InsightsSource can read it.
		// If real digest returned 0 items (e.g. network issues), inject
		// synthetic headlines so we can still verify LLM summarisation.
		digestKey := "news:digest:v1:quick:en"
		if newsResult != nil && newsResult.Metrics.RecordCount > 0 {
			digestJSON, _ := json.Marshal(newsResult.Data)
			_ = rdb.Set(ctx, digestKey, digestJSON, time.Hour)
		} else {
			fmt.Println("[info]  digest empty — injecting synthetic headlines for LLM test")
			synth := []news.ParsedItem{
				{Title: "Federal Reserve holds interest rates steady amid inflation concerns", Source: "Reuters", Hash: "h1"},
				{Title: "OpenAI announces GPT-5 with breakthrough reasoning capabilities", Source: "TechCrunch", Hash: "h2"},
				{Title: "EU passes landmark AI regulation framework effective 2027", Source: "BBC", Hash: "h3"},
				{Title: "Major earthquake strikes Turkey, at least 50 casualties reported", Source: "AP News", Hash: "h4"},
				{Title: "Tesla unveils next-gen humanoid robot for factory automation", Source: "Bloomberg", Hash: "h5"},
				{Title: "Global chip shortage eases as TSMC ramps 3nm production", Source: "Nikkei", Hash: "h6"},
				{Title: "US and China resume trade talks after months of tensions", Source: "WSJ", Hash: "h7"},
				{Title: "SpaceX Starship completes first successful orbital flight", Source: "CNN", Hash: "h8"},
			}
			synthJSON, _ := json.Marshal(synth)
			_ = rdb.Set(ctx, digestKey, synthJSON, time.Hour)
		}

		agentsClient, err := agents.NewDirectAgentsClient(agents.DirectAgentsConfig{
			ProviderType: provider.ProviderType(llmProvider),
			APIKey:       llmKey,
			BaseURL:      llmBaseURL,
			Model:        llmModel,
		})
		if err != nil {
			slog.Error("create LLM client", "err", err)
		} else {
			fmt.Printf("[llm]  provider=%s model=%s\n", llmProvider, llmModel)

			insightsSrc := news.NewInsightsSource(news.InsightsOpts{
				Spec: registry.SourceSpec{
					Domain:       "news",
					Resource:     "insights",
					Version:      1,
					CanonicalKey: "news:insights:v1",
					DataTTL:      24 * time.Hour,
				},
				DigestKey: digestKey,
				RDB:       rdb,
				Agents:    agentsClient,
			})

			insStart := time.Now()
			insResult, err := insightsSrc.Run(ctx)
			insDur := time.Since(insStart)
			if err != nil {
				slog.Error("insights failed", "err", err)
			} else {
				saveResult("output/news_insights.json", insResult, insDur, "insights", "LLM summary")
			}
		}
	}

	// ════════════════════════════════════════════════════════
	//  3. Research — GitHub Trending
	// ════════════════════════════════════════════════════════
	fmt.Println("\n═══ 3. Research GitHub Trending ═════════════════")
	trendingSrc := research.NewTrendingSource(registry.SourceSpec{
		Domain:       "research",
		Resource:     "trending",
		Version:      1,
		CanonicalKey: "research:trending:v1",
		DataTTL:      12 * time.Hour,
	}, httpClient)

	ghStart := time.Now()
	ghResult, err := trendingSrc.Run(ctx)
	ghDur := time.Since(ghStart)
	if err != nil {
		slog.Error("github trending failed", "err", err)
	} else {
		saveResult("output/research_trending.json", ghResult, ghDur, "research", "GitHub trending")
	}

	// ════════════════════════════════════════════════════════
	//  4. Research — Tech Events
	// ════════════════════════════════════════════════════════
	fmt.Println("\n═══ 4. Research Tech Events ═════════════════════")
	techEventsSrc := research.NewTechEventsSource(registry.SourceSpec{
		Domain:       "research",
		Resource:     "tech-events",
		Version:      1,
		CanonicalKey: "research:tech-events:v1",
		DataTTL:      24 * time.Hour,
	}, httpClient)

	teStart := time.Now()
	teResult, err := techEventsSrc.Run(ctx)
	teDur := time.Since(teStart)
	if err != nil {
		slog.Error("tech events failed", "err", err)
	} else {
		saveResult("output/research_tech_events.json", teResult, teDur, "research", "tech events")
	}

	// ════════════════════════════════════════════════════════
	//  5. Research — HackerNews
	// ════════════════════════════════════════════════════════
	fmt.Println("\n═══ 5. Research HackerNews ══════════════════════")
	hnSrc := research.NewHackerNewsSource(registry.SourceSpec{
		Domain:       "research",
		Resource:     "hackernews",
		Version:      1,
		CanonicalKey: "research:hackernews:v1",
		DataTTL:      1 * time.Hour,
	}, httpClient)

	hnStart := time.Now()
	hnResult, err := hnSrc.Run(ctx)
	hnDur := time.Since(hnStart)
	if err != nil {
		slog.Error("hackernews failed", "err", err)
	} else {
		saveResult("output/research_hackernews.json", hnResult, hnDur, "research", "HackerNews")
	}

	// ════════════════════════════════════════════════════════
	//  6. Research — arXiv
	// ════════════════════════════════════════════════════════
	fmt.Println("\n═══ 6. Research arXiv ═══════════════════════════")
	arxivSrc := research.NewArxivSource(registry.SourceSpec{
		Domain:       "research",
		Resource:     "arxiv",
		Version:      1,
		CanonicalKey: "research:arxiv:v1",
		DataTTL:      24 * time.Hour,
	}, httpClient)

	axStart := time.Now()
	axResult, err := arxivSrc.Run(ctx)
	axDur := time.Since(axStart)
	if err != nil {
		slog.Error("arxiv failed", "err", err)
	} else {
		saveResult("output/research_arxiv.json", axResult, axDur, "research", "arXiv papers")
	}

	fmt.Println("\n✓ full smoke test complete")
}

func saveResult(path string, result *registry.FetchResult, dur time.Duration, tag, label string) {
	envelope := seed.SeedEnvelope{
		Seed: &seed.SeedMeta{
			FetchedAt:     time.Now().UnixMilli(),
			RecordCount:   result.Metrics.RecordCount,
			SourceVersion: "v2",
			SchemaVersion: 1,
			State:         "OK",
		},
		Data: result.Data,
	}
	if wErr := writeJSON(path, envelope); wErr != nil {
		slog.Error("write json", "path", path, "err", wErr)
		return
	}
	fmt.Printf("[%-8s] %s → %d records, took %s → %s\n",
		tag, label, result.Metrics.RecordCount, dur.Round(time.Millisecond), path)
}
