// internal/worker/pool.go
package worker

import (
	"fmt"
	"the-unified-document-viewer/internal/models"
)

// Dispatcher starts N workers to process jobs in parallel
func StartWorkerPool(jobQueue chan Job, workerCount int) {
	// This loop creates 'workerCount' number of goroutines (parallel threads)
	for i := 1; i <= workerCount; i++ {
		go func(workerID int) {
			// This range loop keeps the worker alive, waiting for data in the channel
			for job := range jobQueue {
				// This is the "Trigger" point
				processJob(workerID, job)
			}
		}(i)
	}
}

func processJob(workerID int, job Job) {
	var vin string
	switch p := job.Payload.(type) {
	case models.RawSalesData:
		vin = p.VehicleVIN
	case models.RawServiceData:
		vin = p.VehicleVIN
	}
	fmt.Printf("[Worker %d] Transforming & Upserting data for VIN: %s\n", workerID, vin)
}
