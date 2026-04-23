package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a configurable HTTP client for data fetching.
type Client struct {
	http      *http.Client
	userAgent string
}

// ClientOpts configures the fetch Client.
type ClientOpts struct {
	Timeout   time.Duration
	UserAgent string
}

// NewClient creates a fetcher Client.
func NewClient(opts ClientOpts) *Client {
	if opts.Timeout == 0 {
		opts.Timeout = 15 * time.Second
	}
	if opts.UserAgent == "" {
		opts.UserAgent = "uArgus-Monitor/2.0"
	}
	return &Client{
		http:      &http.Client{Timeout: opts.Timeout},
		userAgent: opts.UserAgent,
	}
}

// Get performs an HTTP GET and returns the response body bytes.
func (c *Client) Get(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("fetcher: new request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("fetcher: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("fetcher: read body: %w", err)
	}
	return body, resp.StatusCode, nil
}

// Post performs an HTTP POST with the given body.
func (c *Client) Post(ctx context.Context, url string, contentType string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("fetcher: new request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("fetcher: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("fetcher: read body: %w", err)
	}
	return body, resp.StatusCode, nil
}
