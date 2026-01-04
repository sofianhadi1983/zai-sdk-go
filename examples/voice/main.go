package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/voice"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
	// Create client from environment variables
	client, err := zai.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Clone a voice
	fmt.Println("=== Voice Clone Example ===")
	cloneExample(ctx, client)

	fmt.Println()

	// Example 2: List voices
	fmt.Println("=== Voice List Example ===")
	listExample(ctx, client)

	fmt.Println()

	// Example 3: List voices with filters
	fmt.Println("=== Voice List with Filters Example ===")
	listWithFiltersExample(ctx, client)

	fmt.Println()

	// Example 4: Delete a voice
	fmt.Println("=== Voice Delete Example ===")
	deleteExample(ctx, client)
}

func cloneExample(ctx context.Context, client *zai.Client) {
	// Note: You need to upload an audio file first using the Files API
	// and get the file_id before cloning a voice
	fileID := os.Getenv("VOICE_FILE_ID")
	if fileID == "" {
		fmt.Println("Skipping clone example: VOICE_FILE_ID not set")
		return
	}

	req := voice.NewVoiceCloneRequest(
		"my_custom_voice",              // voice name
		"Hello, this is a test.",       // text corresponding to the audio
		"Welcome to the voice clone.",  // preview text
		fileID,                         // file ID of the uploaded audio
		"voice-clone-v1",               // model
	).SetRequestID("clone_req_123")

	resp, err := client.Voice.Clone(ctx, req)
	if err != nil {
		log.Fatalf("Failed to clone voice: %v", err)
	}

	fmt.Printf("Voice cloned successfully!\n")
	fmt.Printf("Voice ID: %s\n", resp.Voice)
	fmt.Printf("Preview File ID: %s\n", resp.FileID)
	fmt.Printf("File Purpose: %s\n", resp.FilePurpose)
}

func listExample(ctx context.Context, client *zai.Client) {
	req := voice.NewVoiceListRequest()

	resp, err := client.Voice.List(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list voices: %v", err)
	}

	voices := resp.GetVoices()
	fmt.Printf("Found %d voice(s):\n", len(voices))

	for i, v := range voices {
		fmt.Printf("\n%d. Voice: %s\n", i+1, v.Voice)
		fmt.Printf("   Name: %s\n", v.VoiceName)
		fmt.Printf("   Type: %s\n", v.VoiceType)
		fmt.Printf("   Created: %s\n", v.CreateTime)
		fmt.Printf("   Download URL: %s\n", v.DownloadURL)
	}
}

func listWithFiltersExample(ctx context.Context, client *zai.Client) {
	// List only cloned voices
	req := voice.NewVoiceListRequest().
		SetVoiceType("cloned").
		SetRequestID("list_req_456")

	resp, err := client.Voice.List(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list cloned voices: %v", err)
	}

	voices := resp.GetVoices()
	fmt.Printf("Found %d cloned voice(s):\n", len(voices))

	for i, v := range voices {
		fmt.Printf("\n%d. Voice: %s (%s)\n", i+1, v.VoiceName, v.Voice)
		fmt.Printf("   Type: %s\n", v.VoiceType)
		fmt.Printf("   Created: %s\n", v.CreateTime)
	}
}

func deleteExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual voice ID to delete
	voiceID := os.Getenv("VOICE_ID_TO_DELETE")
	if voiceID == "" {
		fmt.Println("Skipping delete example: VOICE_ID_TO_DELETE not set")
		return
	}

	req := voice.NewVoiceDeleteRequest(voiceID).
		SetRequestID("delete_req_789")

	resp, err := client.Voice.Delete(ctx, req)
	if err != nil {
		log.Fatalf("Failed to delete voice: %v", err)
	}

	fmt.Printf("Voice deleted successfully!\n")
	fmt.Printf("Voice ID: %s\n", resp.Voice)
	fmt.Printf("Deleted at: %s\n", resp.UpdateTime)
}
