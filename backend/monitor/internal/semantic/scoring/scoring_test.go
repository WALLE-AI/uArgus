package scoring

import (
	"math"
	"testing"
)

func TestImportanceScorer_WeightsSum(t *testing.T) {
	s := NewImportanceScorer()
	var sum float64
	for _, f := range s.Profile().Factors {
		sum += f.Weight
	}
	if math.Abs(sum-1.0) > 0.001 {
		t.Fatalf("importance weights sum to %f, expected 1.0", sum)
	}
}

func TestDisruptionScorer_WeightsSum(t *testing.T) {
	s := NewDisruptionScorer()
	var sum float64
	for _, f := range s.Profile().Factors {
		sum += f.Weight
	}
	if math.Abs(sum-1.0) > 0.001 {
		t.Fatalf("disruption weights sum to %f, expected 1.0", sum)
	}
}

func TestWeightedScorer_AllZero(t *testing.T) {
	s := NewImportanceScorer()
	score := s.Score(ScoringInput{})
	if score != 0 {
		t.Fatalf("expected 0, got %f", score)
	}
}

func TestWeightedScorer_AllOne(t *testing.T) {
	s := NewImportanceScorer()
	score := s.Score(ScoringInput{
		"severity":      1,
		"tier":          1,
		"corroboration": 1,
		"recency":       1,
	})
	if score != 1.0 {
		t.Fatalf("expected 1.0, got %f", score)
	}
}

func TestWeightedScorer_Clamp(t *testing.T) {
	s := NewImportanceScorer()
	score := s.Score(ScoringInput{
		"severity":      5,
		"tier":          5,
		"corroboration": 5,
		"recency":       5,
	})
	if score != 1.0 {
		t.Fatalf("expected clamped to 1.0, got %f", score)
	}
}

func TestGoalpostScorer_FullScore(t *testing.T) {
	dims := []GoalpostDimension{
		{Name: "gdp", Weight: 1.0, Direction: 1, GoalMin: 0, GoalMax: 100},
	}
	s := NewGoalpostScorer(dims)
	score := s.Score(GoalpostInput{"gdp": 100})
	if score != 100 {
		t.Fatalf("expected 100, got %f", score)
	}
}

func TestGoalpostScorer_ZeroScore(t *testing.T) {
	dims := []GoalpostDimension{
		{Name: "gdp", Weight: 1.0, Direction: 1, GoalMin: 0, GoalMax: 100},
	}
	s := NewGoalpostScorer(dims)
	score := s.Score(GoalpostInput{"gdp": 0})
	if score != 0 {
		t.Fatalf("expected 0, got %f", score)
	}
}

func TestGoalpostScorer_Reverse(t *testing.T) {
	dims := []GoalpostDimension{
		{Name: "mortality", Weight: 1.0, Direction: -1, GoalMin: 0, GoalMax: 100},
	}
	s := NewGoalpostScorer(dims)
	// low mortality = high resilience
	score := s.Score(GoalpostInput{"mortality": 0})
	if score != 100 {
		t.Fatalf("expected 100 for low mortality, got %f", score)
	}
}
