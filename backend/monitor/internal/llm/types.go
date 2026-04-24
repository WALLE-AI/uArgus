package llm

import (
	"encoding/json"
	"time"
)

// ---------------------------------------------------------------------------
// Message types — shared by all LLM providers
// ---------------------------------------------------------------------------

// MessageType identifies the role of a message.
type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
)

// ContentBlockType identifies the kind of content block.
type ContentBlockType string

const (
	ContentBlockText       ContentBlockType = "text"
	ContentBlockThinking   ContentBlockType = "thinking"
	ContentBlockToolUse    ContentBlockType = "tool_use"
	ContentBlockToolResult ContentBlockType = "tool_result"
)

// ContentBlock is one piece of content within a Message.
type ContentBlock struct {
	Type      ContentBlockType `json:"type"`
	Text      string           `json:"text,omitempty"`
	Thinking  string           `json:"thinking,omitempty"`
	Signature string           `json:"signature,omitempty"`
	ID        string           `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	Input     json.RawMessage  `json:"input,omitempty"`
	ToolUseID string           `json:"tool_use_id,omitempty"`
	Content   interface{}      `json:"content,omitempty"`
	IsError   bool             `json:"is_error,omitempty"`
}

// Message represents one turn in a conversation.
type Message struct {
	Type       MessageType    `json:"type"`
	Content    []ContentBlock `json:"content"`
	Model      string         `json:"model,omitempty"`
	StopReason string         `json:"stop_reason,omitempty"`
	Usage      *Usage         `json:"usage,omitempty"`
	IsApiError bool           `json:"is_api_error,omitempty"`
	Timestamp  time.Time      `json:"timestamp,omitempty"`
}

// Usage tracks token consumption.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ThinkingConfig controls extended thinking behavior.
type ThinkingConfig struct {
	Type         string `json:"type"`          // "enabled" | "disabled"
	BudgetTokens int    `json:"budget_tokens"` // max thinking tokens
}

// TaskBudget constrains the total token budget for a task.
type TaskBudget struct {
	TotalTokens int `json:"total_tokens"`
}

// ---------------------------------------------------------------------------
// Streaming event types
// ---------------------------------------------------------------------------

// EventType identifies the kind of streaming event.
type EventType string

const (
	EventRequestStart  EventType = "request_start"
	EventTextDelta     EventType = "text_delta"
	EventThinkingDelta EventType = "thinking_delta"
	EventToolUse       EventType = "tool_use"
	EventAssistant     EventType = "assistant"
	EventError         EventType = "error"
)

// StreamEvent is a single event from a streaming LLM response.
type StreamEvent struct {
	Type       EventType `json:"type"`
	Text       string    `json:"text,omitempty"`
	Message    *Message  `json:"message,omitempty"`
	ToolUse    *ToolUseEvent `json:"tool_use,omitempty"`
	StopReason string    `json:"stop_reason,omitempty"`
	Usage      *Usage    `json:"usage,omitempty"`
	Error      string    `json:"error,omitempty"`
}

// ToolUseEvent carries tool invocation details.
type ToolUseEvent struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}
