// Package files provides types for the Files API.
package files

import "io"

// FilePurpose represents the intended purpose of a file.
type FilePurpose string

const (
	// PurposeAssistants is for files used with Assistants.
	PurposeAssistants FilePurpose = "assistants"
	// PurposeFineTune is for files used for fine-tuning.
	PurposeFineTune FilePurpose = "fine-tune"
	// PurposeBatch is for files used for batch processing.
	PurposeBatch FilePurpose = "batch"
)

// FileStatus represents the status of a file.
type FileStatus string

const (
	// StatusUploaded indicates the file has been successfully uploaded.
	StatusUploaded FileStatus = "uploaded"
	// StatusProcessed indicates the file has been processed.
	StatusProcessed FileStatus = "processed"
	// StatusError indicates there was an error processing the file.
	StatusError FileStatus = "error"
)

// File represents a file that has been uploaded to the API.
type File struct {
	// ID is the unique identifier for the file.
	ID string `json:"id"`

	// Object is the object type, which is always "file".
	Object string `json:"object"`

	// Bytes is the size of the file in bytes.
	Bytes int64 `json:"bytes"`

	// CreatedAt is the Unix timestamp when the file was created.
	CreatedAt int64 `json:"created_at"`

	// Filename is the name of the file.
	Filename string `json:"filename"`

	// Purpose is the intended purpose of the file.
	Purpose FilePurpose `json:"purpose"`

	// Status is the current status of the file.
	Status FileStatus `json:"status,omitempty"`

	// StatusDetails provides additional information about the file status.
	StatusDetails string `json:"status_details,omitempty"`
}

// FileUploadRequest represents a request to upload a file.
type FileUploadRequest struct {
	// File is the file content to upload.
	File io.Reader

	// Filename is the name of the file.
	Filename string

	// Purpose is the intended purpose of the file.
	Purpose FilePurpose
}

// NewFileUploadRequest creates a new file upload request.
//
// Example:
//
//	file, _ := os.Open("data.jsonl")
//	req := files.NewFileUploadRequest(file, "data.jsonl", files.PurposeFineTune)
func NewFileUploadRequest(file io.Reader, filename string, purpose FilePurpose) *FileUploadRequest {
	return &FileUploadRequest{
		File:     file,
		Filename: filename,
		Purpose:  purpose,
	}
}

// FileListResponse represents a list of files.
type FileListResponse struct {
	// Object is the object type, which is always "list".
	Object string `json:"object"`

	// Data is the list of files.
	Data []File `json:"data"`

	// HasMore indicates if there are more files available.
	HasMore bool `json:"has_more,omitempty"`
}

// FileDeleteResponse represents the response when deleting a file.
type FileDeleteResponse struct {
	// ID is the ID of the deleted file.
	ID string `json:"id"`

	// Object is the object type, which is always "file".
	Object string `json:"object"`

	// Deleted indicates whether the file was successfully deleted.
	Deleted bool `json:"deleted"`
}

// FileContentResponse represents the content of a file.
type FileContentResponse struct {
	// Content is the file content.
	Content []byte

	// ContentType is the MIME type of the file.
	ContentType string
}

// GetID returns the file ID.
func (f *File) GetID() string {
	return f.ID
}

// GetFilename returns the file name.
func (f *File) GetFilename() string {
	return f.Filename
}

// GetSize returns the file size in bytes.
func (f *File) GetSize() int64 {
	return f.Bytes
}

// GetPurpose returns the file purpose.
func (f *File) GetPurpose() FilePurpose {
	return f.Purpose
}

// IsUploaded returns true if the file has been successfully uploaded.
func (f *File) IsUploaded() bool {
	return f.Status == StatusUploaded || f.Status == StatusProcessed
}

// HasError returns true if there was an error processing the file.
func (f *File) HasError() bool {
	return f.Status == StatusError
}

// GetFiles returns all files from the list response.
func (r *FileListResponse) GetFiles() []File {
	return r.Data
}

// GetFileIDs returns all file IDs from the list response.
func (r *FileListResponse) GetFileIDs() []string {
	ids := make([]string, len(r.Data))
	for i, file := range r.Data {
		ids[i] = file.ID
	}
	return ids
}

// GetFileByID finds a file by ID in the list response.
func (r *FileListResponse) GetFileByID(id string) *File {
	for i := range r.Data {
		if r.Data[i].ID == id {
			return &r.Data[i]
		}
	}
	return nil
}

// GetFilesByPurpose returns all files with the specified purpose.
func (r *FileListResponse) GetFilesByPurpose(purpose FilePurpose) []File {
	files := make([]File, 0)
	for _, file := range r.Data {
		if file.Purpose == purpose {
			files = append(files, file)
		}
	}
	return files
}

// IsDeleted returns true if the file was successfully deleted.
func (r *FileDeleteResponse) IsDeleted() bool {
	return r.Deleted
}

// GetContent returns the file content.
func (r *FileContentResponse) GetContent() []byte {
	return r.Content
}

// GetContentType returns the content type.
func (r *FileContentResponse) GetContentType() string {
	return r.ContentType
}

// String returns the file content as a string.
func (r *FileContentResponse) String() string {
	return string(r.Content)
}
