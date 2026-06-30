package main

import (
	"fmt"
	"log"
	"net/http"

	"api.poster.com/internal/config"
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
		_ = db
	}

	// 3. Setup Gin Router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 4. Register Route (Ping Endpoint)
	r.GET("/ping", func(c *gin.Context) {
		response.SendSuccess(c, http.StatusOK, "Pong!", gin.H{
			"app": "PlayStation Rental Cashier Platform API",
			"env": cfg.Env,
		})
	})

	// 5. Start Server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server is running on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
