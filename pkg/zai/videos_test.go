package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	videostypes "github.com/z-ai/zai-sdk-go/api/types/videos"
)

func TestVideosService_Create(t *testing.T) {
	t.Parallel()

	t.Run("text-to-video request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/videos/generations", r.URL.Path)
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

			// Verify request body
			var req videostypes.VideoGenerationRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, videostypes.ModelCogVideoX, req.Model)
			assert.Equal(t, "A cat playing", req.Prompt)
			assert.Empty(t, req.ImageURL)

			// Send response
			resp := videostypes.VideoGenerationResponse{
				ID:        "task-abc123",
				Model:     videostypes.ModelCogVideoX,
				RequestID: "req-456",
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
		defer client.Close()

		req := videostypes.NewTextToVideoRequest(videostypes.ModelCogVideoX, "A cat playing")
		task, err := client.Videos.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, task)

		assert.Equal(t, "task-abc123", task.ID)
		assert.Equal(t, videostypes.ModelCogVideoX, task.Model)
		assert.Equal(t, "req-456", task.RequestID)
	})

	t.Run("image-to-video request", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req videostypes.VideoGenerationRequest
			json.NewDecoder(r.Body).Decode(&req)

			assert.Equal(t, "https://example.com/image.jpg", req.ImageURL)
			assert.Empty(t, req.Prompt)

			resp := videostypes.VideoGenerationResponse{
				ID:    "task-xyz789",
				Model: videostypes.ModelCogVideoX,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := videostypes.NewImageToVideoRequest(videostypes.ModelCogVideoX, "https://example.com/image.jpg")
		task, err := client.Videos.Create(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "task-xyz789", task.ID)
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid prompt",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := videostypes.NewTextToVideoRequest(videostypes.ModelCogVideoX, "")
		task, err := client.Videos.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "Invalid prompt")
	})
}

func TestVideosService_Retrieve(t *testing.T) {
	t.Parallel()

	t.Run("completed task", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/async-result/task-abc123", r.URL.Path)

			result := videostypes.VideoResult{
				TaskID:     "task-abc123",
				TaskStatus: videostypes.StatusCompleted,
				VideoResult: []videostypes.VideoData{
					{
						URL:           "https://example.com/video.mp4",
						CoverImageURL: "https://example.com/cover.jpg",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.Retrieve(context.Background(), "task-abc123")
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "task-abc123", result.TaskID)
		assert.True(t, result.IsCompleted())
		assert.Equal(t, "https://example.com/video.mp4", result.GetVideoURL())
		assert.Equal(t, "https://example.com/cover.jpg", result.GetCoverImageURL())
	})

	t.Run("processing task", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := videostypes.VideoResult{
				TaskID:     "task-xyz",
				TaskStatus: videostypes.StatusProcessing,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.Retrieve(context.Background(), "task-xyz")
		require.NoError(t, err)
		assert.True(t, result.IsProcessing())
		assert.False(t, result.IsCompleted())
	})

	t.Run("failed task", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := videostypes.VideoResult{
				TaskID:       "task-failed",
				TaskStatus:   videostypes.StatusFailed,
				ErrorMessage: "Generation failed",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.Retrieve(context.Background(), "task-failed")
		require.NoError(t, err)
		assert.True(t, result.IsFailed())
		assert.True(t, result.HasError())
		assert.Equal(t, "Generation failed", result.GetError())
	})

	t.Run("task not found", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Task not found",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.Retrieve(context.Background(), "task-nonexistent")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Task not found")
	})
}

func TestVideosService_GenerateText(t *testing.T) {
	t.Parallel()

	t.Run("successful generation", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := videostypes.VideoGenerationResponse{
				ID: "task-text-123",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		taskID, err := client.Videos.GenerateText(
			context.Background(),
			videostypes.ModelCogVideoX,
			"A beautiful landscape",
		)
		require.NoError(t, err)
		assert.Equal(t, "task-text-123", taskID)
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid API key",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		taskID, err := client.Videos.GenerateText(
			context.Background(),
			videostypes.ModelCogVideoX,
			"test",
		)
		assert.Error(t, err)
		assert.Empty(t, taskID)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}

func TestVideosService_GenerateFromImage(t *testing.T) {
	t.Parallel()

	t.Run("successful generation", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := videostypes.VideoGenerationResponse{
				ID: "task-image-456",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		taskID, err := client.Videos.GenerateFromImage(
			context.Background(),
			videostypes.ModelCogVideoX,
			"https://example.com/image.jpg",
		)
		require.NoError(t, err)
		assert.Equal(t, "task-image-456", taskID)
	})
}

func TestVideosService_WaitForCompletion(t *testing.T) {
	t.Parallel()

	t.Run("completes successfully", func(t *testing.T) {
		t.Parallel()

		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++

			var result videostypes.VideoResult

			if callCount < 3 {
				// First two calls: still processing
				result = videostypes.VideoResult{
					TaskID:     "task-wait",
					TaskStatus: videostypes.StatusProcessing,
				}
			} else {
				// Third call: completed
				result = videostypes.VideoResult{
					TaskID:     "task-wait",
					TaskStatus: videostypes.StatusCompleted,
					VideoResult: []videostypes.VideoData{
						{URL: "https://example.com/video.mp4"},
					},
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.WaitForCompletion(
			context.Background(),
			"task-wait",
			100*time.Millisecond,
			5*time.Second,
		)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.True(t, result.IsCompleted())
		assert.Equal(t, "https://example.com/video.mp4", result.GetVideoURL())
		assert.GreaterOrEqual(t, callCount, 3)
	})

	t.Run("fails", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := videostypes.VideoResult{
				TaskID:       "task-fail",
				TaskStatus:   videostypes.StatusFailed,
				ErrorMessage: "Generation failed",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.WaitForCompletion(
			context.Background(),
			"task-fail",
			100*time.Millisecond,
			5*time.Second,
		)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.True(t, result.IsFailed())
		assert.Equal(t, "Generation failed", result.GetError())
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := videostypes.VideoResult{
				TaskID:     "task-timeout",
				TaskStatus: videostypes.StatusProcessing,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		result, err := client.Videos.WaitForCompletion(
			context.Background(),
			"task-timeout",
			100*time.Millisecond,
			300*time.Millisecond, // Short timeout for testing
		)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("context canceled", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := videostypes.VideoResult{
				TaskID:     "task-cancel",
				TaskStatus: videostypes.StatusProcessing,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := client.Videos.WaitForCompletion(
			ctx,
			"task-cancel",
			100*time.Millisecond,
			5*time.Second,
		)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestClient_VideosService_Integration(t *testing.T) {
	t.Parallel()

	t.Run("client has videos service", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.Videos)
	})

	t.Run("complete workflow", func(t *testing.T) {
		t.Parallel()

		retrieveCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/videos/generations":
				// Create task
				resp := videostypes.VideoGenerationResponse{
					ID:    "task-workflow",
					Model: videostypes.ModelCogVideoX,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)

			case "/async-result/task-workflow":
				// Retrieve status
				retrieveCount++

				var result videostypes.VideoResult
				if retrieveCount < 2 {
					result = videostypes.VideoResult{
						TaskID:     "task-workflow",
						TaskStatus: videostypes.StatusProcessing,
					}
				} else {
					result = videostypes.VideoResult{
						TaskID:     "task-workflow",
						TaskStatus: videostypes.StatusCompleted,
						VideoResult: []videostypes.VideoData{
							{
								URL:           "https://example.com/workflow.mp4",
								CoverImageURL: "https://example.com/cover.jpg",
							},
						},
					}
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			}
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		ctx := context.Background()

		// Create a video generation task
		taskID, err := client.Videos.GenerateText(ctx, videostypes.ModelCogVideoX, "A sunset")
		require.NoError(t, err)
		assert.Equal(t, "task-workflow", taskID)

		// Wait for completion
		result, err := client.Videos.WaitForCompletion(ctx, taskID, 100*time.Millisecond, 5*time.Second)
		require.NoError(t, err)
		assert.True(t, result.IsCompleted())
		assert.NotEmpty(t, result.GetVideoURL())
	})
}
