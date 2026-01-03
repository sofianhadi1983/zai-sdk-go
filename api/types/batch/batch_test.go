package batch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatch_StatusMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
		checks map[string]bool
	}{
		{
			name:   "validating status",
			status: StatusValidating,
			checks: map[string]bool{
				"IsValidating": true,
				"IsActive":     true,
				"IsTerminal":   false,
			},
		},
		{
			name:   "in_progress status",
			status: StatusInProgress,
			checks: map[string]bool{
				"IsInProgress": true,
				"IsActive":     true,
				"IsTerminal":   false,
			},
		},
		{
			name:   "finalizing status",
			status: StatusFinalizing,
			checks: map[string]bool{
				"IsFinalizing": true,
				"IsActive":     true,
				"IsTerminal":   false,
			},
		},
		{
			name:   "completed status",
			status: StatusCompleted,
			checks: map[string]bool{
				"IsCompleted": true,
				"IsActive":    false,
				"IsTerminal":  true,
			},
		},
		{
			name:   "failed status",
			status: StatusFailed,
			checks: map[string]bool{
				"IsFailed":   true,
				"IsActive":   false,
				"IsTerminal": true,
			},
		},
		{
			name:   "expired status",
			status: StatusExpired,
			checks: map[string]bool{
				"IsExpired":  true,
				"IsActive":   false,
				"IsTerminal": true,
			},
		},
		{
			name:   "cancelling status",
			status: StatusCancelling,
			checks: map[string]bool{
				"IsCancelling": true,
			},
		},
		{
			name:   "cancelled status",
			status: StatusCancelled,
			checks: map[string]bool{
				"IsCancelled": true,
				"IsActive":    false,
				"IsTerminal":  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			batch := &Batch{Status: tt.status}

			if expected, ok := tt.checks["IsValidating"]; ok {
				assert.Equal(t, expected, batch.IsValidating())
			}
			if expected, ok := tt.checks["IsFailed"]; ok {
				assert.Equal(t, expected, batch.IsFailed())
			}
			if expected, ok := tt.checks["IsInProgress"]; ok {
				assert.Equal(t, expected, batch.IsInProgress())
			}
			if expected, ok := tt.checks["IsFinalizing"]; ok {
				assert.Equal(t, expected, batch.IsFinalizing())
			}
			if expected, ok := tt.checks["IsCompleted"]; ok {
				assert.Equal(t, expected, batch.IsCompleted())
			}
			if expected, ok := tt.checks["IsExpired"]; ok {
				assert.Equal(t, expected, batch.IsExpired())
			}
			if expected, ok := tt.checks["IsCancelling"]; ok {
				assert.Equal(t, expected, batch.IsCancelling())
			}
			if expected, ok := tt.checks["IsCancelled"]; ok {
				assert.Equal(t, expected, batch.IsCancelled())
			}
			if expected, ok := tt.checks["IsActive"]; ok {
				assert.Equal(t, expected, batch.IsActive())
			}
			if expected, ok := tt.checks["IsTerminal"]; ok {
				assert.Equal(t, expected, batch.IsTerminal())
			}
		})
	}
}

func TestBatch_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal and unmarshal", func(t *testing.T) {
		t.Parallel()

		createdAt := int64(1609459200)
		completedAt := int64(1609545600)

		batch := &Batch{
			ID:               "batch_123",
			CompletionWindow: "24h",
			CreatedAt:        createdAt,
			Endpoint:         EndpointChatCompletions,
			InputFileID:      "file_input_123",
			Object:           "batch",
			Status:           StatusCompleted,
			CompletedAt:      &completedAt,
			OutputFileID:     "file_output_123",
			RequestCounts: &BatchRequestCounts{
				Total:     100,
				Completed: 95,
				Failed:    5,
			},
		}

		data, err := json.Marshal(batch)
		require.NoError(t, err)

		var decoded Batch
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, batch.ID, decoded.ID)
		assert.Equal(t, batch.CompletionWindow, decoded.CompletionWindow)
		assert.Equal(t, batch.CreatedAt, decoded.CreatedAt)
		assert.Equal(t, batch.Endpoint, decoded.Endpoint)
		assert.Equal(t, batch.InputFileID, decoded.InputFileID)
		assert.Equal(t, batch.Status, decoded.Status)
		assert.Equal(t, batch.CompletedAt, decoded.CompletedAt)
		assert.Equal(t, batch.OutputFileID, decoded.OutputFileID)
		assert.NotNil(t, decoded.RequestCounts)
		assert.Equal(t, batch.RequestCounts.Total, decoded.RequestCounts.Total)
		assert.Equal(t, batch.RequestCounts.Completed, decoded.RequestCounts.Completed)
		assert.Equal(t, batch.RequestCounts.Failed, decoded.RequestCounts.Failed)
	})

	t.Run("unmarshal from API response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "batch_abc123",
			"object": "batch",
			"endpoint": "/v1/chat/completions",
			"input_file_id": "file_xyz",
			"completion_window": "24h",
			"status": "in_progress",
			"created_at": 1700000000,
			"in_progress_at": 1700000100,
			"request_counts": {
				"total": 50,
				"completed": 25,
				"failed": 0
			}
		}`

		var batch Batch
		err := json.Unmarshal([]byte(jsonData), &batch)
		require.NoError(t, err)

		assert.Equal(t, "batch_abc123", batch.ID)
		assert.Equal(t, "batch", batch.Object)
		assert.Equal(t, EndpointChatCompletions, batch.Endpoint)
		assert.Equal(t, "file_xyz", batch.InputFileID)
		assert.Equal(t, "24h", batch.CompletionWindow)
		assert.Equal(t, StatusInProgress, batch.Status)
		assert.Equal(t, int64(1700000000), batch.CreatedAt)
		assert.NotNil(t, batch.InProgressAt)
		assert.Equal(t, int64(1700000100), *batch.InProgressAt)
		assert.NotNil(t, batch.RequestCounts)
		assert.Equal(t, 50, batch.RequestCounts.Total)
		assert.Equal(t, 25, batch.RequestCounts.Completed)
		assert.Equal(t, 0, batch.RequestCounts.Failed)
	})

	t.Run("unmarshal with errors", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "batch_error",
			"object": "batch",
			"endpoint": "/v1/embeddings",
			"input_file_id": "file_xyz",
			"completion_window": "24h",
			"status": "failed",
			"created_at": 1700000000,
			"failed_at": 1700000500,
			"error_file_id": "file_errors",
			"errors": {
				"object": "list",
				"data": [
					{
						"code": "invalid_request",
						"message": "Missing required field",
						"param": "model",
						"line": 10
					}
				]
			}
		}`

		var batch Batch
		err := json.Unmarshal([]byte(jsonData), &batch)
		require.NoError(t, err)

		assert.Equal(t, "batch_error", batch.ID)
		assert.Equal(t, StatusFailed, batch.Status)
		assert.NotNil(t, batch.FailedAt)
		assert.Equal(t, int64(1700000500), *batch.FailedAt)
		assert.Equal(t, "file_errors", batch.ErrorFileID)
		assert.NotNil(t, batch.Errors)
		assert.Equal(t, "list", batch.Errors.Object)
		assert.Len(t, batch.Errors.Data, 1)
		assert.Equal(t, "invalid_request", batch.Errors.Data[0].Code)
		assert.Equal(t, "Missing required field", batch.Errors.Data[0].Message)
		assert.Equal(t, "model", batch.Errors.Data[0].Param)
		assert.Equal(t, 10, batch.Errors.Data[0].Line)
	})
}

func TestBatchCreateRequest(t *testing.T) {
	t.Parallel()

	t.Run("new request", func(t *testing.T) {
		t.Parallel()

		req := NewBatchCreateRequest("24h", EndpointChatCompletions, "file_123")

		assert.Equal(t, "24h", req.CompletionWindow)
		assert.Equal(t, EndpointChatCompletions, req.Endpoint)
		assert.Equal(t, "file_123", req.InputFileID)
		assert.True(t, req.AutoDeleteInputFile)
		assert.Nil(t, req.Metadata)
	})

	t.Run("builder pattern", func(t *testing.T) {
		t.Parallel()

		metadata := map[string]string{
			"user_id": "user_123",
			"batch_name": "test_batch",
		}

		req := NewBatchCreateRequest("24h", EndpointEmbeddings, "file_456").
			SetMetadata(metadata).
			SetAutoDeleteInputFile(false)

		assert.Equal(t, "24h", req.CompletionWindow)
		assert.Equal(t, EndpointEmbeddings, req.Endpoint)
		assert.Equal(t, "file_456", req.InputFileID)
		assert.False(t, req.AutoDeleteInputFile)
		assert.Equal(t, metadata, req.Metadata)
	})

	t.Run("JSON marshaling", func(t *testing.T) {
		t.Parallel()

		req := NewBatchCreateRequest("24h", EndpointChatCompletions, "file_789").
			SetMetadata(map[string]string{"key": "value"})

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded BatchCreateRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, req.CompletionWindow, decoded.CompletionWindow)
		assert.Equal(t, req.Endpoint, decoded.Endpoint)
		assert.Equal(t, req.InputFileID, decoded.InputFileID)
		assert.Equal(t, req.AutoDeleteInputFile, decoded.AutoDeleteInputFile)
		assert.Equal(t, req.Metadata, decoded.Metadata)
	})
}

func TestBatchListResponse(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal list response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"object": "list",
			"data": [
				{
					"id": "batch_1",
					"object": "batch",
					"endpoint": "/v1/chat/completions",
					"input_file_id": "file_1",
					"completion_window": "24h",
					"status": "completed",
					"created_at": 1700000000
				},
				{
					"id": "batch_2",
					"object": "batch",
					"endpoint": "/v1/embeddings",
					"input_file_id": "file_2",
					"completion_window": "24h",
					"status": "in_progress",
					"created_at": 1700000100
				}
			],
			"first_id": "batch_1",
			"last_id": "batch_2",
			"has_more": true
		}`

		var listResp BatchListResponse
		err := json.Unmarshal([]byte(jsonData), &listResp)
		require.NoError(t, err)

		assert.Equal(t, "list", listResp.Object)
		assert.Len(t, listResp.Data, 2)
		assert.Equal(t, "batch_1", listResp.FirstID)
		assert.Equal(t, "batch_2", listResp.LastID)
		assert.True(t, listResp.HasMore)

		batches := listResp.GetBatches()
		assert.Len(t, batches, 2)
		assert.Equal(t, "batch_1", batches[0].ID)
		assert.Equal(t, StatusCompleted, batches[0].Status)
		assert.Equal(t, "batch_2", batches[1].ID)
		assert.Equal(t, StatusInProgress, batches[1].Status)

		assert.True(t, listResp.HasMoreBatches())
	})

	t.Run("empty list response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"object": "list",
			"data": [],
			"has_more": false
		}`

		var listResp BatchListResponse
		err := json.Unmarshal([]byte(jsonData), &listResp)
		require.NoError(t, err)

		assert.Equal(t, "list", listResp.Object)
		assert.Empty(t, listResp.Data)
		assert.False(t, listResp.HasMore)
		assert.False(t, listResp.HasMoreBatches())
	})
}

func TestBatchRequestCounts(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"total": 100,
		"completed": 95,
		"failed": 5
	}`

	var counts BatchRequestCounts
	err := json.Unmarshal([]byte(jsonData), &counts)
	require.NoError(t, err)

	assert.Equal(t, 100, counts.Total)
	assert.Equal(t, 95, counts.Completed)
	assert.Equal(t, 5, counts.Failed)
}

func TestBatchError(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"code": "invalid_request",
		"message": "Invalid model specified",
		"param": "model",
		"line": 42
	}`

	var batchError BatchError
	err := json.Unmarshal([]byte(jsonData), &batchError)
	require.NoError(t, err)

	assert.Equal(t, "invalid_request", batchError.Code)
	assert.Equal(t, "Invalid model specified", batchError.Message)
	assert.Equal(t, "model", batchError.Param)
	assert.Equal(t, 42, batchError.Line)
}

func TestBatchErrors(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"object": "list",
		"data": [
			{
				"code": "invalid_request",
				"message": "Missing field",
				"param": "messages",
				"line": 1
			},
			{
				"code": "rate_limit_exceeded",
				"message": "Too many requests",
				"line": 100
			}
		]
	}`

	var errors BatchErrors
	err := json.Unmarshal([]byte(jsonData), &errors)
	require.NoError(t, err)

	assert.Equal(t, "list", errors.Object)
	assert.Len(t, errors.Data, 2)
	assert.Equal(t, "invalid_request", errors.Data[0].Code)
	assert.Equal(t, "Missing field", errors.Data[0].Message)
	assert.Equal(t, "rate_limit_exceeded", errors.Data[1].Code)
	assert.Equal(t, "Too many requests", errors.Data[1].Message)
}
