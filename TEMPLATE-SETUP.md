# GRGN Stack Template Setup Guide

Welcome to the GRGN Stack template! This guide will help you initialize your new project.

## What is GRGN Stack?

**GRGN** (pronounced "Gur-gen") stands for:

- **G**o - Backend server
- **R**eact - Frontend UI
- **G**raphQL - API layer
- **N**eo4j - Graph database

A production-ready, fully containerized stack with comprehensive testing, CI/CD, and modern development workflows.

---

## Quick Start

### Option 1: Automated Setup (Recommended)

#### Windows (PowerShell)

```powershell
.\init-template.ps1
```

#### Linux/Mac (Bash)

```bash
chmod +x init-template.sh
./init-template.sh
```

The script will:

1. ✅ Prompt you for project details
2. ✅ Replace all template placeholders
3. ✅ Update Go module paths
4. ✅ Configure environment variables
5. ✅ Initialize Git repository
6. ✅ Create .env files
7. ✅ Clean up template-specific files

### Option 2: Manual Setup

If you prefer to set up manually, follow these steps:

#### 1. Update Project Name

Replace `yourusername` and `grgn-stack` in these files:

- `backend/go.mod` - module path
- `pkg/go.mod` - module path
- `go.work` - workspace configuration
- `.env.example` - environment variable prefix
- `pkg/config/config.go` - env prefix and default app name
- All docker-compose files - environment variable prefix
- `backend/main.go` - GraphQL playground title

#### 2. Update GitHub Repository

Replace `github.com/yourusername/grgn-stack` with your actual repository path in:

- `backend/go.mod`
- `pkg/go.mod`
- `README.md` (badges and clone URL)
- All documentation files

#### 3. Update Documentation

Search and replace in all `.md` files:

- `YOUR_USERNAME` → your GitHub username
- `YOUR_REPO` → your repository name

#### 4. Initialize Git

```bash
git init
git add .
git commit -m "Initial commit from GRGN stack template"
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
```

#### 5. Create Environment Files

```bash
cp .env.example .env
cp web/.env.example web/.env
```

Edit both `.env` files with your configuration.

---

## Post-Initialization Setup

After running the initialization script:

### 1. Install Dependencies

```bash
# Root dependencies (Husky for git hooks)
npm install

# Frontend dependencies
cd web
npm install
cd ..
```

### 2. Configure Environment

Edit `.env` and `web/.env` with your settings:

**Backend (.env)**

- Database credentials (Neo4j)
- JWT secrets
- OAuth provider keys (Google, Apple, etc.)
- Server port and host

**Frontend (web/.env)**

- API URL (GraphQL endpoint)
- OAuth client IDs
- Feature flags

### 3. Start Development Environment

```bash
# Start all services with Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

This starts:

- **Backend** (Go/Gin) on port 8080
- **Frontend** (React/Vite) on port 5173
- **Neo4j** database on port 7687 (browser on 7474)

### 4. Verify Setup

Open your browser:

- Frontend: http://localhost:5173
- GraphQL Playground: http://localhost:8080/graphql
- Neo4j Browser: http://localhost:7474

### 5. Run Tests

```bash
# All tests (backend + frontend)
npm test

# Backend only
npm run test:go

# Frontend only
npm run test:web

# With coverage
npm run coverage
```

---

## Project Structure

```
.
├── backend/              # Go backend (Gin + GraphQL)
│   ├── cmd/             # Command-line tools (migrations)
│   ├── internal/        # Internal packages
│   │   ├── database/    # Neo4j connection & migrations
│   │   ├── graphql/     # GraphQL resolvers & generated code
│   │   └── repository/  # Data access layer
│   ├── main.go          # Application entry point
│   └── go.mod
├── web/                 # React frontend
│   ├── src/
│   │   ├── config/      # Configuration
│   │   ├── graphql/     # GraphQL queries
│   │   └── test/        # Test utilities
│   ├── vite.config.ts
│   └── package.json
├── pkg/                 # Shared Go packages (config)
├── schema/              # GraphQL schema & graph models
├── scripts/             # Utility scripts
├── .github/workflows/   # CI/CD pipelines
└── docker-compose*.yml  # Container orchestration
```

---

## Environment Configuration

The stack supports three environments:

### Development

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

### Staging

```bash
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up
```

### Production

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

Each has separate configuration in the respective compose file.

---

## Schema Design Workflow

This stack includes a visual schema design workflow using [Arrows.app](https://arrows.app):

1. **Design visually** in Arrows.app
2. **Export JSON** to `schema/graph-models/`
3. **Tell Copilot** to generate code from your design
4. **Copilot generates**:
   - Neo4j migrations
   - GraphQL schema
   - Resolvers
   - Repository methods

See [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) for details.

---

## Code Generation

The stack uses automatic code generation:

### Backend (gqlgen)

Generates Go types and resolvers from GraphQL schema:

```bash
npm run generate:backend
# or
cd backend && go run github.com/99designs/gqlgen generate
```

### Frontend (GraphQL Code Generator)

Generates TypeScript types and React Query hooks:

```bash
npm run generate:frontend
# or
cd web && npm run generate
```

Run both:

```bash
npm run generate
```

---

## Testing

### Backend Tests

```bash
cd backend
go test -v ./...
```

### Frontend Tests

```bash
cd web
npm test
```

### Coverage Reports

```bash
# Backend
npm run coverage:go
# Opens coverage.html in backend/

# Frontend
npm run coverage:web
# Opens coverage in web/coverage/
```

---

## CI/CD

GitHub Actions workflows are pre-configured:

- **CI**: Runs tests on every push/PR
- **Coverage**: Tracks code coverage
- **Docker Build**: Builds container images
- **Deploy**: Deployment workflows

See [CI-CD.md](CI-CD.md) for setup instructions.

---

## Git Hooks

Husky is configured for pre-commit hooks:

- **Linting**: ESLint on TypeScript files
- **Formatting**: gofmt on Go files

Install hooks:

```bash
npm install  # Automatically runs 'husky' script
```

---

## Documentation

Comprehensive guides are included:

- [README.md](README.md) - Project overview
- [USING-TEMPLATE.md](USING-TEMPLATE.md) - Quick template usage
- [CONFIG.md](CONFIG.md) - Multi-environment configuration
- [DATABASE.md](DATABASE.md) - Neo4j schema & migrations
- [GRAPHQL.md](GRAPHQL.md) - GraphQL schema & code generation
- [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) - Visual schema design
- [TESTING-CI.md](TESTING-CI.md) - Testing strategy
- [COVERAGE.md](COVERAGE.md) - Coverage tracking
- [CI-CD.md](CI-CD.md) - Deployment pipelines
- [HOOKS.md](HOOKS.md) - Git hooks setup

---

## Technology Stack

### Backend

- **Go 1.22+** - Programming language
- **Gin** - Web framework
- **gqlgen** - GraphQL server & code generation
- **Neo4j Go Driver** - Database client
- **Viper** - Configuration management

### Frontend

- **React 18** - UI library
- **TypeScript** - Type safety
- **Mantine UI** - Component library
- **TanStack Query** - Data fetching
- **Vite** - Build tool
- **Vitest** - Testing framework

### Database

- **Neo4j 5** - Graph database

### DevOps

- **Docker** - Containerization
- **Docker Compose** - Orchestration
- **GitHub Actions** - CI/CD
- **Husky** - Git hooks

---

## Troubleshooting

### Go Module Issues

```bash
cd backend && go mod tidy
cd ../pkg && go mod tidy
```

### Docker Issues

```bash
# Clean up and restart
docker-compose down -v
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

### Port Conflicts

If ports are already in use, update in:

- `docker-compose.yml`
- `.env` (backend port)
- `web/.env` (API URL)

### Neo4j Connection Issues

- Check Neo4j is running: `docker ps`
- Verify credentials in `.env`
- Check Neo4j logs: `docker-compose logs neo4j`

---

## Next Steps

1. **Customize Schema**: Design your graph model in [Arrows.app](https://arrows.app)
2. **Add Authentication**: Implement OAuth providers (Google, Apple, etc.)
3. **Build Features**: Add your business logic
4. **Configure CI/CD**: Set up GitHub secrets for deployment
5. **Deploy**: Use provided Docker images for deployment

---

## Getting Help

- Review the documentation in the project root
- Check [Neo4j documentation](https://neo4j.com/docs/)
- Read [gqlgen docs](https://gqlgen.com/)
- Explore [React Query docs](https://tanstack.com/query/latest)

---

**Happy building with GRGN Stack! 🚀**
