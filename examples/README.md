# Z.ai Go SDK Examples

This directory contains example applications demonstrating how to use the Z.ai Go SDK for various use cases.

## Prerequisites

Before running any examples, you need to set your API key as an environment variable:

```bash
export ZAI_API_KEY="your-api-key.your-secret"
```

## Available Examples

### Core APIs

#### [Chat Completions](chat/)
Demonstrates chat completion capabilities including basic chat, streaming, multi-turn conversations, and parameter customization.

```bash
cd chat
go run main.go
```

#### [Embeddings](embeddings/)
Shows how to generate text embeddings for single texts and batches, with various models and dimensions.

```bash
cd embeddings
go run main.go
```

#### [Images](images/)
Examples of image generation from text prompts with different models, sizes, and quality settings.

```bash
cd images
go run main.go
```

### File & Media APIs

#### [Files](files/)
Demonstrates file upload, download, listing, and deletion operations.

```bash
cd files
go run main.go
```

#### [Audio](audio/)
Shows audio transcription capabilities for various audio formats.

```bash
cd audio
go run main.go
```

#### [Videos](videos/)
Examples of video generation with async processing and result retrieval.

```bash
cd videos
go run main.go
```

### Advanced APIs

#### [Assistant](assistant/)
Demonstrates conversational AI assistant features with conversation management and metadata.

```bash
cd assistant
go run main.go
```

#### [Batch](batch/)
Shows batch processing operations with pagination and cancellation.

```bash
cd batch
go run main.go
```

#### [Web Search](websearch/)
Examples of AI-powered web search with intent analysis.

```bash
cd websearch
go run main.go
```

#### [Moderations](moderations/)
Demonstrates content moderation and safety checks.

```bash
cd moderations
go run main.go
```

#### [Tools](tools/)
Shows function calling capabilities with tool execution.

```bash
cd tools
go run main.go
```

### Specialized APIs

#### [Voice](voice/)
Examples of voice cloning and voice management operations.

```bash
cd voice
go run main.go
```

#### [OCR](ocr/)
Demonstrates handwriting recognition with language support.

```bash
cd ocr
go run main.go
```

#### [File Parser](fileparser/)
Shows document parsing capabilities for various file formats (async/sync).

```bash
cd fileparser
go run main.go
```

#### [Web Reader](webreader/)
Examples of web page content extraction with different formats and options.

```bash
cd webreader
go run main.go
```

## Running Examples

### Using Environment Variables

Most examples support loading credentials from environment variables:

```bash
export ZAI_API_KEY="your-api-key.your-secret"
export ZAI_BASE_URL="https://open.bigmodel.cn/api/paas/v4/"  # optional
cd <example-directory>
go run main.go
```

### Using Custom API Keys

You can also modify the example code to use hardcoded API keys (not recommended for production):

```go
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key.your-secret"),
)
```

### For Chinese Users

If you're in mainland China, use the Zhipu client:

```go
client, err := zai.NewZhipuClient(
    zai.WithAPIKey("your-api-key.your-secret"),
)
```

## Example Structure

Each example follows a consistent structure:

1. **Client Creation**: Shows how to create and configure the SDK client
2. **Multiple Examples**: Demonstrates different use cases for the API
3. **Error Handling**: Shows proper error handling patterns
4. **Resource Cleanup**: Uses `defer` to ensure proper cleanup
5. **Comments**: Includes detailed comments explaining each step

## Common Patterns

### Basic Request/Response

```go
req := resource.NewRequest(params)
resp, err := client.Resource.Method(ctx, req)
if err != nil {
    log.Fatal(err)
}
// Use response
```

### Streaming

```go
req := resource.NewRequest(params).SetStream(true)
stream, err := client.Resource.MethodStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for stream.Next() {
    chunk := stream.Current()
    // Process chunk
}

if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

### Builder Pattern

```go
req := resource.NewRequest(params).
    SetOption1(value1).
    SetOption2(value2).
    SetOption3(value3)
```

## Additional Resources

- [API Documentation](https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go)
- [Official Z.ai Docs](https://open.bigmodel.cn/dev/api)
- [GitHub Repository](https://github.com/sofianhadi1983/zai-sdk-go)
- [Python SDK](https://github.com/zhipuai/zhipuai-sdk-python)

## Contributing

If you have additional example ideas or improvements:

1. Create a new directory under `examples/`
2. Add a `main.go` with clear comments and error handling
3. Test the example with real API credentials
4. Submit a pull request

## Support

For issues or questions:
- Open an issue on [GitHub](https://github.com/sofianhadi1983/zai-sdk-go/issues)
- Check the [official Z.ai documentation](https://open.bigmodel.cn/dev/api)
