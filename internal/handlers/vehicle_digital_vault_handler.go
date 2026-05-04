package handlers

import (
	"context"
	"net/http"
	"sync"

	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"
	"the-unified-document-viewer/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type VehicleDigitalVaultHandler struct {
	Repo     *repository.PostgresRepository
	JobQueue chan worker.Job
	Adapter interface {
		FetchSalesByVIN(ctx context.Context, vin string) (models.RawSalesData, error)
		FetchServiceByVIN(ctx context.Context, vin string) (models.RawServiceData, error)
	} // Add adapter for parallel API calls
}

func (h *VehicleDigitalVaultHandler) GetVehicleHistory(c *gin.Context) {
    vin := c.Param("vin")
    if vin == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "VIN is required"})
        return
    }

    // Check if VIN exists in database
    exists, count, err := h.Repo.CheckIfVINExists(vin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // If no records exist, return 404
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{
            "vin": vin,
            "exists": false,
            "message": "No documents found for this VIN",
        })
        return
    }

    // Get all documents for this VIN
    docs, err := h.Repo.GetVehicleDigitalVaultByVIN(vin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Return success with documents
    c.JSON(http.StatusOK, gin.H{
        "vin": vin,
        "exists": true,
        "total_documents": count,
        "documents": docs,
    })
}

// SearchAndSyncByVIN handles POST /vault/search - parallel fetch from sales/service APIs, upsert records, return all by VIN
func (h *VehicleDigitalVaultHandler) SearchAndSyncByVIN(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := telemetry.GetTracer().Start(ctx, "vin.search-sync",
		trace.WithAttributes(
			attribute.String("operation", "parallel-api-sync"),
		),
	)
	defer span.End()

	type VinRequest struct {
		VIN string `json:"vin" binding:"required"`
	}

	var req VinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: VIN required"})
		return
	}

	vin := req.VIN
	if vin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VIN is required"})
		return
	}

	span.SetAttributes(attribute.String("vin", vin))

	// Parallel fetch from Sales and Service APIs using goroutines
	var (
		salesData     models.RawSalesData
		serviceData   models.RawServiceData
		salesErr      error
		serviceErr    error
		wg            sync.WaitGroup
		salesCtx      = ctx
		serviceCtx    = ctx
	)

	wg.Add(2)

	// Goroutine 1: Sales API
	go func() {
		defer wg.Done()
		salesData, salesErr = h.Adapter.FetchSalesByVIN(salesCtx, vin)
	}()

	// Goroutine 2: Service API  
	go func() {
		defer wg.Done()
		serviceData, serviceErr = h.Adapter.FetchServiceByVIN(serviceCtx, vin)
	}()

	wg.Wait()

	// Process sales data
	if salesErr == nil {
		vaultSales := worker.MapSalesToVault(salesData)
		worker.EnrichVaultData(&vaultSales)
		if err := h.Repo.UpsertVehicleDigitalVaultRecord(vaultSales); err != nil {
			span.SetAttributes(attribute.String("upsert.sales.error", err.Error()))
		} else {
			span.SetAttributes(attribute.Bool("upsert.sales.success", true))
		}
	} else {
		span.SetAttributes(attribute.String("fetch.sales.error", salesErr.Error()))
	}

	// Process service data
	if serviceErr == nil {
		vaultService := worker.MapServiceToVault(serviceData)
		worker.EnrichVaultData(&vaultService)
		if err := h.Repo.UpsertVehicleDigitalVaultRecord(vaultService); err != nil {
			span.SetAttributes(attribute.String("upsert.service.error", err.Error()))
		} else {
			span.SetAttributes(attribute.Bool("upsert.service.success", true))
		}
	} else {
		span.SetAttributes(attribute.String("fetch.service.error", serviceErr.Error()))
	}

	// Get all records for this VIN (now including newly created/updated)
	exists, count, err := h.Repo.CheckIfVINExists(vin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query vault records"})
		return
	}

	docs, err := h.Repo.GetVehicleDigitalVaultByVIN(vin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query vault records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vin":              vin,
		"exists":           exists,
		"total_documents":  count,
		"documents":        docs,
		"sales_fetched":    salesErr == nil,
		"service_fetched":  serviceErr == nil,
	})
}
