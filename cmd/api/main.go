package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ecosystem.garyle/service/internal/app/config"
	"ecosystem.garyle/service/internal/app/module"
	"ecosystem.garyle/service/internal/app/module/ota"
	productModule "ecosystem.garyle/service/internal/app/module/wms/master-data/product"
	"ecosystem.garyle/service/internal/infrastructure/database"
	"ecosystem.garyle/service/internal/infrastructure/middleware"
	"ecosystem.garyle/service/pkg/logger"
	"ecosystem.garyle/service/pkg/utils/response"
)

func main() {
	fx.New(
		// Provide core dependencies
		fx.Provide(
			config.NewConfig,
			func() logger.Logger {
				return logger.New(logger.Info)
			},
			newGinEngine,
		),

		// Include modules
		database.Module,
		module.Module,

		// Register lifecycle hooks
		fx.Invoke(
			database.RegisterHooks,
			registerRoutes,
			startServer,
		),
	).Run()
}

// newGinEngine creates and configures a new Gin engine
func newGinEngine(cfg *config.Config, log logger.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(middleware.Logger(log))
	router.Use(gin.Recovery())
	return router
}

// registerRoutes sets up all API routes
func registerRoutes(router *gin.Engine, db *sql.DB) {
	// API v1 routes
	apiV1 := router.Group("/api/v1")

	// Health check endpoint
	apiV1.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"service": "Garyle Ecosystem Service",
			"status":  "healthy",
		}, "Service is healthy")
	})

	// Register feature routes
	ota.RegisterOTAHandlers(db, apiV1)
	productModule.RegisterProductHandlers(db, apiV1)
}

// startServer starts the HTTP server
func startServer(lc fx.Lifecycle, cfg *config.Config, router *gin.Engine, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
			log.Infof("HTTP server listening on %s", serverAddr)

			// Start server in a goroutine so it doesn't block fx lifecycle
			go func() {
				if err := router.Run(serverAddr); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()

			return nil
		},
	})
}
