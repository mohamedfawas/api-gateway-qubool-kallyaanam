package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ProxyHandler handles forwarding requests to microservices.
type ProxyHandler struct {
	serviceURLs map[string]*url.URL
}

// NewProxyHandler creates a new proxy handler with service URLs.
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

// ProxyRequest forwards a request to the specified service.
func (h *ProxyHandler) ProxyRequest(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceURL, exists := h.serviceURLs[serviceName]
		if !exists {
			// Only add the error to the context, don't respond immediately
			c.Error(fmt.Errorf("service configuration not found: %s", serviceName))
			// Let the SimpleErrorHandler middleware handle the response
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(serviceURL)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Scheme = serviceURL.Scheme
			req.URL.Host = serviceURL.Host
			req.URL.Path = serviceURL.Path + c.Request.URL.Path
			if clientIP := c.ClientIP(); clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			// Only add the error to the context, don't respond immediately
			c.Error(fmt.Errorf("proxy error to %s: %v", serviceName, err))
			// Setting status code to make sure SimpleErrorHandler knows it's a bad gateway
			c.Status(http.StatusBadGateway)
			// Let the SimpleErrorHandler middleware handle the response
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
