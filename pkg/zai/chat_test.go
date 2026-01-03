package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/z-ai/zai-sdk-go/api/types/chat"
)

func TestChatService_Create(t *testing.T) {
	t.Parallel()

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()

		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/chat/completions", r.URL.Path)
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Verify request body
			var req chat.ChatCompletionRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)
			assert.Equal(t, "glm-4.7", req.Model)
			assert.Len(t, req.Messages, 1)

			// Send response
			resp := chat.ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: 1700000000,
				Model:   "glm-4.7",
				Choices: []chat.Choice{
					{
						Index: 0,
						Message: chat.Message{
							Role:    chat.RoleAssistant,
							Content: "Hello! How can I help you?",
						},
						FinishReason: "stop",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
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

		// Make request
		req := &chat.ChatCompletionRequest{
			Model: "glm-4.7",
			Messages: []chat.Message{
				chat.NewUserMessage("Hello"),
			},
		}

		resp, err := client.Chat.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, "chatcmpl-123", resp.ID)
		assert.Equal(t, "glm-4.7", resp.Model)
		assert.Len(t, resp.Choices, 1)
		assert.Equal(t, "Hello! How can I help you?", resp.GetContent())
	})

	t.Run("with temperature and max_tokens", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req chat.ChatCompletionRequest
			json.NewDecoder(r.Body).Decode(&req)

			assert.NotNil(t, req.Temperature)
			assert.Equal(t, 0.7, *req.Temperature)
			assert.NotNil(t, req.MaxTokens)
			assert.Equal(t, 100, *req.MaxTokens)

			resp := chat.ChatCompletionResponse{
				ID:      "test",
				Object:  "chat.completion",
				Created: 1700000000,
				Model:   "glm-4",
				Choices: []chat.Choice{
					{
						Index: 0,
						Message: chat.Message{
							Role:    chat.RoleAssistant,
							Content: "Response",
						},
						FinishReason: "stop",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Test"),
			},
		}
		req.SetTemperature(0.7).SetMaxTokens(100)

		resp, err := client.Chat.Create(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "Response", resp.GetContent())
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid request",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Test"),
			},
		}

		resp, err := client.Chat.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Invalid request")
	})
}

func TestChatService_CreateStream(t *testing.T) {
	t.Parallel()

	t.Run("successful stream", func(t *testing.T) {
		t.Parallel()

		// Create mock SSE server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/chat/completions", r.URL.Path)

			// Verify stream parameter is set
			var req chat.ChatCompletionRequest
			json.NewDecoder(r.Body).Decode(&req)
			assert.NotNil(t, req.Stream)
			assert.True(t, *req.Stream)

			// Send SSE chunks
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			chunks := []chat.ChatCompletionChunk{
				{
					ID:      "chatcmpl-123",
					Object:  "chat.completion.chunk",
					Created: 1700000000,
					Model:   "glm-4",
					Choices: []chat.ChunkChoice{
						{
							Index: 0,
							Delta: chat.Delta{
								Role: chat.RoleAssistant,
							},
						},
					},
				},
				{
					ID:      "chatcmpl-123",
					Object:  "chat.completion.chunk",
					Created: 1700000000,
					Model:   "glm-4",
					Choices: []chat.ChunkChoice{
						{
							Index: 0,
							Delta: chat.Delta{
								Content: "Hello",
							},
						},
					},
				},
				{
					ID:      "chatcmpl-123",
					Object:  "chat.completion.chunk",
					Created: 1700000000,
					Model:   "glm-4",
					Choices: []chat.ChunkChoice{
						{
							Index: 0,
							Delta: chat.Delta{
								Content: " world",
							},
						},
					},
				},
				{
					ID:      "chatcmpl-123",
					Object:  "chat.completion.chunk",
					Created: 1700000000,
					Model:   "glm-4",
					Choices: []chat.ChunkChoice{
						{
							Index:        0,
							Delta:        chat.Delta{},
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

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Hello"),
			},
		}

		stream, err := client.Chat.CreateStream(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, stream)
		defer stream.Close()

		// Collect chunks
		var content string
		var chunkCount int
		for stream.Next() {
			chunk := stream.Current()
			content += chunk.GetContent()
			chunkCount++
		}

		assert.NoError(t, stream.Err())
		assert.Equal(t, "Hello world", content)
		assert.Greater(t, chunkCount, 0)
	})

	t.Run("stream error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid API key",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Test"),
			},
		}

		stream, err := client.Chat.CreateStream(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, stream)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}

func TestChatService_StreamContent(t *testing.T) {
	t.Parallel()

	t.Run("successful stream content", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			// Send a few chunks
			contents := []string{"The", " quick", " brown", " fox"}
			for _, content := range contents {
				chunk := chat.ChatCompletionChunk{
					ID:      "test",
					Object:  "chat.completion.chunk",
					Created: 1700000000,
					Model:   "glm-4",
					Choices: []chat.ChunkChoice{
						{
							Index: 0,
							Delta: chat.Delta{
								Content: content,
							},
						},
					},
				}
				data, _ := json.Marshal(chunk)
				w.Write([]byte("data: "))
				w.Write(data)
				w.Write([]byte("\n\n"))
			}

			// Final chunk
			finalChunk := chat.ChatCompletionChunk{
				ID:      "test",
				Object:  "chat.completion.chunk",
				Created: 1700000000,
				Model:   "glm-4",
				Choices: []chat.ChunkChoice{
					{
						Index:        0,
						Delta:        chat.Delta{},
						FinishReason: "stop",
					},
				},
			}
			data, _ := json.Marshal(finalChunk)
			w.Write([]byte("data: "))
			w.Write(data)
			w.Write([]byte("\n\n"))

			w.Write([]byte("data: [DONE]\n\n"))
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Tell me something"),
			},
		}

		content, err := client.Chat.StreamContent(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "The quick brown fox", content)
	})

	t.Run("stream error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Internal error",
				},
			})
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		req := &chat.ChatCompletionRequest{
			Model: "glm-4",
			Messages: []chat.Message{
				chat.NewUserMessage("Test"),
			},
		}

		content, err := client.Chat.StreamContent(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, content)
	})
}

func TestClient_ChatService_Integration(t *testing.T) {
	t.Parallel()

	t.Run("client has chat service", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.Chat)
	})

	t.Run("complete conversation example", func(t *testing.T) {
		t.Parallel()

		// Track conversation
		var conversationLog []string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req chat.ChatCompletionRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Log the conversation
			for _, msg := range req.Messages {
				if content, ok := msg.Content.(string); ok {
					conversationLog = append(conversationLog, string(msg.Role)+": "+content)
				}
			}

			// Generate response based on last message
			lastMsg := req.Messages[len(req.Messages)-1]
			var responseContent string
			if content, ok := lastMsg.Content.(string); ok {
				if strings.Contains(strings.ToLower(content), "hello") {
					responseContent = "Hi there! How can I help you?"
				} else {
					responseContent = "I'm here to assist you."
				}
			}

			resp := chat.ChatCompletionResponse{
				ID:      "test",
				Object:  "chat.completion",
				Created: 1700000000,
				Model:   req.Model,
				Choices: []chat.Choice{
					{
						Index: 0,
						Message: chat.Message{
							Role:    chat.RoleAssistant,
							Content: responseContent,
						},
						FinishReason: "stop",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		// First message
		req1 := &chat.ChatCompletionRequest{
			Model: "glm-4",
		}
		req1.AddSystemMessage("You are a helpful assistant").
			AddUserMessage("Hello!")

		resp1, err := client.Chat.Create(context.Background(), req1)
		require.NoError(t, err)
		assert.Contains(t, resp1.GetContent(), "Hi there")

		// Follow-up message
		req2 := &chat.ChatCompletionRequest{
			Model: "glm-4",
		}
		req2.AddSystemMessage("You are a helpful assistant").
			AddUserMessage("Hello!").
			AddAssistantMessage(resp1.GetContent()).
			AddUserMessage("Thanks!")

		resp2, err := client.Chat.Create(context.Background(), req2)
		require.NoError(t, err)
		assert.NotEmpty(t, resp2.GetContent())

		// Verify conversation was tracked
		assert.Contains(t, conversationLog, "system: You are a helpful assistant")
		assert.Contains(t, conversationLog, "user: Hello!")
		assert.Contains(t, conversationLog, "user: Thanks!")
	})
}
