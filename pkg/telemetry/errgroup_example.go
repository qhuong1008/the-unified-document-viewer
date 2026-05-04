package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

// ============================================================
// REAL-WORLD USAGE EXAMPLES FOR THE PROJECT
// ============================================================

// FetchAllVehicleData fetches data from Sales API, Service API, and VIN database in parallel
// This is the main use case for "The Unified Document Viewer" - collecting vehicle data from multiple sources
func FetchAllVehicleData(ctx context.Context, vin, salesAPIEndpoint, serviceAPIEndpoint string) (*VehicleDataResult, error) {
	tracer := GetTracer()

	// Parent span for the entire operation
	ctx, parentSpan := tracer.Start(ctx, "fetch-all-vehicle-data",
		trace.WithAttributes(
			attribute.String("vin", vin),
			attribute.String("operation.type", "parallel-fetch"),
		),
	)
	defer parentSpan.End()

	g := errgroup.Group{}

	var (
		salesData  *SalesData
		serviceData *ServiceData
		vinData    *VINData
		mu         sync.Mutex
	)

	// ============================================================
	// Parallel Call 1: Sales API
	// ============================================================
	g.Go(func() error {
		salesCtx, span := tracer.Start(ctx, "sales-api-call",
			trace.WithAttributes(
				attribute.String("api.name", "sales-partner"),
				attribute.String("api.endpoint", salesAPIEndpoint),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		startTime := time.Now()

		// Build request with VIN
		url := fmt.Sprintf("%s/vehicles/%s", salesAPIEndpoint, vin)
		req, err := http.NewRequestWithContext(salesCtx, "GET", url, nil)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		// Inject trace context for downstream visibility
		injectTraceContextToHeaders(salesCtx, req.Header)

		resp, err := http.DefaultClient.Do(req)
		duration := time.Since(startTime)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(
				attribute.Int64("duration_ms", duration.Milliseconds()),
				attribute.String("error.type", "network-error"),
			)
			return fmt.Errorf("Sales API call failed: %w", err)
		}
		defer resp.Body.Close()

		// Parse response
		salesData = &SalesData{
			VIN:       vin,
			Timestamp: time.Now(),
		}
		json.NewDecoder(resp.Body).Decode(salesData)

		span.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		mu.Lock()
		defer mu.Unlock()
		return nil
	})

	// ============================================================
	// Parallel Call 2: Service API
	// ============================================================
	g.Go(func() error {
		serviceCtx, span := tracer.Start(ctx, "service-api-call",
			trace.WithAttributes(
				attribute.String("api.name", "service-partner"),
				attribute.String("api.endpoint", serviceAPIEndpoint),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		startTime := time.Now()

		url := fmt.Sprintf("%s/vehicles/%s/service", serviceAPIEndpoint, vin)
		req, err := http.NewRequestWithContext(serviceCtx, "GET", url, nil)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		injectTraceContextToHeaders(serviceCtx, req.Header)

		resp, err := http.DefaultClient.Do(req)
		duration := time.Since(startTime)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(
				attribute.Int64("duration_ms", duration.Milliseconds()),
				attribute.String("error.type", "network-error"),
			)
			return fmt.Errorf("Service API call failed: %w", err)
		}
		defer resp.Body.Close()

		serviceData = &ServiceData{
			VIN:       vin,
			Timestamp: time.Now(),
		}
		json.NewDecoder(resp.Body).Decode(serviceData)

		span.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		mu.Lock()
		defer mu.Unlock()
		return nil
	})

	// ============================================================
	// Parallel Call 3: VIN Database Lookup
	// ============================================================
	g.Go(func() error {
		_, span := tracer.Start(ctx, "vin-database-lookup",
			trace.WithAttributes(
				attribute.String("database.type", "postgres"),
				attribute.String("operation", "vin-lookup"),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Simulate database lookup with timing
		startTime := time.Now()

		// In real implementation, this would be a database call:
		// err := db.QueryRow("SELECT * FROM vehicles WHERE vin = $1", vin).Scan(&vinData...)

		// Simulating the lookup
		vinData = &VINData{
			VIN:   vin,
			Make:  "Toyota",
			Model: "Camry",
			Year: 2024,
		}

		duration := time.Since(startTime)

		span.SetAttributes(
			attribute.Int64("duration_ms", duration.Milliseconds()),
			attribute.String("database.operation", "query"),
		)

		mu.Lock()
		defer mu.Unlock()
		return nil
	})

	// Wait for all parallel calls to complete
	if err := g.Wait(); err != nil {
		parentSpan.SetStatus(codes.Error, "One or more API calls failed")
		parentSpan.RecordError(err)
		return nil, err
	}

	// Combine results
	result := &VehicleDataResult{
		Sales:   salesData,
		Service: serviceData,
		VIN:     vinData,
		Fetched: time.Now(),
	}

	parentSpan.SetAttributes(
		attribute.Bool("sales_api_success", salesData != nil),
		attribute.Bool("service_api_success", serviceData != nil),
		attribute.Bool("vin_lookup_success", vinData != nil),
	)

	return result, nil
}

// injectTraceContextToHeaders injects trace context into HTTP headers for downstream services
func injectTraceContextToHeaders(ctx context.Context, headers http.Header) {
	carrier := headerCarrier{headers}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

type headerCarrier struct {
	http.Header
}

func (h headerCarrier) Get(key string) string {
	if v := h.Header[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h headerCarrier) Set(key, value string) {
	h.Header[key] = []string{value}
}

func (h headerCarrier) Keys() []string {
	keys := make([]string, 0, len(h.Header))
	for k := range h.Header {
		keys = append(keys, k)
	}
	return keys
}

// ============================================================
// DATA STRUCTURES
// ============================================================

type SalesData struct {
	VIN       string    `json:"vin"`
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price,omitempty"`
	Dealer    string    `json:"dealer,omitempty"`
	Location string    `json:"location,omitempty"`
}

type ServiceData struct {
	VIN         string    `json:"vin"`
	Timestamp   time.Time `json:"timestamp"`
	LastService time.Time `json:"last_service,omitempty"`
	Mileage     int       `json:"mileage,omitempty"`
	Notes      string    `json:"notes,omitempty"`
}

type VINData struct {
	VIN   string `json:"vin"`
	Make string `json:"make"`
	Model string `json:"model"`
	Year int    `json:"year"`
}

type VehicleDataResult struct {
	Sales   *SalesData
	Service *ServiceData
	VIN     *VINData
	Fetched time.Time
}

// ============================================================
// LEGACY EXAMPLES (kept for reference)
// ============================================================

// ParallelAPICall demonstrates how to use otel.Tracer to create child spans
// inside an errgroup when making parallel calls to external services.
// This is the key demonstration for the interview requirement:
// "monitoring the performance of each third-party service independently"
func ParallelAPICall(ctx context.Context, salesAPIEndpoint, serviceAPIEndpoint string) error {
	tracer := GetTracer()

	ctx, parentSpan := tracer.Start(ctx, "parallel-api-calls",
		trace.WithAttributes(
			attribute.String("operation", "fetch-external-data"),
			attribute.String("sales_api", salesAPIEndpoint),
			attribute.String("service_api", serviceAPIEndpoint),
		),
	)
	defer parentSpan.End()

	g := errgroup.Group{}

	var mu sync.Mutex

	// Launch parallel call to Sales API
	g.Go(func() error {
		_, salesSpan := tracer.Start(ctx, "Sales API call",
			trace.WithAttributes(
				attribute.String("api.type", "sales"),
				attribute.String("api.endpoint", salesAPIEndpoint),
				attribute.String("api.library", "net/http"),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer salesSpan.End()

		startTime := time.Now()
		resp, err := http.DefaultClient.Get(salesAPIEndpoint)
		duration := time.Since(startTime)

		if err != nil {
			salesSpan.SetStatus(codes.Error, err.Error())
			salesSpan.RecordError(err)
			salesSpan.SetAttributes(
				attribute.String("error.type", "network-error"),
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)
			return fmt.Errorf("Sales API call failed: %w", err)
		}
		defer resp.Body.Close()

		salesSpan.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		_ = mu
		return nil
	})

	// Launch parallel call to Service API
	g.Go(func() error {
		_, serviceSpan := tracer.Start(ctx, "Service API call",
			trace.WithAttributes(
				attribute.String("api.type", "service"),
				attribute.String("api.endpoint", serviceAPIEndpoint),
				attribute.String("api.library", "net/http"),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer serviceSpan.End()

		startTime := time.Now()
		resp, err := http.DefaultClient.Get(serviceAPIEndpoint)
		duration := time.Since(startTime)

		if err != nil {
			serviceSpan.SetStatus(codes.Error, err.Error())
			serviceSpan.RecordError(err)
			serviceSpan.SetAttributes(
				attribute.String("error.type", "network-error"),
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)
			return fmt.Errorf("Service API call failed: %w", err)
		}
		defer resp.Body.Close()

		serviceSpan.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		_ = mu
		return nil
	})

	if err := g.Wait(); err != nil {
		parentSpan.SetStatus(codes.Error, "One or more API calls failed")
		parentSpan.RecordError(err)
		return err
	}

	parentSpan.SetAttributes(
		attribute.Bool("sales_api_success", true),
		attribute.Bool("service_api_success", true),
	)

	return nil
}

// ParallelVendorFetch fetches data from multiple vendor APIs in parallel
// and creates isolated spans for each with full observability
func ParallelVendorFetch(
	ctx context.Context,
	tracer trace.Tracer,
	calls []ExternalAPICall,
) (map[string][]byte, error) {

	results := make(map[string][]byte)
	var mu sync.Mutex
	var wg errgroup.Group

	for _, call := range calls {
		wg.Go(func() error {
			callCtx, span := tracer.Start(ctx, fmt.Sprintf("vendor.%s", call.Name),
				trace.WithAttributes(
					attribute.String("vendor.name", call.Name),
					attribute.String("vendor.endpoint", call.Endpoint),
					attribute.String("http.method", call.Method),
				),
				trace.WithSpanKind(trace.SpanKindClient),
			)
			defer span.End()

			req, err := http.NewRequestWithContext(callCtx, call.Method, call.Endpoint, nil)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				return err
			}

			carrier := httpHeaderCarrier{req.Header}
			otel.GetTextMapPropagator().Inject(callCtx, carrier)

			startTime := time.Now()
			resp, err := http.DefaultClient.Do(req)
			duration := time.Since(startTime)

			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.SetAttributes(
					attribute.String("error.message", err.Error()),
				)
				return err
			}
			defer resp.Body.Close()

			buf := make([]byte, resp.ContentLength)
			n, _ := resp.Body.Read(buf)

			mu.Lock()
			results[call.Name] = buf[:n]
			mu.Unlock()

			span.SetAttributes(
				attribute.Int("http.status_code", resp.StatusCode),
				attribute.Int64("duration_ms", duration.Milliseconds()),
				attribute.Int("response_bytes", n),
			)

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return results, err
	}

	return results, nil
}

type ExternalAPICall struct {
	Name     string
	Endpoint string
	Method   string
}

type httpHeaderCarrier struct {
	http.Header
}

func (h httpHeaderCarrier) Get(key string) string {
	if v := h.Header[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h httpHeaderCarrier) Set(key, value string) {
	h.Header[key] = []string{value}
}

func (h httpHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(h.Header))
	for k := range h.Header {
		keys = append(keys, k)
	}
	return keys
}
