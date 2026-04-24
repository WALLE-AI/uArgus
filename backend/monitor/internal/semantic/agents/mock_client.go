package agents

import "context"

// MockAgentsClient is a test double that returns fixed responses.
type MockAgentsClient struct {
	SummarizeResult string
	ClassifyResult  *AiClassification
	SentimentResult []SentimentResult
	NERResult       [][]NEREntity
	TranslateResult string
	EmbedResult     []float32
}

func (m *MockAgentsClient) Summarize(_ context.Context, _ []string, _ SummarizeOpts) (string, error) {
	return m.SummarizeResult, nil
}

func (m *MockAgentsClient) Classify(_ context.Context, _ string) (*AiClassification, error) {
	return m.ClassifyResult, nil
}

func (m *MockAgentsClient) Sentiment(_ context.Context, texts []string) ([]SentimentResult, error) {
	if m.SentimentResult != nil {
		return m.SentimentResult, nil
	}
	results := make([]SentimentResult, len(texts))
	for i := range results {
		results[i] = SentimentResult{Label: "neutral", Score: 0.5}
	}
	return results, nil
}

func (m *MockAgentsClient) ExtractEntities(_ context.Context, texts []string) ([][]NEREntity, error) {
	if m.NERResult != nil {
		return m.NERResult, nil
	}
	return make([][]NEREntity, len(texts)), nil
}

func (m *MockAgentsClient) Translate(_ context.Context, text string, _ string) (string, error) {
	if m.TranslateResult != "" {
		return m.TranslateResult, nil
	}
	return text, nil
}

func (m *MockAgentsClient) Embed(_ context.Context, _ string) ([]float32, error) {
	return m.EmbedResult, nil
}
