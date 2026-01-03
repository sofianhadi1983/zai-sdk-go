package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/sofianhadi1983/zai-sdk-go/api/types/tools"
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

	// Example 1: Basic web search
	fmt.Println("=== Example 1: Basic Web Search ===")
	basicWebSearchExample(ctx, client)

	// Example 2: Web search with filters
	fmt.Println("\n=== Example 2: Web Search with Filters ===")
	webSearchWithFiltersExample(ctx, client)

	// Example 3: Academic search
	fmt.Println("\n=== Example 3: Academic Search ===")
	academicSearchExample(ctx, client)

	// Example 4: Streaming web search
	fmt.Println("\n=== Example 4: Streaming Web Search ===")
	streamingWebSearchExample(ctx, client)

	// Example 5: Analyzing search results
	fmt.Println("\n=== Example 5: Analyzing Search Results ===")
	analyzeSearchResultsExample(ctx, client)
}

func basicWebSearchExample(ctx context.Context, client *zai.Client) {
	// Create a simple web search request
	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What are the latest breakthroughs in artificial intelligence?",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages)

	resp, err := client.Tools.WebSearch(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Search ID: %s\n", resp.ID)
	fmt.Printf("Request ID: %s\n", resp.RequestID)

	// Display search intents
	intents := resp.GetSearchIntents()
	if len(intents) > 0 {
		fmt.Println("\nSearch Intent Analysis:")
		for _, intent := range intents {
			fmt.Printf("  Optimized Query: %s\n", intent.Query)
			fmt.Printf("  Intent Type: %s\n", intent.Intent)
			fmt.Printf("  Keywords: %s\n", intent.Keywords)
		}
	}

	// Display search results
	results := resp.GetSearchResults()
	fmt.Printf("\nSearch Results (%d found):\n", len(results))
	for i, result := range results {
		fmt.Printf("\n%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)
		fmt.Printf("   Source: %s\n", result.Media)
		fmt.Printf("   Reference: %s\n", result.Refer)
		fmt.Printf("   Snippet: %s\n", truncateString(result.Content, 100))
	}

	// Display recommendations
	recommends := resp.GetSearchRecommendations()
	if len(recommends) > 0 {
		fmt.Println("\nRecommended Queries:")
		for _, recommend := range recommends {
			fmt.Printf("  - %s\n", recommend.Query)
		}
	}
}

func webSearchWithFiltersExample(ctx context.Context, client *zai.Client) {
	// Web search with scope and recent days filter
	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Latest news about quantum computing",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages).
		SetScope("news").
		SetRecentDays(7).
		SetLocation("US").
		SetRequestID("req_filtered_123")

	resp, err := client.Tools.WebSearch(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Search with filters (Scope: news, Recent: 7 days)\n")
	fmt.Printf("Request ID: %s\n\n", resp.RequestID)

	results := resp.GetSearchResults()
	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   Source: %s\n", result.Media)
		fmt.Printf("   Link: %s\n", result.Link)
		fmt.Println()
	}
}

func academicSearchExample(ctx context.Context, client *zai.Client) {
	// Search academic sources
	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Recent research papers on machine learning optimization",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages).
		SetScope("academic").
		SetRecentDays(30)

	resp, err := client.Tools.WebSearch(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Academic Search Results:")

	intents := resp.GetSearchIntents()
	if len(intents) > 0 {
		fmt.Printf("Optimized Query: %s\n\n", intents[0].Query)
	}

	results := resp.GetSearchResults()
	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   %s\n", result.Link)
		fmt.Printf("   Abstract: %s\n\n", truncateString(result.Content, 150))
	}
}

func streamingWebSearchExample(ctx context.Context, client *zai.Client) {
	// Streaming web search
	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What's happening in technology today?",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages).
		SetStream(true).
		SetRecentDays(1)

	stream, err := client.Tools.WebSearchStream(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer stream.Close()

	fmt.Println("Streaming search results...")

	var intentDisplayed bool
	var resultCount int

	for stream.Next() {
		chunk := stream.Current()
		if chunk == nil {
			continue
		}

		// Process chunk
		for _, choice := range chunk.Choices {
			// Display role if present
			if choice.Delta.Role != "" && !intentDisplayed {
				fmt.Printf("\nRole: %s\n", choice.Delta.Role)
			}

			// Process tool calls
			for _, toolCall := range choice.Delta.ToolCalls {
				// Display search intent
				if toolCall.SearchIntent != nil && !intentDisplayed {
					fmt.Println("\n--- Search Intent ---")
					fmt.Printf("Query: %s\n", toolCall.SearchIntent.Query)
					fmt.Printf("Intent: %s\n", toolCall.SearchIntent.Intent)
					fmt.Printf("Keywords: %s\n", toolCall.SearchIntent.Keywords)
					intentDisplayed = true
				}

				// Display search result
				if toolCall.SearchResult != nil {
					resultCount++
					fmt.Printf("\n--- Result %d ---\n", resultCount)
					fmt.Printf("Title: %s\n", toolCall.SearchResult.Title)
					fmt.Printf("Link: %s\n", toolCall.SearchResult.Link)
					fmt.Printf("Source: %s\n", toolCall.SearchResult.Media)
					fmt.Printf("Content: %s\n", truncateString(toolCall.SearchResult.Content, 80))
				}

				// Display recommendation
				if toolCall.SearchRecommend != nil {
					fmt.Printf("\n--- Recommendation ---\n")
					fmt.Printf("Try: %s\n", toolCall.SearchRecommend.Query)
				}
			}

			// Display finish reason
			if choice.FinishReason != "" {
				fmt.Printf("\nFinished: %s\n", choice.FinishReason)
			}
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}

	fmt.Printf("\nTotal results streamed: %d\n", resultCount)
}

func analyzeSearchResultsExample(ctx context.Context, client *zai.Client) {
	// Analyze and categorize search results
	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Best practices for software architecture",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages)

	resp, err := client.Tools.WebSearch(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Search Result Analysis:")

	// Analyze intents
	intents := resp.GetSearchIntents()
	if len(intents) > 0 {
		fmt.Printf("\nIntent Analysis:\n")
		for _, intent := range intents {
			fmt.Printf("  Type: %s\n", intent.Intent)
			fmt.Printf("  Optimized: %s\n", intent.Query)
		}
	}

	// Categorize results by source
	results := resp.GetSearchResults()
	sourceMap := make(map[string][]tools.SearchResult)

	for _, result := range results {
		media := result.Media
		if media == "" {
			media = "Unknown"
		}
		sourceMap[media] = append(sourceMap[media], *result)
	}

	fmt.Printf("\nResults by Source:\n")
	for source, sourceResults := range sourceMap {
		fmt.Printf("\n%s (%d results):\n", source, len(sourceResults))
		for i, result := range sourceResults {
			fmt.Printf("  %d. %s\n", i+1, result.Title)
			fmt.Printf("     %s\n", result.Link)
		}
	}

	// Show recommendations
	recommends := resp.GetSearchRecommendations()
	if len(recommends) > 0 {
		fmt.Printf("\nRelated Searches:\n")
		for i, recommend := range recommends {
			fmt.Printf("  %d. %s\n", i+1, recommend.Query)
		}
	}

	// Summary
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Total Results: %d\n", len(results))
	fmt.Printf("Unique Sources: %d\n", len(sourceMap))
	fmt.Printf("Recommendations: %d\n", len(recommends))
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
