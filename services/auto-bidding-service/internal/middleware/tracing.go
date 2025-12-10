package middleware

import (
	"log/slog"
	"time"

	"category_service/internal/metrics"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("fiber-http")

// TracingMiddleware thêm OpenTelemetry tracing vào Fiber requests
func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract trace context từ headers
		ctx := otel.GetTextMapPropagator().Extract(
			c.Context(),
			propagation.HeaderCarrier(c.GetReqHeaders()),
		)

		// Start span
		ctx, span := tracer.Start(
			ctx,
			c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.route", c.Route().Path),
				attribute.String("http.url", c.OriginalURL()),
				attribute.String("http.user_agent", c.Get("User-Agent")),
				attribute.String("http.client_ip", c.IP()),
			),
		)
		defer span.End()

		// Store context in fiber context
		c.SetUserContext(ctx)

		// Track active requests
		if metrics.HTTPActiveRequests != nil {
			metrics.HTTPActiveRequests.Add(ctx, 1)
			defer metrics.HTTPActiveRequests.Add(ctx, -1)
		}

		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Record duration
		duration := time.Since(start).Milliseconds()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
			attribute.Int64("http.response_size", int64(len(c.Response().Body()))),
			attribute.Int64("http.request_size", int64(len(c.Request().Body()))),
		)

		// Record metrics
		metricAttrs := []attribute.KeyValue{
			attribute.String("method", c.Method()),
			attribute.String("route", c.Route().Path),
			attribute.Int("status", c.Response().StatusCode()),
		}

		if metrics.HTTPRequestCounter != nil {
			metrics.HTTPRequestCounter.Add(ctx, 1,
				metric.WithAttributes(metricAttrs...),
			)
		}

		if metrics.HTTPRequestDuration != nil {
			metrics.HTTPRequestDuration.Record(ctx, float64(duration),
				metric.WithAttributes(metricAttrs...),
			)
		}

		// Log HTTP request với trace context
		statusCode := c.Response().StatusCode()
		logAttrs := []any{
			"method", c.Method(),
			"route", c.Route().Path,
			"status", statusCode,
			"duration_ms", duration,
			"client_ip", c.IP(),
			"trace_id", span.SpanContext().TraceID().String(),
			"span_id", span.SpanContext().SpanID().String(),
		}

		if statusCode >= 500 {
			slog.ErrorContext(ctx, "HTTP request", logAttrs...)
		} else if statusCode >= 400 {
			slog.WarnContext(ctx, "HTTP request", logAttrs...)
		} else {
			slog.InfoContext(ctx, "HTTP request", logAttrs...)
		}

		// Handle errors
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetStatus(codes.Error, fiber.ErrBadRequest.Message)
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return nil
	}
}
