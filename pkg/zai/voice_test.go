package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/voice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVoiceService_Clone(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/voice/clone", r.URL.Path)

		// Verify request body
		var req voice.VoiceCloneRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "my_voice", req.VoiceName)
		assert.Equal(t, "Hello world", req.Text)
		assert.Equal(t, "Preview text", req.Input)
		assert.Equal(t, "file_123", req.FileID)
		assert.Equal(t, "voice-model-v1", req.Model)

		resp := voice.VoiceCloneResponse{
			Voice:       "voice_456",
			FileID:      "file_789",
			FilePurpose: "preview",
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

	req := voice.NewVoiceCloneRequest(
		"my_voice",
		"Hello world",
		"Preview text",
		"file_123",
		"voice-model-v1",
	)

	resp, err := client.Voice.Clone(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "voice_456", resp.Voice)
	assert.Equal(t, "file_789", resp.FileID)
	assert.Equal(t, "preview", resp.FilePurpose)
}

func TestVoiceService_Clone_WithRequestID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req voice.VoiceCloneRequest
		json.NewDecoder(r.Body).Decode(&req)

		assert.Equal(t, "req_123", req.RequestID)

		resp := voice.VoiceCloneResponse{
			Voice:       "voice_456",
			FileID:      "file_789",
			FilePurpose: "preview",
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

	req := voice.NewVoiceCloneRequest(
		"my_voice",
		"text",
		"input",
		"file_123",
		"model",
	).SetRequestID("req_123")

	resp, err := client.Voice.Clone(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestVoiceService_Delete(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/voice/delete", r.URL.Path)

		// Verify request body
		var req voice.VoiceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		assert.Equal(t, "voice_123", req.Voice)

		resp := voice.VoiceDeleteResponse{
			Voice:      "voice_123",
			UpdateTime: "2024-01-01 12:00:00",
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

	req := voice.NewVoiceDeleteRequest("voice_123")

	resp, err := client.Voice.Delete(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "voice_123", resp.Voice)
	assert.Equal(t, "2024-01-01 12:00:00", resp.UpdateTime)
}

func TestVoiceService_Delete_WithRequestID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req voice.VoiceDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)

		assert.Equal(t, "req_456", req.RequestID)

		resp := voice.VoiceDeleteResponse{
			Voice:      "voice_123",
			UpdateTime: "2024-01-01 12:00:00",
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

	req := voice.NewVoiceDeleteRequest("voice_123").SetRequestID("req_456")

	resp, err := client.Voice.Delete(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestVoiceService_List(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/voice/list", r.URL.Path)

		// Verify no query parameters for listing all
		assert.Empty(t, r.URL.Query().Get("voiceType"))
		assert.Empty(t, r.URL.Query().Get("voiceName"))

		resp := voice.VoiceListResponse{
			VoiceList: []voice.VoiceData{
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

	req := voice.NewVoiceListRequest()

	resp, err := client.Voice.List(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	voices := resp.GetVoices()
	assert.Len(t, voices, 2)
	assert.Equal(t, "voice_1", voices[0].Voice)
	assert.Equal(t, "Voice 1", voices[0].VoiceName)
	assert.Equal(t, "cloned", voices[0].VoiceType)
	assert.Equal(t, "voice_2", voices[1].Voice)
	assert.Equal(t, "preset", voices[1].VoiceType)
}

func TestVoiceService_List_WithFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/voice/list", r.URL.Path)

		// Verify query parameters
		assert.Equal(t, "cloned", r.URL.Query().Get("voiceType"))
		assert.Equal(t, "My Voice", r.URL.Query().Get("voiceName"))
		assert.Equal(t, "req_789", r.URL.Query().Get("request_id"))

		resp := voice.VoiceListResponse{
			VoiceList: []voice.VoiceData{
				{
					Voice:       "voice_1",
					VoiceName:   "My Voice",
					VoiceType:   "cloned",
					DownloadURL: "https://example.com/voice1.mp3",
					CreateTime:  "2024-01-01 12:00:00",
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

	req := voice.NewVoiceListRequest().
		SetVoiceType("cloned").
		SetVoiceName("My Voice").
		SetRequestID("req_789")

	resp, err := client.Voice.List(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	voices := resp.GetVoices()
	assert.Len(t, voices, 1)
	assert.Equal(t, "My Voice", voices[0].VoiceName)
	assert.Equal(t, "cloned", voices[0].VoiceType)
}

func TestVoiceService_List_EmptyResult(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := voice.VoiceListResponse{
			VoiceList: []voice.VoiceData{},
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

	req := voice.NewVoiceListRequest()

	resp, err := client.Voice.List(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	voices := resp.GetVoices()
	assert.NotNil(t, voices)
	assert.Len(t, voices, 0)
}
