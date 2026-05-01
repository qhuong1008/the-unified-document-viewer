package handlers

import (
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
)

type VehicleDigitalVaultHandler struct {
	repo     *repository.PostgresRepository
	jobQueue chan worker.Job
}

func (h *VehicleDigitalVaultHandler) GetVehicleHistory(c *gin.Context) {
    vin := c.Param("vin")
    if vin == "" {
        c.JSON(400, gin.H{"error": "VIN is required"})
        return
    }

    docs, err := h.repo.GetVehicleDigitalVaultByVIN(vin)
    if err != nil {
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }

		c.JSON(200, gin.H{
        "vin": vin,
        "total_documents": len(docs),
        "documents": docs,
    })
}