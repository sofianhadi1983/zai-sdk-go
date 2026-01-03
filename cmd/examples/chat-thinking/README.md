# GLM Thinking/Deep Thinking Examples

This example demonstrates how to use GLM's advanced reasoning capabilities for complex problem-solving tasks.

## What is GLM Thinking?

GLM's thinking feature enables the model to:
- Show step-by-step reasoning
- Break down complex problems
- Verify its own work
- Provide detailed explanations of its thought process

## Features Demonstrated

### 1. Basic Thinking (Complex Reasoning)
Uses system prompts to encourage step-by-step reasoning for logic problems.

### 2. Mathematical Reasoning
Shows how GLM can solve mathematical problems with detailed explanations.

### 3. Step-by-Step Analysis
Demonstrates breaking down business scenarios into analytical steps.

### 4. Streaming Thinking Process
Real-time streaming of the model's reasoning process.

## Running the Example

### Prerequisites

Set your API key:
```bash
export ZAI_API_KEY="your-api-key.your-secret"
```

### Run the Example

```bash
cd cmd/examples/chat-thinking
go run main.go
```

## Key Techniques

### 1. System Prompts for Reasoning

```go
messages := []chat.Message{
    chat.NewSystemMessage("Think step-by-step before answering."),
    chat.NewUserMessage("Your complex problem here..."),
}
```

### 2. Temperature Settings

For reasoning tasks, use lower temperatures (0.5-0.7):
```go
req := chat.NewChatCompletionRequest("glm-4-plus", messages).
    SetTemperature(0.5). // More focused reasoning
    SetMaxTokens(2000)
```

### 3. Structured Prompting

Guide the model's thinking process:
```go
chat.NewUserMessage(`Solve this problem step by step:
1. Understand the question
2. Identify key components
3. Analyze step-by-step
4. Reach a conclusion
5. Verify your reasoning`)
```

### 4. Streaming for Real-Time Thinking

```go
stream, err := client.Chat.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
    chunk := stream.Current()
    if chunk != nil {
        fmt.Print(chunk.GetContent()) // See thinking in real-time
    }
}
```

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

## Notes

- GLM's thinking capability is enhanced by clear prompting rather than a specific parameter
- The quality of reasoning depends on the model version (glm-4-plus recommended)
- Streaming allows you to see the thinking process in real-time
- Combining with function calling enables verified computational reasoning
