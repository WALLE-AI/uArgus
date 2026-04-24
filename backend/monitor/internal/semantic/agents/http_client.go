package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPAgentsClient implements AgentsClient via HTTP calls to the agents service.
type HTTPAgentsClient struct {
	baseURL string
	http    *http.Client
}

// NewHTTPAgentsClient creates an HTTPAgentsClient.
func NewHTTPAgentsClient(baseURL string) *HTTPAgentsClient {
	return &HTTPAgentsClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPAgentsClient) Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error) {
	body := map[string]any{
		"texts":     texts,
		"mode":      opts.Mode,
		"maxTokens": opts.MaxTokens,
	}
	resp, err := c.post(ctx, "/summarize", body)
	if err != nil {
		return "", err
	}
	var result struct {
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("agents: unmarshal: %w", err)
	}
	return result.Summary, nil
}

func (c *HTTPAgentsClient) Classify(ctx context.Context, text string) (*AiClassification, error) {
	body := map[string]any{"text": text}
	resp, err := c.post(ctx, "/classify", body)
	if err != nil {
		return nil, err
	}
	var result AiClassification
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("agents: unmarshal: %w", err)
	}
	return &result, nil
}

func (c *HTTPAgentsClient) Sentiment(ctx context.Context, texts []string) ([]SentimentResult, error) {
	body := map[string]any{"texts": texts}
	resp, err := c.post(ctx, "/sentiment", body)
	if err != nil {
		return nil, err
	}
	var result []SentimentResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("agents: unmarshal: %w", err)
	}
	return result, nil
}

func (c *HTTPAgentsClient) ExtractEntities(ctx context.Context, texts []string) ([][]NEREntity, error) {
	body := map[string]any{"texts": texts}
	resp, err := c.post(ctx, "/ner", body)
	if err != nil {
		return nil, err
	}
	var result [][]NEREntity
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("agents: unmarshal: %w", err)
	}
	return result, nil
}

func (c *HTTPAgentsClient) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	body := map[string]any{"text": text, "targetLang": targetLang}
	resp, err := c.post(ctx, "/translate", body)
	if err != nil {
		return "", err
	}
	var result struct {
		Translation string `json:"translation"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("agents: unmarshal: %w", err)
	}
	return result.Translation, nil
}

func (c *HTTPAgentsClient) Embed(ctx context.Context, text string) ([]float32, error) {
	body := map[string]any{"text": text}
	resp, err := c.post(ctx, "/embed", body)
	if err != nil {
		return nil, err
	}
	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("agents: unmarshal: %w", err)
	}
	return result.Embedding, nil
}

func (c *HTTPAgentsClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("agents: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("agents: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("agents: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("agents: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
