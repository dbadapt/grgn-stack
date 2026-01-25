# GRGN Stack Architecture

Visual overview of the GRGN Stack architecture and how components interact.

> **Note:** For the formalized MVC platform implementation specification, domain isolation policies, and internal SDK patterns, please refer to the primary design document: **[backend/mvc_design.md](backend/mvc_design.md)** (Current Development Focus 🚀).

---

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         GRGN Stack                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐ │
│  │              │      │              │      │              │ │
│  │   Frontend   │◄────►│   Backend    │◄────►│   Database   │ │
│  │   (React)    │      │     (Go)     │      │   (Neo4j)    │ │
│  │              │      │              │      │              │ │
│  └──────────────┘      └──────────────┘      └──────────────┘ │
│       │                      │                      │          │
│   TypeScript            GraphQL API            Graph Schema     │
│   Mantine UI            Gin Framework          Cypher Queries   │
│   TanStack Query        gqlgen                 Migrations       │
│       │                      │                      │          │
│   Port 5173             Port 8080               Port 7687       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Data Flow

### Query Flow (MVC)

```
User Interaction
      ↓
VIEW: React Component (web/src/domains/)
      ↓
TanStack Query Hook
      ↓
GraphQL Query (auto-generated)
      ↓
HTTP Request → Backend :8080/graphql
      ↓
CONTROLLER: GraphQL Resolver (resolver.go)
      ↓
CONTROLLER: Handler (*_handler.go)
      ↓
Core Services (via Internal SDK interfaces)
      ↓
Repository Layer (generated/)
      ↓
Neo4j Driver (core/shared/controller/database.go)
      ↓
Cypher Query → Neo4j :7687
      ↓
Graph Data
      ↓
[Return path reverses]
      ↓
MODEL: GraphQL Response (types from .graphql)
      ↓
VIEW: React Component Re-render
```

### Mutation Flow (MVC)

```
User Action (Click, Submit, etc.)
      ↓
VIEW: React Event Handler
      ↓
TanStack Mutation Hook
      ↓
GraphQL Mutation (auto-generated)
      ↓
HTTP POST → Backend :8080/graphql
      ↓
CONTROLLER: GraphQL Resolver (resolver.go)
      ↓
CONTROLLER: Handler (*_handler.go)
      ↓
CONTROLLER: Policy Validation (*_policy.go)
      ↓
Core Services (auth, mailer, etc. via interfaces)
      ↓
Repository Layer (generated/)
      ↓
Neo4j Transaction
      ↓
Cypher CREATE/UPDATE/DELETE
      ↓
Commit Transaction
      ↓
MODEL: Response (GraphQL types)
      ↓
VIEW: Cache Invalidation → UI Update
```

---

## 📁 Layer Architecture

### Frontend (React)

```
web/
│
├── src/
│   ├── App.tsx              # Application root
│   │
│   ├── domains/             # Domain-specific UI (mirrors backend)
│   │   └── {product}/          # e.g., twitter/
│   │       ├── components/        # Domain components
│   │       └── pages/             # Domain pages
│   │
│   ├── components/          # Global reusable UI components
│   │   └── *.tsx               # (inherits from core/shared/view/web)
│   │
│   ├── pages/              # Global page components
│   │   └── *.tsx
│   │
│   ├── graphql/            # GraphQL queries & generated code
│   │   ├── queries.graphql     # Hand-written queries
│   │   └── generated.ts        # Auto-generated types & hooks
│   │
│   ├── config/             # Configuration
│   │   └── env.ts              # Environment variables
│   │
│   ├── hooks/              # Custom React hooks
│   │   └── *.ts
│   │
│   ├── utils/              # Utility functions
│   │   └── *.ts
│   │
│   └── test/               # Test utilities
│       └── setup.ts
```

### Backend (Go) - Modular Monolith

> See [mvc_design.md](backend/mvc_design.md) Section 3 for complete file layout.

```
backend/
│
├── main.go                     # Application entry point
│
├── cmd/                        # Command-line tools
│   ├── grgn/                   # GRGN CLI tool
│   ├── server/                 # HTTP server
│   ├── migrate/                # Migration runner
│   └── worker/                 # Background job runner
│
├── internal/                   # Modular Monolith Domains
│   │
│   ├── core/                   # INFRASTRUCTURE DOMAIN
│   │   ├── shared/             # Global infra (DB, mailer, cache)
│   │   │   ├── model/          # Shared GraphQL scalars
│   │   │   ├── view/           # Base components, admin UI
│   │   │   └── controller/     # SDK implementations
│   │   ├── auth/               # Identity & access
│   │   │   ├── model/          # CoreAuthUser, Session
│   │   │   ├── view/           # Login UI, CLI tools
│   │   │   └── controller/     # Auth handlers
│   │   ├── tenant/             # Multi-tenancy
│   │   └── directory/          # Users, Groups, ACLs
│   │
│   └── {product}/              # PRODUCT DOMAINS (e.g., twitter)
│       ├── shared/             # Product-specific utils
│       └── {app}/              # Individual apps (e.g., tweet, timeline)
│           ├── model/          # GraphQL types (.graphql)
│           ├── view/           # Web, CLI, Jobs
│           ├── controller/     # Business logic, resolvers
│           └── generated/      # Code generation output
│
└── migrations/                 # Central core migrations
    └── *.cypher / *.go         # golang-migrate files
```

### Standalone Packages (pkg/)

```
pkg/
│
├── config/                 # Configuration management
│   └── config.go              # uber-go/config + Viper
│
├── grgn/                   # Core interfaces (importable by external projects)
│   ├── auth.go                # IAuthService interface
│   ├── tenant.go              # ITenantService interface
│   ├── mailer.go              # IMailer interface
│   └── errors.go              # Standard error types
│
└── testing/                # Test utilities
    └── mocks/                 # Interface mocks
```

### Shared Schema

```
schema/
│
├── schema.graphql          # GraphQL API schema
│                              (Single source of truth)
│
└── graph-models/           # Visual graph models
    ├── *.json                 # Arrows.app exports
    └── README.md
```

---

## 🔌 Technology Layers

### Layer 1: Frontend (Presentation)

```
┌─────────────────────────────────────────────┐
│  React Components (Mantine UI)              │
│  ├─ Buttons, Forms, Tables, etc.            │
│  └─ Responsive, accessible, themeable       │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  State Management                           │
│  ├─ TanStack Query (server state)           │
│  ├─ React Hooks (local state)               │
│  └─ Auto caching & refetching               │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  GraphQL Client (auto-generated)            │
│  ├─ Type-safe queries                       │
│  ├─ Type-safe mutations                     │
│  └─ React Query integration                 │
└─────────────────────────────────────────────┘
```

### Layer 2: API (GraphQL)

```
┌─────────────────────────────────────────────┐
│  GraphQL Schema (schema.graphql)            │
│  ├─ Types, Queries, Mutations               │
│  ├─ Input types, Enums                      │
│  └─ Single source of truth                  │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  Code Generation                            │
│  ├─ Backend: gqlgen → Go types/resolvers   │
│  └─ Frontend: graphql-codegen → TS types   │
└─────────────────────────────────────────────┘
```

### Layer 3: Backend (MVC Pattern)

> See [mvc_design.md](backend/mvc_design.md) Section 5 for MVC details.

```
┌─────────────────────────────────────────────┐
│  Gin HTTP Server                            │
│  ├─ Routing                                 │
│  ├─ Middleware (auth, tenant, logging)      │
│  └─ GraphQL endpoint handler                │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  MODEL: GraphQL Schemas (.graphql)          │
│  ├─ types.graphql (entities)                │
│  ├─ enums.graphql (enumerations)            │
│  └─ inputs.graphql (mutations)              │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  CONTROLLER: Business Logic (Go)            │
│  ├─ resolver.go (GraphQL entry point)       │
│  ├─ *_handler.go (use case logic)           │
│  ├─ *_policy.go (business rules)            │
│  └─ Injects core services via interfaces    │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  VIEW: Consumers                            │
│  ├─ view/web/ (React components)            │
│  ├─ view/cli/ (Admin CLI tools)             │
│  ├─ view/jobs/ (Background workers)         │
│  └─ view/mobile/ (Mobile API handlers)      │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  Repository Layer (generated/)              │
│  ├─ Database abstraction                    │
│  ├─ Cypher query builders                   │
│  └─ Transaction management                  │
└─────────────────────────────────────────────┘
```

### Layer 4: Database (Persistence)

```
┌─────────────────────────────────────────────┐
│  Neo4j Driver                               │
│  ├─ Connection pooling                      │
│  ├─ Session management                      │
│  └─ Cypher execution                        │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  Neo4j Graph Database                       │
│  ├─ Nodes (entities)                        │
│  ├─ Relationships (connections)             │
│  └─ Properties (attributes)                 │
└─────────────────────────────────────────────┘
```

---

## 🔄 Code Generation Flow

> See [mvc_design.md](backend/mvc_design.md) Section 10 for complete generation details.

```
┌──────────────────────┐   ┌──────────────────────┐
│  Arrows.app          │   │  model/*.graphql     │
│  (Visual Design)     │   │  (Per-app schemas)   │
└──────────────────────┘   └──────────────────────┘
          │                         │
          ▼                         │
┌──────────────────────┐            │
│  graph-models/*.json │            │
│  (Export to repo)    │────────────┤
└──────────────────────┘            │
                                    ▼
                        ┌──────────────────────┐
                        │  grgn generate       │
                        │  (CLI orchestrates)  │
                        └──────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          ▼                         ▼                         ▼
┌──────────────────────┐   ┌──────────────────────┐   ┌──────────────────────┐
│  gqlgen (Backend)    │   │  Repository Gen      │   │  graphql-codegen     │
│  - models_gen.go     │   │  - interfaces.go     │   │  (Frontend)          │
│  - Resolver stubs    │   │  - neo4j_impl.go     │   │  - generated.ts      │
│  - Input types       │   │  - Type mappers      │   │  - React Query hooks │
└──────────────────────┘   └──────────────────────┘   └──────────────────────┘
          │                         │                         │
          └─────────────────────────┼─────────────────────────┘
                                    ▼
                        ┌──────────────────────┐
                        │  generated/ folder   │
                        │  (per app)           │
                        └──────────────────────┘
```

---

## 🐳 Docker Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Docker Compose                                         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  web         │  │  backend     │  │  neo4j       │ │
│  │              │  │              │  │              │ │
│  │  Node:18     │  │  golang:1.24 │  │  neo4j:5     │ │
│  │  Vite Dev    │  │  Gin Server  │  │  Database    │ │
│  │              │  │              │  │              │ │
│  │  :5173       │  │  :8080       │  │  :7687       │ │
│  │              │  │              │  │  :7474 (UI)  │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│        │                  │                  │         │
│        └──────────────────┼──────────────────┘         │
│                           │                            │
│                     [network]                          │
└─────────────────────────────────────────────────────────┘
```

### Volume Mounts

```
Host                     Container
────────────────────     ──────────────────
./web/src         →      /app/src          (live reload)
./backend         →      /app              (hot reload)
neo4j_data        →      /data             (persistence)
```

---

## 🔐 Authentication Flow (Future)

```
User
  │
  ├─ Login with Google
  │     │
  │     └─► OAuth Flow
  │           │
  │           └─► Google Auth
  │                 │
  │                 └─► JWT Token
  │                       │
  │                       └─► Store in HTTP-only cookie
  │
  └─ Subsequent Requests
        │
        └─► Cookie attached
              │
              └─► Backend validates JWT
                    │
                    ├─ Valid → Process request
                    │            │
                    │            └─► Access Neo4j with user context
                    │
                    └─ Invalid → Return 401
```

---

## 📊 Schema Design Workflow

```
Developer                    Arrows.app              Copilot/AI
    │                            │                        │
    │  1. Design visually        │                        │
    ├──────────────────────────► │                        │
    │                            │                        │
    │  2. Export JSON            │                        │
    │ ◄──────────────────────────┤                        │
    │                            │                        │
    │  3. Save to repo           │                        │
    │  (schema/graph-models/)    │                        │
    │                            │                        │
    │  4. Tell Copilot           │                        │
    ├─────────────────────────────────────────────────────►│
    │  "Implement model X"       │                        │
    │                            │                        │
    │                            │  5. Read JSON          │
    │                            │  6. Generate:          │
    │                            │     - Migrations       │
    │                            │     - GraphQL schema   │
    │                            │     - Resolvers        │
    │                            │     - Repositories     │
    │                            │                        │
    │  7. Review & commit        │                        │
    │ ◄─────────────────────────────────────────────────────┤
    │                            │                        │
```

---

## 🧪 Testing Architecture

```
┌─────────────────────────────────────────────┐
│  Frontend Tests (Vitest)                    │
│  ├─ Component tests                         │
│  ├─ Hook tests                              │
│  ├─ Integration tests                       │
│  └─ Mock GraphQL responses                  │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  Backend Tests (Go testing)                 │
│  ├─ Unit tests (repositories)               │
│  ├─ Integration tests (resolvers)           │
│  ├─ Database tests (with test Neo4j)        │
│  └─ Table-driven tests                      │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  E2E Tests (Future)                         │
│  ├─ Full user flows                         │
│  ├─ Real database                           │
│  └─ Browser automation                      │
└─────────────────────────────────────────────┘
```

---

## 🚀 Deployment Architecture (CI/CD)

```
Developer Push
      │
      ▼
GitHub Repository
      │
      ├─► GitHub Actions (CI)
      │      │
      │      ├─ Run backend tests
      │      ├─ Run frontend tests
      │      ├─ Check linting
      │      ├─ Build Docker images
      │      └─ Calculate coverage
      │
      └─► GitHub Actions (CD)
             │
             ├─ Build production images
             ├─ Push to container registry
             └─ Deploy to environment
                   │
                   ├─ Development
                   ├─ Staging
                   └─ Production
```

---

## 📈 Scalability Considerations

### Horizontal Scaling

```
                    Load Balancer
                         │
        ┌────────────────┼────────────────┐
        ▼                ▼                ▼
    Backend 1        Backend 2        Backend 3
        │                │                │
        └────────────────┼────────────────┘
                         ▼
                   Neo4j Cluster
                 (Causal Cluster)
```

### Caching Strategy

```
Frontend
    │
    └─► TanStack Query Cache (in-memory)
            │
            └─► HTTP Request
                    │
                    ▼
                Backend
                    │
                    ├─► Redis Cache (future)
                    │
                    └─► Neo4j Database
```

---

## 🎯 Key Design Principles

> See [mvc_design.md](backend/mvc_design.md) Section 2 for detailed principles.

1. **MVC Pattern (Redefined)**
   - **Model**: Declarative GraphQL schemas (.graphql files)
   - **View**: Web, CLI, Jobs, Mobile (not just HTML)
   - **Controller**: Business logic, resolvers, policies

2. **Internal SDK Pattern**
   - Product domains consume core services via interfaces
   - Never call external drivers directly
   - Single point of change for infrastructure swaps

3. **Domain Isolation (Configurable)**
   - `strict` / `relaxed` / `open` / `custom` policies
   - Developer-defined in `service_config.yaml`
   - Validated by `grgn` CLI

4. **Type Safety**
   - TypeScript on frontend
   - Go on backend
   - GraphQL schema as contract
   - Domain-prefixed types prevent collisions

5. **Configuration Locality**
   - Each app has `service_config.yaml`
   - Hierarchical inheritance (uber-go/config)
   - No giant global config file

6. **Schema-First Development**
   - GraphQL schemas are single source of truth
   - Code generation for types, resolvers, repositories
   - Visual design with Arrows.app

7. **Multi-Tenancy by Design**
   - Configurable isolation (property vs database)
   - Neo4j Fabric for cross-database queries
   - Tenant context middleware

8. **CLI-Driven Development**
   - `grgn` CLI for scaffolding, validation, deployment
   - `grgn make:*` for code generation
   - `grgn migrate` for schema management (golang-migrate)

---

This architecture provides a solid foundation for building scalable, maintainable full-stack applications with graph database capabilities. For the complete specification, see **[backend/mvc_design.md](backend/mvc_design.md)**.
