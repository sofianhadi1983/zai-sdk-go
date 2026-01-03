package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	batchTypes "github.com/sofianhadi1983/zai-sdk-go/api/types/batch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchService_Create(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/batches", r.URL.Path)

		// Parse request body
		var reqBody batchTypes.BatchCreateRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "24h", reqBody.CompletionWindow)
		assert.Equal(t, batchTypes.EndpointChatCompletions, reqBody.Endpoint)
		assert.Equal(t, "file_abc123", reqBody.InputFileID)
		assert.True(t, reqBody.AutoDeleteInputFile)
		assert.NotNil(t, reqBody.Metadata)
		assert.Equal(t, "user_123", reqBody.Metadata["user_id"])

		// Send mock response
		resp := batchTypes.Batch{
			ID:               "batch_xyz789",
			Object:           "batch",
			CompletionWindow: reqBody.CompletionWindow,
			Endpoint:         reqBody.Endpoint,
			InputFileID:      reqBody.InputFileID,
			Status:           batchTypes.StatusValidating,
			CreatedAt:        1700000000,
			Metadata:         reqBody.Metadata,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := batchTypes.NewBatchCreateRequest(
		"24h",
		batchTypes.EndpointChatCompletions,
		"file_abc123",
	).SetMetadata(map[string]string{"user_id": "user_123"})

	batch, err := client.Batch.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, batch)

	assert.Equal(t, "batch_xyz789", batch.ID)
	assert.Equal(t, "batch", batch.Object)
	assert.Equal(t, batchTypes.StatusValidating, batch.Status)
	assert.True(t, batch.IsValidating())
	assert.Equal(t, "file_abc123", batch.InputFileID)
	assert.Equal(t, "user_123", batch.Metadata["user_id"])
}

func TestBatchService_Retrieve(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/batches/batch_abc123", r.URL.Path)

		// Send mock response
		completedAt := int64(1700001000)
		resp := batchTypes.Batch{
			ID:               "batch_abc123",
			Object:           "batch",
			CompletionWindow: "24h",
			Endpoint:         batchTypes.EndpointChatCompletions,
			InputFileID:      "file_input_123",
			Status:           batchTypes.StatusCompleted,
			CreatedAt:        1700000000,
			CompletedAt:      &completedAt,
			OutputFileID:     "file_output_123",
			RequestCounts: &batchTypes.BatchRequestCounts{
				Total:     100,
				Completed: 95,
				Failed:    5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	batch, err := client.Batch.Retrieve(context.Background(), "batch_abc123")
	require.NoError(t, err)
	require.NotNil(t, batch)

	assert.Equal(t, "batch_abc123", batch.ID)
	assert.Equal(t, batchTypes.StatusCompleted, batch.Status)
	assert.True(t, batch.IsCompleted())
	assert.Equal(t, "file_output_123", batch.OutputFileID)
	assert.NotNil(t, batch.RequestCounts)
	assert.Equal(t, 100, batch.RequestCounts.Total)
	assert.Equal(t, 95, batch.RequestCounts.Completed)
	assert.Equal(t, 5, batch.RequestCounts.Failed)
}

func TestBatchService_Retrieve_EmptyID(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithAPIKey("test-key.test-secret"))
	require.NoError(t, err)

	_, err = client.Batch.Retrieve(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "batch ID cannot be empty")
}

func TestBatchService_List(t *testing.T) {
	t.Parallel()

	t.Run("first page", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/batches", r.URL.Path)

			// Verify query parameters
			query := r.URL.Query()
			assert.Equal(t, "20", query.Get("limit"))
			assert.Empty(t, query.Get("after"))

			// Send mock response
			resp := batchTypes.BatchListResponse{
				Object:  "list",
				FirstID: "batch_1",
				LastID:  "batch_2",
				HasMore: true,
				Data: []batchTypes.Batch{
					{
						ID:               "batch_1",
						Object:           "batch",
						CompletionWindow: "24h",
						Endpoint:         batchTypes.EndpointChatCompletions,
						InputFileID:      "file_1",
						Status:           batchTypes.StatusCompleted,
						CreatedAt:        1700000000,
					},
					{
						ID:               "batch_2",
						Object:           "batch",
						CompletionWindow: "24h",
						Endpoint:         batchTypes.EndpointEmbeddings,
						InputFileID:      "file_2",
						Status:           batchTypes.StatusInProgress,
						CreatedAt:        1700000100,
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		listResp, err := client.Batch.List(context.Background(), "", 20)
		require.NoError(t, err)
		require.NotNil(t, listResp)

		assert.Equal(t, "list", listResp.Object)
		assert.Len(t, listResp.Data, 2)
		assert.True(t, listResp.HasMore)
		assert.True(t, listResp.HasMoreBatches())
		assert.Equal(t, "batch_1", listResp.FirstID)
		assert.Equal(t, "batch_2", listResp.LastID)

		batches := listResp.GetBatches()
		assert.Equal(t, "batch_1", batches[0].ID)
		assert.Equal(t, batchTypes.StatusCompleted, batches[0].Status)
		assert.Equal(t, "batch_2", batches[1].ID)
		assert.Equal(t, batchTypes.StatusInProgress, batches[1].Status)
	})

	t.Run("with cursor pagination", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/batches", r.URL.Path)

			// Verify query parameters
			query := r.URL.Query()
			assert.Equal(t, "10", query.Get("limit"))
			assert.Equal(t, "batch_abc", query.Get("after"))

			// Send mock response
			resp := batchTypes.BatchListResponse{
				Object:  "list",
				HasMore: false,
				Data:    []batchTypes.Batch{},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		listResp, err := client.Batch.List(context.Background(), "batch_abc", 10)
		require.NoError(t, err)
		require.NotNil(t, listResp)

		assert.Equal(t, "list", listResp.Object)
		assert.Empty(t, listResp.Data)
		assert.False(t, listResp.HasMore)
	})
}

func TestBatchService_Cancel(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/batches/batch_abc123/cancel", r.URL.Path)

		// Send mock response
		cancellingAt := int64(1700000500)
		resp := batchTypes.Batch{
			ID:               "batch_abc123",
			Object:           "batch",
			CompletionWindow: "24h",
			Endpoint:         batchTypes.EndpointChatCompletions,
			InputFileID:      "file_123",
			Status:           batchTypes.StatusCancelling,
			CreatedAt:        1700000000,
			CancellingAt:     &cancellingAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	batch, err := client.Batch.Cancel(context.Background(), "batch_abc123")
	require.NoError(t, err)
	require.NotNil(t, batch)

	assert.Equal(t, "batch_abc123", batch.ID)
	assert.Equal(t, batchTypes.StatusCancelling, batch.Status)
	assert.True(t, batch.IsCancelling())
	assert.NotNil(t, batch.CancellingAt)
	assert.Equal(t, int64(1700000500), *batch.CancellingAt)
}

func TestBatchService_Cancel_EmptyID(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithAPIKey("test-key.test-secret"))
	require.NoError(t, err)

	_, err = client.Batch.Cancel(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "batch ID cannot be empty")
}

func TestBatchService_APIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Batch not found",
				"type":    "invalid_request_error",
				"code":    "batch_not_found",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	_, err = client.Batch.Retrieve(context.Background(), "nonexistent_batch")
	require.Error(t, err)
}
