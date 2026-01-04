package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/ocr"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
	// Create client from environment variables
	client, err := zai.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Handwriting OCR with probability scores
	fmt.Println("=== Handwriting OCR Example ===")
	handwritingOCRExample(ctx, client)

	fmt.Println()

	// Example 2: Handwriting OCR with language specification
	fmt.Println("=== Handwriting OCR with Language Example ===")
	handwritingOCRWithLanguageExample(ctx, client)
}

func handwritingOCRExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual image file path
	imagePath := os.Getenv("OCR_IMAGE_PATH")
	if imagePath == "" {
		fmt.Println("Skipping example: OCR_IMAGE_PATH not set")
		fmt.Println("Usage: export OCR_IMAGE_PATH=/path/to/handwriting.jpg")
		return
	}

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image file: %v", err)
	}
	defer file.Close()

	// Create OCR request with probability scores
	req := ocr.NewOCRRequest(file, "handwriting.jpg", ocr.ToolTypeHandWrite).
		SetProbability(true)

	// Perform OCR
	resp, err := client.OCR.HandwritingOCR(ctx, req)
	if err != nil {
		log.Fatalf("Failed to perform OCR: %v", err)
	}

	// Display results
	fmt.Printf("Task ID: %s\n", resp.TaskID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Number of results: %d\n", resp.WordsResultNum)

	if resp.HasResults() {
		fmt.Println("\nRecognized text:")
		fmt.Printf("Full text: %s\n", resp.GetText())

		fmt.Println("\nDetailed results:")
		for i, result := range resp.GetResults() {
			fmt.Printf("\n%d. Text: %s\n", i+1, result.Words)
			fmt.Printf("   Location: (x=%d, y=%d, w=%d, h=%d)\n",
				result.Location.Left,
				result.Location.Top,
				result.Location.Width,
				result.Location.Height)

			if result.Probability != nil {
				fmt.Printf("   Confidence:\n")
				fmt.Printf("     Average: %.2f%%\n", result.Probability.Average*100)
				fmt.Printf("     Minimum: %.2f%%\n", result.Probability.Min*100)
				fmt.Printf("     Variance: %.4f\n", result.Probability.Variance)
			}
		}
	} else {
		fmt.Println("\nNo text detected in the image")
	}
}

func handwritingOCRWithLanguageExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual image file path with Chinese text
	imagePath := os.Getenv("OCR_CHINESE_IMAGE_PATH")
	if imagePath == "" {
		fmt.Println("Skipping example: OCR_CHINESE_IMAGE_PATH not set")
		fmt.Println("Usage: export OCR_CHINESE_IMAGE_PATH=/path/to/chinese_handwriting.jpg")
		return
	}

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image file: %v", err)
	}
	defer file.Close()

	// Create OCR request with Chinese language
	req := ocr.NewOCRRequest(file, "chinese_handwriting.jpg", ocr.ToolTypeHandWrite).
		SetLanguageType("zh-CN").
		SetProbability(true)

	// Perform OCR
	resp, err := client.OCR.HandwritingOCR(ctx, req)
	if err != nil {
		log.Fatalf("Failed to perform OCR: %v", err)
	}

	// Display results
	fmt.Printf("Task ID: %s\n", resp.TaskID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Language: Chinese (zh-CN)\n")

	if resp.HasResults() {
		fmt.Println("\nRecognized Chinese text:")
		fmt.Printf("%s\n", resp.GetText())

		fmt.Printf("\nTotal characters recognized: %d\n", resp.WordsResultNum)
	} else {
		fmt.Println("\nNo text detected in the image")
	}
}

// Additional example functions can be added here for different use cases:
//
// - Batch OCR processing
// - OCR with specific image preprocessing
// - Async OCR for large images
// - OCR result export to different formats
