package logger

import (
	"context"
	"log/slog"
	"os"

	"api_gateway/internal/telemetry"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var Log *slog.Logger

// InitLogger khởi tạo structured logger với slog
func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		// Production: JSON format
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	} else {
		// Development: Text format với màu sắc
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// InitLoggerWithOTel khởi tạo structured logger với OTel bridge
func InitLoggerWithOTel(env string) {
	var baseHandler slog.Handler

	if env == "production" {
		// Production: JSON format
		baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	} else {
		// Development: Text format
		baseHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}

	// Wrap với OTel bridge để gửi logs qua OTLP
	loggerProvider := telemetry.GetLoggerProvider()
	otelHandler := otelslog.NewHandler("category_service-api", otelslog.WithLoggerProvider(loggerProvider))

	// Combine handlers: console output + OTel export
	handler := combinedHandler{
		console: baseHandler,
		otel:    otelHandler,
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// combinedHandler gửi logs đồng thời đến console và OTel
type combinedHandler struct {
	console slog.Handler
	otel    slog.Handler
}

func (h combinedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.console.Enabled(ctx, level)
}

func (h combinedHandler) Handle(ctx context.Context, r slog.Record) error {
	// Log to console
	if err := h.console.Handle(ctx, r); err != nil {
		return err
	}
	// Log to OTel
	return h.otel.Handle(ctx, r)
}

func (h combinedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return combinedHandler{
		console: h.console.WithAttrs(attrs),
		otel:    h.otel.WithAttrs(attrs),
	}
}

func (h combinedHandler) WithGroup(name string) slog.Handler {
	return combinedHandler{
		console: h.console.WithGroup(name),
		otel:    h.otel.WithGroup(name),
	}
}

// WithContext thêm trace_id và span_id từ context vào log
func WithContext(ctx context.Context) *slog.Logger {
	// Sẽ extract trace context từ OpenTelemetry sau
	return Log
}

// Info logs info message
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Debug logs debug message
func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

// Warn logs warning message
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

// Error logs error message
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

// InfoContext logs info with context
func InfoContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).InfoContext(ctx, msg, args...)
}

// ErrorContext logs error with context
func ErrorContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).ErrorContext(ctx, msg, args...)
}
