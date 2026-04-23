package proxy

import "github.com/WALL-AI/uArgus/backend/monitor/internal/config"

// ProxyResolver resolves proxy URLs from configuration.
type ProxyResolver struct {
	cfg *config.Config
}

// NewProxyResolver creates a ProxyResolver.
func NewProxyResolver(cfg *config.Config) *ProxyResolver {
	return &ProxyResolver{cfg: cfg}
}

// RelayURL returns the relay proxy URL, if configured.
func (r *ProxyResolver) RelayURL() string { return r.cfg.ProxyRelayURL }

// ConnectURL returns the CONNECT proxy URL, if configured.
func (r *ProxyResolver) ConnectURL() string { return r.cfg.ProxyConnectURL }

// CurlURL returns the curl proxy URL, if configured.
func (r *ProxyResolver) CurlURL() string { return r.cfg.ProxyCurlURL }

// HasProxy returns true if any proxy is configured.
func (r *ProxyResolver) HasProxy() bool {
	return r.cfg.ProxyRelayURL != "" || r.cfg.ProxyConnectURL != "" || r.cfg.ProxyCurlURL != ""
}
