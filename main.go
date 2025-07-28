package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"video-analysis-service/internal/config"
	"video-analysis-service/internal/database"
	"video-analysis-service/internal/handlers"
	"video-analysis-service/internal/middleware"
	"video-analysis-service/internal/services"
)

// @title Video Analysis Service API
// @version 1.0
// @description A comprehensive video analysis service with face detection and person tracking
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
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	logger.Info("Starting Video Analysis Service...")

	// Initialize database
	var db *database.Database
	var err error

	if cfg.Database.Driver == "sqlite3" {
		db, err = database.Initialize("sqlite", cfg.Database.Name)
	} else {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)
		db, err = database.Initialize("postgres", dsn)
	}

	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize services
	videoService := services.NewVideoService(db.DB, cfg)
	analysisService := services.NewAnalysisService(db.DB, cfg)
	finderService := services.NewFinderService(db.DB, cfg)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoService, logger)
	analysisHandler := handlers.NewAnalysisHandler(analysisService, logger)
	finderHandler := handlers.NewFinderHandler(finderService, logger)

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.ErrorHandler(logger))

	// API routes
	api := router.Group("/api/v1")
	{
		// Video management
		videos := api.Group("/videos")
		{
			videos.POST("/upload", videoHandler.UploadVideo)
			videos.GET("", videoHandler.ListVideos)
			videos.GET("/:id", videoHandler.GetVideo)
			videos.DELETE("/:id", videoHandler.DeleteVideo)
			videos.GET("/:id/download", videoHandler.DownloadVideo)
		}

		// Analysis
		analysis := api.Group("/analysis")
		{
			analysis.POST("/:videoId/start", analysisHandler.StartAnalysis)
			analysis.GET("/:videoId/status", analysisHandler.GetAnalysisStatus)
			analysis.GET("/:videoId/results", analysisHandler.GetAnalysisResults)
			analysis.GET("/:videoId/enhanced-results", analysisHandler.GetEnhancedAnalysisResults)
			analysis.POST("/batch", analysisHandler.StartBatchAnalysis)
		}

		// Person finder
		finder := api.Group("/finder")
		{
			finder.POST("/images/upload", finderHandler.UploadReferenceImage)
			finder.GET("/images", finderHandler.ListReferenceImages)
			finder.DELETE("/images/:id", finderHandler.DeleteReferenceImage)
			finder.POST("/search", finderHandler.SearchPerson)
			finder.GET("/search/:jobId/status", finderHandler.GetSearchStatus)
			finder.GET("/search/:jobId/results", finderHandler.GetSearchResults)
		}

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"service":   "video-analysis-service",
			})
		})
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
