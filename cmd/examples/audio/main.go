package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/z-ai/zai-sdk-go/api/types/audio"
	"github.com/z-ai/zai-sdk-go/pkg/zai"
)

func main() {
	// Create a new client
	client, err := zai.NewClient(
		zai.WithAPIKey("your-api-key.your-secret"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Simple transcription
	fmt.Println("=== Example 1: Simple Transcription ===")
	simpleTranscriptionExample(ctx, client)

	// Example 2: Transcription with language and prompt
	fmt.Println("\n=== Example 2: Transcription with Language & Prompt ===")
	advancedTranscriptionExample(ctx, client)

	// Example 3: Transcription with segments
	fmt.Println("\n=== Example 3: Transcription with Segments ===")
	segmentsExample(ctx, client)

	// Example 4: Different response formats
	fmt.Println("\n=== Example 4: Different Response Formats ===")
	responseFormatsExample(ctx, client)

	// Example 5: Error handling
	fmt.Println("\n=== Example 5: Error Handling ===")
	errorHandlingExample(ctx, client)
}

func simpleTranscriptionExample(ctx context.Context, client *zai.Client) {
	// Open an audio file
	// For this example, we'll use a string reader as a placeholder
	// In real usage: file, _ := os.Open("audio.mp3")
	audioContent := strings.NewReader("simulated audio content")

	// Use the convenience method for simple transcription
	text, err := client.Audio.TranscribeFile(ctx, audioContent, "audio.mp3")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Transcription: %s\n", text)
}

func advancedTranscriptionExample(ctx context.Context, client *zai.Client) {
	// Open audio file
	audioContent := strings.NewReader("simulated audio content")

	// Create a detailed transcription request
	req := audio.NewTranscriptionRequest(
		audioContent,
		"interview.mp3",
		audio.ModelWhisper1,
	)

	// Set optional parameters
	req.SetLanguage("en").                                      // Specify language (ISO-639-1)
		SetPrompt("This is an interview about technology."). // Guide the model's style
		SetResponseFormat(audio.ResponseFormatVerboseJSON).  // Get detailed response
		SetTemperature(0.2)                                  // Lower temperature for more focused output

	// Make the transcription request
	resp, err := client.Audio.Transcribe(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Display results
	fmt.Printf("Transcription: %s\n", resp.GetText())
	fmt.Printf("Language: %s\n", resp.GetLanguage())
	fmt.Printf("Duration: %.2f seconds\n", resp.GetDuration())
}

func segmentsExample(ctx context.Context, client *zai.Client) {
	// Open audio file
	audioContent := strings.NewReader("simulated podcast audio")

	// Get transcription with detailed segments
	resp, err := client.Audio.TranscribeWithSegments(
		ctx,
		audioContent,
		"podcast.mp3",
		"en", // Language (optional, empty string for auto-detection)
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Display basic info
	fmt.Printf("Full Text: %s\n", resp.GetText())
	fmt.Printf("Duration: %.2f seconds\n", resp.GetDuration())
	fmt.Printf("\nSegments:\n")

	// Display each segment with timestamps
	if resp.HasSegments() {
		for i, segment := range resp.GetSegments() {
			fmt.Printf("  [%d] %.2fs - %.2fs (%.2fs): %s\n",
				i+1,
				segment.GetStartTime(),
				segment.GetEndTime(),
				segment.GetDuration(),
				segment.GetText(),
			)
		}

		// You can also get the full transcript by concatenating segments
		fullTranscript := resp.GetFullTranscriptFromSegments()
		fmt.Printf("\nFull transcript from segments: %s\n", fullTranscript)
	}
}

func responseFormatsExample(ctx context.Context, client *zai.Client) {
	audioContent := strings.NewReader("simulated audio")

	// Example 1: JSON format (default)
	fmt.Println("1. JSON Format:")
	jsonReq := audio.NewTranscriptionRequest(audioContent, "audio.mp3", audio.ModelWhisper1)
	jsonResp, err := client.Audio.Transcribe(ctx, jsonReq)
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Text: %s\n", jsonResp.GetText())
	}

	// Example 2: Plain text format
	fmt.Println("\n2. Plain Text Format:")
	audioContent = strings.NewReader("simulated audio") // Reset reader
	textReq := audio.NewTranscriptionRequest(audioContent, "audio.mp3", audio.ModelWhisper1)
	textReq.SetResponseFormat(audio.ResponseFormatText)
	textResp, err := client.Audio.Transcribe(ctx, textReq)
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Text: %s\n", textResp.GetText())
	}

	// Example 3: Verbose JSON with detailed information
	fmt.Println("\n3. Verbose JSON Format:")
	audioContent = strings.NewReader("simulated audio") // Reset reader
	verboseReq := audio.NewTranscriptionRequest(audioContent, "audio.mp3", audio.ModelWhisper1)
	verboseReq.SetResponseFormat(audio.ResponseFormatVerboseJSON)
	verboseResp, err := client.Audio.Transcribe(ctx, verboseReq)
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Text: %s\n", verboseResp.GetText())
		fmt.Printf("   Language: %s\n", verboseResp.GetLanguage())
		fmt.Printf("   Duration: %.2fs\n", verboseResp.GetDuration())
		fmt.Printf("   Has Segments: %v\n", verboseResp.HasSegments())
	}

	// Example 4: WebVTT format (for subtitles)
	fmt.Println("\n4. WebVTT Format:")
	audioContent = strings.NewReader("simulated audio") // Reset reader
	vttReq := audio.NewTranscriptionRequest(audioContent, "audio.mp3", audio.ModelWhisper1)
	vttReq.SetResponseFormat(audio.ResponseFormatVTT)
	vttResp, err := client.Audio.Transcribe(ctx, vttReq)
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   VTT Content:\n%s\n", vttResp.GetText())
	}

	// Example 5: SRT format (for subtitles)
	fmt.Println("\n5. SRT Format:")
	audioContent = strings.NewReader("simulated audio") // Reset reader
	srtReq := audio.NewTranscriptionRequest(audioContent, "audio.mp3", audio.ModelWhisper1)
	srtReq.SetResponseFormat(audio.ResponseFormatSRT)
	srtResp, err := client.Audio.Transcribe(ctx, srtReq)
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   SRT Content:\n%s\n", srtResp.GetText())
	}
}

func errorHandlingExample(ctx context.Context, client *zai.Client) {
	// Example 1: Handle API errors
	fmt.Println("1. API Error Handling:")
	// In real usage, this might fail due to invalid file format
	invalidAudio := strings.NewReader("not a valid audio file")
	req := audio.NewTranscriptionRequest(invalidAudio, "invalid.txt", audio.ModelWhisper1)

	resp, err := client.Audio.Transcribe(ctx, req)
	if err != nil {
		fmt.Printf("   Error occurred (expected): %v\n", err)
		// Check error type and handle accordingly
		// You can check for specific error types from the errors package
	} else {
		fmt.Printf("   Unexpected success: %s\n", resp.GetText())
	}

	// Example 2: Validate response
	fmt.Println("\n2. Response Validation:")
	validAudio := strings.NewReader("simulated valid audio")
	req = audio.NewTranscriptionRequest(validAudio, "valid.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatVerboseJSON)

	resp, err = client.Audio.Transcribe(ctx, req)
	if err != nil {
		log.Printf("   Error: %v", err)
		return
	}

	// Check if we got segments as expected
	if resp.HasSegments() {
		fmt.Printf("   ✓ Received %d segments\n", len(resp.GetSegments()))
	} else {
		fmt.Println("   ✗ No segments in response")
	}

	// Safely access segments
	if text := resp.GetSegmentText(0); text != "" {
		fmt.Printf("   First segment: %s\n", text)
	}
}

// Real-world example: Transcribing an actual audio file
func transcribeRealFile(ctx context.Context, client *zai.Client) {
	// Open a real audio file from disk
	file, err := os.Open("path/to/your/audio.mp3")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	// Create transcription request
	req := audio.NewTranscriptionRequest(
		file,
		"audio.mp3",
		audio.ModelWhisper1,
	)

	// Set language if known (improves accuracy)
	req.SetLanguage("en")

	// Optionally set a prompt to guide the model
	req.SetPrompt("This is a business meeting about project planning.")

	// Get transcription with segments
	req.SetResponseFormat(audio.ResponseFormatVerboseJSON)

	// Make the request
	resp, err := client.Audio.Transcribe(ctx, req)
	if err != nil {
		log.Printf("Transcription error: %v", err)
		return
	}

	// Process the results
	fmt.Printf("Transcription: %s\n", resp.GetText())
	fmt.Printf("Language: %s\n", resp.GetLanguage())
	fmt.Printf("Duration: %.2f seconds\n", resp.GetDuration())

	// Save segments to a file or process them
	if resp.HasSegments() {
		fmt.Println("\nTimestamped segments:")
		for _, segment := range resp.GetSegments() {
			fmt.Printf("[%.2f - %.2f] %s\n",
				segment.GetStartTime(),
				segment.GetEndTime(),
				segment.GetText(),
			)
		}
	}
}

// Example: Generating subtitles from audio
func generateSubtitles(ctx context.Context, client *zai.Client, audioPath string) {
	file, err := os.Open(audioPath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	// Request SRT format
	req := audio.NewTranscriptionRequest(file, "audio.mp3", audio.ModelWhisper1)
	req.SetResponseFormat(audio.ResponseFormatSRT)

	resp, err := client.Audio.Transcribe(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Save SRT content to file
	srtPath := strings.TrimSuffix(audioPath, ".mp3") + ".srt"
	err = os.WriteFile(srtPath, []byte(resp.GetText()), 0644)
	if err != nil {
		log.Printf("Error writing SRT file: %v", err)
		return
	}

	fmt.Printf("Subtitles saved to: %s\n", srtPath)
}
