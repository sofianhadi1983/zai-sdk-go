// Package images provides types for the Images API.
package images

import "github.com/sofianhadi1983/zai-sdk-go/internal/models"

// ImageSize represents the size of the generated image.
type ImageSize string

const (
	// Size1024x1024 generates a 1024x1024 square image (default).
	Size1024x1024 ImageSize = "1024x1024"

	// Recommended sizes from Z.ai API spec:
	// Size768x1344 generates a 768x1344 portrait image.
	Size768x1344 ImageSize = "768x1344"
	// Size864x1152 generates a 864x1152 portrait image.
	Size864x1152 ImageSize = "864x1152"
	// Size1344x768 generates a 1344x768 landscape image.
	Size1344x768 ImageSize = "1344x768"
	// Size1152x864 generates a 1152x864 landscape image.
	Size1152x864 ImageSize = "1152x864"
	// Size1440x720 generates a 1440x720 wide landscape image.
	Size1440x720 ImageSize = "1440x720"
	// Size720x1440 generates a 720x1440 tall portrait image.
	Size720x1440 ImageSize = "720x1440"

	// Legacy sizes (kept for compatibility):
	// Size1792x1024 generates a 1792x1024 landscape image.
	Size1792x1024 ImageSize = "1792x1024"
	// Size1024x1792 generates a 1024x1792 portrait image.
	Size1024x1792 ImageSize = "1024x1792"
)

// ImageQuality represents the quality of the generated image.
type ImageQuality string

const (
	// QualityStandard generates standard quality images.
	QualityStandard ImageQuality = "standard"
	// QualityHD generates high-definition images.
	QualityHD ImageQuality = "hd"
)

// ResponseFormat represents the format of the image data in the response.
type ResponseFormat string

const (
	// ResponseFormatURL returns image URLs.
	ResponseFormatURL ResponseFormat = "url"
	// ResponseFormatB64JSON returns base64-encoded JSON.
	ResponseFormatB64JSON ResponseFormat = "b64_json"
)

// ImageGenerationRequest represents a request to generate images.
type ImageGenerationRequest struct {
	// Model is the model to use for image generation (required).
	Model string `json:"model"`

	// Prompt is the text description of the desired image (required).
	Prompt string `json:"prompt"`

	// Size is the size of the generated images.
	// Defaults to "1024x1024" if not specified.
	Size ImageSize `json:"size,omitempty"`

	// Quality is the quality of the image.
	// Defaults to "standard" if not specified.
	Quality ImageQuality `json:"quality,omitempty"`

	// N is the number of images to generate.
	// Must be between 1 and 10. Defaults to 1 if not specified.
	N *int `json:"n,omitempty"`

	// ResponseFormat is the format in which the generated images are returned.
	// Defaults to "url" if not specified.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// UserID is a unique identifier for the end-user (6-128 characters).
	// Used for abuse detection and monitoring.
	UserID string `json:"user_id,omitempty"`
}

// NewImageGenerationRequest creates a new image generation request with required fields.
//
// Example:
//
//	req := images.NewImageGenerationRequest("cogview-3", "A beautiful sunset over mountains")
func NewImageGenerationRequest(model, prompt string) *ImageGenerationRequest {
	return &ImageGenerationRequest{
		Model:  model,
		Prompt: prompt,
	}
}

// SetSize sets the size of the generated images.
//
// Example:
//
//	req.SetSize(images.Size1024x1792)
func (r *ImageGenerationRequest) SetSize(size ImageSize) *ImageGenerationRequest {
	r.Size = size
	return r
}

// SetQuality sets the quality of the generated images.
//
// Example:
//
//	req.SetQuality(images.QualityHD)
func (r *ImageGenerationRequest) SetQuality(quality ImageQuality) *ImageGenerationRequest {
	r.Quality = quality
	return r
}

// SetN sets the number of images to generate.
//
// Example:
//
//	req.SetN(4)
func (r *ImageGenerationRequest) SetN(n int) *ImageGenerationRequest {
	r.N = &n
	return r
}

// SetResponseFormat sets the format of the image data in the response.
//
// Example:
//
//	req.SetResponseFormat(images.ResponseFormatB64JSON)
func (r *ImageGenerationRequest) SetResponseFormat(format ResponseFormat) *ImageGenerationRequest {
	r.ResponseFormat = format
	return r
}

// SetUserID sets the end-user identifier.
// The user ID should be 6-128 characters and is used for abuse detection.
//
// Example:
//
//	req.SetUserID("user-123456")
func (r *ImageGenerationRequest) SetUserID(userID string) *ImageGenerationRequest {
	r.UserID = userID
	return r
}

// ImageData represents a single generated image.
type ImageData struct {
	// URL is the URL of the generated image (when ResponseFormat is "url").
	URL string `json:"url,omitempty"`

	// B64JSON is the base64-encoded JSON of the generated image (when ResponseFormat is "b64_json").
	B64JSON string `json:"b64_json,omitempty"`

	// RevisedPrompt is the prompt that was used to generate the image, if the API modified the original prompt.
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// GetImageURL returns the URL of the image if available.
func (i *ImageData) GetImageURL() string {
	return i.URL
}

// GetBase64Data returns the base64-encoded image data if available.
func (i *ImageData) GetBase64Data() string {
	return i.B64JSON
}

// ContentFilterItem represents content safety filtering information.
type ContentFilterItem struct {
	// Role indicates the message role ("assistant", "user", or "history").
	Role string `json:"role"`

	// Level indicates the severity level (0-3, where 0 is most severe).
	Level int `json:"level"`
}

// ImageGenerationResponse represents a response from the image generation API.
type ImageGenerationResponse struct {
	// Created is the Unix timestamp when the images were created.
	Created int64 `json:"created"`

	// Data is the list of generated images.
	Data []ImageData `json:"data"`

	// ContentFilter contains safety information about the generated content.
	// Each item indicates the safety level for different parts of the interaction.
	ContentFilter []ContentFilterItem `json:"content_filter,omitempty"`

	// Usage contains token usage information (if available).
	Usage *models.Usage `json:"usage,omitempty"`
}

// GetFirstImage returns the first generated image, or nil if no images were generated.
func (r *ImageGenerationResponse) GetFirstImage() *ImageData {
	if len(r.Data) == 0 {
		return nil
	}
	return &r.Data[0]
}

// GetImageURLs returns all image URLs from the response.
func (r *ImageGenerationResponse) GetImageURLs() []string {
	urls := make([]string, 0, len(r.Data))
	for _, img := range r.Data {
		if img.URL != "" {
			urls = append(urls, img.URL)
		}
	}
	return urls
}

// GetBase64Images returns all base64-encoded images from the response.
func (r *ImageGenerationResponse) GetBase64Images() []string {
	images := make([]string, 0, len(r.Data))
	for _, img := range r.Data {
		if img.B64JSON != "" {
			images = append(images, img.B64JSON)
		}
	}
	return images
}
