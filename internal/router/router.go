package router

import (
	"fmt"
	"net/http"
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
// Update the NewRouter function to include brand routes
func NewRouter(
	respHelper *responses.ResponseHelper,
	healthHandler *handlers.HealthHandler,
	brandHandler *handlers.BrandHandler,
) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(middlewares.CaseConverterMiddleware())
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

	// Handle 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		respHelper.SendError(c, "Route not found", "The requested endpoint does not exist", http.StatusNotFound)
	})

	// Register routes
	router.GET("/ping", healthHandler.HealthCheck)

	// Brand routes
	api := router.Group("/api")
	{
		brands := api.Group("/brands")
		{
			brands.GET("", brandHandler.GetAllBrands)
			brands.GET("/:id", brandHandler.GetBrand)
			brands.POST("", brandHandler.CreateBrand)
			brands.PUT("/:id", brandHandler.UpdateBrand)
			brands.DELETE("/:id", brandHandler.DeleteBrand)
			brands.GET("/active", brandHandler.GetActiveBrands)
			brands.GET("/grouped", brandHandler.GetGroupedBrands)
		}
	}

	return router
}
