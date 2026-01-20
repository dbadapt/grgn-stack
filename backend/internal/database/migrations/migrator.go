package migrations

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Migration represents a single database migration
type Migration struct {
	Version     int
	Description string
	Up          func(ctx context.Context, tx neo4j.ManagedTransaction) error
	Down        func(ctx context.Context, tx neo4j.ManagedTransaction) error
}

// Neo4jDB interface defines required database operations for migrations
type Neo4jDB interface {
	ExecuteWrite(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error)
	ExecuteRead(ctx context.Context, work neo4j.ManagedTransactionWork) (any, error)
}

// Migrator handles database migrations
type Migrator struct {
	db         Neo4jDB
	migrations []Migration
}

// NewMigrator creates a new migration manager
func NewMigrator(db Neo4jDB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []Migration{},
	}
}

// Register adds a migration to the migrator
func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// ensureMigrationTable creates the migration tracking table if it doesn't exist
func (m *Migrator) ensureMigrationTable(ctx context.Context) error {
	_, err := m.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Create constraint for unique migration versions
		query := `
			CREATE CONSTRAINT migration_version_unique IF NOT EXISTS
			FOR (m:Migration) REQUIRE m.version IS UNIQUE
		`
		_, err := tx.Run(ctx, query, nil)
		return nil, err
	})

	return err
}

// getAppliedVersions retrieves all applied migration versions
func (m *Migrator) getAppliedVersions(ctx context.Context) (map[int]bool, error) {
	result, err := m.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (m:Migration) RETURN m.version as version`
		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		versions := make(map[int]bool)
		for result.Next(ctx) {
			record := result.Record()
			version, _ := record.Get("version")
			if v, ok := version.(int64); ok {
				versions[int(v)] = true
			}
		}

		return versions, result.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.(map[int]bool), nil
}

// recordMigration records that a migration has been applied
func (m *Migrator) recordMigration(ctx context.Context, version int, description string) error {
	_, err := m.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (m:Migration {
				version: $version,
				description: $description,
				applied_at: datetime()
			})
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"version":     version,
			"description": description,
		})
		return nil, err
	})

	return err
}

// removeMigrationRecord removes a migration record (for rollback)
func (m *Migrator) removeMigrationRecord(ctx context.Context, version int) error {
	_, err := m.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (m:Migration {version: $version}) DELETE m`
		_, err := tx.Run(ctx, query, map[string]any{"version": version})
		return nil, err
	})

	return err
}

// Up runs all pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	// Ensure migration tracking is set up
	if err := m.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	// Get already applied migrations
	applied, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// Apply pending migrations
	for _, migration := range m.migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d already applied, skipping", migration.Version)
			continue
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Description)

		start := time.Now()

		// Run the migration
		_, err := m.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return nil, migration.Up(ctx, tx)
		})
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		// Record the migration
		if err := m.recordMigration(ctx, migration.Version, migration.Description); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		duration := time.Since(start)
		log.Printf("Migration %d completed in %v", migration.Version, duration)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down(ctx context.Context) error {
	// Get applied migrations
	appliedVersions, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	if len(appliedVersions) == 0 {
		log.Println("No migrations to roll back")
		return nil
	}

	// Find the highest applied version
	maxVersion := 0
	for v := range appliedVersions {
		if v > maxVersion {
			maxVersion = v
		}
	}

	// Find the migration to roll back
	var targetMigration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == maxVersion {
			targetMigration = &m.migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %d not found in registered migrations", maxVersion)
	}

	if targetMigration.Down == nil {
		return fmt.Errorf("migration %d has no down function", maxVersion)
	}

	log.Printf("Rolling back migration %d: %s", targetMigration.Version, targetMigration.Description)

	// Run the rollback
	_, err = m.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, targetMigration.Down(ctx, tx)
	})
	if err != nil {
		return fmt.Errorf("rollback of migration %d failed: %w", maxVersion, err)
	}

	// Remove migration record
	if err := m.removeMigrationRecord(ctx, maxVersion); err != nil {
		return fmt.Errorf("failed to remove migration %d record: %w", maxVersion, err)
	}

	log.Printf("Migration %d rolled back successfully", maxVersion)
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status(ctx context.Context) error {
	applied, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	log.Println("Migration Status:")
	log.Println("=================")

	for _, migration := range m.migrations {
		status := "[ ]"
		if applied[migration.Version] {
			status = "[✓]"
		}
		log.Printf("%s Version %d: %s", status, migration.Version, migration.Description)
	}

	return nil
}
