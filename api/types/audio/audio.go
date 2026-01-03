// Package audio provides types for the Audio API.
package audio

import "io"

// TranscriptionModel represents the audio transcription model.
type TranscriptionModel string

const (
	// ModelWhisper1 is the Whisper-1 model for audio transcription.
	ModelWhisper1 TranscriptionModel = "whisper-1"
)

// ResponseFormat represents the format of the transcription response.
type ResponseFormat string

const (
	// ResponseFormatJSON returns the response in JSON format.
	ResponseFormatJSON ResponseFormat = "json"
	// ResponseFormatText returns the response as plain text.
	ResponseFormatText ResponseFormat = "text"
	// ResponseFormatVerboseJSON returns detailed JSON with timestamps.
	ResponseFormatVerboseJSON ResponseFormat = "verbose_json"
	// ResponseFormatVTT returns the response in WebVTT format.
	ResponseFormatVTT ResponseFormat = "vtt"
	// ResponseFormatSRT returns the response in SRT subtitle format.
	ResponseFormatSRT ResponseFormat = "srt"
)

// TranscriptionRequest represents a request to transcribe audio.
type TranscriptionRequest struct {
	// File is the audio file to transcribe.
	File io.Reader

	// Filename is the name of the audio file.
	Filename string

	// Model is the model to use for transcription (required).
	Model TranscriptionModel

	// Language is the language of the audio (optional, ISO-639-1 format).
	Language string

	// Prompt is an optional text to guide the model's style.
	Prompt string

	// ResponseFormat is the format of the transcript output.
	ResponseFormat ResponseFormat

	// Temperature is the sampling temperature (0 to 1).
	Temperature *float64
}

// NewTranscriptionRequest creates a new transcription request.
//
// Example:
//
//	file, _ := os.Open("audio.mp3")
//	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
func NewTranscriptionRequest(file io.Reader, filename string, model TranscriptionModel) *TranscriptionRequest {
	return &TranscriptionRequest{
		File:           file,
		Filename:       filename,
		Model:          model,
		ResponseFormat: ResponseFormatJSON, // Default to JSON
	}
}

// SetLanguage sets the language of the audio.
//
// Example:
//
//	req.SetLanguage("en")
func (r *TranscriptionRequest) SetLanguage(language string) *TranscriptionRequest {
	r.Language = language
	return r
}

// SetPrompt sets a prompt to guide the transcription.
//
// Example:
//
//	req.SetPrompt("This is a conversation about AI and technology.")
func (r *TranscriptionRequest) SetPrompt(prompt string) *TranscriptionRequest {
	r.Prompt = prompt
	return r
}

// SetResponseFormat sets the response format.
//
// Example:
//
//	req.SetResponseFormat(audio.ResponseFormatVerboseJSON)
func (r *TranscriptionRequest) SetResponseFormat(format ResponseFormat) *TranscriptionRequest {
	r.ResponseFormat = format
	return r
}

// SetTemperature sets the sampling temperature.
//
// Example:
//
//	req.SetTemperature(0.2)
func (r *TranscriptionRequest) SetTemperature(temp float64) *TranscriptionRequest {
	r.Temperature = &temp
	return r
}

// TranscriptionResponse represents the response from audio transcription.
type TranscriptionResponse struct {
	// Text is the transcribed text.
	Text string `json:"text"`

	// Task is the type of task performed.
	Task string `json:"task,omitempty"`

	// Language is the detected language.
	Language string `json:"language,omitempty"`

	// Duration is the audio duration in seconds.
	Duration float64 `json:"duration,omitempty"`

	// Segments contains detailed transcription segments (verbose_json only).
	Segments []TranscriptionSegment `json:"segments,omitempty"`
}

// TranscriptionSegment represents a segment of the transcription with timestamps.
type TranscriptionSegment struct {
	// ID is the segment identifier.
	ID int `json:"id"`

	// Seek is the seek position.
	Seek int `json:"seek,omitempty"`

	// Start is the start time in seconds.
	Start float64 `json:"start"`

	// End is the end time in seconds.
	End float64 `json:"end"`

	// Text is the transcribed text for this segment.
	Text string `json:"text"`

	// Tokens are the token IDs.
	Tokens []int `json:"tokens,omitempty"`

	// Temperature is the sampling temperature used.
	Temperature float64 `json:"temperature,omitempty"`

	// AvgLogprob is the average log probability.
	AvgLogprob float64 `json:"avg_logprob,omitempty"`

	// CompressionRatio is the compression ratio.
	CompressionRatio float64 `json:"compression_ratio,omitempty"`

	// NoSpeechProb is the no-speech probability.
	NoSpeechProb float64 `json:"no_speech_prob,omitempty"`
}

// TranscriptionTextResponse represents a plain text transcription response.
type TranscriptionTextResponse struct {
	// Text is the transcribed text.
	Text string
}

// GetText returns the transcribed text.
func (r *TranscriptionResponse) GetText() string {
	return r.Text
}

// GetLanguage returns the detected language.
func (r *TranscriptionResponse) GetLanguage() string {
	return r.Language
}

// GetDuration returns the audio duration in seconds.
func (r *TranscriptionResponse) GetDuration() float64 {
	return r.Duration
}

// HasSegments returns true if the response contains detailed segments.
func (r *TranscriptionResponse) HasSegments() bool {
	return len(r.Segments) > 0
}

// GetSegments returns all transcription segments.
func (r *TranscriptionResponse) GetSegments() []TranscriptionSegment {
	return r.Segments
}

// GetSegmentText returns the text from a specific segment by index.
func (r *TranscriptionResponse) GetSegmentText(index int) string {
	if index < 0 || index >= len(r.Segments) {
		return ""
	}
	return r.Segments[index].Text
}

// GetFullTranscriptFromSegments concatenates all segment texts.
func (r *TranscriptionResponse) GetFullTranscriptFromSegments() string {
	if len(r.Segments) == 0 {
		return r.Text
	}

	var result string
	for _, segment := range r.Segments {
		result += segment.Text
	}
	return result
}

// GetStartTime returns the start time of the segment in seconds.
func (s *TranscriptionSegment) GetStartTime() float64 {
	return s.Start
}

// GetEndTime returns the end time of the segment in seconds.
func (s *TranscriptionSegment) GetEndTime() float64 {
	return s.End
}

// GetDuration returns the duration of the segment in seconds.
func (s *TranscriptionSegment) GetDuration() float64 {
	return s.End - s.Start
}

// GetText returns the transcribed text for this segment.
func (s *TranscriptionSegment) GetText() string {
	return s.Text
}

// String returns the text from the response.
func (r *TranscriptionTextResponse) String() string {
	return r.Text
}
