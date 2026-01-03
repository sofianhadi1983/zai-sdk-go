package zai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/z-ai/zai-sdk-go/api/types/files"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// FilesService provides access to the Files API.
type FilesService struct {
	client *client.BaseClient
}

// newFilesService creates a new files service.
func newFilesService(baseClient *client.BaseClient) *FilesService {
	return &FilesService{
		client: baseClient,
	}
}

// Upload uploads a file to the API.
//
// Example:
//
//	file, err := os.Open("training_data.jsonl")
//	if err != nil {
//	    // Handle error
//	}
//	defer file.Close()
//
//	req := files.NewFileUploadRequest(file, "training_data.jsonl", files.PurposeFineTune)
//	uploadedFile, err := client.Files.Upload(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Uploaded file ID: %s\n", uploadedFile.ID)
func (s *FilesService) Upload(ctx context.Context, req *files.FileUploadRequest) (*files.File, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the purpose field
	if err := writer.WriteField("purpose", string(req.Purpose)); err != nil {
		return nil, fmt.Errorf("failed to write purpose field: %w", err)
	}

	// Add the file
	part, err := writer.CreateFormFile("file", req.Filename)
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
	apiResp, err := s.client.PostMultipart(ctx, "/files", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	// Parse the response
	var file files.File
	if err := s.client.ParseJSON(apiResp, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// List retrieves a list of files.
//
// Example:
//
//	fileList, err := client.Files.List(ctx)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, file := range fileList.GetFiles() {
//	    fmt.Printf("File: %s (%s)\n", file.Filename, file.ID)
//	}
func (s *FilesService) List(ctx context.Context) (*files.FileListResponse, error) {
	// Make the API request
	apiResp, err := s.client.Get(ctx, "/files", nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp files.FileListResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Retrieve retrieves information about a specific file.
//
// Example:
//
//	file, err := client.Files.Retrieve(ctx, "file-abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("File: %s, Size: %d bytes\n", file.Filename, file.Bytes)
func (s *FilesService) Retrieve(ctx context.Context, fileID string) (*files.File, error) {
	// Make the API request
	path := fmt.Sprintf("/files/%s", fileID)
	apiResp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var file files.File
	if err := s.client.ParseJSON(apiResp, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// Delete deletes a file.
//
// Example:
//
//	deleteResp, err := client.Files.Delete(ctx, "file-abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	if deleteResp.IsDeleted() {
//	    fmt.Println("File deleted successfully")
//	}
func (s *FilesService) Delete(ctx context.Context, fileID string) (*files.FileDeleteResponse, error) {
	// Make the API request
	path := fmt.Sprintf("/files/%s", fileID)
	apiResp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp files.FileDeleteResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RetrieveContent retrieves the content of a file.
//
// Example:
//
//	content, err := client.Files.RetrieveContent(ctx, "file-abc123")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("File content:\n%s\n", content.String())
func (s *FilesService) RetrieveContent(ctx context.Context, fileID string) (*files.FileContentResponse, error) {
	// Make the API request
	path := fmt.Sprintf("/files/%s/content", fileID)
	apiResp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	// Read the response body
	content, err := io.ReadAll(apiResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return &files.FileContentResponse{
		Content:     content,
		ContentType: apiResp.Headers.Get("Content-Type"),
	}, nil
}
