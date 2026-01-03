package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/websearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSearchService_Search(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/web_search", r.URL.Path)

		// Parse request body
		var reqBody websearch.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "latest AI trends 2024", reqBody.SearchQuery)
		assert.Equal(t, 10, reqBody.Count)
		assert.True(t, reqBody.SearchIntent)
		assert.True(t, reqBody.IncludeImage)
		assert.Equal(t, websearch.RecencyFilterWeek, reqBody.SearchRecencyFilter)
		assert.Equal(t, websearch.ContentSizeLarge, reqBody.ContentSize)

		// Send mock response
		resp := websearch.WebSearchResponse{
			ID:        "search_abc123",
			RequestID: "req_xyz789",
			Created:   1700000000,
			SearchIntent: &websearch.SearchIntentResp{
				Query:    "AI trends 2024",
				Intent:   "informational",
				Keywords: "AI, trends, 2024, technology",
			},
			SearchResult: []websearch.SearchResultResp{
				{
					Title:       "Top AI Trends in 2024",
					Link:        "https://example.com/ai-trends-2024",
					Content:     "Discover the most significant artificial intelligence trends shaping 2024...",
					Icon:        "https://example.com/favicon.ico",
					Media:       "Tech News",
					Refer:       "[ref_1]",
					PublishDate: "2024-01-15",
					Images:      []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
				},
				{
					Title:       "Future of Artificial Intelligence",
					Link:        "https://example.com/ai-future",
					Content:     "Exploring the transformative potential of AI in various industries...",
					Icon:        "https://example.com/icon.png",
					Media:       "AI Journal",
					Refer:       "[ref_2]",
					PublishDate: "2024-01-10",
					Images:      []string{"https://example.com/ai-image.jpg"},
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

	req := websearch.NewWebSearchRequest("latest AI trends 2024").
		SetCount(10).
		SetSearchIntent(true).
		SetIncludeImage(true).
		SetRecencyFilter(websearch.RecencyFilterWeek).
		SetContentSize(websearch.ContentSizeLarge)

	resp, err := client.WebSearch.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "search_abc123", resp.ID)
	assert.Equal(t, "req_xyz789", resp.RequestID)
	assert.Equal(t, int64(1700000000), resp.Created)

	// Verify search intent
	assert.True(t, resp.HasIntent())
	assert.NotNil(t, resp.SearchIntent)
	assert.Equal(t, "AI trends 2024", resp.SearchIntent.Query)
	assert.Equal(t, "informational", resp.SearchIntent.Intent)
	assert.Equal(t, "AI, trends, 2024, technology", resp.SearchIntent.Keywords)

	// Verify search results
	results := resp.GetResults()
	assert.Len(t, results, 2)

	// Check first result
	assert.Equal(t, "Top AI Trends in 2024", results[0].Title)
	assert.Equal(t, "https://example.com/ai-trends-2024", results[0].Link)
	assert.Contains(t, results[0].Content, "artificial intelligence")
	assert.Equal(t, "Tech News", results[0].Media)
	assert.Equal(t, "[ref_1]", results[0].Refer)
	assert.Equal(t, "2024-01-15", results[0].PublishDate)
	assert.Len(t, results[0].Images, 2)

	// Check second result
	assert.Equal(t, "Future of Artificial Intelligence", results[1].Title)
	assert.Equal(t, "[ref_2]", results[1].Refer)
	assert.Len(t, results[1].Images, 1)
}

func TestWebSearchService_Search_WithFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/web_search", r.URL.Path)

		// Parse request body
		var reqBody websearch.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify filters
		assert.Equal(t, "arxiv.org", reqBody.SearchDomainFilter)
		assert.Equal(t, websearch.RecencyFilterMonth, reqBody.SearchRecencyFilter)
		assert.Equal(t, websearch.ContentSizeMedium, reqBody.ContentSize)

		// Send mock response
		resp := websearch.WebSearchResponse{
			ID: "search_filtered",
			SearchResult: []websearch.SearchResultResp{
				{
					Title:       "Research Paper on ML",
					Link:        "https://arxiv.org/paper123",
					Content:     "Abstract of the research paper...",
					Icon:        "",
					Media:       "arXiv",
					Refer:       "[ref_1]",
					PublishDate: "2024-01-05",
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

	req := websearch.NewWebSearchRequest("machine learning research").
		SetDomainFilter("arxiv.org").
		SetRecencyFilter(websearch.RecencyFilterMonth).
		SetContentSize(websearch.ContentSizeMedium)

	resp, err := client.WebSearch.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "search_filtered", resp.ID)
	results := resp.GetResults()
	assert.Len(t, results, 1)
	assert.Equal(t, "Research Paper on ML", results[0].Title)
	assert.Equal(t, "https://arxiv.org/paper123", results[0].Link)
}

func TestWebSearchService_Search_WithSensitiveWordCheck(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/web_search", r.URL.Path)

		// Parse request body
		var reqBody websearch.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify sensitive word check
		assert.NotNil(t, reqBody.SensitiveWordCheck)
		assert.Equal(t, websearch.SensitiveWordTypeAll, reqBody.SensitiveWordCheck.Type)
		assert.Equal(t, websearch.SensitiveWordStatusEnable, reqBody.SensitiveWordCheck.Status)

		// Send mock response
		resp := websearch.WebSearchResponse{
			ID:           "search_checked",
			SearchResult: []websearch.SearchResultResp{},
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

	sensitiveCheck := &websearch.SensitiveWordCheck{
		Type:   websearch.SensitiveWordTypeAll,
		Status: websearch.SensitiveWordStatusEnable,
	}

	req := websearch.NewWebSearchRequest("test query").
		SetSensitiveWordCheck(sensitiveCheck)

	resp, err := client.WebSearch.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "search_checked", resp.ID)
	assert.Empty(t, resp.GetResults())
}

func TestWebSearchService_Search_WithoutIntent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/web_search", r.URL.Path)

		// Parse request body
		var reqBody websearch.WebSearchRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		assert.False(t, reqBody.SearchIntent)

		// Send mock response without intent
		resp := websearch.WebSearchResponse{
			ID: "search_no_intent",
			SearchResult: []websearch.SearchResultResp{
				{
					Title:       "Simple Result",
					Link:        "https://example.com",
					Content:     "Content here",
					Icon:        "",
					Media:       "Blog",
					Refer:       "[ref_1]",
					PublishDate: "2024-01-01",
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

	req := websearch.NewWebSearchRequest("simple search")

	resp, err := client.WebSearch.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "search_no_intent", resp.ID)
	assert.False(t, resp.HasIntent())
	assert.Nil(t, resp.SearchIntent)
	assert.Len(t, resp.GetResults(), 1)
}

func TestWebSearchService_Search_EmptyResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/web_search", r.URL.Path)

		// Send mock response with no results
		resp := websearch.WebSearchResponse{
			ID:           "search_empty",
			SearchResult: []websearch.SearchResultResp{},
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

	req := websearch.NewWebSearchRequest("nonexistent query")

	resp, err := client.WebSearch.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "search_empty", resp.ID)
	assert.Empty(t, resp.GetResults())
}

func TestWebSearchService_Search_APIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Invalid search query",
				"type":    "invalid_request_error",
				"code":    "invalid_search_query",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := websearch.NewWebSearchRequest("")

	_, err = client.WebSearch.Search(context.Background(), req)
	require.Error(t, err)
}
