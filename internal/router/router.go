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
// Update the NewRouter function to include brand and category routes
func NewRouter(
	respHelper *responses.ResponseHelper,
	healthHandler *handlers.HealthHandler,
	brandHandler *handlers.BrandHandler,
	categoryHandler *handlers.CategoryHandler,
	mediaHandler *handlers.MediaHandler,
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
		// Add error to context instead of handling directly
		appErr := middlewares.NewNotFoundError("Route not found", "The requested endpoint does not exist")
		_ = c.Error(appErr)
	})

	// Register routes
	router.GET("/ping", healthHandler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Brand routes
		brands := api.Group("/brands")
		{
			brands.GET("", brandHandler.GetAllBrands)
			brands.GET("/:id", brandHandler.GetBrand)
			brands.POST("", brandHandler.CreateBrand)
			brands.PUT("/:id", brandHandler.UpdateBrand)
			brands.DELETE("/:id", brandHandler.DeleteBrand)			
			brands.GET("/grouped", brandHandler.GetGroupedBrands)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.POST("", categoryHandler.CreateCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.GET("/active", categoryHandler.GetActiveCategories)
			categories.GET("/slug/:slug", categoryHandler.FindCategoryBySlug)
		}
		
		// Media routes
		media := api.Group("/media")
		{
			media.GET("", mediaHandler.GetAllMedia)         // index
			media.POST("", mediaHandler.CreateMedia)        // store
			media.DELETE("/:id", mediaHandler.DeleteMedia)  // destroy
			// Remove other routes that don't match Laravel's apiResource except update
		}
	}

	return router
}
