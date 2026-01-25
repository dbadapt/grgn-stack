# GRGN Stack Quick Reference

> Quick commands and tips for working with the GRGN stack

## 🚀 Quick Start

```bash
# Initialize template (first time only)
.\init-template.ps1              # Windows
./init-template.sh               # Linux/Mac

# Install all dependencies
npm run install:all

# Start development environment (Docker)
npm run dev                      # or npm start

# Start services individually
npm run dev:backend              # Backend only
npm run dev:frontend             # Frontend only
```

## 🔗 Access Points

| Service            | URL                           |
| ------------------ | ----------------------------- |
| Frontend           | http://localhost:5173         |
| Backend API        | http://localhost:8080         |
| GraphQL Playground | http://localhost:8080/graphql |
| Neo4j Browser      | http://localhost:7474         |

**Neo4j Credentials:** `neo4j` / (check `.env` for password)

## 📦 Common Commands

### Development

```bash
npm start                        # Start all services (Docker)
npm run dev:backend             # Run backend locally
npm run dev:frontend            # Run frontend locally
```

### Building

```bash
npm run build                   # Build everything
npm run build:backend          # Build Go binary to bin/app
npm run build:frontend         # Build React for production
```

### Testing

```bash
npm test                        # All tests (backend + frontend)
npm run test:backend           # Backend tests only
npm run test:frontend          # Frontend tests only
npm run test:watch             # Watch mode (frontend)
npm run coverage               # Coverage reports
```

### Linting & Formatting

```bash
npm run lint                    # Lint all code
npm run lint:fix               # Fix linting issues
npm run format                 # Format all code
```

### Code Generation

```bash
npm run generate               # Generate all (GraphQL backend + frontend)
npm run generate:backend      # GraphQL Go code (gqlgen)
npm run generate:frontend     # GraphQL TypeScript code
```

### Database

```bash
npm run db:migrate             # Run migrations
npm run db:reset              # Reset database (docker down + up)
```

### Docker

```bash
npm run docker:up              # Start all services
npm run docker:up:detached    # Start detached
npm run docker:down           # Stop services
npm run docker:down:volumes   # Stop and remove volumes
npm run docker:logs           # View all logs
npm run docker:logs:backend   # Backend logs only
npm run docker:logs:frontend  # Frontend logs only
npm run docker:logs:db        # Neo4j logs only
npm run docker:build          # Rebuild containers
npm run docker:restart        # Restart services
npm run docker:clean          # Clean everything
```

### Dependencies

```bash
npm run install:all            # Install all dependencies
npm run update:all            # Update all dependencies
npm run clean                 # Clean deps and build artifacts
```

### CI/CD

```bash
npm run ci                     # Run CI pipeline (lint + test + build)
npm run validate              # Quick validation (lint + test)
```

## 🗄️ Database

### Migrations

```bash
# Run migrations
npm run db:migrate

# Or via script
.\scripts\migrate.ps1          # Windows
./scripts/migrate.sh           # Linux/Mac
```

### Neo4j Cypher Queries

```cypher
// List all nodes
MATCH (n) RETURN n LIMIT 25;

// Count nodes by label
MATCH (n) RETURN labels(n), count(*);

// Delete all (development only!)
MATCH (n) DETACH DELETE n;
```

## 🐳 Docker

```bash
# Build containers
docker-compose build

# Start with logs
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Start detached
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# View logs
docker-compose logs -f backend
docker-compose logs -f web
docker-compose logs -f neo4j

# Stop and remove volumes
docker-compose down -v

# Rebuild and start
docker-compose up --build
```

## 📝 Git Workflow

```bash
# Create feature branch
git checkout -b feature/my-feature

# Commit changes (triggers hooks)
git add .
git commit -m "feat(scope): description"

# Push branch
git push origin feature/my-feature

# Create PR on GitHub
```

### Commit Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Code style
- `refactor` - Code refactoring
- `test` - Tests
- `chore` - Maintenance

## 🔧 Environment Variables

### Backend (.env)

```bash
PROJECT_NAME_SERVER_PORT=8080
PROJECT_NAME_DATABASE_NEO4J_URI=bolt://localhost:7687
PROJECT_NAME_AUTH_JWT_SECRET=your-secret
```

### Frontend (web/.env)

```bash
VITE_API_URL=http://localhost:8080
VITE_API_GRAPHQL_URL=http://localhost:8080/graphql
VITE_ENVIRONMENT=development
```

> Replace `PROJECT_NAME` with your actual project name in uppercase

## 🎨 Schema Design

1. **Design in Arrows.app**
   - Visit https://arrows.app
   - Create nodes and relationships
   - Export as JSON

2. **Save model**

   ```bash
   # Save to:
   services/{domain}/{app}/model/your-model.json
   ```

3. **Generate code**
   - Ask Copilot to read the model
   - Update GraphQL schema
   - Generate migrations
   - Create resolvers

See [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) for details.

## 🧪 Testing Tips

### Backend

```go
// Table-driven tests
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
    }{
        {"case 1", "input", "output"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MyFunction(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Frontend

```typescript
// Component test
describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });
});
```

## 📊 GraphQL

### Query Example

```graphql
query GetUser($id: ID!) {
  user(id: $id) {
    id
    email
    name
  }
}
```

### Mutation Example

```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    email
  }
}
```

### Test in Playground

1. Go to http://localhost:8080/graphql
2. Paste query/mutation
3. Add variables (JSON)
4. Click play button

## 🔍 Troubleshooting

### Port Already in Use

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill
```

### Module Not Found (Go)

```bash
cd backend && go mod tidy
cd ../internal && go mod tidy
```

### Docker Issues

```bash
# Full reset
docker-compose down -v
docker system prune -a
docker-compose up --build
```

### Neo4j Connection Failed

```bash
# Check Neo4j is running
docker-compose ps

# View Neo4j logs
docker-compose logs neo4j

# Restart Neo4j
docker-compose restart neo4j
```

## 📚 Documentation

| File                                     | Purpose                |
| ---------------------------------------- | ---------------------- |
| [README.md](README.md)                   | Main documentation     |
| [TEMPLATE-SETUP.md](TEMPLATE-SETUP.md)   | Template setup guide   |
| [CONFIG.md](CONFIG.md)                   | Configuration guide    |
| [DATABASE.md](DATABASE.md)               | Database schema        |
| [GRAPHQL.md](GRAPHQL.md)                 | GraphQL guide          |
| [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) | Visual design workflow |
| [TESTING-CI.md](TESTING-CI.md)           | Testing & CI/CD        |
| [COVERAGE.md](COVERAGE.md)               | Coverage tracking      |
| [CI-CD.md](CI-CD.md)                     | Deployment guide       |
| [CONTRIBUTING.md](CONTRIBUTING.md)       | How to contribute      |

## 🎯 Common Tasks

### Add New GraphQL Type

1. Edit `services/{domain}/{app}/model/types.graphql`
2. Run `grgn generate:backend`
3. Implement resolver in `services/{domain}/{app}/controller/resolver.go`
4. Run `grgn generate:frontend`
5. Use in React components (distributed in services/)

### Add Database Migration

1. Create file in `services/{domain}/{app}/migrations/` (or root `migrations/`)
2. Implement migration
3. Run `grgn migrate`

### Add New Route

1. Add to GraphQL schema
2. Generate code
3. Implement resolver
4. Add repository method if needed
5. Write tests

---

**Need more help?** Check the full documentation or open an issue!
