package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cobra"
	"github.com/yourusername/grgn-stack/pkg/config"
	shared "github.com/yourusername/grgn-stack/services/core/shared/controller"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with test data",
	Long: `Create test users, tenants, and memberships for development.

This command creates:
- 3 test users (Alice, Bob, Charlie)
- 2 test tenants (Acme Corp, Startup Inc)
- Membership relationships with various roles

Use --clean to clear existing data before seeding.`,
	RunE: runSeed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().Bool("clean", false, "Clear existing data before seeding")
}

func runSeed(cmd *cobra.Command, args []string) error {
	fmt.Println("🌱 Seeding database...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to Neo4j
	db, err := shared.NewNeo4jDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}
	defer db.Close(context.Background())

	ctx := context.Background()

	// Verify connectivity
	if err := db.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to verify database connectivity: %w", err)
	}
	fmt.Println("✅ Connected to Neo4j")

	// Check for --clean flag
	clean, _ := cmd.Flags().GetBool("clean")
	if clean {
		fmt.Println("🧹 Clearing existing data...")
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx, `
				MATCH (n) 
				WHERE NOT n:_Migration 
				DETACH DELETE n
			`, nil)
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("failed to clean data: %w", err)
		}
		fmt.Println("  ✅ Existing data cleared")
	}

	// Create test users
	users := []struct {
		email string
		name  string
	}{
		{"alice@example.com", "Alice Johnson"},
		{"bob@example.com", "Bob Smith"},
		{"charlie@example.com", "Charlie Brown"},
	}

	userIDs := make(map[string]string)

	fmt.Println("\n👥 Creating users...")
	for _, u := range users {
		id := uuid.New().String()
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, `
				MERGE (u:User {email: $email})
				ON CREATE SET 
					u.id = $id,
					u.name = $name,
					u.status = 'ACTIVE',
					u.createdAt = datetime(),
					u.updatedAt = datetime()
				ON MATCH SET
					u.name = $name,
					u.updatedAt = datetime()
				RETURN u.id as id
			`, map[string]any{"id": id, "email": u.email, "name": u.name})
			if err != nil {
				return nil, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return nil, err
			}

			returnedID, _ := record.Get("id")
			return returnedID.(string), nil
		})
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", u.email, err)
		}
		userIDs[u.email] = id
		fmt.Printf("  ✅ %s <%s>\n", u.name, u.email)
	}

	// Create test tenants with memberships
	tenants := []struct {
		name    string
		slug    string
		owner   string
		members []struct {
			email string
			role  string
		}
	}{
		{
			name:  "Acme Corp",
			slug:  "acme",
			owner: "alice@example.com",
			members: []struct {
				email string
				role  string
			}{
				{"bob@example.com", "ADMIN"},
			},
		},
		{
			name:  "Startup Inc",
			slug:  "startup",
			owner: "bob@example.com",
			members: []struct {
				email string
				role  string
			}{
				{"alice@example.com", "MEMBER"},
				{"charlie@example.com", "VIEWER"},
			},
		},
	}

	fmt.Println("\n🏢 Creating tenants...")
	for _, t := range tenants {
		tenantID := uuid.New().String()

		// Create tenant
		_, err := db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx, `
				MERGE (t:Tenant {slug: $slug})
				ON CREATE SET
					t.id = $id,
					t.name = $name,
					t.plan = 'FREE',
					t.status = 'ACTIVE',
					t.isolationMode = 'SHARED',
					t.createdAt = datetime(),
					t.updatedAt = datetime()
				ON MATCH SET
					t.name = $name,
					t.updatedAt = datetime()
				RETURN t
			`, map[string]any{"id": tenantID, "name": t.name, "slug": t.slug})
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("failed to create tenant %s: %w", t.name, err)
		}
		fmt.Printf("  ✅ %s (/%s)\n", t.name, t.slug)

		// Create owner membership
		_, err = db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx, `
				MATCH (u:User {email: $email}), (t:Tenant {slug: $slug})
				MERGE (u)-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t)
				ON CREATE SET
					m.id = $membershipId,
					m.role = 'OWNER',
					m.joinedAt = datetime()
				RETURN m
			`, map[string]any{
				"email":        t.owner,
				"slug":         t.slug,
				"membershipId": uuid.New().String(),
			})
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("failed to create owner membership: %w", err)
		}
		fmt.Printf("    👑 Owner: %s\n", t.owner)

		// Create member memberships
		for _, member := range t.members {
			_, err = db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
				_, err := tx.Run(ctx, `
					MATCH (u:User {email: $email}), (t:Tenant {slug: $slug})
					MERGE (u)-[:HAS_MEMBERSHIP]->(m:Membership)-[:IN_TENANT]->(t)
					ON CREATE SET
						m.id = $membershipId,
						m.role = $role,
						m.joinedAt = datetime()
					RETURN m
				`, map[string]any{
					"email":        member.email,
					"slug":         t.slug,
					"membershipId": uuid.New().String(),
					"role":         member.role,
				})
				return nil, err
			})
			if err != nil {
				return fmt.Errorf("failed to create member membership: %w", err)
			}
			fmt.Printf("    👤 %s: %s\n", member.role, member.email)
		}
	}

	fmt.Println("\n🎉 Seeding complete!")
	fmt.Println("\n📋 Test Data Summary:")
	fmt.Println("   Users:")
	for email, id := range userIDs {
		fmt.Printf("     • %s: %s\n", email, id)
	}

	fmt.Println("\n🧪 Test with GraphQL:")
	fmt.Printf(`
   # Start the server
   go run ./cmd/server

   # In another terminal, test queries:
   
   # Get Alice's tenants
   curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -H "X-User-ID: %s" \
     -d '{"query": "{ myTenants { id name slug memberCount } }"}'

   # Create a new tenant as Alice
   curl -X POST http://localhost:8080/graphql \
     -H "Content-Type: application/json" \
     -H "X-User-ID: %s" \
     -d '{"query": "mutation { createTenant(input: { name: \"New Corp\", slug: \"newcorp\" }) { id name } }"}'
`, userIDs["alice@example.com"], userIDs["alice@example.com"])

	return nil
}
