// api-gateway-qubool-kallyaanam/cmd/main.go
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
)

func main() {
	router := gin.Default()

	// Health Check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "api-gateway",
			"version": "0.1.0",
		})
	})

	// Route definitions for other services
	router.GET("/auth/health", func(c *gin.Context) {
		authServiceURL := os.Getenv("AUTH_SERVICE_URL")
		if authServiceURL == "" {
			authServiceURL = "http://auth-service:8081" // Default in Docker
		}

		resp, err := http.Get(authServiceURL + "/health")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "auth-service",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "auth-service",
				"code":    resp.StatusCode,
			})
			return
		}

		c.Status(resp.StatusCode)
		c.Writer.Write([]byte("Auth Service is UP"))
	})

	router.GET("/user/health", func(c *gin.Context) {
		userServiceURL := os.Getenv("USER_SERVICE_URL")
		if userServiceURL == "" {
			userServiceURL = "http://user-service:8082" // Default in Docker
		}

		resp, err := http.Get(userServiceURL + "/health")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "user-service",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "user-service",
				"code":    resp.StatusCode,
			})
			return
		}

		c.Status(resp.StatusCode)
		c.Writer.Write([]byte("User Service is UP"))
	})

	router.GET("/admin/health", func(c *gin.Context) {
		adminServiceURL := os.Getenv("ADMIN_SERVICE_URL")
		if adminServiceURL == "" {
			adminServiceURL = "http://admin-service:8083" // Default in Docker
		}

		resp, err := http.Get(adminServiceURL + "/health")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "admin-service",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "DOWN",
				"service": "admin-service",
				"code":    resp.StatusCode,
			})
			return
		}

		c.Status(resp.StatusCode)
		c.Writer.Write([]byte("Admin Service is UP"))
	})

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a timeout of 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
