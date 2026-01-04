package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/images"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
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

	// Example 1: Simple image generation
	fmt.Println("=== Example 1: Simple Image Generation ===")
	simpleGenerationExample(ctx, client)

	// Example 2: Image generation with convenience method
	fmt.Println("\n=== Example 2: Quick Image Generation ===")
	quickGenerationExample(ctx, client)

	// Example 3: Generate multiple images at once
	fmt.Println("\n=== Example 3: Multiple Image Generation ===")
	multipleImagesExample(ctx, client)

	// Example 4: Custom parameters (size, quality)
	fmt.Println("\n=== Example 4: Custom Parameters ===")
	customParametersExample(ctx, client)

	// Example 5: Base64 response format
	fmt.Println("\n=== Example 5: Base64 Response Format ===")
	base64FormatExample(ctx, client)
}

func simpleGenerationExample(ctx context.Context, client *zai.Client) {
	// Create an image generation request
	req := images.NewImageGenerationRequest(
		"cogview-3",
		"A beautiful sunset over a mountain range with vibrant orange and purple colors",
	)

	// Generate the image
	resp, err := client.Images.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Get the first generated image
	img := resp.GetFirstImage()
	if img != nil {
		fmt.Printf("Generated image URL: %s\n", img.GetImageURL())
		if img.RevisedPrompt != "" {
			fmt.Printf("Revised prompt: %s\n", img.RevisedPrompt)
		}
	}
}

func quickGenerationExample(ctx context.Context, client *zai.Client) {
	// Use the convenience method for quick single image generation
	imageURL, err := client.Images.Generate(
		ctx,
		"cogview-3",
		"A futuristic city at night with neon lights and flying cars",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Generated image URL: %s\n", imageURL)
}

func multipleImagesExample(ctx context.Context, client *zai.Client) {
	// Generate multiple images at once
	imageURLs, err := client.Images.GenerateMultiple(
		ctx,
		"cogview-3",
		"A serene lake surrounded by autumn trees",
		3, // Generate 3 variations
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Generated %d images:\n", len(imageURLs))
	for i, url := range imageURLs {
		fmt.Printf("  Image %d: %s\n", i+1, url)
	}
}

func customParametersExample(ctx context.Context, client *zai.Client) {
	// Create request with custom parameters using fluent API
	req := images.NewImageGenerationRequest(
		"cogview-3",
		"A detailed portrait of a wise old wizard with a long white beard",
	)
	req.SetSize(images.Size1024x1792). // Portrait orientation
						SetQuality(images.QualityHD).        // High quality
						SetN(2).                              // Generate 2 variations
						SetUserID("example-user-123")        // Optional user ID

	resp, err := client.Images.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Generated %d high-quality portrait images\n", len(resp.Data))
	for i, img := range resp.Data {
		fmt.Printf("  Image %d: %s\n", i+1, img.GetImageURL())
	}

	// Show usage information if available
	if resp.Usage != nil {
		fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
	}
}

func base64FormatExample(ctx context.Context, client *zai.Client) {
	// Request images in base64 format instead of URLs
	req := images.NewImageGenerationRequest(
		"cogview-3",
		"A small icon of a coffee cup",
	)
	req.SetSize(images.Size1024x1024).
		SetResponseFormat(images.ResponseFormatB64JSON)

	resp, err := client.Images.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	img := resp.GetFirstImage()
	if img != nil {
		base64Data := img.GetBase64Data()
		if base64Data != "" {
			fmt.Printf("Received base64-encoded image (length: %d bytes)\n", len(base64Data))
			// You can now decode and save the image
			// Example: decode base64 and save to file
			// data, _ := base64.StdEncoding.DecodeString(base64Data)
			// os.WriteFile("image.png", data, 0644)
		}
	}
}
