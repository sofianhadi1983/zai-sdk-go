package zai

import (
	"context"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/voice"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
)

// VoiceService provides access to the Voice API.
type VoiceService struct {
	client *client.BaseClient
}

// newVoiceService creates a new voice service.
func newVoiceService(baseClient *client.BaseClient) *VoiceService {
	return &VoiceService{
		client: baseClient,
	}
}

// Clone clones a voice with the provided audio sample and parameters.
//
// Example:
//
//	req := voice.NewVoiceCloneRequest(
//	    "my_voice",
//	    "Hello world",
//	    "Preview text",
//	    "file_123",
//	    "voice-model-v1",
//	).SetRequestID("req_123")
//
//	resp, err := client.Voice.Clone(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Cloned voice: %s\n", resp.Voice)
func (s *VoiceService) Clone(ctx context.Context, req *voice.VoiceCloneRequest) (*voice.VoiceCloneResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/voice/clone", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp voice.VoiceCloneResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete deletes a cloned voice by voice ID.
//
// Example:
//
//	req := voice.NewVoiceDeleteRequest("voice_123").
//	    SetRequestID("req_456")
//
//	resp, err := client.Voice.Delete(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Deleted voice: %s at %s\n", resp.Voice, resp.UpdateTime)
func (s *VoiceService) Delete(ctx context.Context, req *voice.VoiceDeleteRequest) (*voice.VoiceDeleteResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/voice/delete", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp voice.VoiceDeleteResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// List lists voices with optional filtering.
//
// Example:
//
//	// List all voices
//	req := voice.NewVoiceListRequest()
//
//	resp, err := client.Voice.List(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	for _, v := range resp.GetVoices() {
//	    fmt.Printf("Voice: %s (%s)\n", v.VoiceName, v.VoiceType)
//	}
//
// Example with filters:
//
//	// List only cloned voices
//	req := voice.NewVoiceListRequest().
//	    SetVoiceType("cloned").
//	    SetRequestID("req_789")
//
//	resp, err := client.Voice.List(ctx, req)
func (s *VoiceService) List(ctx context.Context, req *voice.VoiceListRequest) (*voice.VoiceListResponse, error) {
	// Build query parameters
	query := make(map[string]string)
	if req.VoiceType != "" {
		query["voiceType"] = req.VoiceType
	}
	if req.VoiceName != "" {
		query["voiceName"] = req.VoiceName
	}
	if req.RequestID != "" {
		query["request_id"] = req.RequestID
	}

	// Make the API request
	apiResp, err := s.client.Get(ctx, "/voice/list", query)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp voice.VoiceListResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
