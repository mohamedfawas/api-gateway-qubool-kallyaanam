// api-gateway-qubool-kallyaanam/cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/routes"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/utils"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize logger
	logger, err := utils.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Ensure logs are flushed

	// Set Gin mode
	if !cfg.Logging.Development {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Register middlewares
	middleware.RegisterMiddlewares(router, cfg, logger)

	// Register routes
	routes.RegisterRoutes(router, cfg, logger)

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting API Gateway server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
}
