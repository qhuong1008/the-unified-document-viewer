package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryMiddleware creates a Gin middleware that automatically starts a span
// for each incoming HTTP request. It extracts trace context from incoming headers
// using the configured TextMapPropagator (TraceContext).
func TelemetryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract trace context from incoming HTTP headers
		// This allows Frontend to pass Trace ID via W3C Trace Context headers
		ctx := ExtractTraceContextFromGin(c)

		// Get the tracer
		tracer := GetTracer()

		// Create span name from HTTP method and path
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())

		// Start a new span with the extracted context
		// If there's no trace context, this creates a new trace
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.path", c.Request.URL.Path),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.client_ip", c.ClientIP()),
				attribute.String("user_agent", c.Request.Header.Get("User-Agent")),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)

		// Defer span ending - record request duration correctly
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)

			// Set response status code attribute
			if c.Writer.Status() >= 400 {
				span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP error %d", c.Writer.Status()))
			} else {
				span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
				span.SetStatus(codes.Ok, "")
			}

			// Record actual request duration in milliseconds
			span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))

			span.End()
		}()

		// Store span in Gin context for handlers to access
		c.Set("span", span)
		c.Set("trace_context", ctx)

		// Continue with the request
		c.Next()
	}
}

// ExtractTraceContextFromGin extracts trace context from Gin request headers
// Returns context.Context which can be used directly with tracer.Start
func ExtractTraceContextFromGin(c *gin.Context) context.Context {
	// Create a TextMapCarrier from Gin headers
	carrier := ginTextMapCarrier{headers: c.Request.Header}

	// Use the global TextMapPropagator to extract context
	// This uses TraceContext propagation by default
	ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), carrier)

	return ctx
}

// ginTextMapCarrier implements TextMapCarrier for Gin request headers
type ginTextMapCarrier struct {
	headers map[string][]string
}

func (g ginTextMapCarrier) Get(key string) string {
	if vals, ok := g.headers[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (g ginTextMapCarrier) Set(key string, value string) {
	// This is a read-only carrier for extraction
	// No-op for extraction
}

func (g ginTextMapCarrier) Keys() []string {
	keys := make([]string, 0, len(g.headers))
	for k := range g.headers {
		keys = append(keys, k)
	}
	return keys
}

// GetTraceContext retrieves the trace context from Gin context
// Use this in your handlers to get the context for creating child spans
func GetTraceContext(c *gin.Context) (context.Context, bool) {
	ctxVal, exists := c.Get("trace_context")
	if !exists {
		return nil, false
	}
	return ctxVal.(context.Context), true
}

// GetSpan retrieves the current span from Gin context
func GetSpan(c *gin.Context) (trace.Span, bool) {
	spanVal, exists := c.Get("span")
	if !exists {
		return nil, false
	}
	return spanVal.(trace.Span), true
}
