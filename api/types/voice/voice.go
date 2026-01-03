// Package voice provides types for the Voice API.
package voice

// VoiceCloneRequest represents a request to clone a voice.
type VoiceCloneRequest struct {
	// VoiceName is the name for the cloned voice.
	VoiceName string `json:"voice_name"`

	// Text is the text content corresponding to the sample audio.
	Text string `json:"text"`

	// Input is the target text for preview audio.
	Input string `json:"input"`

	// FileID is the file ID of the uploaded audio file.
	FileID string `json:"file_id"`

	// Model is the model to use for voice cloning.
	Model string `json:"model"`

	// RequestID is an optional request ID for tracking.
	RequestID string `json:"request_id,omitempty"`
}

// NewVoiceCloneRequest creates a new voice clone request.
func NewVoiceCloneRequest(voiceName, text, input, fileID, model string) *VoiceCloneRequest {
	return &VoiceCloneRequest{
		VoiceName: voiceName,
		Text:      text,
		Input:     input,
		FileID:    fileID,
		Model:     model,
	}
}

// SetRequestID sets the request ID.
func (r *VoiceCloneRequest) SetRequestID(requestID string) *VoiceCloneRequest {
	r.RequestID = requestID
	return r
}

// VoiceCloneResponse represents the response from a voice clone operation.
type VoiceCloneResponse struct {
	// Voice is the voice identifier.
	Voice string `json:"voice"`

	// FileID is the audio preview file ID.
	FileID string `json:"file_id"`

	// FilePurpose is the file purpose.
	FilePurpose string `json:"file_purpose"`
}

// VoiceDeleteRequest represents a request to delete a voice.
type VoiceDeleteRequest struct {
	// Voice is the voice to delete.
	Voice string `json:"voice"`

	// RequestID is an optional request ID for tracking.
	RequestID string `json:"request_id,omitempty"`
}

// NewVoiceDeleteRequest creates a new voice delete request.
func NewVoiceDeleteRequest(voice string) *VoiceDeleteRequest {
	return &VoiceDeleteRequest{
		Voice: voice,
	}
}

// SetRequestID sets the request ID.
func (r *VoiceDeleteRequest) SetRequestID(requestID string) *VoiceDeleteRequest {
	r.RequestID = requestID
	return r
}

// VoiceDeleteResponse represents the response from a voice delete operation.
type VoiceDeleteResponse struct {
	// Voice is the voice identifier.
	Voice string `json:"voice"`

	// UpdateTime is the delete time (format: yyyy-MM-dd HH:mm:ss).
	UpdateTime string `json:"update_time"`
}

// VoiceData represents voice data information.
type VoiceData struct {
	// Voice is the voice identifier.
	Voice string `json:"voice"`

	// VoiceName is the voice name.
	VoiceName string `json:"voice_name"`

	// VoiceType is the voice type.
	VoiceType string `json:"voice_type"`

	// DownloadURL is the download URL for the voice.
	DownloadURL string `json:"download_url"`

	// CreateTime is the creation time (format: yyyy-MM-dd HH:mm:ss).
	CreateTime string `json:"create_time"`
}

// VoiceListRequest represents a request to list voices.
type VoiceListRequest struct {
	// VoiceType is an optional filter by voice type.
	VoiceType string `json:"voice_type,omitempty"`

	// VoiceName is an optional filter by voice name.
	VoiceName string `json:"voice_name,omitempty"`

	// RequestID is an optional request ID for tracking.
	RequestID string `json:"request_id,omitempty"`
}

// NewVoiceListRequest creates a new voice list request.
func NewVoiceListRequest() *VoiceListRequest {
	return &VoiceListRequest{}
}

// SetVoiceType sets the voice type filter.
func (r *VoiceListRequest) SetVoiceType(voiceType string) *VoiceListRequest {
	r.VoiceType = voiceType
	return r
}

// SetVoiceName sets the voice name filter.
func (r *VoiceListRequest) SetVoiceName(voiceName string) *VoiceListRequest {
	r.VoiceName = voiceName
	return r
}

// SetRequestID sets the request ID.
func (r *VoiceListRequest) SetRequestID(requestID string) *VoiceListRequest {
	r.RequestID = requestID
	return r
}

// VoiceListResponse represents the response from a voice list operation.
type VoiceListResponse struct {
	// VoiceList contains the list of voices.
	VoiceList []VoiceData `json:"voice_list"`
}

// GetVoices returns the list of voices.
func (r *VoiceListResponse) GetVoices() []VoiceData {
	if r.VoiceList == nil {
		return []VoiceData{}
	}
	return r.VoiceList
}
