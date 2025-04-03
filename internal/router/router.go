package router

import (
	"fmt"
	"time"

	"beautyessentials.com/internal/api/handlers"
	"beautyessentials.com/internal/api/middlewares"
	"beautyessentials.com/internal/api/responses"
	"github.com/gin-gonic/gin"
)

// Router interface defines methods for the router
type Router interface {
	Run(addr ...string) error
}

// NewRouter creates and configures a Gin router
func NewRouter(
	respHelper *responses.ResponseHelper,
	healthHandler *handlers.HealthHandler,
	brandHandler *handlers.BrandHandler,
) *gin.Engine {
	router := gin.Default()

	// Add case converter middleware
	router.Use(middlewares.CaseConverterMiddleware())
	
	// Add error handler middleware
	router.Use(middlewares.ErrorHandler(respHelper))

	// Add custom logger middleware
	router.Use(func(c *gin.Context) {
		// Start timer
		t := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(t)
		status := c.Writer.Status()
		fmt.Printf("[GIN] %s | %3d | %s | %s\n",
			c.Request.Method,
			status,
			latency,
			c.Request.URL.Path,
		)
	})

	// Register routes
	router.GET("/ping", healthHandler.HealthCheck)

	// Brand routes
	router.GET("/api/brands", brandHandler.GetAllBrands)

	return router
}
