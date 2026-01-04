package main

import (
	"context"
	"log"

	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := zai.NewClient(
		zai.WithAPIKey(config.APIKey),
		zai.WithBaseURL(config.BaseURL),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	if err := runAgent(ctx, client, config); err != nil {
		log.Fatalf("Agent error: %v", err)
	}
}

func runAgent(ctx context.Context, client *zai.Client, config *Config) error {
	tools := NewToolRegistry()
	agent := NewAgent(client, tools, config)
	return agent.Run(ctx)
}
