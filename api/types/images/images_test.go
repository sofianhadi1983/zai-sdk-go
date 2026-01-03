package images

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImageGenerationRequest(t *testing.T) {
	t.Parallel()

	req := NewImageGenerationRequest("cogview-3", "A beautiful sunset")

	assert.Equal(t, "cogview-3", req.Model)
	assert.Equal(t, "A beautiful sunset", req.Prompt)
}

func TestImageGenerationRequest_Setters(t *testing.T) {
	t.Parallel()

	t.Run("SetSize", func(t *testing.T) {
		t.Parallel()

		req := &ImageGenerationRequest{}
		req.SetSize(Size1792x1024)

		assert.Equal(t, Size1792x1024, req.Size)
	})

	t.Run("SetQuality", func(t *testing.T) {
		t.Parallel()

		req := &ImageGenerationRequest{}
		req.SetQuality(QualityHD)

		assert.Equal(t, QualityHD, req.Quality)
	})

	t.Run("SetN", func(t *testing.T) {
		t.Parallel()

		req := &ImageGenerationRequest{}
		req.SetN(5)

		require.NotNil(t, req.N)
		assert.Equal(t, 5, *req.N)
	})

	t.Run("SetResponseFormat", func(t *testing.T) {
		t.Parallel()

		req := &ImageGenerationRequest{}
		req.SetResponseFormat(ResponseFormatB64JSON)

		assert.Equal(t, ResponseFormatB64JSON, req.ResponseFormat)
	})

	t.Run("SetUser", func(t *testing.T) {
		t.Parallel()

		req := &ImageGenerationRequest{}
		req.SetUser("user-789")

		assert.Equal(t, "user-789", req.User)
	})

	t.Run("chained setters", func(t *testing.T) {
		t.Parallel()

		req := NewImageGenerationRequest("cogview-3", "test prompt")
		req.SetSize(Size1024x1792).
			SetQuality(QualityHD).
			SetN(3).
			SetResponseFormat(ResponseFormatURL).
			SetUser("user-456")

		assert.Equal(t, Size1024x1792, req.Size)
		assert.Equal(t, QualityHD, req.Quality)
		require.NotNil(t, req.N)
		assert.Equal(t, 3, *req.N)
		assert.Equal(t, ResponseFormatURL, req.ResponseFormat)
		assert.Equal(t, "user-456", req.User)
	})
}

func TestImageGenerationRequest_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal minimal request", func(t *testing.T) {
		t.Parallel()

		req := NewImageGenerationRequest("cogview-3", "A cat")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded ImageGenerationRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "cogview-3", decoded.Model)
		assert.Equal(t, "A cat", decoded.Prompt)
	})

	t.Run("marshal full request", func(t *testing.T) {
		t.Parallel()

		n := 2
		req := &ImageGenerationRequest{
			Model:          "cogview-3",
			Prompt:         "A dog",
			Size:           Size1024x1024,
			Quality:        QualityHD,
			N:              &n,
			ResponseFormat: ResponseFormatURL,
			User:           "user-123",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "cogview-3")
		assert.Contains(t, string(data), "A dog")
		assert.Contains(t, string(data), "1024x1024")
		assert.Contains(t, string(data), "hd")
		assert.Contains(t, string(data), "\"n\":2")
		assert.Contains(t, string(data), "url")
		assert.Contains(t, string(data), "user-123")
	})

	t.Run("omit empty fields", func(t *testing.T) {
		t.Parallel()

		req := NewImageGenerationRequest("cogview-3", "test")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		// Optional fields should be omitted when not set
		assert.NotContains(t, string(data), "size")
		assert.NotContains(t, string(data), "quality")
		assert.NotContains(t, string(data), "\"n\"")
		assert.NotContains(t, string(data), "response_format")
		assert.NotContains(t, string(data), "user")
	})
}

func TestImageData_Helpers(t *testing.T) {
	t.Parallel()

	t.Run("GetImageURL", func(t *testing.T) {
		t.Parallel()

		img := &ImageData{
			URL: "https://example.com/image.png",
		}

		assert.Equal(t, "https://example.com/image.png", img.GetImageURL())
	})

	t.Run("GetImageURL empty", func(t *testing.T) {
		t.Parallel()

		img := &ImageData{}

		assert.Equal(t, "", img.GetImageURL())
	})

	t.Run("GetBase64Data", func(t *testing.T) {
		t.Parallel()

		img := &ImageData{
			B64JSON: "base64encodeddata==",
		}

		assert.Equal(t, "base64encodeddata==", img.GetBase64Data())
	})

	t.Run("GetBase64Data empty", func(t *testing.T) {
		t.Parallel()

		img := &ImageData{}

		assert.Equal(t, "", img.GetBase64Data())
	})
}

func TestImageGenerationResponse_GetFirstImage(t *testing.T) {
	t.Parallel()

	t.Run("with images", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{URL: "https://example.com/image1.png"},
				{URL: "https://example.com/image2.png"},
			},
		}

		img := resp.GetFirstImage()
		require.NotNil(t, img)
		assert.Equal(t, "https://example.com/image1.png", img.URL)
	})

	t.Run("without images", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{},
		}

		img := resp.GetFirstImage()
		assert.Nil(t, img)
	})
}

func TestImageGenerationResponse_GetImageURLs(t *testing.T) {
	t.Parallel()

	t.Run("all URLs", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{URL: "https://example.com/1.png"},
				{URL: "https://example.com/2.png"},
				{URL: "https://example.com/3.png"},
			},
		}

		urls := resp.GetImageURLs()
		require.Len(t, urls, 3)
		assert.Equal(t, "https://example.com/1.png", urls[0])
		assert.Equal(t, "https://example.com/2.png", urls[1])
		assert.Equal(t, "https://example.com/3.png", urls[2])
	})

	t.Run("mixed formats", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{URL: "https://example.com/1.png"},
				{B64JSON: "base64data1"},
				{URL: "https://example.com/2.png"},
			},
		}

		urls := resp.GetImageURLs()
		require.Len(t, urls, 2)
		assert.Equal(t, "https://example.com/1.png", urls[0])
		assert.Equal(t, "https://example.com/2.png", urls[1])
	})

	t.Run("no URLs", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{B64JSON: "base64data1"},
				{B64JSON: "base64data2"},
			},
		}

		urls := resp.GetImageURLs()
		assert.Len(t, urls, 0)
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{},
		}

		urls := resp.GetImageURLs()
		assert.Len(t, urls, 0)
	})
}

func TestImageGenerationResponse_GetBase64Images(t *testing.T) {
	t.Parallel()

	t.Run("all base64", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{B64JSON: "data1=="},
				{B64JSON: "data2=="},
			},
		}

		images := resp.GetBase64Images()
		require.Len(t, images, 2)
		assert.Equal(t, "data1==", images[0])
		assert.Equal(t, "data2==", images[1])
	})

	t.Run("mixed formats", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{URL: "https://example.com/1.png"},
				{B64JSON: "data1=="},
				{URL: "https://example.com/2.png"},
				{B64JSON: "data2=="},
			},
		}

		images := resp.GetBase64Images()
		require.Len(t, images, 2)
		assert.Equal(t, "data1==", images[0])
		assert.Equal(t, "data2==", images[1])
	})

	t.Run("no base64", func(t *testing.T) {
		t.Parallel()

		resp := &ImageGenerationResponse{
			Data: []ImageData{
				{URL: "https://example.com/1.png"},
			},
		}

		images := resp.GetBase64Images()
		assert.Len(t, images, 0)
	})
}

func TestImageGenerationResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal URL response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"created": 1677652288,
			"data": [
				{
					"url": "https://example.com/image.png",
					"revised_prompt": "A detailed image of a sunset"
				}
			]
		}`

		var resp ImageGenerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, int64(1677652288), resp.Created)
		assert.Len(t, resp.Data, 1)

		img := resp.Data[0]
		assert.Equal(t, "https://example.com/image.png", img.URL)
		assert.Equal(t, "A detailed image of a sunset", img.RevisedPrompt)
	})

	t.Run("unmarshal base64 response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"created": 1677652288,
			"data": [
				{
					"b64_json": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
				}
			]
		}`

		var resp ImageGenerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Data, 1)
		img := resp.Data[0]
		assert.NotEmpty(t, img.B64JSON)
		assert.Empty(t, img.URL)
	})

	t.Run("unmarshal multiple images", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"created": 1677652288,
			"data": [
				{
					"url": "https://example.com/1.png"
				},
				{
					"url": "https://example.com/2.png"
				},
				{
					"url": "https://example.com/3.png"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"total_tokens": 10
			}
		}`

		var resp ImageGenerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Data, 3)
		assert.NotNil(t, resp.Usage)
		assert.Equal(t, 10, resp.Usage.TotalTokens)

		urls := resp.GetImageURLs()
		assert.Len(t, urls, 3)
	})
}

func TestImageSize_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ImageSize("1024x1024"), Size1024x1024)
	assert.Equal(t, ImageSize("1792x1024"), Size1792x1024)
	assert.Equal(t, ImageSize("1024x1792"), Size1024x1792)
}

func TestImageQuality_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ImageQuality("standard"), QualityStandard)
	assert.Equal(t, ImageQuality("hd"), QualityHD)
}

func TestResponseFormat_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ResponseFormat("url"), ResponseFormatURL)
	assert.Equal(t, ResponseFormat("b64_json"), ResponseFormatB64JSON)
}

func TestImageGenerationRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Build a complete request
	req := NewImageGenerationRequest("cogview-3", "A futuristic cityscape at night")
	req.SetSize(Size1792x1024).
		SetQuality(QualityHD).
		SetN(4).
		SetResponseFormat(ResponseFormatURL).
		SetUser("user-abc-123")

	// Verify the request is complete
	assert.Equal(t, "cogview-3", req.Model)
	assert.Equal(t, "A futuristic cityscape at night", req.Prompt)
	assert.Equal(t, Size1792x1024, req.Size)
	assert.Equal(t, QualityHD, req.Quality)
	require.NotNil(t, req.N)
	assert.Equal(t, 4, *req.N)
	assert.Equal(t, ResponseFormatURL, req.ResponseFormat)
	assert.Equal(t, "user-abc-123", req.User)

	// Ensure it can be marshaled
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), "cogview-3")
	assert.Contains(t, string(data), "futuristic cityscape")
	assert.Contains(t, string(data), "1792x1024")
	assert.Contains(t, string(data), "hd")
	assert.Contains(t, string(data), "\"n\":4")
}

func TestImageGenerationResponse_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response
	jsonData := `{
		"created": 1677652288,
		"data": [
			{
				"url": "https://cdn.example.com/generated/abc123.png",
				"revised_prompt": "A highly detailed photograph of a futuristic cityscape at night, with neon lights and flying vehicles"
			},
			{
				"url": "https://cdn.example.com/generated/def456.png",
				"revised_prompt": "A highly detailed photograph of a futuristic cityscape at night, with neon lights and flying vehicles, variation 2"
			}
		],
		"usage": {
			"prompt_tokens": 15,
			"total_tokens": 15
		}
	}`

	var resp ImageGenerationResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, int64(1677652288), resp.Created)
	assert.Len(t, resp.Data, 2)

	// Verify helper methods work
	firstImg := resp.GetFirstImage()
	require.NotNil(t, firstImg)
	assert.Contains(t, firstImg.URL, "abc123.png")
	assert.NotEmpty(t, firstImg.RevisedPrompt)

	urls := resp.GetImageURLs()
	require.Len(t, urls, 2)
	assert.Contains(t, urls[0], "abc123.png")
	assert.Contains(t, urls[1], "def456.png")

	// Verify usage
	require.NotNil(t, resp.Usage)
	assert.Equal(t, 15, resp.Usage.TotalTokens)
}
