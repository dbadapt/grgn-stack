// Package grgn provides core interfaces for the GRGN stack.
// These interfaces are standalone and can be imported by external Go projects.
package grgn

import (
	"context"
)

// IDatabase defines the contract for database operations.
// Product domains depend on this interface, not implementations.
type IDatabase interface {
	// Ping performs a connectivity check
	Ping(ctx context.Context) error

	// Close gracefully closes the database connection
	Close(ctx context.Context) error

	// ExecuteRead executes a read transaction
	ExecuteRead(ctx context.Context, work func(ctx context.Context) (any, error)) (any, error)

	// ExecuteWrite executes a write transaction
	ExecuteWrite(ctx context.Context, work func(ctx context.Context) (any, error)) (any, error)
}

// HealthStatus represents the health status of a service component
type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// HealthResponse represents the full health check response
type HealthResponse struct {
	Message     string       `json:"message"`
	Environment string       `json:"environment"`
	Version     string       `json:"version"`
	Database    HealthStatus `json:"database"`
}
