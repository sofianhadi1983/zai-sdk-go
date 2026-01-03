package audio

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranscriptionRequest(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("audio content")
	req := NewTranscriptionRequest(file, "audio.mp3", ModelWhisper1)

	assert.NotNil(t, req.File)
	assert.Equal(t, "audio.mp3", req.Filename)
	assert.Equal(t, ModelWhisper1, req.Model)
	assert.Equal(t, ResponseFormatJSON, req.ResponseFormat) // Default
}

func TestTranscriptionRequest_Setters(t *testing.T) {
	t.Parallel()

	t.Run("SetLanguage", func(t *testing.T) {
		t.Parallel()

		req := &TranscriptionRequest{}
		req.SetLanguage("en")

		assert.Equal(t, "en", req.Language)
	})

	t.Run("SetPrompt", func(t *testing.T) {
		t.Parallel()

		req := &TranscriptionRequest{}
		req.SetPrompt("AI conversation")

		assert.Equal(t, "AI conversation", req.Prompt)
	})

	t.Run("SetResponseFormat", func(t *testing.T) {
		t.Parallel()

		req := &TranscriptionRequest{}
		req.SetResponseFormat(ResponseFormatVerboseJSON)

		assert.Equal(t, ResponseFormatVerboseJSON, req.ResponseFormat)
	})

	t.Run("SetTemperature", func(t *testing.T) {
		t.Parallel()

		req := &TranscriptionRequest{}
		req.SetTemperature(0.5)

		require.NotNil(t, req.Temperature)
		assert.Equal(t, 0.5, *req.Temperature)
	})

	t.Run("chained setters", func(t *testing.T) {
		t.Parallel()

		file := strings.NewReader("test")
		req := NewTranscriptionRequest(file, "test.mp3", ModelWhisper1)
		req.SetLanguage("es").
			SetPrompt("Spanish conversation").
			SetResponseFormat(ResponseFormatVerboseJSON).
			SetTemperature(0.2)

		assert.Equal(t, "es", req.Language)
		assert.Equal(t, "Spanish conversation", req.Prompt)
		assert.Equal(t, ResponseFormatVerboseJSON, req.ResponseFormat)
		require.NotNil(t, req.Temperature)
		assert.Equal(t, 0.2, *req.Temperature)
	})
}

func TestTranscriptionResponse_GetText(t *testing.T) {
	t.Parallel()

	resp := &TranscriptionResponse{
		Text: "Hello world",
	}

	assert.Equal(t, "Hello world", resp.GetText())
}

func TestTranscriptionResponse_GetLanguage(t *testing.T) {
	t.Parallel()

	resp := &TranscriptionResponse{
		Language: "en",
	}

	assert.Equal(t, "en", resp.GetLanguage())
}

func TestTranscriptionResponse_GetDuration(t *testing.T) {
	t.Parallel()

	resp := &TranscriptionResponse{
		Duration: 45.5,
	}

	assert.Equal(t, 45.5, resp.GetDuration())
}

func TestTranscriptionResponse_Segments(t *testing.T) {
	t.Parallel()

	t.Run("HasSegments true", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Segments: []TranscriptionSegment{
				{ID: 0, Text: "Hello"},
				{ID: 1, Text: "world"},
			},
		}

		assert.True(t, resp.HasSegments())
	})

	t.Run("HasSegments false", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Segments: []TranscriptionSegment{},
		}

		assert.False(t, resp.HasSegments())
	})

	t.Run("GetSegments", func(t *testing.T) {
		t.Parallel()

		segments := []TranscriptionSegment{
			{ID: 0, Text: "First"},
			{ID: 1, Text: "Second"},
		}
		resp := &TranscriptionResponse{
			Segments: segments,
		}

		result := resp.GetSegments()
		assert.Len(t, result, 2)
		assert.Equal(t, "First", result[0].Text)
		assert.Equal(t, "Second", result[1].Text)
	})

	t.Run("GetSegmentText valid index", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Segments: []TranscriptionSegment{
				{ID: 0, Text: "Hello"},
				{ID: 1, Text: "world"},
			},
		}

		assert.Equal(t, "Hello", resp.GetSegmentText(0))
		assert.Equal(t, "world", resp.GetSegmentText(1))
	})

	t.Run("GetSegmentText invalid index", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Segments: []TranscriptionSegment{
				{ID: 0, Text: "Hello"},
			},
		}

		assert.Empty(t, resp.GetSegmentText(-1))
		assert.Empty(t, resp.GetSegmentText(5))
	})

	t.Run("GetFullTranscriptFromSegments", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Segments: []TranscriptionSegment{
				{ID: 0, Text: "Hello "},
				{ID: 1, Text: "world"},
				{ID: 2, Text: "!"},
			},
		}

		assert.Equal(t, "Hello world!", resp.GetFullTranscriptFromSegments())
	})

	t.Run("GetFullTranscriptFromSegments no segments", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Text:     "Fallback text",
			Segments: []TranscriptionSegment{},
		}

		assert.Equal(t, "Fallback text", resp.GetFullTranscriptFromSegments())
	})
}

func TestTranscriptionSegment_Helpers(t *testing.T) {
	t.Parallel()

	segment := &TranscriptionSegment{
		ID:    0,
		Start: 10.5,
		End:   15.8,
		Text:  "Hello world",
	}

	assert.Equal(t, 10.5, segment.GetStartTime())
	assert.Equal(t, 15.8, segment.GetEndTime())
	assert.InDelta(t, 5.3, segment.GetDuration(), 0.0001) // Use delta for floating point
	assert.Equal(t, "Hello world", segment.GetText())
}

func TestTranscriptionResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal simple response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"text": "This is a test transcription.",
			"task": "transcribe",
			"language": "en",
			"duration": 30.5
		}`

		var resp TranscriptionResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "This is a test transcription.", resp.Text)
		assert.Equal(t, "transcribe", resp.Task)
		assert.Equal(t, "en", resp.Language)
		assert.Equal(t, 30.5, resp.Duration)
	})

	t.Run("unmarshal verbose response with segments", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"text": "Hello world",
			"language": "en",
			"duration": 5.0,
			"segments": [
				{
					"id": 0,
					"start": 0.0,
					"end": 2.5,
					"text": "Hello ",
					"temperature": 0.0,
					"avg_logprob": -0.5,
					"compression_ratio": 1.2,
					"no_speech_prob": 0.01
				},
				{
					"id": 1,
					"start": 2.5,
					"end": 5.0,
					"text": "world",
					"temperature": 0.0,
					"avg_logprob": -0.3,
					"compression_ratio": 1.1,
					"no_speech_prob": 0.02
				}
			]
		}`

		var resp TranscriptionResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "Hello world", resp.Text)
		assert.True(t, resp.HasSegments())
		assert.Len(t, resp.Segments, 2)

		// Check first segment
		seg1 := resp.Segments[0]
		assert.Equal(t, 0, seg1.ID)
		assert.Equal(t, 0.0, seg1.Start)
		assert.Equal(t, 2.5, seg1.End)
		assert.Equal(t, "Hello ", seg1.Text)

		// Check second segment
		seg2 := resp.Segments[1]
		assert.Equal(t, 1, seg2.ID)
		assert.Equal(t, 2.5, seg2.Start)
		assert.Equal(t, 5.0, seg2.End)
		assert.Equal(t, "world", seg2.Text)
	})

	t.Run("marshal response", func(t *testing.T) {
		t.Parallel()

		resp := &TranscriptionResponse{
			Text:     "Test transcription",
			Language: "en",
			Duration: 10.0,
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		assert.Contains(t, string(data), "Test transcription")
		assert.Contains(t, string(data), "en")
		assert.Contains(t, string(data), "10")
	})
}

func TestTranscriptionTextResponse_String(t *testing.T) {
	t.Parallel()

	resp := &TranscriptionTextResponse{
		Text: "Plain text transcription",
	}

	assert.Equal(t, "Plain text transcription", resp.String())
}

func TestTranscriptionModel_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, TranscriptionModel("whisper-1"), ModelWhisper1)
}

func TestResponseFormat_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ResponseFormat("json"), ResponseFormatJSON)
	assert.Equal(t, ResponseFormat("text"), ResponseFormatText)
	assert.Equal(t, ResponseFormat("verbose_json"), ResponseFormatVerboseJSON)
	assert.Equal(t, ResponseFormat("vtt"), ResponseFormatVTT)
	assert.Equal(t, ResponseFormat("srt"), ResponseFormatSRT)
}

func TestTranscriptionRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Create a complete request with all options
	file := strings.NewReader("audio data")
	req := NewTranscriptionRequest(file, "interview.mp3", ModelWhisper1)
	req.SetLanguage("en").
		SetPrompt("This is an interview about technology.").
		SetResponseFormat(ResponseFormatVerboseJSON).
		SetTemperature(0.2)

	// Verify all fields
	assert.NotNil(t, req.File)
	assert.Equal(t, "interview.mp3", req.Filename)
	assert.Equal(t, ModelWhisper1, req.Model)
	assert.Equal(t, "en", req.Language)
	assert.Equal(t, "This is an interview about technology.", req.Prompt)
	assert.Equal(t, ResponseFormatVerboseJSON, req.ResponseFormat)
	require.NotNil(t, req.Temperature)
	assert.Equal(t, 0.2, *req.Temperature)
}

func TestTranscriptionResponse_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response
	jsonData := `{
		"text": "In this episode, we discuss the future of artificial intelligence and its impact on society.",
		"task": "transcribe",
		"language": "en",
		"duration": 180.5,
		"segments": [
			{
				"id": 0,
				"start": 0.0,
				"end": 3.5,
				"text": "In this episode,",
				"tokens": [123, 456, 789],
				"temperature": 0.0,
				"avg_logprob": -0.35,
				"compression_ratio": 1.25,
				"no_speech_prob": 0.001
			},
			{
				"id": 1,
				"start": 3.5,
				"end": 7.2,
				"text": " we discuss the future of artificial intelligence",
				"tokens": [234, 567, 890],
				"temperature": 0.0,
				"avg_logprob": -0.28,
				"compression_ratio": 1.18,
				"no_speech_prob": 0.002
			},
			{
				"id": 2,
				"start": 7.2,
				"end": 10.5,
				"text": " and its impact on society.",
				"tokens": [345, 678, 901],
				"temperature": 0.0,
				"avg_logprob": -0.32,
				"compression_ratio": 1.22,
				"no_speech_prob": 0.001
			}
		]
	}`

	var resp TranscriptionResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	// Verify basic fields
	assert.Equal(t, "en", resp.GetLanguage())
	assert.Equal(t, 180.5, resp.GetDuration())
	assert.True(t, resp.HasSegments())

	// Verify segments
	segments := resp.GetSegments()
	require.Len(t, segments, 3)

	// Check first segment details
	firstSeg := segments[0]
	assert.Equal(t, "In this episode,", firstSeg.GetText())
	assert.Equal(t, 0.0, firstSeg.GetStartTime())
	assert.Equal(t, 3.5, firstSeg.GetEndTime())
	assert.Equal(t, 3.5, firstSeg.GetDuration())
	assert.Len(t, firstSeg.Tokens, 3)

	// Verify full transcript
	fullText := resp.GetFullTranscriptFromSegments()
	assert.Contains(t, fullText, "In this episode,")
	assert.Contains(t, fullText, "artificial intelligence")
	assert.Contains(t, fullText, "society")

	// Verify individual segment access
	assert.Equal(t, "In this episode,", resp.GetSegmentText(0))
	assert.Equal(t, " we discuss the future of artificial intelligence", resp.GetSegmentText(1))
	assert.Equal(t, " and its impact on society.", resp.GetSegmentText(2))
}
