# Z.ai Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/sofianhadi1983/zai-sdk-go.svg)](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sofianhadi1983/zai-sdk-go)](https://goreportcard.com/report/github.com/sofianhadi1983/zai-sdk-go)
[![Test](https://github.com/sofianhadi1983/zai-sdk-go/actions/workflows/test.yml/badge.svg)](https://github.com/sofianhadi1983/zai-sdk-go/actions/workflows/test.yml)
[![Security](https://github.com/sofianhadi1983/zai-sdk-go/actions/workflows/security.yml/badge.svg)](https://github.com/sofianhadi1983/zai-sdk-go/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/sofianhadi1983/zai-sdk-go/branch/main/graph/badge.svg)](https://codecov.io/gh/sofianhadi1983/zai-sdk-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Official Go SDK for the Z.ai API (智谱AI). This library provides a comprehensive and idiomatic Go interface to all Z.ai API endpoints, including chat completions, embeddings, image generation, and more.

## Features

- ✅ **Complete API Coverage**: All Z.ai API endpoints supported
- 🚀 **Streaming Support**: Real-time streaming for chat and agent interactions
- 🔧 **Type-Safe**: Full Go type definitions for requests and responses
- 🔄 **Automatic Retries**: Built-in retry logic with exponential backoff
- 🔐 **JWT Authentication**: Automatic token generation and caching
- 📦 **Builder Pattern**: Fluent API for easy request construction
- 🧪 **Well-Tested**: Comprehensive test coverage
- 📚 **Rich Examples**: Example code for all API features

### Supported APIs

**Core APIs:**
- **Chat Completions** - Text generation with streaming, function calling, and multimodal support
- **Embeddings** - Text embeddings with batch processing
- **Images** - Image generation with various models

**File & Media APIs:**
- **Files** - File upload, download, and management
- **Audio** - Audio transcription
- **Videos** - Video generation with async processing

**Advanced APIs:**
- **Assistant** - Conversational AI assistants with metadata
- **Batch** - Batch processing with pagination and cancellation
- **Web Search** - AI-powered web search with intent analysis
- **Moderations** - Content moderation and safety checks
- **Tools** - Function calling and tool execution

**Specialized APIs:**
- **Agents** - Agent invocation with streaming and async results
- **Voice** - Voice cloning and management
- **OCR** - Handwriting recognition with language support
- **File Parser** - Document parsing (async/sync)
- **Web Reader** - Web page content extraction

## Installation

```bash
go get github.com/sofianhadi1983/zai-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
    "github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
    // Create client with API key
    client, err := zai.NewClient(
        zai.WithAPIKey("your-api-key.your-secret"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create a chat completion request
    messages := []chat.Message{
        {Role: chat.RoleUser, Content: "Hello! How are you?"},
    }

    req := chat.NewChatCompletionRequest(chat.ModelGLM4Plus, messages)

    // Send the request
    resp, err := client.Chat.Create(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    // Print the response
    fmt.Println(resp.GetContent())
}
```

## Configuration

### Client Initialization

#### With API Key

```go
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
)
```

#### From Environment Variables

```go
// Reads ZAI_API_KEY and optionally ZAI_BASE_URL
client, err := zai.NewClientFromEnv()
```

#### For Chinese Users (Zhipu)

```go
client, err := zai.NewZhipuClient(
    zai.WithAPIKey("your-api-key.your-secret"),
)
```

### Configuration Options

```go
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
    zai.WithBaseURL("https://custom.api.url"),
    zai.WithTimeout(120 * time.Second),
    zai.WithMaxRetries(5),
    zai.WithDisableTokenCache(), // Use raw API key
    zai.WithLogger(customLogger),
)
```

## Usage Examples

### Chat Completions

#### Basic Chat

```go
messages := []chat.Message{
    {Role: chat.RoleUser, Content: "What is the capital of France?"},
}

req := chat.NewChatCompletionRequest(chat.ModelGLM4Plus, messages)

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.GetContent())
```

#### Streaming Chat

```go
messages := []chat.Message{
    {Role: chat.RoleUser, Content: "Write a short story"},
}

req := chat.NewChatCompletionRequest(chat.ModelGLM4Plus, messages).
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

```go
// Define a function/tool
tools := []chat.Tool{
    {
        Type: chat.ToolTypeFunction,
        Function: chat.Function{
            Name:        "get_weather",
            Description: "Get the current weather in a location",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "location": map[string]interface{}{
                        "type":        "string",
                        "description": "City name",
                    },
                },
                "required": []string{"location"},
            },
        },
    },
}

messages := []chat.Message{
    {Role: chat.RoleUser, Content: "What's the weather in Beijing?"},
}

req := chat.NewChatCompletionRequest(chat.ModelGLM4Plus, messages).
    SetTools(tools)

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Handle tool calls
if resp.HasToolCalls() {
    toolCalls := resp.GetToolCalls()
    fmt.Printf("Tool called: %s\n", toolCalls[0].Function.Name)
}
```

#### Multimodal (Image Input)

```go
messages := []chat.Message{
    {
        Role: chat.RoleUser,
        Content: []chat.ContentPart{
            {
                Type: chat.ContentTypeText,
                Text: "What's in this image?",
            },
            {
                Type:     chat.ContentTypeImageURL,
                ImageURL: &chat.ImageURL{URL: "https://example.com/image.jpg"},
            },
        },
    },
}

req := chat.NewChatCompletionRequest(chat.ModelGLM4VPlus, messages)

resp, err := client.Chat.Create(ctx, req)
```

### Embeddings

```go
req := embeddings.NewEmbeddingRequest(
    embeddings.ModelEmbedding3,
    []string{"Hello world", "How are you?"},
)

resp, err := client.Embeddings.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

for i, embedding := range resp.GetEmbeddings() {
    fmt.Printf("Embedding %d has %d dimensions\n", i, len(embedding))
}
```

### Image Generation

```go
req := images.NewImageGenerationRequest(
    "A beautiful sunset over mountains",
    images.ModelCogView3Plus,
)

resp, err := client.Images.Generate(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Image URL: %s\n", resp.GetImageURL())
```

### File Upload

```go
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

### Video Generation

```go
req := videos.NewVideoGenerationRequest(
    "A cat playing piano",
    videos.ModelCogVideoX,
)

resp, err := client.Videos.Generate(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Task ID: %s\n", resp.GetID())

// Poll for results
resultReq := videos.NewRetrieveRequest(resp.GetID())
result, err := client.Videos.Retrieve(ctx, resultReq)
```

### Web Search

```go
req := websearch.NewRequest("latest AI news").
    SetEnableIntent(true)

resp, err := client.WebSearch.Search(ctx, req)
if err != nil {
    log.Fatal(err)
}

for _, item := range resp.GetResults() {
    fmt.Printf("%s: %s\n", item.Title, item.Link)
}
```

### Content Moderation

```go
req := moderation.NewRequest("Text to check for safety")

resp, err := client.Moderations.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.Flagged {
    fmt.Println("Content flagged:", resp.GetCategories())
}
```

### Agent Invocation

```go
messages := []chat.Message{
    {Role: chat.RoleUser, Content: "Translate this to French: Hello"},
}

req := agents.NewAgentInvokeRequest("agent_translation", messages)

resp, err := client.Agents.Invoke(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.GetContent())
```

### Voice Cloning

```go
file, _ := os.Open("voice_sample.mp3")
defer file.Close()

req := voice.NewVoiceCloneRequest(
    "my_voice",
    "Sample text",
    "Preview text",
    "file_123", // uploaded file ID
    "voice-clone-v1",
)

resp, err := client.Voice.Clone(ctx, req)
```

### OCR (Handwriting Recognition)

```go
file, _ := os.Open("handwriting.jpg")
defer file.Close()

req := ocr.NewOCRRequest(file, "handwriting.jpg", ocr.ToolTypeHandWrite).
    SetProbability(true).
    SetLanguageType("zh-CN")

resp, err := client.OCR.HandwritingOCR(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Recognized text:", resp.GetText())
```

### File Parser

```go
file, _ := os.Open("document.pdf")
defer file.Close()

req := fileparser.NewCreateRequest(
    file,
    "document.pdf",
    "pdf",
    fileparser.ToolTypePrime,
)

resp, err := client.FileParser.Create(ctx, req)
```

### Web Reader

```go
req := webreader.NewRequest("https://example.com").
    SetReturnFormat("markdown").
    SetRetainImages(true).
    SetWithLinksSummary(true)

resp, err := client.WebReader.Read(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.HasResult() {
    result := resp.GetResult()
    fmt.Printf("Title: %s\n", result.GetTitle())
    fmt.Printf("Content: %s\n", result.GetContent())
}
```

## Error Handling

The SDK uses Go's standard error handling patterns:

```go
resp, err := client.Chat.Create(ctx, req)
if err != nil {
    // Check for specific error types
    if apiErr, ok := err.(*errors.APIError); ok {
        fmt.Printf("API Error: %s (code: %s)\n", apiErr.Message, apiErr.Code)
    } else if configErr, ok := err.(*errors.ConfigError); ok {
        fmt.Printf("Config Error: %s\n", configErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Environment Variables

- `ZAI_API_KEY` - Your Z.ai API key (format: "key.secret")
- `ZAI_BASE_URL` - (Optional) Custom API base URL

## Advanced Usage

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

config := &client.Config{
    APIKey:     "your-api-key.your-secret",
    HTTPClient: httpClient,
}

baseClient, err := client.NewBaseClient(config)
```

### Retry Configuration

```go
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
    zai.WithMaxRetries(5),        // Max 5 retry attempts
    zai.WithTimeout(120 * time.Second), // 120 second timeout
)
```

### Logging

```go
customLogger := logger.New(logger.LevelDebug)

client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
    zai.WithLogger(customLogger),
)
```

## Examples

Complete working examples are available in the [`cmd/examples`](cmd/examples) directory:

- [Chat](cmd/examples/chat) - Chat completions (basic, streaming, function calling)
- [Embeddings](cmd/examples/embeddings) - Text embeddings
- [Images](cmd/examples/images) - Image generation
- [Files](cmd/examples/files) - File upload and management
- [Videos](cmd/examples/videos) - Video generation
- [Audio](cmd/examples/audio) - Audio transcription
- [Assistant](cmd/examples/assistant) - AI assistants
- [Batch](cmd/examples/batch) - Batch processing
- [Web Search](cmd/examples/websearch) - Web search
- [Moderations](cmd/examples/moderations) - Content moderation
- [Tools](cmd/examples/tools) - Function calling
- [Agents](cmd/examples/agents) - Agent invocation
- [Voice](cmd/examples/voice) - Voice cloning
- [OCR](cmd/examples/ocr) - Handwriting recognition
- [File Parser](cmd/examples/fileparser) - Document parsing
- [Web Reader](cmd/examples/webreader) - Web content extraction

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/sofianhadi1983/zai-sdk-go.git
cd zai-sdk-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test -v ./pkg/zai -run TestChatService
```

### Running Examples

```bash
# Set your API key
export ZAI_API_KEY="your-api-key.your-secret"

# Run an example
go run cmd/examples/chat/main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [Official Z.ai API Documentation](https://open.bigmodel.cn/dev/api)
- [Go Package Documentation](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go)
- [Python SDK](https://github.com/zhipuai/zhipuai-sdk-python)
- [Issue Tracker](https://github.com/sofianhadi1983/zai-sdk-go/issues)

## Support

For issues and questions:
- Open an issue on [GitHub](https://github.com/sofianhadi1983/zai-sdk-go/issues)
- Check the [official Z.ai documentation](https://open.bigmodel.cn/dev/api)

---

**Note**: This is an unofficial Go SDK for the Z.ai API. For the official Python SDK, see [zhipuai-sdk-python](https://github.com/zhipuai/zhipuai-sdk-python).
