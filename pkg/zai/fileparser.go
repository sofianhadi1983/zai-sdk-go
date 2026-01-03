package zai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/fileparser"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
)

// FileParserService provides access to the File Parser API.
type FileParserService struct {
	client *client.BaseClient
}

// newFileParserService creates a new file parser service.
func newFileParserService(baseClient *client.BaseClient) *FileParserService {
	return &FileParserService{
		client: baseClient,
	}
}

// Create creates an asynchronous file parsing task.
//
// Example:
//
//	file, err := os.Open("document.pdf")
//	if err != nil {
//	    // Handle error
//	}
//	defer file.Close()
//
//	req := fileparser.NewCreateRequest(
//	    file,
//	    "document.pdf",
//	    "pdf",
//	    fileparser.ToolTypePrime,
//	)
//
//	resp, err := client.FileParser.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.Success {
//	    fmt.Printf("Task created: %s\n", resp.TaskID)
//	    // Use the TaskID to retrieve results later with Content()
//	}
func (s *FileParserService) Create(ctx context.Context, req *fileparser.CreateRequest) (*fileparser.CreateResponse, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file_type field (required)
	if err := writer.WriteField("file_type", req.FileType); err != nil {
		return nil, fmt.Errorf("failed to write file_type field: %w", err)
	}

	// Add the tool_type field (required)
	if err := writer.WriteField("tool_type", string(req.ToolType)); err != nil {
		return nil, fmt.Errorf("failed to write tool_type field: %w", err)
	}

	// Add the file
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to the form
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Make the API request using PostMultipart
	apiResp, err := s.client.PostMultipart(ctx, "/files/parser/create", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp fileparser.CreateResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Content retrieves the parsing result for a completed task.
//
// Example:
//
//	// Get result as text
//	req := fileparser.NewContentRequest(taskID, fileparser.FormatTypeText)
//	resp, err := client.FileParser.Content(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.HasContent() {
//	    fmt.Printf("Parsed text: %s\n", resp.GetContent())
//	}
//
// Example with download link:
//
//	// Get result as download link
//	req := fileparser.NewContentRequest(taskID, fileparser.FormatTypeDownloadLink)
//	resp, err := client.FileParser.Content(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.HasData() {
//	    // Save binary data to file
//	    os.WriteFile("result.bin", resp.GetData(), 0644)
//	}
func (s *FileParserService) Content(ctx context.Context, req *fileparser.ContentRequest) (*fileparser.ContentResponse, error) {
	// Build the path
	path := fmt.Sprintf("/files/parser/result/%s/%s", req.TaskID, req.FormatType)

	// Make the API request
	apiResp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	// Read the response body
	data, err := io.ReadAll(apiResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	resp := &fileparser.ContentResponse{}

	// Check the format type to determine how to handle the response
	if req.FormatType == fileparser.FormatTypeText {
		// For text format, store as string
		resp.Content = string(data)
	} else {
		// For download_link format, store as binary data
		resp.Data = data
	}

	return resp, nil
}

// CreateSync creates a synchronous file parsing task and returns the result immediately.
//
// Example:
//
//	file, err := os.Open("document.docx")
//	if err != nil {
//	    // Handle error
//	}
//	defer file.Close()
//
//	req := fileparser.NewSyncRequest(file, "document.docx", "docx")
//
//	resp, err := client.FileParser.CreateSync(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.Status && resp.HasContent() {
//	    fmt.Printf("Parsed content: %s\n", resp.GetContent())
//	    fmt.Printf("Download URL: %s\n", resp.GetDownloadURL())
//	}
func (s *FileParserService) CreateSync(ctx context.Context, req *fileparser.SyncRequest) (*fileparser.SyncResponse, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file_type field (required)
	if err := writer.WriteField("file_type", req.FileType); err != nil {
		return nil, fmt.Errorf("failed to write file_type field: %w", err)
	}

	// Add the tool_type field (must be "prime-sync")
	if err := writer.WriteField("tool_type", string(req.ToolType)); err != nil {
		return nil, fmt.Errorf("failed to write tool_type field: %w", err)
	}

	// Add the file
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to the form
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Make the API request using PostMultipart
	apiResp, err := s.client.PostMultipart(ctx, "/files/parser/sync", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp fileparser.SyncResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
