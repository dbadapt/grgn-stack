# GRGN Stack Template

> **G**o + **R**eact + **G**raphQL + **N**eo4j (pronounced "Gur-gen")

![CI](https://github.com/dbadapt/grgn-stack/workflows/CI/badge.svg)
![Coverage](https://github.com/dbadapt/grgn-stack/workflows/Coverage/badge.svg)
![Docker Build](https://github.com/dbadapt/grgn-stack/workflows/Docker%20Build/badge.svg)

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)
![Node Version](https://img.shields.io/badge/Node.js-18+-339933?logo=node.js&logoColor=white)
![Neo4j Version](https://img.shields.io/badge/Neo4j-5-4581C3?logo=neo4j&logoColor=white)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=black)

A modern, production-ready full-stack template for building applications with Go, Neo4j graph database, GraphQL API, and React frontend.

---

## ⚡ Prerequisites

Before using this template, ensure you have:

| Requirement        | Version | Check Command            |
| ------------------ | ------- | ------------------------ |
| **Docker**         | Latest  | `docker --version`       |
| **Docker Compose** | v2+     | `docker compose version` |
| **Node.js**        | 18+     | `node --version`         |
| **Go**             | 1.24+   | `go version`             |

> **Note:** Go 1.24+ is required due to transitive dependency requirements (`golang.org/x/*` packages).

---

## 🎯 Using This Template

**New to this template?** Start here:

### Quick Initialize (Recommended)

**Windows:**

```powershell
.\init-template.ps1
```

**Linux/Mac:**

```bash
chmod +x init-template.sh
./init-template.sh
```

The initialization script will:

- ✅ Set up your project name and repository
- ✅ Update all configuration files
- ✅ Initialize Git repository
- ✅ Create environment files
- ✅ Clean up template artifacts

📖 **For detailed setup instructions, see [TEMPLATE-SETUP.md](TEMPLATE-SETUP.md)**

---

## Features

- 🚀 **Modern Stack**: Go + Neo4j + GraphQL + React + TypeScript
- 📊 **GraphQL API**: Type-safe API with automatic code generation (gqlgen)
- 🎨 **React Frontend**: Mantine UI components + TanStack Query
- 🗄️ **Neo4j Database**: Graph database with migration support
- 🔐 **Authentication Ready**: Multi-provider auth structure prepared (implementation required)
- 🎨 **Visual Schema Design**: Arrows.app integration for collaborative modeling
- 🐳 **Fully Containerized**: Docker Compose for all environments
- ✅ **Testing**: Comprehensive test coverage with CI/CD
- 🔄 **Multi-Environment**: Dev, staging, production configurations
- 📝 **Well Documented**: Complete guides for development workflow

> **Note on Authentication:** This template provides the _structure_ for authentication (environment variables, AuthProvider graph model, config loading) but not the actual OAuth/JWT implementation. You'll implement your chosen auth providers using the provided foundation.

## What's Included vs. What You Build

### ✅ Ready to Use

- User entity with GraphQL queries (`health`, `me`, `user`)
- Neo4j database with migration framework
- Docker Compose for dev/staging/production
- CI/CD pipelines (GitHub Actions)
- Testing infrastructure (Go + Vitest)
- Code generation (gqlgen + graphql-codegen)
- Git hooks (Husky + lint-staged)
- Visual schema design workflow (Arrows.app)

### 🔨 You Implement

- Authentication logic (OAuth, JWT - structure provided)
- Your domain entities and relationships
- Business logic in resolvers
- Additional frontend components
- Deployment secrets and infrastructure

## Tech Stack

### Backend

- **Go 1.24+** with Gin web framework
- **GraphQL** with gqlgen code generation
- **Neo4j** graph database
- **Viper** for configuration management
- Database migrations with custom migrator

### Frontend

- **React 18** + **TypeScript**
- **Mantine UI** component library
- **TanStack Query** for data fetching
- **GraphQL Code Generator** for type-safe queries
- **Vite** for fast builds and HMR
- **Vitest** for testing

### DevOps

- **Docker** & **Docker Compose**
- **GitHub Actions** for CI/CD
- **Git Hooks** (Husky) for code quality
- **Codecov** for coverage tracking

## Quick Start

> **Note:** If you're setting up from the template for the first time, use the initialization script above first!

### Prerequisites

- **Docker** & **Docker Compose**
- **Node.js** 18+
- **Go** 1.24+

### Setup (After Template Initialization)

1. **Install dependencies**

   ```bash
   npm install
   cd web && npm install && cd ..
   ```

2. **Review environment files**

   Edit `.env` and `web/.env` with your configuration:
   - Database credentials
   - API keys
   - OAuth provider IDs

3. **Start services**

   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
   ```

4. **Access application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - GraphQL Playground: http://localhost:8080/graphql
   - Neo4j Browser: http://localhost:7474

## Development

### Running Tests

```bash
# All tests
npm test

# Backend tests only
npm run test:backend

# Frontend tests only
npm run test:frontend

# With coverage
npm run coverage
```

### GraphQL Development

```bash
# Generate GraphQL code (both backend and frontend)
npm run generate

# Backend only
npm run generate:backend

# Frontend only
npm run generate:frontend
```

See [GRAPHQL.md](GRAPHQL.md) for comprehensive GraphQL guide.

## Schema Design & Development

**Visual Design + AI Code Generation Workflow:**

- 🎨 **You**: Design graph models visually in [Arrows.app](https://arrows.app)
- 🤖 **Copilot**: Generates code across all layers automatically

**Documentation:**

- [SCHEMA-QUICKREF.md](SCHEMA-QUICKREF.md) - **START HERE** - Quick reference
- [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) - Complete collaborative workflow
- [schema/graph-models/README.md](schema/graph-models/README.md) - Visual model library
- [DATABASE.md](DATABASE.md) - Neo4j graph database schema
- [GRAPHQL.md](GRAPHQL.md) - GraphQL API schema

**Quick Start:**

1. Open https://arrows.app
2. Import model from `schema/graph-models/`
3. Edit and export JSON
4. Tell Copilot: "I updated [model], implement it"

## Documentation

> **Development Focus:** 🚀 MVC Platform Implementation (See [mvc_design.md](mvc_design.md))

### Getting Started

- **[TEMPLATE-SETUP.md](TEMPLATE-SETUP.md)** - 🎯 Complete template setup guide (START HERE)
- **[USING-TEMPLATE.md](USING-TEMPLATE.md)** - Quick template usage reference
- **[QUICK-REFERENCE.md](QUICK-REFERENCE.md)** - ⚡ Command cheat sheet
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - 🏗️ System architecture & data flow
- **[mvc_design.md](mvc_design.md)** - 🚀 MVC Platform Implementation Specification (CURRENT FOCUS)

### Development Guides

- [SCHEMA-QUICKREF.md](SCHEMA-QUICKREF.md) - Schema design quick reference
- [SCHEMA-WORKFLOW.md](SCHEMA-WORKFLOW.md) - Visual design + Copilot workflow
- [DATABASE.md](DATABASE.md) - Neo4j graph database design
- [GRAPHQL.md](GRAPHQL.md) - GraphQL schema and code generation
- [CONFIG.md](CONFIG.md) - Multi-environment configuration guide

### Testing & Deployment

- [COVERAGE.md](COVERAGE.md) - Code coverage and testing
- [TESTING-CI.md](TESTING-CI.md) - Local CI testing guide
- [CI-CD.md](CI-CD.md) - CI/CD pipeline and deployment

### Contributing

- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute
- [HOOKS.md](HOOKS.md) - Git hooks with Husky

## Project Structure

> **Note:** The structure below reflects the **Target Modular Architecture** (See [mvc_design.md](mvc_design.md)).

```
├── backend/                 # Go backend
│   ├── cmd/                # CLI commands (grgn, migrate, etc.)
│   ├── internal/           # Modular Monolith Domains
│   │   ├── core/           # Infrastructure Domains (Shared, Auth, Tenant)
│   │   └── {product}/      # Product Domains (e.g., twitter, commerce)
│   ├── pkg/                # Standalone importable Go packages
│   ├── main.go             # Application entry point
│   └── gqlgen.yml          # GraphQL codegen config
├── web/                    # React frontend
│   ├── src/
│   │   ├── domains/        # Domain-specific UI components
│   │   ├── graphql/        # GraphQL queries & generated code
│   │   └── config/         # Environment configuration
│   └── vite.config.ts
├── schema/                 # Shared schema definitions
│   ├── schema.graphql      # GraphQL API schema
│   └── graph-models/       # Visual Neo4j models (Arrows.app)
├── scripts/                # Utility scripts
├── .github/workflows/      # CI/CD workflows
└── docker-compose*.yml     # Docker configurations
```

## Database Migrations

```bash
# Run migrations
cd backend
go run ./cmd/migrate

# Create new migration
# Add new file: backend/internal/database/migrations/00X_description.go
# Follow pattern in 001_initial_schema.go
```

See [DATABASE.md](DATABASE.md) for schema design guide.

## Environment Configuration

Three environments supported: **development**, **staging**, **production**

```bash
# Development (default)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Staging
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

See [CONFIG.md](CONFIG.md) for configuration guide.

## Testing & CI/CD

```bash
# Test locally before pushing
npm run test:ci

# Coverage reports
npm run coverage
npm run coverage:backend
npm run coverage:frontend
```

See [TESTING-CI.md](TESTING-CI.md) and [CI-CD.md](CI-CD.md) for details.

## Contributing

This project uses:

- **Conventional Commits** for commit messages
- **Husky** for pre-commit hooks
- **ESLint** for code linting
- **Go fmt** and **go vet** for Go code quality

See [HOOKS.md](HOOKS.md) for details.

## Architecture Highlights

### Schema-First Development

- GraphQL schema (`schema/schema.graphql`) is the single source of truth
- Automatic code generation for both backend (Go) and frontend (TypeScript)
- Visual graph modeling with Arrows.app for Neo4j schema

### Type Safety

- **Backend**: Go's strong typing + generated GraphQL types
- **Frontend**: TypeScript + generated query hooks
- **Database**: Neo4j constraints ensure data integrity

### Scalability

- Graph database excels at complex relationships
- GraphQL eliminates over/under-fetching
- Docker Compose makes multi-environment deployment simple
- Horizontal scaling ready (stateless backend)

## What's Included

### Base Features

- ✅ User authentication structure
- ✅ GraphQL API with health check
- ✅ Database migration system
- ✅ React app with Mantine UI
- ✅ Comprehensive testing setup
- ✅ CI/CD pipelines
- ✅ Multi-environment configuration
- ✅ Visual schema design workflow

### Ready to Build

This template provides the foundation. Add your business logic:

1. Design your graph model in Arrows.app
2. Let Copilot generate migrations, resolvers, and repositories
3. Add your frontend components
4. Deploy with Docker Compose

## License

[MIT](LICENSE)

## Support

- 📖 Read the [documentation](./SCHEMA-QUICKREF.md)
- 🐛 Report issues on [GitHub Issues](https://github.com/dbadapt/grgn-stack/issues)
- 💬 Discuss on [GitHub Discussions](https://github.com/dbadapt/grgn-stack/discussions)

---

**Built with ❤️ using Go, Neo4j, GraphQL, and React**
