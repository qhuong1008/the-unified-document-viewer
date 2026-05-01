package handlers

import (
	"net/http"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
)

type VehicleDigitalVaultHandler struct {
	Repo     *repository.PostgresRepository
	JobQueue chan worker.Job
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
