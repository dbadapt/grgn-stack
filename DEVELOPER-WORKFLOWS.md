# Developer Workflows Guide

## Overview

This guide covers end-to-end schema-first development using the `grgn` CLI tool. From GraphQL schema design to database deployment, this workflow ensures consistent, type-safe development across all services.

---

## 🤝 Team Collaboration

### Branching Strategy

**For schema changes:**
```bash
# Create feature branch for schema work
git checkout -b feature/schema-identity-profiles

# Standard feature branches
git checkout -b feature/user-dashboard
```

### Pull Request Process

**Schema Change PR Checklist:**
- [ ] GraphQL schema changes in `model/*.graphql`
- [ ] Migration file created via `grgn migrate create`
- [ ] All migrations applied via `grgn migrate up`
- [ ] Migration tested locally with `grgn migrate status`
- [ ] Code generated via `npm run generate:backend`
- [ ] Resolvers implemented and tested

### Breaking Changes

- Use semantic versioning for schema changes
- Document breaking changes in CHANGELOG.md
- Deploy migrations in order: dev → staging → production

---

## 🚀 Complete Feature Development Workflow

This workflow takes you from GraphQL schema to a deployed feature.

### Phase 1: Schema Design

**Create your GraphQL types:**
```
services/core/identity/model/
├── enums.graphql
├── types.graphql
└── inputs.graphql
```

**Example - Adding a new field to User:**
```graphql
# services/core/identity/model/types.graphql
type User {
  id: ID!
  email: String!
  name: String
  avatarUrl: String
  status: UserStatus!
  profile: UserProfile          # NEW
  createdAt: DateTime!
  updatedAt: DateTime!
}
```

### Phase 2: Code Generation

**Generate backend Go code:**
```bash
grgn generate:backend
# OR
npm run generate:backend
```

**Optional - Generate frontend types:**
```bash
npm run generate:frontend
```

### Phase 3: Implementation

**Implement resolvers in:**
```
services/core/identity/controller/
├── resolver.go
└── generated/
    ├── generated.go
    └── models_gen.go
```

**Basic resolver structure:**
```go
// services/core/identity/controller/resolver.go
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
    // Your implementation here
    // Delegate to service layer
    return r.userService.GetCurrentUser(ctx)
}
```

### Phase 4: Database Migration

**Create migration:**
```bash
grgn migrate create user_profile_schema --app core/identity
```

**Edit migration file:**
```
services/core/identity/migrations/001_user_profile_schema.cypher
```

**Add Neo4j constraints and indexes:**
```cypher
// Create constraints
CREATE CONSTRAINT user_profile_id_unique IF NOT EXISTS
FOR (p:UserProfile) REQUIRE p.id IS UNIQUE;

// Create indexes
CREATE INDEX user_profile_user_id IF NOT EXISTS
FOR (p:UserProfile) ON (p.userId);
```

### Phase 5: Deploy Migration

**Apply migration to database:**
```bash
grgn migrate up --app core/identity
# OR apply all pending
grgn migrate up
```

**Verify deployment:**
```bash
grgn migrate status --app core/identity
```

**Expected output:**
```
📊 Migration Status
MIGRATION                           STATUS     APPLIED AT
identity/001_user_profile_schema     ✅ Applied  2026-01-26 15:30:22
```

---

## 🔄 Workflow Diagram

```
📝 Schema Design → 🏗️ Code Gen → 🎯 Implement → 🗄️ Migrate → ✅ Verify
      ↓               ↓            ↓           ↓          ↓
   .graphql       grgn gen     Resolvers   grgn up  grgn status
   
   📍 File:                     📍 Files:
   services/                      services/core/identity/
   core/identity/                 ├── controller/
   model/                         ├── repository/
                                  ├── generated/
                                  └── migrations/
```

---

## 🚨 Rollback Workflow

Use when a migration causes issues or needs to be reworked.

### Safe Rollback Procedure

**1. Check current status:**
```bash
grgn migrate status
```

**2. Rollback specific app migration:**
```bash
grgn migrate down --app core/identity
```

**3. Verify rollback:**
```bash
grgn migrate status
```

**4. Fix migration file:**
- Edit the migration file in `migrations/001_*.cypher`
- Fix constraints, indexes, or data issues

**5. Re-deploy:**
```bash
grgn migrate up --app core/identity
```

### Rollback Safety Checklist

Before rolling back:
- [ ] Backup current database state
- [ ] Identify affected users/systems
- [ ] Communicate downtime window
- [ ] Test rollback in staging first

After rollback:
- [ ] Verify data integrity
- [ ] Test affected features
- [ ] Monitor system stability

---

## 📚 Development Patterns

### Adding New Fields to Existing Types

**1. Update GraphQL schema:**
```graphql
type User {
  id: ID!
  email: String!
  phoneNumber: String     # NEW
}
```

**2. Generate code:**
```bash
npm run generate:backend
```

**3. Update resolvers:**
```go
// Add to user creation logic
user.PhoneNumber = input.PhoneNumber
```

**4. Add migration:**
```bash
grgn migrate create user_phone_field --app core/identity
```

### Creating New Relationships

**Example - User ↔ UserProfile:**
```graphql
type User {
  id: ID!
  profile: UserProfile!
}

type UserProfile {
  id: ID!
  user: User!
}
```

**Migration relationship:**
```cypher
// Create relationship index
CREATE INDEX rel_user_has_profile IF NOT EXISTS
FOR ()-[r:HAS_PROFILE]->() ON (r.createdAt);
```

---

## 🛠️ Troubleshooting

### Migration Conflicts

**Issue: Two migrations with same number**
```bash
# Error: duplicate migration ID
grgn migrate up
```

**Solution:**
```bash
# Check migration directory
ls services/core/identity/migrations/

# Renumber manually if needed
mv 002_new_field.cypher 003_new_field.cypher
```

### Database Connection Issues

**Issue: Authentication failure**
```bash
grgn migrate up
Error: Neo4jError: Unauthorized
```

**Solution:**
```bash
# Check .env file
cat .env | grep NEO4J

# Verify Neo4j is running
podman ps | grep neo4j

# Test connection
podman exec grgn-stack_neo4j_1 cypher-shell -u neo4j -p devpassword "RETURN 1"
```

### Schema Generation Failures

**Issue: GraphQL syntax error**
```bash
grgn generate:backend
Error: syntax error: Unexpected Name
```

**Solution:**
```bash
# Validate GraphQL schema
# Check for missing commas, brackets, etc.
# Use GraphQL Playground to test schema at http://localhost:8080/graphql

# Clear and regenerate
rm -rf services/core/identity/generated/
grgn generate:backend
```

### Migration Not Applied

**Issue: Migration shows as pending**
```bash
grgn migrate status
identity/001_user_schema ⏳ Pending
```

**Solution:**
```bash
# Check for syntax errors in migration
podman exec grgn-stack_neo4j_1 cypher-shell \
  -u neo4j -p devpassword \
  -f services/core/identity/migrations/001_user_schema.cypher

# Force re-run
grgn migrate up --app core/identity
```

### Permission Denied on Migrations

**Issue: File permission errors**
```bash
grgn migrate up
Error: permission denied: migrations/001_*.cypher
```

**Solution:**
```bash
# Check file permissions
ls -la services/core/identity/migrations/

# Fix permissions
chmod 644 services/core/identity/migrations/*.cypher
```

---

## 🎯 Quick Reference

### Essential Commands

```bash
# Development workflow
grgn migrate create <name> --app <service>
grgn migrate up [--app <service>]
grgn migrate status [--app <service>]
grgn migrate down [--app <service>]

# Code generation
npm run generate:backend
npm run generate:frontend

# Database seeding
grgn seed [--clean]
```

### File Locations

```
services/core/identity/
├── model/           # GraphQL schemas
├── controller/       # Resolver implementations
├── repository/       # Database access layer
├── service/         # Business logic
└── migrations/       # Database migrations
```

---

## 🚀 Next Steps

After mastering this workflow:

1. **Advanced Queries**: Add DataLoader for N+1 problems
2. **Subscriptions**: Implement real-time updates
3. **Testing**: Add unit tests for resolvers
4. **Performance**: Optimize queries and indexes
5. **Monitoring**: Add metrics and logging

---

**Built with ❤️ for developer productivity**