package main

import (
	"github.com/gin-gonic/gin"
	"the-unified-document-viewer/internal/api"
	"the-unified-document-viewer/internal/worker"
)

func main() {
	// Create a buffered channel to handle bursts of parallel requests
	jobQueue := make(chan worker.Job, 100)

	// Start 5 parallel workers to process incoming webhooks
	worker.StartWorkerPool(jobQueue, 5)

	r := gin.Default()
	handler := &api.WebhookHandler{JobQueue: jobQueue}

	r.POST("/webhooks/sales", handler.HandleSalesWebhook)
	r.POST("/webhooks/service", handler.HandleServiceWebhook)

	r.Run(":8080")
}
