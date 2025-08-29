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

	// Configure CORS for API usage
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	r.Use(cors.New(config))

	// Create upload directories if they don't exist
	os.MkdirAll("../storage/videos", 0755)
	os.MkdirAll("../storage/faces", 0755)
	os.MkdirAll("../storage/data", 0755)
	os.MkdirAll("../storage/temp", 0755)

	// Initialize video storage
	handlers.InitializeStorage()

	// Setup API routes
	setupAPIRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend API server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupAPIRoutes(r *gin.Engine) {
	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", handlers.HealthCheckHandler)

		// Video upload and processing
		api.POST("/upload-video", handlers.UploadVideoHandler)
		api.POST("/search-by-face", handlers.SearchByFaceHandler)

		// Storage management routes
		api.GET("/videos", handlers.ListVideosHandler)
		api.GET("/videos/active", handlers.ListActiveVideosHandler)
		api.GET("/videos/archived", handlers.ListArchivedVideosHandler)
		api.GET("/videos/search", handlers.SearchVideosHandler)
		api.GET("/videos/:id", handlers.GetVideoHandler)
		api.DELETE("/videos/:id", handlers.DeleteVideoHandler)
		api.POST("/videos/:id/restore", handlers.RestoreVideoHandler)
		api.GET("/videos/stats", handlers.GetVideoStatsHandler)
		api.POST("/videos/cleanup", handlers.CleanupOldVideosHandler)
		api.POST("/videos/reset-database", handlers.ResetDatabaseHandler)

		// Search history endpoints
		api.GET("/search-history", handlers.GetSearchHistoryHandler)
		api.GET("/search-history/stats", handlers.GetSearchHistoryStatsHandler)

		// Video preview and file serving
		api.GET("/videos/:id/preview", handlers.GetVideoPreviewHandler)
		api.GET("/videos/:id/file", handlers.GetVideoFileHandler)

		// Face images serving
		api.Static("/faces", "../storage/faces")
	}

	// Root endpoint for API info
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "TrinetraGuard Backend API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "/api/health",
				"upload": "/api/upload-video",
				"search": "/api/search-by-face",
				"videos": "/api/videos",
			},
		})
	})
}
