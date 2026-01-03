// Package embeddings provides types for the Embeddings API.
package embeddings

import "github.com/sofianhadi1983/zai-sdk-go/internal/models"

// EmbeddingRequest represents a request to create embeddings.
type EmbeddingRequest struct {
	// Model is the ID of the model to use.
	// Required. Example: "embedding-2"
	Model string `json:"model"`

	// Input is the text or array of texts to embed.
	// Can be a string or []string.
	// Required.
	Input interface{} `json:"input"`

	// Dimensions is the number of dimensions for the output embeddings.
	// Only supported in some models.
	// Optional.
	Dimensions *int `json:"dimensions,omitempty"`

	// EncodingFormat is the format to return the embeddings in.
	// Can be "float" or "base64".
	// Default: "float"
	EncodingFormat string `json:"encoding_format,omitempty"`

	// User is a unique identifier for the end-user.
	// Optional.
	User string `json:"user,omitempty"`

	// RequestID is a unique identifier for the request.
	// Optional.
	RequestID string `json:"request_id,omitempty"`
}

// SetDimensions sets the dimensions parameter.
func (r *EmbeddingRequest) SetDimensions(dimensions int) *EmbeddingRequest {
	r.Dimensions = &dimensions
	return r
}

// SetEncodingFormat sets the encoding format.
func (r *EmbeddingRequest) SetEncodingFormat(format string) *EmbeddingRequest {
	r.EncodingFormat = format
	return r
}

// SetUser sets the user identifier.
func (r *EmbeddingRequest) SetUser(user string) *EmbeddingRequest {
	r.User = user
	return r
}

// EmbeddingResponse represents the response from an embedding request.
type EmbeddingResponse struct {
	// Object is the object type (always "list").
	Object string `json:"object"`

	// Data is the list of embeddings.
	Data []Embedding `json:"data"`

	// Model is the model used for the embeddings.
	Model string `json:"model"`

	// Usage is the token usage information.
	Usage *models.Usage `json:"usage,omitempty"`
}

// Embedding represents a single embedding.
type Embedding struct {
	// Object is the object type (always "embedding").
	Object string `json:"object"`

	// Embedding is the embedding vector.
	// Can be []float64 or string (base64 encoded) depending on encoding_format.
	Embedding interface{} `json:"embedding"`

	// Index is the index of the embedding in the list.
	Index int `json:"index"`
}

// GetFloatEmbedding returns the embedding as a float64 slice.
// Returns nil if the embedding is not in float format.
func (e *Embedding) GetFloatEmbedding() []float64 {
	if floats, ok := e.Embedding.([]interface{}); ok {
		result := make([]float64, len(floats))
		for i, v := range floats {
			if f, ok := v.(float64); ok {
				result[i] = f
			}
		}
		return result
	}
	if floats, ok := e.Embedding.([]float64); ok {
		return floats
	}
	return nil
}

// GetBase64Embedding returns the embedding as a base64 encoded string.
// Returns empty string if the embedding is not in base64 format.
func (e *Embedding) GetBase64Embedding() string {
	if str, ok := e.Embedding.(string); ok {
		return str
	}
	return ""
}

// GetFirstEmbedding returns the first embedding from the response.
// Returns nil if there are no embeddings.
func (r *EmbeddingResponse) GetFirstEmbedding() *Embedding {
	if len(r.Data) == 0 {
		return nil
	}
	return &r.Data[0]
}

// GetFloatEmbeddings returns all embeddings as float64 slices.
// Skips any embeddings that are not in float format.
func (r *EmbeddingResponse) GetFloatEmbeddings() [][]float64 {
	result := make([][]float64, 0, len(r.Data))
	for _, emb := range r.Data {
		if floats := emb.GetFloatEmbedding(); floats != nil {
			result = append(result, floats)
		}
	}
	return result
}

const (
	// EncodingFormatFloat returns embeddings as float arrays.
	EncodingFormatFloat = "float"

	// EncodingFormatBase64 returns embeddings as base64 encoded strings.
	EncodingFormatBase64 = "base64"
)

// NewEmbeddingRequest creates a new embedding request with a single text input.
func NewEmbeddingRequest(model, text string) *EmbeddingRequest {
	return &EmbeddingRequest{
		Model: model,
		Input: text,
	}
}

// NewBatchEmbeddingRequest creates a new embedding request with multiple text inputs.
func NewBatchEmbeddingRequest(model string, texts []string) *EmbeddingRequest {
	return &EmbeddingRequest{
		Model: model,
		Input: texts,
	}
}
