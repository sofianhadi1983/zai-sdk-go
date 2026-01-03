package main

import (
	"context"
	"fmt"
	"log"

	"github.com/z-ai/zai-sdk-go/api/types/chat"
	"github.com/z-ai/zai-sdk-go/pkg/zai"
)

func main() {
	// Create a new client
	// You can also use zai.NewClientFromEnv() to load from environment variables
	client, err := zai.NewClient(
		zai.WithAPIKey("your-api-key.your-secret"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// For Chinese users, use:
	// client, err := zai.NewZhipuClient(zai.WithAPIKey("your-api-key.your-secret"))

	ctx := context.Background()

	// Example 1: Basic chat completion
	fmt.Println("=== Example 1: Basic Chat ===")
	basicExample(ctx, client)

	fmt.Println("\n=== Example 2: Streaming Chat ===")
	streamingExample(ctx, client)

	fmt.Println("\n=== Example 3: Multi-turn Conversation ===")
	conversationExample(ctx, client)

	fmt.Println("\n=== Example 4: With Parameters ===")
	parametersExample(ctx, client)
}

func basicExample(ctx context.Context, client *zai.Client) {
	// Create a simple chat completion request
	req := &chat.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []chat.Message{
			chat.NewUserMessage("Hello, Z.ai!"),
		},
	}

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetContent())
	fmt.Printf("Model: %s\n", resp.Model)
	if resp.Usage != nil {
		fmt.Printf("Tokens used: %d (prompt: %d, completion: %d)\n",
			resp.Usage.TotalTokens,
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens)
	}
}

func streamingExample(ctx context.Context, client *zai.Client) {
	// Create a streaming request
	req := &chat.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []chat.Message{
			chat.NewUserMessage("Tell me a short joke"),
		},
	}

	stream, err := client.Chat.CreateStream(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer stream.Close()

	fmt.Print("Assistant: ")
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			fmt.Print(chunk.GetContent())
		}
	}
	fmt.Println()

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}
}

func conversationExample(ctx context.Context, client *zai.Client) {
	// Build a multi-turn conversation
	req := &chat.ChatCompletionRequest{
		Model: "glm-4.7",
	}

	// Add conversation history using chained methods
	req.AddSystemMessage("You are a helpful assistant").
		AddUserMessage("What is the capital of France?").
		AddAssistantMessage("The capital of France is Paris.").
		AddUserMessage("What is it famous for?")

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Assistant: %s\n", resp.GetContent())
}

func parametersExample(ctx context.Context, client *zai.Client) {
	// Create a request with custom parameters
	req := &chat.ChatCompletionRequest{
		Model: "glm-4.7",
		Messages: []chat.Message{
			chat.NewSystemMessage("You are a creative writer"),
			chat.NewUserMessage("Write a one-sentence story"),
		},
	}

	// Set parameters using chained methods
	req.SetTemperature(0.9).     // Higher temperature for more creative output
					SetMaxTokens(100).       // Limit response length
					SetTopP(0.95)            // Nucleus sampling

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Story: %s\n", resp.GetContent())

	choice := resp.GetFirstChoice()
	if choice != nil {
		fmt.Printf("Finish reason: %s\n", choice.FinishReason)
	}
}
