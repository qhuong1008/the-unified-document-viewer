package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/utils"
	"the-unified-document-viewer/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SalesServiceAdapter struct {
	baseURL string
	client  *http.Client
}

func NewSalesServiceAdapter(baseURL string) *SalesServiceAdapter {
	return &SalesServiceAdapter{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// FetchSalesData fetches sales data from API (or mock)
func (s *SalesServiceAdapter) FetchSalesData() (map[string]interface{}, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/sales", s.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// FetchServiceData fetches service data from API (or mock)
func (s *SalesServiceAdapter) FetchServiceData() (map[string]interface{}, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/api/service", s.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// FetchSalesByVIN mocks fetching sales data for specific VIN (for VIN search API)
func (s *SalesServiceAdapter) FetchSalesByVIN(ctx context.Context, vin string) (models.RawSalesData, error) {
	ctx, span := telemetry.GetTracer().Start(ctx, "sales.fetch-vin",
		trace.WithAttributes(
			attribute.String("vin", vin),
			attribute.String("api.endpoint", fmt.Sprintf("%s/vehicles/%s", s.baseURL, vin)),
		),
	)
	defer span.End()

	// Mock delay to simulate API call
	time.Sleep(500 * time.Millisecond)

	// Mock realistic sales data for VIN
	mockID := fmt.Sprintf("sales-%s", vin[:min(len(vin), 8)]) // Safe with safe_vin.min
	if len(vin) < 17 {
		return models.RawSalesData{}, fmt.Errorf("invalid VIN length: %d, expected 17 characters", len(vin))
	}

	data := models.RawSalesData{
		ID:           mockID,
		VehicleVIN:   vin,
		DocumentType: "Sales Contract",
		SalesPerson:  "John Smith",
		CreatedAt:    time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		FileURL:      utils.RandomPDFForSales(),
	};

	fmt.Println("sales data: ", data)

	return data, nil
}

// FetchServiceByVIN mocks fetching service data for specific VIN (for VIN search API)
func (s *SalesServiceAdapter) FetchServiceByVIN(ctx context.Context, vin string) (models.RawServiceData, error) {
	ctx, span := telemetry.GetTracer().Start(ctx, "service.fetch-vin",
		trace.WithAttributes(
			attribute.String("vin", vin),
			attribute.String("api.endpoint", fmt.Sprintf("%s/vehicles/%s/service", s.baseURL, vin)),
		),
	)
	defer span.End()

	// Mock delay to simulate API call
	time.Sleep(300 * time.Millisecond)

	// Mock realistic service data for VIN
	mockID := fmt.Sprintf("service-%s", vin[:min(len(vin), 8)]) // Safe with safe_vin.min
	if len(vin) < 17 {
			return models.RawServiceData{}, fmt.Errorf("invalid VIN length: %d, expected 17 characters", len(vin))
	}

	data := models.RawServiceData{
		ID:             mockID,
		VehicleVIN:     vin,
		ServiceType:    "Oil Change & Inspection",
		Technician:     "Mike Johnson",
		CompletionDate: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
		ReportLink:     utils.RandomPDFForService(),
	}
	fmt.Println("service data: ", data)

	return data, nil
}

