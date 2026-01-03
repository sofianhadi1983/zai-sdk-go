// Package batch provides types for the Batch API.
package batch

// BatchError represents an error that occurred during batch processing.
type BatchError struct {
	// Code is the defined business error code
	Code string `json:"code,omitempty"`

	// Line is the line number in the file where the error occurred
	Line int `json:"line,omitempty"`

	// Message is the description of the error in the conversation file
	Message string `json:"message,omitempty"`

	// Param is the parameter that caused the error
	Param string `json:"param,omitempty"`
}

// BatchErrors represents batch errors information.
type BatchErrors struct {
	// Data is the list of batch errors
	Data []BatchError `json:"data,omitempty"`

	// Object is the type identifier, always "list"
	Object string `json:"object,omitempty"`
}

// BatchRequestCounts represents request counts for different states in the batch.
type BatchRequestCounts struct {
	// Completed is the number of requests that have been completed
	Completed int `json:"completed"`

	// Failed is the number of failed requests
	Failed int `json:"failed"`

	// Total is the total number of requests
	Total int `json:"total"`
}

// Batch represents a batch processing object.
type Batch struct {
	// ID is the batch identifier
	ID string `json:"id"`

	// CompletionWindow is the time frame within which the batch should be processed
	// Currently only "24h" is supported
	CompletionWindow string `json:"completion_window"`

	// CreatedAt is the creation time represented by the Unix timestamp (in seconds)
	CreatedAt int64 `json:"created_at"`

	// Endpoint is the address of the Z.ai endpoint
	// Currently "/v1/chat/completions" and "/v1/embeddings" are supported
	Endpoint string `json:"endpoint"`

	// InputFileID is the ID of the input file marked as batch
	InputFileID string `json:"input_file_id"`

	// Object is the type identifier, always "batch"
	Object string `json:"object"`

	// Status is the status of the batch
	Status string `json:"status"`

	// CancelledAt is the cancellation time represented by the Unix timestamp (in seconds)
	CancelledAt *int64 `json:"cancelled_at,omitempty"`

	// CancellingAt is the time when the cancellation request was initiated
	CancellingAt *int64 `json:"cancelling_at,omitempty"`

	// CompletedAt is the completion time represented by the Unix timestamp (in seconds)
	CompletedAt *int64 `json:"completed_at,omitempty"`

	// ErrorFileID contains the output of the request that failed to be executed
	ErrorFileID string `json:"error_file_id,omitempty"`

	// Errors is the batch errors information
	Errors *BatchErrors `json:"errors,omitempty"`

	// ExpiredAt is the expiration time represented by the Unix timestamp (in seconds)
	ExpiredAt *int64 `json:"expired_at,omitempty"`

	// ExpiresAt is the expiration trigger time represented by the Unix timestamp (in seconds)
	ExpiresAt *int64 `json:"expires_at,omitempty"`

	// FailedAt is the failure time represented by the Unix timestamp (in seconds)
	FailedAt *int64 `json:"failed_at,omitempty"`

	// FinalizingAt is the final time represented by the Unix timestamp (in seconds)
	FinalizingAt *int64 `json:"finalizing_at,omitempty"`

	// InProgressAt is the start processing time represented by the Unix timestamp (in seconds)
	InProgressAt *int64 `json:"in_progress_at,omitempty"`

	// Metadata is optional metadata in key:value format to store information in a structured format
	// The key length is 64 characters, and the value is up to 512 characters long
	Metadata map[string]string `json:"metadata,omitempty"`

	// OutputFileID is the ID of the output file for the completed request
	OutputFileID string `json:"output_file_id,omitempty"`

	// RequestCounts is the request count for different states in the batch
	RequestCounts *BatchRequestCounts `json:"request_counts,omitempty"`
}

// Batch status constants
const (
	StatusValidating  = "validating"
	StatusFailed      = "failed"
	StatusInProgress  = "in_progress"
	StatusFinalizing  = "finalizing"
	StatusCompleted   = "completed"
	StatusExpired     = "expired"
	StatusCancelling  = "cancelling"
	StatusCancelled   = "cancelled"
)

// Endpoint constants
const (
	EndpointChatCompletions = "/v1/chat/completions"
	EndpointEmbeddings      = "/v1/embeddings"
)

// IsValidating returns true if the batch is in validating status.
func (b *Batch) IsValidating() bool {
	return b.Status == StatusValidating
}

// IsFailed returns true if the batch has failed.
func (b *Batch) IsFailed() bool {
	return b.Status == StatusFailed
}

// IsInProgress returns true if the batch is being processed.
func (b *Batch) IsInProgress() bool {
	return b.Status == StatusInProgress
}

// IsFinalizing returns true if the batch is finalizing.
func (b *Batch) IsFinalizing() bool {
	return b.Status == StatusFinalizing
}

// IsCompleted returns true if the batch has completed successfully.
func (b *Batch) IsCompleted() bool {
	return b.Status == StatusCompleted
}

// IsExpired returns true if the batch has expired.
func (b *Batch) IsExpired() bool {
	return b.Status == StatusExpired
}

// IsCancelling returns true if the batch is being cancelled.
func (b *Batch) IsCancelling() bool {
	return b.Status == StatusCancelling
}

// IsCancelled returns true if the batch has been cancelled.
func (b *Batch) IsCancelled() bool {
	return b.Status == StatusCancelled
}

// IsActive returns true if the batch is in an active state (validating, in_progress, or finalizing).
func (b *Batch) IsActive() bool {
	return b.IsValidating() || b.IsInProgress() || b.IsFinalizing()
}

// IsTerminal returns true if the batch is in a terminal state (completed, failed, expired, or cancelled).
func (b *Batch) IsTerminal() bool {
	return b.IsCompleted() || b.IsFailed() || b.IsExpired() || b.IsCancelled()
}

// BatchCreateRequest represents a request to create a new batch.
type BatchCreateRequest struct {
	// CompletionWindow is the time frame within which the batch should be processed
	// Currently only "24h" is supported
	CompletionWindow string `json:"completion_window"`

	// Endpoint is the endpoint to be used for all requests in the batch
	// Currently "/v1/chat/completions" and "/v1/embeddings" are supported
	Endpoint string `json:"endpoint"`

	// InputFileID is the ID of an uploaded file that contains requests for the new batch
	// The input file must be formatted as a JSONL file and uploaded with purpose "batch"
	InputFileID string `json:"input_file_id"`

	// Metadata is optional custom metadata for the batch
	Metadata map[string]string `json:"metadata,omitempty"`

	// AutoDeleteInputFile indicates whether to automatically delete the input file after processing
	AutoDeleteInputFile bool `json:"auto_delete_input_file,omitempty"`
}

// NewBatchCreateRequest creates a new batch create request.
//
// Example:
//
//	req := batch.NewBatchCreateRequest(
//	    "24h",
//	    batch.EndpointChatCompletions,
//	    "file_123",
//	)
func NewBatchCreateRequest(completionWindow, endpoint, inputFileID string) *BatchCreateRequest {
	return &BatchCreateRequest{
		CompletionWindow:    completionWindow,
		Endpoint:            endpoint,
		InputFileID:         inputFileID,
		AutoDeleteInputFile: true,
	}
}

// SetMetadata sets the metadata for the batch.
func (r *BatchCreateRequest) SetMetadata(metadata map[string]string) *BatchCreateRequest {
	r.Metadata = metadata
	return r
}

// SetAutoDeleteInputFile sets whether to automatically delete the input file.
func (r *BatchCreateRequest) SetAutoDeleteInputFile(autoDelete bool) *BatchCreateRequest {
	r.AutoDeleteInputFile = autoDelete
	return r
}

// BatchListResponse represents the response from listing batches.
type BatchListResponse struct {
	// Data is the list of batch objects
	Data []Batch `json:"data"`

	// Object is the type identifier, always "list"
	Object string `json:"object"`

	// FirstID is the ID of the first batch in the list
	FirstID string `json:"first_id,omitempty"`

	// LastID is the ID of the last batch in the list
	LastID string `json:"last_id,omitempty"`

	// HasMore indicates whether there are more batches available
	HasMore bool `json:"has_more"`
}

// GetBatches returns the list of batches.
func (r *BatchListResponse) GetBatches() []Batch {
	return r.Data
}

// HasMoreBatches returns whether there are more batches available.
func (r *BatchListResponse) HasMoreBatches() bool {
	return r.HasMore
}
