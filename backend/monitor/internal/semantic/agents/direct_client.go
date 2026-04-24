package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm/provider"
)

// DirectAgentsClient implements AgentsClient by calling an LLM provider directly,
// replacing the HTTP-based HTTPAgentsClient for environments where the agents
// service is not deployed separately.
type DirectAgentsClient struct {
	provider provider.Provider
	model    string
	logger   *slog.Logger
}

// DirectAgentsConfig configures the DirectAgentsClient.
type DirectAgentsConfig struct {
	ProviderType provider.ProviderType
	APIKey       string
	BaseURL      string
	Model        string // e.g., "claude-sonnet-4-20250514", "gpt-4o"
	Logger       *slog.Logger
}

// NewDirectAgentsClient creates a DirectAgentsClient backed by the given LLM provider.
func NewDirectAgentsClient(cfg DirectAgentsConfig) (*DirectAgentsClient, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	p, err := provider.NewProvider(provider.FactoryConfig{
		Type:    cfg.ProviderType,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("direct agents: create provider: %w", err)
	}

	model := cfg.Model
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	return &DirectAgentsClient{
		provider: p,
		model:    model,
		logger:   logger,
	}, nil
}

// Summarize sends headlines to the LLM and returns a summary.
func (c *DirectAgentsClient) Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error) {
	combined := strings.Join(texts, "\n- ")
	maxTokens := opts.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1024
	}
	systemPrompt := "You are a news analyst. Produce a concise summary of the given headlines. Focus on the most important events and their implications."
	if opts.Mode == "brief" {
		systemPrompt += " Keep it under 3 sentences."
	}

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: "Summarize these headlines:\n- " + combined},
				},
			},
		},
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("direct agents: summarize: %w", err)
	}

	return collectText(ch)
}

// Classify sends a text to the LLM for AI classification.
func (c *DirectAgentsClient) Classify(ctx context.Context, text string) (*AiClassification, error) {
	maxTokens := 512
	systemPrompt := `You are a news classifier. Given a headline, return a JSON object with:
- "categories": array of category strings (e.g., "geopolitics", "tech", "finance", "climate", "health", "security")
- "confidence": a float between 0 and 1
Only return valid JSON, no explanation.`

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: text},
				},
			},
		},
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("direct agents: classify: %w", err)
	}

	responseText, err := collectText(ch)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	responseText = strings.TrimSpace(responseText)
	// Strip markdown code fences if present
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var result AiClassification
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		c.logger.Warn("classify: failed to parse LLM response", "response", responseText, "err", err)
		return &AiClassification{Categories: []string{"unknown"}, Confidence: 0.0}, nil
	}
	return &result, nil
}

// Embed is a stub — embedding requires a separate model/endpoint, not supported via chat LLM.
func (c *DirectAgentsClient) Embed(ctx context.Context, text string) ([]float32, error) {
	return nil, fmt.Errorf("direct agents: embedding not supported via chat LLM provider; use a dedicated embedding endpoint")
}

// collectText drains a StreamEvent channel and concatenates all text deltas.
func collectText(ch <-chan llm.StreamEvent) (string, error) {
	var sb strings.Builder
	var lastError string

	for event := range ch {
		switch event.Type {
		case llm.EventTextDelta:
			sb.WriteString(event.Text)
		case llm.EventAssistant:
			if event.Message != nil {
				for _, block := range event.Message.Content {
					if block.Type == llm.ContentBlockText && sb.Len() == 0 {
						sb.WriteString(block.Text)
					}
				}
			}
		case llm.EventError:
			lastError = event.Error
		}
	}

	if sb.Len() == 0 && lastError != "" {
		return "", fmt.Errorf("llm error: %s", lastError)
	}
	return sb.String(), nil
}
