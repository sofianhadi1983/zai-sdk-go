package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/fileparser"
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

	// Example 1: Asynchronous file parsing with different tool types
	fmt.Println("=== Async File Parsing Example ===")
	asyncParsingExample(ctx, client)

	fmt.Println()

	// Example 2: Synchronous file parsing
	fmt.Println("=== Sync File Parsing Example ===")
	syncParsingExample(ctx, client)

	fmt.Println()

	// Example 3: Retrieve parsing results as text
	fmt.Println("=== Get Parsing Results as Text Example ===")
	getResultsTextExample(ctx, client)

	fmt.Println()

	// Example 4: Retrieve parsing results as download link
	fmt.Println("=== Get Parsing Results as Download Link Example ===")
	getResultsDownloadLinkExample(ctx, client)
}

func asyncParsingExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual file path
	filePath := os.Getenv("PARSER_FILE_PATH")
	if filePath == "" {
		fmt.Println("Skipping example: PARSER_FILE_PATH not set")
		fmt.Println("Usage: export PARSER_FILE_PATH=/path/to/document.pdf")
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info for the file type
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	// Create parsing task with Prime tool (most advanced)
	req := fileparser.NewCreateRequest(
		file,
		fileInfo.Name(),
		"pdf", // or "docx", "txt", etc.
		fileparser.ToolTypePrime,
	)

	resp, err := client.FileParser.Create(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create parsing task: %v", err)
	}

	if resp.Success {
		fmt.Printf("Task created successfully!\n")
		fmt.Printf("Task ID: %s\n", resp.TaskID)
		fmt.Printf("Message: %s\n", resp.Message)
		fmt.Printf("\nUse this Task ID to retrieve results with Content() method\n")
		fmt.Printf("Example: client.FileParser.Content(ctx, fileparser.NewContentRequest(\"%s\", fileparser.FormatTypeText))\n", resp.TaskID)
	} else {
		fmt.Printf("Task creation failed: %s\n", resp.Message)
	}
}

func syncParsingExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual file path
	filePath := os.Getenv("PARSER_SYNC_FILE_PATH")
	if filePath == "" {
		fmt.Println("Skipping example: PARSER_SYNC_FILE_PATH not set")
		fmt.Println("Usage: export PARSER_SYNC_FILE_PATH=/path/to/document.docx")
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	// Create synchronous parsing request
	req := fileparser.NewSyncRequest(
		file,
		fileInfo.Name(),
		"docx", // file type
	)

	fmt.Println("Parsing file synchronously (this may take a moment)...")
	resp, err := client.FileParser.CreateSync(ctx, req)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	fmt.Printf("Task ID: %s\n", resp.TaskID)
	fmt.Printf("Status: %t\n", resp.Status)
	fmt.Printf("Message: %s\n", resp.Message)

	if resp.Status && resp.HasContent() {
		fmt.Println("\nParsed Content:")
		fmt.Println("---")
		// Print first 500 characters of content
		content := resp.GetContent()
		if len(content) > 500 {
			fmt.Printf("%s...\n", content[:500])
			fmt.Printf("\n(Content truncated, total length: %d characters)\n", len(content))
		} else {
			fmt.Println(content)
		}
		fmt.Println("---")

		if resp.GetDownloadURL() != "" {
			fmt.Printf("\nDownload URL: %s\n", resp.GetDownloadURL())
		}
	} else {
		fmt.Println("\nNo content available in the response")
	}
}

func getResultsTextExample(ctx context.Context, client *zai.Client) {
	// Note: This requires a task ID from a previous Create() call
	taskID := os.Getenv("PARSER_TASK_ID")
	if taskID == "" {
		fmt.Println("Skipping example: PARSER_TASK_ID not set")
		fmt.Println("Usage: First create a task, then export PARSER_TASK_ID=<task_id>")
		return
	}

	// Small delay to allow processing
	fmt.Println("Waiting for task to complete...")
	time.Sleep(2 * time.Second)

	// Get results as text
	req := fileparser.NewContentRequest(taskID, fileparser.FormatTypeText)

	resp, err := client.FileParser.Content(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get parsing results: %v", err)
	}

	if resp.HasContent() {
		fmt.Println("\nParsed Text Content:")
		fmt.Println("---")
		content := resp.GetContent()
		// Print first 500 characters
		if len(content) > 500 {
			fmt.Printf("%s...\n", content[:500])
			fmt.Printf("\n(Content truncated, total length: %d characters)\n", len(content))
		} else {
			fmt.Println(content)
		}
		fmt.Println("---")
	} else {
		fmt.Println("No text content available")
	}
}

func getResultsDownloadLinkExample(ctx context.Context, client *zai.Client) {
	// Note: This requires a task ID from a previous Create() call
	taskID := os.Getenv("PARSER_TASK_ID_DOWNLOAD")
	if taskID == "" {
		fmt.Println("Skipping example: PARSER_TASK_ID_DOWNLOAD not set")
		fmt.Println("Usage: First create a task, then export PARSER_TASK_ID_DOWNLOAD=<task_id>")
		return
	}

	// Small delay to allow processing
	fmt.Println("Waiting for task to complete...")
	time.Sleep(2 * time.Second)

	// Get results as download link
	req := fileparser.NewContentRequest(taskID, fileparser.FormatTypeDownloadLink)

	resp, err := client.FileParser.Content(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get parsing results: %v", err)
	}

	if resp.HasData() {
		fmt.Printf("Received binary data: %d bytes\n", len(resp.GetData()))

		// Optionally save to file
		outputPath := "parsed_result.bin"
		err := os.WriteFile(outputPath, resp.GetData(), 0644)
		if err != nil {
			log.Fatalf("Failed to save file: %v", err)
		}
		fmt.Printf("Saved to: %s\n", outputPath)
	} else {
		fmt.Println("No binary data available")
	}
}
