package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerInstance trace.Tracer
)

// InitTracer initializes the OpenTelemetry tracer with OTLP HTTP exporter
// that sends traces to OTel Collector at localhost:4318
func InitTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	// Create OTLP HTTP exporter pointing to OTel Collector at localhost:4318
	// Using otlptracehttp which is the modern replacement for otlpexporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service name attribute
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create TracerProvider with the exporter and resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)),
	)

	// Configure global TracerProvider
	otel.SetTracerProvider(tp)

	// Configure TextMapPropagator to extract TraceContext from HTTP headers
	// This allows Frontend to pass Trace ID via HTTP headers (W3C Trace Context)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Create and store the tracer instance
	tracerInstance = otel.Tracer(serviceName)

	return tp, nil
}

// GetTracer returns the configured tracer instance
func GetTracer() trace.Tracer {
	if tracerInstance == nil {
		// Return a no-op tracer if InitTracer hasn't been called
		return otel.Tracer("noop")
	}
	return tracerInstance
}

// GetTracerWithName returns a tracer with a specific name for creating child spans
func GetTracerWithName(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, name, opts...)
}

// Shutdown gracefully shuts down the TracerProvider
func Shutdown(ctx context.Context) error {
	if tracerInstance != nil {
		tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
		if ok {
			return tp.Shutdown(ctx)
		}
	}
	return nil
}

// Helper functions for common tracing operations

// ExtractTraceContext extracts trace context from incoming HTTP headers
// Use this in middleware to get the context from request headers
func ExtractTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	return ctx
}

// InjectTraceContext injects trace context into HTTP headers
// Use this when making outgoing HTTP calls to other services
func InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}
