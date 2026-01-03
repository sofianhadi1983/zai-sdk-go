// Package agents provides types for the Agents API.
package agents

import "github.com/sofianhadi1983/zai-sdk-go/api/types/chat"

// AgentInvokeRequest represents a request to invoke an agent.
type AgentInvokeRequest struct {
	// AgentID is the unique identifier of the agent.
	AgentID string `json:"agent_id"`

	// Messages contains the conversation messages.
	Messages []chat.Message `json:"messages"`

	// Stream indicates whether to stream the response.
	Stream bool `json:"stream,omitempty"`

	// RequestID is a unique identifier for the request.
	RequestID string `json:"request_id,omitempty"`

	// UserID is the user identifier.
	UserID string `json:"user_id,omitempty"`

	// CustomVariables are custom variables for the agent.
	CustomVariables map[string]interface{} `json:"custom_variables,omitempty"`

	// SensitiveWordCheck configures sensitive word filtering.
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`
}

// NewAgentInvokeRequest creates a new agent invoke request.
func NewAgentInvokeRequest(agentID string, messages []chat.Message) *AgentInvokeRequest {
	return &AgentInvokeRequest{
		AgentID:  agentID,
		Messages: messages,
	}
}

// SetStream enables or disables streaming.
func (r *AgentInvokeRequest) SetStream(stream bool) *AgentInvokeRequest {
	r.Stream = stream
	return r
}

// SetRequestID sets the request ID.
func (r *AgentInvokeRequest) SetRequestID(requestID string) *AgentInvokeRequest {
	r.RequestID = requestID
	return r
}

// SetUserID sets the user ID.
func (r *AgentInvokeRequest) SetUserID(userID string) *AgentInvokeRequest {
	r.UserID = userID
	return r
}

// SetCustomVariables sets custom variables.
func (r *AgentInvokeRequest) SetCustomVariables(variables map[string]interface{}) *AgentInvokeRequest {
	r.CustomVariables = variables
	return r
}

// SetSensitiveWordCheck sets sensitive word check configuration.
func (r *AgentInvokeRequest) SetSensitiveWordCheck(check *SensitiveWordCheck) *AgentInvokeRequest {
	r.SensitiveWordCheck = check
	return r
}

// AgentAsyncResultRequest represents a request to get async agent result.
type AgentAsyncResultRequest struct {
	// AgentID is the unique identifier of the agent.
	AgentID string `json:"agent_id"`

	// AsyncID is the asynchronous operation ID.
	AsyncID string `json:"async_id,omitempty"`

	// ConversationID is the conversation identifier.
	ConversationID string `json:"conversation_id,omitempty"`

	// CustomVariables are custom variables for the agent.
	CustomVariables map[string]interface{} `json:"custom_variables,omitempty"`
}

// NewAgentAsyncResultRequest creates a new async result request.
func NewAgentAsyncResultRequest(agentID string) *AgentAsyncResultRequest {
	return &AgentAsyncResultRequest{
		AgentID: agentID,
	}
}

// SetAsyncID sets the async operation ID.
func (r *AgentAsyncResultRequest) SetAsyncID(asyncID string) *AgentAsyncResultRequest {
	r.AsyncID = asyncID
	return r
}

// SetConversationID sets the conversation ID.
func (r *AgentAsyncResultRequest) SetConversationID(conversationID string) *AgentAsyncResultRequest {
	r.ConversationID = conversationID
	return r
}

// SetCustomVariables sets custom variables.
func (r *AgentAsyncResultRequest) SetCustomVariables(variables map[string]interface{}) *AgentAsyncResultRequest {
	r.CustomVariables = variables
	return r
}

// SensitiveWordCheck represents sensitive word filtering configuration.
type SensitiveWordCheck struct {
	// Type is the type of sensitive word check.
	Type string `json:"type"`

	// Status indicates if the check is enabled.
	Status string `json:"status"`
}

// AgentCompletionMessage represents a message in an agent completion response.
type AgentCompletionMessage struct {
	// Role is the role of the message sender.
	Role string `json:"role"`

	// Content is the content of the message.
	Content interface{} `json:"content,omitempty"`
}

// AgentCompletionChoice represents a choice in the agent completion response.
type AgentCompletionChoice struct {
	// Index is the index of this choice.
	Index int `json:"index"`

	// FinishReason is the reason why the generation finished.
	FinishReason string `json:"finish_reason"`

	// Message contains the completion message.
	Message AgentCompletionMessage `json:"message"`
}

// AgentCompletionUsage represents token usage statistics.
type AgentCompletionUsage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
}

// AgentError represents an error in agent completion.
type AgentError struct {
	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Message is the error message.
	Message string `json:"message,omitempty"`
}

// AgentCompletionResponse represents the response from an agent invocation.
type AgentCompletionResponse struct {
	// ID is the unique identifier of the completion.
	ID string `json:"id,omitempty"`

	// AgentID is the unique identifier of the agent.
	AgentID string `json:"agent_id,omitempty"`

	// ConversationID is the unique identifier of the conversation.
	ConversationID string `json:"conversation_id,omitempty"`

	// Status is the status of the completion.
	Status string `json:"status,omitempty"`

	// RequestID is the unique identifier of the request.
	RequestID string `json:"request_id,omitempty"`

	// Choices contains the completion choices.
	Choices []AgentCompletionChoice `json:"choices"`

	// Usage contains token usage statistics.
	Usage *AgentCompletionUsage `json:"usage,omitempty"`

	// Error contains error information if any.
	Error *AgentError `json:"error,omitempty"`
}

// GetChoices returns the completion choices.
func (r *AgentCompletionResponse) GetChoices() []AgentCompletionChoice {
	if r.Choices == nil {
		return []AgentCompletionChoice{}
	}
	return r.Choices
}

// GetContent returns the content from the first choice.
func (r *AgentCompletionResponse) GetContent() interface{} {
	if len(r.Choices) == 0 {
		return nil
	}
	return r.Choices[0].Message.Content
}

// HasError returns true if there is an error.
func (r *AgentCompletionResponse) HasError() bool {
	return r.Error != nil
}

// AgentChoiceDelta represents the delta changes in a streaming choice.
type AgentChoiceDelta struct {
	// Role is the role of the message sender.
	Role string `json:"role,omitempty"`

	// Content is the content delta.
	Content interface{} `json:"content,omitempty"`
}

// AgentStreamChoice represents a choice in the streaming response.
type AgentStreamChoice struct {
	// Index is the index of this choice in the response.
	Index int `json:"index"`

	// FinishReason is the reason why the generation finished.
	FinishReason string `json:"finish_reason,omitempty"`

	// Delta contains the delta changes for this choice.
	Delta AgentChoiceDelta `json:"delta"`
}

// AgentCompletionChunk represents a chunk in the agent streaming response.
type AgentCompletionChunk struct {
	// ID is the unique identifier of the chunk.
	ID string `json:"id,omitempty"`

	// AgentID is the unique identifier of the agent.
	AgentID string `json:"agent_id,omitempty"`

	// ConversationID is the unique identifier of the conversation.
	ConversationID string `json:"conversation_id,omitempty"`

	// Choices contains the list of choices in this chunk.
	Choices []AgentStreamChoice `json:"choices"`

	// Usage contains token usage statistics.
	Usage *AgentCompletionUsage `json:"usage,omitempty"`

	// Error contains error information if any.
	Error *AgentError `json:"error,omitempty"`
}

// GetContent returns the content from the first choice delta.
func (c *AgentCompletionChunk) GetContent() interface{} {
	if len(c.Choices) == 0 {
		return nil
	}
	return c.Choices[0].Delta.Content
}

// HasError returns true if there is an error.
func (c *AgentCompletionChunk) HasError() bool {
	return c.Error != nil
}
