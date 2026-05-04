// internal/api/webhook_handler.go
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/worker"
	"the-unified-document-viewer/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type WebhookHandler struct {
	JobQueue chan worker.Job
}

// POST /webhooks/sales
func (h *WebhookHandler) HandleSalesWebhook(c *gin.Context) {
	h.handleWebhook(c, "sales", func() (interface{}, error) {
		var payload models.RawSalesData
		if err := c.ShouldBindJSON(&payload); err != nil {
			return nil, err
		}
		return payload, nil
	}, worker.SourceSales)
}

// POST /webhooks/service
func (h *WebhookHandler) HandleServiceWebhook(c *gin.Context) {
	h.handleWebhook(c, "service", func() (interface{}, error) {
		var payload models.RawServiceData
		if err := c.ShouldBindJSON(&payload); err != nil {
			return nil, err
		}
		return payload, nil
	}, worker.SourceService)
}

func (h *WebhookHandler) handleWebhook(
	c *gin.Context,
	webhookType string,
	parsePayload func() (interface{}, error),
	sourceType worker.SourceType,
) {
	startTime := time.Now()

	ctx := getOrCreateContext(c)

	ctx, span := startWebhookSpan(ctx, c, webhookType)
	defer span.End()

	payload, err := parsePayload()
	if err != nil {
		span.SetAttributes(attribute.String("error.message", err.Error()))
		span.SetStatus(codes.Error, "Invalid payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid %s Payload", webhookType)})
		return
	}

	addVINAttribute(span, payload)

	h.JobQueue <- worker.Job{Type: sourceType, Payload: payload}

	span.SetAttributes(attribute.Int64("duration_ms", time.Since(startTime).Milliseconds()))

	c.JSON(http.StatusAccepted, gin.H{
		"status": fmt.Sprintf("%s data queued for processing", webhookType),
	})
}

func getOrCreateContext(c *gin.Context) context.Context {
	ctx, exists := telemetry.GetTraceContext(c)
	if exists {
		return ctx
	}
	return c.Request.Context()
}

func startWebhookSpan(ctx context.Context, c *gin.Context, webhookType string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("webhook.%s.receive", webhookType)
	return telemetry.GetTracer().Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("api.source", "webhook"),
			attribute.String("webhook.type", webhookType),
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.Request.URL.Path),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	)
}

// addVINAttribute adds VIN to span if available in payload
func addVINAttribute(span trace.Span, payload interface{}) {
	var vin string
	switch p := payload.(type) {
	case models.RawSalesData:
		vin = p.VehicleVIN
	case models.RawServiceData:
		vin = p.VehicleVIN
	}
	if vin != "" {
		span.SetAttributes(attribute.String("vehicle.vin", vin))
	}
}
