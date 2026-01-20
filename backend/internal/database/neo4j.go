package database

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/yourusername/grgn-stack/pkg/config"
)

// Neo4jDB wraps the Neo4j driver and provides database operations
type Neo4jDB struct {
	driver neo4j.DriverWithContext
	config *config.Config
}

// NewNeo4jDB creates a new Neo4j database connection with connection pooling
func NewNeo4jDB(cfg *config.Config) (*Neo4jDB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Configure connection pool settings
	poolConfig := func(conf *neo4j.Config) {
		conf.MaxConnectionPoolSize = 50
		conf.MaxConnectionLifetime = 1 * time.Hour
		conf.ConnectionAcquisitionTimeout = 2 * time.Minute
		conf.SocketConnectTimeout = 5 * time.Second
		conf.SocketKeepalive = true

		// Adjust pool size for production vs development
		if cfg.IsProduction() {
			conf.MaxConnectionPoolSize = 100
		} else if cfg.Server.Environment == "development" {
			conf.MaxConnectionPoolSize = 10
		}
	}

	// Create the driver
	driver, err := neo4j.NewDriverWithContext(
		cfg.Database.Neo4jURI,
		neo4j.BasicAuth(cfg.Database.Neo4jUsername, cfg.Database.Neo4jPassword, ""),
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	db := &Neo4jDB{
		driver: driver,
		config: cfg,
	}

	return db, nil
}

// VerifyConnectivity checks if the database is accessible and responsive
func (db *Neo4jDB) VerifyConnectivity(ctx context.Context) error {
	if db.driver == nil {
		return fmt.Errorf("driver is not initialized")
	}

	// Verify connectivity with a timeout
	err := db.driver.VerifyConnectivity(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return nil
}

// GetDriver returns the underlying Neo4j driver for advanced usage
func (db *Neo4jDB) GetDriver() neo4j.DriverWithContext {
	return db.driver
}

// ExecuteRead executes a read transaction with automatic retry
func (db *Neo4jDB) ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error) {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("read transaction failed: %w", err)
	}

	return result, nil
}

// ExecuteWrite executes a write transaction with automatic retry
func (db *Neo4jDB) ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error) {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("write transaction failed: %w", err)
	}

	return result, nil
}

// NewSession creates a new session for manual transaction management
func (db *Neo4jDB) NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext {
	return db.driver.NewSession(ctx, config)
}

// Close gracefully closes the database connection and releases resources
func (db *Neo4jDB) Close(ctx context.Context) error {
	if db.driver == nil {
		return nil
	}

	if err := db.driver.Close(ctx); err != nil {
		return fmt.Errorf("failed to close Neo4j driver: %w", err)
	}

	return nil
}

// Ping performs a simple connectivity check (alias for VerifyConnectivity)
func (db *Neo4jDB) Ping(ctx context.Context) error {
	return db.VerifyConnectivity(ctx)
}

// GetServerInfo retrieves information about the connected Neo4j server
func (db *Neo4jDB) GetServerInfo(ctx context.Context) (map[string]any, error) {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "CALL dbms.components() YIELD name, versions, edition", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	record, err := result.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read server info: %w", err)
	}

	info := make(map[string]any)
	info["name"] = record.Values[0]
	info["versions"] = record.Values[1]
	info["edition"] = record.Values[2]

	return info, nil
}
