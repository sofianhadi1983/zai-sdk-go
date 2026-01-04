# GLM Thinking/Deep Thinking Examples

This example demonstrates how to use GLM's advanced reasoning capabilities for complex problem-solving tasks.

## What is GLM Thinking?

GLM's thinking feature enables the model to:
- Show step-by-step reasoning
- Break down complex problems
- Verify its own work
- Provide detailed explanations of its thought process

**GLM-4.7** has thinking **enabled by default** as a native feature. You can control it with the `thinking` parameter.

## Features Demonstrated

### 1. GLM-4.7 Native Thinking (Default)
Demonstrates GLM-4.7's built-in thinking capability with `reasoning_content` output.

### 2. GLM-4.7 with Thinking Disabled
Shows how to disable thinking for faster, more direct responses using `DisableThinking()`.

### 2.5. GLM-4.7 with Preserved Thinking (Multi-turn)
Demonstrates `clear_thinking: false` to maintain reasoning continuity across conversation turns.

### 3. Basic Thinking (Complex Reasoning)
Uses system prompts to encourage step-by-step reasoning for logic problems with GLM-4-Plus.

### 4. Mathematical Reasoning
Shows how GLM can solve mathematical problems with detailed explanations.

### 5. Step-by-Step Analysis
Demonstrates breaking down business scenarios into analytical steps.

### 6. Streaming Thinking Process
Real-time streaming of both `reasoning_content` and final answer with GLM-4.7.

## Running the Example

### Prerequisites

Set your API key:
```bash
export ZAI_API_KEY="your-api-key.your-secret"
```

### Run the Example

```bash
cd examples/chat-thinking
go run main.go
```

## Key Techniques

### 1. Native Thinking Parameter (GLM-4.7)

GLM-4.7 has thinking **enabled by default**. You can explicitly control it:

```go
// Default behavior - thinking is enabled
req := &chat.ChatCompletionRequest{
    Model:    "glm-4.7",
    Messages: messages,
    // Thinking is enabled by default, no need to set it
}

// Explicitly enable thinking (optional, as it's the default)
req.EnableThinking()

// Disable thinking for faster, more direct responses
req.DisableThinking()

// Or set it directly with advanced options
clearThinking := false
req.SetThinking(&chat.ThinkingConfig{
    Type:          chat.ThinkingTypeEnabled,
    ClearThinking: &clearThinking, // Preserve reasoning across turns
})
```

#### Accessing Reasoning Content

When thinking is enabled, the response includes the model's reasoning process:

```go
resp, err := client.Chat.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Get the reasoning process
if reasoning := resp.GetReasoningContent(); reasoning != "" {
    fmt.Println("Reasoning:", reasoning)
}

// Get the final answer
fmt.Println("Answer:", resp.GetContent())
```

#### Preserved Thinking (Multi-turn Conversations)

Set `clear_thinking: false` to maintain reasoning continuity across turns:

```go
req := &chat.ChatCompletionRequest{
    Model:    "glm-4.7",
    Messages: messages,
}

// Enable preserved thinking (sets clear_thinking: false)
req.EnablePreservedThinking()

// First turn
resp1, _ := client.Chat.Create(ctx, req)

// IMPORTANT: Include the complete assistant message with reasoning_content
assistantMsg := resp1.GetFirstChoice().Message
messages = append(messages, assistantMsg)
messages = append(messages, chat.NewUserMessage("Follow-up question..."))

// Second turn - reasoning continuity is maintained
req.Messages = messages
resp2, _ := client.Chat.Create(ctx, req)
```

**When to disable thinking:**
- When you need faster responses
- For simple questions that don't require deep reasoning
- When you want more concise, direct answers

**When to use thinking (default):**
- Complex problem-solving
- Mathematical reasoning
- Step-by-step analysis
- Tasks requiring verification

**When to preserve thinking (clear_thinking: false):**
- Multi-turn problem-solving
- Complex conversations requiring reasoning continuity
- Debugging or step-by-step analysis across multiple messages

### 2. System Prompts for Reasoning (GLM-4-Plus)

For models without native thinking (like GLM-4-Plus), use system prompts:

```go
messages := []chat.Message{
    chat.NewSystemMessage("Think step-by-step before answering."),
    chat.NewUserMessage("Your complex problem here..."),
}
```

### 3. Temperature Settings

For reasoning tasks, use lower temperatures (0.5-0.7):
```go
temp := 0.5 // More focused reasoning
maxTokens := 2000
req := &chat.ChatCompletionRequest{
    Model:       "glm-4-plus",
    Messages:    messages,
    Temperature: &temp,
    MaxTokens:   &maxTokens,
}
```

### 4. Structured Prompting

Guide the model's thinking process:
```go
chat.NewUserMessage(`Solve this problem step by step:
1. Understand the question
2. Identify key components
3. Analyze step-by-step
4. Reach a conclusion
5. Verify your reasoning`)
```

### 5. Streaming for Real-Time Thinking

Stream both reasoning content and the final answer in real-time:

```go
stream, err := client.Chat.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

var fullReasoning string
var fullResponse string

for stream.Next() {
    chunk := stream.Current()
    if chunk == nil {
        continue
    }

    // Stream reasoning content
    if reasoning := chunk.GetReasoningContent(); reasoning != "" {
        fmt.Print(reasoning)
        fullReasoning += reasoning
    }

    // Stream final answer
    if content := chunk.GetContent(); content != "" {
        fmt.Print(content)
        fullResponse += content
    }
}

if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

This allows users to see the model's thinking process as it happens, improving transparency and user experience.

## Best Practices

### 1. Clear Instructions

Always specify that you want step-by-step reasoning:
- "Think step-by-step"
- "Show your work"
- "Explain your reasoning"

### 2. Appropriate Model Selection

Use `glm-4-plus` or `glm-4-air` for best reasoning capabilities.

### 3. Token Allocation

Reasoning tasks need more tokens:
```go
.SetMaxTokens(2000) // Or higher for complex problems
```

### 4. Temperature Control

- **0.3-0.5**: Very focused, deterministic reasoning
- **0.5-0.7**: Balanced reasoning with some creativity
- **0.7-0.9**: More creative problem-solving

## Example Problems

### Logic Puzzles
```
A farmer has chickens and rabbits. Together they have 50 heads
and 140 legs. How many chickens and how many rabbits?
```

### Mathematical Problems
```
Calculate the 15th Fibonacci number and explain the pattern.
```

### Business Analysis
```
Analyze revenue changes across quarters and calculate
overall percentage change.
```

### Algorithm Design
```
Design an algorithm to check if a string is a palindrome.
Explain your approach and complexity.
```

## Output Examples

### Step-by-Step Reasoning
```
Let me solve this step by step:

Step 1: Understanding the problem
- We have chickens (2 legs each) and rabbits (4 legs each)
- Total: 50 heads, 140 legs

Step 2: Setting up equations
- Let c = number of chickens, r = number of rabbits
- c + r = 50 (total heads)
- 2c + 4r = 140 (total legs)

Step 3: Solving
...
```

## Advanced Usage

### Combining with Function Calling

For even more powerful reasoning, combine thinking with tools:

```go
tools := []chat.Tool{
    chat.NewFunctionTool("calculate", "Perform calculations", params),
}

req := chat.NewChatCompletionRequest("glm-4-plus", messages).
    SetTools(tools)
```

The model will:
1. Think through the problem
2. Decide which tools to use
3. Execute calculations
4. Verify results

## Common Use Cases

1. **Educational Content**: Step-by-step problem solving
2. **Code Review**: Detailed analysis and explanations
3. **Research**: Complex analytical tasks
4. **Decision Making**: Multi-factor analysis
5. **Debugging**: Systematic problem diagnosis

## Troubleshooting

### Not Getting Step-by-Step Output?

1. Make sure your system prompt explicitly asks for reasoning
2. Increase max_tokens to allow full explanations
3. Try lower temperature for more structured thinking

### Responses Too Verbose?

1. Add constraints to your prompt: "Be concise but show key steps"
2. Adjust max_tokens
3. Use higher temperature for more natural flow

## Related Examples

- [Basic Chat](../chat/) - Simple chat completions
- [Function Calling](../tools/) - Using tools with reasoning
- [Streaming](../chat/) - Real-time response streaming

## API Reference

- [ChatCompletionRequest](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go/api/types/chat#ChatCompletionRequest)
- [Message Types](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go/api/types/chat#Message)
- [Streaming](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go/pkg/zai#ChatService.CreateStream)

## GLM-4.7 Native Thinking vs Prompting

### Native Thinking (GLM-4.7)

**Enabled by default** - GLM-4.7 has built-in thinking capability:

```go
req := &chat.ChatCompletionRequest{
    Model:    "glm-4.7",
    Messages: messages,
    // Thinking enabled by default
}

// Disable if needed
req.DisableThinking()
```

**Advantages:**
- No system prompts needed
- More consistent reasoning behavior
- Better integrated into the model
- Can be controlled via API parameter

### Prompting Approach (GLM-4-Plus, GLM-4-Air)

For models without native thinking, use system prompts:

```go
messages := []chat.Message{
    chat.NewSystemMessage("Think step-by-step before answering."),
    chat.NewUserMessage("Your problem..."),
}
```

**When to use:**
- With GLM-4-Plus or GLM-4-Air
- When you need custom reasoning instructions
- For fine-grained control over thinking style

## Notes

- **GLM-4.7** has native thinking controlled via the `thinking` parameter
- **Other models** (GLM-4-Plus, GLM-4-Air) use prompting for reasoning
- The quality of reasoning depends on the model version
- Streaming allows you to see the thinking process in real-time
- Combining with function calling enables verified computational reasoning
