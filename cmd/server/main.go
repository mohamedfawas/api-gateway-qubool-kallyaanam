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
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/handlers"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Validate essential config
	if cfg.Auth.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Create the router
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	// Create proxy handler for services
	proxyHandler, err := handlers.NewProxyHandler(cfg.Services)
	if err != nil {
		log.Fatalf("Failed to initialize proxy handler: %v", err)
	}

	// Health check endpoint for Kubernetes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth service routes - no auth required
	authRoutes := router.Group("/auth")
	authRoutes.Any("/*path", proxyHandler.ProxyRequest("auth"))

	// Protected routes - require auth
	protectedRoutes := router.Group("")
	protectedRoutes.Use(middleware.JWTAuth(cfg))

	// User service routes
	protectedRoutes.Any("/user/*path", proxyHandler.ProxyRequest("user"))

	// Admin service routes
	protectedRoutes.Any("/admin/*path", proxyHandler.ProxyRequest("admin"))

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API Gateway on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
