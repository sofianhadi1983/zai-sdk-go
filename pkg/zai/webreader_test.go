package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/webreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebReaderService_Read(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/reader", r.URL.Path)

		// Parse request body
		var req webreader.Request
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "https://example.com", req.URL)

		// Send response
		resp := webreader.Response{
			ReaderResult: &webreader.ReaderData{
				Title:       "Example Domain",
				Description: "Example description",
				URL:         "https://example.com",
				Content:     "This domain is for use in illustrative examples in documents.",
				Images: map[string]string{
					"logo": "https://example.com/logo.png",
				},
				Links: map[string]string{
					"more": "https://www.iana.org/domains/example",
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

	req := webreader.NewRequest("https://example.com")

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.HasResult())
	result := resp.GetResult()
	assert.Equal(t, "Example Domain", result.GetTitle())
	assert.Equal(t, "Example description", result.GetDescription())
	assert.True(t, result.HasContent())
	assert.Contains(t, result.GetContent(), "illustrative examples")
	assert.Len(t, result.GetImages(), 1)
	assert.Len(t, result.GetLinks(), 1)
}

func TestWebReaderService_Read_WithOptions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req webreader.Request
		json.NewDecoder(r.Body).Decode(&req)

		// Verify all options
		assert.Equal(t, "https://blog.example.com", req.URL)
		assert.Equal(t, "req_123", req.RequestID)
		assert.Equal(t, "user_456", req.UserID)
		assert.Equal(t, "30", req.Timeout)
		assert.True(t, req.NoCache)
		assert.Equal(t, "markdown", req.ReturnFormat)
		assert.True(t, req.RetainImages)
		assert.True(t, req.WithImagesSummary)
		assert.True(t, req.WithLinksSummary)

		resp := webreader.Response{
			ReaderResult: &webreader.ReaderData{
				Title:   "Blog Post",
				Content: "# Blog Post\n\nContent here.",
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

	req := webreader.NewRequest("https://blog.example.com").
		SetRequestID("req_123").
		SetUserID("user_456").
		SetTimeout("30").
		SetNoCache(true).
		SetReturnFormat("markdown").
		SetRetainImages(true).
		SetWithImagesSummary(true).
		SetWithLinksSummary(true)

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.HasResult())
	assert.Equal(t, "Blog Post", resp.GetTitle())
	assert.Contains(t, resp.GetContent(), "# Blog Post")
}

func TestWebReaderService_Read_TextFormat(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req webreader.Request
		json.NewDecoder(r.Body).Decode(&req)

		assert.Equal(t, "text", req.ReturnFormat)

		resp := webreader.Response{
			ReaderResult: &webreader.ReaderData{
				Title:   "Text Article",
				Content: "Plain text content without any markdown formatting.",
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

	req := webreader.NewRequest("https://example.com/article").
		SetReturnFormat("text")

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "Plain text content without any markdown formatting.", resp.GetContent())
}

func TestWebReaderService_Read_NoResult(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send empty response
		resp := webreader.Response{}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := webreader.NewRequest("https://example.com")

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, req)

	assert.False(t, resp.HasResult())
	assert.Nil(t, resp.GetResult())
	assert.Equal(t, "", resp.GetTitle())
	assert.Equal(t, "", resp.GetContent())
}

func TestWebReaderService_Read_WithMetadata(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := webreader.Response{
			ReaderResult: &webreader.ReaderData{
				Title:         "Article with Metadata",
				Content:       "Article content",
				PublishedTime: "2024-01-01T12:00:00Z",
				Metadata: map[string]interface{}{
					"author": "John Doe",
					"tags":   []interface{}{"tech", "golang"},
				},
				External: map[string]interface{}{
					"source": "external_api",
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

	req := webreader.NewRequest("https://example.com/article")

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.HasResult())
	result := resp.GetResult()
	assert.Equal(t, "2024-01-01T12:00:00Z", result.PublishedTime)
	assert.NotNil(t, result.Metadata)
	assert.NotNil(t, result.External)
	assert.Equal(t, "John Doe", result.Metadata["author"])
	assert.Equal(t, "external_api", result.External["source"])
}

func TestWebReaderService_Read_NoGFM(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req webreader.Request
		json.NewDecoder(r.Body).Decode(&req)

		assert.True(t, req.NoGFM)
		assert.True(t, req.KeepImgDataURL)

		resp := webreader.Response{
			ReaderResult: &webreader.ReaderData{
				Title:   "Article",
				Content: "Content without GFM",
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

	req := webreader.NewRequest("https://example.com").
		SetNoGFM(true).
		SetKeepImgDataURL(true)

	resp, err := client.WebReader.Read(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
