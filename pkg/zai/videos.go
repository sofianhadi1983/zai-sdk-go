package zai

import (
	"context"
	"fmt"
	"time"

	"github.com/z-ai/zai-sdk-go/api/types/videos"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// VideosService provides access to the Videos API.
type VideosService struct {
	client *client.BaseClient
}

// newVideosService creates a new videos service.
func newVideosService(baseClient *client.BaseClient) *VideosService {
	return &VideosService{
		client: baseClient,
	}
}

// Create submits a video generation task.
//
// Example for text-to-video:
//
//	req := videos.NewTextToVideoRequest("cogvideox", "A cat playing with a ball")
//	task, err := client.Videos.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Task ID: %s\n", task.GetTaskID())
//
// Example for image-to-video:
//
//	req := videos.NewImageToVideoRequest("cogvideox", "https://example.com/image.jpg")
//	task, err := client.Videos.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
func (s *VideosService) Create(ctx context.Context, req *videos.VideoGenerationRequest) (*videos.VideoGenerationResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/videos/generations", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp videos.VideoGenerationResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Retrieve retrieves the status and result of a video generation task.
//
// Example:
//
//	result, err := client.Videos.Retrieve(ctx, "task-abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	if result.IsCompleted() {
//	    fmt.Printf("Video URL: %s\n", result.GetVideoURL())
//	} else if result.IsProcessing() {
//	    fmt.Println("Video is still being generated...")
//	} else if result.IsFailed() {
//	    fmt.Printf("Generation failed: %s\n", result.GetError())
//	}
func (s *VideosService) Retrieve(ctx context.Context, taskID string) (*videos.VideoResult, error) {
	// Make the API request
	path := fmt.Sprintf("/async-result/%s", taskID)
	apiResp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var result videos.VideoResult
	if err := s.client.ParseJSON(apiResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateText is a convenience method for text-to-video generation.
// Returns the task ID for later retrieval.
//
// Example:
//
//	taskID, err := client.Videos.GenerateText(ctx, "cogvideox", "A sunset over the ocean")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Video generation started. Task ID: %s\n", taskID)
func (s *VideosService) GenerateText(ctx context.Context, model videos.VideoModel, prompt string) (string, error) {
	req := videos.NewTextToVideoRequest(model, prompt)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetTaskID(), nil
}

// GenerateFromImage is a convenience method for image-to-video generation.
// Returns the task ID for later retrieval.
//
// Example:
//
//	taskID, err := client.Videos.GenerateFromImage(ctx, "cogvideox", "https://example.com/image.jpg")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Video generation started. Task ID: %s\n", taskID)
func (s *VideosService) GenerateFromImage(ctx context.Context, model videos.VideoModel, imageURL string) (string, error) {
	req := videos.NewImageToVideoRequest(model, imageURL)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetTaskID(), nil
}

// WaitForCompletion waits for a video generation task to complete.
// It polls the task status at regular intervals until completion or failure.
//
// Example:
//
//	result, err := client.Videos.WaitForCompletion(ctx, "task-abc123", 5*time.Second, 2*time.Minute)
//	if err != nil {
//	    // Handle error
//	}
//
//	if result.IsCompleted() {
//	    fmt.Printf("Video ready: %s\n", result.GetVideoURL())
//	}
func (s *VideosService) WaitForCompletion(ctx context.Context, taskID string, pollInterval, timeout time.Duration) (*videos.VideoResult, error) {
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		// Check deadline
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for video generation to complete")
		}

		// Check if context is done
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Retrieve current status
		result, err := s.Retrieve(ctx, taskID)
		if err != nil {
			return nil, err
		}

		// Check if task is complete or failed
		if result.IsCompleted() || result.IsFailed() {
			return result, nil
		}

		// Wait for next poll
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			// Continue polling
		}
	}
}
