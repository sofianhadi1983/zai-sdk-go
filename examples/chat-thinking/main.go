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
	client, err := zai.NewClient(
		zai.WithAPIKey("your-secret"),
		zai.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	fmt.Println("=== GLM Thinking/Deep Thinking Examples ===")

	// Example 1: GLM-4.7 with native thinking (enabled by default)
	fmt.Println("Example 1: GLM-4.7 Native Thinking (Default)")
	glm47NativeThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 2: GLM-4.7 with thinking disabled
	fmt.Println("Example 2: GLM-4.7 with Thinking Disabled")
	glm47DisabledThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 2.5: GLM-4.7 with preserved thinking (multi-turn)
	fmt.Println("Example 2.5: GLM-4.7 with Preserved Thinking (Multi-turn)")
	glm47PreservedThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 3: Basic thinking with GLM-4-Plus
	fmt.Println("Example 3: Complex reasoning task")
	basicThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 4: Mathematical problem solving
	fmt.Println("Example 4: Mathematical reasoning")
	mathThinkingExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 5: Step-by-step analysis
	fmt.Println("Example 5: Step-by-step analysis")
	stepByStepExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 6: Streaming thinking process
	fmt.Println("Example 6: Streaming thinking (real-time)")
	streamingThinkingExample(ctx, client)
}

func glm47NativeThinkingExample(ctx context.Context, client *zai.Client) {
	// GLM-4.7 has thinking enabled by default
	// No need to explicitly enable it, but we can for clarity
	messages := []chat.Message{
		chat.NewUserMessage(`Solve this problem step by step:

A farmer has chickens and rabbits. Together they have 50 heads and 140 legs.
How many chickens and how many rabbits does the farmer have?

Please show your thinking process.`),
	}

	temp := 0.7
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4.7",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
		// Thinking is enabled by default for GLM-4.7, but we can explicitly enable it:
		// Thinking: &chat.ThinkingConfig{Type: chat.ThinkingTypeEnabled},
	}

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Display reasoning content if available
	if reasoning := resp.GetReasoningContent(); reasoning != "" {
		fmt.Println("Reasoning Process:")
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println(reasoning)
		fmt.Println(strings.Repeat("-", 60))
	}

	fmt.Println("\nFinal Answer:")
	fmt.Println(resp.GetContent())
	fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)
}

func glm47DisabledThinkingExample(ctx context.Context, client *zai.Client) {
	// Demonstrate disabling thinking for GLM-4.7
	messages := []chat.Message{
		chat.NewUserMessage(`Solve this problem:

A farmer has chickens and rabbits. Together they have 50 heads and 140 legs.
How many chickens and how many rabbits does the farmer have?`),
	}

	temp := 0.7
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4.7",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	// Disable thinking for faster, more direct responses
	req.DisableThinking()

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Response (with thinking disabled):")
	fmt.Println(resp.GetContent())
	fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)
}

func glm47PreservedThinkingExample(ctx context.Context, client *zai.Client) {
	// Demonstrate preserved thinking across multiple turns
	// This maintains reasoning continuity by preserving reasoning_content

	// Turn 1: Initial reasoning task
	messages := []chat.Message{
		chat.NewUserMessage(`Let's work on a complex problem together.

First, help me understand: If a store sells apples for $2 each and oranges for $3 each,
and someone spent $23 buying 9 fruits total, how many of each did they buy?`),
	}

	temp := 0.7
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4.7",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}
	// Enable preserved thinking for multi-turn reasoning continuity
	req.EnablePreservedThinking()

	fmt.Println("Turn 1: Initial Problem")
	resp1, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Assistant's Response:")
	fmt.Println(resp1.GetContent())
	if reasoning := resp1.GetReasoningContent(); reasoning != "" {
		fmt.Printf("\n[Reasoning preserved: %d chars]\n", len(reasoning))
	}

	// Turn 2: Follow-up question - preserve the reasoning_content from Turn 1
	assistantMsg := resp1.GetFirstChoice().Message
	messages = append(messages, assistantMsg) // Include full message with reasoning_content
	messages = append(messages, chat.NewUserMessage("Now, can you verify this answer by checking the totals?"))

	req.Messages = messages

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("Turn 2: Follow-up Verification")
	resp2, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Assistant's Response:")
	fmt.Println(resp2.GetContent())
	if reasoning := resp2.GetReasoningContent(); reasoning != "" {
		fmt.Printf("\n[Reasoning preserved: %d chars]\n", len(reasoning))
	}

	fmt.Println("\nNote: The model maintained reasoning continuity across both turns.")
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

	temp := 0.7
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4-plus",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

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

	temp := 0.5 // Lower temperature for more focused reasoning
	maxTokens := 1500
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4-plus",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

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

	temp := 0.6
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4-plus",
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Step-by-Step Analysis:")
	fmt.Println(resp.GetContent())
}

func streamingThinkingExample(ctx context.Context, client *zai.Client) {
	// Stream both reasoning content and final answer from GLM-4.7
	messages := []chat.Message{
		chat.NewUserMessage(`Design a simple algorithm to determine if a string is a palindrome.

Explain your thinking:
1. What approach will you use?
2. What are the edge cases?
3. What's the time complexity?`),
	}

	streamFlag := true
	temp := 0.7
	maxTokens := 2000
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4.7",
		Messages:    messages,
		Stream:      &streamFlag,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
		// Thinking enabled by default for GLM-4.7
	}

	stream, err := client.Chat.CreateStream(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer stream.Close()

	var fullReasoning string
	var fullResponse string
	var showingReasoning bool

	for stream.Next() {
		chunk := stream.Current()
		if chunk == nil {
			continue
		}

		// Check for reasoning content
		if reasoning := chunk.GetReasoningContent(); reasoning != "" {
			if !showingReasoning {
				fmt.Println("Reasoning Process (streaming):")
				fmt.Println(strings.Repeat("-", 60))
				showingReasoning = true
			}
			fmt.Print(reasoning)
			fullReasoning += reasoning
		}

		// Check for regular content
		if content := chunk.GetContent(); content != "" {
			if showingReasoning {
				fmt.Println("\n" + strings.Repeat("-", 60))
				fmt.Println("\nFinal Answer (streaming):")
				fmt.Println(strings.Repeat("-", 60))
				showingReasoning = false
			}
			fmt.Print(content)
			fullResponse += content
		}
	}

	if err := stream.Err(); err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf("\nReasoning: %d chars, Answer: %d chars\n", len(fullReasoning), len(fullResponse))
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

	temp := 0.7
	req := &chat.ChatCompletionRequest{
		Model:       "glm-4-plus",
		Messages:    messages,
		Tools:       tools,
		Temperature: &temp,
	}

	resp, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Reasoning with Tools:")

	// Check if the response has tool calls
	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		fmt.Println("\nModel's reasoning led to tool calls:")
		for i, toolCall := range resp.Choices[0].Message.ToolCalls {
			fmt.Printf("\nTool Call %d:\n", i+1)
			fmt.Printf("  Function: %s\n", toolCall.Function.Name)
			fmt.Printf("  Arguments: %s\n", toolCall.Function.Arguments)
		}
	}

	fmt.Println("\nResponse:")
	fmt.Println(resp.GetContent())
}
