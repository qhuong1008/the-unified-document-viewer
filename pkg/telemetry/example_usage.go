package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ============================================================
// HOW TO USE THE TELEMETRY PACKAGE - COMPLETE GUIDE
// ============================================================

// This file demonstrates how to integrate and use the telemetry package
// in your Go applications.

// --------------------------------------------------------
// EXAMPLE 1: How to set up Telemetry in main.go
// --------------------------------------------------------

func Example_initializeTelemetry() {
	ctx := context.Background()

	// Initialize the tracer - call this at the start of your application
	tp, err := InitTracer(ctx, "my-service-name")
	if err != nil {
		fmt.Printf("Failed to initialize telemetry: %v\n", err)
		return
	}
	defer tp.Shutdown(ctx)

	fmt.Println("Telemetry initialized successfully!")
	// Output: Telemetry initialized successfully!
}

// --------------------------------------------------------
// EXAMPLE 2: How to add the Middleware to Gin
// --------------------------------------------------------

func Example_telemetryMiddleware() {
	// In your main.go or router setup:
	r := gin.Default()

	// Add this line to enable automatic request tracing:
	r.Use(TelemetryMiddleware())

	// Now all HTTP requests will be automatically traced!
	// The middleware extracts trace context from incoming headers
	// and creates spans for each request.
	fmt.Println("Middleware added to Gin router!")
	// Output: Middleware added to Gin router!
}

// --------------------------------------------------------
// EXAMPLE 3: How to fetch vehicle data with tracing
// --------------------------------------------------------

func Example_fetchAllVehicleData() {
	// This is how you would use it in a Gin handler
	// to fetch data from multiple APIs in parallel

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TelemetryMiddleware())

	// GET /vehicle/:vin - Fetch all data for a vehicle
	r.GET("/vehicle/:vin", func(c *gin.Context) {
		vin := c.Param("vin")

		// Get the trace context from the middleware
		ctx, ok := GetTraceContext(c)
		if !ok {
			ctx = context.Background()
		}

		// Call the function that fetches from all 3 sources in parallel
		// (Sales API, Service API, VIN Database)
		result, err := FetchAllVehicleData(ctx, vin,
			"https://api.sales-partner.com",
			"https://api.service-partner.com",
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"vin":     vin,
			"sales":  result.Sales,
			"service": result.Service,
			"vehicle": result.VIN,
		})
	})

	// Test the endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/vehicle/1HGBH41JXMN109186", nil)
	r.ServeHTTP(w, req)

	fmt.Printf("Response Status: %d\n", w.Code)
	// Output: Response Status: 200
}

// --------------------------------------------------------
// EXAMPLE 4: How to make traced HTTP calls to external APIs
// --------------------------------------------------------

func Example_tracedHTTPClient() {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TelemetryMiddleware())

	// Handler that calls external APIs with trace propagation
	r.GET("/external/:vin", func(c *gin.Context) {
		vin := c.Param("vin")

		ctx, ok := GetTraceContext(c)
		if !ok {
			ctx = context.Background()
		}

		tracer := GetTracer()

		// Create a span for the external call
		_, span := tracer.Start(ctx, "external-api-call",
			trace.WithAttributes(
				attribute.String("external.api", "sales"),
				attribute.String("vin", vin),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Build the request
		url := fmt.Sprintf("https://api.sales-partner.com/vehicles/%s", vin)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Inject trace context into request headers
		// This allows the downstream service to see the trace!
	carrier := httpHeaderCarrier{req.Header}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

		// Make the HTTP call
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

		c.JSON(http.StatusOK, gin.H{"vin": vin, "status": "fetched"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/external/1HGBH41JXMN109186", nil)
	r.ServeHTTP(w, req)

	fmt.Printf("Response Status: %d\n", w.Code)
	// Output: Response Status: 200
}

// ============================================================
// STEP-BY-STEP INTEGRATION GUIDE
// ============================================================

/*
QUICK START GUIDE:

1️⃣  Initialize Telemetry in main.go:
--------------------------------
func main() {
    ctx := context.Background()
    tp, err := telemetry.InitTracer(ctx, "unified-document-viewer")
    if err != nil {
        log.Printf("Warning: %v", err)
    }
    defer tp.Shutdown(ctx)
    
    // ... rest of your code
}

2️⃣  Add Middleware to Gin:
------------------------
    r := gin.Default()
    r.Use(telemetry.TelemetryMiddleware())

3️⃣  Use in Handlers:
-------------------
    // Option A: Use FetchAllVehicleData for parallel API calls
    result, err := telemetry.FetchAllVehicleData(ctx, vin, salesURL, serviceURL)
    
    // Option B: Create custom spans
    ctx, span := tracer.Start(ctx, "my-operation")
    defer span.End()
    
4️⃣  Test locally:
--------------
    // Using the provided docker-compose.yml
    docker-compose up -d
    
    // Start OTel Collector will be at localhost:4318
    
    // Then run your app and check http://localhost:4318 for traces!

5️⃣  Frontend Integration:
------------------------
    When calling your API from frontend, add these headers:
    
    - traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
    
    This allows trace context to flow from frontend → backend!
    
6️⃣  Viewing Traces:
----------------
    // Install Jaeger or Tempo
    docker run -d -p 16686:16686 -p 6831:6831/udp jaegertracing/all-in-one:latest
    
    // Or use OTel Collector with Prometheus
    // Check docker-compose.yml for the full stack!
*/
