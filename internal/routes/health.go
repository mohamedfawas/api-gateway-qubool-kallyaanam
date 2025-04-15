// Package routes provides routing definitions for the API Gateway
package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/utils"
)

// Constants for health status
const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

// registerHealthRoutes registers health check endpoints
func registerHealthRoutes(router *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	router.GET("/health", aggregatedHealthCheck(cfg, logger))
}

// aggregatedHealthCheck provides a comprehensive health check endpoint for the API Gateway
// and all its dependent services
func aggregatedHealthCheck(cfg *config.Config, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create the aggregated health response
		health := models.AggregatedHealth{
			Gateway: models.ServiceHealth{
				Name:      "api-gateway",
				Status:    StatusUp,
				Version:   "0.1.0", // This should come from a version package or build info
				Timestamp: time.Now(),
			},
			Services: make([]models.ServiceHealth, 0, 3), // Pre-allocate for our services
		}

		// Check all services health in parallel
		var wg sync.WaitGroup
		var mutex sync.Mutex // To protect concurrent writes to the services slice

		// Helper function to fetch health status for a service
		checkServiceHealth := func(name, url string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()

			serviceHealth := models.ServiceHealth{
				Name:      name,
				Status:    StatusDown, // Default to down
				Timestamp: time.Now(),
			}

			// Create request
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/health", nil)
			if err != nil {
				logger.Error("Failed to create request for service health check",
					zap.String("service", name),
					zap.String("url", url),
					zap.Error(err),
				)
				mutex.Lock()
				health.Services = append(health.Services, serviceHealth)
				mutex.Unlock()
				return
			}

			// Send request
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				logger.Error("Service health check failed",
					zap.String("service", name),
					zap.String("url", url),
					zap.Error(err),
				)
				mutex.Lock()
				health.Services = append(health.Services, serviceHealth)
				mutex.Unlock()
				return
			}
			defer resp.Body.Close()

			// Parse response
			if resp.StatusCode == http.StatusOK {
				// Try to decode the response body
				var healthResp map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&healthResp); err == nil {
					// Extract status and details
					if status, ok := healthResp["status"].(string); ok {
						serviceHealth.Status = status
					} else {
						serviceHealth.Status = StatusUp // Default to UP if status is present but not a string
					}

					if version, ok := healthResp["version"].(string); ok {
						serviceHealth.Version = version
					}

					// Extract any other details as strings
					details := make(map[string]string)
					for k, v := range healthResp {
						if k != "status" && k != "version" {
							// Convert any value to string
							details[k] = stringify(v)
						}
					}

					if len(details) > 0 {
						serviceHealth.Details = details
					}
				} else {
					// If we can't parse the response, still mark as UP but log the error
					serviceHealth.Status = StatusUp
					logger.Warn("Couldn't parse service health response",
						zap.String("service", name),
						zap.Error(err),
					)
				}
			} else {
				logger.Warn("Service health check returned non-200 status",
					zap.String("service", name),
					zap.Int("statusCode", resp.StatusCode),
				)
			}

			mutex.Lock()
			health.Services = append(health.Services, serviceHealth)
			mutex.Unlock()
		}

		// Start health checks for each service
		wg.Add(3)
		go checkServiceHealth("auth-service", cfg.Services.AuthServiceURL)
		go checkServiceHealth("user-service", cfg.Services.UserServiceURL)
		go checkServiceHealth("admin-service", cfg.Services.AdminServiceURL)

		// Wait for all checks to complete
		wg.Wait()

		// Check if any service is down
		allServicesUp := true
		for _, service := range health.Services {
			if service.Status != StatusUp {
				allServicesUp = false
				break
			}
		}

		if !allServicesUp {
			// If any service is down, the gateway is still up but with degraded functionality
			health.Gateway.Details = map[string]string{
				"message": "Some services are unavailable",
			}
		}

		// Use the new response utility instead of direct JSON formatting
		utils.RespondWithSuccess(c, "Health check successful", health)
	}
}

// stringify converts an interface value to a string
func stringify(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	default:
		// Try to marshal to JSON
		if data, err := json.Marshal(v); err == nil {
			return string(data)
		}
		return "unknown"
	}
}
