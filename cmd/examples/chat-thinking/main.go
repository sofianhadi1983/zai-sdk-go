package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
	// Create client from environment variables
	client, err := zai.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	fmt.Println("=== GLM Thinking/Deep Thinking Examples ===\n")

	// Example 1: Basic thinking with GLM-4-Plus
	fmt.Println("Example 1: Complex reasoning task")
	basicThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 2: Mathematical problem solving
	fmt.Println("Example 2: Mathematical reasoning")
	mathThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 3: Step-by-step analysis
	fmt.Println("Example 3: Step-by-step analysis")
	stepByStepExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 4: Streaming thinking process
	fmt.Println("Example 4: Streaming thinking (real-time)")
	streamingThinkingExample(ctx, client)
}

func basicThinkingExample(ctx context.Context, client *zai.Client) {
	// For deep thinking, use GLM-4-Plus or GLM-4-Air with specific prompting
	messages := []chat.Message{
		chat.NewSystemMessage("You are a helpful AI assistant. Think step-by-step before answering."),
		chat.NewUserMessage(`Solve this problem step by step:

A farmer has chickens and rabbits. Together they have 50 heads and 140 legs.
How many chickens and how many rabbits does the farmer have?

Please show your thinking process.`),
	}

	req := chat.NewChatCompletionRequest("glm-4-plus", messages).
		SetTemperature(0.7).
		SetMaxTokens(2000)

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Response:")
	fmt.Println(resp.GetContent())
	fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)
}

func mathThinkingExample(ctx context.Context, client *zai.Client) {
	messages := []chat.Message{
		chat.NewSystemMessage(`You are a mathematical reasoning assistant.
Always show your work step-by-step, explaining your reasoning at each step.`),
		chat.NewUserMessage(`Calculate the 15th Fibonacci number and explain the pattern.

Requirements:
1. Show the calculation process
2. Explain the Fibonacci sequence pattern
3. Verify your answer`),
	}

	req := chat.NewChatCompletionRequest("glm-4-plus", messages).
		SetTemperature(0.5). // Lower temperature for more focused reasoning
		SetMaxTokens(1500)

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Mathematical Reasoning:")
	fmt.Println(resp.GetContent())
}

func stepByStepExample(ctx context.Context, client *zai.Client) {
	messages := []chat.Message{
		chat.NewSystemMessage(`You are an expert analyst. Break down complex problems into steps.

For each problem:
1. Understand the question
2. Identify key components
3. Analyze step-by-step
4. Reach a conclusion
5. Verify your reasoning`),
		chat.NewUserMessage(`Analyze this scenario:

A company's revenue increased by 20% in Q1, decreased by 10% in Q2,
increased by 15% in Q3, and decreased by 5% in Q4.

If they started with $1,000,000 revenue at the beginning of the year,
what was their final revenue? What was the overall percentage change?`),
	}

	req := chat.NewChatCompletionRequest("glm-4-plus", messages).
		SetTemperature(0.6).
		SetMaxTokens(2000)

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Step-by-Step Analysis:")
	fmt.Println(resp.GetContent())
}

func streamingThinkingExample(ctx context.Context, client *zai.Client) {
	messages := []chat.Message{
		chat.NewSystemMessage("Think through problems carefully, showing your reasoning process."),
		chat.NewUserMessage(`Design a simple algorithm to determine if a string is a palindrome.

Explain your thinking:
1. What approach will you use?
2. What are the edge cases?
3. What's the time complexity?
4. Show example code`),
	}

	req := chat.NewChatCompletionRequest("glm-4-plus", messages).
		SetStream(true).
		SetTemperature(0.7).
		SetMaxTokens(2000)

	stream, err := client.Chat.CreateStream(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer stream.Close()

	fmt.Println("Streaming thinking process:")
	fmt.Println(strings.Repeat("-", 60))

	var fullResponse string
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			content := chunk.GetContent()
			fmt.Print(content)
			fullResponse += content
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf("\nTotal characters streamed: %d\n", len(fullResponse))
}

// Example with function calling for complex reasoning
func reasoningWithToolsExample(ctx context.Context, client *zai.Client) {
	// Define a tool for calculations
	tools := []chat.Tool{
		chat.NewFunctionTool(
			"calculate",
			"Perform mathematical calculations",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]interface{}{
						"type":        "string",
						"description": "Mathematical expression to evaluate",
					},
				},
				"required": []string{"expression"},
			},
		),
	}

	messages := []chat.Message{
		chat.NewSystemMessage("You can use tools to help with calculations. Think through the problem first."),
		chat.NewUserMessage("Calculate the compound interest on $10,000 at 5% annual rate for 10 years, compounded annually."),
	}

	req := chat.NewChatCompletionRequest("glm-4-plus", messages).
		SetTools(tools).
		SetTemperature(0.7)

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Reasoning with Tools:")
	if resp.HasToolCalls() {
		fmt.Println("\nModel's reasoning led to tool calls:")
		for i, toolCall := range resp.GetToolCalls() {
			fmt.Printf("\nTool Call %d:\n", i+1)
			fmt.Printf("  Function: %s\n", toolCall.Function.Name)
			fmt.Printf("  Arguments: %s\n", toolCall.Function.Arguments)
		}
	}
	fmt.Println("\nResponse:")
	fmt.Println(resp.GetContent())
}
