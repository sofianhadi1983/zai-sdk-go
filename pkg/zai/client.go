// Package zai provides the main SDK client for Z.ai API.
package zai

import (
	"os"
	"time"

	"github.com/sofianhadi1983/zai-sdk-go/internal/client"
	"github.com/sofianhadi1983/zai-sdk-go/internal/constants"
	"github.com/sofianhadi1983/zai-sdk-go/internal/logger"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai/errors"
)

// Client is the main SDK client for Z.ai API.
type Client struct {
	baseClient *client.BaseClient
	config     *ClientConfig

	// Chat provides access to the Chat Completions API.
	Chat *ChatService

	// Embeddings provides access to the Embeddings API.
	Embeddings *EmbeddingsService

	// Images provides access to the Images API.
	Images *ImagesService

	// Files provides access to the Files API.
	Files *FilesService

	// Videos provides access to the Videos API.
	Videos *VideosService

	// Audio provides access to the Audio API.
	Audio *AudioService

	// Assistant provides access to the Assistant API.
	Assistant *AssistantService

	// Batch provides access to the Batch API.
	Batch *BatchService

	// WebSearch provides access to the Web Search API.
	WebSearch *WebSearchService

	// Moderations provides access to the Moderations API.
	Moderations *ModerationsService

	// Tools provides access to the Tools API.
	Tools *ToolsService

	// Agents provides access to the Agents API.
	Agents *AgentsService

	// Voice provides access to the Voice API.
	Voice *VoiceService

	// OCR provides access to the OCR API.
	OCR *OCRService

	// FileParser provides access to the File Parser API.
	FileParser *FileParserService

	// WebReader provides access to the Web Reader API.
	WebReader *WebReaderService
}

// ClientConfig holds configuration for the SDK client.
type ClientConfig struct {
	// APIKey is the API key for authentication (format: "key.secret").
	APIKey string

	// BaseURL is the base URL for API requests.
	// If empty, uses the default Z.ai API URL.
	BaseURL string

	// Timeout is the request timeout.
	// If zero, uses the default timeout (120 seconds).
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	// If zero, uses the default (3 retries).
	MaxRetries int

	// DisableTokenCache disables JWT token caching.
	// When true, uses raw API key for authentication.
	DisableTokenCache bool

	// Logger is a custom logger.
	// If nil, uses the default logger.
	Logger *logger.Logger
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*ClientConfig)

// WithAPIKey sets the API key for the client.
//
// The API key should be in the format "key.secret" as provided
// by the Z.ai platform.
//
// Example:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("abc123.xyz789"),
//	)
func WithAPIKey(apiKey string) ClientOption {
	return func(c *ClientConfig) {
		c.APIKey = apiKey
	}
}

// WithBaseURL sets the base URL for API requests.
//
// Use this option when you need to use a custom API endpoint,
// such as a proxy or a different regional endpoint.
//
// Example:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-key"),
//	    zai.WithBaseURL("https://api.custom.com/v1"),
//	)
func WithBaseURL(baseURL string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the request timeout.
//
// This controls how long the client will wait for a response
// before timing out. Default is 120 seconds.
//
// Example:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-key"),
//	    zai.WithTimeout(60 * time.Second),
//	)
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
//
// The client will automatically retry failed requests up to
// this number of times with exponential backoff. Default is 3.
//
// Example:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-key"),
//	    zai.WithMaxRetries(5),
//	)
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *ClientConfig) {
		c.MaxRetries = maxRetries
	}
}

// WithDisableTokenCache disables JWT token caching.
//
// By default, the client caches JWT tokens to reduce overhead.
// Use this option if you want to disable caching and generate
// a new token for each request.
//
// Example:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-key"),
//	    zai.WithDisableTokenCache(),
//	)
func WithDisableTokenCache() ClientOption {
	return func(c *ClientConfig) {
		c.DisableTokenCache = true
	}
}

// WithLogger sets a custom logger.
//
// Use this option to provide your own logger for debugging
// and monitoring API requests.
//
// Example:
//
//	customLogger := logger.New(logger.LevelDebug)
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-key"),
//	    zai.WithLogger(customLogger),
//	)
func WithLogger(logger *logger.Logger) ClientOption {
	return func(c *ClientConfig) {
		c.Logger = logger
	}
}

// NewClient creates a new Z.ai SDK client for overseas users.
// The default base URL is https://open.bigmodel.cn/api/paas/v4/
//
// Basic usage:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-api-key.your-secret"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// With additional options:
//
//	client, err := zai.NewClient(
//	    zai.WithAPIKey("your-api-key.your-secret"),
//	    zai.WithBaseURL("https://custom.api.url"),
//	    zai.WithTimeout(120 * time.Second),
//	    zai.WithMaxRetries(5),
//	    zai.WithLogger(customLogger),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
func NewClient(opts ...ClientOption) (*Client, error) {
	config := &ClientConfig{
		BaseURL: constants.ZaiBaseURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return newClient(config)
}

// NewZhipuClient creates a new Z.ai SDK client for Chinese users.
// The default base URL is https://open.bigmodel.cn/api/paas/v4/
//
// This client is optimized for users in mainland China and uses
// the domestic API endpoint.
//
// Example:
//
//	client, err := zai.NewZhipuClient(
//	    zai.WithAPIKey("your-api-key.your-secret"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Use the client for chat completions
//	messages := []chat.Message{
//	    {Role: chat.RoleUser, Content: "你好！"},
//	}
//	req := chat.NewChatCompletionRequest(chat.ModelGLM4Plus, messages)
//	resp, err := client.Chat.Create(context.Background(), req)
func NewZhipuClient(opts ...ClientOption) (*Client, error) {
	config := &ClientConfig{
		BaseURL: constants.ZhipuBaseURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return newClient(config)
}

// NewClientFromEnv creates a new client from environment variables.
// Reads ZAI_API_KEY and optionally ZAI_BASE_URL from environment.
//
// Environment variables:
//   - ZAI_API_KEY: Required. Your API key in format "key.secret"
//   - ZAI_BASE_URL: Optional. Custom base URL for API requests
//
// Example:
//
//	// Set environment variables first:
//	// export ZAI_API_KEY="your-api-key.your-secret"
//	// export ZAI_BASE_URL="https://custom.api.url"  # optional
//
//	client, err := zai.NewClientFromEnv()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// You can also override settings with options:
//
//	client, err := zai.NewClientFromEnv(
//	    zai.WithTimeout(60 * time.Second),
//	    zai.WithMaxRetries(3),
//	)
func NewClientFromEnv(opts ...ClientOption) (*Client, error) {
	apiKey := os.Getenv("ZAI_API_KEY")
	baseURL := os.Getenv("ZAI_BASE_URL")

	config := &ClientConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	// Apply additional options
	for _, opt := range opts {
		opt(config)
	}

	// Set default base URL if not provided
	if config.BaseURL == "" {
		config.BaseURL = constants.ZaiBaseURL
	}

	return newClient(config)
}

// NewZhipuClientFromEnv creates a new Chinese client from environment variables.
// Reads ZAI_API_KEY from environment and uses Zhipu base URL.
//
// This is a convenience function for Chinese users that automatically
// uses the domestic API endpoint (https://open.bigmodel.cn/api/paas/v4/).
//
// Environment variables:
//   - ZAI_API_KEY: Required. Your API key in format "key.secret"
//
// Example:
//
//	// Set environment variable first:
//	// export ZAI_API_KEY="your-api-key.your-secret"
//
//	client, err := zai.NewZhipuClientFromEnv()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
func NewZhipuClientFromEnv(opts ...ClientOption) (*Client, error) {
	apiKey := os.Getenv("ZAI_API_KEY")

	config := &ClientConfig{
		APIKey:  apiKey,
		BaseURL: constants.ZhipuBaseURL,
	}

	// Apply additional options
	for _, opt := range opts {
		opt(config)
	}

	return newClient(config)
}

// newClient creates a new client from the given configuration.
func newClient(config *ClientConfig) (*Client, error) {
	// Validate configuration
	if config.APIKey == "" {
		return nil, errors.NewConfigError("APIKey", "API key is required")
	}

	// Create internal base client config
	baseConfig := &client.Config{
		APIKey:            config.APIKey,
		BaseURL:           config.BaseURL,
		Timeout:           config.Timeout,
		MaxRetries:        config.MaxRetries,
		DisableTokenCache: config.DisableTokenCache,
		Logger:            config.Logger,
	}

	// Create base client
	baseClient, err := client.NewBaseClient(baseConfig)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseClient: baseClient,
		config:     config,
	}

	// Initialize services
	c.Chat = newChatService(baseClient)
	c.Embeddings = newEmbeddingsService(baseClient)
	c.Images = newImagesService(baseClient)
	c.Files = newFilesService(baseClient)
	c.Videos = newVideosService(baseClient)
	c.Audio = newAudioService(baseClient)
	c.Assistant = newAssistantService(baseClient)
	c.Batch = newBatchService(baseClient)
	c.WebSearch = newWebSearchService(baseClient)
	c.Moderations = newModerationsService(baseClient)
	c.Tools = newToolsService(baseClient)
	c.Agents = newAgentsService(baseClient)
	c.Voice = newVoiceService(baseClient)
	c.OCR = newOCRService(baseClient)
	c.FileParser = newFileParserService(baseClient)
	c.WebReader = newWebReaderService(baseClient)

	return c, nil
}

// GetConfig returns the client configuration.
//
// This method allows you to inspect the current client configuration
// including API key, base URL, timeout, and other settings.
//
// Example:
//
//	config := client.GetConfig()
//	fmt.Printf("Base URL: %s\n", config.BaseURL)
//	fmt.Printf("Timeout: %v\n", config.Timeout)
//	fmt.Printf("Max Retries: %d\n", config.MaxRetries)
func (c *Client) GetConfig() *ClientConfig {
	return c.config
}

// GetLogger returns the client logger.
//
// Use this method to access the logger for custom logging or debugging.
//
// Example:
//
//	logger := client.GetLogger()
//	logger.Debug("Custom debug message")
//	logger.Info("Custom info message")
func (c *Client) GetLogger() *logger.Logger {
	return c.baseClient.GetLogger()
}

// Close closes the client and releases resources.
//
// This method should be called when you're done using the client
// to ensure proper cleanup of resources, especially HTTP connections.
// It's recommended to use defer to ensure cleanup happens.
//
// Example:
//
//	client, err := zai.NewClient(zai.WithAPIKey("your-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Use client...
func (c *Client) Close() {
	if c.baseClient != nil {
		c.baseClient.Close()
	}
}
