package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/database"
	"video-analysis-service/internal/handlers"
	"video-analysis-service/internal/middleware"
	"video-analysis-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Video Analysis Service API
// @version 1.0
// @description A comprehensive video analysis service for person detection, tracking, and face matching
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Initialize services
	videoService := services.NewVideoService(db.DB, cfg)
	analysisService := services.NewAnalysisService(db.DB, cfg)
	finderService := services.NewFinderService(db.DB, cfg)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoService, logger)
	analysisHandler := handlers.NewAnalysisHandler(analysisService, logger)
	finderHandler := handlers.NewFinderHandler(finderService, logger)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.Timeout(30 * time.Second))

	// API routes
	api := router.Group("/api/v1")
	{
		// Video management endpoints
		videos := api.Group("/videos")
		{
			videos.POST("/upload", videoHandler.UploadVideo)
			videos.GET("", videoHandler.ListVideos)
			videos.GET("/:id", videoHandler.GetVideo)
			videos.DELETE("/:id", videoHandler.DeleteVideo)
			videos.GET("/:id/download", videoHandler.DownloadVideo)
		}

		// Analysis endpoints
		analysis := api.Group("/analysis")
		{
			analysis.POST("/:videoId/start", analysisHandler.StartAnalysis)
			analysis.GET("/:videoId/status", analysisHandler.GetAnalysisStatus)
			analysis.GET("/:videoId/results", analysisHandler.GetAnalysisResults)
			analysis.POST("/batch", analysisHandler.StartBatchAnalysis)
		}

		// Finder endpoints
		finder := api.Group("/finder")
		{
			finder.POST("/upload", finderHandler.UploadReferenceImage)
			finder.GET("/images", finderHandler.ListReferenceImages)
			finder.POST("/search", finderHandler.SearchPerson)
			finder.GET("/search/:searchId/status", finderHandler.GetSearchStatus)
			finder.GET("/search/:searchId/results", finderHandler.GetSearchResults)
			finder.DELETE("/images/:id", finderHandler.DeleteReferenceImage)
		}

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "Video Analysis Service",
				"version": "1.0.0",
				"time":    time.Now().UTC(),
			})
		})
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
