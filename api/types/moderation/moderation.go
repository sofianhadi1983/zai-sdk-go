// Package moderation provides types for the Moderations API.
package moderation

// ModerationRequest represents a request to moderate content.
type ModerationRequest struct {
	// Model is the moderation model to use.
	Model string `json:"model"`

	// Input is the content to moderate.
	// Can be a string, array of strings, or structured input.
	Input interface{} `json:"input"`
}

// NewModerationRequest creates a new moderation request.
func NewModerationRequest(model string, input interface{}) *ModerationRequest {
	return &ModerationRequest{
		Model: model,
		Input: input,
	}
}

// NewTextModerationRequest creates a new moderation request for text input.
func NewTextModerationRequest(model string, text string) *ModerationRequest {
	return &ModerationRequest{
		Model: model,
		Input: map[string]interface{}{
			"type": "text",
			"text": text,
		},
	}
}

// NewBatchTextModerationRequest creates a new moderation request for multiple texts.
func NewBatchTextModerationRequest(model string, texts []string) *ModerationRequest {
	return &ModerationRequest{
		Model: model,
		Input: texts,
	}
}

// ModerationResponse represents the response from the Moderations API.
type ModerationResponse struct {
	// ID is the unique identifier for the moderation request.
	ID string `json:"id"`

	// Model is the model used to generate the moderation results.
	Model string `json:"model"`

	// Results contains the moderation results for each input.
	Results []ModerationResult `json:"results"`
}

// ModerationResult represents a single moderation result.
type ModerationResult struct {
	// Flagged indicates whether any category was flagged.
	Flagged bool `json:"flagged"`

	// Categories contains the flagged status for each category.
	Categories ModerationCategories `json:"categories"`

	// CategoryScores contains the confidence scores for each category.
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

// ModerationCategories contains the flagged status for each moderation category.
type ModerationCategories struct {
	// Harassment: Content that expresses, incites, or promotes harassing language.
	Harassment bool `json:"harassment"`

	// HarassmentThreatening: Harassment content that also includes violence or serious harm.
	HarassmentThreatening bool `json:"harassment/threatening"`

	// Hate: Content that expresses, incites, or promotes hate based on protected characteristics.
	Hate bool `json:"hate"`

	// HateThreatening: Hateful content that also includes violence or serious harm.
	HateThreatening bool `json:"hate/threatening"`

	// SelfHarm: Content that promotes, encourages, or depicts acts of self-harm.
	SelfHarm bool `json:"self-harm"`

	// SelfHarmInstructions: Content that encourages or gives instructions for self-harm.
	SelfHarmInstructions bool `json:"self-harm/instructions"`

	// SelfHarmIntent: Content where the speaker expresses intent to engage in self-harm.
	SelfHarmIntent bool `json:"self-harm/intent"`

	// Sexual: Content meant to arouse sexual excitement.
	Sexual bool `json:"sexual"`

	// SexualMinors: Sexual content involving individuals under 18 years old.
	SexualMinors bool `json:"sexual/minors"`

	// Violence: Content that depicts death, violence, or physical injury.
	Violence bool `json:"violence"`

	// ViolenceGraphic: Content that depicts violence or injury in graphic detail.
	ViolenceGraphic bool `json:"violence/graphic"`
}

// ModerationCategoryScores contains confidence scores for each moderation category.
type ModerationCategoryScores struct {
	// Harassment score (0.0 to 1.0)
	Harassment float64 `json:"harassment"`

	// HarassmentThreatening score (0.0 to 1.0)
	HarassmentThreatening float64 `json:"harassment/threatening"`

	// Hate score (0.0 to 1.0)
	Hate float64 `json:"hate"`

	// HateThreatening score (0.0 to 1.0)
	HateThreatening float64 `json:"hate/threatening"`

	// SelfHarm score (0.0 to 1.0)
	SelfHarm float64 `json:"self-harm"`

	// SelfHarmInstructions score (0.0 to 1.0)
	SelfHarmInstructions float64 `json:"self-harm/instructions"`

	// SelfHarmIntent score (0.0 to 1.0)
	SelfHarmIntent float64 `json:"self-harm/intent"`

	// Sexual score (0.0 to 1.0)
	Sexual float64 `json:"sexual"`

	// SexualMinors score (0.0 to 1.0)
	SexualMinors float64 `json:"sexual/minors"`

	// Violence score (0.0 to 1.0)
	Violence float64 `json:"violence"`

	// ViolenceGraphic score (0.0 to 1.0)
	ViolenceGraphic float64 `json:"violence/graphic"`
}

// GetResults returns the moderation results.
func (r *ModerationResponse) GetResults() []ModerationResult {
	if r.Results == nil {
		return []ModerationResult{}
	}
	return r.Results
}

// IsFlagged returns true if any result was flagged.
func (r *ModerationResponse) IsFlagged() bool {
	for _, result := range r.Results {
		if result.Flagged {
			return true
		}
	}
	return false
}

// HasCategory checks if any result was flagged for a specific category.
func (c *ModerationCategories) HasCategory(check func(*ModerationCategories) bool) bool {
	return check(c)
}

// IsSafe returns true if the result is not flagged.
func (r *ModerationResult) IsSafe() bool {
	return !r.Flagged
}
