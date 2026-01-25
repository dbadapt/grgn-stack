package shared

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// IDatabase defines the contract for database operations.
// This interface is used by other domains to interact with the database
// without depending on the concrete Neo4j implementation.
type IDatabase interface {
	// Ping performs a connectivity check
	Ping(ctx context.Context) error

	// Close gracefully closes the database connection
	Close(ctx context.Context) error

	// VerifyConnectivity checks if the database is accessible
	VerifyConnectivity(ctx context.Context) error

	// ExecuteRead executes a read transaction with automatic retry
	ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error)

	// ExecuteWrite executes a write transaction with automatic retry
	ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error)

	// NewSession creates a new session for manual transaction management
	NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext

	// GetDriver returns the underlying driver for advanced usage
	GetDriver() neo4j.DriverWithContext
}

// Ensure Neo4jDB implements IDatabase
var _ IDatabase = (*Neo4jDB)(nil)
