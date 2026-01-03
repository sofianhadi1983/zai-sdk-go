package webreader

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	t.Parallel()

	url := "https://example.com"
	req := NewRequest(url)

	assert.NotNil(t, req)
	assert.Equal(t, url, req.URL)
	assert.Empty(t, req.RequestID)
	assert.Empty(t, req.UserID)
	assert.Empty(t, req.Timeout)
	assert.False(t, req.NoCache)
	assert.Empty(t, req.ReturnFormat)
	assert.False(t, req.RetainImages)
	assert.False(t, req.NoGFM)
	assert.False(t, req.KeepImgDataURL)
	assert.False(t, req.WithImagesSummary)
	assert.False(t, req.WithLinksSummary)
}

func TestRequest_BuilderMethods(t *testing.T) {
	t.Parallel()

	req := NewRequest("https://example.com").
		SetRequestID("req_123").
		SetUserID("user_456").
		SetTimeout("30").
		SetNoCache(true).
		SetReturnFormat("markdown").
		SetRetainImages(true).
		SetNoGFM(false).
		SetKeepImgDataURL(true).
		SetWithImagesSummary(true).
		SetWithLinksSummary(true)

	assert.Equal(t, "https://example.com", req.URL)
	assert.Equal(t, "req_123", req.RequestID)
	assert.Equal(t, "user_456", req.UserID)
	assert.Equal(t, "30", req.Timeout)
	assert.True(t, req.NoCache)
	assert.Equal(t, "markdown", req.ReturnFormat)
	assert.True(t, req.RetainImages)
	assert.False(t, req.NoGFM)
	assert.True(t, req.KeepImgDataURL)
	assert.True(t, req.WithImagesSummary)
	assert.True(t, req.WithLinksSummary)
}

func TestRequest_JSON(t *testing.T) {
	t.Parallel()

	req := NewRequest("https://example.com").
		SetRequestID("req_123").
		SetReturnFormat("text")

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded Request
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.URL, decoded.URL)
	assert.Equal(t, req.RequestID, decoded.RequestID)
	assert.Equal(t, req.ReturnFormat, decoded.ReturnFormat)
}

func TestReaderData_GetMethods(t *testing.T) {
	t.Parallel()

	data := &ReaderData{
		Title:       "Test Page",
		Description: "A test page description",
		Content:     "This is the page content.",
		Images: map[string]string{
			"img1": "https://example.com/img1.jpg",
			"img2": "https://example.com/img2.jpg",
		},
		Links: map[string]string{
			"link1": "https://example.com/page1",
			"link2": "https://example.com/page2",
		},
	}

	assert.Equal(t, "Test Page", data.GetTitle())
	assert.Equal(t, "A test page description", data.GetDescription())
	assert.Equal(t, "This is the page content.", data.GetContent())
	assert.True(t, data.HasContent())
	assert.Len(t, data.GetImages(), 2)
	assert.Len(t, data.GetLinks(), 2)
	assert.Equal(t, "https://example.com/img1.jpg", data.GetImages()["img1"])
	assert.Equal(t, "https://example.com/page1", data.GetLinks()["link1"])
}

func TestReaderData_EmptyMaps(t *testing.T) {
	t.Parallel()

	data := &ReaderData{
		Title:   "Test",
		Content: "Content",
	}

	images := data.GetImages()
	links := data.GetLinks()

	assert.NotNil(t, images)
	assert.NotNil(t, links)
	assert.Len(t, images, 0)
	assert.Len(t, links, 0)
}

func TestReaderData_HasContent(t *testing.T) {
	t.Parallel()

	t.Run("with content", func(t *testing.T) {
		data := &ReaderData{
			Content: "Some content",
		}
		assert.True(t, data.HasContent())
	})

	t.Run("without content", func(t *testing.T) {
		data := &ReaderData{}
		assert.False(t, data.HasContent())
	})
}

func TestReaderData_JSON(t *testing.T) {
	t.Parallel()

	data := ReaderData{
		Title:         "Test Page",
		Description:   "Description",
		URL:           "https://example.com",
		Content:       "Page content",
		PublishedTime: "2024-01-01T12:00:00Z",
		Images: map[string]string{
			"img1": "https://example.com/img.jpg",
		},
		Links: map[string]string{
			"link1": "https://example.com/link",
		},
		Metadata: map[string]interface{}{
			"author": "John Doe",
		},
	}

	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	var decoded ReaderData
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, data.Title, decoded.Title)
	assert.Equal(t, data.Description, decoded.Description)
	assert.Equal(t, data.URL, decoded.URL)
	assert.Equal(t, data.Content, decoded.Content)
	assert.Equal(t, data.PublishedTime, decoded.PublishedTime)
	assert.Equal(t, data.Images, decoded.Images)
	assert.Equal(t, data.Links, decoded.Links)
}

func TestResponse_GetMethods(t *testing.T) {
	t.Parallel()

	resp := &Response{
		ReaderResult: &ReaderData{
			Title:   "Test Page",
			Content: "Page content",
		},
	}

	assert.True(t, resp.HasResult())
	assert.NotNil(t, resp.GetResult())
	assert.Equal(t, "Test Page", resp.GetTitle())
	assert.Equal(t, "Page content", resp.GetContent())
}

func TestResponse_NoResult(t *testing.T) {
	t.Parallel()

	resp := &Response{}

	assert.False(t, resp.HasResult())
	assert.Nil(t, resp.GetResult())
	assert.Equal(t, "", resp.GetTitle())
	assert.Equal(t, "", resp.GetContent())
}

func TestResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := Response{
		ReaderResult: &ReaderData{
			Title:       "Test Page",
			Description: "Description",
			Content:     "Content",
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded Response
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.ReaderResult)
	assert.Equal(t, resp.ReaderResult.Title, decoded.ReaderResult.Title)
	assert.Equal(t, resp.ReaderResult.Description, decoded.ReaderResult.Description)
	assert.Equal(t, resp.ReaderResult.Content, decoded.ReaderResult.Content)
}

func TestResponse_JSON_NoResult(t *testing.T) {
	t.Parallel()

	resp := Response{}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded Response
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Nil(t, decoded.ReaderResult)
	assert.False(t, decoded.HasResult())
}
