// Package logger provides structured logging for the Z.ai SDK.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// contextKey is the type for context keys used in logging.
type contextKey string

const (
	// TraceIDKey is the context key for trace IDs.
	TraceIDKey contextKey = "trace_id"

	// RequestIDKey is the context key for request IDs.
	RequestIDKey contextKey = "request_id"

	// UserIDKey is the context key for user IDs.
	UserIDKey contextKey = "user_id"
)

// Level represents the severity level of a log message.
type Level = slog.Level

const (
	// LevelDebug represents debug-level messages.
	LevelDebug = slog.LevelDebug

	// LevelInfo represents info-level messages.
	LevelInfo = slog.LevelInfo

	// LevelWarn represents warning-level messages.
	LevelWarn = slog.LevelWarn

	// LevelError represents error-level messages.
	LevelError = slog.LevelError
)

// Config holds logger configuration.
type Config struct {
	// Level is the minimum log level to output.
	Level Level

	// Format determines the output format ("json" or "text").
	Format string

	// Output is the writer to send logs to (default: os.Stdout).
	Output io.Writer

	// AddSource adds source file and line number to logs.
	AddSource bool
}

// DefaultConfig returns the default logger configuration.
func DefaultConfig() *Config {
	return &Config{
		Level:     LevelInfo,
		Format:    "text",
		Output:    os.Stdout,
		AddSource: false,
	}
}

// Logger wraps slog.Logger with additional functionality.
type Logger struct {
	*slog.Logger
}

// New creates a new logger with the given configuration.
func New(cfg *Config) *Logger {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Default returns the default logger.
var defaultLogger = New(DefaultConfig())

// Default returns the default logger instance.
func Default() *Logger {
	return defaultLogger
}

// SetDefault sets the default logger.
func SetDefault(logger *Logger) {
	defaultLogger = logger
	slog.SetDefault(logger.Logger)
}

// Context-aware logging methods

// WithContext returns a logger with context attributes added.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := extractContextAttrs(ctx)
	if len(attrs) == 0 {
		return l
	}
	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

// DebugContext logs a debug message with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logContext(ctx, LevelDebug, msg, args...)
}

// InfoContext logs an info message with context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logContext(ctx, LevelInfo, msg, args...)
}

// WarnContext logs a warning message with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logContext(ctx, LevelWarn, msg, args...)
}

// ErrorContext logs an error message with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logContext(ctx, LevelError, msg, args...)
}

// logContext logs a message with context attributes.
func (l *Logger) logContext(ctx context.Context, level Level, msg string, args ...any) {
	attrs := extractContextAttrs(ctx)
	args = append(attrs, args...)
	l.Logger.Log(ctx, level, msg, args...)
}

// extractContextAttrs extracts logging attributes from context.
func extractContextAttrs(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}

	var attrs []any

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok && id != "" {
			attrs = append(attrs, slog.String("trace_id", id))
		}
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok && id != "" {
			attrs = append(attrs, slog.String("request_id", id))
		}
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok && id != "" {
			attrs = append(attrs, slog.String("user_id", id))
		}
	}

	return attrs
}

// Helper functions for adding context values

// WithTraceID adds a trace ID to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetTraceID retrieves the trace ID from the context.
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID retrieves the user ID from the context.
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// Package-level convenience functions that use the default logger

// Debug logs a debug message using the default logger.
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger.
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger.
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// DebugContext logs a debug message with context using the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context using the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context using the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context using the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}

// With creates a new logger with the given attributes.
func With(args ...any) *Logger {
	return &Logger{
		Logger: defaultLogger.Logger.With(args...),
	}
}
