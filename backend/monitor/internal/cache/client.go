package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the single Redis interaction interface — reads and writes share it.
type Client interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Pipeline(ctx context.Context, cmds []Cmd) ([]Result, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)
	Info(ctx context.Context) (map[string]string, error)
}

// Cmd represents one Redis command in a pipeline.
type Cmd struct {
	Op   string // "SET", "GET", "DEL", "EXPIRE", …
	Args []any
}

// Result is one response from a pipeline.
type Result struct {
	Value []byte
	Err   error
}

// ClientConfig holds timeout / retry settings.
type ClientConfig struct {
	ReadTimeout     time.Duration // default 1.5s
	WriteTimeout    time.Duration // default 5s
	PipelineTimeout time.Duration // default 5s
	SeedTimeout     time.Duration // default 15s
	MaxRetries      int           // default 1
}

// DefaultClientConfig returns production defaults.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		ReadTimeout:     1500 * time.Millisecond,
		WriteTimeout:    5 * time.Second,
		PipelineTimeout: 5 * time.Second,
		SeedTimeout:     15 * time.Second,
		MaxRetries:      1,
	}
}

// ── Upstash REST implementation ─────────────────────────────

type upstashClient struct {
	url   string
	token string
	cfg   ClientConfig
	http  *http.Client
}

// NewUpstashClient creates a Client backed by Upstash REST API.
func NewUpstashClient(url, token string, cfg ClientConfig) Client {
	return &upstashClient{
		url:   url,
		token: token,
		cfg:   cfg,
		http:  &http.Client{Timeout: cfg.ReadTimeout},
	}
}

func (c *upstashClient) do(ctx context.Context, timeout time.Duration, body any) (json.RawMessage, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cache: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := c.http
	if timeout != c.cfg.ReadTimeout {
		client = &http.Client{Timeout: timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cache: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cache: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cache: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var envelope struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return envelope.Result, nil
}

func (c *upstashClient) Get(ctx context.Context, key string) ([]byte, error) {
	raw, err := c.do(ctx, c.cfg.ReadTimeout, []string{"GET", key})
	if err != nil {
		return nil, err
	}
	// Upstash returns JSON null for missing keys
	if string(raw) == "null" {
		return nil, nil
	}
	// unwrap JSON string
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return raw, nil // return raw bytes if not a string
	}
	return []byte(s), nil
}

func (c *upstashClient) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	cmd := []any{"SET", key, string(value)}
	if ttl > 0 {
		cmd = append(cmd, "EX", int(ttl.Seconds()))
	}
	_, err := c.do(ctx, c.cfg.WriteTimeout, cmd)
	return err
}

func (c *upstashClient) Del(ctx context.Context, keys ...string) error {
	cmd := make([]any, 0, 1+len(keys))
	cmd = append(cmd, "DEL")
	for _, k := range keys {
		cmd = append(cmd, k)
	}
	_, err := c.do(ctx, c.cfg.WriteTimeout, cmd)
	return err
}

func (c *upstashClient) Pipeline(ctx context.Context, cmds []Cmd) ([]Result, error) {
	batch := make([][]any, len(cmds))
	for i, cmd := range cmds {
		row := make([]any, 0, 1+len(cmd.Args))
		row = append(row, cmd.Op)
		row = append(row, cmd.Args...)
		batch[i] = row
	}

	payload, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/pipeline", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: c.cfg.PipelineTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cache: pipeline request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cache: pipeline HTTP %d: %s", resp.StatusCode, string(body))
	}

	var results []struct {
		Result json.RawMessage `json:"result"`
		Error  string          `json:"error"`
	}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("cache: pipeline unmarshal: %w", err)
	}

	out := make([]Result, len(results))
	for i, r := range results {
		if r.Error != "" {
			out[i] = Result{Err: fmt.Errorf("cache: pipeline[%d]: %s", i, r.Error)}
		} else {
			out[i] = Result{Value: r.Result}
		}
	}
	return out, nil
}

func (c *upstashClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	_, err := c.do(ctx, c.cfg.WriteTimeout, []any{"EXPIRE", key, int(ttl.Seconds())})
	return err
}

func (c *upstashClient) Info(ctx context.Context) (map[string]string, error) {
	raw, err := c.do(ctx, c.cfg.ReadTimeout, []string{"INFO"})
	if err != nil {
		return nil, err
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		// try parsing as key-value pairs from structured response
		result := map[string]string{"raw": string(raw)}
		return result, nil
	}
	return parseInfoString(s), nil
}

func parseInfoString(s string) map[string]string {
	m := make(map[string]string)
	for _, line := range bytes.Split([]byte(s), []byte("\r\n")) {
		l := string(line)
		if l == "" || l[0] == '#' {
			continue
		}
		if idx := bytes.IndexByte(line, ':'); idx > 0 {
			m[string(line[:idx])] = string(line[idx+1:])
		}
	}
	return m
}

func (c *upstashClient) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	cmd := make([]any, 0, 3+len(keys)+len(args))
	cmd = append(cmd, "EVAL", script, len(keys))
	for _, k := range keys {
		cmd = append(cmd, k)
	}
	cmd = append(cmd, args...)
	raw, err := c.do(ctx, c.cfg.WriteTimeout, cmd)
	if err != nil {
		return nil, err
	}
	var v any
	_ = json.Unmarshal(raw, &v)
	return v, nil
}
