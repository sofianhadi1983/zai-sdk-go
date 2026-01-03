package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	imagestypes "github.com/sofianhadi1983/zai-sdk-go/api/types/images"
)

func TestImagesService_Create(t *testing.T) {
	t.Parallel()

	t.Run("successful image generation", func(t *testing.T) {
		t.Parallel()

		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/images/generations", r.URL.Path)
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			var req imagestypes.ImageGenerationRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "cogview-3", req.Model)
			assert.Equal(t, "A beautiful sunset", req.Prompt)

			// Send response
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{
						URL:           "https://example.com/image.png",
						RevisedPrompt: "A beautiful sunset over the ocean",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		// Create client
		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		// Make request
		req := imagestypes.NewImageGenerationRequest("cogview-3", "A beautiful sunset")

		resp, err := client.Images.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, int64(1677652288), resp.Created)
		assert.Len(t, resp.Data, 1)

		img := resp.GetFirstImage()
		require.NotNil(t, img)
		assert.Equal(t, "https://example.com/image.png", img.URL)
		assert.Equal(t, "A beautiful sunset over the ocean", img.RevisedPrompt)
	})

	t.Run("with custom parameters", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req imagestypes.ImageGenerationRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Verify custom parameters
			assert.Equal(t, imagestypes.Size1792x1024, req.Size)
			assert.Equal(t, imagestypes.QualityHD, req.Quality)
			assert.NotNil(t, req.N)
			assert.Equal(t, 2, *req.N)
			assert.Equal(t, imagestypes.ResponseFormatURL, req.ResponseFormat)

			// Send response with multiple images
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{URL: "https://example.com/1.png"},
					{URL: "https://example.com/2.png"},
				},
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

		req := imagestypes.NewImageGenerationRequest("cogview-3", "test")
		req.SetSize(imagestypes.Size1792x1024).
			SetQuality(imagestypes.QualityHD).
			SetN(2).
			SetResponseFormat(imagestypes.ResponseFormatURL)

		resp, err := client.Images.Create(context.Background(), req)
		require.NoError(t, err)
		assert.Len(t, resp.Data, 2)
	})

	t.Run("base64 response format", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req imagestypes.ImageGenerationRequest
			json.NewDecoder(r.Body).Decode(&req)

			assert.Equal(t, imagestypes.ResponseFormatB64JSON, req.ResponseFormat)

			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{
						B64JSON: "base64encodedimagedata==",
					},
				},
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

		req := imagestypes.NewImageGenerationRequest("cogview-3", "test")
		req.SetResponseFormat(imagestypes.ResponseFormatB64JSON)

		resp, err := client.Images.Create(context.Background(), req)
		require.NoError(t, err)

		img := resp.GetFirstImage()
		require.NotNil(t, img)
		assert.Equal(t, "base64encodedimagedata==", img.B64JSON)
		assert.Empty(t, img.URL)
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

		req := imagestypes.NewImageGenerationRequest("cogview-3", "")

		resp, err := client.Images.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Invalid prompt")
	})
}

func TestImagesService_Generate(t *testing.T) {
	t.Parallel()

	t.Run("successful generation with URL", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{URL: "https://example.com/generated.png"},
				},
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

		imageURL, err := client.Images.Generate(
			context.Background(),
			"cogview-3",
			"A cat playing piano",
		)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/generated.png", imageURL)
	})

	t.Run("successful generation with base64", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{B64JSON: "base64imagedata=="},
				},
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

		imageData, err := client.Images.Generate(
			context.Background(),
			"cogview-3",
			"A dog in a park",
		)
		require.NoError(t, err)
		assert.Equal(t, "base64imagedata==", imageData)
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data:    []imagestypes.ImageData{},
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

		imageURL, err := client.Images.Generate(
			context.Background(),
			"cogview-3",
			"test",
		)
		require.NoError(t, err)
		assert.Empty(t, imageURL)
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

		imageURL, err := client.Images.Generate(
			context.Background(),
			"cogview-3",
			"test",
		)
		assert.Error(t, err)
		assert.Empty(t, imageURL)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}

func TestImagesService_GenerateMultiple(t *testing.T) {
	t.Parallel()

	t.Run("successful multiple generation", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req imagestypes.ImageGenerationRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Verify count is set
			assert.NotNil(t, req.N)
			assert.Equal(t, 3, *req.N)

			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{URL: "https://example.com/1.png"},
					{URL: "https://example.com/2.png"},
					{URL: "https://example.com/3.png"},
				},
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

		imageURLs, err := client.Images.GenerateMultiple(
			context.Background(),
			"cogview-3",
			"A landscape painting",
			3,
		)
		require.NoError(t, err)
		require.Len(t, imageURLs, 3)
		assert.Equal(t, "https://example.com/1.png", imageURLs[0])
		assert.Equal(t, "https://example.com/2.png", imageURLs[1])
		assert.Equal(t, "https://example.com/3.png", imageURLs[2])
	})

	t.Run("multiple base64 images", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data: []imagestypes.ImageData{
					{B64JSON: "data1=="},
					{B64JSON: "data2=="},
				},
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

		images, err := client.Images.GenerateMultiple(
			context.Background(),
			"cogview-3",
			"test",
			2,
		)
		require.NoError(t, err)
		require.Len(t, images, 2)
		assert.Equal(t, "data1==", images[0])
		assert.Equal(t, "data2==", images[1])
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data:    []imagestypes.ImageData{},
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

		images, err := client.Images.GenerateMultiple(
			context.Background(),
			"cogview-3",
			"test",
			3,
		)
		require.NoError(t, err)
		assert.Len(t, images, 0)
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Rate limit exceeded",
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

		images, err := client.Images.GenerateMultiple(
			context.Background(),
			"cogview-3",
			"test",
			3,
		)
		assert.Error(t, err)
		assert.Nil(t, images)
		assert.Contains(t, err.Error(), "Rate limit exceeded")
	})
}

func TestClient_ImagesService_Integration(t *testing.T) {
	t.Parallel()

	t.Run("client has images service", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.Images)
	})

	t.Run("complete workflow", func(t *testing.T) {
		t.Parallel()

		// Simulate a complete image generation workflow
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req imagestypes.ImageGenerationRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Create response based on request
			data := make([]imagestypes.ImageData, 0)
			count := 1
			if req.N != nil {
				count = *req.N
			}

			for i := 0; i < count; i++ {
				if req.ResponseFormat == imagestypes.ResponseFormatB64JSON {
					data = append(data, imagestypes.ImageData{
						B64JSON:       "base64data==",
						RevisedPrompt: req.Prompt + " (enhanced)",
					})
				} else {
					data = append(data, imagestypes.ImageData{
						URL:           "https://example.com/image.png",
						RevisedPrompt: req.Prompt + " (enhanced)",
					})
				}
			}

			resp := imagestypes.ImageGenerationResponse{
				Created: 1677652288,
				Data:    data,
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

		// Test single image generation
		url, err := client.Images.Generate(context.Background(), "cogview-3", "A sunset")
		require.NoError(t, err)
		assert.NotEmpty(t, url)

		// Test multiple image generation
		urls, err := client.Images.GenerateMultiple(context.Background(), "cogview-3", "A landscape", 2)
		require.NoError(t, err)
		assert.Len(t, urls, 2)
	})
}
