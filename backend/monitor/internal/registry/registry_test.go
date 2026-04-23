package registry

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// ── mock source ─────────────────────────────────────────────

type mockSource struct {
	name string
	spec SourceSpec
	deps []string
	run  func(ctx context.Context) (*FetchResult, error)
}

func (m *mockSource) Name() string                                       { return m.name }
func (m *mockSource) Spec() SourceSpec                                   { return m.spec }
func (m *mockSource) Dependencies() []string                             { return m.deps }
func (m *mockSource) Run(ctx context.Context) (*FetchResult, error) {
	if m.run != nil {
		return m.run(ctx)
	}
	return &FetchResult{Metrics: FetchMetrics{RecordCount: 1, Duration: time.Millisecond}}, nil
}

func newMock(name string, deps ...string) *mockSource {
	return &mockSource{
		name: name,
		spec: SourceSpec{
			Domain:           "test",
			Resource:         name,
			Version:          1,
			Schedule:         IntervalSchedule{Every: time.Hour}, // won't actually fire in test
			DataTTL:          time.Hour,
			MaxStaleDuration: 10 * time.Minute,
		},
		deps: deps,
	}
}

// ── tests ───────────────────────────────────────────────────

func TestRegister_DuplicateName(t *testing.T) {
	r := New()
	r.Register(newMock("a"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate name")
		}
	}()
	r.Register(newMock("a"))
}

func TestBoot_DependencyOrder(t *testing.T) {
	r := New()
	// B depends on A
	r.Register(newMock("a"))
	r.Register(newMock("b", "a"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := r.Boot(ctx); err != nil {
		t.Fatal(err)
	}
	defer r.Shutdown(ctx)

	if len(r.order) != 2 || r.order[0] != "a" || r.order[1] != "b" {
		t.Fatalf("expected [a, b], got %v", r.order)
	}
}

func TestBoot_CycleDetection(t *testing.T) {
	r := New()
	r.Register(newMock("a", "b"))
	r.Register(newMock("b", "a"))

	ctx := context.Background()
	if err := r.Boot(ctx); err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestHealth_StateTransitions(t *testing.T) {
	ht := NewHealthTracker()
	src := newMock("x")
	ht.Register(src)

	// initially healthy
	snap := ht.Snapshot()
	if snap["x"].State != HealthStateHealthy {
		t.Fatalf("expected healthy, got %s", snap["x"].State)
	}

	// 3 failures → failing
	ht.RecordFailure("x")
	ht.RecordFailure("x")
	ht.RecordFailure("x")
	snap = ht.Snapshot()
	if snap["x"].State != HealthStateFailing {
		t.Fatalf("expected failing, got %s", snap["x"].State)
	}

	// success resets
	ht.RecordSuccess("x", FetchMetrics{RecordCount: 10, Duration: time.Millisecond})
	snap = ht.Snapshot()
	if snap["x"].State != HealthStateHealthy {
		t.Fatalf("expected healthy after success, got %s", snap["x"].State)
	}
}

func TestDispatch_RecordsMetrics(t *testing.T) {
	var callCount atomic.Int32
	r := New()
	src := newMock("d")
	src.spec.Schedule = OnDemandSchedule{}
	src.run = func(ctx context.Context) (*FetchResult, error) {
		callCount.Add(1)
		return &FetchResult{Metrics: FetchMetrics{RecordCount: 5, Duration: 10 * time.Millisecond}}, nil
	}
	r.Register(src)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := r.Boot(ctx); err != nil {
		t.Fatal(err)
	}

	// manually dispatch
	r.dispatchRun(ctx, src)
	if callCount.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", callCount.Load())
	}
	snap := r.HealthSnapshot()
	if snap["d"].LastRecordCount != 5 {
		t.Fatalf("expected 5 records, got %d", snap["d"].LastRecordCount)
	}
}

func TestDispatch_FailureRecorded(t *testing.T) {
	r := New()
	src := newMock("f")
	src.spec.Schedule = OnDemandSchedule{}
	src.run = func(ctx context.Context) (*FetchResult, error) {
		return nil, errors.New("boom")
	}
	r.Register(src)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = r.Boot(ctx)

	r.dispatchRun(ctx, src)
	snap := r.HealthSnapshot()
	if snap["f"].ConsecutiveFails != 1 {
		t.Fatalf("expected 1 fail, got %d", snap["f"].ConsecutiveFails)
	}
}
