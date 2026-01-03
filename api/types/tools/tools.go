// Package tools provides types for the Tools API.
package tools

import "github.com/sofianhadi1983/zai-sdk-go/api/types/chat"

// WebSearchRequest represents a request to perform web search using AI models.
type WebSearchRequest struct {
	// Model is the model to use for web search (e.g., "web-search-pro").
	Model string `json:"model"`

	// Messages contains the conversation context.
	// The current version supports single-turn conversations with User Message.
	// The tool will understand the User Message and perform a search.
	Messages []chat.Message `json:"messages"`

	// Stream indicates whether to stream the response.
	Stream bool `json:"stream,omitempty"`

	// RequestID is a unique identifier for the request.
	RequestID string `json:"request_id,omitempty"`

	// Scope specifies the search scope (e.g., entire web, academic).
	// Default is the entire web.
	Scope string `json:"scope,omitempty"`

	// Location specifies the user's location to improve relevance.
	Location string `json:"location,omitempty"`

	// RecentDays specifies returning search results updated in N days (1-30).
	RecentDays int `json:"recent_days,omitempty"`
}

// NewWebSearchRequest creates a new web search request.
func NewWebSearchRequest(model string, messages []chat.Message) *WebSearchRequest {
	return &WebSearchRequest{
		Model:    model,
		Messages: messages,
	}
}

// SetStream enables or disables streaming.
func (r *WebSearchRequest) SetStream(stream bool) *WebSearchRequest {
	r.Stream = stream
	return r
}

// SetRequestID sets the request ID.
func (r *WebSearchRequest) SetRequestID(requestID string) *WebSearchRequest {
	r.RequestID = requestID
	return r
}

// SetScope sets the search scope.
func (r *WebSearchRequest) SetScope(scope string) *WebSearchRequest {
	r.Scope = scope
	return r
}

// SetLocation sets the user's location.
func (r *WebSearchRequest) SetLocation(location string) *WebSearchRequest {
	r.Location = location
	return r
}

// SetRecentDays sets the recent days filter (1-30).
func (r *WebSearchRequest) SetRecentDays(days int) *WebSearchRequest {
	r.RecentDays = days
	return r
}

// SearchIntent represents search intent analysis.
type SearchIntent struct {
	// Index is the search round (default is 0).
	Index int `json:"index"`

	// Query is the optimized search query.
	Query string `json:"query"`

	// Intent is the determined intent type.
	Intent string `json:"intent"`

	// Keywords are the extracted search keywords.
	Keywords string `json:"keywords"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	// Index is the search round (default is 0).
	Index int `json:"index"`

	// Title is the result title.
	Title string `json:"title"`

	// Link is the result URL.
	Link string `json:"link"`

	// Content is the result content/snippet.
	Content string `json:"content"`

	// Icon is the result icon URL.
	Icon string `json:"icon"`

	// Media is the source media name.
	Media string `json:"media"`

	// Refer is the reference number (e.g., "[ref_1]").
	Refer string `json:"refer"`
}

// SearchRecommend represents a recommended search query.
type SearchRecommend struct {
	// Index is the search round (default is 0).
	Index int `json:"index"`

	// Query is the recommended query.
	Query string `json:"query"`
}

// WebSearchMessageToolCall represents a tool call in the web search message.
type WebSearchMessageToolCall struct {
	// ID is the tool call identifier.
	ID string `json:"id"`

	// Type is the type of tool call.
	Type string `json:"type"`

	// SearchIntent contains search intent information.
	SearchIntent *SearchIntent `json:"search_intent,omitempty"`

	// SearchResult contains search result data.
	SearchResult *SearchResult `json:"search_result,omitempty"`

	// SearchRecommend contains search recommendations.
	SearchRecommend *SearchRecommend `json:"search_recommend,omitempty"`
}

// WebSearchMessage represents a message in the web search response.
type WebSearchMessage struct {
	// Role is the message role (e.g., "assistant").
	Role string `json:"role"`

	// ToolCalls contains the tool calls made.
	ToolCalls []WebSearchMessageToolCall `json:"tool_calls,omitempty"`
}

// WebSearchChoice represents a choice in the web search response.
type WebSearchChoice struct {
	// Index is the choice index.
	Index int `json:"index"`

	// FinishReason is the reason why generation finished.
	FinishReason string `json:"finish_reason"`

	// Message contains the response message.
	Message WebSearchMessage `json:"message"`
}

// WebSearchResponse represents the response from a web search request.
type WebSearchResponse struct {
	// ID is the unique identifier for the response.
	ID string `json:"id,omitempty"`

	// Created is the creation timestamp.
	Created int64 `json:"created,omitempty"`

	// RequestID is the request identifier.
	RequestID string `json:"request_id,omitempty"`

	// Choices contains the response choices.
	Choices []WebSearchChoice `json:"choices"`
}

// GetChoices returns the response choices.
func (r *WebSearchResponse) GetChoices() []WebSearchChoice {
	if r.Choices == nil {
		return []WebSearchChoice{}
	}
	return r.Choices
}

// GetToolCalls returns all tool calls from the first choice.
func (r *WebSearchResponse) GetToolCalls() []WebSearchMessageToolCall {
	if len(r.Choices) == 0 {
		return []WebSearchMessageToolCall{}
	}
	return r.Choices[0].Message.ToolCalls
}

// GetSearchIntents returns all search intents from tool calls.
func (r *WebSearchResponse) GetSearchIntents() []*SearchIntent {
	intents := []*SearchIntent{}
	for _, toolCall := range r.GetToolCalls() {
		if toolCall.SearchIntent != nil {
			intents = append(intents, toolCall.SearchIntent)
		}
	}
	return intents
}

// GetSearchResults returns all search results from tool calls.
func (r *WebSearchResponse) GetSearchResults() []*SearchResult {
	results := []*SearchResult{}
	for _, toolCall := range r.GetToolCalls() {
		if toolCall.SearchResult != nil {
			results = append(results, toolCall.SearchResult)
		}
	}
	return results
}

// GetSearchRecommendations returns all search recommendations from tool calls.
func (r *WebSearchResponse) GetSearchRecommendations() []*SearchRecommend {
	recommends := []*SearchRecommend{}
	for _, toolCall := range r.GetToolCalls() {
		if toolCall.SearchRecommend != nil {
			recommends = append(recommends, toolCall.SearchRecommend)
		}
	}
	return recommends
}

// ChoiceDeltaToolCall represents a tool call delta in streaming responses.
type ChoiceDeltaToolCall struct {
	// Index is the index of the tool call.
	Index int `json:"index"`

	// ID is the unique identifier for the tool call.
	ID string `json:"id,omitempty"`

	// Type is the type of the tool call.
	Type string `json:"type,omitempty"`

	// SearchIntent contains search intent information.
	SearchIntent *SearchIntent `json:"search_intent,omitempty"`

	// SearchResult contains search result data.
	SearchResult *SearchResult `json:"search_result,omitempty"`

	// SearchRecommend contains search recommendations.
	SearchRecommend *SearchRecommend `json:"search_recommend,omitempty"`
}

// ChoiceDelta represents the delta changes in a streaming choice.
type ChoiceDelta struct {
	// Role is the role of the message sender.
	Role string `json:"role,omitempty"`

	// ToolCalls contains the list of tool call deltas.
	ToolCalls []ChoiceDeltaToolCall `json:"tool_calls,omitempty"`
}

// WebSearchStreamChoice represents a choice in the streaming response.
type WebSearchStreamChoice struct {
	// Index is the index of this choice in the response.
	Index int `json:"index"`

	// FinishReason is the reason why the generation finished.
	FinishReason string `json:"finish_reason,omitempty"`

	// Delta contains the delta changes for this choice.
	Delta ChoiceDelta `json:"delta"`
}

// WebSearchChunk represents a chunk in the web search streaming response.
type WebSearchChunk struct {
	// ID is the unique identifier for the chunk.
	ID string `json:"id,omitempty"`

	// Created is the timestamp when the chunk was created.
	Created int64 `json:"created,omitempty"`

	// Choices contains the list of choices in this chunk.
	Choices []WebSearchStreamChoice `json:"choices"`
}

// TokenizerRequest represents a request to count tokens.
type TokenizerRequest struct {
	// Model is the model to use for tokenization (required).
	// Options: "glm-4.6", "glm-4.6v", "glm-4.5"
	Model string `json:"model"`

	// Messages contains the conversation messages to tokenize (required).
	// Minimum one message required. Supports user, system, and assistant messages.
	Messages []chat.Message `json:"messages"`

	// Tools is an optional list of function definitions (max 128).
	Tools []chat.Tool `json:"tools,omitempty"`

	// RequestID is a unique identifier for the request.
	// If not provided, the platform will generate one automatically.
	RequestID string `json:"request_id,omitempty"`

	// UserID is a unique identifier for the end-user.
	UserID string `json:"user_id,omitempty"`
}

// NewTokenizerRequest creates a new tokenizer request.
//
// Example:
//
//	messages := []chat.Message{
//	    chat.NewUserMessage("Hello, world!"),
//	}
//	req := tools.NewTokenizerRequest("glm-4.6", messages)
func NewTokenizerRequest(model string, messages []chat.Message) *TokenizerRequest {
	return &TokenizerRequest{
		Model:    model,
		Messages: messages,
	}
}

// SetTools sets the tools for tokenization.
func (r *TokenizerRequest) SetTools(tools []chat.Tool) *TokenizerRequest {
	r.Tools = tools
	return r
}

// SetRequestID sets the request ID.
func (r *TokenizerRequest) SetRequestID(requestID string) *TokenizerRequest {
	r.RequestID = requestID
	return r
}

// SetUserID sets the user ID.
func (r *TokenizerRequest) SetUserID(userID string) *TokenizerRequest {
	r.UserID = userID
	return r
}

// TokenizerUsage represents token usage statistics.
type TokenizerUsage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// ImageTokens is the number of tokens for images.
	ImageTokens int `json:"image_tokens,omitempty"`

	// VideoTokens is the number of tokens for video.
	VideoTokens int `json:"video_tokens,omitempty"`

	// TotalTokens is the total number of tokens.
	TotalTokens int `json:"total_tokens"`
}

// TokenizerResponse represents the response from the tokenizer API.
type TokenizerResponse struct {
	// ID is the task sequence number from the platform.
	ID string `json:"id"`

	// Usage contains the token count statistics.
	Usage TokenizerUsage `json:"usage"`

	// Created is the Unix timestamp.
	Created int64 `json:"created,omitempty"`

	// RequestID is the client-submitted or platform-generated identifier.
	RequestID string `json:"request_id,omitempty"`
}
