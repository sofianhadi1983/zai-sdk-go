package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/moderation"
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

	// Example 1: Basic text moderation
	fmt.Println("=== Example 1: Basic Text Moderation ===")
	basicModerationExample(ctx, client)

	// Example 2: Batch moderation
	fmt.Println("\n=== Example 2: Batch Text Moderation ===")
	batchModerationExample(ctx, client)

	// Example 3: Using convenience methods
	fmt.Println("\n=== Example 3: Convenience Methods ===")
	convenienceModerationExample(ctx, client)

	// Example 4: Checking specific categories
	fmt.Println("\n=== Example 4: Checking Specific Categories ===")
	categoryCheckExample(ctx, client)

	// Example 5: Analyzing category scores
	fmt.Println("\n=== Example 5: Analyzing Category Scores ===")
	scoreAnalysisExample(ctx, client)
}

func basicModerationExample(ctx context.Context, client *zai.Client) {
	// Create a moderation request for a single text
	req := moderation.NewTextModerationRequest(
		"moderation",
		"This is a sample text to check for content safety.",
	)

	resp, err := client.Moderations.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Request ID: %s\n", resp.ID)
	fmt.Printf("Model: %s\n", resp.Model)

	// Check overall flagged status
	if resp.IsFlagged() {
		fmt.Println("⚠️  Content flagged as potentially harmful")
	} else {
		fmt.Println("✓ Content appears safe")
	}

	// Display results for each input
	for i, result := range resp.GetResults() {
		fmt.Printf("\nResult %d:\n", i+1)
		fmt.Printf("  Flagged: %v\n", result.Flagged)
		fmt.Printf("  Safe: %v\n", result.IsSafe())

		if result.Flagged {
			displayFlaggedCategories(result.Categories)
		}
	}
}

func batchModerationExample(ctx context.Context, client *zai.Client) {
	// Moderate multiple texts at once
	texts := []string{
		"Hello, how are you today?",
		"This is a professional business email.",
		"Let's discuss the project requirements.",
	}

	req := moderation.NewBatchTextModerationRequest("moderation", texts)

	resp, err := client.Moderations.Create(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Moderated %d texts\n", len(resp.GetResults()))
	fmt.Printf("Overall flagged: %v\n\n", resp.IsFlagged())

	// Display results for each text
	for i, result := range resp.GetResults() {
		fmt.Printf("Text %d: \"%s\"\n", i+1, texts[i])
		if result.IsSafe() {
			fmt.Println("  Status: ✓ Safe")
		} else {
			fmt.Println("  Status: ⚠️  Flagged")
			displayFlaggedCategories(result.Categories)
		}
		fmt.Println()
	}
}

func convenienceModerationExample(ctx context.Context, client *zai.Client) {
	// Using CheckText convenience method
	fmt.Println("Checking single text:")
	resp, err := client.Moderations.CheckText(
		ctx,
		"moderation",
		"Sample text for quick moderation check",
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.IsFlagged() {
		fmt.Println("  ⚠️  Content flagged")
	} else {
		fmt.Println("  ✓ Content safe")
	}

	// Using CheckBatch convenience method
	fmt.Println("\nChecking batch of texts:")
	texts := []string{
		"First text to moderate",
		"Second text to moderate",
		"Third text to moderate",
	}

	respBatch, err := client.Moderations.CheckBatch(ctx, "moderation", texts)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	safeCount := 0
	flaggedCount := 0
	for _, result := range respBatch.GetResults() {
		if result.IsSafe() {
			safeCount++
		} else {
			flaggedCount++
		}
	}

	fmt.Printf("  Results: %d safe, %d flagged\n", safeCount, flaggedCount)
}

func categoryCheckExample(ctx context.Context, client *zai.Client) {
	// Check for specific content categories
	testCases := []struct {
		description string
		text        string
	}{
		{
			description: "Professional content",
			text:        "Let's schedule a meeting to discuss the quarterly report.",
		},
		{
			description: "Potentially harmful content",
			text:        "Example of content that might be flagged by moderation.",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\nChecking: %s\n", tc.description)

		resp, err := client.Moderations.CheckText(ctx, "moderation", tc.text)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		if len(resp.GetResults()) == 0 {
			continue
		}

		result := resp.GetResults()[0]
		categories := result.Categories

		// Check specific categories
		categoriesChecked := []struct {
			name    string
			flagged bool
		}{
			{"Harassment", categories.Harassment},
			{"Hate", categories.Hate},
			{"Self-Harm", categories.SelfHarm},
			{"Sexual", categories.Sexual},
			{"Violence", categories.Violence},
		}

		flaggedCategories := []string{}
		for _, cat := range categoriesChecked {
			if cat.flagged {
				flaggedCategories = append(flaggedCategories, cat.name)
			}
		}

		if len(flaggedCategories) > 0 {
			fmt.Printf("  Flagged categories: %v\n", flaggedCategories)
		} else {
			fmt.Println("  No categories flagged ✓")
		}
	}
}

func scoreAnalysisExample(ctx context.Context, client *zai.Client) {
	// Analyze category scores to understand confidence levels
	text := "Sample text for detailed score analysis"

	resp, err := client.Moderations.CheckText(ctx, "moderation", text)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if len(resp.GetResults()) == 0 {
		return
	}

	result := resp.GetResults()[0]
	scores := result.CategoryScores

	fmt.Println("Category Scores (0.0 = safe, 1.0 = harmful):")

	scoreList := []struct {
		name  string
		score float64
	}{
		{"Harassment", scores.Harassment},
		{"Harassment/Threatening", scores.HarassmentThreatening},
		{"Hate", scores.Hate},
		{"Hate/Threatening", scores.HateThreatening},
		{"Self-Harm", scores.SelfHarm},
		{"Self-Harm/Instructions", scores.SelfHarmInstructions},
		{"Self-Harm/Intent", scores.SelfHarmIntent},
		{"Sexual", scores.Sexual},
		{"Sexual/Minors", scores.SexualMinors},
		{"Violence", scores.Violence},
		{"Violence/Graphic", scores.ViolenceGraphic},
	}

	// Display scores and highlight high-risk categories
	for _, s := range scoreList {
		indicator := "✓"
		if s.score > 0.5 {
			indicator = "⚠️"
		} else if s.score > 0.8 {
			indicator = "🚫"
		}

		fmt.Printf("  %s %-25s: %.4f\n", indicator, s.name, s.score)
	}

	// Summary
	fmt.Println("\nSummary:")
	if result.Flagged {
		fmt.Println("  Overall: ⚠️  Content flagged for review")
	} else {
		fmt.Println("  Overall: ✓ Content appears safe")
	}
}

// Helper function to display flagged categories
func displayFlaggedCategories(categories moderation.ModerationCategories) {
	flagged := []string{}

	if categories.Harassment {
		flagged = append(flagged, "Harassment")
	}
	if categories.HarassmentThreatening {
		flagged = append(flagged, "Harassment/Threatening")
	}
	if categories.Hate {
		flagged = append(flagged, "Hate")
	}
	if categories.HateThreatening {
		flagged = append(flagged, "Hate/Threatening")
	}
	if categories.SelfHarm {
		flagged = append(flagged, "Self-Harm")
	}
	if categories.SelfHarmInstructions {
		flagged = append(flagged, "Self-Harm/Instructions")
	}
	if categories.SelfHarmIntent {
		flagged = append(flagged, "Self-Harm/Intent")
	}
	if categories.Sexual {
		flagged = append(flagged, "Sexual")
	}
	if categories.SexualMinors {
		flagged = append(flagged, "Sexual/Minors")
	}
	if categories.Violence {
		flagged = append(flagged, "Violence")
	}
	if categories.ViolenceGraphic {
		flagged = append(flagged, "Violence/Graphic")
	}

	if len(flagged) > 0 {
		fmt.Printf("  Flagged Categories: %v\n", flagged)
	}
}
