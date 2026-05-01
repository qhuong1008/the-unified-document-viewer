package worker

import (
	"fmt"
	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"
	"time"
)

// StartWorkerPool initializes N workers for parallel processing
func StartWorkerPool(jobQueue chan Job, workerCount int, repository *repository.PostgresRepository) {
	for i := 1; i <= workerCount; i++ {
		go func(workerID int) {
			for job := range jobQueue {
				ExecuteJob(workerID, job, repository)
			}
		}(i)
	}
}

// ExecuteJob orchestrates the Transformation, Enrichment, and Persistence process
func ExecuteJob(workerID int, job Job, repository *repository.PostgresRepository) {
	fmt.Printf("[Worker %d] [START] Processing job of type %s\n", workerID, job.Type)
	
	// Bước 1: Transformation
	vaultData, ok := transformJobPayload(job)
	if !ok {
		fmt.Printf("[Worker %d] [ERROR] Failed to transform job of type %s\n", workerID, job.Type)
		return
	}

	// Add 5-second timeout to test parallel execution
	fmt.Printf("[Worker %d] [TIMEOUT] Starting 5-second processing delay...\n", workerID)
	time.Sleep(5 * time.Second)
	fmt.Printf("[Worker %d] [TIMEOUT] Finished 5-second delay\n", workerID)

	// Bước 2: Enrichment
	EnrichVaultData(&vaultData)

	// Bước 3: Persistence (Lưu trữ bền vững)
	if err := repository.UpsertVehicleDigitalVaultRecord(vaultData); err != nil {
		fmt.Printf("[Worker %d] [DATABASE ERROR] Failed to sync VIN %s: %v\n", workerID, vaultData.VIN, err)
		return
	}

	// Log phục vụ chiến lược Observability theo yêu cầu của Keyloop
	fmt.Printf("[Worker %d] [SUCCESS] Synced %s record for VIN %s. Category: %s (ID: %s)\n",
		workerID, vaultData.SourceSystem, vaultData.VIN, vaultData.DocCategory, vaultData.ExternalID)
}

func transformJobPayload(job Job) (models.VehicleDigitalVault, bool) {
	switch job.Type {
	case SourceSales:
		if raw, ok := job.Payload.(models.RawSalesData); ok {
			return MapSalesToVault(raw), true
		}
	case SourceService:
		if raw, ok := job.Payload.(models.RawServiceData); ok {
			return MapServiceToVault(raw), true
		}
	}
	return models.VehicleDigitalVault{}, false
}