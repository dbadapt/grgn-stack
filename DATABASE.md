# Database Design - Neo4j Graph Model

## Overview

This template uses Neo4j graph database to model domain entities. The graph structure allows flexible relationships between entities and efficient traversal queries.

> **Note:** This document shows the **base template schema**. Add your own domain-specific nodes by extending the GraphQL schemas in `services/{domain}/{app}/model/`.

## Design Tools

### Neo4j Browser

- Access: http://localhost:7474
- View current schema: `CALL db.schema.visualization()`
- Explore relationships and data

## Graph Model

### Node Types

#### TODO: Define Core Domain Nodes
The core domain nodes (e.g., User, Tenant, etc.) will be defined here once the schema design is finalized.

---

## Relationships

#### TODO: Define Relationships
Relationships between domain nodes will be defined here.

---

## Sample Queries

#### TODO: Add Sample Queries
Add Cypher queries for common operations once the model is defined.

---

## Extending the Schema

### Adding Your Domain Nodes

1. **Define types** in `services/{domain}/{app}/model/*.graphql`
2. **Create migration** in `services/{domain}/{app}/migrations/` (or root `migrations/` for core)
3. **Run code generation** with `npm run generate`
4. **Update this document** with your new nodes

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
MATCH (n) RETURN n

# Profile with actual execution stats
PROFILE
MATCH (n) RETURN n
```

---

## Migrations

Migrations are located in `services/{domain}/{app}/migrations/` (or root `migrations/` for core).

### Current Migrations

- (No active migrations)

### Creating New Migrations

See [Database Schema Management](mvc_design.md#11-database-schema-management) for instructions.
