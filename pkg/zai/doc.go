// Package zai provides the official Go SDK for Z.ai Open Platform APIs.
//
// The SDK enables Go developers to easily integrate with Z.ai's large language
// model APIs, including chat completions, embeddings, image generation, and more.
//
// # Installation
//
//	go get github.com/z-ai/zai-sdk-go
//
// # Quick Start
//
// For international users:
//
//	client, err := zai.NewClient(zai.WithAPIKey("your-api-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// For users in mainland China:
//
//	client, err := zai.NewZhipuClient(zai.WithAPIKey("your-api-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Configuration
//
// The SDK can be configured using functional options:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-api-key"),
//	    zai.WithBaseURL("https://custom.api.url"),
//	    zai.WithTimeout(30*time.Second),
//	    zai.WithMaxRetries(3),
//	)
//
// Or using environment variables:
//
//	export ZAI_API_KEY="your-api-key"
//	export ZAI_BASE_URL="https://api.z.ai/api/paas/v4/"
//
//	client, err := zai.NewClientFromEnv()
//
// # Error Handling
//
// All API errors implement the error interface and can be checked using
// the errors package:
//
//	import "github.com/z-ai/zai-sdk-go/pkg/zai/errors"
//
//	resp, err := client.Chat.Create(ctx, req)
//	if err != nil {
//	    if errors.IsAuthenticationError(err) {
//	        // Handle authentication error
//	    } else if errors.IsRateLimitError(err) {
//	        // Handle rate limit
//	    }
//	}
//
// # Context Propagation
//
// All API methods accept a context.Context for cancellation and timeout control:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	resp, err := client.Chat.Create(ctx, req)
//
// # More Information
//
// For detailed documentation and examples, see:
//   - API Reference: https://pkg.go.dev/github.com/z-ai/zai-sdk-go
//   - Official Docs: https://docs.z.ai/
//   - GitHub: https://github.com/z-ai/zai-sdk-go
package zai
