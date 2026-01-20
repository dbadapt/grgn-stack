# GraphQL Setup Guide

## Overview

GRGN Stack uses **GraphQL** as its primary API layer with automatic code generation for both backend (Go) and frontend (TypeScript/React).

## Architecture

- **Backend**: gqlgen generates Go resolvers from GraphQL schema
- **Frontend**: graphql-codegen generates TypeScript types and React Query hooks
- **Schema**: Single source of truth in `schema/schema.graphql`
- **Playground**: Interactive GraphQL IDE available in development

## Quick Start

### 1. Schema Development

Edit the GraphQL schema:

```bash
schema/schema.graphql
```

### 2. Generate Code

**Both backend and frontend:**

```bash
npm run generate
```

**Backend only:**

```bash
npm run generate:backend
# or
cd backend && go run github.com/99designs/gqlgen generate
```

**Frontend only:**

```bash
npm run generate:frontend
# or
cd web && npm run generate
```

### 3. Implement Resolvers

After generating, implement the resolver functions in:

```
backend/internal/graphql/resolver/schema.resolvers.go
```

## GraphQL Playground

Access the interactive playground in development:

```
http://localhost:8080/graphql
```

**Sample Query:**

```graphql
query {
  health
  me {
    id
    email
    name
  }
}
```

## Frontend Usage

### 1. Define Queries

Create `.graphql` files in `web/src/graphql/`:

```graphql
# queries.graphql
query GetCurrentUser {
  me {
    id
    email
    name
  }
}
```

### 2. Generate TypeScript Code

```bash
cd web && npm run generate
```

This generates:

- TypeScript types for all schema types
- React Query hooks for all queries/mutations
- Type-safe GraphQL operations

### 3. Use in Components

```tsx
import { useGetCurrentUserQuery } from './graphql/generated';

function UserProfile() {
  const { data, isLoading, error } = useGetCurrentUserQuery();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h1>{data.me.name}</h1>
      <p>{data.me.email}</p>
    </div>
  );
}
```

## Configuration Files

### Backend: gqlgen.yml

Located at `backend/gqlgen.yml`:

```yaml
schema:
  - ../schema/*.graphql
model:
  filename: internal/graphql/model/models_gen.go
  package: model
resolver:
  layout: follow-schema
  dir: internal/graphql/resolver
  package: resolver
exec:
  filename: internal/graphql/generated.go
  package: graphql
```

### Frontend: codegen.yml

Located at `web/codegen.yml`:

```yaml
schema: http://localhost:8080/graphql
documents:
  - 'src/**/*.graphql'
  - 'src/**/*.tsx'
generates:
  src/graphql/generated.ts:
    plugins:
      - typescript
      - typescript-operations
      - typescript-react-query
```

## Schema Organization

### Current Schema Structure

```
schema/
└── schema.graphql    # Main GraphQL schema
```

### Types Defined (Template Baseline)

- **User**: Basic user account information
- **PaginationInput**: Standard pagination input
- **Queries**: `health`, `me`, `user(id: ID!)`
- **Mutations**: `_empty` (placeholder - add your mutations here)
- **Subscriptions**: `_empty` (placeholder - add real-time updates here)

> **Note:** This is the minimal starting schema. Add your domain-specific types, queries, and mutations following the schema-first workflow described below.

## Development Workflow

### Adding New Fields

1. **Update Schema**

   ```bash
   # Edit schema/schema.graphql
   type User {
     id: ID!
     email: String!
     phoneNumber: String  # NEW FIELD
   }
   ```

2. **Regenerate Code**

   ```bash
   npm run generate
   ```

3. **Implement Resolver**

   ```go
   // backend/internal/graphql/resolver/schema.resolvers.go
   func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
       // Fetch user from database
       return &model.User{
           ID: id,
           Email: "user@example.com",
           PhoneNumber: stringPtr("+1234567890"), // NEW
       }, nil
   }
   ```

4. **Use in Frontend**
   ```graphql
   query GetUser($id: ID!) {
     user(id: $id) {
       id
       email
       phoneNumber # Automatically typed
     }
   }
   ```

### Adding New Queries/Mutations

1. **Add to Schema**

   ```graphql
   type Query {
     searchUsers(query: String!): [User!]!
   }
   ```

2. **Generate**

   ```bash
   npm run generate:backend
   ```

3. **Implement**
   ```go
   func (r *queryResolver) SearchUsers(ctx context.Context, query string) ([]*model.User, error) {
       // Implementation
   }
   ```

## Testing

### GraphQL Playground Testing

Visit http://localhost:8080/graphql and try:

```graphql
# Test health check
query {
  health
}

# Test user query
query {
  me {
    id
    email
    name
    createdAt
  }
}
```

### Frontend Testing

Create test queries in `web/src/graphql/queries.graphql` and run:

```bash
cd web && npm run generate
```

Then use the generated hooks in your components.

## Best Practices

### 1. Schema Design

- Use clear, descriptive names
- Define enums for fixed sets of values
- Use input types for mutations
- Add comments/descriptions to fields

### 2. Resolver Implementation

- Keep resolvers thin, delegate to services
- Handle errors gracefully
- Use context for authentication
- Implement pagination for lists

### 3. Frontend Usage

- Define queries in `.graphql` files
- Use React Query hooks for caching
- Handle loading and error states
- Leverage TypeScript types

### 4. Code Generation

- Run `npm run generate` after schema changes
- Commit generated files to version control
- Review generated code in PRs
- Don't manually edit generated files

## Troubleshooting

### Backend Generation Fails

```bash
# Clear generated files and regenerate
cd backend
rm -rf internal/graphql
go run github.com/99designs/gqlgen generate
```

### Frontend Generation Fails

```bash
# Ensure backend is running
cd backend && go run .

# In another terminal
cd web && npm run generate
```

### Schema Conflicts

- Ensure gqlgen.yml and codegen.yml point to correct schema
- Check for syntax errors in schema.graphql
- Verify imports and package names

## Next Steps

1. **Implement Authentication**: Add auth middleware and context
2. **Connect to Neo4j**: Implement resolvers with database queries
3. **Add Validation**: Input validation and error handling
4. **Enable Subscriptions**: Real-time updates with WebSockets
5. **Add DataLoader**: Optimize N+1 query problems
6. **Setup Testing**: Unit tests for resolvers

## Resources

- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Code Generator](https://the-guild.dev/graphql/codegen)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [React Query](https://tanstack.com/query/latest)
