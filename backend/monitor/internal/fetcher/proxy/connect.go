package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// ConnectTunnel creates an HTTPS CONNECT tunnel through a proxy.
func ConnectTunnel(ctx context.Context, proxyAddr, targetHost string, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("proxy connect: dial: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodConnect, "http://"+targetHost, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	req.Host = targetHost

	if err := req.Write(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("proxy connect: write CONNECT: %w", err)
	}

	resp, err := http.ReadResponse(nil, req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("proxy connect: read response: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		conn.Close()
		return nil, fmt.Errorf("proxy connect: status %d", resp.StatusCode)
	}
	return conn, nil
}
