package zai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/z-ai/zai-sdk-go/api/types/audio"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// AudioService provides access to the Audio API.
type AudioService struct {
	client *client.BaseClient
}

// newAudioService creates a new audio service.
func newAudioService(baseClient *client.BaseClient) *AudioService {
	return &AudioService{
		client: baseClient,
	}
}

// Transcribe transcribes audio to text.
//
// Example:
//
//	file, err := os.Open("audio.mp3")
//	if err != nil {
//	    // Handle error
//	}
//	defer file.Close()
//
//	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
//	req.SetLanguage("en").SetResponseFormat(audio.ResponseFormatVerboseJSON)
//
//	resp, err := client.Audio.Transcribe(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Transcription: %s\n", resp.GetText())
//	fmt.Printf("Language: %s\n", resp.GetLanguage())
//	fmt.Printf("Duration: %.2f seconds\n", resp.GetDuration())
func (s *AudioService) Transcribe(ctx context.Context, req *audio.TranscriptionRequest) (*audio.TranscriptionResponse, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the model field
	if err := writer.WriteField("model", string(req.Model)); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// Add optional fields
	if req.Language != "" {
		if err := writer.WriteField("language", req.Language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %w", err)
		}
	}

	if req.Prompt != "" {
		if err := writer.WriteField("prompt", req.Prompt); err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %w", err)
		}
	}

	if req.ResponseFormat != "" {
		if err := writer.WriteField("response_format", string(req.ResponseFormat)); err != nil {
			return nil, fmt.Errorf("failed to write response_format field: %w", err)
		}
	}

	if req.Temperature != nil {
		if err := writer.WriteField("temperature", fmt.Sprintf("%f", *req.Temperature)); err != nil {
			return nil, fmt.Errorf("failed to write temperature field: %w", err)
		}
	}

	// Add the audio file
	part, err := writer.CreateFormFile("file", req.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to the form
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Make the API request using PostMultipart
	apiResp, err := s.client.PostMultipart(ctx, "/audio/transcriptions", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	// Handle different response formats
	if req.ResponseFormat == audio.ResponseFormatText ||
		req.ResponseFormat == audio.ResponseFormatVTT ||
		req.ResponseFormat == audio.ResponseFormatSRT {
		// For text-based formats, read as plain text
		content, err := io.ReadAll(apiResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return &audio.TranscriptionResponse{
			Text: string(content),
		}, nil
	}

	// For JSON formats, parse as JSON
	var resp audio.TranscriptionResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// TranscribeFile is a convenience method to transcribe an audio file with default settings.
// Returns the transcribed text.
//
// Example:
//
//	file, _ := os.Open("interview.mp3")
//	text, err := client.Audio.TranscribeFile(ctx, file, "interview.mp3")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Transcription: %s\n", text)
func (s *AudioService) TranscribeFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	req := audio.NewTranscriptionRequest(file, filename, audio.ModelWhisper1)

	resp, err := s.Transcribe(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetText(), nil
}

// TranscribeWithSegments transcribes audio and returns detailed segments with timestamps.
// This is useful for generating subtitles or analyzing speech patterns.
//
// Example:
//
//	file, _ := os.Open("podcast.mp3")
//	resp, err := client.Audio.TranscribeWithSegments(ctx, file, "podcast.mp3", "en")
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, segment := range resp.GetSegments() {
//	    fmt.Printf("[%.2f - %.2f] %s\n",
//	        segment.GetStartTime(),
//	        segment.GetEndTime(),
//	        segment.GetText())
//	}
func (s *AudioService) TranscribeWithSegments(ctx context.Context, file io.Reader, filename, language string) (*audio.TranscriptionResponse, error) {
	req := audio.NewTranscriptionRequest(file, filename, audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatVerboseJSON)

	if language != "" {
		req.SetLanguage(language)
	}

	return s.Transcribe(ctx, req)
}
