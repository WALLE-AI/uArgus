package fusion

import (
	"context"
	"testing"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

type memClient struct{ data map[string][]byte }

func newMem() *memClient                                                        { return &memClient{data: make(map[string][]byte)} }
func (m *memClient) Get(_ context.Context, k string) ([]byte, error)            { return m.data[k], nil }
func (m *memClient) Set(_ context.Context, k string, v []byte, _ time.Duration) error { m.data[k] = v; return nil }
func (m *memClient) Del(_ context.Context, _ ...string) error                   { return nil }
func (m *memClient) Pipeline(_ context.Context, _ []cache.Cmd) ([]cache.Result, error) { return nil, nil }
func (m *memClient) Expire(_ context.Context, _ string, _ time.Duration) error  { return nil }
func (m *memClient) Eval(_ context.Context, _ string, _ []string, _ ...any) (any, error) { return nil, nil }

func TestCorrelationFuser_TwoDomains(t *testing.T) {
	mc := newMem()
	mc.data["news:digest:v1"] = []byte(`{"items":[1]}`)
	mc.data["market:quotes:v1"] = []byte(`{"items":[2]}`)

	f := NewCorrelationFuser(mc)
	cards, err := f.Fuse(context.Background(), []FuseInput{
		{Domain: "news", Key: "news:digest:v1"},
		{Domain: "market", Key: "market:quotes:v1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) == 0 {
		t.Fatal("expected at least 1 correlation card")
	}
}

func TestCrossSourceFuser(t *testing.T) {
	mc := newMem()
	mc.data["key1"] = []byte(`{"a":1}`)
	mc.data["key2"] = []byte(`{"b":2}`)

	f := NewCrossSourceFuser(mc)
	signals, err := f.Fuse(context.Background(), []string{"key1", "key2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(signals) == 0 {
		t.Fatal("expected at least 1 cross-source signal")
	}
}
