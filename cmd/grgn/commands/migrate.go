package commands

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cobra"
	"github.com/yourusername/grgn-stack/pkg/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage Neo4j database migrations for all apps in the GRGN stack.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE:  runMigrateUp,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE:  runMigrateStatus,
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Long: `Create a new migration file for a specific app.

Examples:
  grgn migrate create add_user_roles --app core/identity
  grgn migrate create add_tenant_settings --app core/tenant`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrateCreate,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration (if supported)",
	Long: `Rollback the last applied migration.

Note: Neo4j migrations are typically not reversible. This command will 
mark the migration as unapplied but won't undo schema changes.
Use with caution and consider creating a new migration instead.`,
	RunE: runMigrateDown,
}

var (
	appFilter string
)

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	// Add flags
	migrateUpCmd.Flags().StringVar(&appFilter, "app", "", "Filter by app (e.g., core/identity)")
	migrateStatusCmd.Flags().StringVar(&appFilter, "app", "", "Filter by app (e.g., core/identity)")
	migrateCreateCmd.Flags().StringVar(&appFilter, "app", "", "App to create migration for (required, e.g., core/identity)")
	migrateCreateCmd.MarkFlagRequired("app")
	migrateDownCmd.Flags().StringVar(&appFilter, "app", "", "Filter by app (e.g., core/identity)")
}

// Migration represents a single migration file
type Migration struct {
	ID       string // e.g., "core/identity/001_user_schema"
	App      string // e.g., "core/identity"
	Filename string // e.g., "001_user_schema.cypher"
	Path     string // Full path to file
	Checksum string // SHA256 of file contents
}

// AppliedMigration represents a migration that has been applied
type AppliedMigration struct {
	ID        string
	AppliedAt time.Time
	Checksum  string
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Running migrations...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to Neo4j
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		cfg.Database.Neo4jURI,
		neo4j.BasicAuth(cfg.Database.Neo4jUsername, cfg.Database.Neo4jPassword, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %w", err)
	}
	defer driver.Close(ctx)

	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}
	fmt.Println("✅ Connected to Neo4j")

	// Ensure migration tracking exists
	if err := ensureMigrationTracking(ctx, driver); err != nil {
		return fmt.Errorf("failed to ensure migration tracking: %w", err)
	}

	// Discover migrations
	migrations, err := discoverMigrations()
	if err != nil {
		return fmt.Errorf("failed to discover migrations: %w", err)
	}

	// Filter by app if specified
	if appFilter != "" {
		var filtered []Migration
		for _, m := range migrations {
			if m.App == appFilter {
				filtered = append(filtered, m)
			}
		}
		migrations = filtered
	}

	if len(migrations) == 0 {
		fmt.Println("📭 No migrations found")
		return nil
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, driver)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find pending migrations
	appliedMap := make(map[string]AppliedMigration)
	for _, a := range applied {
		appliedMap[a.ID] = a
	}

	var pending []Migration
	for _, m := range migrations {
		if _, ok := appliedMap[m.ID]; !ok {
			pending = append(pending, m)
		}
	}

	if len(pending) == 0 {
		fmt.Println("✅ All migrations are up to date")
		return nil
	}

	fmt.Printf("📋 Found %d pending migration(s)\n", len(pending))

	// Apply pending migrations
	for _, m := range pending {
		fmt.Printf("\n⏳ Applying: %s\n", m.ID)

		if err := applyMigration(ctx, driver, m); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", m.ID, err)
		}

		fmt.Printf("✅ Applied: %s\n", m.ID)
	}

	fmt.Printf("\n🎉 Successfully applied %d migration(s)\n", len(pending))
	return nil
}

func runMigrateStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("📊 Migration Status")
	fmt.Println()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to Neo4j
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		cfg.Database.Neo4jURI,
		neo4j.BasicAuth(cfg.Database.Neo4jUsername, cfg.Database.Neo4jPassword, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %w", err)
	}
	defer driver.Close(ctx)

	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	// Discover migrations
	migrations, err := discoverMigrations()
	if err != nil {
		return fmt.Errorf("failed to discover migrations: %w", err)
	}

	// Filter by app if specified
	if appFilter != "" {
		var filtered []Migration
		for _, m := range migrations {
			if m.App == appFilter {
				filtered = append(filtered, m)
			}
		}
		migrations = filtered
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, driver)
	if err != nil {
		// If migration tracking doesn't exist yet, treat as no applied migrations
		applied = []AppliedMigration{}
	}

	appliedMap := make(map[string]AppliedMigration)
	for _, a := range applied {
		appliedMap[a.ID] = a
	}

	// Print status
	fmt.Printf("%-40s %-10s %-20s\n", "MIGRATION", "STATUS", "APPLIED AT")
	fmt.Println(strings.Repeat("-", 72))

	for _, m := range migrations {
		if a, ok := appliedMap[m.ID]; ok {
			fmt.Printf("%-40s %-10s %-20s\n", m.ID, "✅ Applied", a.AppliedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%-40s %-10s %-20s\n", m.ID, "⏳ Pending", "-")
		}
	}

	return nil
}

func discoverMigrations() ([]Migration, error) {
	var migrations []Migration

	// Search patterns for migrations
	patterns := []string{
		"services/core/*/migrations/*.cypher",
		"services/*/*/migrations/*.cypher",
		"migrations/*.cypher",
	}

	seen := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, path := range matches {
			// Normalize path
			path = filepath.ToSlash(path)

			if seen[path] {
				continue
			}
			seen[path] = true

			// Parse migration info
			m, err := parseMigration(path)
			if err != nil {
				fmt.Printf("⚠️  Skipping invalid migration: %s (%v)\n", path, err)
				continue
			}

			migrations = append(migrations, m)
		}
	}

	// Sort by ID
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	return migrations, nil
}

func parseMigration(path string) (Migration, error) {
	// Read file for checksum
	content, err := os.ReadFile(path)
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate checksum
	hash := sha256.Sum256(content)
	checksum := fmt.Sprintf("%x", hash)

	// Extract app and filename
	// Path format: services/core/identity/migrations/001_user_schema.cypher
	parts := strings.Split(filepath.ToSlash(path), "/")

	var app, filename string

	// Find migrations directory and work backwards
	for i, part := range parts {
		if part == "migrations" && i > 0 && i < len(parts)-1 {
			// App is everything between services/ and /migrations
			if i >= 2 && parts[i-2] == "services" {
				app = parts[i-2+1] + "/" + parts[i-1]
			} else if i >= 1 {
				app = parts[i-1]
			}
			filename = parts[i+1]
			break
		}
	}

	if app == "" || filename == "" {
		return Migration{}, fmt.Errorf("invalid migration path structure")
	}

	// Remove .cypher extension for ID
	name := strings.TrimSuffix(filename, ".cypher")
	id := app + "/" + name

	return Migration{
		ID:       id,
		App:      app,
		Filename: filename,
		Path:     path,
		Checksum: checksum,
	}, nil
}

func ensureMigrationTracking(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.Run(ctx, `
		CREATE CONSTRAINT migration_id_unique IF NOT EXISTS
		FOR (m:Migration) REQUIRE m.id IS UNIQUE
	`, nil)

	return err
}

func getAppliedMigrations(ctx context.Context, driver neo4j.DriverWithContext) ([]AppliedMigration, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (m:Migration)
		RETURN m.id AS id, m.appliedAt AS appliedAt, m.checksum AS checksum
		ORDER BY m.id
	`, nil)
	if err != nil {
		return nil, err
	}

	var applied []AppliedMigration
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		appliedAt, _ := record.Get("appliedAt")
		checksum, _ := record.Get("checksum")

		a := AppliedMigration{
			ID:       id.(string),
			Checksum: checksum.(string),
		}

		// Handle Neo4j time type
		if t, ok := appliedAt.(time.Time); ok {
			a.AppliedAt = t
		}

		applied = append(applied, a)
	}

	return applied, result.Err()
}

func applyMigration(ctx context.Context, driver neo4j.DriverWithContext, m Migration) error {
	// Read migration file
	content, err := os.ReadFile(m.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Parse and execute statements
	statements := parseCypherStatements(string(content))

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := session.Run(ctx, stmt, nil)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w\nStatement: %s", err, stmt)
		}
	}

	// Record migration as applied
	_, err = session.Run(ctx, `
		CREATE (m:Migration {
			id: $id,
			appliedAt: datetime(),
			checksum: $checksum
		})
	`, map[string]any{
		"id":       m.ID,
		"checksum": m.Checksum,
	})

	return err
}

func parseCypherStatements(content string) []string {
	var statements []string
	var current strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comment-only lines
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || trimmed == "" {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		// Check if statement ends with semicolon
		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSuffix(strings.TrimSpace(current.String()), ";")
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}

	// Handle final statement without semicolon
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func runMigrateCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	if appFilter == "" {
		return fmt.Errorf("--app flag is required (e.g., --app core/identity)")
	}

	// Determine the migrations directory
	migrationsDir := filepath.Join("services", appFilter, "migrations")

	// Ensure the migrations directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Find the next migration number
	nextNum := 1
	entries, err := os.ReadDir(migrationsDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".cypher") {
				// Extract number from filename like "001_name.cypher"
				parts := strings.SplitN(entry.Name(), "_", 2)
				if len(parts) > 0 {
					var num int
					if _, err := fmt.Sscanf(parts[0], "%d", &num); err == nil {
						if num >= nextNum {
							nextNum = num + 1
						}
					}
				}
			}
		}
	}

	// Create the migration file
	filename := fmt.Sprintf("%03d_%s.cypher", nextNum, name)
	filePath := filepath.Join(migrationsDir, filename)

	// Generate template content
	content := fmt.Sprintf(`// ============================================
// Migration: %s/%03d_%s
// Description: [Add description here]
// Created: %s
// ============================================

// ----- CONSTRAINTS -----

// Example: Create a unique constraint
// CREATE CONSTRAINT example_id_unique IF NOT EXISTS
// FOR (e:Example) REQUIRE e.id IS UNIQUE;

// ----- INDEXES -----

// Example: Create an index
// CREATE INDEX example_status IF NOT EXISTS
// FOR (e:Example) ON (e.status);

// ----- DATA MIGRATIONS -----

// Example: Update existing data
// MATCH (e:Example) WHERE e.oldField IS NOT NULL
// SET e.newField = e.oldField
// REMOVE e.oldField;
`, appFilter, nextNum, name, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	fmt.Printf("✅ Created migration: %s\n", filePath)
	fmt.Printf("\n📝 Next steps:\n")
	fmt.Printf("   1. Edit %s to add your schema changes\n", filePath)
	fmt.Printf("   2. Run 'grgn migrate up' to apply the migration\n")
	fmt.Printf("   3. Run 'grgn migrate status' to verify\n")

	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	fmt.Println("⚠️  Rolling back last migration...")
	fmt.Println("   Note: This only marks the migration as unapplied.")
	fmt.Println("   Schema changes in Neo4j are NOT automatically reversed.")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to Neo4j
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		cfg.Database.Neo4jURI,
		neo4j.BasicAuth(cfg.Database.Neo4jUsername, cfg.Database.Neo4jPassword, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %w", err)
	}
	defer driver.Close(ctx)

	// Verify connectivity
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, driver)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		fmt.Println("📭 No migrations to rollback")
		return nil
	}

	// Filter by app if specified
	if appFilter != "" {
		var filtered []AppliedMigration
		for _, a := range applied {
			if strings.HasPrefix(a.ID, appFilter+"/") {
				filtered = append(filtered, a)
			}
		}
		applied = filtered
	}

	if len(applied) == 0 {
		fmt.Println("📭 No migrations to rollback for the specified app")
		return nil
	}

	// Get the last applied migration
	last := applied[len(applied)-1]

	fmt.Printf("\n🔙 Rolling back: %s\n", last.ID)
	fmt.Printf("   Applied at: %s\n", last.AppliedAt.Format("2006-01-02 15:04:05"))

	// Remove the migration record
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err = session.Run(ctx, `
		MATCH (m:Migration {id: $id})
		DELETE m
	`, map[string]any{"id": last.ID})
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	fmt.Printf("✅ Migration record removed: %s\n", last.ID)
	fmt.Println("\n⚠️  Remember: Schema changes have NOT been reversed.")
	fmt.Println("   You may need to manually clean up constraints/indexes if needed.")

	return nil
}
