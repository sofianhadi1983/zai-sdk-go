# Migration Guide: Python SDK to Go SDK

This guide helps developers migrate from the Z.ai Python SDK to the Go SDK.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Client Initialization](#client-initialization)
- [API Differences](#api-differences)
  - [Chat Completions](#chat-completions)
  - [Embeddings](#embeddings)
  - [Image Generation](#image-generation)
  - [File Operations](#file-operations)
  - [Other APIs](#other-apis)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## Overview

The Z.ai Go SDK provides the same functionality as the Python SDK but follows Go conventions and idioms. The main differences are:

1. **Type Safety**: Go SDK uses strong typing instead of dictionaries
2. **Builder Pattern**: Fluent API for request construction
3. **Context Support**: All methods accept `context.Context` for cancellation
4. **Streaming**: Go SDK uses channel-based streaming
5. **Error Handling**: Explicit error returns instead of exceptions

## Installation

### Python SDK
```bash
pip install zhipuai
```

### Go SDK
```bash
go get github.com/sofianhadi1983/zai-sdk-go
```

## Client Initialization

### Python SDK
```python
from zhipuai import ZhipuAI

# Initialize client
client = ZhipuAI(api_key="your-api-key.your-secret")
```

### Go SDK
```go
import (
    "github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

// Initialize client
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### From Environment Variables

**Python:**
```python
import os
from zhipuai import ZhipuAI

client = ZhipuAI(api_key=os.getenv("ZHIPUAI_API_KEY"))
```

**Go:**
```go
// Set environment variable ZAI_API_KEY
client, err := zai.NewClientFromEnv()
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## API Differences

### Chat Completions

#### Basic Chat

**Python:**
```python
response = client.chat.completions.create(
    model="glm-4.7",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)
print(response.choices[0].message.content)
```

**Go:**
```go
import "github.com/sofianhadi1983/zai-sdk-go/api/types/chat"

messages := []chat.Message{
    chat.NewUserMessage("Hello!"),
}
req := chat.NewChatCompletionRequest(chat.ModelGLM47, messages)

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.GetContent())
```

#### Streaming Chat

**Python:**
```python
response = client.chat.completions.create(
    model="glm-4.7",
    messages=[{"role": "user", "content": "Tell me a story"}],
    stream=True
)
for chunk in response:
    if chunk.choices:
        print(chunk.choices[0].delta.content, end="")
```

**Go:**
```go
messages := []chat.Message{
    chat.NewUserMessage("Tell me a story"),
}
req := chat.NewChatCompletionRequest(chat.ModelGLM47, messages).
    SetStream(true)

stream, err := client.Chat.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
    chunk := stream.Current()
    if chunk != nil {
        fmt.Print(chunk.GetContent())
    }
}

if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

#### Function Calling

**Python:**
```python
tools = [{
    "type": "function",
    "function": {
        "name": "get_weather",
        "description": "Get weather for a location",
        "parameters": {
            "type": "object",
            "properties": {
                "location": {"type": "string"}
            }
        }
    }
}]

response = client.chat.completions.create(
    model="glm-4.7",
    messages=[{"role": "user", "content": "What's the weather in Beijing?"}],
    tools=tools
)
```

**Go:**
```go
tools := []chat.Tool{
    chat.NewFunctionTool(
        "get_weather",
        "Get weather for a location",
        map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "location": map[string]interface{}{
                    "type": "string",
                },
            },
        },
    ),
}

messages := []chat.Message{
    chat.NewUserMessage("What's the weather in Beijing?"),
}
req := chat.NewChatCompletionRequest(chat.ModelGLM47, messages).
    SetTools(tools)

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.HasToolCalls() {
    for _, toolCall := range resp.GetToolCalls() {
        fmt.Printf("Function: %s\n", toolCall.Function.Name)
    }
}
```

### Embeddings

**Python:**
```python
response = client.embeddings.create(
    model="embedding-2",
    input="Hello world"
)
embedding = response.data[0].embedding
```

**Go:**
```go
import "github.com/sofianhadi1983/zai-sdk-go/api/types/embeddings"

req := embeddings.NewEmbeddingRequest(
    embeddings.ModelEmbedding3,
    "Hello world",
)
resp, err := client.Embeddings.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}
embedding := resp.GetFirstEmbedding().GetFloatEmbedding()
```

#### Batch Embeddings

**Python:**
```python
response = client.embeddings.create(
    model="embedding-2",
    input=["Hello", "World"]
)
for data in response.data:
    print(len(data.embedding))
```

**Go:**
```go
texts := []string{"Hello", "World"}
req := embeddings.NewBatchEmbeddingRequest(
    embeddings.ModelEmbedding3,
    texts,
)
resp, err := client.Embeddings.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}
for _, emb := range resp.GetFloatEmbeddings() {
    fmt.Printf("%d dimensions\n", len(emb))
}
```

### Image Generation

**Python:**
```python
response = client.images.generations(
    model="cogview-3",
    prompt="A beautiful sunset"
)
image_url = response.data[0].url
```

**Go:**
```go
import "github.com/sofianhadi1983/zai-sdk-go/api/types/images"

req := images.NewImageGenerationRequest(
    images.ModelCogView3Plus,
    "A beautiful sunset",
)
resp, err := client.Images.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}
imageURL := resp.GetImageURL()
```

### File Operations

#### Upload File

**Python:**
```python
with open("document.pdf", "rb") as f:
    response = client.files.create(
        file=f,
        purpose="fine-tune"
    )
print(response.id)
```

**Go:**
```go
import "github.com/sofianhadi1983/zai-sdk-go/api/types/files"

file, err := os.Open("document.pdf")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

req := files.NewCreateRequest(file, "document.pdf", files.PurposeFineTune)
resp, err := client.Files.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File ID: %s\n", resp.ID)
```

### Other APIs

#### Web Search

**Python:**
```python
response = client.web_search(
    query="latest AI news",
    enable_intent=True
)
for result in response.data:
    print(result.title)
```

**Go:**
```go
import "github.com/sofianhadi1983/zai-sdk-go/api/types/websearch"

req := websearch.NewRequest("latest AI news").
    SetEnableIntent(true)

resp, err := client.WebSearch.Search(ctx, req)
if err != nil {
    log.Fatal(err)
}
for _, item := range resp.GetResults() {
    fmt.Println(item.Title)
}
```

## Error Handling

### Python SDK
```python
try:
    response = client.chat.completions.create(...)
except Exception as e:
    print(f"Error: {e}")
```

### Go SDK
```go
import "github.com/sofianhadi1983/zai-sdk-go/pkg/zai/errors"

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    if apiErr, ok := err.(*errors.APIError); ok {
        fmt.Printf("API Error: %s (code: %s)\n", apiErr.Message, apiErr.Code)
    } else if errors.IsAuthenticationError(err) {
        fmt.Println("Authentication failed")
    } else if errors.IsRateLimitError(err) {
        fmt.Println("Rate limit exceeded")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Best Practices

### 1. Context Management

**Always use context for cancellation and timeouts:**

```go
import "context"

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.Chat.Create(ctx, req)
```

### 2. Resource Cleanup

**Always close clients and streams:**

```go
client, err := zai.NewClient(zai.WithAPIKey("key"))
if err != nil {
    log.Fatal(err)
}
defer client.Close()  // Always close the client

stream, err := client.Chat.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()  // Always close streams
```

### 3. Builder Pattern

**Use method chaining for cleaner code:**

```go
req := chat.NewChatCompletionRequest(chat.ModelGLM47, messages).
    SetTemperature(0.8).
    SetMaxTokens(1000).
    SetTopP(0.9)
```

### 4. Type Safety

**Leverage Go's type system:**

```go
// Use constants instead of strings
req := chat.NewChatCompletionRequest(chat.ModelGLM47, messages)

// Use typed message constructors
msg := chat.NewUserMessage("Hello")
```

### 5. Error Handling

**Check errors explicitly:**

```go
resp, err := client.Chat.Create(ctx, req)
if err != nil {
    // Handle error
    return
}
// Use response
```

## Common Pitfalls

### 1. Forgetting Error Checks

**❌ Wrong:**
```go
resp, _ := client.Chat.Create(ctx, req)  // Ignoring errors
```

**✅ Correct:**
```go
resp, err := client.Chat.Create(ctx, req)
if err != nil {
    return err
}
```

### 2. Not Closing Resources

**❌ Wrong:**
```go
stream, _ := client.Chat.CreateStream(ctx, req)
for stream.Next() {
    // ...
}
// Stream not closed!
```

**✅ Correct:**
```go
stream, err := client.Chat.CreateStream(ctx, req)
if err != nil {
    return err
}
defer stream.Close()

for stream.Next() {
    // ...
}
```

### 3. Missing Context

**❌ Wrong:**
```go
// No timeout or cancellation
resp, err := client.Chat.Create(context.Background(), req)
```

**✅ Correct:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
resp, err := client.Chat.Create(ctx, req)
```

## Additional Resources

- [Go SDK Documentation](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go)
- [Python SDK Documentation](https://github.com/zhipuai/zhipuai-sdk-python)
- [Official Z.ai API Docs](https://open.bigmodel.cn/dev/api)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)

## Getting Help

If you have questions or issues:

- Check the [examples directory](examples/) for working code
- Open an issue on [GitHub](https://github.com/sofianhadi1983/zai-sdk-go/issues)
- Review the [Go SDK API documentation](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go)
