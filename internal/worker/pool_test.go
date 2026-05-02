package worker

import (
	"errors"
	"testing"
	"time"

	"the-unified-document-viewer/internal/models"
)

type MockRepository struct {
	upsertFunc   func(data models.VehicleDigitalVault) error
	getByVINFunc func(vin string) ([]models.VehicleDigitalVault, error)
	checkVINFunc func(vin string) (bool, int64, error)
}

func (m *MockRepository) UpsertVehicleDigitalVaultRecord(data models.VehicleDigitalVault) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(data)
	}
	return nil
}

func (m *MockRepository) UpsertVaultRecord(data models.VehicleDigitalVault) error {
	return m.UpsertVehicleDigitalVaultRecord(data)
}

func (m *MockRepository) GetVehicleDigitalVaultByVIN(vin string) ([]models.VehicleDigitalVault, error) {
	if m.getByVINFunc != nil {
		return m.getByVINFunc(vin)
	}
	return nil, nil
}

func (m *MockRepository) CheckIfVINExists(vin string) (bool, int64, error) {
	if m.checkVINFunc != nil {
		return m.checkVINFunc(vin)
	}
	return false, 0, nil
}

func TestTransformJobPayload_Sales(t *testing.T) {
	job := Job{
		Type: SourceSales,
		Payload: models.RawSalesData{
			ID:            "TEST-SALES-001",
			VehicleVIN:    "1HGBH41JXMN109186",
			DocumentType:  "Sales Contract",
			SalesPerson:   "John Doe",
			CreatedAt:     time.Now(),
			FileURL:       "https://example.com/contract.pdf",
		},
	}

	vault, ok := transformJobPayload(job)

	if !ok {
		t.Error("returned ok=false, want true")
	}
	if vault.ExternalID != "TEST-SALES-001" {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, "TEST-SALES-001")
	}
	if vault.VIN != "1HGBH41JXMN109186" {
		t.Errorf("VIN = %q, want %q", vault.VIN, "1HGBH41JXMN109186")
	}
	if vault.SourceSystem != string(SourceSales) {
		t.Errorf("SourceSystem = %q, want %q", vault.SourceSystem, SourceSales)
	}
}

func TestTransformJobPayload_Service(t *testing.T) {
	job := Job{
		Type: SourceService,
		Payload: models.RawServiceData{
			ID:             "TEST-SERVICE-001",
			VehicleVIN:     "1HGBH41JXMN109186",
			ServiceType:    "Oil Change",
			Technician:      "Jane Smith",
			CompletionDate: time.Now(),
			ReportLink:     "https://example.com/report.pdf",
		},
	}

	vault, ok := transformJobPayload(job)

	if !ok {
		t.Error("returned ok=false, want true")
	}
	if vault.ExternalID != "TEST-SERVICE-001" {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, "TEST-SERVICE-001")
	}
	if vault.VIN != "1HGBH41JXMN109186" {
		t.Errorf("VIN = %q, want %q", vault.VIN, "1HGBH41JXMN109186")
	}
	if vault.SourceSystem != string(SourceService) {
		t.Errorf("SourceSystem = %q, want %q", vault.SourceSystem, SourceService)
	}
}

func TestTransformJobPayload_InvalidType(t *testing.T) {
	job := Job{
		Type:    "INVALID",
		Payload: nil,
	}

	_, ok := transformJobPayload(job)

	if ok {
		t.Error("returned ok=true for invalid job type, want false")
	}
}

func TestTransformJobPayload_WrongPayloadType(t *testing.T) {
	job := Job{
		Type:    SourceSales,
		Payload: "invalid payload",
	}

	_, ok := transformJobPayload(job)

	if ok {
		t.Error("returned ok=true for wrong payload type, want false")
	}
}

func TestExecuteJob_SalesType(t *testing.T) {
	job := Job{
		Type: SourceSales,
		Payload: models.RawSalesData{
			ID:            "EXEC-SALES-001",
			VehicleVIN:    "1HGBH41JXMN109186",
			DocumentType:  "Sales Contract",
			SalesPerson:   "Test Person",
			CreatedAt:     time.Now(),
			FileURL:       "https://example.com/test.pdf",
		},
	}

	vault, ok := transformJobPayload(job)
	if !ok {
		t.Error("transformJobPayload failed for sales job")
	}
	if vault.ExternalID != "EXEC-SALES-001" {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, "EXEC-SALES-001")
	}
}

func TestExecuteJob_ServiceType(t *testing.T) {
	job := Job{
		Type: SourceService,
		Payload: models.RawServiceData{
			ID:             "EXEC-SERVICE-001",
			VehicleVIN:     "1HGBH41JXMN109186",
			ServiceType:    "Brake Repair",
			Technician:     "Test Tech",
			CompletionDate: time.Now(),
			ReportLink:    "https://example.com/test.pdf",
		},
	}

	vault, ok := transformJobPayload(job)
	if !ok {
		t.Error("transformJobPayload failed for service job")
	}
	if vault.ExternalID != "EXEC-SERVICE-001" {
		t.Errorf("ExternalID = %q, want %q", vault.ExternalID, "EXEC-SERVICE-001")
	}
}

func TestExecuteJob_InvalidPayload(t *testing.T) {
	job := Job{
		Type:    SourceSales,
		Payload: "invalid",
	}

	_, ok := transformJobPayload(job)
	if ok {
		t.Error("should return false for invalid payload type")
	}
}

func TestJobType_Constants(t *testing.T) {
	if SourceSales != "SALES" {
		t.Errorf("SourceSales = %q, want %q", SourceSales, "SALES")
	}
	if SourceService != "SERVICE" {
		t.Errorf("SourceService = %q, want %q", SourceService, "SERVICE")
	}
}

func TestJob_PayloadInterface(t *testing.T) {
	salesJob := Job{
		Type:    SourceSales,
		Payload: models.RawSalesData{ID: "PAYLOAD-1"},
	}

	serviceJob := Job{
		Type:    SourceService,
		Payload: models.RawServiceData{ID: "PAYLOAD-2"},
	}

	if _, ok := salesJob.Payload.(models.RawSalesData); !ok {
		t.Error("Sales job payload should be RawSalesData")
	}
	if _, ok := serviceJob.Payload.(models.RawServiceData); !ok {
		t.Error("Service job payload should be RawServiceData")
	}
}

var testDBError = errors.New("database error")

func TestRepository_UpsertError(t *testing.T) {
	repo := &MockRepository{
		upsertFunc: func(data models.VehicleDigitalVault) error {
			return testDBError
		},
	}

	vault := models.VehicleDigitalVault{
		ExternalID: "ERROR-TEST-001",
		VIN:        "1HGBH41JXMN109186",
	}

	err := repo.UpsertVehicleDigitalVaultRecord(vault)
	if err != testDBError {
		t.Errorf("error = %v, want %v", err, testDBError)
	}
}

func TestGetVehicleDigitalVaultByVIN_Ordering(t *testing.T) {
	repo := &MockRepository{
		getByVINFunc: func(vin string) ([]models.VehicleDigitalVault, error) {
			return []models.VehicleDigitalVault{
				{ExternalID: "DOC-1", VIN: vin, SyncedAt: time.Now().Add(-2 * time.Hour)},
				{ExternalID: "DOC-2", VIN: vin, SyncedAt: time.Now().Add(-1 * time.Hour)},
				{ExternalID: "DOC-3", VIN: vin, SyncedAt: time.Now()},
			}, nil
		},
	}

	docs, err := repo.GetVehicleDigitalVaultByVIN("1HGBH41JXMN109186")
	if err != nil {
		t.Errorf("error = %v, want nil", err)
	}
	if len(docs) != 3 {
		t.Errorf("returned %d documents, want 3", len(docs))
	}
}

func TestCheckIfVINExists_Count(t *testing.T) {
	repo := &MockRepository{
		checkVINFunc: func(vin string) (bool, int64, error) {
			return true, 5, nil
		},
	}

	exists, count, err := repo.CheckIfVINExists("1HGBH41JXMN109186")
	if err != nil {
		t.Errorf("error = %v, want nil", err)
	}
	if !exists {
		t.Error("returned exists=false, want true")
	}
	if count != 5 {
		t.Errorf("count = %d, want 5", count)
	}
}
