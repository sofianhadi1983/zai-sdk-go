package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/sofianhadi1983/zai-sdk-go/api/types/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolsService_WebSearch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/tools", r.URL.Path)

		// Parse request body
		var reqBody tools.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "web-search-pro", reqBody.Model)
		assert.False(t, reqBody.Stream)
		assert.Len(t, reqBody.Messages, 1)
		assert.Equal(t, "What is AI?", reqBody.Messages[0].Content)

		// Send mock response
		resp := tools.WebSearchResponse{
			ID:        "ws_abc123",
			Created:   1700000000,
			RequestID: "req_xyz789",
			Choices: []tools.WebSearchChoice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: tools.WebSearchMessage{
						Role: "assistant",
						ToolCalls: []tools.WebSearchMessageToolCall{
							{
								ID:   "call_1",
								Type: "web_search",
								SearchIntent: &tools.SearchIntent{
									Index:    0,
									Query:    "artificial intelligence definition",
									Intent:   "informational",
									Keywords: "AI, definition, technology",
								},
							},
							{
								ID:   "call_2",
								Type: "web_search",
								SearchResult: &tools.SearchResult{
									Index:   0,
									Title:   "What is Artificial Intelligence?",
									Link:    "https://example.com/ai-intro",
									Content: "Artificial intelligence is the simulation of human intelligence...",
									Icon:    "https://example.com/favicon.ico",
									Media:   "Tech Encyclopedia",
									Refer:   "[ref_1]",
								},
							},
							{
								ID:   "call_3",
								Type: "web_search",
								SearchRecommend: &tools.SearchRecommend{
									Index: 0,
									Query: "AI applications and examples",
								},
							},
						},
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

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What is AI?",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages)

	resp, err := client.Tools.WebSearch(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "ws_abc123", resp.ID)
	assert.Equal(t, int64(1700000000), resp.Created)
	assert.Equal(t, "req_xyz789", resp.RequestID)

	// Verify choices
	choices := resp.GetChoices()
	assert.Len(t, choices, 1)
	assert.Equal(t, "stop", choices[0].FinishReason)

	// Verify tool calls
	toolCalls := resp.GetToolCalls()
	assert.Len(t, toolCalls, 3)

	// Verify search intents
	intents := resp.GetSearchIntents()
	assert.Len(t, intents, 1)
	assert.Equal(t, "artificial intelligence definition", intents[0].Query)
	assert.Equal(t, "informational", intents[0].Intent)

	// Verify search results
	results := resp.GetSearchResults()
	assert.Len(t, results, 1)
	assert.Equal(t, "What is Artificial Intelligence?", results[0].Title)
	assert.Equal(t, "https://example.com/ai-intro", results[0].Link)

	// Verify search recommendations
	recommends := resp.GetSearchRecommendations()
	assert.Len(t, recommends, 1)
	assert.Equal(t, "AI applications and examples", recommends[0].Query)
}

func TestToolsService_WebSearch_WithFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/tools", r.URL.Path)

		// Parse request body
		var reqBody tools.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify filters
		assert.Equal(t, "academic", reqBody.Scope)
		assert.Equal(t, "US", reqBody.Location)
		assert.Equal(t, 7, reqBody.RecentDays)
		assert.Equal(t, "req_filter123", reqBody.RequestID)

		// Send mock response
		resp := tools.WebSearchResponse{
			ID: "ws_filtered",
			Choices: []tools.WebSearchChoice{
				{
					Index: 0,
					Message: tools.WebSearchMessage{
						Role: "assistant",
						ToolCalls: []tools.WebSearchMessageToolCall{
							{
								ID:   "call_1",
								Type: "web_search",
								SearchResult: &tools.SearchResult{
									Index:   0,
									Title:   "Recent Academic Paper on AI",
									Link:    "https://arxiv.org/paper123",
									Content: "This paper discusses...",
									Media:   "arXiv",
									Refer:   "[ref_1]",
								},
							},
						},
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

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Latest AI research",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages).
		SetScope("academic").
		SetLocation("US").
		SetRecentDays(7).
		SetRequestID("req_filter123")

	resp, err := client.Tools.WebSearch(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "ws_filtered", resp.ID)
	results := resp.GetSearchResults()
	assert.Len(t, results, 1)
	assert.Equal(t, "Recent Academic Paper on AI", results[0].Title)
}

func TestToolsService_WebSearchStream(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/tools", r.URL.Path)

		// Parse request body
		var reqBody tools.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify streaming is enabled
		assert.True(t, reqBody.Stream)

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		require.True(t, ok)

		// Send streaming chunks
		chunks := []tools.WebSearchChunk{
			{
				ID:      "ws_chunk_1",
				Created: 1700000000,
				Choices: []tools.WebSearchStreamChoice{
					{
						Index: 0,
						Delta: tools.ChoiceDelta{
							Role: "assistant",
							ToolCalls: []tools.ChoiceDeltaToolCall{
								{
									Index: 0,
									ID:    "call_1",
									Type:  "web_search",
									SearchIntent: &tools.SearchIntent{
										Index:    0,
										Query:    "streaming search query",
										Intent:   "informational",
										Keywords: "stream, test",
									},
								},
							},
						},
					},
				},
			},
			{
				ID:      "ws_chunk_2",
				Created: 1700000001,
				Choices: []tools.WebSearchStreamChoice{
					{
						Index: 0,
						Delta: tools.ChoiceDelta{
							ToolCalls: []tools.ChoiceDeltaToolCall{
								{
									Index: 0,
									SearchResult: &tools.SearchResult{
										Index:   0,
										Title:   "Streaming Result",
										Link:    "https://stream.example.com",
										Content: "Streaming content...",
										Media:   "Stream Media",
										Refer:   "[ref_1]",
									},
								},
							},
						},
					},
				},
			},
			{
				ID:      "ws_chunk_3",
				Created: 1700000002,
				Choices: []tools.WebSearchStreamChoice{
					{
						Index:        0,
						FinishReason: "stop",
						Delta: tools.ChoiceDelta{
							ToolCalls: []tools.ChoiceDeltaToolCall{
								{
									Index: 0,
									SearchRecommend: &tools.SearchRecommend{
										Index: 0,
										Query: "try this query",
									},
								},
							},
						},
					},
				},
			},
		}

		for _, chunk := range chunks {
			data, _ := json.Marshal(chunk)
			_, _ = w.Write([]byte("data: "))
			_, _ = w.Write(data)
			_, _ = w.Write([]byte("\n\n"))
			flusher.Flush()
		}

		// Send done event
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Stream search test",
		},
	}

	req := tools.NewWebSearchRequest("web-search-pro", messages).
		SetStream(true)

	stream, err := client.Tools.WebSearchStream(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, stream)
	defer stream.Close()

	// Collect chunks
	var chunks []tools.WebSearchChunk
	for stream.Next() {
		chunk := stream.Current()
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil && !strings.Contains(err.Error(), "[DONE]") {
		require.NoError(t, err)
	}

	// Verify we received all chunks
	require.Len(t, chunks, 3)

	// Verify first chunk (search intent)
	assert.Equal(t, "ws_chunk_1", chunks[0].ID)
	assert.Len(t, chunks[0].Choices, 1)
	assert.Equal(t, "assistant", chunks[0].Choices[0].Delta.Role)
	assert.NotNil(t, chunks[0].Choices[0].Delta.ToolCalls[0].SearchIntent)
	assert.Equal(t, "streaming search query", chunks[0].Choices[0].Delta.ToolCalls[0].SearchIntent.Query)

	// Verify second chunk (search result)
	assert.Equal(t, "ws_chunk_2", chunks[1].ID)
	assert.NotNil(t, chunks[1].Choices[0].Delta.ToolCalls[0].SearchResult)
	assert.Equal(t, "Streaming Result", chunks[1].Choices[0].Delta.ToolCalls[0].SearchResult.Title)

	// Verify third chunk (search recommendation and finish)
	assert.Equal(t, "ws_chunk_3", chunks[2].ID)
	assert.Equal(t, "stop", chunks[2].Choices[0].FinishReason)
	assert.NotNil(t, chunks[2].Choices[0].Delta.ToolCalls[0].SearchRecommend)
	assert.Equal(t, "try this query", chunks[2].Choices[0].Delta.ToolCalls[0].SearchRecommend.Query)
}

func TestToolsService_WebSearch_APIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Invalid web search request",
				"type":    "invalid_request_error",
				"code":    "invalid_search_request",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := tools.NewWebSearchRequest("web-search-pro", []chat.Message{})

	_, err = client.Tools.WebSearch(context.Background(), req)
	require.Error(t, err)
}
