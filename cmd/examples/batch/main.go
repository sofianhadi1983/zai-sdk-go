package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/z-ai/zai-sdk-go/api/types/batch"
	"github.com/z-ai/zai-sdk-go/pkg/zai"
)

func main() {
	// Create a new client
	client, err := zai.NewClient(
		zai.WithAPIKey("your-api-key.your-secret"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Create a batch processing job
	fmt.Println("=== Example 1: Create Batch ==")
	createBatchExample(ctx, client)

	// Example 2: Retrieve batch status
	fmt.Println("\n=== Example 2: Retrieve Batch Status ===")
	retrieveBatchExample(ctx, client)

	// Example 3: List all batches
	fmt.Println("\n=== Example 3: List Batches ===")
	listBatchesExample(ctx, client)

	// Example 4: Monitor batch progress
	fmt.Println("\n=== Example 4: Monitor Batch Progress ===")
	monitorBatchExample(ctx, client)

	// Example 5: Cancel a batch
	fmt.Println("\n=== Example 5: Cancel Batch ===")
	cancelBatchExample(ctx, client)

	// Example 6: List with pagination
	fmt.Println("\n=== Example 6: Paginated Batch List ===")
	paginatedListExample(ctx, client)
}

func createBatchExample(ctx context.Context, client *zai.Client) {
	// Create a batch request
	// Note: You must first upload a JSONL file with purpose "batch"
	// Each line in the file should be a separate API request
	req := batch.NewBatchCreateRequest(
		"24h",                           // Completion window
		batch.EndpointChatCompletions,   // Endpoint
		"file_abc123",                   // Input file ID (from Files.Upload)
	).SetMetadata(map[string]string{
		"user_id": "user_123",
		"batch_name": "test_batch",
	})

	batchJob, err := client.Batch.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Batch created successfully!\n")
	fmt.Printf("  ID: %s\n", batchJob.ID)
	fmt.Printf("  Status: %s\n", batchJob.Status)
	fmt.Printf("  Input File: %s\n", batchJob.InputFileID)
	fmt.Printf("  Created At: %v\n", time.Unix(batchJob.CreatedAt, 0))

	if batchJob.Metadata != nil {
		fmt.Printf("  Metadata:\n")
		for key, value := range batchJob.Metadata {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}
}

func retrieveBatchExample(ctx context.Context, client *zai.Client) {
	// Retrieve a specific batch by ID
	batchID := "batch_abc123"
	batchJob, err := client.Batch.Retrieve(ctx, batchID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Batch Details:\n")
	fmt.Printf("  ID: %s\n", batchJob.ID)
	fmt.Printf("  Status: %s\n", batchJob.Status)
	fmt.Printf("  Endpoint: %s\n", batchJob.Endpoint)
	fmt.Printf("  Completion Window: %s\n", batchJob.CompletionWindow)

	// Check status
	if batchJob.IsCompleted() {
		fmt.Printf("  ✓ Batch completed successfully\n")
		if batchJob.OutputFileID != "" {
			fmt.Printf("  Output File: %s\n", batchJob.OutputFileID)
		}
	} else if batchJob.IsInProgress() {
		fmt.Printf("  ⏳ Batch is in progress\n")
	} else if batchJob.IsFailed() {
		fmt.Printf("  ✗ Batch failed\n")
		if batchJob.ErrorFileID != "" {
			fmt.Printf("  Error File: %s\n", batchJob.ErrorFileID)
		}
	}

	// Show request counts if available
	if batchJob.RequestCounts != nil {
		fmt.Printf("  Request Counts:\n")
		fmt.Printf("    Total: %d\n", batchJob.RequestCounts.Total)
		fmt.Printf("    Completed: %d\n", batchJob.RequestCounts.Completed)
		fmt.Printf("    Failed: %d\n", batchJob.RequestCounts.Failed)
		if batchJob.RequestCounts.Total > 0 {
			completionRate := float64(batchJob.RequestCounts.Completed) / float64(batchJob.RequestCounts.Total) * 100
			fmt.Printf("    Completion Rate: %.2f%%\n", completionRate)
		}
	}

	// Show error details if available
	if batchJob.Errors != nil && len(batchJob.Errors.Data) > 0 {
		fmt.Printf("  Errors:\n")
		for i, batchErr := range batchJob.Errors.Data {
			fmt.Printf("    Error %d:\n", i+1)
			fmt.Printf("      Code: %s\n", batchErr.Code)
			fmt.Printf("      Message: %s\n", batchErr.Message)
			if batchErr.Line > 0 {
				fmt.Printf("      Line: %d\n", batchErr.Line)
			}
			if batchErr.Param != "" {
				fmt.Printf("      Parameter: %s\n", batchErr.Param)
			}
		}
	}
}

func listBatchesExample(ctx context.Context, client *zai.Client) {
	// List all batches (first page)
	resp, err := client.Batch.List(ctx, "", 10)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d batches:\n", len(resp.GetBatches()))
	for _, batchJob := range resp.GetBatches() {
		fmt.Printf("\nBatch: %s\n", batchJob.ID)
		fmt.Printf("  Status: %s\n", batchJob.Status)
		fmt.Printf("  Endpoint: %s\n", batchJob.Endpoint)
		fmt.Printf("  Created: %v\n", time.Unix(batchJob.CreatedAt, 0))

		// Show completion status
		statusIndicator := "○"
		if batchJob.IsCompleted() {
			statusIndicator = "✓"
		} else if batchJob.IsFailed() {
			statusIndicator = "✗"
		} else if batchJob.IsInProgress() {
			statusIndicator = "⏳"
		}
		fmt.Printf("  %s %s\n", statusIndicator, batchJob.Status)
	}

	if resp.HasMoreBatches() {
		fmt.Printf("\nMore batches available. Use cursor: %s\n", resp.LastID)
	}
}

func monitorBatchExample(ctx context.Context, client *zai.Client) {
	// Monitor a batch until completion
	batchID := "batch_abc123"

	fmt.Printf("Monitoring batch: %s\n", batchID)

	// Poll every 10 seconds for up to 5 minutes
	maxAttempts := 30
	pollInterval := 10 * time.Second

	for i := 0; i < maxAttempts; i++ {
		batchJob, err := client.Batch.Retrieve(ctx, batchID)
		if err != nil {
			log.Printf("Error retrieving batch: %v", err)
			return
		}

		fmt.Printf("[%s] Status: %s", time.Now().Format("15:04:05"), batchJob.Status)

		if batchJob.RequestCounts != nil {
			progress := float64(batchJob.RequestCounts.Completed) / float64(batchJob.RequestCounts.Total) * 100
			fmt.Printf(" - Progress: %d/%d (%.1f%%)",
				batchJob.RequestCounts.Completed,
				batchJob.RequestCounts.Total,
				progress)
		}
		fmt.Println()

		// Check if batch reached terminal state
		if batchJob.IsTerminal() {
			if batchJob.IsCompleted() {
				fmt.Printf("✓ Batch completed successfully!\n")
				fmt.Printf("  Output File: %s\n", batchJob.OutputFileID)
			} else if batchJob.IsFailed() {
				fmt.Printf("✗ Batch failed\n")
				if batchJob.ErrorFileID != "" {
					fmt.Printf("  Error File: %s\n", batchJob.ErrorFileID)
				}
			} else if batchJob.IsCancelled() {
				fmt.Printf("⊗ Batch was cancelled\n")
			} else if batchJob.IsExpired() {
				fmt.Printf("⏱ Batch expired\n")
			}
			break
		}

		// Wait before next poll
		if i < maxAttempts-1 {
			time.Sleep(pollInterval)
		}
	}
}

func cancelBatchExample(ctx context.Context, client *zai.Client) {
	// Cancel an in-progress batch
	batchID := "batch_abc123"

	batchJob, err := client.Batch.Cancel(ctx, batchID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Batch cancellation initiated\n")
	fmt.Printf("  ID: %s\n", batchJob.ID)
	fmt.Printf("  Status: %s\n", batchJob.Status)

	if batchJob.IsCancelling() {
		fmt.Printf("  ⏳ Cancellation in progress...\n")
		if batchJob.CancellingAt != nil {
			fmt.Printf("  Cancelling At: %v\n", time.Unix(*batchJob.CancellingAt, 0))
		}
	} else if batchJob.IsCancelled() {
		fmt.Printf("  ✓ Batch has been cancelled\n")
		if batchJob.CancelledAt != nil {
			fmt.Printf("  Cancelled At: %v\n", time.Unix(*batchJob.CancelledAt, 0))
		}
	}
}

func paginatedListExample(ctx context.Context, client *zai.Client) {
	// Demonstrate cursor-based pagination
	var allBatches []batch.Batch
	cursor := ""
	pageSize := 5

	fmt.Printf("Fetching all batches using pagination (page size: %d)\n", pageSize)

	for page := 1; ; page++ {
		fmt.Printf("Fetching page %d...\n", page)

		resp, err := client.Batch.List(ctx, cursor, pageSize)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		batches := resp.GetBatches()
		fmt.Printf("  Retrieved %d batches\n", len(batches))

		allBatches = append(allBatches, batches...)

		if !resp.HasMoreBatches() {
			fmt.Printf("  ✓ No more batches\n")
			break
		}

		cursor = resp.LastID
		fmt.Printf("  → Next cursor: %s\n", cursor)
	}

	fmt.Printf("\nTotal batches retrieved: %d\n", len(allBatches))

	// Group by status
	statusCounts := make(map[string]int)
	for _, batchJob := range allBatches {
		statusCounts[batchJob.Status]++
	}

	fmt.Printf("\nBatches by status:\n")
	for status, count := range statusCounts {
		fmt.Printf("  %s: %d\n", status, count)
	}
}

// Error handling example
func errorHandlingExample(ctx context.Context, client *zai.Client) {
	// Example 1: Handle invalid batch ID
	fmt.Println("1. Testing with invalid batch ID:")
	_, err := client.Batch.Retrieve(ctx, "invalid_id")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 2: Handle empty batch ID
	fmt.Println("\n2. Testing with empty batch ID:")
	_, err = client.Batch.Retrieve(ctx, "")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 3: Check batch completion with proper error handling
	fmt.Println("\n3. Checking batch completion status:")
	batchJob, err := client.Batch.Retrieve(ctx, "batch_abc123")
	if err != nil {
		fmt.Printf("   Error retrieving batch: %v\n", err)
		return
	}

	if batchJob.IsFailed() {
		fmt.Printf("   ✗ Batch failed\n")
		if batchJob.Errors != nil {
			fmt.Printf("   Error count: %d\n", len(batchJob.Errors.Data))
			for i, batchErr := range batchJob.Errors.Data {
				fmt.Printf("   Error %d: %s - %s\n", i+1, batchErr.Code, batchErr.Message)
			}
		}
	} else if batchJob.IsCompleted() {
		fmt.Printf("   ✓ Batch completed successfully\n")
	}
}
