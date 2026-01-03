package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/z-ai/zai-sdk-go/api/types/videos"
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

	// Example 1: Text-to-video generation
	fmt.Println("=== Example 1: Text-to-Video Generation ===")
	textToVideoExample(ctx, client)

	// Example 2: Image-to-video generation
	fmt.Println("\n=== Example 2: Image-to-Video Generation ===")
	imageToVideoExample(ctx, client)

	// Example 3: Check generation status
	fmt.Println("\n=== Example 3: Check Generation Status ===")
	checkStatusExample(ctx, client)

	// Example 4: Wait for completion
	fmt.Println("\n=== Example 4: Wait for Completion ===")
	waitForCompletionExample(ctx, client)

	// Example 5: Handle generation errors
	fmt.Println("\n=== Example 5: Handle Errors ===")
	errorHandlingExample(ctx, client)
}

func textToVideoExample(ctx context.Context, client *zai.Client) {
	// Create a text-to-video request
	req := videos.NewTextToVideoRequest(
		videos.ModelCogVideoX,
		"A cat playing with a ball of yarn in a sunny garden",
	)

	// Submit the video generation task
	task, err := client.Videos.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Video generation task created!\n")
	fmt.Printf("Task ID: %s\n", task.GetTaskID())
	fmt.Printf("Model: %s\n", task.Model)
	fmt.Printf("\nUse this task ID to check the status later.\n")
}

func imageToVideoExample(ctx context.Context, client *zai.Client) {
	// Create an image-to-video request
	// This animates a static image into a video
	req := videos.NewImageToVideoRequest(
		videos.ModelCogVideoX,
		"https://example.com/your-image.jpg",
	)
	req.SetUser("user-example-123")

	// Submit the task
	task, err := client.Videos.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Image-to-video task created!\n")
	fmt.Printf("Task ID: %s\n", task.ID)
}

func checkStatusExample(ctx context.Context, client *zai.Client) {
	// Use the convenience method to generate a video
	taskID, err := client.Videos.GenerateText(
		ctx,
		videos.ModelCogVideoX,
		"A sunrise over mountains",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Task ID: %s\n", taskID)

	// Check the status
	result, err := client.Videos.Retrieve(ctx, taskID)
	if err != nil {
		log.Printf("Error retrieving status: %v", err)
		return
	}

	fmt.Printf("Task Status: %s\n", result.TaskStatus)

	if result.IsProcessing() {
		fmt.Println("Video is still being generated...")
		fmt.Println("Check back later or use WaitForCompletion()")
	} else if result.IsCompleted() {
		fmt.Printf("Video is ready!\n")
		fmt.Printf("Video URL: %s\n", result.GetVideoURL())
		fmt.Printf("Cover Image URL: %s\n", result.GetCoverImageURL())
	} else if result.IsFailed() {
		fmt.Printf("Generation failed: %s\n", result.GetError())
	}
}

func waitForCompletionExample(ctx context.Context, client *zai.Client) {
	// Generate a video
	taskID, err := client.Videos.GenerateText(
		ctx,
		videos.ModelCogVideoX,
		"A peaceful lake at sunset with flying birds",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Task ID: %s\n", taskID)
	fmt.Println("Waiting for video generation to complete...")
	fmt.Println("This may take several minutes...")

	// Wait for the video to be generated
	// Poll every 5 seconds, timeout after 5 minutes
	result, err := client.Videos.WaitForCompletion(
		ctx,
		taskID,
		5*time.Second,  // Poll interval
		5*time.Minute,  // Timeout
	)
	if err != nil {
		log.Printf("Error waiting for completion: %v", err)
		return
	}

	if result.IsCompleted() {
		fmt.Println("\n✓ Video generation completed!")
		fmt.Printf("Video URL: %s\n", result.GetVideoURL())
		fmt.Printf("Cover Image: %s\n", result.GetCoverImageURL())

		// Get all video URLs if multiple were generated
		allURLs := result.GetAllVideoURLs()
		fmt.Printf("Total videos generated: %d\n", len(allURLs))
	} else if result.IsFailed() {
		fmt.Printf("\n✗ Video generation failed: %s\n", result.GetError())
	}
}

func errorHandlingExample(ctx context.Context, client *zai.Client) {
	// Demonstrate error handling for different scenarios

	// Example 1: Invalid prompt (empty)
	fmt.Println("1. Testing with empty prompt:")
	_, err := client.Videos.GenerateText(ctx, videos.ModelCogVideoX, "")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 2: Check failed task
	fmt.Println("\n2. Checking a failed task:")
	// In a real scenario, you would have a task ID that failed
	// This is just to demonstrate how to handle errors
	taskID, _ := client.Videos.GenerateText(ctx, videos.ModelCogVideoX, "test prompt")
	result, err := client.Videos.Retrieve(ctx, taskID)
	if err != nil {
		fmt.Printf("   Error retrieving task: %v\n", err)
		return
	}

	if result.HasError() {
		fmt.Printf("   Task has error: %s\n", result.GetError())
	}

	// Example 3: Timeout handling
	fmt.Println("\n3. Demonstrating timeout (with short timeout):")
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err = client.Videos.WaitForCompletion(
		ctx,
		taskID,
		50*time.Millisecond,
		1*time.Second,
	)
	if err != nil {
		fmt.Printf("   Expected timeout/cancel: %v\n", err)
	}
}

// Helper function to demonstrate manual polling
func manualPollingExample(ctx context.Context, client *zai.Client, taskID string) {
	fmt.Println("Manually polling for completion...")

	maxAttempts := 60 // Max 5 minutes at 5 second intervals
	pollInterval := 5 * time.Second

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Retrieve current status
		result, err := client.Videos.Retrieve(ctx, taskID)
		if err != nil {
			log.Printf("Error checking status: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		// Check if done
		if result.IsCompleted() {
			fmt.Printf("✓ Video ready: %s\n", result.GetVideoURL())
			return
		} else if result.IsFailed() {
			fmt.Printf("✗ Generation failed: %s\n", result.GetError())
			return
		}

		// Still processing, wait and try again
		fmt.Printf("Still processing... (attempt %d/%d)\n", attempt+1, maxAttempts)
		time.Sleep(pollInterval)
	}

	fmt.Println("Timeout: Video generation took too long")
}
