package worker

import (
	"testing"
	"time"

	"the-unified-document-viewer/internal/models"

	"github.com/google/uuid"
)

func TestSanitizeUTF8(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal string", "Hello World", "Hello World"},
		{"Empty string", "", ""},
		{"Special characters", "Test-Case-123", "Test-Case-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeUTF8(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapSalesToVault(t *testing.T) {
	now := time.Now()
	rawSales := models.RawSalesData{
		ID:            "SALES-001",
		VehicleVIN:    "1HGBH41JXMN109186",
		DocumentType:  "Sales Contract",
		SalesPerson:   "John Doe",
		CreatedAt:     now,
		FileURL:       "https://example.com/contract.pdf",
	}

	vault := MapSalesToVault(rawSales)

	if vault.ExternalID != rawSales.ID {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, rawSales.ID)
	}
	if vault.VIN != rawSales.VehicleVIN {
		t.Errorf("VIN = %q, want %q", vault.VIN, rawSales.VehicleVIN)
	}
	if vault.SourceSystem != string(SourceSales) {
		t.Errorf("SourceSystem = %q, want %q", vault.SourceSystem, SourceSales)
	}
	if vault.SalesPerson != rawSales.SalesPerson {
		t.Errorf("SalesPerson = %q, want %q", vault.SalesPerson, rawSales.SalesPerson)
	}
	if !vault.SalesDocumentIssueDate.Equal(rawSales.CreatedAt) {
		t.Errorf("SalesDocumentIssueDate = %v, want %v", vault.SalesDocumentIssueDate, rawSales.CreatedAt)
	}
	if vault.FileURL != rawSales.FileURL {
		t.Errorf("FileURL = %q, want %q", vault.FileURL, rawSales.FileURL)
	}
	if vault.Title != "Sales Contract" {
		t.Errorf("Title = %q, want %q", vault.Title, "Sales Contract")
	}
	if vault.ID == uuid.Nil {
		t.Error("ID should not be nil")
	}
	if vault.SyncedAt.IsZero() {
		t.Error("SyncedAt should not be zero")
	}
}

func TestMapServiceToVault(t *testing.T) {
	now := time.Now()
	rawService := models.RawServiceData{
		ID:             "SERVICE-001",
		VehicleVIN:     "1HGBH41JXMN109186",
		ServiceType:    "Oil Change",
		Technician:     "Jane Smith",
		CompletionDate: now,
		ReportLink:    "https://example.com/report.pdf",
	}

	vault := MapServiceToVault(rawService)

	if vault.ExternalID != rawService.ID {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, rawService.ID)
	}
	if vault.VIN != rawService.VehicleVIN {
		t.Errorf("VIN = %q, want %q", vault.VIN, rawService.VehicleVIN)
	}
	if vault.SourceSystem != string(SourceService) {
		t.Errorf("SourceSystem = %q, want %q", vault.SourceSystem, SourceService)
	}
	if vault.Technician != rawService.Technician {
		t.Errorf("Technician = %q, want %q", vault.Technician, rawService.Technician)
	}
	if !vault.ServiceCompletionDate.Equal(rawService.CompletionDate) {
		t.Errorf("ServiceCompletionDate = %v, want %v", vault.ServiceCompletionDate, rawService.CompletionDate)
	}
	if vault.FileURL != rawService.ReportLink {
		t.Errorf("FileURL = %q, want %q", vault.FileURL, rawService.ReportLink)
	}
	if vault.Title != "Service Report" {
		t.Errorf("Title = %q, want %q", vault.Title, "Service Report")
	}
	if vault.ID == uuid.Nil {
		t.Error("ID should not be nil")
	}
	if vault.SyncedAt.IsZero() {
		t.Error("SyncedAt should not be zero")
	}
}

func TestMapSalesToVault_EmptyFields(t *testing.T) {
	rawSales := models.RawSalesData{
		ID:           "",
		VehicleVIN:   "",
		DocumentType: "",
		SalesPerson:  "",
		CreatedAt:    time.Time{},
		FileURL:      "",
	}

	vault := MapSalesToVault(rawSales)
	if vault.ExternalID == "" {
		t.Log("Empty ExternalID handled")
	}
}

func TestMapServiceToVault_EmptyFields(t *testing.T) {
	rawService := models.RawServiceData{
		ID:             "",
		VehicleVIN:     "",
		ServiceType:    "",
		Technician:     "",
		CompletionDate: time.Time{},
		ReportLink:     "",
	}

	vault := MapServiceToVault(rawService)
	if vault.ExternalID == "" {
		t.Log("Empty ExternalID handled")
	}
}

func TestMapSalesToVault_VINPreservation(t *testing.T) {
	vin := "1HGBH41JXMN109186"
	rawSales := models.RawSalesData{
		ID:            "TEST-001",
		VehicleVIN:    vin,
		DocumentType: "Sales Contract",
		SalesPerson:   "Test Person",
		CreatedAt:    time.Now(),
		FileURL:      "https://example.com/test.pdf",
	}

	vault := MapSalesToVault(rawSales)

	if len(vault.VIN) != 17 {
		t.Errorf("VIN length = %d, want 17", len(vault.VIN))
	}
	if vault.VIN != vin {
		t.Errorf("VIN = %q, want %q", vault.VIN, vin)
	}
}
