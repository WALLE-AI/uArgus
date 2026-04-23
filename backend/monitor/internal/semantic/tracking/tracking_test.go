package tracking

import "testing"

// ── test helpers ────────────────────────────────────────────

type testItem struct{ hash string }

func (i testItem) TrackHash() string { return i.hash }

// ── dedup tests ─────────────────────────────────────────────

func TestDedup(t *testing.T) {
	items := []Trackable{
		testItem{"aaa"},
		testItem{"bbb"},
		testItem{"aaa"},
		testItem{"ccc"},
	}
	deduped := HashDedup{}.Dedup(items)
	if len(deduped) != 3 {
		t.Fatalf("expected 3 unique items, got %d", len(deduped))
	}
}

func TestHash(t *testing.T) {
	h1 := Hash("hello")
	h2 := Hash("hello")
	h3 := Hash("world")
	if h1 != h2 {
		t.Fatal("same input should produce same hash")
	}
	if h1 == h3 {
		t.Fatal("different input should produce different hash")
	}
	if len(h1) != 32 {
		t.Fatalf("expected 32 hex chars, got %d", len(h1))
	}
}

// ── stage derivation tests ──────────────────────────────────

func TestDeriveStage(t *testing.T) {
	tests := []struct {
		count int
		ageMs int64
		want  string
	}{
		{1, 0, "BREAKING"},
		{2, 3600000, "DEVELOPING"},    // 1h
		{5, 43200000, "SUSTAINED"},    // 12h
		{10, 172800000, "FADING"},     // 48h
	}
	for _, tt := range tests {
		now := int64(1000000000000)
		info := TrackInfo{
			FirstSeen: now - tt.ageMs,
			LastSeen:  now,
			Count:     tt.count,
		}
		got := deriveStage(info)
		if got != tt.want {
			t.Errorf("count=%d age=%dms: want %s, got %s", tt.count, tt.ageMs, tt.want, got)
		}
	}
}

// ── trend derivation tests ──────────────────────────────────

func TestDeriveTrend(t *testing.T) {
	if deriveTrend([]float64{1, 2, 3}) != "improving" {
		t.Fatal("expected improving")
	}
	if deriveTrend([]float64{3, 2, 1}) != "declining" {
		t.Fatal("expected declining")
	}
	if deriveTrend([]float64{1, 3, 2}) != "stable" {
		t.Fatal("expected stable")
	}
	if deriveTrend([]float64{1}) != "stable" {
		t.Fatal("expected stable for short series")
	}
}
