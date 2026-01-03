// Package videos provides types for the Videos API.
package videos

// VideoModel represents the video generation model.
type VideoModel string

const (
	// ModelCogVideoX is the CogVideoX model for video generation.
	ModelCogVideoX VideoModel = "cogvideox"
)

// TaskStatus represents the status of a video generation task.
type TaskStatus string

const (
	// StatusSubmitted indicates the task has been submitted.
	StatusSubmitted TaskStatus = "submitted"
	// StatusProcessing indicates the task is being processed.
	StatusProcessing TaskStatus = "processing"
	// StatusCompleted indicates the task has completed successfully.
	StatusCompleted TaskStatus = "completed"
	// StatusFailed indicates the task has failed.
	StatusFailed TaskStatus = "failed"
)

// VideoGenerationRequest represents a request to generate a video.
type VideoGenerationRequest struct {
	// Model is the model to use for video generation (required).
	Model VideoModel `json:"model"`

	// Prompt is the text description of the desired video (required for text-to-video).
	Prompt string `json:"prompt,omitempty"`

	// ImageURL is the URL of the image to animate (required for image-to-video).
	ImageURL string `json:"image_url,omitempty"`

	// User is a unique identifier representing your end-user.
	User string `json:"user,omitempty"`
}

// NewTextToVideoRequest creates a new text-to-video generation request.
//
// Example:
//
//	req := videos.NewTextToVideoRequest("cogvideox", "A cat playing with a ball")
func NewTextToVideoRequest(model VideoModel, prompt string) *VideoGenerationRequest {
	return &VideoGenerationRequest{
		Model:  model,
		Prompt: prompt,
	}
}

// NewImageToVideoRequest creates a new image-to-video generation request.
//
// Example:
//
//	req := videos.NewImageToVideoRequest("cogvideox", "https://example.com/image.jpg")
func NewImageToVideoRequest(model VideoModel, imageURL string) *VideoGenerationRequest {
	return &VideoGenerationRequest{
		Model:    model,
		ImageURL: imageURL,
	}
}

// SetUser sets the user identifier.
//
// Example:
//
//	req.SetUser("user-123")
func (r *VideoGenerationRequest) SetUser(user string) *VideoGenerationRequest {
	r.User = user
	return r
}

// VideoTask represents a video generation task.
type VideoTask struct {
	// ID is the unique identifier for the task.
	ID string `json:"id"`

	// Model is the model used for generation.
	Model VideoModel `json:"model"`

	// Status is the current status of the task.
	Status TaskStatus `json:"task_status"`

	// RequestID is the request identifier.
	RequestID string `json:"request_id,omitempty"`
}

// VideoGenerationResponse represents the response from creating a video generation task.
type VideoGenerationResponse struct {
	// ID is the task ID.
	ID string `json:"id"`

	// Model is the model used.
	Model VideoModel `json:"model,omitempty"`

	// RequestID is the request identifier.
	RequestID string `json:"request_id,omitempty"`
}

// VideoResult represents a completed video generation result.
type VideoResult struct {
	// TaskID is the ID of the task.
	TaskID string `json:"task_id"`

	// TaskStatus is the status of the task.
	TaskStatus TaskStatus `json:"task_status"`

	// RequestID is the request identifier.
	RequestID string `json:"request_id,omitempty"`

	// VideoResult contains the generated video data.
	VideoResult []VideoData `json:"video_result,omitempty"`

	// ErrorMessage contains error information if the task failed.
	ErrorMessage string `json:"error_message,omitempty"`
}

// VideoData represents video generation data.
type VideoData struct {
	// URL is the URL of the generated video.
	URL string `json:"url,omitempty"`

	// CoverImageURL is the URL of the video cover image.
	CoverImageURL string `json:"cover_image_url,omitempty"`
}

// GetTaskID returns the task ID.
func (r *VideoGenerationResponse) GetTaskID() string {
	return r.ID
}

// IsSubmitted returns true if the task status is submitted.
func (t *VideoTask) IsSubmitted() bool {
	return t.Status == StatusSubmitted
}

// IsProcessing returns true if the task is being processed.
func (t *VideoTask) IsProcessing() bool {
	return t.Status == StatusProcessing
}

// IsCompleted returns true if the task has completed successfully.
func (t *VideoTask) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// IsFailed returns true if the task has failed.
func (t *VideoTask) IsFailed() bool {
	return t.Status == StatusFailed
}

// IsCompleted returns true if the video generation completed successfully.
func (r *VideoResult) IsCompleted() bool {
	return r.TaskStatus == StatusCompleted
}

// IsFailed returns true if the video generation failed.
func (r *VideoResult) IsFailed() bool {
	return r.TaskStatus == StatusFailed
}

// IsProcessing returns true if the video is still being generated.
func (r *VideoResult) IsProcessing() bool {
	return r.TaskStatus == StatusProcessing || r.TaskStatus == StatusSubmitted
}

// GetFirstVideo returns the first generated video, or nil if no videos were generated.
func (r *VideoResult) GetFirstVideo() *VideoData {
	if len(r.VideoResult) == 0 {
		return nil
	}
	return &r.VideoResult[0]
}

// GetVideoURL returns the URL of the first generated video.
func (r *VideoResult) GetVideoURL() string {
	video := r.GetFirstVideo()
	if video == nil {
		return ""
	}
	return video.URL
}

// GetCoverImageURL returns the cover image URL of the first generated video.
func (r *VideoResult) GetCoverImageURL() string {
	video := r.GetFirstVideo()
	if video == nil {
		return ""
	}
	return video.CoverImageURL
}

// GetAllVideoURLs returns all video URLs from the result.
func (r *VideoResult) GetAllVideoURLs() []string {
	urls := make([]string, 0, len(r.VideoResult))
	for _, video := range r.VideoResult {
		if video.URL != "" {
			urls = append(urls, video.URL)
		}
	}
	return urls
}

// HasError returns true if there is an error message.
func (r *VideoResult) HasError() bool {
	return r.ErrorMessage != ""
}

// GetError returns the error message if any.
func (r *VideoResult) GetError() string {
	return r.ErrorMessage
}

// GetURL returns the video URL.
func (v *VideoData) GetURL() string {
	return v.URL
}

// GetCoverImageURL returns the cover image URL.
func (v *VideoData) GetCoverImageURL() string {
	return v.CoverImageURL
}
