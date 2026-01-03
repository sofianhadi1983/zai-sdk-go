// Package assistant provides types for the Assistant API.
package assistant

import (
	"encoding/json"
	"fmt"
)

// MessageTextContent represents text content for conversation messages.
type MessageTextContent struct {
	// Type is the content type, currently supports "text"
	Type string `json:"type"`

	// Text is the text content of the message
	Text string `json:"text"`
}

// MessageContent represents content in a message (can be text or other types).
type MessageContent interface {
	isMessageContent()
}

// Ensure MessageTextContent implements MessageContent
func (MessageTextContent) isMessageContent() {}

// ConversationMessage represents a message in a conversation.
type ConversationMessage struct {
	// Role is the message role (e.g., "user", "assistant")
	Role string `json:"role"`

	// Content is the list of message content items
	Content []MessageContent `json:"content"`
}

// AssistantAttachment represents file attachments for assistant conversations.
type AssistantAttachment struct {
	// FileID is the file identifier for attachment
	FileID string `json:"file_id"`
}

// TranslateParameters represents translation configuration.
type TranslateParameters struct {
	// FromLanguage is the source language code
	FromLanguage string `json:"from_language,omitempty"`

	// ToLanguage is the target language code
	ToLanguage string `json:"to_language,omitempty"`
}

// ExtraParameters represents extra parameters for assistant functionality.
type ExtraParameters struct {
	// Translate is the translation configuration
	Translate *TranslateParameters `json:"translate,omitempty"`
}

// ConversationRequest represents a request to create or continue a conversation.
type ConversationRequest struct {
	// AssistantID is the assistant identifier
	AssistantID string `json:"assistant_id"`

	// Messages is the list of conversation messages
	Messages []ConversationMessage `json:"messages"`

	// Model is the model name (optional, defaults to GLM-4-Assistant)
	Model string `json:"model,omitempty"`

	// Stream enables streaming SSE responses
	Stream bool `json:"stream"`

	// ConversationID is the conversation ID (creates new if not provided)
	ConversationID string `json:"conversation_id,omitempty"`

	// Attachments are optional file attachments
	Attachments []AssistantAttachment `json:"attachments,omitempty"`

	// Metadata is optional metadata extension field
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// RequestID is optional request identifier
	RequestID string `json:"request_id,omitempty"`

	// UserID is optional user identifier
	UserID string `json:"user_id,omitempty"`

	// ExtraParameters contains additional parameters
	ExtraParameters *ExtraParameters `json:"extra_parameters,omitempty"`
}

// NewConversationRequest creates a new conversation request.
//
// Example:
//
//	req := assistant.NewConversationRequest("asst_123", messages)
func NewConversationRequest(assistantID string, messages []ConversationMessage) *ConversationRequest {
	return &ConversationRequest{
		AssistantID: assistantID,
		Messages:    messages,
	}
}

// SetModel sets the model name.
func (r *ConversationRequest) SetModel(model string) *ConversationRequest {
	r.Model = model
	return r
}

// SetStream enables streaming responses.
func (r *ConversationRequest) SetStream(stream bool) *ConversationRequest {
	r.Stream = stream
	return r
}

// SetConversationID sets the conversation ID.
func (r *ConversationRequest) SetConversationID(id string) *ConversationRequest {
	r.ConversationID = id
	return r
}

// SetAttachments sets file attachments.
func (r *ConversationRequest) SetAttachments(attachments []AssistantAttachment) *ConversationRequest {
	r.Attachments = attachments
	return r
}

// SetMetadata sets metadata.
func (r *ConversationRequest) SetMetadata(metadata map[string]interface{}) *ConversationRequest {
	r.Metadata = metadata
	return r
}

// SetRequestID sets the request ID.
func (r *ConversationRequest) SetRequestID(id string) *ConversationRequest {
	r.RequestID = id
	return r
}

// SetUserID sets the user ID.
func (r *ConversationRequest) SetUserID(id string) *ConversationRequest {
	r.UserID = id
	return r
}

// SetExtraParameters sets extra parameters.
func (r *ConversationRequest) SetExtraParameters(params *ExtraParameters) *ConversationRequest {
	r.ExtraParameters = params
	return r
}

// ErrorInfo represents error information for assistant operations.
type ErrorInfo struct {
	// Code is the error code identifier
	Code string `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`
}

// TextContentBlock represents a text content block in assistant messages.
type TextContentBlock struct {
	// Content is the text content
	Content string `json:"content"`

	// Role is the role of the message sender
	Role string `json:"role"`

	// Type is the type identifier, always "content"
	Type string `json:"type"`
}

// Ensure TextContentBlock implements MessageContent
func (TextContentBlock) isMessageContent() {}

// ToolsDeltaBlock represents tool execution content in assistant messages.
type ToolsDeltaBlock struct {
	// Type is the type identifier for tools
	Type string `json:"type"`

	// ToolCallID is the tool call identifier
	ToolCallID string `json:"tool_call_id,omitempty"`

	// ToolName is the name of the tool being called
	ToolName string `json:"tool_name,omitempty"`

	// ToolOutput is the output from the tool
	ToolOutput string `json:"tool_output,omitempty"`
}

// Ensure ToolsDeltaBlock implements MessageContent
func (ToolsDeltaBlock) isMessageContent() {}

// AssistantChoice represents assistant response choice information.
type AssistantChoice struct {
	// Index is the choice result index
	Index int `json:"index"`

	// Delta is the current conversation output message content
	Delta MessageContent `json:"delta"`

	// FinishReason is the inference end reason:
	// - "stop": natural end or stop words
	// - "sensitive": content intercepted by security audit
	// - "network_error": service exception
	FinishReason string `json:"finish_reason"`

	// Metadata is the metadata extension field
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UnmarshalJSON custom unmarshaler for AssistantChoice to handle polymorphic Delta field.
func (c *AssistantChoice) UnmarshalJSON(data []byte) error {
	// Define an auxiliary type to avoid recursion
	type Alias AssistantChoice
	aux := &struct {
		Delta json.RawMessage `json:"delta"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the delta based on its type field
	var typeCheck struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(aux.Delta, &typeCheck); err != nil {
		return err
	}

	switch typeCheck.Type {
	case "content":
		var textBlock TextContentBlock
		if err := json.Unmarshal(aux.Delta, &textBlock); err != nil {
			return err
		}
		c.Delta = textBlock
	case "tools":
		var toolsBlock ToolsDeltaBlock
		if err := json.Unmarshal(aux.Delta, &toolsBlock); err != nil {
			return err
		}
		c.Delta = toolsBlock
	default:
		return fmt.Errorf("unknown delta type: %s", typeCheck.Type)
	}

	return nil
}

// CompletionUsage represents token usage statistics.
type CompletionUsage struct {
	// PromptTokens is the number of input tokens
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of output tokens
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens used
	TotalTokens int `json:"total_tokens"`
}

// AssistantCompletion represents the assistant completion response.
type AssistantCompletion struct {
	// ID is the unique request identifier
	ID string `json:"id"`

	// ConversationID is the conversation identifier
	ConversationID string `json:"conversation_id"`

	// AssistantID is the assistant identifier
	AssistantID string `json:"assistant_id"`

	// Created is the request creation time as Unix timestamp
	Created int64 `json:"created"`

	// Status is the response status:
	// - "completed": generation finished
	// - "in_progress": generating
	// - "failed": generation exception
	Status string `json:"status"`

	// LastError contains error information if generation failed
	LastError *ErrorInfo `json:"last_error,omitempty"`

	// Choices is the list of response choices with incremental information
	Choices []AssistantChoice `json:"choices"`

	// Metadata is optional metadata extension field
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Usage contains token usage statistics
	Usage *CompletionUsage `json:"usage,omitempty"`
}

// GetText returns the text content from the first choice's delta.
func (r *AssistantCompletion) GetText() string {
	if len(r.Choices) == 0 {
		return ""
	}

	delta := r.Choices[0].Delta
	if textBlock, ok := delta.(TextContentBlock); ok {
		return textBlock.Content
	}

	return ""
}

// IsCompleted returns true if the generation is completed.
func (r *AssistantCompletion) IsCompleted() bool {
	return r.Status == "completed"
}

// IsInProgress returns true if the generation is in progress.
func (r *AssistantCompletion) IsInProgress() bool {
	return r.Status == "in_progress"
}

// IsFailed returns true if the generation failed.
func (r *AssistantCompletion) IsFailed() bool {
	return r.Status == "failed"
}

// GetError returns the error message if generation failed.
func (r *AssistantCompletion) GetError() string {
	if r.LastError != nil {
		return r.LastError.Message
	}
	return ""
}

// AssistantSupport represents information about an assistant.
type AssistantSupport struct {
	// AssistantID is the assistant identifier
	AssistantID string `json:"assistant_id"`

	// CreatedAt is the assistant creation timestamp
	CreatedAt int64 `json:"created_at"`

	// UpdatedAt is the last update timestamp
	UpdatedAt int64 `json:"updated_at"`

	// Name is the assistant display name
	Name string `json:"name"`

	// Avatar is the assistant avatar URL or identifier
	Avatar string `json:"avatar"`

	// Description is the assistant description text
	Description string `json:"description"`

	// Status is the assistant status (currently only "publish")
	Status string `json:"status"`

	// Tools is the list of tool names supported by the assistant
	Tools []string `json:"tools"`

	// StarterPrompts are recommended startup prompts
	StarterPrompts []string `json:"starter_prompts"`
}

// AssistantSupportResponse represents the response for assistant support query.
type AssistantSupportResponse struct {
	// Code is the response status code
	Code int `json:"code"`

	// Message is the response message
	Message string `json:"msg"`

	// Data is the list of available assistants
	Data []AssistantSupport `json:"data"`
}

// GetAssistants returns the list of assistants.
func (r *AssistantSupportResponse) GetAssistants() []AssistantSupport {
	return r.Data
}

// ConversationUsage represents usage information for a conversation.
type ConversationUsage struct {
	// ID is the unique conversation identifier
	ID string `json:"id"`

	// AssistantID is the assistant identifier
	AssistantID string `json:"assistant_id"`

	// CreateTime is the conversation creation timestamp
	CreateTime int64 `json:"create_time"`

	// UpdateTime is the last update timestamp
	UpdateTime int64 `json:"update_time"`

	// Usage contains token usage statistics
	Usage CompletionUsage `json:"usage"`
}

// ConversationUsageList represents a list of conversation usage data.
type ConversationUsageList struct {
	// AssistantID is the assistant identifier
	AssistantID string `json:"assistant_id"`

	// HasMore indicates whether there are more pages available
	HasMore bool `json:"has_more"`

	// ConversationList is the list of conversation usage records
	ConversationList []ConversationUsage `json:"conversation_list"`
}

// ConversationUsageResponse represents the response for conversation usage list query.
type ConversationUsageResponse struct {
	// Code is the response status code
	Code int `json:"code"`

	// Message is the response message
	Message string `json:"msg"`

	// Data contains the conversation usage data
	Data ConversationUsageList `json:"data"`
}

// GetConversations returns the list of conversations.
func (r *ConversationUsageResponse) GetConversations() []ConversationUsage {
	return r.Data.ConversationList
}

// HasMore returns whether there are more pages available.
func (r *ConversationUsageResponse) HasMore() bool {
	return r.Data.HasMore
}
