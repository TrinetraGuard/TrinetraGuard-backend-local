package main

import (
	"log"
	"os"

	"video-processing-backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Create upload directories if they don't exist
	os.MkdirAll("videos", 0755)
	os.MkdirAll("faces", 0755)
	os.MkdirAll("storage", 0755)

	// Initialize video storage
	handlers.InitializeStorage()

	// Setup routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Serve the frontend
	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	// Serve the storage page
	r.GET("/storage", func(c *gin.Context) {
		c.File("storage.html")
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/upload-video", handlers.UploadVideoHandler)
		api.GET("/health", handlers.HealthCheckHandler)
		
		// Storage management routes
		api.GET("/videos", handlers.ListVideosHandler)
		api.GET("/videos/:id", handlers.GetVideoHandler)
		api.DELETE("/videos/:id", handlers.DeleteVideoHandler)
		api.GET("/videos/stats", handlers.GetVideoStatsHandler)
		api.POST("/videos/cleanup", handlers.CleanupOldVideosHandler)
	}

	// Serve static files
	r.Static("/faces", "./faces")
}
