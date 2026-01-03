package zai

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/ocr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRService_HandwritingOCR(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/files/ocr", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		require.NoError(t, err)

		// Verify required fields
		assert.Equal(t, "hand_write", r.FormValue("tool_type"))

		// Verify file was uploaded
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "test.jpg", header.Filename)

		// Read file content
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "test image data", string(content))

		// Send response
		resp := ocr.OCRResponse{
			TaskID:         "task_123",
			Message:        "success",
			Status:         "completed",
			WordsResultNum: 2,
			WordsResult: []ocr.WordsResult{
				{
					Location: ocr.Location{
						Left:   10,
						Top:    20,
						Width:  100,
						Height: 50,
					},
					Words: "Hello",
					Probability: &ocr.Probability{
						Average:  0.95,
						Variance: 0.02,
						Min:      0.90,
					},
				},
				{
					Location: ocr.Location{
						Left:   120,
						Top:    20,
						Width:  100,
						Height: 50,
					},
					Words: "World",
					Probability: &ocr.Probability{
						Average:  0.98,
						Variance: 0.01,
						Min:      0.95,
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

	file := strings.NewReader("test image data")
	req := ocr.NewOCRRequest(file, "test.jpg", ocr.ToolTypeHandWrite)

	resp, err := client.OCR.HandwritingOCR(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "task_123", resp.TaskID)
	assert.Equal(t, "success", resp.Message)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, 2, resp.WordsResultNum)
	assert.True(t, resp.HasResults())

	results := resp.GetResults()
	assert.Len(t, results, 2)
	assert.Equal(t, "Hello", results[0].Words)
	assert.Equal(t, "World", results[1].Words)

	text := resp.GetText()
	assert.Equal(t, "Hello World", text)
}

func TestOCRService_HandwritingOCR_WithOptions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify all fields
		assert.Equal(t, "hand_write", r.FormValue("tool_type"))
		assert.Equal(t, "zh-CN", r.FormValue("language_type"))
		assert.Equal(t, "true", r.FormValue("probability"))

		resp := ocr.OCRResponse{
			TaskID:         "task_456",
			Message:        "success",
			Status:         "completed",
			WordsResultNum: 1,
			WordsResult: []ocr.WordsResult{
				{
					Location: ocr.Location{Left: 10, Top: 20, Width: 100, Height: 50},
					Words:    "你好",
					Probability: &ocr.Probability{
						Average:  0.95,
						Variance: 0.02,
						Min:      0.90,
					},
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

	file := strings.NewReader("test image data")
	req := ocr.NewOCRRequest(file, "chinese.jpg", ocr.ToolTypeHandWrite).
		SetLanguageType("zh-CN").
		SetProbability(true)

	resp, err := client.OCR.HandwritingOCR(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "task_456", resp.TaskID)
	assert.True(t, resp.HasResults())
	assert.Equal(t, "你好", resp.GetText())
}

func TestOCRService_HandwritingOCR_NoProbability(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify probability field is not set
		_, exists := r.MultipartForm.Value["probability"]
		assert.False(t, exists, "probability field should not be present when false")

		resp := ocr.OCRResponse{
			TaskID:         "task_789",
			Message:        "success",
			Status:         "completed",
			WordsResultNum: 1,
			WordsResult: []ocr.WordsResult{
				{
					Location: ocr.Location{Left: 10, Top: 20, Width: 100, Height: 50},
					Words:    "Test",
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

	file := strings.NewReader("test image data")
	req := ocr.NewOCRRequest(file, "test.jpg", ocr.ToolTypeHandWrite)

	resp, err := client.OCR.HandwritingOCR(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "Test", resp.GetText())
	// Probability should be nil when not requested
	if len(resp.WordsResult) > 0 {
		assert.Nil(t, resp.WordsResult[0].Probability)
	}
}

func TestOCRService_HandwritingOCR_NoResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ocr.OCRResponse{
			TaskID:         "task_empty",
			Message:        "no text found",
			Status:         "completed",
			WordsResultNum: 0,
			WordsResult:    []ocr.WordsResult{},
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

	file := strings.NewReader("test image data")
	req := ocr.NewOCRRequest(file, "empty.jpg", ocr.ToolTypeHandWrite)

	resp, err := client.OCR.HandwritingOCR(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, resp.HasResults())
	assert.Equal(t, "", resp.GetText())
	assert.Len(t, resp.GetResults(), 0)
}

func TestOCRService_HandwritingOCR_MultipartCreation(t *testing.T) {
	t.Parallel()

	// Test that multipart form is created correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		assert.Contains(t, contentType, "multipart/form-data")
		assert.Contains(t, contentType, "boundary=")

		// Parse multipart
		boundary := strings.Split(contentType, "boundary=")[1]
		reader := multipart.NewReader(r.Body, boundary)
		form, err := reader.ReadForm(10 << 20)
		require.NoError(t, err)

		// Verify form values exist
		assert.NotEmpty(t, form.Value["tool_type"])
		assert.NotEmpty(t, form.File["file"])

		resp := ocr.OCRResponse{
			TaskID:         "task_multipart",
			Message:        "success",
			Status:         "completed",
			WordsResultNum: 0,
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

	file := strings.NewReader("test data")
	req := ocr.NewOCRRequest(file, "test.png", ocr.ToolTypeHandWrite)

	resp, err := client.OCR.HandwritingOCR(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
