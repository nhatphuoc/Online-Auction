package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type OTelConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTelEndpoint   string // OTel Collector endpoint
}

type OTelShutdown func(context.Context) error

// Global logger provider for slog bridge
var globalLoggerProvider *sdklog.LoggerProvider

// InitOTel khởi tạo OpenTelemetry với tracing, metrics, và logs
func InitOTel(ctx context.Context, cfg OTelConfig) (OTelShutdown, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup Trace Provider
	tracerProvider, err := setupTracerProvider(ctx, res, cfg.OTelEndpoint)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tracerProvider)

	// Setup Meter Provider
	meterProvider, err := setupMeterProvider(ctx, res, cfg.OTelEndpoint)
	if err != nil {
		return nil, err
	}
	otel.SetMeterProvider(meterProvider)

	// Setup Log Provider
	loggerProvider, err := setupLoggerProvider(ctx, res, cfg.OTelEndpoint)
	if err != nil {
		return nil, err
	}
	globalLoggerProvider = loggerProvider

	// Setup Propagator
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	slog.Info("OpenTelemetry initialized",
		"service", cfg.ServiceName,
		"endpoint", cfg.OTelEndpoint,
	)

	// Return shutdown function
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var shutdownErr error
		if err := tracerProvider.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("tracer provider shutdown: %w", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			if shutdownErr != nil {
				shutdownErr = fmt.Errorf("%v; meter provider shutdown: %w", shutdownErr, err)
			} else {
				shutdownErr = fmt.Errorf("meter provider shutdown: %w", err)
			}
		}
		if err := loggerProvider.Shutdown(ctx); err != nil {
			if shutdownErr != nil {
				shutdownErr = fmt.Errorf("%v; logger provider shutdown: %w", shutdownErr, err)
			} else {
				shutdownErr = fmt.Errorf("logger provider shutdown: %w", err)
			}
		}
		return shutdownErr
	}, nil
}

// setupTracerProvider tạo TracerProvider với OTLP exporter
func setupTracerProvider(ctx context.Context, res *resource.Resource, endpoint string) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(), // Sử dụng insecure cho local development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()), // Sample tất cả traces
	)

	return tp, nil
}

// setupMeterProvider tạo MeterProvider với OTLP exporter
func setupMeterProvider(ctx context.Context, res *resource.Resource, endpoint string) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(), // Sử dụng insecure cho local development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)

	return mp, nil
}

// setupLoggerProvider tạo LoggerProvider với OTLP exporter
func setupLoggerProvider(ctx context.Context, res *resource.Resource, endpoint string) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // Sử dụng insecure cho local development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	return lp, nil
}

// GetLoggerProvider returns the global logger provider
func GetLoggerProvider() *sdklog.LoggerProvider {
	return globalLoggerProvider
}
