package tools

import (
	"encoding/json"
	"testing"

	"github.com/z-ai/zai-sdk-go/api/types/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebSearchRequest(t *testing.T) {
	t.Parallel()

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What is artificial intelligence?",
		},
	}

	req := NewWebSearchRequest("web-search-pro", messages)

	assert.Equal(t, "web-search-pro", req.Model)
	assert.Len(t, req.Messages, 1)
	assert.Equal(t, chat.RoleUser, req.Messages[0].Role)
}

func TestWebSearchRequest_BuilderPattern(t *testing.T) {
	t.Parallel()

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Latest news about AI",
		},
	}

	req := NewWebSearchRequest("web-search-pro", messages).
		SetStream(true).
		SetRequestID("req_123").
		SetScope("academic").
		SetLocation("US").
		SetRecentDays(7)

	assert.Equal(t, "web-search-pro", req.Model)
	assert.True(t, req.Stream)
	assert.Equal(t, "req_123", req.RequestID)
	assert.Equal(t, "academic", req.Scope)
	assert.Equal(t, "US", req.Location)
	assert.Equal(t, 7, req.RecentDays)
}

func TestWebSearchRequest_JSONMarshaling(t *testing.T) {
	t.Parallel()

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "test query",
		},
	}

	req := NewWebSearchRequest("web-search-pro", messages).
		SetScope("web").
		SetRecentDays(30)

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded WebSearchRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Model, decoded.Model)
	assert.Equal(t, req.Scope, decoded.Scope)
	assert.Equal(t, req.RecentDays, decoded.RecentDays)
}

func TestWebSearchResponse(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal from API response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "ws_abc123",
			"created": 1700000000,
			"request_id": "req_xyz789",
			"choices": [
				{
					"index": 0,
					"finish_reason": "stop",
					"message": {
						"role": "assistant",
						"tool_calls": [
							{
								"id": "call_1",
								"type": "web_search",
								"search_intent": {
									"index": 0,
									"query": "artificial intelligence 2024",
									"intent": "informational",
									"keywords": "AI, 2024, technology"
								}
							},
							{
								"id": "call_2",
								"type": "web_search",
								"search_result": {
									"index": 0,
									"title": "AI Trends in 2024",
									"link": "https://example.com/ai-2024",
									"content": "Latest developments in AI...",
									"icon": "https://example.com/icon.png",
									"media": "Tech News",
									"refer": "[ref_1]"
								}
							},
							{
								"id": "call_3",
								"type": "web_search",
								"search_recommend": {
									"index": 0,
									"query": "AI breakthroughs 2024"
								}
							}
						]
					}
				}
			]
		}`

		var resp WebSearchResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "ws_abc123", resp.ID)
		assert.Equal(t, int64(1700000000), resp.Created)
		assert.Equal(t, "req_xyz789", resp.RequestID)

		// Check choices
		choices := resp.GetChoices()
		assert.Len(t, choices, 1)
		assert.Equal(t, 0, choices[0].Index)
		assert.Equal(t, "stop", choices[0].FinishReason)
		assert.Equal(t, "assistant", choices[0].Message.Role)

		// Check tool calls
		toolCalls := resp.GetToolCalls()
		assert.Len(t, toolCalls, 3)

		// Verify search intent
		assert.NotNil(t, toolCalls[0].SearchIntent)
		assert.Equal(t, "artificial intelligence 2024", toolCalls[0].SearchIntent.Query)
		assert.Equal(t, "informational", toolCalls[0].SearchIntent.Intent)

		// Verify search result
		assert.NotNil(t, toolCalls[1].SearchResult)
		assert.Equal(t, "AI Trends in 2024", toolCalls[1].SearchResult.Title)
		assert.Equal(t, "https://example.com/ai-2024", toolCalls[1].SearchResult.Link)

		// Verify search recommendation
		assert.NotNil(t, toolCalls[2].SearchRecommend)
		assert.Equal(t, "AI breakthroughs 2024", toolCalls[2].SearchRecommend.Query)
	})

	t.Run("empty choices", func(t *testing.T) {
		t.Parallel()

		resp := &WebSearchResponse{}
		assert.Empty(t, resp.GetChoices())
		assert.Empty(t, resp.GetToolCalls())
	})
}

func TestWebSearchResponse_HelperMethods(t *testing.T) {
	t.Parallel()

	resp := &WebSearchResponse{
		Choices: []WebSearchChoice{
			{
				Index: 0,
				Message: WebSearchMessage{
					Role: "assistant",
					ToolCalls: []WebSearchMessageToolCall{
						{
							ID:   "call_1",
							Type: "web_search",
							SearchIntent: &SearchIntent{
								Index:    0,
								Query:    "test query",
								Intent:   "informational",
								Keywords: "test, keywords",
							},
						},
						{
							ID:   "call_2",
							Type: "web_search",
							SearchResult: &SearchResult{
								Index:   0,
								Title:   "Test Result",
								Link:    "https://example.com",
								Content: "Test content",
								Media:   "Test Media",
								Refer:   "[ref_1]",
							},
						},
						{
							ID:   "call_3",
							Type: "web_search",
							SearchRecommend: &SearchRecommend{
								Index: 0,
								Query: "recommended query",
							},
						},
					},
				},
			},
		},
	}

	t.Run("GetSearchIntents", func(t *testing.T) {
		t.Parallel()

		intents := resp.GetSearchIntents()
		assert.Len(t, intents, 1)
		assert.Equal(t, "test query", intents[0].Query)
		assert.Equal(t, "informational", intents[0].Intent)
	})

	t.Run("GetSearchResults", func(t *testing.T) {
		t.Parallel()

		results := resp.GetSearchResults()
		assert.Len(t, results, 1)
		assert.Equal(t, "Test Result", results[0].Title)
		assert.Equal(t, "https://example.com", results[0].Link)
	})

	t.Run("GetSearchRecommendations", func(t *testing.T) {
		t.Parallel()

		recommends := resp.GetSearchRecommendations()
		assert.Len(t, recommends, 1)
		assert.Equal(t, "recommended query", recommends[0].Query)
	})
}

func TestWebSearchChunk(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal streaming chunk", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "ws_chunk_123",
			"created": 1700000000,
			"choices": [
				{
					"index": 0,
					"delta": {
						"role": "assistant",
						"tool_calls": [
							{
								"index": 0,
								"id": "call_1",
								"type": "web_search",
								"search_intent": {
									"index": 0,
									"query": "streaming query",
									"intent": "commercial",
									"keywords": "stream, test"
								}
							}
						]
					}
				}
			]
		}`

		var chunk WebSearchChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		assert.Equal(t, "ws_chunk_123", chunk.ID)
		assert.Equal(t, int64(1700000000), chunk.Created)
		assert.Len(t, chunk.Choices, 1)

		choice := chunk.Choices[0]
		assert.Equal(t, 0, choice.Index)
		assert.Equal(t, "assistant", choice.Delta.Role)
		assert.Len(t, choice.Delta.ToolCalls, 1)

		toolCall := choice.Delta.ToolCalls[0]
		assert.Equal(t, "call_1", toolCall.ID)
		assert.Equal(t, "web_search", toolCall.Type)
		assert.NotNil(t, toolCall.SearchIntent)
		assert.Equal(t, "streaming query", toolCall.SearchIntent.Query)
	})

	t.Run("chunk with search result delta", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "ws_chunk_456",
			"created": 1700000001,
			"choices": [
				{
					"index": 0,
					"finish_reason": "stop",
					"delta": {
						"tool_calls": [
							{
								"index": 0,
								"search_result": {
									"index": 0,
									"title": "Streaming Result",
									"link": "https://stream.example.com",
									"content": "Streaming content...",
									"icon": "https://icon.png",
									"media": "Stream Media",
									"refer": "[ref_2]"
								}
							}
						]
					}
				}
			]
		}`

		var chunk WebSearchChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		assert.Equal(t, "ws_chunk_456", chunk.ID)
		assert.Len(t, chunk.Choices, 1)
		assert.Equal(t, "stop", chunk.Choices[0].FinishReason)

		toolCall := chunk.Choices[0].Delta.ToolCalls[0]
		assert.NotNil(t, toolCall.SearchResult)
		assert.Equal(t, "Streaming Result", toolCall.SearchResult.Title)
		assert.Equal(t, "https://stream.example.com", toolCall.SearchResult.Link)
	})

	t.Run("chunk with recommendation delta", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "ws_chunk_789",
			"choices": [
				{
					"index": 0,
					"delta": {
						"tool_calls": [
							{
								"index": 0,
								"search_recommend": {
									"index": 0,
									"query": "try this query instead"
								}
							}
						]
					}
				}
			]
		}`

		var chunk WebSearchChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		toolCall := chunk.Choices[0].Delta.ToolCalls[0]
		assert.NotNil(t, toolCall.SearchRecommend)
		assert.Equal(t, "try this query instead", toolCall.SearchRecommend.Query)
	})
}

func TestSearchIntent(t *testing.T) {
	t.Parallel()

	intent := SearchIntent{
		Index:    0,
		Query:    "optimized search query",
		Intent:   "transactional",
		Keywords: "buy, purchase, order",
	}

	data, err := json.Marshal(intent)
	require.NoError(t, err)

	var decoded SearchIntent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, intent.Query, decoded.Query)
	assert.Equal(t, intent.Intent, decoded.Intent)
	assert.Equal(t, intent.Keywords, decoded.Keywords)
}

func TestSearchResult(t *testing.T) {
	t.Parallel()

	result := SearchResult{
		Index:   0,
		Title:   "Example Article",
		Link:    "https://example.com/article",
		Content: "Article content here...",
		Icon:    "https://example.com/favicon.ico",
		Media:   "Example News",
		Refer:   "[ref_5]",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded SearchResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, result.Title, decoded.Title)
	assert.Equal(t, result.Link, decoded.Link)
	assert.Equal(t, result.Content, decoded.Content)
	assert.Equal(t, result.Media, decoded.Media)
}

func TestSearchRecommend(t *testing.T) {
	t.Parallel()

	recommend := SearchRecommend{
		Index: 0,
		Query: "alternative search query",
	}

	data, err := json.Marshal(recommend)
	require.NoError(t, err)

	var decoded SearchRecommend
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, recommend.Query, decoded.Query)
}
