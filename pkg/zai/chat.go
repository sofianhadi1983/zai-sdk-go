package zai

import (
	"context"

	"github.com/z-ai/zai-sdk-go/api/types/chat"
	"github.com/z-ai/zai-sdk-go/internal/client"
	"github.com/z-ai/zai-sdk-go/internal/streaming"
)

// ChatService provides access to the Chat Completions API.
type ChatService struct {
	client *client.BaseClient
}

// newChatService creates a new chat service.
func newChatService(baseClient *client.BaseClient) *ChatService {
	return &ChatService{
		client: baseClient,
	}
}

// Create creates a chat completion.
//
// Example:
//
//	req := &chat.ChatCompletionRequest{
//	    Model: "glm-4.7",
//	    Messages: []chat.Message{
//	        chat.NewUserMessage("Hello!"),
//	    },
//	}
//
//	resp, err := client.Chat.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Println(resp.GetContent())
func (s *ChatService) Create(ctx context.Context, req *chat.ChatCompletionRequest) (*chat.ChatCompletionResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/chat/completions", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp chat.ChatCompletionResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateStream creates a streaming chat completion.
// Returns a stream of chat completion chunks.
//
// Example:
//
//	req := &chat.ChatCompletionRequest{
//	    Model: "glm-4.7",
//	    Messages: []chat.Message{
//	        chat.NewUserMessage("Tell me a story"),
//	    },
//	}
//	req.SetStream(true)
//
//	stream, err := client.Chat.CreateStream(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//	defer stream.Close()
//
//	for stream.Next() {
//	    chunk := stream.Current()
//	    fmt.Print(chunk.GetContent())
//	}
//
//	if err := stream.Err(); err != nil {
//	    // Handle stream error
//	}
func (s *ChatService) CreateStream(ctx context.Context, req *chat.ChatCompletionRequest) (*streaming.Stream[chat.ChatCompletionChunk], error) {
	// Ensure stream is enabled
	stream := true
	req.Stream = &stream

	// Make the streaming request
	streamResp, err := s.client.Stream(ctx, "/chat/completions", req)
	if err != nil {
		return nil, err
	}

	// Create typed stream
	return client.NewTypedStream[chat.ChatCompletionChunk](streamResp, ctx), nil
}

// StreamContent is a convenience method that streams content and collects it into a string.
// Returns the complete content and any error that occurred.
//
// Example:
//
//	req := &chat.ChatCompletionRequest{
//	    Model: "glm-4.7",
//	    Messages: []chat.Message{
//	        chat.NewUserMessage("Hello!"),
//	    },
//	}
//
//	content, err := client.Chat.StreamContent(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Println(content)
func (s *ChatService) StreamContent(ctx context.Context, req *chat.ChatCompletionRequest) (string, error) {
	stream, err := s.CreateStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var content string
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			content += chunk.GetContent()
		}
	}

	if err := stream.Err(); err != nil {
		return content, err
	}

	return content, nil
}
