// Package models provides data models for the API Gateway
package models

import "time"

// ServiceHealth represents the health status of a microservice
type ServiceHealth struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Version   string            `json:"version,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// AggregatedHealth represents the overall health status of the API Gateway
// and its dependent services
type AggregatedHealth struct {
	Gateway  ServiceHealth   `json:"gateway"`
	Services []ServiceHealth `json:"services"`
}
