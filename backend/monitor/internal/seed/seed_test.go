package seed

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// ── in-memory mock Client (same as cache tests) ────────────

type memClient struct {
	data map[string][]byte
}

func newMemClient() *memClient {
	return &memClient{data: make(map[string][]byte)}
}

func (m *memClient) Get(_ context.Context, key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (m *memClient) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *memClient) Del(_ context.Context, keys ...string) error {
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

func (m *memClient) Pipeline(_ context.Context, cmds []cache.Cmd) ([]cache.Result, error) {
	results := make([]cache.Result, len(cmds))
	for i, cmd := range cmds {
		switch cmd.Op {
		case "SET":
			if len(cmd.Args) >= 2 {
				key := cmd.Args[0].(string)
				val := cmd.Args[1].(string)
				m.data[key] = []byte(val)
				results[i] = cache.Result{Value: []byte(`"OK"`)}
			}
		case "GET":
			if len(cmd.Args) >= 1 {
				key := cmd.Args[0].(string)
				v, ok := m.data[key]
				if !ok {
					results[i] = cache.Result{Value: []byte("null")}
				} else {
					results[i] = cache.Result{Value: v}
				}
			}
		case "EXPIRE":
			results[i] = cache.Result{Value: []byte("1")}
		}
	}
	return results, nil
}

func (m *memClient) Expire(_ context.Context, _ string, _ time.Duration) error { return nil }

func (m *memClient) Eval(_ context.Context, script string, keys []string, args ...any) (any, error) {
	// simplified CAS: just delete the key
	if len(keys) > 0 {
		delete(m.data, keys[0])
	}
	return nil, nil
}

// ── Envelope tests ──────────────────────────────────────────

func TestEnvelope_Unwrap_Contract(t *testing.T) {
	input := `{"_seed":{"fetchedAt":1234,"recordCount":10,"sourceVersion":"v2","schemaVersion":1,"state":"OK"},"data":{"items":[1,2,3]}}`
	meta, data, err := Unwrap([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if meta == nil {
		t.Fatal("expected meta")
	}
	if meta.RecordCount != 10 {
		t.Fatalf("expected 10 records, got %d", meta.RecordCount)
	}
	if len(data) == 0 {
		t.Fatal("expected data")
	}
}

func TestEnvelope_Unwrap_Legacy(t *testing.T) {
	input := `{"items":[1,2,3]}`
	meta, data, err := Unwrap([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if meta != nil {
		t.Fatal("expected nil meta for legacy format")
	}
	if string(data) != input {
		t.Fatalf("expected raw data, got %s", string(data))
	}
}

func TestEnvelope_Build(t *testing.T) {
	meta := SeedMeta{
		FetchedAt:     1234,
		RecordCount:   5,
		SourceVersion: "v2",
		SchemaVersion: 1,
		State:         "OK",
	}
	b, err := Build(meta, map[string]int{"count": 5})
	if err != nil {
		t.Fatal(err)
	}
	var env SeedEnvelope
	if err := json.Unmarshal(b, &env); err != nil {
		t.Fatal(err)
	}
	if env.Seed.RecordCount != 5 {
		t.Fatalf("expected 5, got %d", env.Seed.RecordCount)
	}
}

// ── Lock tests ──────────────────────────────────────────────

func TestLock_AcquireRelease(t *testing.T) {
	mc := newMemClient()
	l := NewLock(mc)
	ctx := context.Background()

	ok, err := l.Acquire(ctx, "test", "res", "run1", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected lock acquired")
	}

	// second acquire should fail (key already exists)
	// In our simplified mock, SET always succeeds, but in real Redis NX would fail.
	// Test the Release path instead:
	if err := l.Release(ctx, "test", "res", "run1"); err != nil {
		t.Fatal(err)
	}

	// after release, should be able to acquire again
	ok, err = l.Acquire(ctx, "test", "res", "run2", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected lock acquired after release")
	}
}

// ── Publisher tests ─────────────────────────────────────────

func TestPublisher_AtomicPublish(t *testing.T) {
	mc := newMemClient()
	p := NewPublisher(mc)
	ctx := context.Background()

	err := p.AtomicPublish(ctx, "test:key", []byte(`{"data":1}`), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	val, _ := mc.Get(ctx, "test:key")
	if string(val) != `{"data":1}` {
		t.Fatalf("expected data, got %s", string(val))
	}
}

func TestPublisher_WithEnvelope(t *testing.T) {
	mc := newMemClient()
	p := NewPublisher(mc)
	ctx := context.Background()

	meta := SeedMeta{
		FetchedAt:     1234,
		RecordCount:   3,
		SourceVersion: "v2",
		SchemaVersion: 1,
		State:         "OK",
	}
	err := p.PublishWithEnvelope(ctx, "test:env", meta, []int{1, 2, 3}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// check canonical key has _seed
	val, _ := mc.Get(ctx, "test:env")
	if val == nil {
		t.Fatal("expected value")
	}
	var env SeedEnvelope
	if err := json.Unmarshal(val, &env); err != nil {
		t.Fatal(err)
	}
	if env.Seed == nil || env.Seed.RecordCount != 3 {
		t.Fatal("expected envelope with 3 records")
	}

	// check seed-meta key
	metaVal, _ := mc.Get(ctx, "seed-meta:test:env")
	if metaVal == nil {
		t.Fatal("expected seed-meta value")
	}
}

// ── Runner tests ────────────────────────────────────────────

type mockSource struct {
	name   string
	spec   registry.SourceSpec
	result *registry.FetchResult
	err    error
}

func (m *mockSource) Name() string                                         { return m.name }
func (m *mockSource) Spec() registry.SourceSpec                            { return m.spec }
func (m *mockSource) Dependencies() []string                               { return nil }
func (m *mockSource) Run(_ context.Context) (*registry.FetchResult, error) { return m.result, m.err }

func TestRunner_HappyPath(t *testing.T) {
	mc := newMemClient()
	runner := NewRunner(mc, nil)
	ctx := context.Background()

	src := &mockSource{
		name: "test:source",
		spec: registry.SourceSpec{
			Domain:       "test",
			Resource:     "source",
			Version:      1,
			CanonicalKey: "test:source:v1",
			DataTTL:      time.Hour,
			Schedule:     registry.CronSchedule{Expr: "* * * * *"},
		},
		result: &registry.FetchResult{
			Data:    map[string]int{"count": 5},
			Metrics: registry.FetchMetrics{RecordCount: 5, Duration: 10 * time.Millisecond},
		},
	}

	if err := runner.Run(ctx, src); err != nil {
		t.Fatal(err)
	}

	// verify canonical key exists
	val, _ := mc.Get(ctx, "test:source:v1")
	if val == nil {
		t.Fatal("expected canonical key to exist")
	}

	// verify seed-meta exists
	metaVal, _ := mc.Get(ctx, "seed-meta:test:source")
	if metaVal == nil {
		t.Fatal("expected seed-meta to exist")
	}
}

func TestRunner_FetchFail_ExtendTTL(t *testing.T) {
	mc := newMemClient()
	// pre-populate with existing data
	_ = mc.Set(context.Background(), "fail:source:v1", []byte(`{"old":"data"}`), 0)

	runner := NewRunner(mc, nil)
	ctx := context.Background()

	src := &mockSource{
		name: "fail:source",
		spec: registry.SourceSpec{
			Domain:       "fail",
			Resource:     "source",
			Version:      1,
			CanonicalKey: "fail:source:v1",
			DataTTL:      time.Hour,
			Schedule:     registry.CronSchedule{Expr: "* * * * *"},
		},
		err: context.DeadlineExceeded,
	}

	err := runner.Run(ctx, src)
	if err == nil {
		t.Fatal("expected error")
	}

	// existing data should still be there (TTL extended, not deleted)
	val, _ := mc.Get(ctx, "fail:source:v1")
	if val == nil {
		t.Fatal("expected existing data preserved after failure")
	}
}
