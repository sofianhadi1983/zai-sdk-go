package zai

import (
	"context"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/agents"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
	streaming "github.com/sofianhadi1983/zai-sdk-go/internal/streaming"
)

// AgentsService provides access to the Agents API.
type AgentsService struct {
	client *client.BaseClient
}

// newAgentsService creates a new agents service.
func newAgentsService(baseClient *client.BaseClient) *AgentsService {
	return &AgentsService{
		client: baseClient,
	}
}

// Invoke invokes an agent with the given request.
//
// Example:
//
//	messages := []chat.Message{
//	    {Role: chat.RoleUser, Content: "Tell me a joke"},
//	}
//
//	req := agents.NewAgentInvokeRequest("general_translation", messages).
//	    SetUserID("user_123")
//
//	resp, err := client.Agents.Invoke(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	content := resp.GetContent()
//	fmt.Println(content)
func (s *AgentsService) Invoke(ctx context.Context, req *agents.AgentInvokeRequest) (*agents.AgentCompletionResponse, error) {
	// Ensure streaming is disabled
	req.Stream = false

	// Make the API request
	apiResp, err := s.client.Post(ctx, "/v1/agents", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp agents.AgentCompletionResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// InvokeStream invokes an agent with streaming response.
//
// Example:
//
//	messages := []chat.Message{
//	    {Role: chat.RoleUser, Content: "Write a story"},
//	}
//
//	req := agents.NewAgentInvokeRequest("general_translation", messages).
//	    SetStream(true)
//
//	stream, err := client.Agents.InvokeStream(ctx, req)
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
//	    // Handle error
//	}
func (s *AgentsService) InvokeStream(ctx context.Context, req *agents.AgentInvokeRequest) (*streaming.Stream[agents.AgentCompletionChunk], error) {
	// Ensure streaming is enabled
	req.Stream = true

	// Make the streaming request
	streamResp, err := s.client.Stream(ctx, "/v1/agents", req)
	if err != nil {
		return nil, err
	}

	// Create typed stream
	return client.NewTypedStream[agents.AgentCompletionChunk](streamResp, ctx), nil
}

// AsyncResult retrieves the result of an async agent invocation.
//
// Example:
//
//	req := agents.NewAgentAsyncResultRequest("agent_123").
//	    SetAsyncID("async_456").
//	    SetConversationID("conv_789")
//
//	resp, err := client.Agents.AsyncResult(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.Status == "completed" {
//	    content := resp.GetContent()
//	    fmt.Println(content)
//	}
func (s *AgentsService) AsyncResult(ctx context.Context, req *agents.AgentAsyncResultRequest) (*agents.AgentCompletionResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/v1/agents/async-result", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp agents.AgentCompletionResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
