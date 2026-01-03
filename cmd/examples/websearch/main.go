package main

import (
	"context"
	"fmt"
	"log"

	"github.com/z-ai/zai-sdk-go/api/types/websearch"
	"github.com/z-ai/zai-sdk-go/pkg/zai"
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

	// Example 1: Basic web search
	fmt.Println("=== Example 1: Basic Web Search ===")
	basicSearchExample(ctx, client)

	// Example 2: Search with intent analysis
	fmt.Println("\n=== Example 2: Search with Intent Analysis ===")
	searchWithIntentExample(ctx, client)

	// Example 3: Search with filters
	fmt.Println("\n=== Example 3: Search with Filters ===")
	searchWithFiltersExample(ctx, client)

	// Example 4: Search with images
	fmt.Println("\n=== Example 4: Search with Images ===")
	searchWithImagesExample(ctx, client)

	// Example 5: Domain-specific search
	fmt.Println("\n=== Example 5: Domain-Specific Search ===")
	domainSpecificSearchExample(ctx, client)

	// Example 6: Recent content search
	fmt.Println("\n=== Example 6: Recent Content Search ===")
	recentContentSearchExample(ctx, client)
}

func basicSearchExample(ctx context.Context, client *zai.Client) {
	// Perform a simple web search
	req := websearch.NewWebSearchRequest("latest artificial intelligence breakthroughs")

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Search ID: %s\n", resp.ID)
	fmt.Printf("Found %d results:\n\n", len(resp.GetResults()))

	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)
		fmt.Printf("   Source: %s\n", result.Media)
		fmt.Printf("   Published: %s\n", result.PublishDate)
		fmt.Printf("   Snippet: %s\n", truncateString(result.Content, 100))
		fmt.Printf("   Reference: %s\n\n", result.Refer)
	}
}

func searchWithIntentExample(ctx context.Context, client *zai.Client) {
	// Search with intent analysis enabled
	req := websearch.NewWebSearchRequest("how to learn machine learning").
		SetSearchIntent(true).
		SetCount(5)

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Display search intent analysis
	if resp.HasIntent() {
		fmt.Println("Search Intent Analysis:")
		fmt.Printf("  Optimized Query: %s\n", resp.SearchIntent.Query)
		fmt.Printf("  Intent Type: %s\n", resp.SearchIntent.Intent)
		fmt.Printf("  Keywords: %s\n\n", resp.SearchIntent.Keywords)
	}

	// Display search results
	fmt.Printf("Search Results (%d found):\n\n", len(resp.GetResults()))
	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   %s\n", result.Link)
		fmt.Printf("   %s\n\n", truncateString(result.Content, 80))
	}
}

func searchWithFiltersExample(ctx context.Context, client *zai.Client) {
	// Search with various filters
	req := websearch.NewWebSearchRequest("quantum computing research").
		SetCount(10).
		SetRecencyFilter(websearch.RecencyFilterMonth).
		SetContentSize(websearch.ContentSizeLarge)

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Filtered Search Results (Recent: %s, Size: %s):\n\n",
		websearch.RecencyFilterMonth, websearch.ContentSizeLarge)

	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   Published: %s\n", result.PublishDate)
		fmt.Printf("   Content Length: %d characters\n", len(result.Content))
		fmt.Printf("   Link: %s\n\n", result.Link)
	}
}

func searchWithImagesExample(ctx context.Context, client *zai.Client) {
	// Search with image inclusion
	req := websearch.NewWebSearchRequest("latest smartphone releases 2024").
		SetIncludeImage(true).
		SetCount(5)

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Search Results with Images:")

	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)

		if len(result.Images) > 0 {
			fmt.Printf("   Images (%d):\n", len(result.Images))
			for j, imgURL := range result.Images {
				fmt.Printf("     %d. %s\n", j+1, imgURL)
			}
		} else {
			fmt.Println("   No images available")
		}
		fmt.Println()
	}
}

func domainSpecificSearchExample(ctx context.Context, client *zai.Client) {
	// Search within a specific domain
	req := websearch.NewWebSearchRequest("machine learning papers").
		SetDomainFilter("arxiv.org").
		SetCount(5).
		SetRecencyFilter(websearch.RecencyFilterWeek)

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Domain-Specific Search (arxiv.org only):")

	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)
		fmt.Printf("   Source: %s\n", result.Media)
		fmt.Printf("   Published: %s\n", result.PublishDate)
		fmt.Printf("   Abstract: %s\n\n", truncateString(result.Content, 150))
	}
}

func recentContentSearchExample(ctx context.Context, client *zai.Client) {
	// Search for recent content with different time filters
	queries := []struct {
		query  string
		filter string
	}{
		{"AI news", websearch.RecencyFilterDay},
		{"tech updates", websearch.RecencyFilterWeek},
		{"industry trends", websearch.RecencyFilterMonth},
	}

	for _, q := range queries {
		req := websearch.NewWebSearchRequest(q.query).
			SetRecencyFilter(q.filter).
			SetCount(3)

		resp, err := client.WebSearch.Search(ctx, req)
		if err != nil {
			log.Printf("Error searching '%s': %v", q.query, err)
			continue
		}

		fmt.Printf("Recent Content - '%s' (Filter: %s):\n", q.query, q.filter)
		for i, result := range resp.GetResults() {
			fmt.Printf("  %d. %s (%s)\n", i+1, result.Title, result.PublishDate)
		}
		fmt.Println()
	}
}

// Advanced examples

func searchWithSensitiveWordCheckExample(ctx context.Context, client *zai.Client) {
	// Search with sensitive word filtering enabled
	sensitiveCheck := &websearch.SensitiveWordCheck{
		Type:   websearch.SensitiveWordTypeAll,
		Status: websearch.SensitiveWordStatusEnable,
	}

	req := websearch.NewWebSearchRequest("controversial topic").
		SetSensitiveWordCheck(sensitiveCheck).
		SetCount(10)

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Search with Sensitive Word Filtering:")
	fmt.Printf("Results returned: %d\n", len(resp.GetResults()))

	for i, result := range resp.GetResults() {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   %s\n\n", result.Link)
	}
}

func comprehensiveSearchExample(ctx context.Context, client *zai.Client) {
	// Comprehensive search with all options
	req := websearch.NewWebSearchRequest("artificial intelligence ethics").
		SetSearchEngine("google").
		SetCount(15).
		SetSearchIntent(true).
		SetIncludeImage(true).
		SetRecencyFilter(websearch.RecencyFilterMonth).
		SetContentSize(websearch.ContentSizeLarge).
		SetRequestID("req_comprehensive_123").
		SetUserID("user_456")

	resp, err := client.WebSearch.Search(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Comprehensive Search Results:")
	fmt.Printf("Request ID: %s\n", resp.RequestID)
	fmt.Printf("Search ID: %s\n\n", resp.ID)

	// Display intent if available
	if resp.HasIntent() {
		fmt.Println("Intent Analysis:")
		fmt.Printf("  Query: %s\n", resp.SearchIntent.Query)
		fmt.Printf("  Intent: %s\n", resp.SearchIntent.Intent)
		fmt.Printf("  Keywords: %s\n\n", resp.SearchIntent.Keywords)
	}

	// Display detailed results
	fmt.Printf("Results: %d found\n\n", len(resp.GetResults()))
	for i, result := range resp.GetResults() {
		fmt.Printf("─── Result %d ───\n", i+1)
		fmt.Printf("Title: %s\n", result.Title)
		fmt.Printf("URL: %s\n", result.Link)
		fmt.Printf("Source: %s\n", result.Media)
		fmt.Printf("Published: %s\n", result.PublishDate)
		fmt.Printf("Reference: %s\n", result.Refer)
		fmt.Printf("Content: %s\n", truncateString(result.Content, 200))

		if len(result.Images) > 0 {
			fmt.Printf("Images: %d attached\n", len(result.Images))
		}
		fmt.Println()
	}
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
