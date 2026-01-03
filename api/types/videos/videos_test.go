package videos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTextToVideoRequest(t *testing.T) {
	t.Parallel()

	req := NewTextToVideoRequest(ModelCogVideoX, "A cat playing with a ball")

	assert.Equal(t, ModelCogVideoX, req.Model)
	assert.Equal(t, "A cat playing with a ball", req.Prompt)
	assert.Empty(t, req.ImageURL)
}

func TestNewImageToVideoRequest(t *testing.T) {
	t.Parallel()

	req := NewImageToVideoRequest(ModelCogVideoX, "https://example.com/image.jpg")

	assert.Equal(t, ModelCogVideoX, req.Model)
	assert.Equal(t, "https://example.com/image.jpg", req.ImageURL)
	assert.Empty(t, req.Prompt)
}

func TestVideoGenerationRequest_SetUser(t *testing.T) {
	t.Parallel()

	req := &VideoGenerationRequest{}
	req.SetUser("user-123")

	assert.Equal(t, "user-123", req.User)
}

func TestVideoGenerationRequest_Chaining(t *testing.T) {
	t.Parallel()

	req := NewTextToVideoRequest(ModelCogVideoX, "test prompt")
	req.SetUser("user-456")

	assert.Equal(t, ModelCogVideoX, req.Model)
	assert.Equal(t, "test prompt", req.Prompt)
	assert.Equal(t, "user-456", req.User)
}

func TestVideoGenerationRequest_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal text-to-video request", func(t *testing.T) {
		t.Parallel()

		req := NewTextToVideoRequest(ModelCogVideoX, "A sunset")
		req.SetUser("user-789")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "cogvideox")
		assert.Contains(t, string(data), "A sunset")
		assert.Contains(t, string(data), "user-789")
	})

	t.Run("marshal image-to-video request", func(t *testing.T) {
		t.Parallel()

		req := NewImageToVideoRequest(ModelCogVideoX, "https://example.com/img.jpg")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "cogvideox")
		assert.Contains(t, string(data), "https://example.com/img.jpg")
		assert.NotContains(t, string(data), "prompt")
	})

	t.Run("omit empty fields", func(t *testing.T) {
		t.Parallel()

		req := &VideoGenerationRequest{
			Model:  ModelCogVideoX,
			Prompt: "test",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		// ImageURL and User should be omitted
		assert.NotContains(t, string(data), "image_url")
		assert.NotContains(t, string(data), "user")
	})
}

func TestVideoTask_StatusChecks(t *testing.T) {
	t.Parallel()

	t.Run("submitted status", func(t *testing.T) {
		t.Parallel()

		task := &VideoTask{Status: StatusSubmitted}
		assert.True(t, task.IsSubmitted())
		assert.False(t, task.IsProcessing())
		assert.False(t, task.IsCompleted())
		assert.False(t, task.IsFailed())
	})

	t.Run("processing status", func(t *testing.T) {
		t.Parallel()

		task := &VideoTask{Status: StatusProcessing}
		assert.False(t, task.IsSubmitted())
		assert.True(t, task.IsProcessing())
		assert.False(t, task.IsCompleted())
		assert.False(t, task.IsFailed())
	})

	t.Run("completed status", func(t *testing.T) {
		t.Parallel()

		task := &VideoTask{Status: StatusCompleted}
		assert.False(t, task.IsSubmitted())
		assert.False(t, task.IsProcessing())
		assert.True(t, task.IsCompleted())
		assert.False(t, task.IsFailed())
	})

	t.Run("failed status", func(t *testing.T) {
		t.Parallel()

		task := &VideoTask{Status: StatusFailed}
		assert.False(t, task.IsSubmitted())
		assert.False(t, task.IsProcessing())
		assert.False(t, task.IsCompleted())
		assert.True(t, task.IsFailed())
	})
}

func TestVideoGenerationResponse_GetTaskID(t *testing.T) {
	t.Parallel()

	resp := &VideoGenerationResponse{
		ID: "task-abc123",
	}

	assert.Equal(t, "task-abc123", resp.GetTaskID())
}

func TestVideoGenerationResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "task-xyz789",
			"model": "cogvideox",
			"request_id": "req-123"
		}`

		var resp VideoGenerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "task-xyz789", resp.ID)
		assert.Equal(t, ModelCogVideoX, resp.Model)
		assert.Equal(t, "req-123", resp.RequestID)
	})

	t.Run("marshal response", func(t *testing.T) {
		t.Parallel()

		resp := &VideoGenerationResponse{
			ID:        "task-abc",
			Model:     ModelCogVideoX,
			RequestID: "req-456",
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		assert.Contains(t, string(data), "task-abc")
		assert.Contains(t, string(data), "cogvideox")
		assert.Contains(t, string(data), "req-456")
	})
}

func TestVideoResult_StatusChecks(t *testing.T) {
	t.Parallel()

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{TaskStatus: StatusCompleted}
		assert.True(t, result.IsCompleted())
		assert.False(t, result.IsFailed())
		assert.False(t, result.IsProcessing())
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{TaskStatus: StatusFailed}
		assert.False(t, result.IsCompleted())
		assert.True(t, result.IsFailed())
		assert.False(t, result.IsProcessing())
	})

	t.Run("processing", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{TaskStatus: StatusProcessing}
		assert.False(t, result.IsCompleted())
		assert.False(t, result.IsFailed())
		assert.True(t, result.IsProcessing())
	})

	t.Run("submitted", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{TaskStatus: StatusSubmitted}
		assert.False(t, result.IsCompleted())
		assert.False(t, result.IsFailed())
		assert.True(t, result.IsProcessing())
	})
}

func TestVideoResult_GetFirstVideo(t *testing.T) {
	t.Parallel()

	t.Run("with videos", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{
				{URL: "https://example.com/video1.mp4"},
				{URL: "https://example.com/video2.mp4"},
			},
		}

		video := result.GetFirstVideo()
		require.NotNil(t, video)
		assert.Equal(t, "https://example.com/video1.mp4", video.URL)
	})

	t.Run("without videos", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{},
		}

		video := result.GetFirstVideo()
		assert.Nil(t, video)
	})
}

func TestVideoResult_GetVideoURL(t *testing.T) {
	t.Parallel()

	t.Run("with video", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{
				{URL: "https://example.com/video.mp4"},
			},
		}

		assert.Equal(t, "https://example.com/video.mp4", result.GetVideoURL())
	})

	t.Run("without video", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{},
		}

		assert.Empty(t, result.GetVideoURL())
	})
}

func TestVideoResult_GetCoverImageURL(t *testing.T) {
	t.Parallel()

	t.Run("with cover image", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{
				{
					URL:           "https://example.com/video.mp4",
					CoverImageURL: "https://example.com/cover.jpg",
				},
			},
		}

		assert.Equal(t, "https://example.com/cover.jpg", result.GetCoverImageURL())
	})

	t.Run("without cover image", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{},
		}

		assert.Empty(t, result.GetCoverImageURL())
	})
}

func TestVideoResult_GetAllVideoURLs(t *testing.T) {
	t.Parallel()

	t.Run("multiple videos", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{
				{URL: "https://example.com/v1.mp4"},
				{URL: "https://example.com/v2.mp4"},
				{URL: "https://example.com/v3.mp4"},
			},
		}

		urls := result.GetAllVideoURLs()
		require.Len(t, urls, 3)
		assert.Equal(t, "https://example.com/v1.mp4", urls[0])
		assert.Equal(t, "https://example.com/v2.mp4", urls[1])
		assert.Equal(t, "https://example.com/v3.mp4", urls[2])
	})

	t.Run("no videos", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{},
		}

		urls := result.GetAllVideoURLs()
		assert.Len(t, urls, 0)
	})

	t.Run("mixed empty URLs", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			VideoResult: []VideoData{
				{URL: "https://example.com/v1.mp4"},
				{URL: ""},
				{URL: "https://example.com/v2.mp4"},
			},
		}

		urls := result.GetAllVideoURLs()
		require.Len(t, urls, 2)
		assert.Equal(t, "https://example.com/v1.mp4", urls[0])
		assert.Equal(t, "https://example.com/v2.mp4", urls[1])
	})
}

func TestVideoResult_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("has error", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			TaskStatus:   StatusFailed,
			ErrorMessage: "Generation failed",
		}

		assert.True(t, result.HasError())
		assert.Equal(t, "Generation failed", result.GetError())
	})

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			TaskStatus: StatusCompleted,
		}

		assert.False(t, result.HasError())
		assert.Empty(t, result.GetError())
	})
}

func TestVideoResult_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal completed result", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"task_id": "task-123",
			"task_status": "completed",
			"request_id": "req-456",
			"video_result": [
				{
					"url": "https://example.com/video.mp4",
					"cover_image_url": "https://example.com/cover.jpg"
				}
			]
		}`

		var result VideoResult
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Equal(t, "task-123", result.TaskID)
		assert.Equal(t, StatusCompleted, result.TaskStatus)
		assert.Equal(t, "req-456", result.RequestID)
		assert.Len(t, result.VideoResult, 1)
		assert.Equal(t, "https://example.com/video.mp4", result.VideoResult[0].URL)
		assert.Equal(t, "https://example.com/cover.jpg", result.VideoResult[0].CoverImageURL)
	})

	t.Run("unmarshal failed result", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"task_id": "task-789",
			"task_status": "failed",
			"error_message": "Invalid prompt"
		}`

		var result VideoResult
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Equal(t, "task-789", result.TaskID)
		assert.Equal(t, StatusFailed, result.TaskStatus)
		assert.Equal(t, "Invalid prompt", result.ErrorMessage)
		assert.True(t, result.IsFailed())
		assert.True(t, result.HasError())
	})

	t.Run("marshal result", func(t *testing.T) {
		t.Parallel()

		result := &VideoResult{
			TaskID:     "task-abc",
			TaskStatus: StatusCompleted,
			VideoResult: []VideoData{
				{URL: "https://example.com/v.mp4"},
			},
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		assert.Contains(t, string(data), "task-abc")
		assert.Contains(t, string(data), "completed")
		assert.Contains(t, string(data), "https://example.com/v.mp4")
	})
}

func TestVideoData_Helpers(t *testing.T) {
	t.Parallel()

	video := &VideoData{
		URL:           "https://example.com/video.mp4",
		CoverImageURL: "https://example.com/cover.jpg",
	}

	assert.Equal(t, "https://example.com/video.mp4", video.GetURL())
	assert.Equal(t, "https://example.com/cover.jpg", video.GetCoverImageURL())
}

func TestVideoModel_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, VideoModel("cogvideox"), ModelCogVideoX)
}

func TestTaskStatus_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, TaskStatus("submitted"), StatusSubmitted)
	assert.Equal(t, TaskStatus("processing"), StatusProcessing)
	assert.Equal(t, TaskStatus("completed"), StatusCompleted)
	assert.Equal(t, TaskStatus("failed"), StatusFailed)
}

func TestVideoGenerationRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Text-to-video with all options
	req := NewTextToVideoRequest(ModelCogVideoX, "A futuristic cityscape")
	req.SetUser("user-example-123")

	assert.Equal(t, ModelCogVideoX, req.Model)
	assert.Equal(t, "A futuristic cityscape", req.Prompt)
	assert.Equal(t, "user-example-123", req.User)

	// Ensure it can be marshaled
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), "cogvideox")
	assert.Contains(t, string(data), "futuristic cityscape")
}

func TestVideoResult_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response
	jsonData := `{
		"task_id": "task-real-world-123",
		"task_status": "completed",
		"request_id": "req-987654",
		"video_result": [
			{
				"url": "https://cdn.example.com/videos/generated-abc123.mp4",
				"cover_image_url": "https://cdn.example.com/covers/cover-abc123.jpg"
			}
		]
	}`

	var result VideoResult
	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "task-real-world-123", result.TaskID)
	assert.Equal(t, StatusCompleted, result.TaskStatus)
	assert.True(t, result.IsCompleted())
	assert.False(t, result.IsFailed())
	assert.False(t, result.HasError())

	// Verify helper methods work
	firstVideo := result.GetFirstVideo()
	require.NotNil(t, firstVideo)
	assert.Contains(t, firstVideo.URL, "generated-abc123.mp4")
	assert.Contains(t, firstVideo.CoverImageURL, "cover-abc123.jpg")

	assert.Equal(t, firstVideo.URL, result.GetVideoURL())
	assert.Equal(t, firstVideo.CoverImageURL, result.GetCoverImageURL())

	allURLs := result.GetAllVideoURLs()
	require.Len(t, allURLs, 1)
	assert.Equal(t, firstVideo.URL, allURLs[0])
}
