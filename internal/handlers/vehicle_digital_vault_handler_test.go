package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockAdapter struct {
	salesData  models.RawSalesData
	serviceData models.RawServiceData
}

func (m *mockAdapter) FetchSalesByVIN(ctx context.Context, vin string) (models.RawSalesData, error) {
	return m.salesData, nil
}

func (m *mockAdapter) FetchServiceByVIN(ctx context.Context, vin string) (models.RawServiceData, error) {
	return m.serviceData, nil
}

func setupTestHandler() (*gin.Engine, *VehicleDigitalVaultHandler, func()) {
	r := gin.Default()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.VehicleDigitalVault{})

	repo := repository.NewPostgresRepository(db)
	jobQueue := make(chan worker.Job, 10)
	
	mock := &mockAdapter{}
	handler := &VehicleDigitalVaultHandler{
		Repo: repo,
		JobQueue: jobQueue,
		Adapter: mock,
	}

	r.POST("/vault/search", handler.SearchAndSyncByVIN)

	return r, handler, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestSearchAndSyncByVIN(t *testing.T) {
	r, handler, teardown := setupTestHandler()
	defer teardown()

	// Mock data
	mock := handler.Adapter.(*mockAdapter)
	mock.salesData = models.RawSalesData{
		ID: "test-sales-1",
		VehicleVIN: "TESTVIN123",
		DocumentType: "Test Contract",
		SalesPerson: "Test Sales",
		FileURL: "test-sales.pdf",
	}

	mock.serviceData = models.RawServiceData{
		ID: "test-service-1",
		VehicleVIN: "TESTVIN123",
		ServiceType: "Test Service",
		Technician: "Test Tech",
		ReportLink: "test-service.pdf",
	}

	// Test request
	body, _ := json.Marshal(map[string]string{"vin": "TESTVIN123"})
	req := httptest.NewRequest("POST", "/vault/search", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["exists"].(bool), "Should have records after sync")
	assert.Equal(t, float64(2), response["total_documents"].(float64), "Should have 2 documents")
	assert.True(t, response["sales_fetched"].(bool))
	assert.True(t, response["service_fetched"].(bool))
}

func TestSearchAndSyncByVINInvalidJSON(t *testing.T) {
	r, _, teardown := setupTestHandler()
	defer teardown()

	req := httptest.NewRequest("POST", "/vault/search", bytes.NewReader([]byte(`invalid json`)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchAndSyncByVINEmptyVin(t *testing.T) {
	r, _, teardown := setupTestHandler()
	defer teardown()

	body, _ := json.Marshal(map[string]string{"vin": ""})
	req := httptest.NewRequest("POST", "/vault/search", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

