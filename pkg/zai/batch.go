package zai

import (
	"context"
	"fmt"

	"github.com/z-ai/zai-sdk-go/api/types/batch"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// BatchService provides access to the Batch API.
type BatchService struct {
	client *client.BaseClient
}

// newBatchService creates a new batch service.
func newBatchService(baseClient *client.BaseClient) *BatchService {
	return &BatchService{
		client: baseClient,
	}
}

// Create creates a new batch processing job.
//
// Example:
//
//	req := batch.NewBatchCreateRequest(
//	    "24h",
//	    batch.EndpointChatCompletions,
//	    "file_abc123",
//	).SetMetadata(map[string]string{"user_id": "user_123"})
//
//	batchJob, err := client.Batch.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Batch ID: %s, Status: %s\n", batchJob.ID, batchJob.Status)
func (s *BatchService) Create(ctx context.Context, req *batch.BatchCreateRequest) (*batch.Batch, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/batches", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp batch.Batch
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Retrieve retrieves a batch by ID.
//
// Example:
//
//	batchJob, err := client.Batch.Retrieve(ctx, "batch_abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	if batchJob.IsCompleted() {
//	    fmt.Printf("Batch completed. Output file: %s\n", batchJob.OutputFileID)
//	} else if batchJob.IsInProgress() {
//	    fmt.Printf("Batch in progress. Completed: %d/%d\n",
//	        batchJob.RequestCounts.Completed, batchJob.RequestCounts.Total)
//	}
func (s *BatchService) Retrieve(ctx context.Context, batchID string) (*batch.Batch, error) {
	if batchID == "" {
		return nil, fmt.Errorf("batch ID cannot be empty")
	}

	// Make the API request
	apiResp, err := s.client.Get(ctx, fmt.Sprintf("/batches/%s", batchID), nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp batch.Batch
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// List lists batches with cursor-based pagination.
//
// Example:
//
//	// Get first page
//	resp, err := client.Batch.List(ctx, "", 20)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, batchJob := range resp.GetBatches() {
//	    fmt.Printf("Batch %s: %s\n", batchJob.ID, batchJob.Status)
//	}
//
//	// Get next page if available
//	if resp.HasMoreBatches() {
//	    nextResp, err := client.Batch.List(ctx, resp.LastID, 20)
//	    // Process next page...
//	}
func (s *BatchService) List(ctx context.Context, after string, limit int) (*batch.BatchListResponse, error) {
	// Build query parameters
	query := make(map[string]string)
	if after != "" {
		query["after"] = after
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}

	// Make the API request
	apiResp, err := s.client.Get(ctx, "/batches", query)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp batch.BatchListResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Cancel cancels an in-progress batch.
//
// Example:
//
//	batchJob, err := client.Batch.Cancel(ctx, "batch_abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	if batchJob.IsCancelling() {
//	    fmt.Println("Batch cancellation initiated")
//	} else if batchJob.IsCancelled() {
//	    fmt.Println("Batch has been cancelled")
//	}
func (s *BatchService) Cancel(ctx context.Context, batchID string) (*batch.Batch, error) {
	if batchID == "" {
		return nil, fmt.Errorf("batch ID cannot be empty")
	}

	// Make the API request
	apiResp, err := s.client.Post(ctx, fmt.Sprintf("/batches/%s/cancel", batchID), nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp batch.Batch
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
