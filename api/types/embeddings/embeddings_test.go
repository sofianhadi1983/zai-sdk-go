package embeddings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmbeddingRequest(t *testing.T) {
	t.Parallel()

	req := NewEmbeddingRequest("embedding-2", "Hello world")

	assert.Equal(t, "embedding-2", req.Model)
	assert.Equal(t, "Hello world", req.Input)
}

func TestNewBatchEmbeddingRequest(t *testing.T) {
	t.Parallel()

	texts := []string{"Hello", "World"}
	req := NewBatchEmbeddingRequest("embedding-2", texts)

	assert.Equal(t, "embedding-2", req.Model)
	assert.Equal(t, texts, req.Input)
}

func TestEmbeddingRequest_Setters(t *testing.T) {
	t.Parallel()

	t.Run("SetDimensions", func(t *testing.T) {
		t.Parallel()

		req := &EmbeddingRequest{}
		req.SetDimensions(1024)

		require.NotNil(t, req.Dimensions)
		assert.Equal(t, 1024, *req.Dimensions)
	})

	t.Run("SetEncodingFormat", func(t *testing.T) {
		t.Parallel()

		req := &EmbeddingRequest{}
		req.SetEncodingFormat(EncodingFormatBase64)

		assert.Equal(t, EncodingFormatBase64, req.EncodingFormat)
	})

	t.Run("SetUser", func(t *testing.T) {
		t.Parallel()

		req := &EmbeddingRequest{}
		req.SetUser("user-123")

		assert.Equal(t, "user-123", req.User)
	})

	t.Run("chained setters", func(t *testing.T) {
		t.Parallel()

		req := NewEmbeddingRequest("embedding-2", "test")
		req.SetDimensions(512).
			SetEncodingFormat(EncodingFormatFloat).
			SetUser("user-456")

		require.NotNil(t, req.Dimensions)
		assert.Equal(t, 512, *req.Dimensions)
		assert.Equal(t, EncodingFormatFloat, req.EncodingFormat)
		assert.Equal(t, "user-456", req.User)
	})
}

func TestEmbeddingRequest_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal single text", func(t *testing.T) {
		t.Parallel()

		req := NewEmbeddingRequest("embedding-2", "Hello")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded EmbeddingRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "embedding-2", decoded.Model)
		assert.Equal(t, "Hello", decoded.Input)
	})

	t.Run("marshal batch texts", func(t *testing.T) {
		t.Parallel()

		texts := []string{"Hello", "World"}
		req := NewBatchEmbeddingRequest("embedding-2", texts)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "Hello")
		assert.Contains(t, string(data), "World")
	})

	t.Run("marshal with dimensions", func(t *testing.T) {
		t.Parallel()

		req := NewEmbeddingRequest("embedding-2", "test")
		req.SetDimensions(256)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "dimensions")
		assert.Contains(t, string(data), "256")
	})

	t.Run("omit empty fields", func(t *testing.T) {
		t.Parallel()

		req := NewEmbeddingRequest("embedding-2", "test")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		// Dimensions should be omitted when nil
		assert.NotContains(t, string(data), "dimensions")
		assert.NotContains(t, string(data), "user")
	})
}

func TestEmbedding_GetFloatEmbedding(t *testing.T) {
	t.Parallel()

	t.Run("with float slice interface", func(t *testing.T) {
		t.Parallel()

		emb := &Embedding{
			Embedding: []interface{}{1.0, 2.0, 3.0},
		}

		floats := emb.GetFloatEmbedding()
		require.NotNil(t, floats)
		assert.Equal(t, []float64{1.0, 2.0, 3.0}, floats)
	})

	t.Run("with float64 slice", func(t *testing.T) {
		t.Parallel()

		emb := &Embedding{
			Embedding: []float64{4.0, 5.0, 6.0},
		}

		floats := emb.GetFloatEmbedding()
		require.NotNil(t, floats)
		assert.Equal(t, []float64{4.0, 5.0, 6.0}, floats)
	})

	t.Run("with base64 string", func(t *testing.T) {
		t.Parallel()

		emb := &Embedding{
			Embedding: "YmFzZTY0ZW5jb2RlZA==",
		}

		floats := emb.GetFloatEmbedding()
		assert.Nil(t, floats)
	})
}

func TestEmbedding_GetBase64Embedding(t *testing.T) {
	t.Parallel()

	t.Run("with base64 string", func(t *testing.T) {
		t.Parallel()

		emb := &Embedding{
			Embedding: "YmFzZTY0ZW5jb2RlZA==",
		}

		base64 := emb.GetBase64Embedding()
		assert.Equal(t, "YmFzZTY0ZW5jb2RlZA==", base64)
	})

	t.Run("with float slice", func(t *testing.T) {
		t.Parallel()

		emb := &Embedding{
			Embedding: []float64{1.0, 2.0, 3.0},
		}

		base64 := emb.GetBase64Embedding()
		assert.Equal(t, "", base64)
	})
}

func TestEmbeddingResponse_GetFirstEmbedding(t *testing.T) {
	t.Parallel()

	t.Run("with embeddings", func(t *testing.T) {
		t.Parallel()

		resp := &EmbeddingResponse{
			Data: []Embedding{
				{Index: 0, Embedding: []float64{1.0, 2.0}},
				{Index: 1, Embedding: []float64{3.0, 4.0}},
			},
		}

		emb := resp.GetFirstEmbedding()
		require.NotNil(t, emb)
		assert.Equal(t, 0, emb.Index)
	})

	t.Run("without embeddings", func(t *testing.T) {
		t.Parallel()

		resp := &EmbeddingResponse{
			Data: []Embedding{},
		}

		emb := resp.GetFirstEmbedding()
		assert.Nil(t, emb)
	})
}

func TestEmbeddingResponse_GetFloatEmbeddings(t *testing.T) {
	t.Parallel()

	t.Run("all float embeddings", func(t *testing.T) {
		t.Parallel()

		resp := &EmbeddingResponse{
			Data: []Embedding{
				{Index: 0, Embedding: []float64{1.0, 2.0}},
				{Index: 1, Embedding: []float64{3.0, 4.0}},
			},
		}

		floats := resp.GetFloatEmbeddings()
		require.Len(t, floats, 2)
		assert.Equal(t, []float64{1.0, 2.0}, floats[0])
		assert.Equal(t, []float64{3.0, 4.0}, floats[1])
	})

	t.Run("mixed format embeddings", func(t *testing.T) {
		t.Parallel()

		resp := &EmbeddingResponse{
			Data: []Embedding{
				{Index: 0, Embedding: []float64{1.0, 2.0}},
				{Index: 1, Embedding: "base64string"},
				{Index: 2, Embedding: []float64{3.0, 4.0}},
			},
		}

		floats := resp.GetFloatEmbeddings()
		require.Len(t, floats, 2)
		assert.Equal(t, []float64{1.0, 2.0}, floats[0])
		assert.Equal(t, []float64{3.0, 4.0}, floats[1])
	})

	t.Run("no embeddings", func(t *testing.T) {
		t.Parallel()

		resp := &EmbeddingResponse{
			Data: []Embedding{},
		}

		floats := resp.GetFloatEmbeddings()
		assert.Len(t, floats, 0)
	})
}

func TestEmbeddingResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal float embeddings", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"object": "list",
			"data": [
				{
					"object": "embedding",
					"embedding": [0.1, 0.2, 0.3],
					"index": 0
				},
				{
					"object": "embedding",
					"embedding": [0.4, 0.5, 0.6],
					"index": 1
				}
			],
			"model": "embedding-2",
			"usage": {
				"prompt_tokens": 10,
				"total_tokens": 10
			}
		}`

		var resp EmbeddingResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "list", resp.Object)
		assert.Equal(t, "embedding-2", resp.Model)
		assert.Len(t, resp.Data, 2)

		require.NotNil(t, resp.Usage)
		assert.Equal(t, 10, resp.Usage.PromptTokens)
		assert.Equal(t, 10, resp.Usage.TotalTokens)

		// Check embeddings
		emb1 := resp.Data[0]
		assert.Equal(t, "embedding", emb1.Object)
		assert.Equal(t, 0, emb1.Index)

		floats := emb1.GetFloatEmbedding()
		require.NotNil(t, floats)
		assert.Equal(t, 0.1, floats[0])
		assert.Equal(t, 0.2, floats[1])
		assert.Equal(t, 0.3, floats[2])
	})

	t.Run("unmarshal base64 embeddings", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"object": "list",
			"data": [
				{
					"object": "embedding",
					"embedding": "YmFzZTY0ZW5jb2RlZA==",
					"index": 0
				}
			],
			"model": "embedding-2"
		}`

		var resp EmbeddingResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Data, 1)

		emb := resp.Data[0]
		base64 := emb.GetBase64Embedding()
		assert.Equal(t, "YmFzZTY0ZW5jb2RlZA==", base64)
	})
}

func TestEncodingFormat_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "float", EncodingFormatFloat)
	assert.Equal(t, "base64", EncodingFormatBase64)
}

func TestEmbeddingRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Build a complete request
	req := NewBatchEmbeddingRequest("embedding-2", []string{
		"The quick brown fox",
		"jumps over the lazy dog",
	})

	req.SetDimensions(1024).
		SetEncodingFormat(EncodingFormatFloat).
		SetUser("user-123")

	// Verify the request is complete
	assert.Equal(t, "embedding-2", req.Model)
	assert.Equal(t, []string{"The quick brown fox", "jumps over the lazy dog"}, req.Input)
	require.NotNil(t, req.Dimensions)
	assert.Equal(t, 1024, *req.Dimensions)
	assert.Equal(t, EncodingFormatFloat, req.EncodingFormat)
	assert.Equal(t, "user-123", req.User)

	// Ensure it can be marshaled
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), "embedding-2")
	assert.Contains(t, string(data), "quick brown fox")
	assert.Contains(t, string(data), "1024")
}

func TestEmbeddingResponse_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response with high-dimensional embeddings
	jsonData := `{
		"object": "list",
		"data": [
			{
				"object": "embedding",
				"embedding": [0.0023, -0.009, 0.015, -0.011, 0.008],
				"index": 0
			}
		],
		"model": "embedding-2",
		"usage": {
			"prompt_tokens": 8,
			"total_tokens": 8
		}
	}`

	var resp EmbeddingResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "embedding-2", resp.Model)
	assert.Equal(t, "list", resp.Object)

	// Verify helper methods work
	emb := resp.GetFirstEmbedding()
	require.NotNil(t, emb)
	assert.Equal(t, "embedding", emb.Object)

	floats := emb.GetFloatEmbedding()
	require.Len(t, floats, 5)
	assert.Equal(t, 0.0023, floats[0])

	// Verify usage
	require.NotNil(t, resp.Usage)
	assert.Equal(t, 8, resp.Usage.TotalTokens)
}
