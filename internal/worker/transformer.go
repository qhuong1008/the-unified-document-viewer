package worker

import (
	"strings"
	"the-unified-document-viewer/internal/models"
	"time"

	"github.com/google/uuid"
)

func sanitizeUTF8(s string) string {
	return strings.ToValidUTF8(s, " ")
}

func MapSalesToVault(raw models.RawSalesData) models.VehicleDigitalVault {
	return models.VehicleDigitalVault{
		ID:           uuid.New(),
		ExternalID:   sanitizeUTF8(raw.ID),
		SourceSystem: string(SourceSales),
		VIN:          sanitizeUTF8(raw.VehicleVIN),
		Title:        sanitizeUTF8("Sales Contract"),
		DocCategory:  "Commercial",
		SalesPerson:  sanitizeUTF8(raw.SalesPerson),
		SalesDocumentIssueDate:  raw.CreatedAt,
		FileURL:      sanitizeUTF8(raw.FileURL),
		SyncedAt:     time.Now(),
	}
}

func MapServiceToVault(raw models.RawServiceData) models.VehicleDigitalVault {
	return models.VehicleDigitalVault{
		ID:           uuid.New(),
		ExternalID:   sanitizeUTF8(raw.ID),
		SourceSystem: string(SourceService),
		VIN:          sanitizeUTF8(raw.VehicleVIN),
		Title:        sanitizeUTF8("Service Report"),
		DocCategory:  "Technical",
		Technician:   sanitizeUTF8(raw.Technician),
		ServiceCompletionDate:  raw.CompletionDate,
		FileURL:      sanitizeUTF8(raw.ReportLink),
		SyncedAt:     time.Now(),
	}
}

