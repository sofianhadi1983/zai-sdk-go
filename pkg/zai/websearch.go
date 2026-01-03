package zai

import (
	"context"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/websearch"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
)

// WebSearchService provides access to the Web Search API.
type WebSearchService struct {
	client *client.BaseClient
}

// newWebSearchService creates a new web search service.
func newWebSearchService(baseClient *client.BaseClient) *WebSearchService {
	return &WebSearchService{
		client: baseClient,
	}
}

// Search performs a web search and returns results.
//
// Example:
//
//	req := websearch.NewWebSearchRequest("latest AI breakthroughs").
//	    SetCount(10).
//	    SetRecencyFilter(websearch.RecencyFilterWeek).
//	    SetSearchIntent(true).
//	    SetIncludeImage(true)
//
//	resp, err := client.WebSearch.Search(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Process search intent if available
//	if resp.HasIntent() {
//	    fmt.Printf("Search Intent: %s\n", resp.SearchIntent.Intent)
//	    fmt.Printf("Optimized Query: %s\n", resp.SearchIntent.Query)
//	}
//
//	// Process search results
//	for i, result := range resp.GetResults() {
//	    fmt.Printf("%d. %s\n", i+1, result.Title)
//	    fmt.Printf("   URL: %s\n", result.Link)
//	    fmt.Printf("   Content: %s\n", result.Content)
//	    fmt.Printf("   Published: %s\n", result.PublishDate)
//	}
func (s *WebSearchService) Search(ctx context.Context, req *websearch.WebSearchRequest) (*websearch.WebSearchResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/web_search", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp websearch.WebSearchResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
