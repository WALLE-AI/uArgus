package scoring

import "math"

// Factor defines a named weight in a scoring profile.
type Factor struct {
	Name   string
	Weight float64
}

// ScoringProfile describes a complete scoring configuration.
type ScoringProfile struct {
	Name    string
	Factors []Factor
	Clamp   [2]float64 // [min, max]
}

// ScoringInput maps factor names to their raw values.
type ScoringInput map[string]float64

// WeightedScorer computes a weighted sum clamped to the profile's range.
type WeightedScorer struct {
	profile ScoringProfile
}

// NewCustomScorer creates a WeightedScorer from a caller-defined profile.
func NewCustomScorer(p ScoringProfile) *WeightedScorer {
	return &WeightedScorer{profile: p}
}

// Score computes the weighted sum of input factors clamped to [min, max].
func (s *WeightedScorer) Score(input ScoringInput) float64 {
	var sum float64
	for _, f := range s.profile.Factors {
		sum += input[f.Name] * f.Weight
	}
	return math.Max(s.profile.Clamp[0], math.Min(s.profile.Clamp[1], sum))
}

// Profile returns the scorer's profile (for inspection/testing).
func (s *WeightedScorer) Profile() ScoringProfile {
	return s.profile
}
