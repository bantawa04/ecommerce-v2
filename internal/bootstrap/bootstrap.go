package bootstrap

import (
	"context"
	"log"
	"net/http"

	"beautyessentials.com/internal/api/handlers"
	"beautyessentials.com/internal/api/responses"
	"beautyessentials.com/internal/config"
	repoImpl "beautyessentials.com/internal/repository/implementations"
	"beautyessentials.com/internal/router"
	serviceImpl "beautyessentials.com/internal/service/implementations"
	"beautyessentials.com/internal/service/external" // Add this import
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module exported for initializing application
var Module = fx.Options(
	ConfigModule,
	RepositoryModule,
	ServiceModule,
	HandlerModule,
	RouterModule,
	fx.Invoke(bootstrap),
)

// ConfigModule provides configuration dependencies
var ConfigModule = fx.Options( 
	fx.Provide(config.LoadConfig),
	fx.Provide(config.InitDatabase),
	fx.Provide(responses.NewResponseHelper),
)

// RepositoryModule provides repository dependencies
var RepositoryModule = fx.Options(
	fx.Provide(repoImpl.NewHealthRepository),
	fx.Provide(repoImpl.NewBrandRepository),
	fx.Provide(repoImpl.NewCategoryRepository),
	fx.Provide(repoImpl.NewMediaRepository), // Add media repository
)

// ServiceModule provides service dependencies
var ServiceModule = fx.Options(
	fx.Provide(serviceImpl.NewHealthService),
	fx.Provide(serviceImpl.NewBrandService),
	fx.Provide(serviceImpl.NewCategoryService),
	fx.Provide(serviceImpl.NewMediaService), // Add media service
	fx.Provide(external.NewImageKitService), // Add ImageKit service for media uploads
)

// HandlerModule provides handler dependencies
var HandlerModule = fx.Options(
	fx.Provide(handlers.NewHealthHandler),
	fx.Provide(handlers.NewBrandHandler),
	fx.Provide(handlers.NewCategoryHandler),
	fx.Provide(handlers.NewMediaHandler), // Add media handler
)

// RouterModule provides router dependencies
var RouterModule = fx.Options(
	fx.Provide(router.NewRouter),
	fx.Provide(newHTTPServer),
)

// BuildApp constructs the fx application with all dependencies
func BuildApp() *fx.App {
	return fx.New(Module)
}

// newHTTPServer creates an HTTP server with the provided router and configuration
func newHTTPServer(router *gin.Engine, cfg *config.Config) *http.Server {
	serverConfig := cfg.Server()
	return &http.Server{
		Addr:         ":" + serverConfig.Port,
		Handler:      router,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
	}
}

// bootstrap registers lifecycle hooks and initializes the application
func bootstrap(
	lifecycle fx.Lifecycle,
	router *gin.Engine,
	server *http.Server,
	cfg *config.Config,
	db *gorm.DB,
) {
	// Register server start and stop hooks
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting Application")
			log.Println("------------------------")
			log.Println("-- Beauty Essentials API --")
			log.Println("------------------------")

			// Start the server in a goroutine
			go func() {
				serverConfig := cfg.Server()
				log.Printf("Server starting on port %s", serverConfig.Port)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			
			// Close database connection
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				_ = sqlDB.Close()
			}
			
			return server.Shutdown(ctx)
		},
	})
}
