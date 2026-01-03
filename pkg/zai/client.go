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
func WithAPIKey(apiKey string) ClientOption {
	return func(c *ClientConfig) {
		c.APIKey = apiKey
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *ClientConfig) {
		c.MaxRetries = maxRetries
	}
}

// WithDisableTokenCache disables JWT token caching.
func WithDisableTokenCache() ClientOption {
	return func(c *ClientConfig) {
		c.DisableTokenCache = true
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *logger.Logger) ClientOption {
	return func(c *ClientConfig) {
		c.Logger = logger
	}
}

// NewClient creates a new Z.ai SDK client for overseas users.
// The default base URL is https://open.bigmodel.cn/api/paas/v4/
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

	return c, nil
}

// GetConfig returns the client configuration.
func (c *Client) GetConfig() *ClientConfig {
	return c.config
}

// GetLogger returns the client logger.
func (c *Client) GetLogger() *logger.Logger {
	return c.baseClient.GetLogger()
}

// Close closes the client and releases resources.
func (c *Client) Close() {
	if c.baseClient != nil {
		c.baseClient.Close()
	}
}
