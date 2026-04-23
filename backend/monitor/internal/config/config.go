package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the monitor service.
type Config struct {
	// Redis
	RedisURL   string
	RedisToken string

	// Agents
	AgentsURL string

	// Proxy
	ProxyRelayURL   string
	ProxyConnectURL string
	ProxyCurlURL    string

	// Server
	Port int
	Env  string // "development" | "production"

	// Feature flags
	SidecarMode bool

	// Timeouts (derived defaults)
	RedisReadTimeout     time.Duration
	RedisWriteTimeout    time.Duration
	RedisPipelineTimeout time.Duration
	RedisSeedTimeout     time.Duration
}

// Load reads configuration from environment variables and returns a validated Config.
func Load() (*Config, error) {
	c := &Config{
		RedisURL:   os.Getenv("UPSTASH_REDIS_REST_URL"),
		RedisToken: os.Getenv("UPSTASH_REDIS_REST_TOKEN"),

		AgentsURL: envOr("AGENTS_URL", "http://localhost:3001"),

		ProxyRelayURL:   os.Getenv("PROXY_RELAY_URL"),
		ProxyConnectURL: os.Getenv("PROXY_CONNECT_URL"),
		ProxyCurlURL:    os.Getenv("PROXY_CURL_URL"),

		Port: envInt("PORT", 8090),
		Env:  envOr("ENV", "development"),

		SidecarMode: envBool("SIDECAR_MODE", false),

		RedisReadTimeout:     1500 * time.Millisecond,
		RedisWriteTimeout:    5 * time.Second,
		RedisPipelineTimeout: 5 * time.Second,
		RedisSeedTimeout:     15 * time.Second,
	}

	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if c.RedisURL == "" {
		return fmt.Errorf("config: UPSTASH_REDIS_REST_URL is required")
	}
	if c.RedisToken == "" {
		return fmt.Errorf("config: UPSTASH_REDIS_REST_TOKEN is required")
	}
	return nil
}

func (c *Config) IsProd() bool {
	return c.Env == "production"
}

// ── helpers ─────────────────────────────────────────────────

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
