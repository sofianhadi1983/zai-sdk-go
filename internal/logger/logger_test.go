package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *Config
		want   Level
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
			want:   LevelInfo,
		},
		{
			name: "custom config",
			config: &Config{
				Level:  LevelDebug,
				Format: "json",
			},
			want: LevelDebug,
		},
		{
			name: "text format",
			config: &Config{
				Level:  LevelWarn,
				Format: "text",
			},
			want: LevelWarn,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := New(tt.config)
			if logger == nil {
				t.Fatal("New() returned nil logger")
			}
			if logger.Logger == nil {
				t.Fatal("Logger.Logger is nil")
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.Level != LevelInfo {
		t.Errorf("DefaultConfig().Level = %v, want %v", cfg.Level, LevelInfo)
	}

	if cfg.Format != "text" {
		t.Errorf("DefaultConfig().Format = %q, want %q", cfg.Format, "text")
	}

	if cfg.AddSource != false {
		t.Errorf("DefaultConfig().AddSource = %v, want false", cfg.AddSource)
	}
}

func TestLogger_JSONOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	logger.Info("test message", slog.String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("JSON output missing message: %s", output)
	}

	// Verify it's valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}

	// Verify fields
	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
		t.Errorf("JSON output incorrect msg field: %v", logEntry["msg"])
	}

	if key, ok := logEntry["key"].(string); !ok || key != "value" {
		t.Errorf("JSON output incorrect key field: %v", logEntry["key"])
	}
}

func TestLogger_TextOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "text",
		Output: &buf,
	})

	logger.Info("test message", slog.String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Text output missing message: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Text output missing attribute: %s", output)
	}
}

func TestLogger_LogLevels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		logLevel  Level
		logFunc   func(*Logger, string)
		shouldLog bool
	}{
		{
			name:      "debug when level is debug",
			logLevel:  LevelDebug,
			logFunc:   func(l *Logger, msg string) { l.Debug(msg) },
			shouldLog: true,
		},
		{
			name:      "debug when level is info",
			logLevel:  LevelInfo,
			logFunc:   func(l *Logger, msg string) { l.Debug(msg) },
			shouldLog: false,
		},
		{
			name:      "info when level is info",
			logLevel:  LevelInfo,
			logFunc:   func(l *Logger, msg string) { l.Info(msg) },
			shouldLog: true,
		},
		{
			name:      "warn when level is warn",
			logLevel:  LevelWarn,
			logFunc:   func(l *Logger, msg string) { l.Warn(msg) },
			shouldLog: true,
		},
		{
			name:      "warn when level is error",
			logLevel:  LevelError,
			logFunc:   func(l *Logger, msg string) { l.Warn(msg) },
			shouldLog: false,
		},
		{
			name:      "error when level is error",
			logLevel:  LevelError,
			logFunc:   func(l *Logger, msg string) { l.Error(msg) },
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := New(&Config{
				Level:  tt.logLevel,
				Format: "text",
				Output: &buf,
			})

			tt.logFunc(logger, "test message")

			hasOutput := buf.Len() > 0
			if hasOutput != tt.shouldLog {
				t.Errorf("shouldLog = %v, but hasOutput = %v, output: %q", tt.shouldLog, hasOutput, buf.String())
			}
		})
	}
}

func TestLogger_WithContext(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "json",
		Output: &buf,
	})

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")
	ctx = WithUserID(ctx, "user-789")

	logger.InfoContext(ctx, "test message")

	output := buf.String()
	if !strings.Contains(output, "trace-123") {
		t.Errorf("Output missing trace ID: %s", output)
	}
	if !strings.Contains(output, "req-456") {
		t.Errorf("Output missing request ID: %s", output)
	}
	if !strings.Contains(output, "user-789") {
		t.Errorf("Output missing user ID: %s", output)
	}
}

func TestContextHelpers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Test WithTraceID and GetTraceID
	ctx = WithTraceID(ctx, "trace-123")
	if got := GetTraceID(ctx); got != "trace-123" {
		t.Errorf("GetTraceID() = %q, want %q", got, "trace-123")
	}

	// Test WithRequestID and GetRequestID
	ctx = WithRequestID(ctx, "req-456")
	if got := GetRequestID(ctx); got != "req-456" {
		t.Errorf("GetRequestID() = %q, want %q", got, "req-456")
	}

	// Test WithUserID and GetUserID
	ctx = WithUserID(ctx, "user-789")
	if got := GetUserID(ctx); got != "user-789" {
		t.Errorf("GetUserID() = %q, want %q", got, "user-789")
	}

	// Test empty context
	emptyCtx := context.Background()
	if got := GetTraceID(emptyCtx); got != "" {
		t.Errorf("GetTraceID(emptyCtx) = %q, want empty string", got)
	}
	if got := GetRequestID(emptyCtx); got != "" {
		t.Errorf("GetRequestID(emptyCtx) = %q, want empty string", got)
	}
	if got := GetUserID(emptyCtx); got != "" {
		t.Errorf("GetUserID(emptyCtx) = %q, want empty string", got)
	}
}

func TestLogger_ContextMethods(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:  LevelDebug,
		Format: "text",
		Output: &buf,
	})

	ctx := WithTraceID(context.Background(), "trace-123")

	tests := []struct {
		name    string
		logFunc func()
		want    string
	}{
		{
			name:    "DebugContext",
			logFunc: func() { logger.DebugContext(ctx, "debug message") },
			want:    "debug message",
		},
		{
			name:    "InfoContext",
			logFunc: func() { logger.InfoContext(ctx, "info message") },
			want:    "info message",
		},
		{
			name:    "WarnContext",
			logFunc: func() { logger.WarnContext(ctx, "warn message") },
			want:    "warn message",
		},
		{
			name:    "ErrorContext",
			logFunc: func() { logger.ErrorContext(ctx, "error message") },
			want:    "error message",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("Output missing message %q: %s", tt.want, output)
			}
			if !strings.Contains(output, "trace-123") {
				t.Errorf("Output missing trace ID: %s", output)
			}
		})
	}
}

func TestPackageLevelFunctions(t *testing.T) {
	// Note: Not parallel because we're using the default logger

	var buf bytes.Buffer
	SetDefault(New(&Config{
		Level:  LevelDebug,
		Format: "text",
		Output: &buf,
	}))

	tests := []struct {
		name    string
		logFunc func()
		want    string
	}{
		{
			name:    "Debug",
			logFunc: func() { Debug("debug message") },
			want:    "debug message",
		},
		{
			name:    "Info",
			logFunc: func() { Info("info message") },
			want:    "info message",
		},
		{
			name:    "Warn",
			logFunc: func() { Warn("warn message") },
			want:    "warn message",
		},
		{
			name:    "Error",
			logFunc: func() { Error("error message") },
			want:    "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("Output missing message %q: %s", tt.want, output)
			}
		})
	}
}

func TestPackageLevelContextFunctions(t *testing.T) {
	// Note: Not parallel because we're using the default logger

	var buf bytes.Buffer
	SetDefault(New(&Config{
		Level:  LevelDebug,
		Format: "text",
		Output: &buf,
	}))

	ctx := WithTraceID(context.Background(), "trace-123")

	tests := []struct {
		name    string
		logFunc func()
		want    string
	}{
		{
			name:    "DebugContext",
			logFunc: func() { DebugContext(ctx, "debug message") },
			want:    "debug message",
		},
		{
			name:    "InfoContext",
			logFunc: func() { InfoContext(ctx, "info message") },
			want:    "info message",
		},
		{
			name:    "WarnContext",
			logFunc: func() { WarnContext(ctx, "warn message") },
			want:    "warn message",
		},
		{
			name:    "ErrorContext",
			logFunc: func() { ErrorContext(ctx, "error message") },
			want:    "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("Output missing message %q: %s", tt.want, output)
			}
			if !strings.Contains(output, "trace-123") {
				t.Errorf("Output missing trace ID: %s", output)
			}
		})
	}
}

func TestLogger_With(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "text",
		Output: &buf,
	})

	childLogger := &Logger{
		Logger: logger.Logger.With(slog.String("component", "test")),
	}

	childLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "component=test") {
		t.Errorf("Output missing component attribute: %s", output)
	}
}

func TestWith(t *testing.T) {
	// Note: Not parallel because we're using the default logger

	var buf bytes.Buffer
	SetDefault(New(&Config{
		Level:  LevelInfo,
		Format: "text",
		Output: &buf,
	}))

	logger := With(slog.String("service", "zai-sdk"))
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "service=zai-sdk") {
		t.Errorf("Output missing service attribute: %s", output)
	}
}

func TestLogger_AddSource(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := New(&Config{
		Level:     LevelInfo,
		Format:    "text",
		Output:    &buf,
		AddSource: true,
	})

	logger.Info("test message")

	output := buf.String()
	// Should contain source file information
	if !strings.Contains(output, "logger_test.go") {
		t.Errorf("Output missing source file: %s", output)
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()

	logger := Default()
	if logger == nil {
		t.Fatal("Default() returned nil")
	}
	if logger.Logger == nil {
		t.Fatal("Default().Logger is nil")
	}
}

func TestExtractContextAttrs_NilContext(t *testing.T) {
	t.Parallel()

	attrs := extractContextAttrs(nil)
	if len(attrs) != 0 {
		t.Errorf("extractContextAttrs(nil) returned %d attrs, want 0", len(attrs))
	}
}

func TestExtractContextAttrs_EmptyContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attrs := extractContextAttrs(ctx)
	if len(attrs) != 0 {
		t.Errorf("extractContextAttrs(empty) returned %d attrs, want 0", len(attrs))
	}
}

func TestExtractContextAttrs_PartialContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	// Only trace ID, no request ID or user ID

	attrs := extractContextAttrs(ctx)
	if len(attrs) != 1 { // one slog.Attr
		t.Errorf("extractContextAttrs() returned %d attrs, want 1", len(attrs))
	}
}

// Benchmark tests
func BenchmarkLogger_Info(b *testing.B) {
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "text",
		Output: io.Discard,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", slog.Int("iteration", i))
	}
}

func BenchmarkLogger_InfoContext(b *testing.B) {
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "text",
		Output: io.Discard,
	})

	ctx := WithTraceID(context.Background(), "trace-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "benchmark message", slog.Int("iteration", i))
	}
}

func BenchmarkLogger_JSONOutput(b *testing.B) {
	logger := New(&Config{
		Level:  LevelInfo,
		Format: "json",
		Output: io.Discard,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", slog.Int("iteration", i))
	}
}
