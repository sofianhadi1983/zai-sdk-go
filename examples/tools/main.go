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

	// Example 6: Token counting
	fmt.Println("\n=== Example 6: Token Counting (Tokenizer) ===")
	tokenizerExample(ctx, client)
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

func tokenizerExample(ctx context.Context, client *zai.Client) {
	// Example: Count tokens before making an actual API call
	// This helps estimate costs and stay within token limits

	// Example 1: Simple message tokenization
	messages := []chat.Message{
		chat.NewSystemMessage("You are a helpful AI assistant specialized in explaining complex topics."),
		chat.NewUserMessage("Can you explain quantum computing in simple terms?"),
	}

	req := tools.NewTokenizerRequest("glm-4.6", messages)

	resp, err := client.Tools.Tokenizer(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Simple Message Tokenization:")
	fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
	fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	fmt.Printf("  Request ID: %s\n", resp.RequestID)

	// Example 2: Tokenization with function tools
	fmt.Println("\nTokenization with Function Tools:")

	toolDef := chat.NewFunctionTool(
		"get_weather",
		"Get the current weather in a given location",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "City name, e.g., San Francisco",
				},
				"unit": map[string]interface{}{
					"type": "string",
					"enum": []string{"celsius", "fahrenheit"},
				},
			},
			"required": []string{"location"},
		},
	)

	messagesWithTools := []chat.Message{
		chat.NewSystemMessage("You are a weather assistant."),
		chat.NewUserMessage("What's the weather like in Tokyo?"),
	}

	reqWithTools := tools.NewTokenizerRequest("glm-4.6", messagesWithTools).
		SetTools([]chat.Tool{toolDef})

	respWithTools, err := client.Tools.Tokenizer(ctx, reqWithTools)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("  Prompt tokens (with tool definitions): %d\n", respWithTools.Usage.PromptTokens)
	fmt.Printf("  Total tokens: %d\n", respWithTools.Usage.TotalTokens)
	fmt.Printf("  Extra tokens from tool definition: %d\n",
		respWithTools.Usage.PromptTokens-resp.Usage.PromptTokens)

	// Example 3: Multi-turn conversation tokenization
	fmt.Println("\nMulti-turn Conversation Tokenization:")

	conversation := []chat.Message{
		chat.NewSystemMessage("You are a coding tutor."),
		chat.NewUserMessage("How do I write a for loop in Python?"),
		chat.NewAssistantMessage("In Python, you can write a for loop like this:\n\nfor i in range(10):\n    print(i)"),
		chat.NewUserMessage("Can you explain what range() does?"),
	}

	convReq := tools.NewTokenizerRequest("glm-4.6", conversation)

	convResp, err := client.Tools.Tokenizer(ctx, convReq)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("  Messages in conversation: %d\n", len(conversation))
	fmt.Printf("  Total prompt tokens: %d\n", convResp.Usage.PromptTokens)
	fmt.Printf("  Average tokens per message: %.1f\n",
		float64(convResp.Usage.PromptTokens)/float64(len(conversation)))

	// Cost estimation (example pricing)
	fmt.Println("\n--- Cost Estimation Example ---")
	pricePerMToken := 0.01 // Example: $0.01 per 1M tokens
	estimatedCost := float64(convResp.Usage.TotalTokens) * pricePerMToken / 1_000_000
	fmt.Printf("  Estimated cost for this request: $%.6f\n", estimatedCost)
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
