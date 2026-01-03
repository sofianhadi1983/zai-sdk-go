package websearch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSearchRequest(t *testing.T) {
	t.Parallel()

	t.Run("new request", func(t *testing.T) {
		t.Parallel()

		req := NewWebSearchRequest("artificial intelligence")

		assert.Equal(t, "artificial intelligence", req.SearchQuery)
		assert.Equal(t, SearchEnginePrime, req.SearchEngine)
		assert.Equal(t, 0, req.Count)
		assert.False(t, req.SearchIntent)
		assert.False(t, req.IncludeImage)
	})

	t.Run("builder pattern", func(t *testing.T) {
		t.Parallel()

		req := NewWebSearchRequest("machine learning").
			SetSearchEngine("google").
			SetCount(10).
			SetRecencyFilter(RecencyFilterOneWeek).
			SetContentSize(ContentSizeLarge).
			SetSearchIntent(true).
			SetIncludeImage(true).
			SetDomainFilter("arxiv.org").
			SetRequestID("req_123").
			SetUserID("user_456")

		assert.Equal(t, "machine learning", req.SearchQuery)
		assert.Equal(t, "google", req.SearchEngine)
		assert.Equal(t, 10, req.Count)
		assert.Equal(t, RecencyFilterOneWeek, req.SearchRecencyFilter)
		assert.Equal(t, ContentSizeLarge, req.ContentSize)
		assert.True(t, req.SearchIntent)
		assert.True(t, req.IncludeImage)
		assert.Equal(t, "arxiv.org", req.SearchDomainFilter)
		assert.Equal(t, "req_123", req.RequestID)
		assert.Equal(t, "user_456", req.UserID)
	})

	t.Run("with sensitive word check", func(t *testing.T) {
		t.Parallel()

		sensitiveCheck := &SensitiveWordCheck{
			Type:   SensitiveWordTypeAll,
			Status: SensitiveWordStatusEnable,
		}

		req := NewWebSearchRequest("test query").
			SetSensitiveWordCheck(sensitiveCheck)

		assert.NotNil(t, req.SensitiveWordCheck)
		assert.Equal(t, SensitiveWordTypeAll, req.SensitiveWordCheck.Type)
		assert.Equal(t, SensitiveWordStatusEnable, req.SensitiveWordCheck.Status)
	})

	t.Run("JSON marshaling", func(t *testing.T) {
		t.Parallel()

		req := NewWebSearchRequest("test search").
			SetCount(5).
			SetSearchIntent(true)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded WebSearchRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, req.SearchQuery, decoded.SearchQuery)
		assert.Equal(t, req.Count, decoded.Count)
		assert.Equal(t, req.SearchIntent, decoded.SearchIntent)
	})
}

func TestWebSearchResponse(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal from API response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "search_abc123",
			"request_id": "req_xyz789",
			"created": 1700000000,
			"search_intent": {
				"query": "AI trends 2024",
				"intent": "informational",
				"keywords": "AI, trends, 2024"
			},
			"search_result": [
				{
					"title": "Top AI Trends in 2024",
					"link": "https://example.com/ai-trends",
					"content": "Discover the latest trends in artificial intelligence...",
					"icon": "https://example.com/icon.png",
					"media": "Example News",
					"refer": "[ref_1]",
					"publish_date": "2024-01-15",
					"images": ["https://example.com/image1.jpg"]
				},
				{
					"title": "Future of AI",
					"link": "https://example.com/ai-future",
					"content": "Exploring the future possibilities of AI...",
					"icon": "https://example.com/icon2.png",
					"media": "Tech Blog",
					"refer": "[ref_2]",
					"publish_date": "2024-01-10",
					"images": []
				}
			]
		}`

		var resp WebSearchResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "search_abc123", resp.ID)
		assert.Equal(t, "req_xyz789", resp.RequestID)
		assert.Equal(t, int64(1700000000), resp.Created)

		// Check search intent
		assert.True(t, resp.HasIntent())
		assert.NotNil(t, resp.SearchIntent)
		assert.Equal(t, "AI trends 2024", resp.SearchIntent.Query)
		assert.Equal(t, "informational", resp.SearchIntent.Intent)
		assert.Equal(t, "AI, trends, 2024", resp.SearchIntent.Keywords)

		// Check search results
		results := resp.GetResults()
		assert.Len(t, results, 2)

		// First result
		assert.Equal(t, "Top AI Trends in 2024", results[0].Title)
		assert.Equal(t, "https://example.com/ai-trends", results[0].Link)
		assert.Contains(t, results[0].Content, "artificial intelligence")
		assert.Equal(t, "Example News", results[0].Media)
		assert.Equal(t, "[ref_1]", results[0].Refer)
		assert.Equal(t, "2024-01-15", results[0].PublishDate)
		assert.Len(t, results[0].Images, 1)

		// Second result
		assert.Equal(t, "Future of AI", results[1].Title)
		assert.Equal(t, "[ref_2]", results[1].Refer)
		assert.Empty(t, results[1].Images)
	})

	t.Run("response without intent", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "search_123",
			"search_result": [
				{
					"title": "Test Result",
					"link": "https://test.com",
					"content": "Test content",
					"icon": "",
					"media": "Test",
					"refer": "[ref_1]",
					"publish_date": "2024-01-01"
				}
			]
		}`

		var resp WebSearchResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.False(t, resp.HasIntent())
		assert.Nil(t, resp.SearchIntent)
		assert.Len(t, resp.GetResults(), 1)
	})

	t.Run("empty search results", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "search_empty",
			"search_result": []
		}`

		var resp WebSearchResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Empty(t, resp.GetResults())
		assert.False(t, resp.HasIntent())
	})
}

func TestSensitiveWordCheck(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		t.Parallel()

		check := &SensitiveWordCheck{
			Type:   SensitiveWordTypeAll,
			Status: SensitiveWordStatusEnable,
		}

		data, err := json.Marshal(check)
		require.NoError(t, err)

		var decoded SensitiveWordCheck
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, check.Type, decoded.Type)
		assert.Equal(t, check.Status, decoded.Status)
	})

	t.Run("constants", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "ALL", SensitiveWordTypeAll)
		assert.Equal(t, "ENABLE", SensitiveWordStatusEnable)
		assert.Equal(t, "DISABLE", SensitiveWordStatusDisable)
	})
}

func TestSearchIntentResp(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"query": "best AI tools",
		"intent": "commercial",
		"keywords": "AI, tools, best"
	}`

	var intent SearchIntentResp
	err := json.Unmarshal([]byte(jsonData), &intent)
	require.NoError(t, err)

	assert.Equal(t, "best AI tools", intent.Query)
	assert.Equal(t, "commercial", intent.Intent)
	assert.Equal(t, "AI, tools, best", intent.Keywords)
}

func TestSearchResultResp(t *testing.T) {
	t.Parallel()

	t.Run("with images", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"title": "Article Title",
			"link": "https://example.com/article",
			"content": "Article content here",
			"icon": "https://example.com/favicon.ico",
			"media": "News Site",
			"refer": "[ref_5]",
			"publish_date": "2024-01-20",
			"images": [
				"https://example.com/img1.jpg",
				"https://example.com/img2.jpg"
			]
		}`

		var result SearchResultResp
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Equal(t, "Article Title", result.Title)
		assert.Equal(t, "https://example.com/article", result.Link)
		assert.Equal(t, "Article content here", result.Content)
		assert.Equal(t, "News Site", result.Media)
		assert.Equal(t, "[ref_5]", result.Refer)
		assert.Equal(t, "2024-01-20", result.PublishDate)
		assert.Len(t, result.Images, 2)
	})

	t.Run("without images", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"title": "Simple Article",
			"link": "https://simple.com",
			"content": "Content",
			"icon": "",
			"media": "Blog",
			"refer": "[ref_1]",
			"publish_date": "2024-01-01"
		}`

		var result SearchResultResp
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Equal(t, "Simple Article", result.Title)
		assert.Empty(t, result.Images)
	})
}

func TestRecencyFilterConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "oneDay", RecencyFilterOneDay)
	assert.Equal(t, "oneWeek", RecencyFilterOneWeek)
	assert.Equal(t, "oneMonth", RecencyFilterOneMonth)
	assert.Equal(t, "oneYear", RecencyFilterOneYear)
	assert.Equal(t, "noLimit", RecencyFilterNoLimit)
}

func TestContentSizeConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "small", ContentSizeSmall)
	assert.Equal(t, "medium", ContentSizeMedium)
	assert.Equal(t, "large", ContentSizeLarge)
}
