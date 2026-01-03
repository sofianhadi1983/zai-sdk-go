package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	embeddingstypes "github.com/sofianhadi1983/zai-sdk-go/api/types/embeddings"
)

func TestEmbeddingsService_Create(t *testing.T) {
	t.Parallel()

	t.Run("single text embedding", func(t *testing.T) {
		t.Parallel()

		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/embeddings", r.URL.Path)
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			var req embeddingstypes.EmbeddingRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "embedding-2", req.Model)
			assert.Equal(t, "Hello world", req.Input)

			// Send response
			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data: []embeddingstypes.Embedding{
					{
						Object:    "embedding",
						Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
						Index:     0,
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
		req := embeddingstypes.NewEmbeddingRequest("embedding-2", "Hello world")

		resp, err := client.Embeddings.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, "list", resp.Object)
		assert.Equal(t, "embedding-2", resp.Model)
		assert.Len(t, resp.Data, 1)

		emb := resp.GetFirstEmbedding()
		require.NotNil(t, emb)
		floats := emb.GetFloatEmbedding()
		assert.Equal(t, []float64{0.1, 0.2, 0.3, 0.4, 0.5}, floats)
	})

	t.Run("batch text embeddings", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req embeddingstypes.EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Verify batch input
			inputs, ok := req.Input.([]interface{})
			assert.True(t, ok)
			assert.Len(t, inputs, 2)

			// Send response with multiple embeddings
			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data: []embeddingstypes.Embedding{
					{
						Object:    "embedding",
						Embedding: []float64{0.1, 0.2},
						Index:     0,
					},
					{
						Object:    "embedding",
						Embedding: []float64{0.3, 0.4},
						Index:     1,
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

		texts := []string{"Hello", "World"}
		req := embeddingstypes.NewBatchEmbeddingRequest("embedding-2", texts)

		resp, err := client.Embeddings.Create(context.Background(), req)
		require.NoError(t, err)
		assert.Len(t, resp.Data, 2)

		floats := resp.GetFloatEmbeddings()
		assert.Len(t, floats, 2)
		assert.Equal(t, []float64{0.1, 0.2}, floats[0])
		assert.Equal(t, []float64{0.3, 0.4}, floats[1])
	})

	t.Run("with dimensions parameter", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req embeddingstypes.EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)

			assert.NotNil(t, req.Dimensions)
			assert.Equal(t, 512, *req.Dimensions)

			// Return embedding with specified dimensions
			embedding := make([]float64, 512)
			for i := range embedding {
				embedding[i] = float64(i) * 0.001
			}

			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data: []embeddingstypes.Embedding{
					{
						Object:    "embedding",
						Embedding: embedding,
						Index:     0,
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

		req := embeddingstypes.NewEmbeddingRequest("embedding-2", "test")
		req.SetDimensions(512)

		resp, err := client.Embeddings.Create(context.Background(), req)
		require.NoError(t, err)

		emb := resp.GetFirstEmbedding()
		floats := emb.GetFloatEmbedding()
		assert.Len(t, floats, 512)
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid model",
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

		req := embeddingstypes.NewEmbeddingRequest("invalid-model", "test")

		resp, err := client.Embeddings.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Invalid model")
	})
}

func TestEmbeddingsService_CreateSingle(t *testing.T) {
	t.Parallel()

	t.Run("successful embedding", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data: []embeddingstypes.Embedding{
					{
						Object:    "embedding",
						Embedding: []float64{1.0, 2.0, 3.0},
						Index:     0,
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

		embedding, err := client.Embeddings.CreateSingle(
			context.Background(),
			"embedding-2",
			"Hello world",
		)
		require.NoError(t, err)
		assert.Equal(t, []float64{1.0, 2.0, 3.0}, embedding)
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data:   []embeddingstypes.Embedding{},
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

		embedding, err := client.Embeddings.CreateSingle(
			context.Background(),
			"embedding-2",
			"test",
		)
		require.NoError(t, err)
		assert.Nil(t, embedding)
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

		embedding, err := client.Embeddings.CreateSingle(
			context.Background(),
			"embedding-2",
			"test",
		)
		assert.Error(t, err)
		assert.Nil(t, embedding)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}

func TestEmbeddingsService_CreateBatch(t *testing.T) {
	t.Parallel()

	t.Run("successful batch", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req embeddingstypes.EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Verify it's a batch request
			inputs, ok := req.Input.([]interface{})
			assert.True(t, ok)
			assert.Len(t, inputs, 3)

			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data: []embeddingstypes.Embedding{
					{
						Object:    "embedding",
						Embedding: []float64{1.0, 2.0},
						Index:     0,
					},
					{
						Object:    "embedding",
						Embedding: []float64{3.0, 4.0},
						Index:     1,
					},
					{
						Object:    "embedding",
						Embedding: []float64{5.0, 6.0},
						Index:     2,
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

		texts := []string{"Hello", "World", "Test"}
		embeddings, err := client.Embeddings.CreateBatch(
			context.Background(),
			"embedding-2",
			texts,
		)
		require.NoError(t, err)
		assert.Len(t, embeddings, 3)
		assert.Equal(t, []float64{1.0, 2.0}, embeddings[0])
		assert.Equal(t, []float64{3.0, 4.0}, embeddings[1])
		assert.Equal(t, []float64{5.0, 6.0}, embeddings[2])
	})

	t.Run("empty batch", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data:   []embeddingstypes.Embedding{},
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

		embeddings, err := client.Embeddings.CreateBatch(
			context.Background(),
			"embedding-2",
			[]string{},
		)
		require.NoError(t, err)
		assert.Len(t, embeddings, 0)
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

		embeddings, err := client.Embeddings.CreateBatch(
			context.Background(),
			"embedding-2",
			[]string{"test1", "test2"},
		)
		assert.Error(t, err)
		assert.Nil(t, embeddings)
		assert.Contains(t, err.Error(), "Rate limit exceeded")
	})
}

func TestClient_EmbeddingsService_Integration(t *testing.T) {
	t.Parallel()

	t.Run("client has embeddings service", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.Embeddings)
	})

	t.Run("semantic similarity example", func(t *testing.T) {
		t.Parallel()

		// Simulate a semantic similarity use case
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req embeddingstypes.EmbeddingRequest
			json.NewDecoder(r.Body).Decode(&req)

			inputs, _ := req.Input.([]interface{})
			embeddings := make([]embeddingstypes.Embedding, len(inputs))

			// Generate mock embeddings (in real usage, these would be semantically meaningful)
			for i := range inputs {
				// Similar texts would have similar embeddings
				var emb []float64
				if i == 0 || i == 1 { // "cat" and "kitten" are similar
					emb = []float64{0.9, 0.1, 0.05}
				} else { // "dog" is different
					emb = []float64{0.1, 0.9, 0.05}
				}

				embeddings[i] = embeddingstypes.Embedding{
					Object:    "embedding",
					Embedding: emb,
					Index:     i,
				}
			}

			resp := embeddingstypes.EmbeddingResponse{
				Object: "list",
				Model:  "embedding-2",
				Data:   embeddings,
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

		// Get embeddings for three texts
		texts := []string{"cat", "kitten", "dog"}
		embeds, err := client.Embeddings.CreateBatch(
			context.Background(),
			"embedding-2",
			texts,
		)
		require.NoError(t, err)
		assert.Len(t, embeds, 3)

		// Verify embeddings were created
		for i, emb := range embeds {
			assert.NotNil(t, emb)
			assert.Greater(t, len(emb), 0, "Embedding %d should have dimensions", i)
		}

		// In a real scenario, you would compute cosine similarity here
		// to find that "cat" and "kitten" are more similar than "cat" and "dog"
	})
}
