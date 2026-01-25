# GraphQL Setup Guide

## Overview

GRGN Stack uses **GraphQL** as its primary API layer with automatic code generation for both backend (Go) and frontend (TypeScript/React).

## Architecture

- **Backend**: gqlgen generates Go resolvers from GraphQL schemas
- **Frontend**: graphql-codegen generates TypeScript types and React Query hooks
- **Schema**: Single source of truth colocated in `services/{domain}/{app}/model/`
- **Playground**: Interactive GraphQL IDE available in development

## Quick Start

### 1. Schema Development

Edit the GraphQL schema in your app's model directory:

```bash
services/{domain}/{app}/model/types.graphql
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
services/{domain}/{app}/controller/resolver.go
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

Create `.graphql` files in `services/{domain}/{app}/view/web/`:

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

Located in each app directory:

```yaml
schema:
  - model/*.graphql
  - ../../core/shared/model/*.graphql
model:
  filename: controller/generated/models_gen.go
  package: generated
resolver:
  layout: follow-schema
  dir: controller
  package: controller
exec:
  filename: controller/generated/generated.go
  package: generated
```

### Frontend: codegen.yml

```yaml
schema: http://localhost:8080/graphql
documents:
  - 'services/**/view/web/**/*.graphql'
  - 'services/**/view/web/**/*.tsx'
generates:
  services/core/shared/view/web/generated.ts:
    plugins:
      - typescript
      - typescript-operations
      - typescript-react-query
```

## Schema Organization

### Current Schema Structure

```
services/
├── core/
│   ├── shared/model/        # Global scalars & common types
│   └── auth/model/          # Identity & Access types
└── twitter/
    └── tweet/model/         # Tweet specific types
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
   # Edit services/core/auth/model/types.graphql
   type CoreAuthUser {
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
   // services/core/auth/controller/resolver.go
   func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
       // Fetch user from database
       return &model.User{
           ID: "123",
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

Create test queries in `services/**/view/web/*.graphql` and run:

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
rm -rf services/**/controller/generated
grgn generate:backend
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
