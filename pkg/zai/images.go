package zai

import (
	"context"

	"github.com/z-ai/zai-sdk-go/api/types/images"
	"github.com/z-ai/zai-sdk-go/internal/client"
)

// ImagesService provides access to the Images API.
type ImagesService struct {
	client *client.BaseClient
}

// newImagesService creates a new images service.
func newImagesService(baseClient *client.BaseClient) *ImagesService {
	return &ImagesService{
		client: baseClient,
	}
}

// Create generates images based on the provided prompt.
//
// Example:
//
//	req := images.NewImageGenerationRequest("cogview-3", "A beautiful sunset over mountains")
//	req.SetSize(images.Size1024x1792).SetQuality(images.QualityHD)
//
//	resp, err := client.Images.Create(ctx, req)
//	if err != nil {
//	    // Handle error
//	}
//
//	firstImage := resp.GetFirstImage()
//	if firstImage != nil {
//	    fmt.Printf("Image URL: %s\n", firstImage.GetImageURL())
//	}
func (s *ImagesService) Create(ctx context.Context, req *images.ImageGenerationRequest) (*images.ImageGenerationResponse, error) {
	// Make the API request
	apiResp, err := s.client.Post(ctx, "/images/generations", req)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var resp images.ImageGenerationResponse
	if err := s.client.ParseJSON(apiResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Generate is a convenience method for generating a single image from a text prompt.
// Returns the URL or base64 data of the first generated image.
//
// Example:
//
//	imageURL, err := client.Images.Generate(ctx, "cogview-3", "A cat playing piano")
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Generated image: %s\n", imageURL)
func (s *ImagesService) Generate(ctx context.Context, model, prompt string) (string, error) {
	req := images.NewImageGenerationRequest(model, prompt)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return "", err
	}

	firstImage := resp.GetFirstImage()
	if firstImage == nil {
		return "", nil
	}

	// Return URL if available, otherwise return base64 data
	if firstImage.URL != "" {
		return firstImage.URL, nil
	}

	return firstImage.B64JSON, nil
}

// GenerateMultiple is a convenience method for generating multiple images from a text prompt.
// Returns URLs or base64 data of all generated images.
//
// Example:
//
//	imageURLs, err := client.Images.GenerateMultiple(ctx, "cogview-3", "A dog in a park", 3)
//	if err != nil {
//	    // Handle error
//	}
//
//	for i, url := range imageURLs {
//	    fmt.Printf("Image %d: %s\n", i+1, url)
//	}
func (s *ImagesService) GenerateMultiple(ctx context.Context, model, prompt string, count int) ([]string, error) {
	req := images.NewImageGenerationRequest(model, prompt)
	req.SetN(count)

	resp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Return URLs if available, otherwise return base64 data
	urls := resp.GetImageURLs()
	if len(urls) > 0 {
		return urls, nil
	}

	return resp.GetBase64Images(), nil
}
