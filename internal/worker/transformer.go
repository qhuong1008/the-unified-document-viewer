// internal/worker/transformer.go
package worker

import (
	"fmt"
	"the-unified-document-viewer/internal/models"
	"time"

	"github.com/google/uuid"
)

func MapSalesToVault(raw models.RawSalesData) models.VehicleDigitalVault {
	return models.VehicleDigitalVault{
		ID:           uuid.New(),
		ExternalID:   fmt.Sprintf("SALES-%d", raw.ID), 
		VIN:          raw.VehicleVIN,
		SourceSystem: "SALES",
		Title:        fmt.Sprintf("Sales Contract: %s", raw.DocumentType),
		DocCategory:  "Commercial", 
		EventDate:    raw.CreatedAt, 
		FileURL:      raw.FileURL,
		SyncedAt:     time.Now(),
	}
}

func MapServiceToVault(raw models.RawServiceData) models.VehicleDigitalVault {
	return models.VehicleDigitalVault{
		ID:           uuid.New(),
		ExternalID:   fmt.Sprintf("SERVICE-%d", raw.ID),
		VIN:          raw.VehicleVIN,
		SourceSystem: "SERVICE",
		Title:        fmt.Sprintf("Service Report: %s", raw.ServiceType),
		DocCategory:  "Technical", 
		EventDate:    raw.CompletionDate, 
		FileURL:      raw.ReportLink,
		SyncedAt:     time.Now(),
	}
}