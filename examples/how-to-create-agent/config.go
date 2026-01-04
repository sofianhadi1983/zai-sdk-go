package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey           string `mapstructure:"api_key"`
	BaseURL          string `mapstructure:"base_url"`
	MaxMessages      int    `mapstructure:"max_messages"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	LogLevel         string `mapstructure:"log_level"`
	SaveConversation bool   `mapstructure:"save_conversation"`
	EnableStreaming  bool   `mapstructure:"enable_streaming"`
	EnableColors     bool   `mapstructure:"enable_colors"`
	ShowProgress     bool   `mapstructure:"show_progress"`
	TypingSpeedMs    int    `mapstructure:"typing_speed_ms"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.BindEnv("api_key", "ZAI_API_KEY")
	viper.BindEnv("base_url", "ZAI_BASE_URL")
	viper.BindEnv("max_messages", "ZAI_MAX_MESSAGES")
	viper.BindEnv("enable_logging", "ZAI_ENABLE_LOGGING")
	viper.BindEnv("log_level", "ZAI_LOG_LEVEL")
	viper.BindEnv("save_conversation", "ZAI_SAVE_CONVERSATION")
	viper.BindEnv("enable_streaming", "ZAI_ENABLE_STREAMING")
	viper.BindEnv("enable_colors", "ZAI_ENABLE_COLORS")
	viper.BindEnv("show_progress", "ZAI_SHOW_PROGRESS")
	viper.BindEnv("typing_speed_ms", "ZAI_TYPING_SPEED_MS")

	viper.SetDefault("base_url", "https://api.z.ai/api/coding/paas/v4")
	viper.SetDefault("max_messages", 50)
	viper.SetDefault("enable_logging", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("save_conversation", false)
	viper.SetDefault("enable_streaming", true)
	viper.SetDefault("enable_colors", true)
	viper.SetDefault("show_progress", true)
	viper.SetDefault("typing_speed_ms", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required (set in config.yaml or ZAI_API_KEY env var)")
	}

	return &config, nil
}
