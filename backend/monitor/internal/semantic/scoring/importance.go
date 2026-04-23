package scoring

// NewImportanceScorer returns the news importance scoring profile.
// severity*0.55 + tier*0.2 + corroboration*0.15 + recency*0.1
func NewImportanceScorer() *WeightedScorer {
	return NewCustomScorer(ScoringProfile{
		Name: "importance",
		Factors: []Factor{
			{Name: "severity", Weight: 0.55},
			{Name: "tier", Weight: 0.20},
			{Name: "corroboration", Weight: 0.15},
			{Name: "recency", Weight: 0.10},
		},
		Clamp: [2]float64{0, 1},
	})
}
