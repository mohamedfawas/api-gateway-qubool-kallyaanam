package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/models"
)

// ProxyHandler handles forwarding requests to microservices
type ProxyHandler struct {
	serviceURLs map[string]*url.URL
}

// NewProxyHandler creates a new proxy handler with service URLs
func NewProxyHandler(services map[string]string) (*ProxyHandler, error) {
	serviceURLs := make(map[string]*url.URL)

	for name, urlStr := range services {
		serviceURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		serviceURLs[name] = serviceURL
	}

	return &ProxyHandler{
		serviceURLs: serviceURLs,
	}, nil
}

// ProxyRequest forwards a request to the specified service
func (h *ProxyHandler) ProxyRequest(serviceName string) gin.HandlerFunc {
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

		// Create a reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(serviceURL)

		// Customize director to modify requests as needed
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			req.URL.Scheme = serviceURL.Scheme
			req.URL.Host = serviceURL.Host

			// Preserve original path
			req.URL.Path = serviceURL.Path + c.Request.URL.Path

			// Forward client IP for tracking purposes
			if clientIP := c.ClientIP(); clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		// Handle errors when the service is unavailable
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			c.JSON(http.StatusBadGateway, models.NewErrorResponse(
				http.StatusBadGateway,
				"Service unavailable",
				err.Error(),
			))
		}

		// Serve the proxy request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
