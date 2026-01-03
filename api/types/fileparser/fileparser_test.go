package fileparser

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolType_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ToolType("lite"), ToolTypeLite)
	assert.Equal(t, ToolType("expert"), ToolTypeExpert)
	assert.Equal(t, ToolType("prime"), ToolTypePrime)
	assert.Equal(t, ToolType("prime-sync"), ToolTypePrimeSync)
}

func TestFormatType_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, FormatType("text"), FormatTypeText)
	assert.Equal(t, FormatType("download_link"), FormatTypeDownloadLink)
}

func TestNewCreateRequest(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("test pdf data")
	fileName := "test.pdf"
	fileType := "pdf"

	req := NewCreateRequest(file, fileName, fileType, ToolTypePrime)

	assert.NotNil(t, req)
	assert.Equal(t, file, req.File)
	assert.Equal(t, fileName, req.FileName)
	assert.Equal(t, fileType, req.FileType)
	assert.Equal(t, ToolTypePrime, req.ToolType)
}

func TestCreateResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := CreateResponse{
		TaskID:  "task_123",
		Message: "Task created successfully",
		Success: true,
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded CreateResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp.TaskID, decoded.TaskID)
	assert.Equal(t, resp.Message, decoded.Message)
	assert.Equal(t, resp.Success, decoded.Success)
}

func TestNewContentRequest(t *testing.T) {
	t.Parallel()

	taskID := "task_123"
	formatType := FormatTypeText

	req := NewContentRequest(taskID, formatType)

	assert.NotNil(t, req)
	assert.Equal(t, taskID, req.TaskID)
	assert.Equal(t, formatType, req.FormatType)
}

func TestContentResponse_Methods(t *testing.T) {
	t.Parallel()

	t.Run("with text content", func(t *testing.T) {
		resp := ContentResponse{
			Content: "This is parsed text content",
		}

		assert.Equal(t, "This is parsed text content", resp.GetContent())
		assert.True(t, resp.HasContent())
		assert.False(t, resp.HasData())
		assert.Nil(t, resp.GetData())
	})

	t.Run("with binary data", func(t *testing.T) {
		resp := ContentResponse{
			Data: []byte("binary data"),
		}

		assert.Equal(t, []byte("binary data"), resp.GetData())
		assert.True(t, resp.HasData())
		assert.False(t, resp.HasContent())
		assert.Equal(t, "", resp.GetContent())
	})

	t.Run("empty response", func(t *testing.T) {
		resp := ContentResponse{}

		assert.False(t, resp.HasContent())
		assert.False(t, resp.HasData())
		assert.Equal(t, "", resp.GetContent())
		assert.Nil(t, resp.GetData())
	})
}

func TestNewSyncRequest(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("test docx data")
	fileName := "test.docx"
	fileType := "docx"

	req := NewSyncRequest(file, fileName, fileType)

	assert.NotNil(t, req)
	assert.Equal(t, file, req.File)
	assert.Equal(t, fileName, req.FileName)
	assert.Equal(t, fileType, req.FileType)
	assert.Equal(t, ToolTypePrimeSync, req.ToolType)
}

func TestSyncResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := SyncResponse{
		TaskID:           "task_456",
		Message:          "Parsing completed",
		Status:           true,
		Content:          "Parsed document content",
		ParsingResultURL: "https://example.com/result.txt",
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded SyncResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp.TaskID, decoded.TaskID)
	assert.Equal(t, resp.Message, decoded.Message)
	assert.Equal(t, resp.Status, decoded.Status)
	assert.Equal(t, resp.Content, decoded.Content)
	assert.Equal(t, resp.ParsingResultURL, decoded.ParsingResultURL)
}

func TestSyncResponse_Methods(t *testing.T) {
	t.Parallel()

	t.Run("with content", func(t *testing.T) {
		resp := SyncResponse{
			TaskID:           "task_789",
			Message:          "Success",
			Status:           true,
			Content:          "Document text content",
			ParsingResultURL: "https://example.com/download",
		}

		assert.Equal(t, "Document text content", resp.GetContent())
		assert.Equal(t, "https://example.com/download", resp.GetDownloadURL())
		assert.True(t, resp.HasContent())
	})

	t.Run("without content", func(t *testing.T) {
		resp := SyncResponse{
			TaskID:           "task_empty",
			Message:          "No content",
			Status:           false,
			Content:          "",
			ParsingResultURL: "",
		}

		assert.Equal(t, "", resp.GetContent())
		assert.Equal(t, "", resp.GetDownloadURL())
		assert.False(t, resp.HasContent())
	})
}

func TestToolTypeValues(t *testing.T) {
	t.Parallel()

	// Test that tool types have the expected string values
	tests := []struct {
		toolType ToolType
		expected string
	}{
		{ToolTypeLite, "lite"},
		{ToolTypeExpert, "expert"},
		{ToolTypePrime, "prime"},
		{ToolTypePrimeSync, "prime-sync"},
	}

	for _, tt := range tests {
		t.Run(string(tt.toolType), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.toolType))
		})
	}
}

func TestFormatTypeValues(t *testing.T) {
	t.Parallel()

	// Test that format types have the expected string values
	tests := []struct {
		formatType FormatType
		expected   string
	}{
		{FormatTypeText, "text"},
		{FormatTypeDownloadLink, "download_link"},
	}

	for _, tt := range tests {
		t.Run(string(tt.formatType), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.formatType))
		})
	}
}
