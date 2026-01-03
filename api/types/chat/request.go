package chat

// ChatCompletionRequest represents a request to create a chat completion.
type ChatCompletionRequest struct {
	// Model is the ID of the model to use.
	// Required. Examples: "glm-4.7", "glm-4", "glm-3-turbo"
	Model string `json:"model"`

	// Messages is the list of messages in the conversation.
	// Required.
	Messages []Message `json:"messages"`

	// Temperature controls randomness in the output (0.0 to 1.0).
	// Higher values make output more random, lower values more deterministic.
	// Default: 0.95
	Temperature *float64 `json:"temperature,omitempty"`

	// TopP controls nucleus sampling (0.0 to 1.0).
	// Alternative to temperature for controlling randomness.
	// Default: 0.7
	TopP *float64 `json:"top_p,omitempty"`

	// Stream indicates whether to stream the response.
	// If true, tokens will be sent as server-sent events.
	Stream *bool `json:"stream,omitempty"`

	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens *int `json:"max_tokens,omitempty"`

	// Stop is a list of sequences where the API will stop generating.
	Stop []string `json:"stop,omitempty"`

	// Tools is a list of tools the model may call.
	Tools []Tool `json:"tools,omitempty"`

	// ToolChoice controls which (if any) tool is called by the model.
	ToolChoice interface{} `json:"tool_choice,omitempty"`

	// ResponseFormat specifies the format of the model's output.
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	// N is the number of completions to generate.
	// Default: 1
	N *int `json:"n,omitempty"`

	// PresencePenalty penalizes new tokens based on whether they appear in the text so far.
	// Range: -2.0 to 2.0
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty penalizes new tokens based on their frequency in the text so far.
	// Range: -2.0 to 2.0
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`

	// LogitBias modifies the likelihood of specified tokens appearing.
	LogitBias map[string]float64 `json:"logit_bias,omitempty"`

	// User is a unique identifier for the end-user.
	User string `json:"user,omitempty"`

	// RequestID is a unique identifier for the request.
	RequestID string `json:"request_id,omitempty"`

	// DoSample controls whether to use sampling.
	DoSample *bool `json:"do_sample,omitempty"`

	// Seed is the random seed for deterministic generation.
	Seed *int `json:"seed,omitempty"`

	// Extra fields for model-specific parameters.
	Extra map[string]interface{} `json:"-"`
}

// SetTemperature sets the temperature parameter.
func (r *ChatCompletionRequest) SetTemperature(temp float64) *ChatCompletionRequest {
	r.Temperature = &temp
	return r
}

// SetTopP sets the top_p parameter.
func (r *ChatCompletionRequest) SetTopP(topP float64) *ChatCompletionRequest {
	r.TopP = &topP
	return r
}

// SetStream sets whether to stream the response.
func (r *ChatCompletionRequest) SetStream(stream bool) *ChatCompletionRequest {
	r.Stream = &stream
	return r
}

// SetMaxTokens sets the maximum number of tokens to generate.
func (r *ChatCompletionRequest) SetMaxTokens(maxTokens int) *ChatCompletionRequest {
	r.MaxTokens = &maxTokens
	return r
}

// AddMessage adds a message to the conversation.
func (r *ChatCompletionRequest) AddMessage(message Message) *ChatCompletionRequest {
	r.Messages = append(r.Messages, message)
	return r
}

// AddUserMessage adds a user message to the conversation.
func (r *ChatCompletionRequest) AddUserMessage(content string) *ChatCompletionRequest {
	return r.AddMessage(NewUserMessage(content))
}

// AddSystemMessage adds a system message to the conversation.
func (r *ChatCompletionRequest) AddSystemMessage(content string) *ChatCompletionRequest {
	return r.AddMessage(NewSystemMessage(content))
}

// AddAssistantMessage adds an assistant message to the conversation.
func (r *ChatCompletionRequest) AddAssistantMessage(content string) *ChatCompletionRequest {
	return r.AddMessage(NewAssistantMessage(content))
}

// AddTool adds a tool to the request.
func (r *ChatCompletionRequest) AddTool(tool Tool) *ChatCompletionRequest {
	r.Tools = append(r.Tools, tool)
	return r
}

// SetToolChoice sets the tool choice parameter.
func (r *ChatCompletionRequest) SetToolChoice(choice ToolChoice) *ChatCompletionRequest {
	r.ToolChoice = string(choice)
	return r
}

// SetResponseFormat sets the response format.
func (r *ChatCompletionRequest) SetResponseFormat(format ResponseFormat) *ChatCompletionRequest {
	r.ResponseFormat = &format
	return r
}
