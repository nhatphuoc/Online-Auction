package metrics

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	// HTTP Metrics
	HTTPRequestCounter  metric.Int64Counter
	HTTPRequestDuration metric.Float64Histogram
	HTTPActiveRequests  metric.Int64UpDownCounter

	// gRPC Metrics
	GRPCRequestCounter  metric.Int64Counter
	GRPCRequestDuration metric.Float64Histogram
	GRPCActiveRequests  metric.Int64UpDownCounter

	// Database Metrics
	DBQueryDuration     metric.Float64Histogram
	DBConnectionsActive metric.Int64UpDownCounter
)

// InitMetrics khởi tạo tất cả metrics
func InitMetrics(ctx context.Context) error {
	meter := otel.Meter("final4-api")

	var err error

	// HTTP Metrics
	HTTPRequestCounter, err = meter.Int64Counter(
		"http.server.requests",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	HTTPRequestDuration, err = meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	HTTPActiveRequests, err = meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	// gRPC Metrics
	GRPCRequestCounter, err = meter.Int64Counter(
		"grpc.server.requests",
		metric.WithDescription("Total number of gRPC requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	GRPCRequestDuration, err = meter.Float64Histogram(
		"grpc.server.duration",
		metric.WithDescription("gRPC request duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	GRPCActiveRequests, err = meter.Int64UpDownCounter(
		"grpc.server.active_requests",
		metric.WithDescription("Number of active gRPC requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	// Database Metrics
	DBQueryDuration, err = meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}

	DBConnectionsActive, err = meter.Int64UpDownCounter(
		"db.connections.active",
		metric.WithDescription("Number of active database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return err
	}

	slog.Info("Metrics initialized successfully")
	return nil
}
