package agents

import "context"

// SummarizeOpts configures a summarisation request.
type SummarizeOpts struct {
	Mode       string // "brief" | "analysis" | "translate" | "" (default)
	Variant    string // "full" | "tech"
	Lang       string // ISO language code, e.g. "en", "fr", "zh"
	GeoContext string // optional geopolitical intelligence context
	MaxTokens  int
}

// AiClassification is the result of an AI classification request.
type AiClassification struct {
	Categories []string `json:"categories"`
	Confidence float64  `json:"confidence"`
}

// SentimentResult holds the sentiment analysis result for a single text.
type SentimentResult struct {
	Label string  `json:"label"` // "positive" | "negative" | "neutral"
	Score float64 `json:"score"` // 0.0 - 1.0
}

// NEREntity represents a named entity extracted from text.
type NEREntity struct {
	Text       string  `json:"text"`
	Type       string  `json:"type"`       // PER, ORG, LOC, MISC
	Confidence float64 `json:"confidence"` // 0.0 - 1.0
}

// AgentsClient is the interface for LLM/Embedding backend calls.
type AgentsClient interface {
	Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error)
	Classify(ctx context.Context, text string) (*AiClassification, error)
	Sentiment(ctx context.Context, texts []string) ([]SentimentResult, error)
	ExtractEntities(ctx context.Context, texts []string) ([][]NEREntity, error)
	Translate(ctx context.Context, text string, targetLang string) (string, error)
	Embed(ctx context.Context, text string) ([]float32, error)
}
