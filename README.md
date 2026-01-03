# Z.ai SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/z-ai/zai-sdk-go.svg)](https://pkg.go.dev/github.com/z-ai/zai-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/z-ai/zai-sdk-go)](https://goreportcard.com/report/github.com/z-ai/zai-sdk-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Go SDK for [Z.ai Open Platform](https://docs.z.ai/), making it easy for developers to integrate Z.ai's large language model APIs into their Go applications.

**Status:** 🚧 Under Development - Porting from [Python SDK v0.2.0](https://github.com/zhipuai/zhipuai-sdk-python-v4)

## Features

### 🤖 Chat Completions
- Standard chat completions with various models including `glm-4.7`
- Real-time streaming responses
- Function/tool calling capabilities
- Character role-playing with `charglm-3`
- Multimodal chat with vision models

### 🧠 Embeddings
- High-quality text embeddings
- Configurable dimensions
- Batch processing support

### 🎥 Video Generation
- Text-to-video generation
- Image-to-video generation
- Customizable quality, duration, and resolution

### 🎵 Audio Processing
- Speech transcription
- Multiple audio format support

### 🤝 Assistant API
- Structured conversation management
- Streaming conversations
- Rich metadata support

### 🔧 Additional Features
- Web search integration
- File management
- Batch operations
- Content moderation
- Image generation

## Installation

```bash
go get github.com/z-ai/zai-sdk-go
```

### Requirements

- Go 1.23 or higher

## Quick Start

### Get API Key

- **Overseas regions**: Visit [Z.ai Open Platform](https://docs.z.ai/)
- **Mainland China**: Visit [Zhipu AI Open Platform](https://www.bigmodel.cn/)

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/z-ai/zai-sdk-go/pkg/zai"
    "github.com/z-ai/zai-sdk-go/api/types/chat"
)

func main() {
    // For overseas users
    client, err := zai.NewClient(
        zai.WithAPIKey("your-api-key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // For Chinese users, use NewZhipuClient instead
    // client, err := zai.NewZhipuClient(zai.WithAPIKey("your-api-key"))

    ctx := context.Background()

    req := &chat.ChatCompletionRequest{
        Model: "glm-4.7",
        Messages: []chat.Message{
            {Role: "user", Content: "Hello, Z.ai!"},
        },
    }

    resp, err := client.Chat.Create(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

### Configuration

The SDK supports multiple configuration methods:

#### Environment Variables

```bash
export ZAI_API_KEY="your-api-key"
export ZAI_BASE_URL="https://api.z.ai/api/paas/v4/"  # Optional
```

```go
client, err := zai.NewClientFromEnv()
```

#### Code Configuration

```go
client, err := zai.NewClient(
    zai.WithAPIKey("your-api-key"),
    zai.WithBaseURL("https://api.z.ai/api/paas/v4/"),
    zai.WithTimeout(30*time.Second),
    zai.WithMaxRetries(3),
)
```

## Examples

See the [examples](./cmd/examples/) directory for complete examples:

- [Chat Completions](./cmd/examples/chat/)
- [Streaming Chat](./cmd/examples/chat/)
- [Embeddings](./cmd/examples/embeddings/)
- [Image Generation](./cmd/examples/images/)
- [File Upload](./cmd/examples/files/)

## Documentation

- [API Reference](https://pkg.go.dev/github.com/z-ai/zai-sdk-go)
- [Z.ai Official Documentation](https://docs.z.ai/)
- [Migration Guide from Python SDK](./MIGRATION.md)

## Development

### Building

```bash
make build
```

### Testing

```bash
# Run unit tests
make test

# Run tests with coverage
make test-cover

# Run integration tests (requires API key)
ZAI_API_KEY=your-key make test-integration
```

### Linting

```bash
make lint
```

### Formatting

```bash
make format
```

## Architecture

This SDK follows Clean Architecture principles with clear separation of concerns:

- `pkg/zai/` - Public SDK interface
- `internal/` - Internal implementation details
- `api/types/` - Request/Response type definitions
- `cmd/examples/` - Example applications

For detailed architecture information, see [CLAUDE.md](../z-ai-sdk-python/CLAUDE.md).

## Error Handling

```go
import "github.com/z-ai/zai-sdk-go/pkg/zai/errors"

resp, err := client.Chat.Create(ctx, req)
if err != nil {
    switch {
    case errors.IsAuthenticationError(err):
        // Handle authentication error
    case errors.IsRateLimitError(err):
        // Handle rate limit
    case errors.IsServerError(err):
        // Handle server error
    default:
        // Handle other errors
    }
}
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Support

For questions and support:
- GitHub Issues: [z-ai/zai-sdk-go/issues](https://github.com/z-ai/zai-sdk-go/issues)
- Email: user_feedback@z.ai
- Documentation: [docs.z.ai](https://docs.z.ai/)

## Acknowledgments

This Go SDK is a port of the official [Python SDK](https://github.com/zhipuai/zhipuai-sdk-python-v4) maintained by the Z.ai team.
