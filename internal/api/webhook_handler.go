// internal/api/webhook_handler.go
package api

import (
	"fmt"
	"net/http"
	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	JobQueue chan worker.Job
}

// POST /webhooks/sales
func (h *WebhookHandler) HandleSalesWebhook(c *gin.Context) {
	fmt.Printf("check")
	var payload models.RawSalesData
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sales Payload"})
		return
	}

	// Async Dispatch: Push job to queue and return 202 immediately
	h.JobQueue <- worker.Job{Type: worker.SourceSales, Payload: payload}
	c.JSON(http.StatusAccepted, gin.H{"status": "Sales data queued for processing"})
}

// POST /webhooks/service
func (h *WebhookHandler) HandleServiceWebhook(c *gin.Context) {
	var payload models.RawServiceData
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Service Payload"})
		return
	}

	// Async Dispatch
	h.JobQueue <- worker.Job{Type: worker.SourceService, Payload: payload}
	c.JSON(http.StatusAccepted, gin.H{"status": "Service data queued for processing"})
}
