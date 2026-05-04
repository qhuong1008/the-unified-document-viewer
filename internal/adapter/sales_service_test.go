package adapter

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"

	"the-unified-document-viewer/internal/utils"
)

func TestFetchSalesByVIN(t *testing.T) {
	adapter := NewSalesServiceAdapter("https://test.com")
	ctx := context.Background()

	data, err := adapter.FetchSalesByVIN(ctx, "1HGCM82633A004352")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data.VehicleVIN != "1HGCM82633A004352" {
		t.Errorf("Expected VIN 1HGCM82633A004352, got %s", data.VehicleVIN)
	}

	if data.ID == "" || !strings.HasPrefix(data.ID, "sales-") {
		t.Errorf("Expected ID starting with 'sales-', got %s", data.ID)
	}

	if data.DocumentType != "Sales Contract" {
		t.Errorf("Expected 'Sales Contract', got %s", data.DocumentType)
	}

	if data.SalesPerson != "John Smith" {
		t.Errorf("Expected 'John Smith', got %s", data.SalesPerson)
	}

	if data.FileURL == "" {
		t.Error("Expected non-empty FileURL from random sales PDF")
	}
}

func TestFetchServiceByVIN(t *testing.T) {
	adapter := NewSalesServiceAdapter("https://test.com")
	ctx := context.Background()

	data, err := adapter.FetchServiceByVIN(ctx, "1HGCM82633A004352")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data.VehicleVIN != "1HGCM82633A004352" {
		t.Errorf("Expected VIN 1HGCM82633A004352, got %s", data.VehicleVIN)
	}

	if data.ID == "" || !strings.HasPrefix(data.ID, "service-") {
		t.Errorf("Expected ID starting with 'service-', got %s", data.ID)
	}

	if data.ServiceType != "Oil Change & Inspection" {
		t.Errorf("Expected 'Oil Change & Inspection', got %s", data.ServiceType)
	}

	if data.Technician != "Mike Johnson" {
		t.Errorf("Expected 'Mike Johnson', got %s", data.Technician)
	}

	if data.ReportLink == "" {
		t.Error("Expected non-empty ReportLink from random service PDF")
	}
}

func TestRandomPDFLinks(t *testing.T) {
	// Test that random functions return valid URLs
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	url1 := utils.RandomPDFForSales()
	url2 := utils.RandomPDFForService()

	if url1 == "" {
		t.Error("RandomPDFForSales returned empty URL")
	}

	if url2 == "" {
		t.Error("RandomPDFForService returned empty URL")
	}

	// Check they are from PDFLinks slice
	found1 := false
	for _, link := range utils.PDFLinks {
		if link == url1 {
			found1 = true
			break
		}
	}
	if !found1 {
		t.Errorf("Sales URL not in PDFLinks: %s", url1)
	}

	found2 := false
	for _, link := range utils.PDFLinks {
		if link == url2 {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Errorf("Service URL not in PDFLinks: %s", url2)
	}
}

