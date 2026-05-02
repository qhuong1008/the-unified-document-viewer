package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandleSalesWebhook_ValidPayload(t *testing.T) {
	jobQueue := make(chan worker.Job, 1)
	handler := &WebhookHandler{JobQueue: jobQueue}

	payload := models.RawSalesData{
		ID:            "SALES-TEST-001",
		VehicleVIN:    "1HGBH41JXMN109186",
		DocumentType:  "Sales Contract",
		SalesPerson:   "Test Salesperson",
		FileURL:       "https://example.com/contract.pdf",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/webhooks/sales", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.HandleSalesWebhook(c)

	if w.Code != http.StatusAccepted {
		t.Errorf("status = %d, want %d", w.Code, http.StatusAccepted)
	}

	select {
	case job := <-jobQueue:
		if job.Type != worker.SourceSales {
			t.Errorf("Job type = %q, want %q", job.Type, worker.SourceSales)
		}
		if _, ok := job.Payload.(models.RawSalesData); !ok {
			t.Error("Job payload should be RawSalesData type")
		}
	default:
		t.Error("Job was not pushed to queue")
	}
}

func TestHandleSalesWebhook_InvalidPayload(t *testing.T) {
	jobQueue := make(chan worker.Job, 1)
	handler := &WebhookHandler{JobQueue: jobQueue}

	invalidPayload := []byte(`{invalid json`)
	body := bytes.NewBuffer(invalidPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/webhooks/sales", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler.HandleSalesWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleServiceWebhook_ValidPayload(t *testing.T) {
	jobQueue := make(chan worker.Job, 1)
	handler := &WebhookHandler{JobQueue: jobQueue}

	payload := models.RawServiceData{
		ID:             "SERVICE-TEST-001",
		VehicleVIN:     "1HGBH41JXMN109186",
		ServiceType:    "Oil Change",
		Technician:     "Test Technician",
		ReportLink:     "https://example.com/report.pdf",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/webhooks/service", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.HandleServiceWebhook(c)

	if w.Code != http.StatusAccepted {
		t.Errorf("status = %d, want %d", w.Code, http.StatusAccepted)
	}

	select {
	case job := <-jobQueue:
		if job.Type != worker.SourceService {
			t.Errorf("Job type = %q, want %q", job.Type, worker.SourceService)
		}
		if _, ok := job.Payload.(models.RawServiceData); !ok {
			t.Error("Job payload should be RawServiceData type")
		}
	default:
		t.Error("Job was not pushed to queue")
	}
}

func TestHandleServiceWebhook_InvalidPayload(t *testing.T) {
	jobQueue := make(chan worker.Job, 1)
	handler := &WebhookHandler{JobQueue: jobQueue}

	invalidPayload := []byte(`{invalid json`)
	body := bytes.NewBuffer(invalidPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/webhooks/service", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler.HandleServiceWebhook(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestWebhookHandler_ConcurrentProcessing(t *testing.T) {
	jobQueue := make(chan worker.Job, 100)
	handler := &WebhookHandler{JobQueue: jobQueue}

	for i := 0; i < 10; i++ {
		payload := models.RawSalesData{
			ID:            "SALES-CONCURRENT-001",
			VehicleVIN:    "1HGBH41JXMN109186",
			DocumentType: "Sales Contract",
			SalesPerson:   "Test Person",
			FileURL:      "https://example.com/test.pdf",
		}
		body, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/webhooks/sales", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		go handler.HandleSalesWebhook(c)
	}

	count := 0
	for count < 10 {
		select {
		case <-jobQueue:
			count++
		}
	}

	if count != 10 {
		t.Errorf("Expected 10 jobs, got %d", count)
	}
}
