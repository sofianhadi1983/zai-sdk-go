package zai

import (
	"context"

	"github.com/z-ai/zai-sdk-go/api/types/moderation"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// ModerationsService provides access to the Moderations API.
type ModerationsService struct {
	client *client.BaseClient
}

// newModerationsService creates a new moderations service.
func newModerationsService(baseClient *client.BaseClient) *ModerationsService {
	return &ModerationsService{
		client: baseClient,
	}
}

// Create performs content moderation and returns the results.
//
// Example:
//
//	req := moderation.NewTextModerationRequest("moderation", "content to check")
//
//	resp, err := client.Moderations.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.IsFlagged() {
//	    for _, result := range resp.GetResults() {
//	        if result.Categories.Hate {
//	            fmt.Println("Hate speech detected")
//	        }
//	        if result.Categories.Violence {
//	            fmt.Println("Violence detected")
//	        }
//	    }
//	}
func (s *ModerationsService) Create(ctx context.Context, req *moderation.ModerationRequest) (*moderation.ModerationResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/moderations", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp moderation.ModerationResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CheckText is a convenience method for checking a single text string.
//
// Example:
//
//	resp, err := client.Moderations.CheckText(ctx, "moderation", "text to check")
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.IsFlagged() {
//	    fmt.Println("Content flagged")
//	}
func (s *ModerationsService) CheckText(ctx context.Context, model string, text string) (*moderation.ModerationResponse, error) {
	req := moderation.NewTextModerationRequest(model, text)
	return s.Create(ctx, req)
}

// CheckBatch is a convenience method for checking multiple text strings at once.
//
// Example:
//
//	texts := []string{"text1", "text2", "text3"}
//	resp, err := client.Moderations.CheckBatch(ctx, "moderation", texts)
//	if err != nil {
//	    // Handle error
//	}
//
//	for i, result := range resp.GetResults() {
//	    if result.Flagged {
//	        fmt.Printf("Text %d flagged\n", i)
//	    }
//	}
func (s *ModerationsService) CheckBatch(ctx context.Context, model string, texts []string) (*moderation.ModerationResponse, error) {
	req := moderation.NewBatchTextModerationRequest(model, texts)
	return s.Create(ctx, req)
}
