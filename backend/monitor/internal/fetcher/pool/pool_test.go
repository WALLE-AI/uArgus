package pool

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

func TestBoundedPool_AllSettled(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	results := BoundedPool(context.Background(), 2, items, func(_ context.Context, n int) (int, error) {
		return n * 10, nil
	})
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if r.Value != items[r.Index]*10 {
			t.Fatalf("result[%d] = %d, want %d", r.Index, r.Value, items[r.Index]*10)
		}
	}
}

func TestBoundedPool_ConcurrencyLimit(t *testing.T) {
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	items := make([]int, 20)
	_ = BoundedPool(context.Background(), 5, items, func(_ context.Context, _ int) (int, error) {
		cur := concurrent.Add(1)
		for {
			old := maxConcurrent.Load()
			if cur <= old || maxConcurrent.CompareAndSwap(old, cur) {
				break
			}
		}
		concurrent.Add(-1)
		return 0, nil
	})

	if maxConcurrent.Load() > 5 {
		t.Fatalf("max concurrent %d exceeded limit 5", maxConcurrent.Load())
	}
}

func TestBoundedPool_PartialFailure(t *testing.T) {
	items := []int{1, 2, 3}
	results := BoundedPool(context.Background(), 3, items, func(_ context.Context, n int) (int, error) {
		if n == 2 {
			return 0, fmt.Errorf("fail")
		}
		return n, nil
	})
	if results[1].Err == nil {
		t.Fatal("expected error for item 2")
	}
	if results[0].Err != nil || results[2].Err != nil {
		t.Fatal("expected items 1 and 3 to succeed")
	}
}

func TestFanOut_Success(t *testing.T) {
	err := FanOut(context.Background(),
		func(_ context.Context) error { return nil },
		func(_ context.Context) error { return nil },
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestFanOut_Error(t *testing.T) {
	err := FanOut(context.Background(),
		func(_ context.Context) error { return nil },
		func(_ context.Context) error { return fmt.Errorf("boom") },
	)
	if err == nil {
		t.Fatal("expected error")
	}
}
