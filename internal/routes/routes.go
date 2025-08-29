package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trinetraguard/backend/internal/handlers"
	"github.com/trinetraguard/backend/internal/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Setup CORS
	setupCORS(router)

	// Setup middleware
	setupMiddleware(router)

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API routes
	setupAPIRoutes(router)

	// 404 handler
	router.NoRoute(notFoundHandler)
}

// setupCORS configures CORS middleware
func setupCORS(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
}

// setupMiddleware configures application middleware
func setupMiddleware(router *gin.Engine) {
	// Request logging
	router.Use(middleware.Logger())

	// Request ID middleware
	router.Use(middleware.RequestID())

	// Error handling
	router.Use(middleware.ErrorHandler())
}

// setupAPIRoutes configures API routes with versioning
func setupAPIRoutes(router *gin.Engine) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		setupUserRoutes(v1)
		setupExampleRoutes(v1)
		setupSystemRoutes(v1)
		setupLostPersonRoutes(v1)
	}
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler()

	users := api.Group("/users")
	{
		users.GET("/", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("/", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}

// setupExampleRoutes configures example routes
func setupExampleRoutes(api *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler()

	examples := api.Group("/examples")
	{
		examples.GET("/", userHandler.GetExamples)
		examples.POST("/", userHandler.CreateExample)
	}
}

// setupSystemRoutes configures system-related routes
func setupSystemRoutes(api *gin.RouterGroup) {
	system := api.Group("/system")
	{
		system.GET("/status", systemStatus)
		system.GET("/info", systemInfo)
	}
}

// setupLostPersonRoutes configures lost person-related routes
func setupLostPersonRoutes(api *gin.RouterGroup) {
	lostPersonHandler := handlers.NewLostPersonHandler()

	lostPersons := api.Group("/lost-persons")
	{
		lostPersons.POST("/", lostPersonHandler.CreateLostPerson)
		lostPersons.GET("/", lostPersonHandler.GetAllLostPersons)
		lostPersons.GET("/search", lostPersonHandler.SearchLostPersons)
		lostPersons.GET("/:id", lostPersonHandler.GetLostPersonByID)
		lostPersons.PUT("/:id", lostPersonHandler.UpdateLostPerson)
		lostPersons.DELETE("/:id", lostPersonHandler.DeleteLostPerson)
	}

	// Image serving route
	api.GET("/images/:path", lostPersonHandler.GetImage)
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"message":   "TrinetraGuard Backend is running",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// systemStatus returns system status information
func systemStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "healthy",
		"uptime":    "running",
		"timestamp": time.Now().UTC(),
		"services": gin.H{
			"api":      "healthy",
			"database": "not_configured",
		},
	})
}

// systemInfo returns system information
func systemInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":        "TrinetraGuard Backend",
		"version":     "1.0.0",
		"environment": "development",
		"go_version":  "1.21",
		"framework":   "gin",
		"timestamp":   time.Now().UTC(),
	})
}

// notFoundHandler handles 404 requests
func notFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"error":   "Not Found",
		"message": "The requested resource was not found",
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})
}
