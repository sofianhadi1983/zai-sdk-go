// Package fileparser provides types for the File Parser API.
package fileparser

import "io"

// ToolType represents the type of parsing tool to use.
type ToolType string

const (
	// ToolTypeLite is the lite parsing tool.
	ToolTypeLite ToolType = "lite"

	// ToolTypeExpert is the expert parsing tool.
	ToolTypeExpert ToolType = "expert"

	// ToolTypePrime is the prime parsing tool (async).
	ToolTypePrime ToolType = "prime"

	// ToolTypePrimeSync is the prime parsing tool (sync).
	ToolTypePrimeSync ToolType = "prime-sync"
)

// FormatType represents the format of the parsing result.
type FormatType string

const (
	// FormatTypeText returns the result as text.
	FormatTypeText FormatType = "text"

	// FormatTypeDownloadLink returns a download link for the result.
	FormatTypeDownloadLink FormatType = "download_link"
)

// CreateRequest represents a request to create a file parsing task.
type CreateRequest struct {
	// File is the file to parse (required).
	File io.Reader

	// FileName is the name of the file being uploaded (required).
	FileName string

	// FileType specifies the type of file (required).
	FileType string

	// ToolType specifies the parsing tool to use (required).
	ToolType ToolType
}

// NewCreateRequest creates a new file parser create request.
func NewCreateRequest(file io.Reader, fileName, fileType string, toolType ToolType) *CreateRequest {
	return &CreateRequest{
		File:     file,
		FileName: fileName,
		FileType: fileType,
		ToolType: toolType,
	}
}

// CreateResponse represents the response from creating a file parsing task.
type CreateResponse struct {
	// TaskID is the parsing task identifier.
	TaskID string `json:"task_id"`

	// Message is the status message.
	Message string `json:"message"`

	// Success indicates whether the task was created successfully.
	Success bool `json:"success"`
}

// ContentRequest represents a request to get parsing results.
type ContentRequest struct {
	// TaskID is the task identifier (required).
	TaskID string

	// FormatType specifies the format of the result (required).
	FormatType FormatType
}

// NewContentRequest creates a new content request.
func NewContentRequest(taskID string, formatType FormatType) *ContentRequest {
	return &ContentRequest{
		TaskID:     taskID,
		FormatType: formatType,
	}
}

// ContentResponse represents the response from getting parsing content.
// This can be either text content or binary data depending on the format type.
type ContentResponse struct {
	// Content is the parsed content (when FormatType is Text).
	Content string

	// Data is the raw binary data (when FormatType is DownloadLink).
	Data []byte
}

// GetContent returns the content as a string.
func (r *ContentResponse) GetContent() string {
	return r.Content
}

// GetData returns the raw binary data.
func (r *ContentResponse) GetData() []byte {
	return r.Data
}

// HasContent returns true if the response has text content.
func (r *ContentResponse) HasContent() bool {
	return r.Content != ""
}

// HasData returns true if the response has binary data.
func (r *ContentResponse) HasData() bool {
	return len(r.Data) > 0
}

// SyncRequest represents a request to create a synchronous file parsing task.
type SyncRequest struct {
	// File is the file to parse (required).
	File io.Reader

	// FileName is the name of the file being uploaded (required).
	FileName string

	// FileType specifies the type of file (required).
	FileType string

	// ToolType must be ToolTypePrimeSync for sync operations.
	ToolType ToolType
}

// NewSyncRequest creates a new synchronous file parser request.
func NewSyncRequest(file io.Reader, fileName, fileType string) *SyncRequest {
	return &SyncRequest{
		File:     file,
		FileName: fileName,
		FileType: fileType,
		ToolType: ToolTypePrimeSync,
	}
}

// SyncResponse represents the response from a synchronous file parsing task.
type SyncResponse struct {
	// TaskID is the parsing task identifier.
	TaskID string `json:"task_id"`

	// Message is the status message.
	Message string `json:"message"`

	// Status indicates whether the parsing was successful.
	Status bool `json:"status"`

	// Content is the parsed result text content.
	Content string `json:"content"`

	// ParsingResultURL is the download link for the parsed result.
	ParsingResultURL string `json:"parsing_result_url"`
}

// GetContent returns the parsed content.
func (r *SyncResponse) GetContent() string {
	return r.Content
}

// GetDownloadURL returns the download URL for the result.
func (r *SyncResponse) GetDownloadURL() string {
	return r.ParsingResultURL
}

// HasContent returns true if the response has content.
func (r *SyncResponse) HasContent() bool {
	return r.Content != ""
}
