# Database Design - Neo4j Graph Model

## Overview

This template uses Neo4j graph database to model users, authentication, and domain entities. The graph structure allows flexible relationships between entities and efficient traversal queries.

> **Note:** This document shows the **base template schema**. The User and AuthProvider nodes are included by default. Add your own domain-specific nodes using the [Schema Workflow](SCHEMA-WORKFLOW.md).

## Visual Design Tools

### Recommended: Arrows.app

- Online tool: https://arrows.app
- Create and save graph models as JSON
- Export to Cypher for implementation
- Save designs in `/schema/graph-models/` directory

### Neo4j Browser

- Access: http://localhost:7474
- View current schema: `CALL db.schema.visualization()`
- Explore relationships and data

## Graph Model

### Node Types

#### 1. User

Represents registered users in the platform.

**Labels:** `User`

**Properties:**

```
- id: String (UUID, unique)
- email: String (unique, required)
- name: String (optional)
- passwordHash: String (for local auth)
- emailVerified: Boolean
- createdAt: DateTime
- updatedAt: DateTime
```

**Constraints:**

- UNIQUE: `User.id`
- UNIQUE: `User.email`
- EXISTS: `User.email`

**Indexes:**

- `User.email` (for fast lookup)
- `User.createdAt` (for sorting)

---

#### 2. AuthProvider

Represents external authentication providers (Google, Apple, SAML, etc.)

**Labels:** `AuthProvider`

**Properties:**

```
- id: String (UUID)
- provider: String (google, apple, saml, local)
- providerId: String (external user ID from provider)
- providerEmail: String
- accessToken: String (encrypted)
- refreshToken: String (encrypted)
- tokenExpiry: DateTime
- createdAt: DateTime
- updatedAt: DateTime
```

**Constraints:**

- UNIQUE: `AuthProvider.id`
- COMPOSITE UNIQUE: `(AuthProvider.provider, AuthProvider.providerId)`

---

## Relationships

### User ➔ AuthProvider

```cypher
(:User)-[:AUTHENTICATED_BY]->(:AuthProvider)
```

**Properties:**

- `linkedAt: DateTime` - When the auth provider was linked
- `isPrimary: Boolean` - Whether this is the primary auth method

**Meaning:** User can authenticate using this provider

---

## Sample Queries

### Create User with Auth Provider

```cypher
CREATE (u:User {
  id: randomUUID(),
  email: 'user@example.com',
  name: 'John Doe',
  emailVerified: true,
  createdAt: datetime(),
  updatedAt: datetime()
})
CREATE (a:AuthProvider {
  id: randomUUID(),
  provider: 'google',
  providerId: 'google-user-id-123',
  providerEmail: 'user@gmail.com',
  createdAt: datetime(),
  updatedAt: datetime()
})
CREATE (u)-[:AUTHENTICATED_BY {linkedAt: datetime(), isPrimary: true}]->(a)
RETURN u, a
```

### Find User by Auth Provider

```cypher
MATCH (u:User)-[:AUTHENTICATED_BY]->(a:AuthProvider)
WHERE a.provider = 'google' AND a.providerId = $providerId
RETURN u
```

### Get All Auth Providers for User

```cypher
MATCH (u:User {id: $userId})-[:AUTHENTICATED_BY]->(a:AuthProvider)
RETURN a
ORDER BY a.createdAt
```

---

## Extending the Schema

### Adding Your Domain Nodes

1. **Design visually** in [Arrows.app](https://arrows.app)
2. **Export JSON** to `schema/graph-models/your-model.json`
3. **Create migration** in `backend/internal/database/migrations/`
4. **Update this document** with your new nodes

### Example: Adding a Post Node

```cypher
// Migration Up
CREATE CONSTRAINT post_id_unique IF NOT EXISTS
FOR (p:Post) REQUIRE p.id IS UNIQUE;

CREATE INDEX post_created_at IF NOT EXISTS
FOR (p:Post) ON (p.createdAt);

// Relationship
(:User)-[:AUTHORED]->(:Post)
```

---

## Best Practices

### 1. Naming Conventions

- **Nodes:** PascalCase (`User`, `AuthProvider`, `BlogPost`)
- **Relationships:** UPPER_SNAKE_CASE (`AUTHENTICATED_BY`, `CREATED_BY`)
- **Properties:** camelCase (`createdAt`, `emailVerified`)

### 2. Use Relationship Direction Meaningfully

- Good: `(:User)-[:AUTHORED]->(:Post)`
- Bad: `(:User)-[:HAS]->(:Post)` (too generic)

### 3. Avoid Deep Traversals

- Keep important data within 3 hops
- Denormalize when needed for performance

### 4. Index Heavy-Read Properties

- Email, IDs that are frequently looked up
- DateTime properties used for sorting

### 5. Use Composite Keys for External References

```cypher
CREATE CONSTRAINT auth_provider_unique IF NOT EXISTS
FOR (a:AuthProvider)
REQUIRE (a.provider, a.providerId) IS UNIQUE
```

---

## Tools & Commands

### View Current Schema

```bash
# Start Neo4j browser
http://localhost:7474

# Run in Neo4j browser:
CALL db.schema.visualization()

# List all constraints:
SHOW CONSTRAINTS

# List all indexes:
SHOW INDEXES
```

### Performance Analysis

```cypher
# Explain query execution
EXPLAIN
MATCH (u:User)-[:AUTHENTICATED_BY]->(a:AuthProvider)
WHERE a.provider = 'google'
RETURN u, a

# Profile with actual execution stats
PROFILE
MATCH (u:User)-[:AUTHENTICATED_BY]->(a:AuthProvider)
WHERE a.provider = 'google'
RETURN u, a
```

---

## Migrations

Migrations are located in `backend/internal/database/migrations/`.

### Current Migrations

- ✅ `001_initial_schema.go` - User node with basic constraints

### Creating New Migrations

See [migrations README](backend/internal/database/migrations/README.md) for instructions.
