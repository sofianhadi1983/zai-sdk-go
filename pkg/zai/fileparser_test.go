package zai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/fileparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileParserService_Create(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/files/parser/create", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify required fields
		assert.Equal(t, "pdf", r.FormValue("file_type"))
		assert.Equal(t, "prime", r.FormValue("tool_type"))

		// Verify file was uploaded
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "test.pdf", header.Filename)

		// Read file content
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "test pdf data", string(content))

		// Send response
		resp := fileparser.CreateResponse{
			TaskID:  "task_123",
			Message: "Task created successfully",
			Success: true,
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

	file := strings.NewReader("test pdf data")
	req := fileparser.NewCreateRequest(file, "test.pdf", "pdf", fileparser.ToolTypePrime)

	resp, err := client.FileParser.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "task_123", resp.TaskID)
	assert.Equal(t, "Task created successfully", resp.Message)
	assert.True(t, resp.Success)
}

func TestFileParserService_Create_WithDifferentToolTypes(t *testing.T) {
	t.Parallel()

	toolTypes := []fileparser.ToolType{
		fileparser.ToolTypeLite,
		fileparser.ToolTypeExpert,
		fileparser.ToolTypePrime,
	}

	for _, toolType := range toolTypes {
		t.Run(string(toolType), func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := r.ParseMultipartForm(10 << 20)
				require.NoError(t, err)

				assert.Equal(t, string(toolType), r.FormValue("tool_type"))

				resp := fileparser.CreateResponse{
					TaskID:  "task_456",
					Message: "Success",
					Success: true,
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
			req := fileparser.NewCreateRequest(file, "test.docx", "docx", toolType)

			resp, err := client.FileParser.Create(context.Background(), req)
			require.NoError(t, err)
			assert.True(t, resp.Success)
		})
	}
}

func TestFileParserService_Content_Text(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/files/parser/result/task_123/text", r.URL.Path)

		// Return text content
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is the parsed text content from the document."))
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := fileparser.NewContentRequest("task_123", fileparser.FormatTypeText)

	resp, err := client.FileParser.Content(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.HasContent())
	assert.False(t, resp.HasData())
	assert.Equal(t, "This is the parsed text content from the document.", resp.GetContent())
}

func TestFileParserService_Content_DownloadLink(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/files/parser/result/task_456/download_link", r.URL.Path)

		// Return binary data
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{0x01, 0x02, 0x03, 0x04})
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := fileparser.NewContentRequest("task_456", fileparser.FormatTypeDownloadLink)

	resp, err := client.FileParser.Content(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, resp.HasContent())
	assert.True(t, resp.HasData())
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, resp.GetData())
}

func TestFileParserService_CreateSync(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/files/parser/sync", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify required fields
		assert.Equal(t, "docx", r.FormValue("file_type"))
		assert.Equal(t, "prime-sync", r.FormValue("tool_type"))

		// Verify file was uploaded
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "document.docx", header.Filename)

		// Send response
		resp := fileparser.SyncResponse{
			TaskID:           "task_789",
			Message:          "Parsing completed",
			Status:           true,
			Content:          "This is the parsed document content.",
			ParsingResultURL: "https://example.com/result.txt",
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

	file := strings.NewReader("test docx data")
	req := fileparser.NewSyncRequest(file, "document.docx", "docx")

	resp, err := client.FileParser.CreateSync(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "task_789", resp.TaskID)
	assert.Equal(t, "Parsing completed", resp.Message)
	assert.True(t, resp.Status)
	assert.True(t, resp.HasContent())
	assert.Equal(t, "This is the parsed document content.", resp.GetContent())
	assert.Equal(t, "https://example.com/result.txt", resp.GetDownloadURL())
}

func TestFileParserService_CreateSync_EmptyContent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fileparser.SyncResponse{
			TaskID:           "task_empty",
			Message:          "No content found",
			Status:           false,
			Content:          "",
			ParsingResultURL: "",
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

	file := strings.NewReader("empty file")
	req := fileparser.NewSyncRequest(file, "empty.pdf", "pdf")

	resp, err := client.FileParser.CreateSync(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, resp.Status)
	assert.False(t, resp.HasContent())
	assert.Equal(t, "", resp.GetContent())
	assert.Equal(t, "", resp.GetDownloadURL())
}
