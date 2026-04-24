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

// ChromaClient interacts with a Chroma vector database via its HTTP API.
type ChromaClient struct {
	baseURL    string // e.g. "http://localhost:8000"
	collection string // default collection name
	http       *http.Client
}

// ChromaConfig configures a ChromaClient.
type ChromaConfig struct {
	BaseURL    string // Chroma server URL
	Collection string // default collection name
}

// NewChromaClient creates a ChromaClient.
func NewChromaClient(cfg ChromaConfig) *ChromaClient {
	collection := cfg.Collection
	if collection == "" {
		collection = "news_embeddings"
	}
	return &ChromaClient{
		baseURL:    cfg.BaseURL,
		collection: collection,
		http:       &http.Client{Timeout: 30 * time.Second},
	}
}

// ---------------------------------------------------------------------------
// Collection management
// ---------------------------------------------------------------------------

// CreateCollection creates a collection in Chroma (idempotent via get_or_create).
func (c *ChromaClient) CreateCollection(ctx context.Context, name string) (string, error) {
	if name == "" {
		name = c.collection
	}
	body := map[string]any{
		"name":             name,
		"get_or_create":    true,
	}
	resp, err := c.post(ctx, "/api/v1/collections", body)
	if err != nil {
		return "", fmt.Errorf("chroma: create collection: %w", err)
	}
	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("chroma: unmarshal collection: %w", err)
	}
	return result.ID, nil
}

// getCollectionID looks up a collection ID by name.
func (c *ChromaClient) getCollectionID(ctx context.Context, name string) (string, error) {
	if name == "" {
		name = c.collection
	}
	url := fmt.Sprintf("%s/api/v1/collections/%s", c.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("chroma: get collection: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("chroma: get collection HTTP %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.ID, nil
}

// ---------------------------------------------------------------------------
// Upsert & Query
// ---------------------------------------------------------------------------

// UpsertInput holds data for a single vector upsert.
type UpsertInput struct {
	ID        string            `json:"id"`
	Embedding []float32         `json:"embedding"`
	Document  string            `json:"document,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Upsert adds or updates vectors in the collection.
func (c *ChromaClient) Upsert(ctx context.Context, collection string, items []UpsertInput) error {
	collID, err := c.getCollectionID(ctx, collection)
	if err != nil {
		return err
	}

	ids := make([]string, len(items))
	embeddings := make([][]float32, len(items))
	documents := make([]string, len(items))
	metadatas := make([]map[string]string, len(items))

	for i, item := range items {
		ids[i] = item.ID
		embeddings[i] = item.Embedding
		documents[i] = item.Document
		metadatas[i] = item.Metadata
	}

	body := map[string]any{
		"ids":        ids,
		"embeddings": embeddings,
		"documents":  documents,
		"metadatas":  metadatas,
	}

	url := fmt.Sprintf("/api/v1/collections/%s/upsert", collID)
	_, err = c.post(ctx, url, body)
	if err != nil {
		return fmt.Errorf("chroma: upsert: %w", err)
	}
	return nil
}

// QueryResult holds a single query result from Chroma.
type QueryResult struct {
	ID       string            `json:"id"`
	Document string            `json:"document"`
	Metadata map[string]string `json:"metadata"`
	Distance float64           `json:"distance"`
}

// Query performs a vector similarity search.
func (c *ChromaClient) Query(ctx context.Context, collection string, queryEmbedding []float32, nResults int) ([]QueryResult, error) {
	collID, err := c.getCollectionID(ctx, collection)
	if err != nil {
		return nil, err
	}

	if nResults <= 0 {
		nResults = 10
	}

	body := map[string]any{
		"query_embeddings": [][]float32{queryEmbedding},
		"n_results":        nResults,
		"include":          []string{"documents", "metadatas", "distances"},
	}

	url := fmt.Sprintf("/api/v1/collections/%s/query", collID)
	resp, err := c.post(ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf("chroma: query: %w", err)
	}

	// Chroma returns nested arrays: ids[][], documents[][], distances[][], metadatas[][]
	var raw struct {
		IDs       [][]string              `json:"ids"`
		Documents [][]string              `json:"documents"`
		Distances [][]float64             `json:"distances"`
		Metadatas [][]map[string]string   `json:"metadatas"`
	}
	if err := json.Unmarshal(resp, &raw); err != nil {
		return nil, fmt.Errorf("chroma: unmarshal query: %w", err)
	}

	if len(raw.IDs) == 0 || len(raw.IDs[0]) == 0 {
		return nil, nil
	}

	results := make([]QueryResult, len(raw.IDs[0]))
	for i := range raw.IDs[0] {
		r := QueryResult{ID: raw.IDs[0][i]}
		if len(raw.Documents) > 0 && len(raw.Documents[0]) > i {
			r.Document = raw.Documents[0][i]
		}
		if len(raw.Metadatas) > 0 && len(raw.Metadatas[0]) > i {
			r.Metadata = raw.Metadatas[0][i]
		}
		if len(raw.Distances) > 0 && len(raw.Distances[0]) > i {
			r.Distance = raw.Distances[0][i]
		}
		results[i] = r
	}

	return results, nil
}

// DeleteCollection deletes a collection by name.
func (c *ChromaClient) DeleteCollection(ctx context.Context, name string) error {
	if name == "" {
		name = c.collection
	}
	url := fmt.Sprintf("%s/api/v1/collections/%s", c.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("chroma: delete collection: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chroma: delete collection HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------

func (c *ChromaClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("chroma: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chroma: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("chroma: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("chroma: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
