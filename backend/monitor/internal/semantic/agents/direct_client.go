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
	embedder *EmbeddingClient
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
// Prompts are mode-aware (brief/analysis/translate/default) and variant-aware (tech/full).
func (c *DirectAgentsClient) Summarize(ctx context.Context, texts []string, opts SummarizeOpts) (string, error) {
	maxTokens := opts.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1024
	}

	prompts := BuildSummarizePrompts(texts, opts)

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: prompts.SystemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: prompts.UserPrompt},
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
	systemPrompt, userPrompt := BuildClassifyPrompt(text)

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: userPrompt},
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

	responseText = stripCodeFences(responseText)

	var result AiClassification
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		c.logger.Warn("classify: failed to parse LLM response", "response", responseText, "err", err)
		return &AiClassification{Categories: []string{"unknown"}, Confidence: 0.0}, nil
	}
	return &result, nil
}

// Sentiment performs batch sentiment analysis via LLM.
func (c *DirectAgentsClient) Sentiment(ctx context.Context, texts []string) ([]SentimentResult, error) {
	maxTokens := 1024
	systemPrompt, userPrompt := BuildSentimentPrompt(texts)

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: userPrompt},
				},
			},
		},
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("direct agents: sentiment: %w", err)
	}

	responseText, err := collectText(ch)
	if err != nil {
		return nil, err
	}

	responseText = stripCodeFences(responseText)

	var results []SentimentResult
	if err := json.Unmarshal([]byte(responseText), &results); err != nil {
		c.logger.Warn("sentiment: failed to parse LLM response", "response", responseText, "err", err)
		// Return neutral fallback for all inputs
		fallback := make([]SentimentResult, len(texts))
		for i := range fallback {
			fallback[i] = SentimentResult{Label: "neutral", Score: 0.5}
		}
		return fallback, nil
	}
	return results, nil
}

// ExtractEntities performs batch NER via LLM.
func (c *DirectAgentsClient) ExtractEntities(ctx context.Context, texts []string) ([][]NEREntity, error) {
	maxTokens := 2048
	systemPrompt, userPrompt := BuildNERPrompt(texts)

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: userPrompt},
				},
			},
		},
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("direct agents: ner: %w", err)
	}

	responseText, err := collectText(ch)
	if err != nil {
		return nil, err
	}

	responseText = stripCodeFences(responseText)

	var results [][]NEREntity
	if err := json.Unmarshal([]byte(responseText), &results); err != nil {
		c.logger.Warn("ner: failed to parse LLM response", "response", responseText, "err", err)
		return make([][]NEREntity, len(texts)), nil
	}
	return results, nil
}

// Translate translates text to targetLang via LLM.
func (c *DirectAgentsClient) Translate(ctx context.Context, text string, targetLang string) (string, error) {
	maxTokens := 1024
	prompts := BuildSummarizePrompts([]string{text}, SummarizeOpts{
		Mode:    "translate",
		Variant: targetLang,
	})

	ch, err := c.provider.CallModel(ctx, provider.CallModelParams{
		Model:        c.model,
		SystemPrompt: prompts.SystemPrompt,
		Messages: []llm.Message{
			{
				Type: llm.MessageTypeUser,
				Content: []llm.ContentBlock{
					{Type: llm.ContentBlockText, Text: prompts.UserPrompt},
				},
			},
		},
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("direct agents: translate: %w", err)
	}

	return collectText(ch)
}

// Embed generates a vector embedding for the given text.
// Requires an EmbeddingClient to be set; returns an error if not configured.
func (c *DirectAgentsClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if c.embedder == nil {
		return nil, fmt.Errorf("direct agents: embedding not configured; set an EmbeddingClient")
	}
	return c.embedder.Embed(ctx, text)
}

// SetEmbedder injects an EmbeddingClient for Embed() calls.
func (c *DirectAgentsClient) SetEmbedder(e *EmbeddingClient) {
	c.embedder = e
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

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

// stripCodeFences removes markdown code fences from LLM JSON responses.
func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
