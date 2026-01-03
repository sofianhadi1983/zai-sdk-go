package zai

import (
	"context"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/webreader"
	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
)

// WebReaderService provides access to the Web Reader API.
type WebReaderService struct {
	client *client.BaseClient
}

// newWebReaderService creates a new web reader service.
func newWebReaderService(baseClient *client.BaseClient) *WebReaderService {
	return &WebReaderService{
		client: baseClient,
	}
}

// Read reads and extracts content from a web page.
//
// Example:
//
//	req := webreader.NewRequest("https://example.com").
//	    SetReturnFormat("markdown").
//	    SetRetainImages(true).
//	    SetWithLinksSummary(true)
//
//	resp, err := client.WebReader.Read(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	if resp.HasResult() {
//	    result := resp.GetResult()
//	    fmt.Printf("Title: %s\n", result.GetTitle())
//	    fmt.Printf("Content: %s\n", result.GetContent())
//	    fmt.Printf("Images: %d\n", len(result.GetImages()))
//	    fmt.Printf("Links: %d\n", len(result.GetLinks()))
//	}
//
// Example with request tracking:
//
//	req := webreader.NewRequest("https://blog.example.com/article").
//	    SetRequestID("req_123").
//	    SetUserID("user_456").
//	    SetReturnFormat("text").
//	    SetTimeout("30")
//
//	resp, err := client.WebReader.Read(ctx, req)
//
// Example with caching disabled:
//
//	req := webreader.NewRequest("https://news.example.com").
//	    SetNoCache(true).
//	    SetReturnFormat("markdown").
//	    SetWithImagesSummary(true)
//
//	resp, err := client.WebReader.Read(ctx, req)
func (s *WebReaderService) Read(ctx context.Context, req *webreader.Request) (*webreader.Response, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/reader", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp webreader.Response
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
