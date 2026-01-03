// Package ocr provides types for the OCR API.
package ocr

import "io"

// ToolType represents the type of OCR tool to use.
type ToolType string

const (
	// ToolTypeHandWrite is the handwriting recognition tool.
	ToolTypeHandWrite ToolType = "hand_write"
)

// OCRRequest represents a request for OCR processing.
type OCRRequest struct {
	// File is the image file to process (required).
	File io.Reader

	// FileName is the name of the file being uploaded (required).
	FileName string

	// ToolType specifies the OCR tool to use (required).
	ToolType ToolType

	// LanguageType specifies the language for recognition (optional).
	LanguageType string

	// Probability indicates whether to include probability scores (optional).
	Probability bool
}

// NewOCRRequest creates a new OCR request.
func NewOCRRequest(file io.Reader, fileName string, toolType ToolType) *OCRRequest {
	return &OCRRequest{
		File:     file,
		FileName: fileName,
		ToolType: toolType,
	}
}

// SetLanguageType sets the language type.
func (r *OCRRequest) SetLanguageType(languageType string) *OCRRequest {
	r.LanguageType = languageType
	return r
}

// SetProbability sets whether to include probability scores.
func (r *OCRRequest) SetProbability(probability bool) *OCRRequest {
	r.Probability = probability
	return r
}

// Location represents the location of recognized text in the image.
type Location struct {
	// Left is the x-coordinate of the top-left corner.
	Left int `json:"left"`

	// Top is the y-coordinate of the top-left corner.
	Top int `json:"top"`

	// Width is the width of the bounding box.
	Width int `json:"width"`

	// Height is the height of the bounding box.
	Height int `json:"height"`
}

// Probability represents probability scores for recognized text.
type Probability struct {
	// Average is the average confidence score.
	Average float64 `json:"average"`

	// Variance is the variance in confidence scores.
	Variance float64 `json:"variance"`

	// Min is the minimum confidence score.
	Min float64 `json:"min"`
}

// WordsResult represents a single word recognition result.
type WordsResult struct {
	// Location is the bounding box of the recognized text.
	Location Location `json:"location"`

	// Words is the recognized text.
	Words string `json:"words"`

	// Probability contains confidence scores (if requested).
	Probability *Probability `json:"probability,omitempty"`
}

// OCRResponse represents the response from an OCR operation.
type OCRResponse struct {
	// TaskID is the task or result identifier.
	TaskID string `json:"task_id"`

	// Message is the status message.
	Message string `json:"message"`

	// Status is the OCR task status.
	Status string `json:"status"`

	// WordsResultNum is the number of recognition results.
	WordsResultNum int `json:"words_result_num"`

	// WordsResult contains the recognition results.
	WordsResult []WordsResult `json:"words_result,omitempty"`
}

// GetResults returns the recognition results.
// Returns an empty slice if no results are available.
func (r *OCRResponse) GetResults() []WordsResult {
	if r.WordsResult == nil {
		return []WordsResult{}
	}
	return r.WordsResult
}

// HasResults returns true if the response contains recognition results.
func (r *OCRResponse) HasResults() bool {
	return r.WordsResultNum > 0 && len(r.WordsResult) > 0
}

// GetText returns all recognized text concatenated together.
func (r *OCRResponse) GetText() string {
	if !r.HasResults() {
		return ""
	}

	var text string
	for i, result := range r.WordsResult {
		if i > 0 {
			text += " "
		}
		text += result.Words
	}
	return text
}
