package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	openai "github.com/sashabaranov/go-openai"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/llm"
)

// ---------------------------------------------------------------------------
// OpenAICompatProvider — implements Provider for OpenAI-compatible APIs
// (GPT, Ollama, vLLM, etc.)
// ---------------------------------------------------------------------------

// OpenAICompatProvider wraps go-openai to implement Provider for any
// OpenAI-compatible endpoint.
type OpenAICompatProvider struct {
	client *openai.Client
	logger *slog.Logger
}

// OpenAICompatConfig holds configuration for creating an OpenAICompatProvider.
type OpenAICompatConfig struct {
	APIKey  string
	BaseURL string // e.g., "http://localhost:11434/v1" for Ollama
	Logger  *slog.Logger
}

// NewOpenAICompatProvider creates a new OpenAI-compatible provider.
func NewOpenAICompatProvider(cfg OpenAICompatConfig) *OpenAICompatProvider {
	config := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	}
	client := openai.NewClientWithConfig(config)
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &OpenAICompatProvider{
		client: client,
		logger: logger,
	}
}

// CallModel implements Provider.CallModel using OpenAI-compatible streaming.
func (p *OpenAICompatProvider) CallModel(ctx context.Context, params CallModelParams) (<-chan llm.StreamEvent, error) {
	ch := make(chan llm.StreamEvent, 128)

	req, err := p.buildRequest(params)
	if err != nil {
		close(ch)
		return nil, fmt.Errorf("building openai request: %w", err)
	}

	go func() {
		defer close(ch)
		p.streamResponse(ctx, req, ch)
	}()

	return ch, nil
}

// buildRequest converts CallModelParams to an OpenAI ChatCompletionRequest.
func (p *OpenAICompatProvider) buildRequest(params CallModelParams) (openai.ChatCompletionRequest, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(params.Messages)+1)

	// System prompt
	if params.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: params.SystemPrompt,
		})
	}

	// Convert messages
	for _, msg := range params.Messages {
		oaiMsg := convertMessageToOpenAI(msg)
		if oaiMsg.Role != "" {
			messages = append(messages, oaiMsg)
		}
	}

	model := params.Model
	if model == "" {
		model = "gpt-4o"
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	if params.MaxOutputTokens != nil {
		req.MaxTokens = *params.MaxOutputTokens
	}

	// Convert tools
	if len(params.Tools) > 0 {
		oaiTools := make([]openai.Tool, 0, len(params.Tools))
		for _, t := range params.Tools {
			oaiTools = append(oaiTools, openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  t.InputSchema,
				},
			})
		}
		req.Tools = oaiTools
	}

	return req, nil
}

// streamResponse handles the OpenAI streaming response.
func (p *OpenAICompatProvider) streamResponse(ctx context.Context, req openai.ChatCompletionRequest, ch chan<- llm.StreamEvent) {
	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		ch <- llm.StreamEvent{
			Type:  llm.EventError,
			Error: err.Error(),
		}
		return
	}
	defer stream.Close()

	ch <- llm.StreamEvent{Type: llm.EventRequestStart}

	var fullContent string
	var toolCalls []openai.ToolCall

	for {
		response, err := stream.Recv()
		if err != nil {
			// Check if stream ended normally
			if err.Error() == "EOF" {
				break
			}
			p.logger.Error("openai stream error", "error", err)
			ch <- llm.StreamEvent{
				Type:  llm.EventError,
				Error: err.Error(),
			}
			return
		}

		for _, choice := range response.Choices {
			delta := choice.Delta

			// Text content
			if delta.Content != "" {
				fullContent += delta.Content
				ch <- llm.StreamEvent{
					Type: llm.EventTextDelta,
					Text: delta.Content,
				}
			}

			// Tool calls
			for _, tc := range delta.ToolCalls {
				idx := 0
				if tc.Index != nil {
					idx = *tc.Index
				}
				for len(toolCalls) <= idx {
					toolCalls = append(toolCalls, openai.ToolCall{})
				}
				if tc.ID != "" {
					toolCalls[idx].ID = tc.ID
				}
				if tc.Function.Name != "" {
					toolCalls[idx].Function.Name = tc.Function.Name
				}
				toolCalls[idx].Function.Arguments += tc.Function.Arguments
			}
		}
	}

	// Build final assistant message
	content := make([]llm.ContentBlock, 0)
	if fullContent != "" {
		content = append(content, llm.ContentBlock{
			Type: llm.ContentBlockText,
			Text: fullContent,
		})
	}

	// Convert accumulated tool calls
	for _, tc := range toolCalls {
		inputBytes := json.RawMessage(tc.Function.Arguments)
		content = append(content, llm.ContentBlock{
			Type:  llm.ContentBlockToolUse,
			ID:    tc.ID,
			Name:  tc.Function.Name,
			Input: inputBytes,
		})
	}

	stopReason := "end_turn"
	if len(toolCalls) > 0 {
		stopReason = "tool_use"
	}

	assistantMsg := llm.Message{
		Type:       llm.MessageTypeAssistant,
		Content:    content,
		Model:      req.Model,
		StopReason: stopReason,
	}

	ch <- llm.StreamEvent{
		Type:    llm.EventAssistant,
		Message: &assistantMsg,
	}
}

// convertMessageToOpenAI converts an internal Message to OpenAI format.
func convertMessageToOpenAI(msg llm.Message) openai.ChatCompletionMessage {
	switch msg.Type {
	case llm.MessageTypeAssistant:
		text := ""
		var toolCalls []openai.ToolCall
		for _, b := range msg.Content {
			switch b.Type {
			case llm.ContentBlockText:
				text += b.Text
			case llm.ContentBlockToolUse:
				toolCalls = append(toolCalls, openai.ToolCall{
					ID:   b.ID,
					Type: openai.ToolTypeFunction,
					Function: openai.FunctionCall{
						Name:      b.Name,
						Arguments: string(b.Input),
					},
				})
			}
		}
		return openai.ChatCompletionMessage{
			Role:      openai.ChatMessageRoleAssistant,
			Content:   text,
			ToolCalls: toolCalls,
		}

	case llm.MessageTypeUser:
		// Check if it's a tool result
		for _, b := range msg.Content {
			if b.Type == llm.ContentBlockToolResult {
				contentStr, _ := b.Content.(string)
				return openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    contentStr,
					ToolCallID: b.ToolUseID,
				}
			}
		}
		// Regular user message
		text := ""
		for _, b := range msg.Content {
			if b.Type == llm.ContentBlockText {
				text += b.Text
			}
		}
		return openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		}

	default:
		return openai.ChatCompletionMessage{}
	}
}
