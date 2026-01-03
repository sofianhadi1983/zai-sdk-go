package zai

import (
	"context"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/assistant"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
	"github.com/sofianhadi1983/zai-sdk-go/internal/streaming"
)

// AssistantService provides access to the Assistant API.
type AssistantService struct {
	client *client.BaseClient
}

// newAssistantService creates a new assistant service.
func newAssistantService(baseClient *client.BaseClient) *AssistantService {
	return &AssistantService{
		client: baseClient,
	}
}

// Conversation creates or continues a conversation with an assistant.
//
// Example:
//
//	messages := []assistant.ConversationMessage{
//	    {
//	        Role: "user",
//	        Content: []assistant.MessageContent{
//	            assistant.MessageTextContent{Type: "text", Text: "Hello!"},
//	        },
//	    },
//	}
//	req := assistant.NewConversationRequest("asst_123", messages)
//	resp, err := client.Assistant.Conversation(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Println(resp.GetText())
func (s *AssistantService) Conversation(ctx context.Context, req *assistant.ConversationRequest) (*assistant.AssistantCompletion, error) {
	// Ensure stream is set to false for non-streaming requests
	req.Stream = false

	// Make the API request
	apiResp, err := s.client.Post(ctx, "/assistant", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp assistant.AssistantCompletion
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ConversationStream creates a streaming conversation with an assistant.
//
// Example:
//
//	messages := []assistant.ConversationMessage{
//	    {
//	        Role: "user",
//	        Content: []assistant.MessageContent{
//	            assistant.MessageTextContent{Type: "text", Text: "Tell me a story"},
//	        },
//	    },
//	}
//	req := assistant.NewConversationRequest("asst_123", messages)
//	req.SetStream(true)
//
//	stream, err := client.Assistant.ConversationStream(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//	defer stream.Close()
//
//	for stream.Next() {
//	    chunk := stream.Value()
//	    fmt.Print(chunk.GetText())
//	}
//	if err := stream.Err(); err != nil {
//	    // Handle error
//	}
func (s *AssistantService) ConversationStream(ctx context.Context, req *assistant.ConversationRequest) (*streaming.Stream[assistant.AssistantCompletion], error) {
	// Ensure stream is set to true
	req.Stream = true

	// Make the streaming request
	streamResp, err := s.client.Stream(ctx, "/assistant", req)
	if err != nil {
		return nil, err
	}

	// Create typed stream
	return client.NewTypedStream[assistant.AssistantCompletion](streamResp, ctx), nil
}

// QuerySupport retrieves information about available assistants.
//
// Example:
//
//	// Get all available assistants
//	resp, err := client.Assistant.QuerySupport(ctx, nil)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, asst := range resp.GetAssistants() {
//	    fmt.Printf("Assistant: %s - %s\n", asst.Name, asst.Description)
//	}
//
//	// Get specific assistants
//	resp, err = client.Assistant.QuerySupport(ctx, []string{"asst_123", "asst_456"})
func (s *AssistantService) QuerySupport(ctx context.Context, assistantIDs []string) (*assistant.AssistantSupportResponse, error) {
	body := map[string]interface{}{}
	if assistantIDs != nil {
		body["assistant_id_list"] = assistantIDs
	}

	// Make the API request
	apiResp, err := s.client.Post(ctx, "/assistant/list", body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp assistant.AssistantSupportResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// QueryConversationUsage retrieves conversation usage history for an assistant.
//
// Example:
//
//	// Get first page of conversation history
//	resp, err := client.Assistant.QueryConversationUsage(ctx, "asst_123", 1, 10)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, conv := range resp.GetConversations() {
//	    fmt.Printf("Conversation %s: %d total tokens\n",
//	        conv.ID, conv.Usage.TotalTokens)
//	}
//
//	if resp.HasMore() {
//	    // Fetch next page
//	    resp, err = client.Assistant.QueryConversationUsage(ctx, "asst_123", 2, 10)
//	}
func (s *AssistantService) QueryConversationUsage(ctx context.Context, assistantID string, page, pageSize int) (*assistant.ConversationUsageResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	body := map[string]interface{}{
		"assistant_id": assistantID,
		"page":         page,
		"page_size":    pageSize,
	}

	// Make the API request
	apiResp, err := s.client.Post(ctx, "/assistant/conversation/list", body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp assistant.ConversationUsageResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateConversation is a convenience method to create a new conversation.
//
// Example:
//
//	text := "Explain quantum computing in simple terms"
//	resp, err := client.Assistant.CreateConversation(ctx, "asst_123", text)
func (s *AssistantService) CreateConversation(ctx context.Context, assistantID, text string) (*assistant.AssistantCompletion, error) {
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{
					Type: "text",
					Text: text,
				},
			},
		},
	}

	req := assistant.NewConversationRequest(assistantID, messages)
	return s.Conversation(ctx, req)
}

// ContinueConversation is a convenience method to continue an existing conversation.
//
// Example:
//
//	text := "Can you elaborate on that?"
//	resp, err := client.Assistant.ContinueConversation(ctx, "asst_123", "conv_456", text)
func (s *AssistantService) ContinueConversation(ctx context.Context, assistantID, conversationID, text string) (*assistant.AssistantCompletion, error) {
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{
					Type: "text",
					Text: text,
				},
			},
		},
	}

	req := assistant.NewConversationRequest(assistantID, messages)
	req.SetConversationID(conversationID)
	return s.Conversation(ctx, req)
}
