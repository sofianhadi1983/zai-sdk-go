package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/agents"
	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentsService_Invoke(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/agents", r.URL.Path)

		resp := agents.AgentCompletionResponse{
			ID:             "comp_123",
			AgentID:        "agent_translation",
			ConversationID: "conv_456",
			Status:         "completed",
			Choices: []agents.AgentCompletionChoice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: agents.AgentCompletionMessage{
						Role:    "assistant",
						Content: "Agent response content",
					},
				},
			},
			Usage: &agents.AgentCompletionUsage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	messages := []chat.Message{
		{Role: chat.RoleUser, Content: "test message"},
	}

	req := agents.NewAgentInvokeRequest("agent_translation", messages)

	resp, err := client.Agents.Invoke(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "comp_123", resp.ID)
	assert.Equal(t, "agent_translation", resp.AgentID)
	assert.Equal(t, "completed", resp.Status)
	assert.Len(t, resp.GetChoices(), 1)
	assert.Equal(t, "Agent response content", resp.GetContent())
	assert.False(t, resp.HasError())
}

func TestAgentsService_InvokeStream(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/agents", r.URL.Path)

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		chunks := []agents.AgentCompletionChunk{
			{
				ID:      "chunk_1",
				AgentID: "agent_test",
				Choices: []agents.AgentStreamChoice{
					{
						Index: 0,
						Delta: agents.AgentChoiceDelta{
							Role:    "assistant",
							Content: "Hello",
						},
					},
				},
			},
			{
				ID: "chunk_2",
				Choices: []agents.AgentStreamChoice{
					{
						Index: 0,
						Delta: agents.AgentChoiceDelta{
							Content: " world",
						},
					},
				},
			},
		}

		for _, chunk := range chunks {
			data, _ := json.Marshal(chunk)
			w.Write([]byte("data: "))
			w.Write(data)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		}

		w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	messages := []chat.Message{
		{Role: chat.RoleUser, Content: "test"},
	}

	req := agents.NewAgentInvokeRequest("agent_test", messages).SetStream(true)

	stream, err := client.Agents.InvokeStream(context.Background(), req)
	require.NoError(t, err)
	defer stream.Close()

	var chunks []agents.AgentCompletionChunk
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
	}

	if err := stream.Err(); err != nil && !strings.Contains(err.Error(), "[DONE]") {
		require.NoError(t, err)
	}

	require.Len(t, chunks, 2)
	assert.Equal(t, "Hello", chunks[0].GetContent())
	assert.Equal(t, " world", chunks[1].GetContent())
}

func TestAgentsService_AsyncResult(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/agents/async-result", r.URL.Path)

		resp := agents.AgentCompletionResponse{
			ID:             "async_result_123",
			AgentID:        "agent_xyz",
			ConversationID: "conv_789",
			Status:         "completed",
			Choices: []agents.AgentCompletionChoice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: agents.AgentCompletionMessage{
						Role:    "assistant",
						Content: "Async result content",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := agents.NewAgentAsyncResultRequest("agent_xyz").
		SetAsyncID("async_456").
		SetConversationID("conv_789")

	resp, err := client.Agents.AsyncResult(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "async_result_123", resp.ID)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, "Async result content", resp.GetContent())
}
