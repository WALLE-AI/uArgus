package seed

import (
	"context"
	"fmt"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"time"
)

// luaCASRelease is the Lua script for CAS lock release —
// only deletes the key if the stored value matches the caller's runID.
const luaCASRelease = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
else
  return 0
end
`

// Lock provides distributed locking via Redis SET NX PX + Lua CAS release.
type Lock struct {
	client cache.Client
}

// NewLock creates a Lock.
func NewLock(client cache.Client) *Lock {
	return &Lock{client: client}
}

// Acquire tries to take the lock. Returns true if acquired.
func (l *Lock) Acquire(ctx context.Context, domain, resource, runID string, ttl time.Duration) (bool, error) {
	key := lockKey(domain, resource)
	// SET key runID NX PX ttlMs
	cmd := []cache.Cmd{{
		Op:   "SET",
		Args: []any{key, runID, "NX", "PX", int(ttl.Milliseconds())},
	}}
	results, err := l.client.Pipeline(ctx, cmd)
	if err != nil {
		return false, fmt.Errorf("seed lock acquire: %w", err)
	}
	// Upstash returns "OK" on success, null on failure
	if len(results) > 0 && results[0].Err == nil && string(results[0].Value) != "null" {
		return true, nil
	}
	return false, nil
}

// Release uses Lua CAS to release only if we still own the lock.
func (l *Lock) Release(ctx context.Context, domain, resource, runID string) error {
	key := lockKey(domain, resource)
	_, err := l.client.Eval(ctx, luaCASRelease, []string{key}, runID)
	if err != nil {
		return fmt.Errorf("seed lock release: %w", err)
	}
	return nil
}

func lockKey(domain, resource string) string {
	return fmt.Sprintf("seed-lock:%s:%s", domain, resource)
}
