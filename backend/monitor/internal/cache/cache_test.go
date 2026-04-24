package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// ── in-memory mock Client ───────────────────────────────────

type memClient struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newMemClient() *memClient {
	return &memClient{data: make(map[string][]byte)}
}

func (m *memClient) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (m *memClient) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *memClient) Del(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

func (m *memClient) Pipeline(_ context.Context, cmds []Cmd) ([]Result, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	results := make([]Result, len(cmds))
	for i, cmd := range cmds {
		switch cmd.Op {
		case "SET":
			if len(cmd.Args) >= 2 {
				key := cmd.Args[0].(string)
				val := cmd.Args[1].(string)
				m.data[key] = []byte(val)
				results[i] = Result{Value: []byte(`"OK"`)}
			}
		case "GET":
			if len(cmd.Args) >= 1 {
				key := cmd.Args[0].(string)
				v, ok := m.data[key]
				if !ok {
					results[i] = Result{Value: []byte("null")}
				} else {
					results[i] = Result{Value: v}
				}
			}
		case "EXPIRE":
			results[i] = Result{Value: []byte("1")}
		}
	}
	return results, nil
}

func (m *memClient) Expire(_ context.Context, _ string, _ time.Duration) error { return nil }

func (m *memClient) Eval(_ context.Context, _ string, _ []string, _ ...any) (any, error) {
	return nil, nil
}

func (m *memClient) Info(_ context.Context) (map[string]string, error) {
	return map[string]string{"db0": "keys=0"}, nil
}

// ── tests ───────────────────────────────────────────────────

func TestFetchThrough_Singleflight(t *testing.T) {
	mc := newMemClient()
	ft := NewFetchThrough[string](mc, nil)

	var fetchCount atomic.Int32
	fetcher := func(ctx context.Context) (*string, error) {
		fetchCount.Add(1)
		time.Sleep(50 * time.Millisecond)
		s := "result"
		return &s, nil
	}

	// launch 10 concurrent fetches for the same key
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ft.Fetch(context.Background(), FetchOpts[string]{
				Key:     "test-key",
				TTL:     time.Minute,
				Fetcher: fetcher,
			})
			if err != nil {
				t.Errorf("fetch error: %v", err)
			}
		}()
	}
	wg.Wait()

	if fetchCount.Load() != 1 {
		t.Fatalf("expected 1 fetch call (singleflight), got %d", fetchCount.Load())
	}
}

func TestFetchThrough_NegativeSentinel(t *testing.T) {
	mc := newMemClient()
	ft := NewFetchThrough[string](mc, nil)

	fetcher := func(ctx context.Context) (*string, error) {
		return nil, nil // nil result → negative cache
	}

	val, err := ft.Fetch(context.Background(), FetchOpts[string]{
		Key:         "neg-key",
		TTL:         time.Minute,
		NegativeTTL: time.Minute,
		Fetcher:     fetcher,
	})
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatal("expected nil for negative cache")
	}

	// second fetch should hit the negative sentinel
	val, err = ft.Fetch(context.Background(), FetchOpts[string]{
		Key:     "neg-key",
		TTL:     time.Minute,
		Fetcher: func(ctx context.Context) (*string, error) { t.Fatal("should not call fetcher"); return nil, nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatal("expected nil from negative sentinel")
	}
}

func TestKeyRegistry_FromSpecs(t *testing.T) {
	specs := []registry.SourceSpec{{
		Domain:           "news",
		Resource:         "digest",
		Version:          1,
		CanonicalKey:     "news:digest:v1:full:en",
		DataTTL:          24 * time.Hour,
		MaxStaleDuration: 30 * time.Minute,
		Schedule:         registry.CronSchedule{Expr: "*/5 * * * *"},
	}}

	kr := NewKeyRegistryFromSpecs(specs)
	entry, ok := kr.Get("news:digest:v1:full:en")
	if !ok {
		t.Fatal("key not found")
	}
	if entry.SeedMetaKey != "seed-meta:news:digest" {
		t.Fatalf("expected seed-meta:news:digest, got %s", entry.SeedMetaKey)
	}
	if entry.DataTTL != 24*time.Hour {
		t.Fatalf("expected 24h TTL, got %v", entry.DataTTL)
	}
}

func TestKeyRegistry_Validate(t *testing.T) {
	specs := []registry.SourceSpec{{
		Domain:       "bad",
		Resource:     "spec",
		Version:      1,
		CanonicalKey: "bad:spec:v1",
		DataTTL:      0, // invalid
		Schedule:     registry.CronSchedule{Expr: "* * * * *"},
	}}

	kr := NewKeyRegistryFromSpecs(specs)
	if err := kr.Validate(); err == nil {
		t.Fatal("expected validation error for TTL<=0")
	}
}

func TestSidecar_LRU(t *testing.T) {
	sc := NewSidecarCache(3, 1024*1024, time.Hour)
	defer sc.Stop()

	sc.Set("a", []byte("1"), time.Minute)
	sc.Set("b", []byte("2"), time.Minute)
	sc.Set("c", []byte("3"), time.Minute)
	sc.Set("d", []byte("4"), time.Minute) // should evict "a"

	if _, ok := sc.Get("a"); ok {
		t.Fatal("expected 'a' to be evicted")
	}
	if _, ok := sc.Get("d"); !ok {
		t.Fatal("expected 'd' to exist")
	}

	stats := sc.Stats()
	if stats.Entries != 3 {
		t.Fatalf("expected 3 entries, got %d", stats.Entries)
	}
}

func TestSidecar_TTLExpiry(t *testing.T) {
	sc := NewSidecarCache(100, 1024*1024, time.Hour)
	defer sc.Stop()

	sc.Set("x", []byte("val"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	if _, ok := sc.Get("x"); ok {
		t.Fatal("expected expired entry to be gone")
	}
}
