# Schema Design Workflow

## Overview

This document describes the collaborative workflow for designing and implementing schemas in GRGN Stack, enabling you to work visually while GitHub Copilot handles code generation.

---

## The Two-Schema System

GRGN Stack uses **two interconnected schemas**:

1. **Neo4j Graph Schema** (DATABASE.md) - Nodes and relationships in the graph database
2. **GraphQL Schema** (schema/schema.graphql) - API types and operations

Both must stay in sync, but serve different purposes:

- **Graph = Storage**: How data is structured and related in Neo4j
- **GraphQL = API**: How clients query and mutate data

---

## Collaborative Workflow

### Your Role: Visual Design

**Tool: [Arrows.app](https://arrows.app)**

1. **Create/Edit Graph Models**
   - Go to https://arrows.app
   - Design nodes, relationships, and properties
   - Use colors to group related entities
   - Add property constraints (unique, required)

2. **Save Your Designs**
   - Export as JSON from Arrows.app
   - Save to `schema/graph-models/[descriptive-name].json`
   - Commit to git for version control

3. **Share Context with Copilot**
   - Describe what you designed in natural language
   - Mention the JSON file name
   - Explain business requirements or constraints

### Copilot's Role: Code Generation

**Copilot will:**

1. **Read your Arrows.app JSON** from `schema/graph-models/`
2. **Update DATABASE.md** with formal documentation
3. **Create/Update GraphQL Schema** in `schema/schema.graphql`
4. **Generate Migration Files** in `backend/internal/database/migrations/`
5. **Implement Resolvers** in `backend/internal/graphql/resolver/`
6. **Generate Repository Code** in `backend/internal/repository/`

---

## Step-by-Step Example

### Step 1: Design in Arrows.app

**Scenario:** Add authentication providers to the system

1. Open https://arrows.app
2. Create nodes:
   - `User` node (if not exists)
   - `AuthProvider` node with properties:
     - id (string)
     - provider (string)
     - providerId (string)
     - accessToken (string)
     - createdAt (datetime)

3. Create relationship:
   - `User` -[:AUTHENTICATED_BY]-> `AuthProvider`
   - Add properties: linkedAt, isPrimary

4. Export as JSON → Save to `schema/graph-models/auth-providers.json`

### Step 2: Tell Copilot

In chat or comment:

```
I've created a new graph model for authentication providers in
schema/graph-models/auth-providers.json. Please:

1. Read the arrows.app model
2. Update DATABASE.md with the new nodes/relationships
3. Add AuthProvider type to schema.graphql
4. Create a migration for the auth provider constraints
5. Generate repository methods for linking auth providers
```

### Step 3: Copilot Generates Code

Copilot will:

- Parse your JSON model
- Generate consistent code across all layers
- Create migration file: `002_auth_providers.go`
- Update GraphQL types and resolvers
- Implement repository methods

### Step 4: Review & Iterate

- Review generated code
- Test with `npm run generate` (GraphQL)
- Run migration: `./bin/migrate`
- If changes needed, update Arrows.app model and repeat

---

## Workflow Best Practices

### For You (Visual Design)

✅ **DO:**

- Use descriptive node labels (PascalCase: `ScavengerHunt`)
- Use descriptive relationship types (UPPERCASE: `OWNS`, `AUTHENTICATED_BY`)
- Add all important properties to nodes
- Use colors in Arrows.app to group related concepts
- Save intermediate versions as you work
- Document business rules in model descriptions

❌ **DON'T:**

- Skip property definitions (add them to nodes)
- Use generic names like "Thing" or "Item"
- Forget relationship directions
- Mix visual concerns with data structure

### For Copilot (Code Generation)

✅ **DO:**

- Read Arrows.app JSON before generating
- Maintain consistency across GraphQL and Neo4j
- Generate migrations for all schema changes
- Update documentation (DATABASE.md)
- Create tests for new repository methods
- Follow existing naming conventions

❌ **DON'T:**

- Generate code without reading the visual model
- Skip migration files
- Create breaking changes to existing schema
- Forget to update both schemas (Neo4j + GraphQL)

---

## Files & Locations

| Purpose             | Location                                    | Who Updates      | Format      |
| ------------------- | ------------------------------------------- | ---------------- | ----------- |
| Visual models       | `schema/graph-models/*.json`                | You (Arrows.app) | JSON        |
| Graph documentation | `DATABASE.md`                               | Copilot          | Markdown    |
| GraphQL schema      | `schema/schema.graphql`                     | Copilot          | GraphQL SDL |
| Migrations          | `backend/internal/database/migrations/*.go` | Copilot          | Go          |
| Resolvers           | `backend/internal/graphql/resolver/*.go`    | Copilot          | Go          |
| Repositories        | `backend/internal/repository/*.go`          | Copilot          | Go          |

---

## Common Scenarios

### Scenario 1: Add New Entity Type

**You:**

1. Design in Arrows.app (new node type)
2. Add relationships to existing nodes
3. Save JSON
4. Tell Copilot: "I've added X entity in [filename].json, please implement"

**Copilot:**

- Updates DATABASE.md
- Adds GraphQL type
- Creates migration
- Generates repository
- Implements resolvers

### Scenario 2: Modify Existing Relationship

**You:**

1. Edit model in Arrows.app
2. Change relationship properties or direction
3. Re-export JSON (overwrite old file)
4. Tell Copilot: "I've updated the X relationship, please create migration"

**Copilot:**

- Creates migration to alter relationship
- Updates documentation
- Modifies affected repository methods
- Updates GraphQL schema if needed

### Scenario 3: Complex Multi-Entity Feature

**Example:** Implement complete order/payment flow

**You:**

1. Create comprehensive Arrows.app model:
   - User -[:PLACED]-> Order
   - Order -[:CONTAINS]-> QRCode
   - Order -[:PAYMENT_VIA]-> PaymentProvider
2. Save as `order-flow.json`
3. Describe business rules to Copilot

**Copilot:**

- Breaks down into multiple migrations
- Implements in layers (entities, then relationships)
- Creates transaction handling code
- Adds error cases
- Generates comprehensive tests

---

## Quick Reference Commands

### View Current Schema (Neo4j)

```cypher
CALL db.schema.visualization()
```

### Generate GraphQL Code

```bash
npm run generate
```

### Run Database Migration

```bash
./bin/migrate
```

### View All Constraints

```cypher
SHOW CONSTRAINTS
```

### Export Current Neo4j Schema

```cypher
CALL apoc.meta.graph()
```

---

## Tips for Effective Collaboration

### When Describing to Copilot

Be specific about:

1. **What changed** - New entity? Modified relationship? New properties?
2. **Where to look** - JSON filename in graph-models/
3. **Context** - Business rules, constraints, validation needs
4. **Scope** - Just migration? Full implementation? Tests too?

**Example:**

```
I've added a ScavengerHunt model in schema/graph-models/scavenger-hunt.json.
It has checkpoints that users can scan. Each checkpoint has points and order.
Please:
1. Create migration for ScavengerHunt and Checkpoint nodes
2. Add to GraphQL schema with queries for active hunts
3. Implement repository methods for creating hunts and tracking progress
4. Add tests for progress tracking logic
```

### When Copilot is Stuck

If Copilot doesn't understand the visual model:

1. Check JSON exports correctly from Arrows.app
2. Describe the model in plain English
3. Point to similar existing code as reference
4. Ask Copilot to read the JSON file first

---

## Database vs GraphQL Schema Mapping

### Neo4j Node → GraphQL Type

**Neo4j:**

```cypher
(:User {
  id: String,
  email: String,
  name: String,
  createdAt: DateTime
})
```

**GraphQL:**

```graphql
type User {
  id: ID!
  email: String!
  name: String
  createdAt: Time!
}
```

### Neo4j Relationship → GraphQL Field

**Neo4j:**

```cypher
(:User)-[:OWNS]->(:QRCode)
```

**GraphQL:**

```graphql
type User {
  qrCodes: [QRCode!]!
}

type QRCode {
  user: User!
}
```

### Complex Queries

**Neo4j Cypher:**

```cypher
MATCH (u:User)-[:OWNS]->(qr:QRCode)
WHERE qr.status = 'ACTIVE'
RETURN u, qr
```

**GraphQL Query:**

```graphql
query {
  myQRCodes(status: ACTIVE) {
    items {
      id
      code
      status
    }
  }
}
```

---

## Current Status

### ✅ Completed

- Initial User and QRCode nodes
- Basic OWNS relationship
- GraphQL schema for core entities
- Migration framework setup

### 🚧 In Progress

- Setting up database instance (Docker Compose)

### ⏳ Planned

- Auth provider model (Google, Apple, SAML)
- Order and payment flow
- Application-specific nodes:
  - Inventory items
  - Reservations
  - Scavenger hunts with checkpoints

---

## Next Steps

1. **Create Core Model**
   - Open Arrows.app
   - Import existing User/QRCode structure
   - Save as `schema/graph-models/core-model.json`

2. **Design Auth Providers**
   - Add AuthProvider nodes
   - Link to Users
   - Save as `schema/graph-models/auth-model.json`

3. **Implement Database Instance**
   - Add Neo4j to docker-compose
   - Configure for development
   - Test connection

4. **Continue with Order Flow**
   - Design in Arrows.app
   - Let Copilot implement
   - Test end-to-end

---

## Resources

- [Arrows.app](https://arrows.app) - Visual graph modeling tool
- [Neo4j Documentation](https://neo4j.com/docs/)
- [GraphQL Schema Guide](https://graphql.org/learn/schema/)
- [gqlgen Documentation](https://gqlgen.com/) - Go GraphQL generator
- [DATABASE.md](./DATABASE.md) - Full Neo4j schema documentation
- [GRAPHQL.md](./GRAPHQL.md) - GraphQL setup guide

---

## Troubleshooting

### Arrows.app Export Issues

- **Problem:** JSON missing properties
- **Solution:** Ensure properties are added to nodes in Arrows.app before export

### Schema Out of Sync

- **Problem:** GraphQL and Neo4j don't match
- **Solution:** Ask Copilot to "reconcile schemas" by reading DATABASE.md and schema.graphql

### Migration Conflicts

- **Problem:** Migration fails due to existing constraints
- **Solution:** Create migration to drop old constraint first, then add new one

### Generated Code Not Working

- **Problem:** Resolvers don't compile
- **Solution:** Run `npm run generate` to regenerate GraphQL code, check resolver implementation

---

## Version History

- **v1.0** (2026-01-19) - Initial collaborative workflow established
- Focus: Visual design (you) + Code generation (Copilot)
- Tools: Arrows.app for design, Neo4j + GraphQL for implementation
