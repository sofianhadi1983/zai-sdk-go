package main

import (
	"context"
	"fmt"
	"log"

	"github.com/z-ai/zai-sdk-go/api/types/assistant"
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

	// Example 1: Query available assistants
	fmt.Println("=== Example 1: Query Available Assistants ===")
	queryAssistantsExample(ctx, client)

	// Example 2: Simple conversation
	fmt.Println("\n=== Example 2: Simple Conversation ===")
	simpleConversationExample(ctx, client)

	// Example 3: Continue existing conversation
	fmt.Println("\n=== Example 3: Continue Conversation ===")
	continueConversationExample(ctx, client)

	// Example 4: Streaming conversation
	fmt.Println("\n=== Example 4: Streaming Conversation ===")
	streamingConversationExample(ctx, client)

	// Example 5: Conversation with attachments
	fmt.Println("\n=== Example 5: Conversation with Attachments ===")
	conversationWithAttachmentsExample(ctx, client)

	// Example 6: Query conversation usage
	fmt.Println("\n=== Example 6: Query Conversation Usage ===")
	queryConversationUsageExample(ctx, client)
}

func queryAssistantsExample(ctx context.Context, client *zai.Client) {
	// Query all available assistants
	resp, err := client.Assistant.QuerySupport(ctx, nil)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d assistants:\n", len(resp.GetAssistants()))
	for _, asst := range resp.GetAssistants() {
		fmt.Printf("\nAssistant ID: %s\n", asst.AssistantID)
		fmt.Printf("  Name: %s\n", asst.Name)
		fmt.Printf("  Description: %s\n", asst.Description)
		fmt.Printf("  Tools: %v\n", asst.Tools)
		fmt.Printf("  Starter Prompts:\n")
		for _, prompt := range asst.StarterPrompts {
			fmt.Printf("    - %s\n", prompt)
		}
	}

	// Query specific assistants
	resp, err = client.Assistant.QuerySupport(ctx, []string{"asst_123", "asst_456"})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("\nQueried specific assistants: %d found\n", len(resp.GetAssistants()))
}

func simpleConversationExample(ctx context.Context, client *zai.Client) {
	// Create a simple conversation using convenience method
	resp, err := client.Assistant.CreateConversation(
		ctx,
		"asst_123", // Your assistant ID
		"Explain quantum computing in simple terms",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Conversation ID: %s\n", resp.ConversationID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", resp.GetText())

	if resp.Usage != nil {
		fmt.Printf("Tokens used: %d (prompt: %d, completion: %d)\n",
			resp.Usage.TotalTokens,
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens)
	}
}

func continueConversationExample(ctx context.Context, client *zai.Client) {
	// Continue an existing conversation
	resp, err := client.Assistant.ContinueConversation(
		ctx,
		"asst_123",    // Your assistant ID
		"conv_456",    // Existing conversation ID
		"Can you elaborate on quantum superposition?",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Conversation ID: %s\n", resp.ConversationID)
	fmt.Printf("Response: %s\n", resp.GetText())
}

func streamingConversationExample(ctx context.Context, client *zai.Client) {
	// Create a streaming conversation
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{
					Type: "text",
					Text: "Tell me a short story about a robot learning to paint",
				},
			},
		},
	}

	req := assistant.NewConversationRequest("asst_123", messages)
	req.SetStream(true)

	stream, err := client.Assistant.ConversationStream(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer stream.Close()

	fmt.Print("Streaming response: ")
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			fmt.Print(chunk.GetText())

			// Check if completed
			if chunk.IsCompleted() {
				fmt.Println("\n[Completed]")
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}
}

func conversationWithAttachmentsExample(ctx context.Context, client *zai.Client) {
	// Create a conversation with file attachments
	// First, upload a file (see files example)
	// file, _ := client.Files.Upload(ctx, fileReader, "document.pdf", "assistants")

	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{
					Type: "text",
					Text: "Please analyze this document and summarize the main points",
				},
			},
		},
	}

	req := assistant.NewConversationRequest("asst_123", messages)

	// Add file attachments
	attachments := []assistant.AssistantAttachment{
		{FileID: "file_123"}, // Use actual file ID from upload
	}
	req.SetAttachments(attachments)

	// Add metadata
	metadata := map[string]interface{}{
		"user_session": "session_789",
		"context":      "document_analysis",
	}
	req.SetMetadata(metadata)

	resp, err := client.Assistant.Conversation(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetText())
}

func queryConversationUsageExample(ctx context.Context, client *zai.Client) {
	// Query conversation history and usage
	page := 1
	pageSize := 10

	resp, err := client.Assistant.QueryConversationUsage(ctx, "asst_123", page, pageSize)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Conversation History (Page %d):\n", page)
	for _, conv := range resp.GetConversations() {
		fmt.Printf("\nConversation ID: %s\n", conv.ID)
		fmt.Printf("  Created: %d\n", conv.CreateTime)
		fmt.Printf("  Updated: %d\n", conv.UpdateTime)
		fmt.Printf("  Token Usage:\n")
		fmt.Printf("    Prompt: %d\n", conv.Usage.PromptTokens)
		fmt.Printf("    Completion: %d\n", conv.Usage.CompletionTokens)
		fmt.Printf("    Total: %d\n", conv.Usage.TotalTokens)
	}

	if resp.HasMore() {
		fmt.Println("\nMore conversations available. Fetch next page...")
		// Fetch next page
		nextPage := page + 1
		nextResp, err := client.Assistant.QueryConversationUsage(ctx, "asst_123", nextPage, pageSize)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		fmt.Printf("Page %d: %d conversations\n", nextPage, len(nextResp.GetConversations()))
	}
}

// Advanced example: Multi-turn conversation with translation
func multiTurnConversationWithTranslation(ctx context.Context, client *zai.Client) {
	// Create conversation with translation
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{
					Type: "text",
					Text: "Hello, how are you?",
				},
			},
		},
	}

	req := assistant.NewConversationRequest("asst_123", messages)

	// Set translation parameters
	extraParams := &assistant.ExtraParameters{
		Translate: &assistant.TranslateParameters{
			FromLanguage: "en",
			ToLanguage:   "zh",
		},
	}
	req.SetExtraParameters(extraParams)

	resp, err := client.Assistant.Conversation(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Translated Response: %s\n", resp.GetText())
}

// Error handling example
func errorHandlingExample(ctx context.Context, client *zai.Client) {
	// Example 1: Handle invalid assistant ID
	fmt.Println("1. Testing with invalid assistant ID:")
	resp, err := client.Assistant.CreateConversation(ctx, "invalid_id", "Hello")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 2: Handle failed conversation
	fmt.Println("\n2. Checking conversation status:")
	resp, err = client.Assistant.CreateConversation(ctx, "asst_123", "Test")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}

	if resp.IsFailed() {
		fmt.Printf("   Conversation failed: %s\n", resp.GetError())
	} else if resp.IsCompleted() {
		fmt.Printf("   ✓ Conversation completed successfully\n")
	} else if resp.IsInProgress() {
		fmt.Printf("   Conversation still in progress\n")
	}
}
