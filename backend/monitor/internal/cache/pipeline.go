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
