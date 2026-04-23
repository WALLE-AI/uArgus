package classify

import "testing"

func TestClassify_Earthquake(t *testing.T) {
	c := NewNewsClassifier()
	r := c.Classify("Major earthquake hits Tokyo, tsunami warning issued")
	if r.Severity < 6 {
		t.Fatalf("expected severity ≥6, got %d", r.Severity)
	}
	if len(r.Categories) == 0 {
		t.Fatal("expected categories")
	}
	found := false
	for _, cat := range r.Categories {
		if cat == "disaster" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'disaster' category, got %v", r.Categories)
	}
}

func TestClassify_NoMatch(t *testing.T) {
	c := NewNewsClassifier()
	r := c.Classify("Local bakery opens new branch")
	if len(r.Categories) != 0 {
		t.Fatalf("expected no categories, got %v", r.Categories)
	}
}

func TestClassify_Conflict(t *testing.T) {
	c := NewNewsClassifier()
	r := c.Classify("Missile strikes reported in conflict zone, troops deployed")
	if r.Severity != 7 {
		t.Fatalf("expected severity 7, got %d", r.Severity)
	}
}

func TestCustomClassifier(t *testing.T) {
	levels := []KeywordLevel{
		{Name: "custom", Severity: 5, Keywords: []string{"foobar"}},
	}
	c := NewCustomClassifier(levels)
	r := c.Classify("Something about foobar happening")
	if len(r.Categories) != 1 || r.Categories[0] != "custom" {
		t.Fatalf("expected [custom], got %v", r.Categories)
	}
}
