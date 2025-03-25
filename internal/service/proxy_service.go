package service

import (
	"net/url"
)

// ProxyService defines the interface for proxy operations
type ProxyService interface {
	GetServiceURL(serviceName string) (*url.URL, bool)
}

// proxyService implements ProxyService
type proxyService struct {
	serviceURLs map[string]*url.URL
}

// NewProxyService creates a new proxy service
func NewProxyService(services map[string]string) (ProxyService, error) {
	serviceURLs := make(map[string]*url.URL)
	for name, urlStr := range services {
		serviceURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		serviceURLs[name] = serviceURL
	}

	return &proxyService{
		serviceURLs: serviceURLs,
	}, nil
}

// GetServiceURL returns the URL for a given service
func (s *proxyService) GetServiceURL(serviceName string) (*url.URL, bool) {
	url, exists := s.serviceURLs[serviceName]
	return url, exists
}
