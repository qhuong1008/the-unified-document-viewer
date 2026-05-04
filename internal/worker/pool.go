package worker

import (
	"context"
	"fmt"
	"time"

	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	// Create a background context for the job
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start parent span for the entire job processing
	startTime := time.Now()
	ctx, parentSpan := telemetry.GetTracer().Start(ctx, "job.execute",
		trace.WithAttributes(
			attribute.String("job.type", string(job.Type)),
			attribute.Int("worker.id", workerID),
			attribute.String("operation.type", "job-processing"),
		),
	)
	defer parentSpan.End()

	fmt.Printf("[Worker %d] [START] Processing job of type %s\n", workerID, job.Type)

// === Bước 1: Transformation ===
	_, transformSpan := telemetry.GetTracer().Start(ctx, "job.transform",
		trace.WithAttributes(
			attribute.String("operation", "transform"),
			attribute.String("job.type", string(job.Type)),
		),
	)
	defer transformSpan.End()

	vaultData, ok := transformJobPayload(job)
	if !ok {
		transformSpan.SetStatus(codes.Error, "Failed to transform job payload")
		transformSpan.SetAttributes(attribute.Bool("transform.success", false))
		fmt.Printf("[Worker %d] [ERROR] Failed to transform job of type %s\n", workerID, job.Type)
		parentSpan.SetStatus(codes.Error, "Transform failed")
		parentSpan.SetAttributes(attribute.Bool("job.success", false))
		return
	}

	transformSpan.SetAttributes(attribute.Bool("transform.success", true))

	// Add VIN attribute for filtering across all spans
	if vaultData.VIN != "" {
		parentSpan.SetAttributes(attribute.String("vehicle.vin", vaultData.VIN))
		transformSpan.SetAttributes(attribute.String("vehicle.vin", vaultData.VIN))
	}

// Add source attribute
	parentSpan.SetAttributes(attribute.String("api.source", string(job.Type)))

	// Add 5-second timeout to test parallel execution
	fmt.Printf("[Worker %d] [TIMEOUT] Starting 5-second processing delay...\n", workerID)
	// time.Sleep(5 * time.Second)
	fmt.Printf("[Worker %d] [TIMEOUT] Finished 5-second delay\n", workerID)
	fmt.Printf("vaultData: %+v ", vaultData)

// === Bước 2: Enrichment ===
	_, enrichSpan := telemetry.GetTracer().Start(ctx, "job.enrich",
		trace.WithAttributes(
			attribute.String("operation", "enrich"),
			attribute.String("vehicle.vin", vaultData.VIN),
		),
	)
	defer enrichSpan.End()

	EnrichVaultData(&vaultData)

	// Add enriched source attribute
	enrichSpan.SetAttributes(attribute.String("api.source", "enriched"))
	enrichSpan.SetAttributes(attribute.Bool("enrich.success", true))

// === Bước 3: Persistence (Lưu trữ bền vững) ===
	_, persistSpan := telemetry.GetTracer().Start(ctx, "database.persist",
		trace.WithAttributes(
			attribute.String("database.operation", "upsert"),
			attribute.String("database.table", "vehicle_digital_vault"),
			attribute.String("vehicle.vin", vaultData.VIN),
		),
	)
	defer persistSpan.End()

	if err := repository.UpsertVehicleDigitalVaultRecord(vaultData); err != nil {
		persistSpan.SetStatus(codes.Error, err.Error())
		persistSpan.SetAttributes(attribute.String("error.message", err.Error()))
		persistSpan.SetAttributes(attribute.Bool("persist.success", false))
		fmt.Printf("[Worker %d] [DATABASE ERROR] Failed to sync VIN %s: %v\n", workerID, vaultData.VIN, err)
		parentSpan.SetStatus(codes.Error, "Persistence failed")
		parentSpan.SetAttributes(attribute.Bool("job.success", false))
		return
	}

	// Record success
	persistSpan.SetAttributes(attribute.Bool("persist.success", true))

	// Record total duration
	duration := time.Since(startTime)
	parentSpan.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))
	parentSpan.SetAttributes(attribute.Bool("job.success", true))

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
