package proxy

import (
	"context"
	"fmt"
)

// FetchFunc is a function that performs an HTTP fetch and returns (body, statusCode, err).
type FetchFunc func(ctx context.Context, url string) ([]byte, int, error)

// Strategy defines how to attempt direct vs proxy fetches.
type Strategy interface {
	Execute(ctx context.Context, url string, direct, proxied FetchFunc) ([]byte, int, error)
}

// DirectFirst tries direct, falls back to proxy on failure.
type DirectFirst struct{}

func (DirectFirst) Execute(ctx context.Context, url string, direct, proxied FetchFunc) ([]byte, int, error) {
	body, status, err := direct(ctx, url)
	if err == nil && status < 400 {
		return body, status, nil
	}
	return proxied(ctx, url)
}

// ProxyFirst tries proxy first, falls back to direct.
type ProxyFirst struct{}

func (ProxyFirst) Execute(ctx context.Context, url string, direct, proxied FetchFunc) ([]byte, int, error) {
	body, status, err := proxied(ctx, url)
	if err == nil && status < 400 {
		return body, status, nil
	}
	return direct(ctx, url)
}

// CurlOnly always uses the curl proxy.
type CurlOnly struct{}

func (CurlOnly) Execute(ctx context.Context, url string, _, proxied FetchFunc) ([]byte, int, error) {
	return proxied(ctx, url)
}

// TwoLegCascade tries two proxy endpoints in sequence.
type TwoLegCascade struct {
	SecondProxy FetchFunc
}

func (t TwoLegCascade) Execute(ctx context.Context, url string, direct, proxied FetchFunc) ([]byte, int, error) {
	body, status, err := proxied(ctx, url)
	if err == nil && status < 400 {
		return body, status, nil
	}
	if t.SecondProxy != nil {
		body, status, err = t.SecondProxy(ctx, url)
		if err == nil && status < 400 {
			return body, status, nil
		}
	}
	return nil, 0, fmt.Errorf("proxy: both legs failed for %s", url)
}
