package zai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	filestypes "github.com/z-ai/zai-sdk-go/api/types/files"
)

func TestFilesService_Upload(t *testing.T) {
	t.Parallel()

	t.Run("successful upload", func(t *testing.T) {
		t.Parallel()

		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/files", r.URL.Path)
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
			assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

			// Parse multipart form
			err := r.ParseMultipartForm(32 << 20) // 32 MB
			assert.NoError(t, err)

			// Verify purpose field
			purpose := r.FormValue("purpose")
			assert.Equal(t, "fine-tune", purpose)

			// Verify file
			file, header, err := r.FormFile("file")
			assert.NoError(t, err)
			assert.Equal(t, "training.jsonl", header.Filename)

			// Read file content
			content, err := io.ReadAll(file)
			assert.NoError(t, err)
			assert.Contains(t, string(content), "test data")
			file.Close()

			// Send response
			resp := filestypes.File{
				ID:        "file-abc123",
				Object:    "file",
				Bytes:     int64(len(content)),
				CreatedAt: 1677652288,
				Filename:  "training.jsonl",
				Purpose:   filestypes.PurposeFineTune,
				Status:    filestypes.StatusUploaded,
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

		// Upload file
		fileContent := strings.NewReader("test data")
		req := filestypes.NewFileUploadRequest(fileContent, "training.jsonl", filestypes.PurposeFineTune)

		uploadedFile, err := client.Files.Upload(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, uploadedFile)

		assert.Equal(t, "file-abc123", uploadedFile.ID)
		assert.Equal(t, "training.jsonl", uploadedFile.Filename)
		assert.Equal(t, filestypes.PurposeFineTune, uploadedFile.Purpose)
		assert.Equal(t, filestypes.StatusUploaded, uploadedFile.Status)
	})

	t.Run("API error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid file format",
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

		fileContent := strings.NewReader("invalid data")
		req := filestypes.NewFileUploadRequest(fileContent, "bad.txt", filestypes.PurposeFineTune)

		uploadedFile, err := client.Files.Upload(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, uploadedFile)
		assert.Contains(t, err.Error(), "Invalid file format")
	})
}

func TestFilesService_List(t *testing.T) {
	t.Parallel()

	t.Run("successful list", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/files", r.URL.Path)

			resp := filestypes.FileListResponse{
				Object: "list",
				Data: []filestypes.File{
					{
						ID:       "file-1",
						Filename: "file1.jsonl",
						Purpose:  filestypes.PurposeFineTune,
						Status:   filestypes.StatusUploaded,
					},
					{
						ID:       "file-2",
						Filename: "file2.txt",
						Purpose:  filestypes.PurposeAssistants,
						Status:   filestypes.StatusProcessed,
					},
				},
				HasMore: false,
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

		fileList, err := client.Files.List(context.Background())
		require.NoError(t, err)
		require.NotNil(t, fileList)

		assert.Equal(t, "list", fileList.Object)
		assert.Len(t, fileList.Data, 2)
		assert.False(t, fileList.HasMore)

		assert.Equal(t, "file-1", fileList.Data[0].ID)
		assert.Equal(t, "file-2", fileList.Data[1].ID)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := filestypes.FileListResponse{
				Object:  "list",
				Data:    []filestypes.File{},
				HasMore: false,
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

		fileList, err := client.Files.List(context.Background())
		require.NoError(t, err)
		assert.Len(t, fileList.Data, 0)
	})

	t.Run("API error", func(t *testing.T) {
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

		fileList, err := client.Files.List(context.Background())
		assert.Error(t, err)
		assert.Nil(t, fileList)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}

func TestFilesService_Retrieve(t *testing.T) {
	t.Parallel()

	t.Run("successful retrieve", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/files/file-abc123", r.URL.Path)

			resp := filestypes.File{
				ID:        "file-abc123",
				Object:    "file",
				Bytes:     2048,
				CreatedAt: 1677652288,
				Filename:  "data.jsonl",
				Purpose:   filestypes.PurposeFineTune,
				Status:    filestypes.StatusProcessed,
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

		file, err := client.Files.Retrieve(context.Background(), "file-abc123")
		require.NoError(t, err)
		require.NotNil(t, file)

		assert.Equal(t, "file-abc123", file.ID)
		assert.Equal(t, "data.jsonl", file.Filename)
		assert.Equal(t, int64(2048), file.Bytes)
		assert.Equal(t, filestypes.PurposeFineTune, file.Purpose)
		assert.Equal(t, filestypes.StatusProcessed, file.Status)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "File not found",
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

		file, err := client.Files.Retrieve(context.Background(), "file-nonexistent")
		assert.Error(t, err)
		assert.Nil(t, file)
		assert.Contains(t, err.Error(), "File not found")
	})
}

func TestFilesService_Delete(t *testing.T) {
	t.Parallel()

	t.Run("successful delete", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/files/file-abc123", r.URL.Path)

			resp := filestypes.FileDeleteResponse{
				ID:      "file-abc123",
				Object:  "file",
				Deleted: true,
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

		deleteResp, err := client.Files.Delete(context.Background(), "file-abc123")
		require.NoError(t, err)
		require.NotNil(t, deleteResp)

		assert.Equal(t, "file-abc123", deleteResp.ID)
		assert.True(t, deleteResp.Deleted)
		assert.True(t, deleteResp.IsDeleted())
	})

	t.Run("delete failed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Cannot delete file in use",
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

		deleteResp, err := client.Files.Delete(context.Background(), "file-in-use")
		assert.Error(t, err)
		assert.Nil(t, deleteResp)
		assert.Contains(t, err.Error(), "Cannot delete file in use")
	})
}

func TestFilesService_RetrieveContent(t *testing.T) {
	t.Parallel()

	t.Run("successful content retrieval", func(t *testing.T) {
		t.Parallel()

		fileContent := "line1\nline2\nline3"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/files/file-abc123/content", r.URL.Path)

			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fileContent))
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		content, err := client.Files.RetrieveContent(context.Background(), "file-abc123")
		require.NoError(t, err)
		require.NotNil(t, content)

		assert.Equal(t, []byte(fileContent), content.Content)
		assert.Equal(t, fileContent, content.String())
		assert.Equal(t, "text/plain", content.ContentType)
	})

	t.Run("JSON content", func(t *testing.T) {
		t.Parallel()

		jsonContent := `{"key": "value"}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(jsonContent))
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		content, err := client.Files.RetrieveContent(context.Background(), "file-json")
		require.NoError(t, err)
		assert.Equal(t, "application/json", content.ContentType)
		assert.Equal(t, jsonContent, content.String())
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "File not found",
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

		content, err := client.Files.RetrieveContent(context.Background(), "file-nonexistent")
		assert.Error(t, err)
		assert.Nil(t, content)
	})
}

func TestClient_FilesService_Integration(t *testing.T) {
	t.Parallel()

	t.Run("client has files service", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client.Files)
	})

	t.Run("complete workflow", func(t *testing.T) {
		t.Parallel()

		var uploadedFileID string

		// Simulate a complete file management workflow
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/files":
				// Upload
				uploadedFileID = "file-workflow-123"
				resp := filestypes.File{
					ID:        uploadedFileID,
					Object:    "file",
					Bytes:     100,
					CreatedAt: 1677652288,
					Filename:  "workflow.txt",
					Purpose:   filestypes.PurposeAssistants,
					Status:    filestypes.StatusUploaded,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)

			case r.Method == http.MethodGet && r.URL.Path == "/files":
				// List
				resp := filestypes.FileListResponse{
					Object: "list",
					Data: []filestypes.File{
						{ID: uploadedFileID, Filename: "workflow.txt"},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)

			case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/content"):
				// Retrieve content
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("file content"))

			case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/files/"):
				// Retrieve
				resp := filestypes.File{
					ID:       uploadedFileID,
					Filename: "workflow.txt",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)

			case r.Method == http.MethodDelete:
				// Delete
				resp := filestypes.FileDeleteResponse{
					ID:      uploadedFileID,
					Object:  "file",
					Deleted: true,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}
		}))
		defer server.Close()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL(server.URL),
		)
		require.NoError(t, err)
		defer client.Close()

		ctx := context.Background()

		// Upload a file
		uploadReq := filestypes.NewFileUploadRequest(
			strings.NewReader("test content"),
			"workflow.txt",
			filestypes.PurposeAssistants,
		)
		uploadedFile, err := client.Files.Upload(ctx, uploadReq)
		require.NoError(t, err)
		assert.NotEmpty(t, uploadedFile.ID)

		// List files
		fileList, err := client.Files.List(ctx)
		require.NoError(t, err)
		assert.Len(t, fileList.Data, 1)

		// Retrieve file info
		file, err := client.Files.Retrieve(ctx, uploadedFile.ID)
		require.NoError(t, err)
		assert.Equal(t, "workflow.txt", file.Filename)

		// Retrieve file content
		content, err := client.Files.RetrieveContent(ctx, uploadedFile.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, content.Content)

		// Delete file
		deleteResp, err := client.Files.Delete(ctx, uploadedFile.ID)
		require.NoError(t, err)
		assert.True(t, deleteResp.IsDeleted())
	})
}
