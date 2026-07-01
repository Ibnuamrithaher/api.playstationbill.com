package main

import (
	"fmt"
	"log"
	"net/http"

	"api.poster.com/internal/config"
	deliveryHTTP "api.poster.com/internal/delivery/http"
	"api.poster.com/internal/domain"
	"api.poster.com/internal/repository"
	"api.poster.com/internal/service"
	"api.poster.com/pkg/database"
	"api.poster.com/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v. Continuing without DB connection.", err)
	} else {
		// Run Auto Migration
		if err := db.AutoMigrate(&domain.Category{}); err != nil {
			log.Printf("Warning: Auto-migration failed: %v", err)
		}
	}

	// 3. Setup Gin Router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Serve uploaded files statically
	r.Static("/public", "./public")

	// 4. Register Route (Ping Endpoint)
	r.GET("/ping", func(c *gin.Context) {
		response.SendSuccess(c, http.StatusOK, "Pong!", gin.H{
			"app": "PlayStation Rental Cashier Platform API",
			"env": cfg.Env,
		})
	})

	// Register Category Routes if DB is available
	if db != nil {
		categoryRepo := repository.NewCategoryRepository(db)
		categoryService := service.NewCategoryService(categoryRepo)
		categoryHandler := deliveryHTTP.NewCategoryHandler(categoryService)

		r.POST("/api/category", categoryHandler.Create)
	} else {
		r.POST("/api/category", func(c *gin.Context) {
			response.SendError(c, http.StatusInternalServerError, "Database connection not available", "DB is nil")
		})
	}

	// 5. Start Server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server is running on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

