# Database Migrations

This directory contains the database migration system for GRGN Stack's Neo4j database.

## Overview

The migration system provides version control for database schema changes, including:

- Node constraints (uniqueness, existence)
- Indexes for performance
- Initial data seeding
- Schema evolution tracking

## Migration Tool

### Building

```bash
go build -o ./bin/migrate ./cmd/migrate
```

### Usage

**Check migration status:**

```bash
grgn migrate:status
```

**Run pending migrations:**

```bash
grgn migrate
```

**Rollback last migration:**

```bash
grgn migrate:down
```

## Creating New Migrations

1. Create a new file in `services/{domain}/{app}/migrations/` (or root `migrations/` for core) following the naming convention:

   ```
   XXX_description.go
   ```

   where XXX is the migration version number (e.g., `002_add_payment_tables.go`)

2. Define your migration:

   ```go
   package migrations

   import (
       "context"
       "github.com/neo4j/neo4j-go-driver/v5/neo4j"
   )

   var MigrationXXX = Migration{
       Version:     XXX,
       Description: "Description of what this migration does",

       Up: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
           // Your forward migration Cypher queries
           _, err := tx.Run(ctx, `
               CREATE CONSTRAINT my_constraint IF NOT EXISTS
               FOR (n:NodeType) REQUIRE n.property IS UNIQUE
           `, nil)
           return err
       },

       Down: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
           // Your rollback Cypher queries
           _, err := tx.Run(ctx, `
               DROP CONSTRAINT my_constraint IF EXISTS
           `, nil)
           return err
       },
   }
   ```

3. Register your migration in `registry.go`:
   ```go
   func GetAllMigrations() []Migration {
       return []Migration{
           Migration001InitialSchema,
           MigrationXXX,  // Add your new migration here
       }
   }
   ```

## Existing Migrations

### 001_initial_schema.go

Initial database schema including:

**Constraints:**

- `User.email` - Unique
- `User.id` - Unique
- `QRCode.id` - Unique
- `QRCode.code` - Unique

**Indexes:**

- `User.createdAt` - Time-based queries
- `QRCode.userId` - User-to-QRCode lookups
- `QRCode.status` - Status filtering
- `QRCode.applicationType` - Application type filtering
- `QRCode(userId, status)` - Composite index for common queries
- `QRCode.createdAt` - Time-based queries

## Migration Tracking

Migrations are tracked in the database using `Migration` nodes:

```cypher
(:Migration {
    version: Int,
    description: String,
    applied_at: DateTime
})
```

## Best Practices

1. **Always provide a Down function** - Enables rollback capability
2. **Use IF NOT EXISTS / IF EXISTS** - Makes migrations idempotent
3. **Test migrations** - Run up/down/up cycle in development
4. **One logical change per migration** - Easier to understand and rollback
5. **Never modify applied migrations** - Create a new migration instead
6. **Document breaking changes** - Add comments for schema changes
7. **Consider performance** - Large data migrations may need batching

## Environment-Specific Migrations

The migration tool uses the same configuration system as the main application:

- Development: `.env.development`
- Staging: `.env.staging`
- Production: `.env.production`

Set `GRGN_STACK_SERVER_ENVIRONMENT` to control which config is loaded.

## Integration with CI/CD

Migrations should run automatically in the CI/CD pipeline before deployment:

```yaml
# Example GitHub Actions step
- name: Run Database Migrations
  run: |
    ./bin/migrate -command=up
  env:
    GRGN_STACK_DATABASE_NEO4J_URI: ${{ secrets.NEO4J_URI }}
    GRGN_STACK_DATABASE_NEO4J_USERNAME: ${{ secrets.NEO4J_USERNAME }}
    GRGN_STACK_DATABASE_NEO4J_PASSWORD: ${{ secrets.NEO4J_PASSWORD }}
```

## Troubleshooting

**Migration fails mid-execution:**

- Check Neo4j logs for constraint violations or syntax errors
- Manually inspect the database state
- Fix the issue and re-run (migrations are idempotent)
- If needed, manually rollback and create a fixed migration

**Cannot connect to Neo4j:**

- Verify connection string in environment config
- Check Neo4j is running: `docker ps`
- Test connection: `./bin/migrate -command=status`

**Constraint already exists:**

- Ensure you're using `IF NOT EXISTS` clauses
- Check if migration was partially applied
- Review migration tracking nodes: `MATCH (m:Migration) RETURN m`
