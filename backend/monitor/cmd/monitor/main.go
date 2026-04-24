package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/config"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm/provider"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/news"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/research"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/seed"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/semantic/agents"
)

func main() {
	// ── config ──────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// ── logger ──────────────────────────────────────────────
	level := slog.LevelInfo
	if !cfg.IsProd() {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
	slog.Info("monitor starting", "env", cfg.Env, "port", cfg.Port)

	// ── graceful shutdown context ───────────────────────────
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ── cache + metrics ─────────────────────────────────────
	rdb := cache.NewUpstashClient(cfg.RedisURL, cfg.RedisToken, cache.DefaultClientConfig())
	promReg := prometheus.NewRegistry()
	promReg.MustRegister(prometheus.NewGoCollector(), prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	metrics := cache.NewMetrics(promReg)

	// ── seed runner ─────────────────────────────────────────
	runner := seed.NewRunner(rdb, metrics)

	// ── registry ────────────────────────────────────────────
	reg := registry.New(runner)

	// ── agents client (direct LLM or HTTP fallback) ────────
	var agentsClient agents.AgentsClient
	if cfg.LLMProvider != "" {
		direct, err := agents.NewDirectAgentsClient(agents.DirectAgentsConfig{
			ProviderType: provider.ProviderType(cfg.LLMProvider),
			APIKey:       cfg.LLMAPIKey,
			BaseURL:      cfg.LLMBaseURL,
			Model:        cfg.LLMModel,
		})
		if err != nil {
			slog.Error("failed to create direct LLM client, falling back to HTTP", "err", err)
			agentsClient = agents.NewHTTPAgentsClient(cfg.AgentsURL)
		} else {
			slog.Info("using direct LLM provider", "type", cfg.LLMProvider, "model", cfg.LLMModel)
			agentsClient = direct
		}
	} else {
		agentsClient = agents.NewHTTPAgentsClient(cfg.AgentsURL)
	}

	// ── register sources ────────────────────────────────────
	news.RegisterAll(reg, rdb, agentsClient, cfg.ProxyRelayURL)
	research.RegisterAll(reg, rdb)

	// ── boot ────────────────────────────────────────────────
	if err := reg.Boot(ctx); err != nil {
		slog.Error("registry boot failed", "err", err)
		os.Exit(1)
	}

	// ── HTTP server (/healthz + /metrics) ───────────────────
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		snap := reg.HealthSnapshot()
		redisInfo, _ := rdb.Info(r.Context())

		// Compute aggregate status from individual source health states.
		aggregate := "ok"
		var degraded, failing []string
		for name, s := range snap {
			switch s.State {
			case registry.HealthStateFailing:
				failing = append(failing, name)
			case registry.HealthStateDegraded:
				degraded = append(degraded, name)
			}
		}
		if len(failing) > 0 {
			aggregate = "failing"
		} else if len(degraded) > 0 {
			aggregate = "degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		if aggregate == "failing" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status":   aggregate,
			"degraded": degraded,
			"failing":  failing,
			"sources":  snap,
			"redis":    redisInfo,
			"ts":       time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.Handle("/metrics", promhttp.HandlerFor(promReg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	go func() {
		slog.Info("HTTP server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "err", err)
		}
	}()

	// ── wait for shutdown signal ────────────────────────────
	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = reg.Shutdown(shutdownCtx)
	_ = srv.Shutdown(shutdownCtx)

	slog.Info("monitor stopped")
}
