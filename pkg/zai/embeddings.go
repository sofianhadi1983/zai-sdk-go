package zai

import (
	"context"

	"github.com/z-ai/zai-sdk-go/api/types/embeddings"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// EmbeddingsService provides access to the Embeddings API.
type EmbeddingsService struct {
	client *client.BaseClient
}

// newEmbeddingsService creates a new embeddings service.
func newEmbeddingsService(baseClient *client.BaseClient) *EmbeddingsService {
	return &EmbeddingsService{
		client: baseClient,
	}
}

// Create creates embeddings for the given input text(s).
//
// Example with single text:
//
//	req := embeddings.NewEmbeddingRequest("embedding-2", "Hello world")
//	resp, err := client.Embeddings.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	emb := resp.GetFirstEmbedding()
//	floats := emb.GetFloatEmbedding()
//	fmt.Printf("Embedding dimension: %d\n", len(floats))
//
// Example with batch:
//
//	texts := []string{"Hello", "World"}
//	req := embeddings.NewBatchEmbeddingRequest("embedding-2", texts)
//	resp, err := client.Embeddings.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, emb := range resp.Data {
//	    floats := emb.GetFloatEmbedding()
//	    fmt.Printf("Embedding %d: %d dimensions\n", emb.Index, len(floats))
//	}
func (s *EmbeddingsService) Create(ctx context.Context, req *embeddings.EmbeddingRequest) (*embeddings.EmbeddingResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/embeddings", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp embeddings.EmbeddingResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateSingle is a convenience method for creating embeddings for a single text.
// Returns the embedding vector directly.
//
// Example:
//
//	embedding, err := client.Embeddings.CreateSingle(ctx, "embedding-2", "Hello world")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Got embedding with %d dimensions\n", len(embedding))
func (s *EmbeddingsService) CreateSingle(ctx context.Context, model, text string) ([]float64, error) {
	req := embeddings.NewEmbeddingRequest(model, text)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	emb := resp.GetFirstEmbedding()
	if emb == nil {
		return nil, nil
	}

	return emb.GetFloatEmbedding(), nil
}

// CreateBatch is a convenience method for creating embeddings for multiple texts.
// Returns a slice of embedding vectors.
//
// Example:
//
//	texts := []string{"Hello", "World", "AI"}
//	embeddings, err := client.Embeddings.CreateBatch(ctx, "embedding-2", texts)
//	if err != nil {
//	    // Handle error
//	}
//
//	for i, embedding := range embeddings {
//	    fmt.Printf("Text %d: %d dimensions\n", i, len(embedding))
//	}
func (s *EmbeddingsService) CreateBatch(ctx context.Context, model string, texts []string) ([][]float64, error) {
	req := embeddings.NewBatchEmbeddingRequest(model, texts)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetFloatEmbeddings(), nil
}
