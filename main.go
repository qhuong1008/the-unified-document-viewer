package main

import (
	"log"
	"the-unified-document-viewer/internal/api"
	"the-unified-document-viewer/internal/database"
	"the-unified-document-viewer/internal/handlers"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Khởi tạo kết nối PostgreSQL
	dsn := "host=localhost user=postgres password=postgres dbname=the_unified_document_viewer port=5432 sslmode=disable"
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v", err)
	}

	// 2. Khởi tạo Repository Layer
	// Đây là nơi chứa logic Upsert để chống trùng lặp dữ liệu (Idempotency)
	repository := repository.NewPostgresRepository(db)

	// 3. Thiết lập hàng đợi công việc (Job Queue)
	// Buffer 100 giúp hệ thống chịu được các đợt bùng nổ request (burst traffic)
	jobQueue := make(chan worker.Job, 100)

	// 4. Khởi chạy Worker Pool
	// Quan trọng: Truyền repo vào để Worker có thể lưu dữ liệu sau khi Transform & Enrich
	worker.StartWorkerPool(jobQueue, 5, repository)

	// 5. Khởi tạo REST API với Gin
	r := gin.Default()
	
	// Inject jobQueue vào handler để đẩy job từ Webhook sang Worker
	handler := &api.WebhookHandler{JobQueue: jobQueue}

	// Route tiếp nhận dữ liệu từ các hệ thống nguồn (Scenario D)
	r.POST("/webhooks/sales", handler.HandleSalesWebhook)
	r.POST("/webhooks/service", handler.HandleServiceWebhook)

// Route lấy dữ liệu đã hợp nhất (Sẽ triển khai ở bước sau)
	// Thay thế dòng comment bằng:
vaultHandler := &handlers.VehicleDigitalVaultHandler{Repo: repository}
	r.GET("/vault/:vin", vaultHandler.GetVehicleHistory)
	log.Println("Server đang chạy tại port :8080...")
	r.Run(":8080")
}