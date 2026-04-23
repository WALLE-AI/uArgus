package scoring

// NewDisruptionScorer returns the chokepoint disruption scoring profile.
// threatLevel*0.4 + warningCount*0.3 + severity*0.2 + anomaly*0.1
func NewDisruptionScorer() *WeightedScorer {
	return NewCustomScorer(ScoringProfile{
		Name: "disruption",
		Factors: []Factor{
			{Name: "threatLevel", Weight: 0.40},
			{Name: "warningCount", Weight: 0.30},
			{Name: "severity", Weight: 0.20},
			{Name: "anomaly", Weight: 0.10},
		},
		Clamp: [2]float64{0, 1},
	})
}
