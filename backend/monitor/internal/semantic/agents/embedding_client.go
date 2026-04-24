package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// EmbeddingClient calls an OpenAI-compatible /v1/embeddings endpoint (e.g. BGE-M3
// served by vLLM, Ollama, or any OpenAI-compat embedding server).
type EmbeddingClient struct {
	baseURL string // e.g. "http://localhost:11434"
	model   string // e.g. "bge-m3"
	apiKey  string // optional bearer token
	http    *http.Client
}

// EmbeddingConfig configures an EmbeddingClient.
type EmbeddingConfig struct {
	BaseURL string // base URL of the embedding service (no trailing /v1/embeddings)
	Model   string // embedding model name
	APIKey  string // optional API key
}

// NewEmbeddingClient creates an EmbeddingClient.
// BaseURL may be either a service root (e.g. "http://localhost:11434") or
// the full endpoint URL (e.g. "https://api.siliconflow.cn/v1/embeddings").
// In the latter case the "/v1/embeddings" suffix is NOT appended again.
func NewEmbeddingClient(cfg EmbeddingConfig) *EmbeddingClient {
	model := cfg.Model
	if model == "" {
		model = "bge-m3"
	}
	base := strings.TrimRight(cfg.BaseURL, "/")
	return &EmbeddingClient{
		baseURL: base,
		model:   model,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// embeddingRequest is the request body for the /v1/embeddings endpoint.
type embeddingRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"` // string or []string
}

// embeddingResponse is the response body from the /v1/embeddings endpoint.
type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Embed generates a vector embedding for a single text.
func (c *EmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
	results, err := c.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("embedding: empty response")
	}
	return results[0], nil
}

// EmbedBatch generates vector embeddings for multiple texts in a single API call.
func (c *EmbeddingClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody := embeddingRequest{
		Model: c.model,
		Input: texts,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("embedding: marshal: %w", err)
	}

	embURL := c.baseURL
	if !strings.HasSuffix(embURL, "/embeddings") {
		embURL = c.baseURL + "/v1/embeddings"
	}
	url := embURL
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("embedding: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("embedding: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("embedding: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result embeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("embedding: unmarshal: %w", err)
	}

	// Re-order by index to match input order
	embeddings := make([][]float32, len(texts))
	for _, d := range result.Data {
		if d.Index >= 0 && d.Index < len(embeddings) {
			embeddings[d.Index] = d.Embedding
		}
	}

	return embeddings, nil
}
