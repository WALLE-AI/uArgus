package preprocess

import (
	"context"
	"testing"
)

// ── test stages ─────────────────────────────────────────────

type addFieldStage struct{ key, val string }

func (s addFieldStage) Name() string { return "add-field" }
func (s addFieldStage) Process(_ context.Context, data any) (any, error) {
	m := data.(map[string]any)
	m[s.key] = s.val
	return m, nil
}

// ── tests ───────────────────────────────────────────────────

func TestPipeline_SequentialExecution(t *testing.T) {
	p := NewPipeline(
		addFieldStage{key: "a", val: "1"},
		addFieldStage{key: "b", val: "2"},
	)
	out, err := p.Run(context.Background(), map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	m := out.(map[string]any)
	if m["a"] != "1" || m["b"] != "2" {
		t.Fatalf("expected a=1,b=2, got %v", m)
	}
}

func TestFormatMapper(t *testing.T) {
	fm := NewFormatMapper("test", []FieldMapping{
		{From: "old_name", To: "new_name"},
	})
	input := map[string]any{"old_name": "value"}
	out, err := fm.Process(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	m := out.(map[string]any)
	if m["new_name"] != "value" {
		t.Fatalf("expected renamed field, got %v", m)
	}
	if _, ok := m["old_name"]; ok {
		t.Fatal("old field should be removed")
	}
}

func TestNormalizeTimestamp(t *testing.T) {
	result := NormalizeTimestamp("Mon, 02 Jan 2006 15:04:05 MST")
	s, ok := result.(string)
	if !ok {
		t.Fatal("expected string result")
	}
	if s != "2006-01-02T15:04:05Z" {
		t.Fatalf("expected RFC3339, got %s", s)
	}
}

func TestValidateRequired(t *testing.T) {
	v := NewValidateRequired("title", "link")
	_, err := v.Process(context.Background(), map[string]any{"title": "x"})
	if err == nil {
		t.Fatal("expected error for missing 'link'")
	}
	_, err = v.Process(context.Background(), map[string]any{"title": "x", "link": "y"})
	if err != nil {
		t.Fatal(err)
	}
}
