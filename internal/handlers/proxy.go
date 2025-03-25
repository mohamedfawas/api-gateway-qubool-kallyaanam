package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/internal/service"
	"github.com/mohamedfawas/api-gateway-qubool-kallyaanam/pkg/logging"
	"go.uber.org/zap"
)

// ProxyHandler handles forwarding requests to microservices.
type ProxyHandler struct {
	proxyService service.ProxyService
	logger       *zap.Logger
}

// NewProxyHandler creates a new proxy handler with service URLs.
func NewProxyHandler(services map[string]string) (*ProxyHandler, error) {
	proxyService, err := service.NewProxyService(services) // creates a new proxy service instance with the provided services map
	if err != nil {
		return nil, err
	}

	return &ProxyHandler{
		proxyService: proxyService,
		logger:       logging.Logger(),
	}, nil
}

// ProxyRequest forwards a request to the specified service.
func (h *ProxyHandler) ProxyRequest(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceURL, exists := h.proxyService.GetServiceURL(serviceName)
		if !exists {
			// Only add the error to the context, don't respond immediately
			c.Error(fmt.Errorf("service configuration not found: %s", serviceName))
			// Let the SimpleErrorHandler middleware handle the response
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(serviceURL) //  forwards requests to the specified service URL

		originalDirector := proxy.Director // 'Director' function that modifies the outbound request before it’s sent
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Scheme = serviceURL.Scheme                  // sets the protocol (either http or https) to match that of the target service , If serviceURL is https://api.example.com/v1, the scheme becomes https
			req.URL.Host = serviceURL.Host                      // sets the host to the host part of the target service URL, If serviceURL is https://api.example.com/v1, the host becomes api.example.com
			req.URL.Path = serviceURL.Path + c.Request.URL.Path // appends the original request path to the target service path, If the original request path is /users/123, the final path becomes /api/v1/users/123
			if clientIP := c.ClientIP(); clientIP != "" {       // checks if the client’s IP address is available. If it is, the IP is added to the request headers under X-Forwarded-For
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			// Log the error
			h.logger.Error("proxy error",
				zap.String("service", serviceName),
				zap.Error(err),
			)

			// Only add the error to the context, don't respond immediately
			c.Error(fmt.Errorf("proxy error to %s: %v", serviceName, err))
			// Setting status code to make sure SimpleErrorHandler knows it's a bad gateway
			c.Status(http.StatusBadGateway)
			// Let the SimpleErrorHandler middleware handle the response
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
