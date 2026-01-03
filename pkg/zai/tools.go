package zai

import (
	"context"

	"github.com/z-ai/zai-sdk-go/api/types/tools"
	"github.com/z-ai/zai-sdk-go/internal/client"
	streaming "github.com/z-ai/zai-sdk-go/internal/streaming"
)

// ToolsService provides access to the Tools API.
type ToolsService struct {
	client *client.BaseClient
}

// newToolsService creates a new tools service.
func newToolsService(baseClient *client.BaseClient) *ToolsService {
	return &ToolsService{
		client: baseClient,
	}
}

// WebSearch performs web search using AI models.
//
// Example (non-streaming):
//
//	messages := []chat.Message{
//	    {
//	        Role:    chat.RoleUser,
//	        Content: "Latest AI breakthroughs in 2024",
//	    },
//	}
//
//	req := tools.NewWebSearchRequest("web-search-pro", messages).
//	    SetRecentDays(7).
//	    SetScope("academic")
//
//	resp, err := client.Tools.WebSearch(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Process search results
//	for _, result := range resp.GetSearchResults() {
//	    fmt.Printf("Title: %s\n", result.Title)
//	    fmt.Printf("Link: %s\n", result.Link)
//	    fmt.Printf("Content: %s\n", result.Content)
//	}
//
//	// Process search intents
//	for _, intent := range resp.GetSearchIntents() {
//	    fmt.Printf("Optimized Query: %s\n", intent.Query)
//	    fmt.Printf("Intent Type: %s\n", intent.Intent)
//	}
func (s *ToolsService) WebSearch(ctx context.Context, req *tools.WebSearchRequest) (*tools.WebSearchResponse, error) {
	// Ensure streaming is disabled for non-streaming request
	req.Stream = false

	// Make the API request
	apiResp, err := s.client.Post(ctx, "/tools", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp tools.WebSearchResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// WebSearchStream performs streaming web search using AI models.
//
// Example (streaming):
//
//	messages := []chat.Message{
//	    {
//	        Role:    chat.RoleUser,
//	        Content: "What's happening in AI today?",
//	    },
//	}
//
//	req := tools.NewWebSearchRequest("web-search-pro", messages).
//	    SetStream(true).
//	    SetRecentDays(1)
//
//	stream, err := client.Tools.WebSearchStream(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//	defer stream.Close()
//
//	// Process streaming chunks
//	for stream.Next() {
//	    chunk := stream.Current()
//	    // Process chunk
//	    for _, choice := range chunk.Choices {
//	        for _, toolCall := range choice.Delta.ToolCalls {
//	            if toolCall.SearchResult != nil {
//	                fmt.Printf("Result: %s\n", toolCall.SearchResult.Title)
//	            }
//	        }
//	    }
//	}
//
//	if err := stream.Err(); err != nil {
//	    // Handle error
//	}
func (s *ToolsService) WebSearchStream(ctx context.Context, req *tools.WebSearchRequest) (*streaming.Stream[tools.WebSearchChunk], error) {
	// Ensure streaming is enabled
	req.Stream = true

	// Make the streaming request
	streamResp, err := s.client.Stream(ctx, "/tools", req)
	if err != nil {
		return nil, err
	}

	// Create typed stream
	return client.NewTypedStream[tools.WebSearchChunk](streamResp, ctx), nil
}
