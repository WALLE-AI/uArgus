package preprocess

import (
	"context"
	"fmt"
)

// Stage is one step in a preprocessing pipeline.
type Stage interface {
	Name() string
	Process(ctx context.Context, data any) (any, error)
}

// Pipeline executes a sequence of Stages.
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates a Pipeline from the given stages.
func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run executes all stages in order. Each stage receives the previous stage's output.
func (p *Pipeline) Run(ctx context.Context, input any) (any, error) {
	current := input
	for _, s := range p.stages {
		out, err := s.Process(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("preprocess %s: %w", s.Name(), err)
		}
		current = out
	}
	return current, nil
}

// Append adds more stages to the pipeline.
func (p *Pipeline) Append(stages ...Stage) {
	p.stages = append(p.stages, stages...)
}
