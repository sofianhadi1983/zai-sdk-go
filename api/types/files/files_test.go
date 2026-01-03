package files

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileUploadRequest(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("test content")
	req := NewFileUploadRequest(file, "test.txt", PurposeFineTune)

	assert.NotNil(t, req.File)
	assert.Equal(t, "test.txt", req.Filename)
	assert.Equal(t, PurposeFineTune, req.Purpose)
}

func TestFile_Getters(t *testing.T) {
	t.Parallel()

	file := &File{
		ID:       "file-123",
		Filename: "data.jsonl",
		Bytes:    1024,
		Purpose:  PurposeFineTune,
		Status:   StatusUploaded,
	}

	assert.Equal(t, "file-123", file.GetID())
	assert.Equal(t, "data.jsonl", file.GetFilename())
	assert.Equal(t, int64(1024), file.GetSize())
	assert.Equal(t, PurposeFineTune, file.GetPurpose())
}

func TestFile_IsUploaded(t *testing.T) {
	t.Parallel()

	t.Run("uploaded status", func(t *testing.T) {
		t.Parallel()

		file := &File{Status: StatusUploaded}
		assert.True(t, file.IsUploaded())
	})

	t.Run("processed status", func(t *testing.T) {
		t.Parallel()

		file := &File{Status: StatusProcessed}
		assert.True(t, file.IsUploaded())
	})

	t.Run("error status", func(t *testing.T) {
		t.Parallel()

		file := &File{Status: StatusError}
		assert.False(t, file.IsUploaded())
	})

	t.Run("empty status", func(t *testing.T) {
		t.Parallel()

		file := &File{}
		assert.False(t, file.IsUploaded())
	})
}

func TestFile_HasError(t *testing.T) {
	t.Parallel()

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		file := &File{
			Status:        StatusError,
			StatusDetails: "Processing failed",
		}
		assert.True(t, file.HasError())
	})

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		file := &File{Status: StatusUploaded}
		assert.False(t, file.HasError())
	})
}

func TestFile_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal file", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "file-abc123",
			"object": "file",
			"bytes": 2048,
			"created_at": 1677652288,
			"filename": "training_data.jsonl",
			"purpose": "fine-tune",
			"status": "uploaded"
		}`

		var file File
		err := json.Unmarshal([]byte(jsonData), &file)
		require.NoError(t, err)

		assert.Equal(t, "file-abc123", file.ID)
		assert.Equal(t, "file", file.Object)
		assert.Equal(t, int64(2048), file.Bytes)
		assert.Equal(t, int64(1677652288), file.CreatedAt)
		assert.Equal(t, "training_data.jsonl", file.Filename)
		assert.Equal(t, PurposeFineTune, file.Purpose)
		assert.Equal(t, StatusUploaded, file.Status)
	})

	t.Run("unmarshal file with error", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "file-error",
			"object": "file",
			"bytes": 512,
			"created_at": 1677652288,
			"filename": "bad_data.jsonl",
			"purpose": "fine-tune",
			"status": "error",
			"status_details": "Invalid JSON format"
		}`

		var file File
		err := json.Unmarshal([]byte(jsonData), &file)
		require.NoError(t, err)

		assert.Equal(t, StatusError, file.Status)
		assert.Equal(t, "Invalid JSON format", file.StatusDetails)
		assert.True(t, file.HasError())
	})

	t.Run("marshal file", func(t *testing.T) {
		t.Parallel()

		file := &File{
			ID:        "file-123",
			Object:    "file",
			Bytes:     1024,
			CreatedAt: 1677652288,
			Filename:  "test.txt",
			Purpose:   PurposeAssistants,
			Status:    StatusProcessed,
		}

		data, err := json.Marshal(file)
		require.NoError(t, err)

		assert.Contains(t, string(data), "file-123")
		assert.Contains(t, string(data), "test.txt")
		assert.Contains(t, string(data), "assistants")
		assert.Contains(t, string(data), "processed")
	})
}

func TestFileListResponse_GetFiles(t *testing.T) {
	t.Parallel()

	resp := &FileListResponse{
		Object: "list",
		Data: []File{
			{ID: "file-1", Filename: "file1.txt"},
			{ID: "file-2", Filename: "file2.txt"},
		},
	}

	files := resp.GetFiles()
	assert.Len(t, files, 2)
	assert.Equal(t, "file-1", files[0].ID)
	assert.Equal(t, "file-2", files[1].ID)
}

func TestFileListResponse_GetFileIDs(t *testing.T) {
	t.Parallel()

	resp := &FileListResponse{
		Data: []File{
			{ID: "file-1"},
			{ID: "file-2"},
			{ID: "file-3"},
		},
	}

	ids := resp.GetFileIDs()
	require.Len(t, ids, 3)
	assert.Equal(t, "file-1", ids[0])
	assert.Equal(t, "file-2", ids[1])
	assert.Equal(t, "file-3", ids[2])
}

func TestFileListResponse_GetFileByID(t *testing.T) {
	t.Parallel()

	t.Run("file exists", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{
				{ID: "file-1", Filename: "file1.txt"},
				{ID: "file-2", Filename: "file2.txt"},
				{ID: "file-3", Filename: "file3.txt"},
			},
		}

		file := resp.GetFileByID("file-2")
		require.NotNil(t, file)
		assert.Equal(t, "file-2", file.ID)
		assert.Equal(t, "file2.txt", file.Filename)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{
				{ID: "file-1"},
			},
		}

		file := resp.GetFileByID("file-999")
		assert.Nil(t, file)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{},
		}

		file := resp.GetFileByID("file-1")
		assert.Nil(t, file)
	})
}

func TestFileListResponse_GetFilesByPurpose(t *testing.T) {
	t.Parallel()

	t.Run("filter by purpose", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{
				{ID: "file-1", Purpose: PurposeFineTune},
				{ID: "file-2", Purpose: PurposeAssistants},
				{ID: "file-3", Purpose: PurposeFineTune},
				{ID: "file-4", Purpose: PurposeBatch},
			},
		}

		fineTuneFiles := resp.GetFilesByPurpose(PurposeFineTune)
		require.Len(t, fineTuneFiles, 2)
		assert.Equal(t, "file-1", fineTuneFiles[0].ID)
		assert.Equal(t, "file-3", fineTuneFiles[1].ID)
	})

	t.Run("no matching files", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{
				{ID: "file-1", Purpose: PurposeFineTune},
			},
		}

		assistantFiles := resp.GetFilesByPurpose(PurposeAssistants)
		assert.Len(t, assistantFiles, 0)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Data: []File{},
		}

		files := resp.GetFilesByPurpose(PurposeFineTune)
		assert.Len(t, files, 0)
	})
}

func TestFileListResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal file list", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"object": "list",
			"data": [
				{
					"id": "file-1",
					"object": "file",
					"bytes": 1024,
					"created_at": 1677652288,
					"filename": "file1.txt",
					"purpose": "assistants",
					"status": "uploaded"
				},
				{
					"id": "file-2",
					"object": "file",
					"bytes": 2048,
					"created_at": 1677652289,
					"filename": "file2.txt",
					"purpose": "fine-tune",
					"status": "processed"
				}
			],
			"has_more": false
		}`

		var resp FileListResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "list", resp.Object)
		assert.Len(t, resp.Data, 2)
		assert.False(t, resp.HasMore)

		assert.Equal(t, "file-1", resp.Data[0].ID)
		assert.Equal(t, PurposeAssistants, resp.Data[0].Purpose)

		assert.Equal(t, "file-2", resp.Data[1].ID)
		assert.Equal(t, PurposeFineTune, resp.Data[1].Purpose)
	})

	t.Run("marshal file list", func(t *testing.T) {
		t.Parallel()

		resp := &FileListResponse{
			Object: "list",
			Data: []File{
				{
					ID:       "file-abc",
					Filename: "test.txt",
					Purpose:  PurposeBatch,
				},
			},
			HasMore: true,
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		assert.Contains(t, string(data), "list")
		assert.Contains(t, string(data), "file-abc")
		assert.Contains(t, string(data), "batch")
		assert.Contains(t, string(data), "\"has_more\":true")
	})
}

func TestFileDeleteResponse_IsDeleted(t *testing.T) {
	t.Parallel()

	t.Run("deleted", func(t *testing.T) {
		t.Parallel()

		resp := &FileDeleteResponse{
			ID:      "file-123",
			Object:  "file",
			Deleted: true,
		}

		assert.True(t, resp.IsDeleted())
	})

	t.Run("not deleted", func(t *testing.T) {
		t.Parallel()

		resp := &FileDeleteResponse{
			ID:      "file-123",
			Object:  "file",
			Deleted: false,
		}

		assert.False(t, resp.IsDeleted())
	})
}

func TestFileDeleteResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal delete response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "file-deleted",
			"object": "file",
			"deleted": true
		}`

		var resp FileDeleteResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "file-deleted", resp.ID)
		assert.Equal(t, "file", resp.Object)
		assert.True(t, resp.Deleted)
		assert.True(t, resp.IsDeleted())
	})

	t.Run("marshal delete response", func(t *testing.T) {
		t.Parallel()

		resp := &FileDeleteResponse{
			ID:      "file-xyz",
			Object:  "file",
			Deleted: true,
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		assert.Contains(t, string(data), "file-xyz")
		assert.Contains(t, string(data), "\"deleted\":true")
	})
}

func TestFileContentResponse_Helpers(t *testing.T) {
	t.Parallel()

	t.Run("GetContent", func(t *testing.T) {
		t.Parallel()

		content := []byte("file content here")
		resp := &FileContentResponse{
			Content:     content,
			ContentType: "text/plain",
		}

		assert.Equal(t, content, resp.GetContent())
	})

	t.Run("GetContentType", func(t *testing.T) {
		t.Parallel()

		resp := &FileContentResponse{
			Content:     []byte("data"),
			ContentType: "application/json",
		}

		assert.Equal(t, "application/json", resp.GetContentType())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		resp := &FileContentResponse{
			Content: []byte("test content"),
		}

		assert.Equal(t, "test content", resp.String())
	})
}

func TestFilePurpose_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, FilePurpose("assistants"), PurposeAssistants)
	assert.Equal(t, FilePurpose("fine-tune"), PurposeFineTune)
	assert.Equal(t, FilePurpose("batch"), PurposeBatch)
}

func TestFileStatus_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, FileStatus("uploaded"), StatusUploaded)
	assert.Equal(t, FileStatus("processed"), StatusProcessed)
	assert.Equal(t, FileStatus("error"), StatusError)
}

func TestFileUploadRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Create a complete upload request
	fileContent := strings.NewReader("line1\nline2\nline3")
	req := NewFileUploadRequest(fileContent, "training.jsonl", PurposeFineTune)

	// Verify the request is complete
	assert.NotNil(t, req.File)
	assert.Equal(t, "training.jsonl", req.Filename)
	assert.Equal(t, PurposeFineTune, req.Purpose)
}

func TestFileListResponse_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response
	jsonData := `{
		"object": "list",
		"data": [
			{
				"id": "file-abc123",
				"object": "file",
				"bytes": 120000,
				"created_at": 1677652288,
				"filename": "training_data.jsonl",
				"purpose": "fine-tune",
				"status": "processed"
			},
			{
				"id": "file-def456",
				"object": "file",
				"bytes": 45000,
				"created_at": 1677652300,
				"filename": "assistant_knowledge.txt",
				"purpose": "assistants",
				"status": "uploaded"
			}
		],
		"has_more": false
	}`

	var resp FileListResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "list", resp.Object)
	assert.Len(t, resp.Data, 2)
	assert.False(t, resp.HasMore)

	// Verify helper methods work
	files := resp.GetFiles()
	assert.Len(t, files, 2)

	ids := resp.GetFileIDs()
	assert.Contains(t, ids, "file-abc123")
	assert.Contains(t, ids, "file-def456")

	// Test filtering
	fineTuneFiles := resp.GetFilesByPurpose(PurposeFineTune)
	assert.Len(t, fineTuneFiles, 1)
	assert.Equal(t, "training_data.jsonl", fineTuneFiles[0].Filename)

	// Test lookup
	file := resp.GetFileByID("file-def456")
	require.NotNil(t, file)
	assert.Equal(t, "assistant_knowledge.txt", file.Filename)
	assert.True(t, file.IsUploaded())
	assert.False(t, file.HasError())
}
