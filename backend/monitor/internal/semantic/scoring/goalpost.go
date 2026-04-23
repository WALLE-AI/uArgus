package scoring

// GoalpostDimension defines one resilience dimension with direction and goal range.
type GoalpostDimension struct {
	Name      string
	Weight    float64
	Direction int     // +1 = higher is better, -1 = lower is better
	GoalMin   float64 // lower goalpost
	GoalMax   float64 // upper goalpost
}

// NewGoalpostScorer returns a resilience scorer using goalpost normalisation.
// Each dimension is normalised to 0-100 via (value-GoalMin)/(GoalMax-GoalMin) * Direction,
// then weighted and summed.
func NewGoalpostScorer(dims []GoalpostDimension) *GoalpostScorer {
	return &GoalpostScorer{dims: dims}
}

// GoalpostScorer scores resilience across multiple dimensions.
type GoalpostScorer struct {
	dims []GoalpostDimension
}

// GoalpostInput maps dimension names to raw values.
type GoalpostInput map[string]float64

// Score computes a 0-100 resilience score.
func (g *GoalpostScorer) Score(input GoalpostInput) float64 {
	var totalWeight float64
	var weighted float64

	for _, d := range g.dims {
		v, ok := input[d.Name]
		if !ok {
			continue
		}
		rng := d.GoalMax - d.GoalMin
		if rng == 0 {
			continue
		}

		normalised := (v - d.GoalMin) / rng
		if d.Direction < 0 {
			normalised = 1 - normalised
		}
		if normalised < 0 {
			normalised = 0
		}
		if normalised > 1 {
			normalised = 1
		}

		weighted += normalised * d.Weight * 100
		totalWeight += d.Weight
	}

	if totalWeight == 0 {
		return 0
	}
	return weighted / totalWeight
}
