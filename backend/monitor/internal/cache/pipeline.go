package cache

import (
	"context"
	"time"
)

// PipelineSet is a convenience for batching multiple SET commands.
func PipelineSet(ctx context.Context, c Client, entries map[string][]byte, ttl time.Duration) error {
	cmds := make([]Cmd, 0, len(entries))
	for key, val := range entries {
		args := []any{key, string(val)}
		if ttl > 0 {
			args = append(args, "EX", int(ttl.Seconds()))
		}
		cmds = append(cmds, Cmd{Op: "SET", Args: args})
	}
	results, err := c.Pipeline(ctx, cmds)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

// PipelineGet fetches multiple keys in a single pipeline. Returns a map of key → raw bytes.
// Missing keys are omitted from the result (no error).
func PipelineGet(ctx context.Context, c Client, keys []string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	cmds := make([]Cmd, len(keys))
	for i, k := range keys {
		cmds[i] = Cmd{Op: "GET", Args: []any{k}}
	}
	results, err := c.Pipeline(ctx, cmds)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]byte, len(keys))
	for i, r := range results {
		if r.Err == nil && r.Value != nil && string(r.Value) != "null" {
			out[keys[i]] = r.Value
		}
	}
	return out, nil
}

// PipelineExpire extends TTL on multiple keys in a single pipeline.
func PipelineExpire(ctx context.Context, c Client, keys []string, ttl time.Duration) error {
	cmds := make([]Cmd, len(keys))
	for i, k := range keys {
		cmds[i] = Cmd{Op: "EXPIRE", Args: []any{k, int(ttl.Seconds())}}
	}
	results, err := c.Pipeline(ctx, cmds)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}
