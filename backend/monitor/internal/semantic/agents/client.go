package agents

import "context"

// SummarizeOpts configures a summarisation request.
type SummarizeOpts struct {
	Mode      string // "brief" | "detailed"
	MaxTokens int
}

// AiClassification is the result of an AI classification request.
type AiClassification struct {
	Categories []string `json:"categories"`
	Confidence float64  `json:"confidence"`
}

// AgentsClient is the interface for LLM/Embedding backend calls.
type AgentsClient interface {
	Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error)
	Classify(ctx context.Context, text string) (*AiClassification, error)
	Embed(ctx context.Context, text string) ([]float32, error)
}
