package voice

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVoiceCloneRequest(t *testing.T) {
	t.Parallel()

	voiceName := "my_voice"
	text := "Hello world"
	input := "Preview text"
	fileID := "file_123"
	model := "voice-model-v1"

	req := NewVoiceCloneRequest(voiceName, text, input, fileID, model)

	assert.NotNil(t, req)
	assert.Equal(t, voiceName, req.VoiceName)
	assert.Equal(t, text, req.Text)
	assert.Equal(t, input, req.Input)
	assert.Equal(t, fileID, req.FileID)
	assert.Equal(t, model, req.Model)
	assert.Empty(t, req.RequestID)
}

func TestVoiceCloneRequest_SetRequestID(t *testing.T) {
	t.Parallel()

	req := NewVoiceCloneRequest("voice", "text", "input", "file_123", "model")
	req.SetRequestID("req_456")

	assert.Equal(t, "req_456", req.RequestID)

	// Test method chaining
	req2 := NewVoiceCloneRequest("voice", "text", "input", "file_123", "model").
		SetRequestID("req_789")

	assert.Equal(t, "req_789", req2.RequestID)
}

func TestVoiceCloneRequest_JSON(t *testing.T) {
	t.Parallel()

	req := NewVoiceCloneRequest("my_voice", "text", "input", "file_123", "model").
		SetRequestID("req_456")

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded VoiceCloneRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.VoiceName, decoded.VoiceName)
	assert.Equal(t, req.Text, decoded.Text)
	assert.Equal(t, req.Input, decoded.Input)
	assert.Equal(t, req.FileID, decoded.FileID)
	assert.Equal(t, req.Model, decoded.Model)
	assert.Equal(t, req.RequestID, decoded.RequestID)
}

func TestVoiceCloneResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := VoiceCloneResponse{
		Voice:       "voice_123",
		FileID:      "file_456",
		FilePurpose: "preview",
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded VoiceCloneResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp.Voice, decoded.Voice)
	assert.Equal(t, resp.FileID, decoded.FileID)
	assert.Equal(t, resp.FilePurpose, decoded.FilePurpose)
}

func TestNewVoiceDeleteRequest(t *testing.T) {
	t.Parallel()

	voice := "voice_123"

	req := NewVoiceDeleteRequest(voice)

	assert.NotNil(t, req)
	assert.Equal(t, voice, req.Voice)
	assert.Empty(t, req.RequestID)
}

func TestVoiceDeleteRequest_SetRequestID(t *testing.T) {
	t.Parallel()

	req := NewVoiceDeleteRequest("voice_123")
	req.SetRequestID("req_456")

	assert.Equal(t, "req_456", req.RequestID)

	// Test method chaining
	req2 := NewVoiceDeleteRequest("voice_123").
		SetRequestID("req_789")

	assert.Equal(t, "req_789", req2.RequestID)
}

func TestVoiceDeleteRequest_JSON(t *testing.T) {
	t.Parallel()

	req := NewVoiceDeleteRequest("voice_123").
		SetRequestID("req_456")

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded VoiceDeleteRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.Voice, decoded.Voice)
	assert.Equal(t, req.RequestID, decoded.RequestID)
}

func TestVoiceDeleteResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := VoiceDeleteResponse{
		Voice:      "voice_123",
		UpdateTime: "2024-01-01 12:00:00",
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded VoiceDeleteResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp.Voice, decoded.Voice)
	assert.Equal(t, resp.UpdateTime, decoded.UpdateTime)
}

func TestVoiceData_JSON(t *testing.T) {
	t.Parallel()

	data := VoiceData{
		Voice:       "voice_123",
		VoiceName:   "My Voice",
		VoiceType:   "cloned",
		DownloadURL: "https://example.com/voice.mp3",
		CreateTime:  "2024-01-01 12:00:00",
	}

	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	var decoded VoiceData
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, data.Voice, decoded.Voice)
	assert.Equal(t, data.VoiceName, decoded.VoiceName)
	assert.Equal(t, data.VoiceType, decoded.VoiceType)
	assert.Equal(t, data.DownloadURL, decoded.DownloadURL)
	assert.Equal(t, data.CreateTime, decoded.CreateTime)
}

func TestNewVoiceListRequest(t *testing.T) {
	t.Parallel()

	req := NewVoiceListRequest()

	assert.NotNil(t, req)
	assert.Empty(t, req.VoiceType)
	assert.Empty(t, req.VoiceName)
	assert.Empty(t, req.RequestID)
}

func TestVoiceListRequest_SetMethods(t *testing.T) {
	t.Parallel()

	req := NewVoiceListRequest().
		SetVoiceType("cloned").
		SetVoiceName("My Voice").
		SetRequestID("req_123")

	assert.Equal(t, "cloned", req.VoiceType)
	assert.Equal(t, "My Voice", req.VoiceName)
	assert.Equal(t, "req_123", req.RequestID)
}

func TestVoiceListRequest_JSON(t *testing.T) {
	t.Parallel()

	req := NewVoiceListRequest().
		SetVoiceType("preset").
		SetVoiceName("Test Voice")

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var decoded VoiceListRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, req.VoiceType, decoded.VoiceType)
	assert.Equal(t, req.VoiceName, decoded.VoiceName)
}

func TestVoiceListResponse_GetVoices(t *testing.T) {
	t.Parallel()

	t.Run("with voices", func(t *testing.T) {
		resp := VoiceListResponse{
			VoiceList: []VoiceData{
				{
					Voice:       "voice_1",
					VoiceName:   "Voice 1",
					VoiceType:   "cloned",
					DownloadURL: "https://example.com/voice1.mp3",
					CreateTime:  "2024-01-01 12:00:00",
				},
				{
					Voice:       "voice_2",
					VoiceName:   "Voice 2",
					VoiceType:   "preset",
					DownloadURL: "https://example.com/voice2.mp3",
					CreateTime:  "2024-01-02 12:00:00",
				},
			},
		}

		voices := resp.GetVoices()
		assert.Len(t, voices, 2)
		assert.Equal(t, "voice_1", voices[0].Voice)
		assert.Equal(t, "voice_2", voices[1].Voice)
	})

	t.Run("with nil list", func(t *testing.T) {
		resp := VoiceListResponse{
			VoiceList: nil,
		}

		voices := resp.GetVoices()
		assert.NotNil(t, voices)
		assert.Len(t, voices, 0)
	})

	t.Run("with empty list", func(t *testing.T) {
		resp := VoiceListResponse{
			VoiceList: []VoiceData{},
		}

		voices := resp.GetVoices()
		assert.NotNil(t, voices)
		assert.Len(t, voices, 0)
	})
}

func TestVoiceListResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := VoiceListResponse{
		VoiceList: []VoiceData{
			{
				Voice:       "voice_1",
				VoiceName:   "Voice 1",
				VoiceType:   "cloned",
				DownloadURL: "https://example.com/voice1.mp3",
				CreateTime:  "2024-01-01 12:00:00",
			},
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded VoiceListResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded.VoiceList, 1)
	assert.Equal(t, resp.VoiceList[0].Voice, decoded.VoiceList[0].Voice)
	assert.Equal(t, resp.VoiceList[0].VoiceName, decoded.VoiceList[0].VoiceName)
}
