package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/files"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
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

	// Example 1: Upload a file
	fmt.Println("=== Example 1: Upload a File ===")
	uploadedFile := uploadFileExample(ctx, client)

	// Example 2: List all files
	fmt.Println("\n=== Example 2: List Files ===")
	listFilesExample(ctx, client)

	// Example 3: Retrieve file information
	if uploadedFile != nil {
		fmt.Println("\n=== Example 3: Retrieve File Info ===")
		retrieveFileExample(ctx, client, uploadedFile.ID)

		// Example 4: Retrieve file content
		fmt.Println("\n=== Example 4: Retrieve File Content ===")
		retrieveContentExample(ctx, client, uploadedFile.ID)

		// Example 5: Delete a file
		fmt.Println("\n=== Example 5: Delete a File ===")
		deleteFileExample(ctx, client, uploadedFile.ID)
	}

	// Example 6: Upload from real file
	fmt.Println("\n=== Example 6: Upload from Real File ===")
	uploadRealFileExample(ctx, client)

	// Example 7: Filter files by purpose
	fmt.Println("\n=== Example 7: Filter Files by Purpose ===")
	filterFilesByPurposeExample(ctx, client)
}

func uploadFileExample(ctx context.Context, client *zai.Client) *files.File {
	// Create a sample file content
	// In a real scenario, you would read from an actual file
	fileContent := strings.NewReader(`{"prompt": "Hello", "completion": "World"}
{"prompt": "What is AI?", "completion": "Artificial Intelligence"}
{"prompt": "Go is", "completion": "a programming language"}`)

	// Create upload request
	req := files.NewFileUploadRequest(
		fileContent,
		"training_data.jsonl",
		files.PurposeFineTune,
	)

	// Upload the file
	uploadedFile, err := client.Files.Upload(ctx, req)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return nil
	}

	fmt.Printf("File uploaded successfully!\n")
	fmt.Printf("  ID: %s\n", uploadedFile.ID)
	fmt.Printf("  Filename: %s\n", uploadedFile.Filename)
	fmt.Printf("  Size: %d bytes\n", uploadedFile.Bytes)
	fmt.Printf("  Purpose: %s\n", uploadedFile.Purpose)
	fmt.Printf("  Status: %s\n", uploadedFile.Status)

	return uploadedFile
}

func listFilesExample(ctx context.Context, client *zai.Client) {
	// List all files
	fileList, err := client.Files.List(ctx)
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return
	}

	fmt.Printf("Found %d files:\n", len(fileList.Data))
	for i, file := range fileList.GetFiles() {
		fmt.Printf("  %d. %s (%s) - %s\n",
			i+1,
			file.Filename,
			file.ID,
			file.Purpose,
		)
	}

	if fileList.HasMore {
		fmt.Println("  (more files available)")
	}
}

func retrieveFileExample(ctx context.Context, client *zai.Client, fileID string) {
	// Retrieve file information
	file, err := client.Files.Retrieve(ctx, fileID)
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		return
	}

	fmt.Printf("File information:\n")
	fmt.Printf("  ID: %s\n", file.GetID())
	fmt.Printf("  Filename: %s\n", file.GetFilename())
	fmt.Printf("  Size: %d bytes\n", file.GetSize())
	fmt.Printf("  Purpose: %s\n", file.GetPurpose())
	fmt.Printf("  Status: %s\n", file.Status)
	fmt.Printf("  Created: %d\n", file.CreatedAt)

	if file.IsUploaded() {
		fmt.Println("  ✓ File is ready to use")
	}

	if file.HasError() {
		fmt.Printf("  ✗ Error: %s\n", file.StatusDetails)
	}
}

func retrieveContentExample(ctx context.Context, client *zai.Client, fileID string) {
	// Retrieve file content
	content, err := client.Files.RetrieveContent(ctx, fileID)
	if err != nil {
		log.Printf("Error retrieving content: %v", err)
		return
	}

	fmt.Printf("File content (%s):\n", content.GetContentType())
	fmt.Printf("Content length: %d bytes\n", len(content.GetContent()))

	// Display first 200 characters
	contentStr := content.String()
	if len(contentStr) > 200 {
		fmt.Printf("Preview: %s...\n", contentStr[:200])
	} else {
		fmt.Printf("Content:\n%s\n", contentStr)
	}
}

func deleteFileExample(ctx context.Context, client *zai.Client, fileID string) {
	// Delete the file
	deleteResp, err := client.Files.Delete(ctx, fileID)
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return
	}

	if deleteResp.IsDeleted() {
		fmt.Printf("File %s deleted successfully\n", deleteResp.ID)
	} else {
		fmt.Printf("Failed to delete file %s\n", deleteResp.ID)
	}
}

func uploadRealFileExample(ctx context.Context, client *zai.Client) {
	// This example shows how to upload a real file from disk
	// First, create a temporary file for demonstration
	tmpFile, err := os.CreateTemp("", "example-*.jsonl")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write some sample data
	_, err = tmpFile.WriteString(`{"prompt": "Example 1", "completion": "Response 1"}
{"prompt": "Example 2", "completion": "Response 2"}`)
	if err != nil {
		log.Printf("Error writing to temp file: %v", err)
		return
	}
	tmpFile.Close()

	// Now upload the file
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	req := files.NewFileUploadRequest(
		file,
		"real_file_example.jsonl",
		files.PurposeFineTune,
	)

	uploadedFile, err := client.Files.Upload(ctx, req)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return
	}

	fmt.Printf("Real file uploaded successfully!\n")
	fmt.Printf("  ID: %s\n", uploadedFile.ID)
	fmt.Printf("  Filename: %s\n", uploadedFile.Filename)

	// Clean up by deleting the uploaded file
	_, _ = client.Files.Delete(ctx, uploadedFile.ID)
}

func filterFilesByPurposeExample(ctx context.Context, client *zai.Client) {
	// List all files
	fileList, err := client.Files.List(ctx)
	if err != nil {
		log.Printf("Error listing files: %v", err)
		return
	}

	// Filter by different purposes
	purposes := []files.FilePurpose{
		files.PurposeFineTune,
		files.PurposeAssistants,
		files.PurposeBatch,
	}

	for _, purpose := range purposes {
		purposeFiles := fileList.GetFilesByPurpose(purpose)
		fmt.Printf("%s files: %d\n", purpose, len(purposeFiles))
		for _, file := range purposeFiles {
			fmt.Printf("  - %s (%s)\n", file.Filename, file.ID)
		}
	}

	// Get all file IDs
	allIDs := fileList.GetFileIDs()
	fmt.Printf("\nTotal files: %d\n", len(allIDs))

	// Find a specific file by ID
	if len(allIDs) > 0 {
		firstFile := fileList.GetFileByID(allIDs[0])
		if firstFile != nil {
			fmt.Printf("First file: %s\n", firstFile.Filename)
		}
	}
}
