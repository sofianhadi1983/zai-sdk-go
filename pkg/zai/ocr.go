package zai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/ocr"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
)

// OCRService provides access to the OCR API.
type OCRService struct {
	client *client.BaseClient
}

// newOCRService creates a new OCR service.
func newOCRService(baseClient *client.BaseClient) *OCRService {
	return &OCRService{
		client: baseClient,
	}
}

// HandwritingOCR performs handwriting OCR on an image.
//
// Example:
//
//	file, err := os.Open("handwriting.jpg")
//	if err != nil {
//	    // Handle error
//	}
//	defer file.Close()
//
//	req := ocr.NewOCRRequest(file, "handwriting.jpg", ocr.ToolTypeHandWrite).
//	    SetProbability(true).
//	    SetLanguageType("zh-CN")
//
//	resp, err := client.OCR.HandwritingOCR(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.HasResults() {
//	    fmt.Printf("Recognized text: %s\n", resp.GetText())
//	    for _, result := range resp.GetResults() {
//	        fmt.Printf("  - %s (location: %+v)\n", result.Words, result.Location)
//	    }
//	}
func (s *OCRService) HandwritingOCR(ctx context.Context, req *ocr.OCRRequest) (*ocr.OCRResponse, error) {
	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the tool_type field (required)
	if err := writer.WriteField("tool_type", string(req.ToolType)); err != nil {
		return nil, fmt.Errorf("failed to write tool_type field: %w", err)
	}

	// Add optional fields
	if req.LanguageType != "" {
		if err := writer.WriteField("language_type", req.LanguageType); err != nil {
			return nil, fmt.Errorf("failed to write language_type field: %w", err)
		}
	}

	if req.Probability {
		if err := writer.WriteField("probability", strconv.FormatBool(req.Probability)); err != nil {
			return nil, fmt.Errorf("failed to write probability field: %w", err)
		}
	}

	// Add the image file
	part, err := writer.CreateFormFile("file", req.FileName)
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
	apiResp, err := s.client.PostMultipart(ctx, "/files/ocr", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp ocr.OCRResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
