package models

// Usage represents token usage information for API requests.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`

	// PromptTokensDetails provides detailed breakdown of prompt tokens.
	PromptTokensDetails *PromptTokensDetails `json:"prompt_tokens_details,omitempty"`

	// CompletionTokensDetails provides detailed breakdown of completion tokens.
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

// PromptTokensDetails provides detailed information about prompt tokens.
type PromptTokensDetails struct {
	// CachedTokens is the number of cached tokens reused.
	CachedTokens int `json:"cached_tokens,omitempty"`

	// AudioTokens is the number of audio tokens in the prompt.
	AudioTokens int `json:"audio_tokens,omitempty"`

	// TextTokens is the number of text tokens in the prompt.
	TextTokens int `json:"text_tokens,omitempty"`

	// ImageTokens is the number of image tokens in the prompt.
	ImageTokens int `json:"image_tokens,omitempty"`
}

// CompletionTokensDetails provides detailed information about completion tokens.
type CompletionTokensDetails struct {
	// ReasoningTokens is the number of reasoning tokens used.
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`

	// AudioTokens is the number of audio tokens in the completion.
	AudioTokens int `json:"audio_tokens,omitempty"`

	// TextTokens is the number of text tokens in the completion.
	TextTokens int `json:"text_tokens,omitempty"`
}

// IsEmpty returns true if usage information is empty.
func (u *Usage) IsEmpty() bool {
	return u == nil || (u.PromptTokens == 0 && u.CompletionTokens == 0 && u.TotalTokens == 0)
}

// HasCachedTokens returns true if cached tokens were used.
func (u *Usage) HasCachedTokens() bool {
	return u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0
}

// HasReasoningTokens returns true if reasoning tokens were used.
func (u *Usage) HasReasoningTokens() bool {
	return u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0
}

// GetCachedTokens returns the number of cached tokens, or 0 if none.
func (u *Usage) GetCachedTokens() int {
	if u.PromptTokensDetails == nil {
		return 0
	}
	return u.PromptTokensDetails.CachedTokens
}

// GetReasoningTokens returns the number of reasoning tokens, or 0 if none.
func (u *Usage) GetReasoningTokens() int {
	if u.CompletionTokensDetails == nil {
		return 0
	}
	return u.CompletionTokensDetails.ReasoningTokens
}
