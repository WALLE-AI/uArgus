package fallback

import (
	"context"
	"fmt"
	"testing"
)

func TestChain_FirstSuccess(t *testing.T) {
	chain := NewChain(
		Provider[string]{Name: "a", Fn: func(_ context.Context) (string, error) { return "a-ok", nil }},
		Provider[string]{Name: "b", Fn: func(_ context.Context) (string, error) { return "b-ok", nil }},
	)
	val, err := chain.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if val != "a-ok" {
		t.Fatalf("expected a-ok, got %s", val)
	}
}

func TestChain_Fallback(t *testing.T) {
	chain := NewChain(
		Provider[string]{Name: "fail", Fn: func(_ context.Context) (string, error) { return "", fmt.Errorf("fail") }},
		Provider[string]{Name: "ok", Fn: func(_ context.Context) (string, error) { return "recovered", nil }},
	)
	val, err := chain.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if val != "recovered" {
		t.Fatalf("expected recovered, got %s", val)
	}
}

func TestChain_AllFail(t *testing.T) {
	chain := NewChain(
		Provider[string]{Name: "a", Fn: func(_ context.Context) (string, error) { return "", fmt.Errorf("a fail") }},
		Provider[string]{Name: "b", Fn: func(_ context.Context) (string, error) { return "", fmt.Errorf("b fail") }},
	)
	_, err := chain.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error when all providers fail")
	}
}
