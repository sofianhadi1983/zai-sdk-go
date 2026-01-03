package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sofianhadi1983/zai-sdk-go/api/types/assistant"
)

func TestAssistantService_Conversation(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/assistant", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

		// Parse request body
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "asst_123", reqBody["assistant_id"])
		assert.False(t, reqBody["stream"].(bool))

		// Send response
		resp := assistant.AssistantCompletion{
			ID:             "req_789",
			ConversationID: "conv_456",
			AssistantID:    "asst_123",
			Created:        1609459200,
			Status:         "completed",
			Choices: []assistant.AssistantChoice{
				{
					Index: 0,
					Delta: assistant.TextContentBlock{
						Content: "Hello! How can I help you?",
						Role:    "assistant",
						Type:    "content",
					},
					FinishReason: "stop",
				},
			},
			Usage: &assistant.CompletionUsage{
				PromptTokens:     10,
				CompletionTokens: 8,
				TotalTokens:      18,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{Type: "text", Text: "Hello"},
			},
		},
	}
	req := assistant.NewConversationRequest("asst_123", messages)

	// Make request
	resp, err := client.Assistant.Conversation(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "req_789", resp.ID)
	assert.Equal(t, "conv_456", resp.ConversationID)
	assert.Equal(t, "asst_123", resp.AssistantID)
	assert.True(t, resp.IsCompleted())
	assert.Equal(t, "Hello! How can I help you?", resp.GetText())
	require.NotNil(t, resp.Usage)
	assert.Equal(t, 18, resp.Usage.TotalTokens)
}

func TestAssistantService_ConversationStream(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/assistant", r.URL.Path)

		// Parse request body
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify stream is enabled
		assert.True(t, reqBody["stream"].(bool))

		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send multiple chunks
		chunks := []assistant.AssistantCompletion{
			{
				ID:             "req_789",
				ConversationID: "conv_456",
				AssistantID:    "asst_123",
				Status:         "in_progress",
				Choices: []assistant.AssistantChoice{
					{
						Index: 0,
						Delta: assistant.TextContentBlock{
							Content: "Hello",
							Role:    "assistant",
							Type:    "content",
						},
						FinishReason: "",
					},
				},
			},
			{
				ID:             "req_789",
				ConversationID: "conv_456",
				AssistantID:    "asst_123",
				Status:         "in_progress",
				Choices: []assistant.AssistantChoice{
					{
						Index: 0,
						Delta: assistant.TextContentBlock{
							Content: "!",
							Role:    "assistant",
							Type:    "content",
						},
						FinishReason: "",
					},
				},
			},
			{
				ID:             "req_789",
				ConversationID: "conv_456",
				AssistantID:    "asst_123",
				Status:         "completed",
				Choices: []assistant.AssistantChoice{
					{
						Index:        0,
						Delta:        assistant.TextContentBlock{Type: "content"},
						FinishReason: "stop",
					},
				},
			},
		}

		for _, chunk := range chunks {
			data, _ := json.Marshal(chunk)
			w.Write([]byte("data: "))
			w.Write(data)
			w.Write([]byte("\n\n"))
		}

		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	messages := []assistant.ConversationMessage{
		{
			Role: "user",
			Content: []assistant.MessageContent{
				assistant.MessageTextContent{Type: "text", Text: "Hi"},
			},
		},
	}
	req := assistant.NewConversationRequest("asst_123", messages)

	// Make streaming request
	stream, err := client.Assistant.ConversationStream(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, stream)
	defer stream.Close()

	// Collect chunks
	var chunks []assistant.AssistantCompletion
	for stream.Next() {
		chunk := stream.Current()
		require.NotNil(t, chunk)
		chunks = append(chunks, *chunk)
	}
	require.NoError(t, stream.Err())

	// Verify chunks
	require.Len(t, chunks, 3)
	assert.Equal(t, "Hello", chunks[0].GetText())
	assert.Equal(t, "!", chunks[1].GetText())
	assert.True(t, chunks[2].IsCompleted())
}

func TestAssistantService_QuerySupport(t *testing.T) {
	t.Parallel()

	t.Run("query all assistants", func(t *testing.T) {
		t.Parallel()

		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/assistant/list", r.URL.Path)

			// Send response
			resp := assistant.AssistantSupportResponse{
				Code:    200,
				Message: "success",
				Data: []assistant.AssistantSupport{
					{
						AssistantID:    "asst_1",
						Name:           "Code Helper",
						Description:    "Helps with coding tasks",
						Status:         "publish",
						Tools:          []string{"code_interpreter"},
						StarterPrompts: []string{"Help me debug this"},
					},
					{
						AssistantID:    "asst_2",
						Name:           "Math Tutor",
						Description:    "Assists with math problems",
						Status:         "publish",
						Tools:          []string{"calculator"},
						StarterPrompts: []string{"Solve this equation"},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		// Create client
		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		// Query all assistants
		resp, err := client.Assistant.QuerySupport(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify response
		assistants := resp.GetAssistants()
		require.Len(t, assistants, 2)
		assert.Equal(t, "asst_1", assistants[0].AssistantID)
		assert.Equal(t, "Code Helper", assistants[0].Name)
		assert.Equal(t, "asst_2", assistants[1].AssistantID)
		assert.Equal(t, "Math Tutor", assistants[1].Name)
	})

	t.Run("query specific assistants", func(t *testing.T) {
		t.Parallel()

		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse request body
			var reqBody map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)

			// Verify assistant IDs
			ids := reqBody["assistant_id_list"].([]interface{})
			assert.Len(t, ids, 2)

			// Send response
			resp := assistant.AssistantSupportResponse{
				Code:    200,
				Message: "success",
				Data: []assistant.AssistantSupport{
					{
						AssistantID: "asst_1",
						Name:        "Code Helper",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		// Create client
		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		// Query specific assistants
		resp, err := client.Assistant.QuerySupport(context.Background(), []string{"asst_1", "asst_2"})
		require.NoError(t, err)
		require.NotNil(t, resp)

		assistants := resp.GetAssistants()
		require.Len(t, assistants, 1)
		assert.Equal(t, "asst_1", assistants[0].AssistantID)
	})
}

func TestAssistantService_QueryConversationUsage(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/assistant/conversation/list", r.URL.Path)

		// Parse request body
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "asst_123", reqBody["assistant_id"])
		assert.Equal(t, float64(1), reqBody["page"])
		assert.Equal(t, float64(10), reqBody["page_size"])

		// Send response
		resp := assistant.ConversationUsageResponse{
			Code:    200,
			Message: "success",
			Data: assistant.ConversationUsageList{
				AssistantID: "asst_123",
				HasMore:     true,
				ConversationList: []assistant.ConversationUsage{
					{
						ID:          "conv_1",
						AssistantID: "asst_123",
						CreateTime:  1609459200,
						UpdateTime:  1609545600,
						Usage: assistant.CompletionUsage{
							PromptTokens:     100,
							CompletionTokens: 50,
							TotalTokens:      150,
						},
					},
					{
						ID:          "conv_2",
						AssistantID: "asst_123",
						CreateTime:  1609632000,
						UpdateTime:  1609718400,
						Usage: assistant.CompletionUsage{
							PromptTokens:     200,
							CompletionTokens: 100,
							TotalTokens:      300,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Query conversation usage
	resp, err := client.Assistant.QueryConversationUsage(context.Background(), "asst_123", 1, 10)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.True(t, resp.HasMore())
	conversations := resp.GetConversations()
	require.Len(t, conversations, 2)

	conv1 := conversations[0]
	assert.Equal(t, "conv_1", conv1.ID)
	assert.Equal(t, "asst_123", conv1.AssistantID)
	assert.Equal(t, 150, conv1.Usage.TotalTokens)

	conv2 := conversations[1]
	assert.Equal(t, "conv_2", conv2.ID)
	assert.Equal(t, 300, conv2.Usage.TotalTokens)
}

func TestAssistantService_CreateConversation(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send response
		resp := assistant.AssistantCompletion{
			ID:             "req_789",
			ConversationID: "conv_new",
			AssistantID:    "asst_123",
			Status:         "completed",
			Choices: []assistant.AssistantChoice{
				{
					Index: 0,
					Delta: assistant.TextContentBlock{
						Content: "Here's an explanation of quantum computing...",
						Role:    "assistant",
						Type:    "content",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create conversation
	resp, err := client.Assistant.CreateConversation(
		context.Background(),
		"asst_123",
		"Explain quantum computing",
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "conv_new", resp.ConversationID)
	assert.True(t, resp.IsCompleted())
	assert.Contains(t, resp.GetText(), "quantum computing")
}

func TestAssistantService_ContinueConversation(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify conversation ID is set
		assert.Equal(t, "conv_existing", reqBody["conversation_id"])

		// Send response
		resp := assistant.AssistantCompletion{
			ID:             "req_790",
			ConversationID: "conv_existing",
			AssistantID:    "asst_123",
			Status:         "completed",
			Choices: []assistant.AssistantChoice{
				{
					Index: 0,
					Delta: assistant.TextContentBlock{
						Content: "Here's more detail...",
						Role:    "assistant",
						Type:    "content",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Continue conversation
	resp, err := client.Assistant.ContinueConversation(
		context.Background(),
		"asst_123",
		"conv_existing",
		"Can you elaborate?",
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "conv_existing", resp.ConversationID)
	assert.True(t, resp.IsCompleted())
}

func TestAssistantService_APIError(t *testing.T) {
	t.Parallel()

	// Mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "Invalid assistant ID", "code": "invalid_request"}}`))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Try to create conversation
	resp, err := client.Assistant.CreateConversation(
		context.Background(),
		"invalid_id",
		"Hello",
	)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestNewAssistantService(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithAPIKey("test-key.test-secret"))
	require.NoError(t, err)
	defer client.Close()

	// Verify Assistant service is initialized
	assert.NotNil(t, client.Assistant)
}
