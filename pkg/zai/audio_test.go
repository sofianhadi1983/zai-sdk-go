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
	"github.com/z-ai/zai-sdk-go/api/types/audio"
)

func TestAudioService_Transcribe_JSON(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/audio/transcriptions", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		require.NoError(t, err)

		// Verify form fields
		assert.Equal(t, "whisper-1", r.FormValue("model"))
		assert.Equal(t, "en", r.FormValue("language"))
		assert.Equal(t, "AI conversation", r.FormValue("prompt"))
		assert.Equal(t, "json", r.FormValue("response_format"))
		assert.Equal(t, "0.200000", r.FormValue("temperature"))

		// Verify file
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "audio.mp3", header.Filename)

		// Read file content
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "audio content", string(content))

		// Send JSON response
		resp := audio.TranscriptionResponse{
			Text:     "Hello world",
			Task:     "transcribe",
			Language: "en",
			Duration: 5.0,
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Create request
	file := strings.NewReader("audio content")
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetLanguage("en").
		SetPrompt("AI conversation").
		SetResponseFormat(audio.ResponseFormatJSON).
		SetTemperature(0.2)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "Hello world", resp.GetText())
	assert.Equal(t, "en", resp.GetLanguage())
	assert.Equal(t, 5.0, resp.GetDuration())
}

func TestAudioService_Transcribe_VerboseJSON(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, http.MethodPost, r.Method)

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		assert.Equal(t, "verbose_json", r.FormValue("response_format"))

		// Send verbose JSON response with segments
		resp := audio.TranscriptionResponse{
			Text:     "Hello world",
			Language: "en",
			Duration: 5.0,
			Segments: []audio.TranscriptionSegment{
				{
					ID:    0,
					Start: 0.0,
					End:   2.5,
					Text:  "Hello ",
				},
				{
					ID:    1,
					Start: 2.5,
					End:   5.0,
					Text:  "world",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Create request
	file := strings.NewReader("audio content")
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatVerboseJSON)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "Hello world", resp.GetText())
	assert.True(t, resp.HasSegments())
	assert.Len(t, resp.GetSegments(), 2)

	// Verify segments
	segments := resp.GetSegments()
	assert.Equal(t, "Hello ", segments[0].GetText())
	assert.Equal(t, 0.0, segments[0].GetStartTime())
	assert.Equal(t, 2.5, segments[0].GetEndTime())

	assert.Equal(t, "world", segments[1].GetText())
	assert.Equal(t, 2.5, segments[1].GetStartTime())
	assert.Equal(t, 5.0, segments[1].GetEndTime())
}

func TestAudioService_Transcribe_Text(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		assert.Equal(t, "text", r.FormValue("response_format"))

		// Send plain text response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("This is plain text transcription."))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	file := strings.NewReader("audio content")
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatText)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "This is plain text transcription.", resp.GetText())
}

func TestAudioService_Transcribe_VTT(t *testing.T) {
	t.Parallel()

	vttContent := `WEBVTT

00:00:00.000 --> 00:00:02.500
Hello

00:00:02.500 --> 00:00:05.000
world`

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		assert.Equal(t, "vtt", r.FormValue("response_format"))

		// Send VTT response
		w.Header().Set("Content-Type", "text/vtt")
		w.Write([]byte(vttContent))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	file := strings.NewReader("audio content")
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatVTT)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, vttContent, resp.GetText())
}

func TestAudioService_Transcribe_SRT(t *testing.T) {
	t.Parallel()

	srtContent := `1
00:00:00,000 --> 00:00:02,500
Hello

2
00:00:02,500 --> 00:00:05,000
world`

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		assert.Equal(t, "srt", r.FormValue("response_format"))

		// Send SRT response
		w.Header().Set("Content-Type", "application/x-subrip")
		w.Write([]byte(srtContent))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	file := strings.NewReader("audio content")
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatSRT)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, srtContent, resp.GetText())
}

func TestAudioService_TranscribeFile(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify default settings
		assert.Equal(t, "whisper-1", r.FormValue("model"))
		assert.Equal(t, "json", r.FormValue("response_format"))

		// Verify file
		file, header, err := r.FormFile("file")
		require.NoError(t, err)
		defer file.Close()
		assert.Equal(t, "interview.mp3", header.Filename)

		// Send JSON response
		resp := audio.TranscriptionResponse{
			Text: "This is a transcription of the interview.",
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Make request
	file := strings.NewReader("interview audio")
	text, err := client.Audio.TranscribeFile(context.Background(), file, "interview.mp3")
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "This is a transcription of the interview.", text)
}

func TestAudioService_TranscribeWithSegments(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify settings
		assert.Equal(t, "whisper-1", r.FormValue("model"))
		assert.Equal(t, "verbose_json", r.FormValue("response_format"))
		assert.Equal(t, "en", r.FormValue("language"))

		// Send verbose JSON response with segments
		resp := audio.TranscriptionResponse{
			Text:     "In this episode, we discuss AI.",
			Language: "en",
			Duration: 10.5,
			Segments: []audio.TranscriptionSegment{
				{
					ID:    0,
					Start: 0.0,
					End:   3.5,
					Text:  "In this episode,",
				},
				{
					ID:    1,
					Start: 3.5,
					End:   7.2,
					Text:  " we discuss",
				},
				{
					ID:    2,
					Start: 7.2,
					End:   10.5,
					Text:  " AI.",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Make request
	file := strings.NewReader("podcast audio")
	resp, err := client.Audio.TranscribeWithSegments(
		context.Background(),
		file,
		"podcast.mp3",
		"en",
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify response
	assert.Equal(t, "In this episode, we discuss AI.", resp.GetText())
	assert.Equal(t, "en", resp.GetLanguage())
	assert.Equal(t, 10.5, resp.GetDuration())
	assert.True(t, resp.HasSegments())
	assert.Len(t, resp.GetSegments(), 3)

	// Verify segments
	assert.Equal(t, "In this episode,", resp.GetSegmentText(0))
	assert.Equal(t, " we discuss", resp.GetSegmentText(1))
	assert.Equal(t, " AI.", resp.GetSegmentText(2))

	// Verify full transcript
	fullText := resp.GetFullTranscriptFromSegments()
	assert.Equal(t, "In this episode, we discuss AI.", fullText)
}

func TestAudioService_TranscribeWithSegments_NoLanguage(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify language is not set
		assert.Empty(t, r.FormValue("language"))

		// Send response
		resp := audio.TranscriptionResponse{
			Text: "Auto-detected language",
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Make request without language
	file := strings.NewReader("audio")
	resp, err := client.Audio.TranscribeWithSegments(
		context.Background(),
		file,
		"test.mp3",
		"", // Empty language
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "Auto-detected language", resp.GetText())
}

func TestAudioService_Transcribe_APIError(t *testing.T) {
	t.Parallel()

	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "Invalid audio file", "code": "invalid_request"}}`))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)
	defer client.Close()

	// Create request
	file := strings.NewReader("invalid audio")
	req := audio.NewTranscriptionRequest(file, "invalid.mp3", audio.ModelWhisper1)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAudioService_Transcribe_MinimalRequest(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify only required fields are set
		assert.Equal(t, "whisper-1", r.FormValue("model"))
		assert.Equal(t, "json", r.FormValue("response_format")) // Default
		assert.Empty(t, r.FormValue("language"))
		assert.Empty(t, r.FormValue("prompt"))
		assert.Empty(t, r.FormValue("temperature"))

		// Send response
		resp := audio.TranscriptionResponse{
			Text: "Minimal transcription",
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Create minimal request
	file := strings.NewReader("audio")
	req := audio.NewTranscriptionRequest(file, "test.mp3", audio.ModelWhisper1)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "Minimal transcription", resp.GetText())
}

func TestAudioService_Transcribe_AllOptionalFields(t *testing.T) {
	t.Parallel()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		require.NoError(t, err)

		// Verify all optional fields are set
		assert.Equal(t, "whisper-1", r.FormValue("model"))
		assert.Equal(t, "es", r.FormValue("language"))
		assert.Equal(t, "Spanish conversation about technology", r.FormValue("prompt"))
		assert.Equal(t, "verbose_json", r.FormValue("response_format"))
		assert.Equal(t, "0.300000", r.FormValue("temperature"))

		// Send response
		resp := audio.TranscriptionResponse{
			Text:     "Conversación en español",
			Language: "es",
			Duration: 15.0,
		}

		w.Header().Set("Content-Type", "application/json")
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

	// Create request with all optional fields
	file := strings.NewReader("spanish audio")
	req := audio.NewTranscriptionRequest(file, "spanish.mp3", audio.ModelWhisper1)
	req.SetLanguage("es").
		SetPrompt("Spanish conversation about technology").
		SetResponseFormat(audio.ResponseFormatVerboseJSON).
		SetTemperature(0.3)

	// Make request
	resp, err := client.Audio.Transcribe(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "Conversación en español", resp.GetText())
	assert.Equal(t, "es", resp.GetLanguage())
	assert.Equal(t, 15.0, resp.GetDuration())
}

func TestAudioService_Transcribe_ContextCancellation(t *testing.T) {
	t.Parallel()

	// Mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will never be reached due to context cancellation
		resp := audio.TranscriptionResponse{Text: "Should not reach here"}
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

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create request
	file := strings.NewReader("audio")
	req := audio.NewTranscriptionRequest(file, "test.mp3", audio.ModelWhisper1)

	// Make request with cancelled context
	resp, err := client.Audio.Transcribe(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestNewAudioService(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithAPIKey("test-key.test-secret"))
	require.NoError(t, err)
	defer client.Close()

	// Verify Audio service is initialized
	assert.NotNil(t, client.Audio)
}
