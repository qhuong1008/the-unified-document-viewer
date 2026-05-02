package main

import (
	"log"
	"time"

	"the-unified-document-viewer/internal/api"
	"the-unified-document-viewer/internal/auth"
	"the-unified-document-viewer/internal/database"
	"the-unified-document-viewer/internal/handlers"
	"the-unified-document-viewer/internal/middleware"
	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"
	"the-unified-document-viewer/internal/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=the_unified_document_viewer port=5432 sslmode=disable"
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v", err)
	}

	vaultRepo := repository.NewPostgresRepository(db)

	userRepo := repository.NewUserRepository(db)
	if err := userRepo.AutoMigrate(); err != nil {
		log.Printf("Warning: Failed to auto-migrate users table: %v", err)
	}

	defaultUsername := "admin"
	_, err = userRepo.FindByUsername(defaultUsername)
	if err != nil {
		hashedPassword, _ := auth.HashPassword("admin123")
		defaultUser := &models.User{
			Username: defaultUsername,
			Password: hashedPassword,
		}
		if err := userRepo.Create(defaultUser); err != nil {
			log.Printf("Warning: Failed to create default user: %v", err)
		} else {
			log.Println("Created default admin user (username: admin, password: admin123)")
		}
	}

	jwtManager := auth.NewJWTManager("your-secret-key-here", 24*time.Hour)

	jobQueue := make(chan worker.Job, 100)

	worker.StartWorkerPool(jobQueue, 5, vaultRepo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	webhookHandler := &api.WebhookHandler{JobQueue: jobQueue}
	authHandler := api.NewAuthHandler(jwtManager, userRepo)
	vaultHandler := &handlers.VehicleDigitalVaultHandler{Repo: vaultRepo}

	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.RefreshToken)

	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleware(jwtManager))
	{
		protected.POST("/webhooks/sales", webhookHandler.HandleSalesWebhook)
		protected.POST("/webhooks/service", webhookHandler.HandleServiceWebhook)
		protected.GET("/vault/:vin", vaultHandler.GetVehicleHistory)
	}

	// Public proxy endpoint to bypass CORS when accessing external files
	r.GET("/file-proxy", api.FileProxyHandler())

	log.Println("Server đang chạy tại port :8080...")
	r.Run(":8080")
}
