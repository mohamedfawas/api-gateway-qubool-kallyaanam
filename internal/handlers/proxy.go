package handlers

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

// ProxyHandler handles forwarding requests to the appropriate microservice
type ProxyHandler struct {
	serviceURLs map[string]*url.URL
}

// NewProxyHandler creates a new proxy handler with service URL mappings
func NewProxyHandler(cfg *config.Config) (*ProxyHandler, error) {
	urls := make(map[string]*url.URL)

	// Parse service URLs from config
	for service, endpoint := range cfg.Services {
		serviceURL, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		urls[service] = serviceURL
	}

	return &ProxyHandler{
		serviceURLs: urls,
	}, nil
}

// ProxyToService returns a handler that forwards requests to the specified service
func (h *ProxyHandler) ProxyToService(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceURL, exists := h.serviceURLs[serviceName]
		if !exists {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
				http.StatusInternalServerError,
				"Service configuration not found",
				nil,
			))
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(serviceURL)

		// Customize director to maintain path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Extract path parameter (removing the first /* match)
			path := c.Param("path")

			req.URL.Scheme = serviceURL.Scheme
			req.URL.Host = serviceURL.Host
			req.URL.Path = serviceURL.Path + path

			// Forward user context if available
			if userID, exists := c.Get("userID"); exists {
				req.Header.Set("X-User-ID", userID.(string))
			}

			if userRole, exists := c.Get("userRole"); exists {
				req.Header.Set("X-User-Role", userRole.(string))
			}

			// Forward request ID
			if requestID, exists := c.Get("requestID"); exists {
				req.Header.Set("X-Request-ID", requestID.(string))
			}

			// Forward the original caller's IP
			if clientIP, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		// Handle errors from downstream services
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			c.JSON(http.StatusBadGateway, models.NewErrorResponse(
				http.StatusBadGateway,
				"Service unavailable",
				err.Error(),
			))
		}

		// Serve the request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
