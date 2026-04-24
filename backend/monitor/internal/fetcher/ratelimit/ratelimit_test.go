package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestFixedInterval_Pacing(t *testing.T) {
	rl := NewFixedInterval(50 * time.Millisecond)

	start := time.Now()
	for i := 0; i < 3; i++ {
		if err := rl.Wait(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	elapsed := time.Since(start)
	// 3 calls with 50ms interval → at least 100ms total (first call is immediate)
	if elapsed < 90*time.Millisecond {
		t.Fatalf("expected ≥90ms, got %v", elapsed)
	}
}

func TestExponentialBackoff(t *testing.T) {
	cases := []struct {
		n    int
		want time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{10, 30 * time.Second}, // capped
	}
	for _, c := range cases {
		got := ExponentialBackoff(1*time.Second, c.n, 30*time.Second)
		if got != c.want {
			t.Errorf("Exponential(n=%d) = %v, want %v", c.n, got, c.want)
		}
	}
}

func TestLinearBackoff(t *testing.T) {
	got := LinearBackoff(1*time.Second, 2, 10*time.Second)
	if got != 3*time.Second {
		t.Fatalf("expected 3s, got %v", got)
	}
	got = LinearBackoff(1*time.Second, 20, 10*time.Second)
	if got != 10*time.Second {
		t.Fatalf("expected capped at 10s, got %v", got)
	}
}
