package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/webreader"
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

	// Example 1: Basic web page reading
	fmt.Println("=== Basic Web Reader Example ===")
	basicExample(ctx, client)

	fmt.Println()

	// Example 2: Web reading with markdown format
	fmt.Println("=== Markdown Format Example ===")
	markdownExample(ctx, client)

	fmt.Println()

	// Example 3: Web reading with images and links summary
	fmt.Println("=== Images and Links Summary Example ===")
	summaryExample(ctx, client)

	fmt.Println()

	// Example 4: Web reading with cache disabled
	fmt.Println("=== No Cache Example ===")
	noCacheExample(ctx, client)
}

func basicExample(ctx context.Context, client *zai.Client) {
	// Note: Replace with an actual URL
	url := os.Getenv("WEB_READER_URL")
	if url == "" {
		url = "https://example.com"
		fmt.Printf("Using default URL: %s\n", url)
		fmt.Println("(Set WEB_READER_URL environment variable to use a different URL)")
	}

	req := webreader.NewRequest(url)

	resp, err := client.WebReader.Read(ctx, req)
	if err != nil {
		log.Fatalf("Failed to read web page: %v", err)
	}

	if resp.HasResult() {
		result := resp.GetResult()
		fmt.Printf("Title: %s\n", result.GetTitle())
		fmt.Printf("URL: %s\n", result.URL)
		fmt.Printf("Description: %s\n", result.GetDescription())

		if result.HasContent() {
			content := result.GetContent()
			// Print first 500 characters
			if len(content) > 500 {
				fmt.Printf("\nContent (first 500 chars):\n%s...\n", content[:500])
				fmt.Printf("\nTotal content length: %d characters\n", len(content))
			} else {
				fmt.Printf("\nContent:\n%s\n", content)
			}
		}
	} else {
		fmt.Println("No result returned from web reader")
	}
}

func markdownExample(ctx context.Context, client *zai.Client) {
	url := os.Getenv("WEB_READER_MARKDOWN_URL")
	if url == "" {
		fmt.Println("Skipping example: WEB_READER_MARKDOWN_URL not set")
		fmt.Println("Usage: export WEB_READER_MARKDOWN_URL=https://blog.example.com/article")
		return
	}

	req := webreader.NewRequest(url).
		SetReturnFormat("markdown").
		SetRetainImages(true).
		SetRequestID("markdown_req_123")

	resp, err := client.WebReader.Read(ctx, req)
	if err != nil {
		log.Fatalf("Failed to read web page: %v", err)
	}

	if resp.HasResult() {
		result := resp.GetResult()
		fmt.Printf("Title: %s\n", result.GetTitle())

		if result.PublishedTime != "" {
			fmt.Printf("Published: %s\n", result.PublishedTime)
		}

		if result.HasContent() {
			fmt.Println("\nMarkdown Content:")
			content := result.GetContent()
			// Print first 800 characters of markdown
			if len(content) > 800 {
				fmt.Printf("%s...\n", content[:800])
				fmt.Printf("\n(Content truncated, total length: %d characters)\n", len(content))
			} else {
				fmt.Println(content)
			}
		}
	}
}

func summaryExample(ctx context.Context, client *zai.Client) {
	url := os.Getenv("WEB_READER_SUMMARY_URL")
	if url == "" {
		fmt.Println("Skipping example: WEB_READER_SUMMARY_URL not set")
		fmt.Println("Usage: export WEB_READER_SUMMARY_URL=https://news.example.com")
		return
	}

	req := webreader.NewRequest(url).
		SetReturnFormat("text").
		SetWithImagesSummary(true).
		SetWithLinksSummary(true).
		SetUserID("user_789")

	resp, err := client.WebReader.Read(ctx, req)
	if err != nil {
		log.Fatalf("Failed to read web page: %v", err)
	}

	if resp.HasResult() {
		result := resp.GetResult()
		fmt.Printf("Title: %s\n", result.GetTitle())
		fmt.Printf("Description: %s\n", result.GetDescription())

		// Display images summary
		images := result.GetImages()
		if len(images) > 0 {
			fmt.Printf("\nImages found: %d\n", len(images))
			count := 0
			for key, url := range images {
				if count < 5 { // Show first 5 images
					fmt.Printf("  - %s: %s\n", key, url)
					count++
				}
			}
			if len(images) > 5 {
				fmt.Printf("  ... and %d more images\n", len(images)-5)
			}
		}

		// Display links summary
		links := result.GetLinks()
		if len(links) > 0 {
			fmt.Printf("\nLinks found: %d\n", len(links))
			count := 0
			for key, url := range links {
				if count < 5 { // Show first 5 links
					fmt.Printf("  - %s: %s\n", key, url)
					count++
				}
			}
			if len(links) > 5 {
				fmt.Printf("  ... and %d more links\n", len(links)-5)
			}
		}

		// Display metadata if available
		if result.Metadata != nil && len(result.Metadata) > 0 {
			fmt.Println("\nMetadata:")
			for key, value := range result.Metadata {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
}

func noCacheExample(ctx context.Context, client *zai.Client) {
	url := os.Getenv("WEB_READER_NOCACHE_URL")
	if url == "" {
		fmt.Println("Skipping example: WEB_READER_NOCACHE_URL not set")
		fmt.Println("Usage: export WEB_READER_NOCACHE_URL=https://dynamic.example.com")
		return
	}

	req := webreader.NewRequest(url).
		SetNoCache(true).
		SetReturnFormat("markdown").
		SetTimeout("30").
		SetRequestID("nocache_req_456")

	fmt.Println("Reading web page with cache disabled...")
	resp, err := client.WebReader.Read(ctx, req)
	if err != nil {
		log.Fatalf("Failed to read web page: %v", err)
	}

	if resp.HasResult() {
		result := resp.GetResult()
		fmt.Printf("Title: %s\n", result.GetTitle())
		fmt.Printf("Content length: %d characters\n", len(result.GetContent()))
		fmt.Println("\nNote: Content was fetched fresh without using cache")
	}
}
