package agents

import "context"

// MockAgentsClient is a test double that returns fixed responses.
type MockAgentsClient struct {
	SummarizeResult string
	ClassifyResult  *AiClassification
	EmbedResult     []float32
}

func (m *MockAgentsClient) Summarize(_ context.Context, _ []string, _ SummarizeOpts) (string, error) {
	return m.SummarizeResult, nil
}

func (m *MockAgentsClient) Classify(_ context.Context, _ string) (*AiClassification, error) {
	return m.ClassifyResult, nil
}

func (m *MockAgentsClient) Embed(_ context.Context, _ string) ([]float32, error) {
	return m.EmbedResult, nil
}
