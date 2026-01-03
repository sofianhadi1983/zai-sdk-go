package chat

import "github.com/z-ai/zai-sdk-go/internal/models"

// ChatCompletionResponse represents the response from a chat completion request.
type ChatCompletionResponse struct {
	// ID is the unique identifier for the completion.
	ID string `json:"id"`

	// Object is the object type (always "chat.completion").
	Object string `json:"object"`

	// Created is the Unix timestamp of when the completion was created.
	Created int64 `json:"created"`

	// Model is the model used for the completion.
	Model string `json:"model"`

	// Choices is the list of completion choices.
	Choices []Choice `json:"choices"`

	// Usage is the token usage information.
	Usage *models.Usage `json:"usage,omitempty"`

	// SystemFingerprint is a unique identifier for the model configuration.
	SystemFingerprint string `json:"system_fingerprint,omitempty"`

	// Extra fields for model-specific data.
	Extra map[string]interface{} `json:"-"`
}

// Choice represents a completion choice.
type Choice struct {
	// Index is the index of this choice in the list.
	Index int `json:"index"`

	// Message is the generated message.
	Message Message `json:"message"`

	// FinishReason is the reason the model stopped generating.
	// Possible values: "stop", "length", "tool_calls", "content_filter", "function_call"
	FinishReason string `json:"finish_reason"`

	// LogProbs is the log probabilities for the choice.
	LogProbs *LogProbs `json:"logprobs,omitempty"`
}

// LogProbs represents log probability information.
type LogProbs struct {
	// Content is the log probabilities for each token.
	Content []TokenLogProb `json:"content,omitempty"`
}

// TokenLogProb represents log probability information for a token.
type TokenLogProb struct {
	// Token is the token string.
	Token string `json:"token"`

	// LogProb is the log probability of the token.
	LogProb float64 `json:"logprob"`

	// Bytes is the byte representation of the token.
	Bytes []int `json:"bytes,omitempty"`

	// TopLogProbs is the top log probabilities at this position.
	TopLogProbs []TopLogProb `json:"top_logprobs,omitempty"`
}

// TopLogProb represents a top log probability.
type TopLogProb struct {
	// Token is the token string.
	Token string `json:"token"`

	// LogProb is the log probability of the token.
	LogProb float64 `json:"logprob"`

	// Bytes is the byte representation of the token.
	Bytes []int `json:"bytes,omitempty"`
}

// GetFirstChoice returns the first choice from the response.
// Returns nil if there are no choices.
func (r *ChatCompletionResponse) GetFirstChoice() *Choice {
	if len(r.Choices) == 0 {
		return nil
	}
	return &r.Choices[0]
}

// GetContent returns the content of the first choice's message.
// Returns empty string if there are no choices or the content is not a string.
func (r *ChatCompletionResponse) GetContent() string {
	choice := r.GetFirstChoice()
	if choice == nil {
		return ""
	}

	if content, ok := choice.Message.Content.(string); ok {
		return content
	}

	return ""
}

// ChatCompletionChunk represents a chunk in a streaming chat completion.
type ChatCompletionChunk struct {
	// ID is the unique identifier for the completion.
	ID string `json:"id"`

	// Object is the object type (always "chat.completion.chunk").
	Object string `json:"object"`

	// Created is the Unix timestamp of when the chunk was created.
	Created int64 `json:"created"`

	// Model is the model used for the completion.
	Model string `json:"model"`

	// Choices is the list of chunk choices.
	Choices []ChunkChoice `json:"choices"`

	// SystemFingerprint is a unique identifier for the model configuration.
	SystemFingerprint string `json:"system_fingerprint,omitempty"`

	// Usage is the token usage information (only in the final chunk).
	Usage *models.Usage `json:"usage,omitempty"`
}

// ChunkChoice represents a choice in a streaming chunk.
type ChunkChoice struct {
	// Index is the index of this choice in the list.
	Index int `json:"index"`

	// Delta is the incremental message content.
	Delta Delta `json:"delta"`

	// FinishReason is the reason the model stopped generating.
	// Only present in the final chunk.
	FinishReason string `json:"finish_reason,omitempty"`

	// LogProbs is the log probabilities for the choice.
	LogProbs *LogProbs `json:"logprobs,omitempty"`
}

// Delta represents incremental message content in a streaming chunk.
type Delta struct {
	// Role is the role of the message author (only in the first chunk).
	Role Role `json:"role,omitempty"`

	// Content is the incremental content.
	Content string `json:"content,omitempty"`

	// ToolCalls are incremental tool calls.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// FunctionCall is the incremental function call (deprecated).
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// GetContent returns the content from the first choice's delta.
// Returns empty string if there are no choices.
func (c *ChatCompletionChunk) GetContent() string {
	if len(c.Choices) == 0 {
		return ""
	}
	return c.Choices[0].Delta.Content
}

// IsFinished returns true if this chunk indicates the completion is finished.
func (c *ChatCompletionChunk) IsFinished() bool {
	if len(c.Choices) == 0 {
		return false
	}
	return c.Choices[0].FinishReason != ""
}
