package main

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/embeddings"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
	// Create a new client
	client, err := zai.NewClient(
		zai.WithAPIKey("your-api-key.your-secret"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Single text embedding
	fmt.Println("=== Example 1: Single Text Embedding ===")
	singleEmbeddingExample(ctx, client)

	fmt.Println("\n=== Example 2: Batch Embeddings ===")
	batchEmbeddingExample(ctx, client)

	fmt.Println("\n=== Example 3: Semantic Similarity ===")
	semanticSimilarityExample(ctx, client)

	fmt.Println("\n=== Example 4: With Custom Dimensions ===")
	customDimensionsExample(ctx, client)
}

func singleEmbeddingExample(ctx context.Context, client *zai.Client) {
	// Create embedding for a single text using convenience method
	embedding, err := client.Embeddings.CreateSingle(
		ctx,
		"embedding-2",
		"Hello, world!",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Generated embedding with %d dimensions\n", len(embedding))
	fmt.Printf("First 5 values: %v\n", embedding[:5])
}

func batchEmbeddingExample(ctx context.Context, client *zai.Client) {
	// Create embeddings for multiple texts at once
	texts := []string{
		"The quick brown fox",
		"jumps over the lazy dog",
		"Machine learning is fascinating",
	}

	embeddings, err := client.Embeddings.CreateBatch(
		ctx,
		"embedding-2",
		texts,
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Generated %d embeddings\n", len(embeddings))
	for i, emb := range embeddings {
		fmt.Printf("Text %d (\"%s\"): %d dimensions\n",
			i+1, texts[i], len(emb))
	}
}

func semanticSimilarityExample(ctx context.Context, client *zai.Client) {
	// Compare semantic similarity between texts
	texts := []string{
		"The cat sits on the mat",
		"A feline rests on the rug",
		"Dogs are playing in the park",
	}

	embeddings, err := client.Embeddings.CreateBatch(
		ctx,
		"embedding-2",
		texts,
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Compute cosine similarities
	fmt.Println("Cosine similarities:")
	for i := 0; i < len(texts); i++ {
		for j := i + 1; j < len(texts); j++ {
			similarity := cosineSimilarity(embeddings[i], embeddings[j])
			fmt.Printf("  \"%s\" vs \"%s\": %.4f\n",
				truncate(texts[i], 30),
				truncate(texts[j], 30),
				similarity)
		}
	}
}

func customDimensionsExample(ctx context.Context, client *zai.Client) {
	// Create embedding with custom dimensions
	req := embeddings.NewEmbeddingRequest("embedding-2", "Custom dimensions example")
	req.SetDimensions(512)

	resp, err := client.Embeddings.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	emb := resp.GetFirstEmbedding()
	if emb != nil {
		floats := emb.GetFloatEmbedding()
		fmt.Printf("Generated embedding with custom dimensions: %d\n", len(floats))
		fmt.Printf("Model used: %s\n", resp.Model)

		if resp.Usage != nil {
			fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
		}
	}
}

// cosineSimilarity computes the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// truncate truncates a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
