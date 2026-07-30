package contracts

import (
	"context"
	"io"
)

// Role represents the role of a message in a conversation.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a single message in a chat completion.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// Choice represents a single generated choice in a chat completion.
type Choice struct {
	Message Message `json:"message"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Response represents the complete response from a chat completion.
type Response struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// ChatOptions defines parameters for the chat completion request.
type ChatOptions struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
	Seed        *int    `json:"seed,omitempty"`
}

// Provider defines the contract for an AI model provider.
type Provider interface {
	Name() string
	Chat(ctx context.Context, model string, messages []Message, opts ChatOptions) (*Response, error)
	ChatStream(ctx context.Context, model string, messages []Message, opts ChatOptions) (io.ReadCloser, error)
}
