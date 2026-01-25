# GRGN Stack Architecture

Visual overview of the GRGN Stack architecture and how components interact.

> **Note:** For the formalized MVC platform implementation specification, domain isolation policies, and internal SDK patterns, please refer to the primary design document: **[mvc_design.md](mvc_design.md)** (Current Development Focus 🚀).

---

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         GRGN Stack                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐ │
│  │   Frontend   │◄────►│   Services   │◄────►│   Database   │ │
│  │ (Distributed)│      │ (Mod Monolith)│      │   (Neo4j)    │ │
│  │              │      │              │      │              │ │
│  └──────────────┘      └──────────────┘      └──────────────┘ │
│       │                      │                      │          │
│   TypeScript            GraphQL API            Graph Schema     │
│   React Components      MVC Pattern            Cypher Queries   │
│   TanStack Query        grgn CLI               Migrations       │
│       │                      │                      │          │
│   Vite Dev              Go 1.24+               Neo4j Fabric     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Data Flow

### Query Flow (MVC)

```
User Interaction
      ↓
VIEW: React Component (services/{domain}/{app}/view/web/)
      ↓
TanStack Query Hook
      ↓
GraphQL Query (auto-generated)
      ↓
HTTP Request → API Gateway (cmd/server/)
      ↓
CONTROLLER: GraphQL Resolver (resolver.go)
      ↓
CONTROLLER: Handler (*_handler.go)
      ↓
Core Services (via Internal SDK interfaces)
      ↓
Repository Layer (generated/)
      ↓
Neo4j Driver (services/core/shared/controller/database.go)
      ↓
Cypher Query → Neo4j Fabric
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
HTTP POST → API Gateway
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

## 📁 Project Structure

> See [mvc_design.md](mvc_design.md) Section 3 for complete file layout.

```
/
├── cmd/                        # ENTRY POINTS
│   ├── grgn/                   # GRGN CLI tool
│   ├── server/                 # HTTP server (main.go)
│   ├── migrate/                # Migration runner
│   └── worker/                 # Background job runner
│
├── pkg/                        # STANDALONE PACKAGES
│   ├── config/                 # Configuration management
│   ├── grgn/                   # Core interfaces
│   └── testing/                # Test utilities
│
├── migrations/                 # CENTRAL INFRASTRUCTURE MIGRATIONS
│   └── *.cypher / *.go         # golang-migrate files
│
├── services/                   # MODULAR MONOLITH DOMAINS
│   │
│   ├── core/                   # INFRASTRUCTURE DOMAIN
│   │   ├── shared/             # Global infra (DB, mailer, cache)
│   │   │   ├── model/          # Shared GraphQL scalars
│   │   │   ├── view/           # Base components (React), admin UI
│   │   │   └── controller/     # SDK implementations
│   │   ├── auth/               # Identity & access
│   │   ├── tenant/             # Multi-tenancy
│   │   └── directory/          # Users, Groups, ACLs
│   │
│   └── {product}/              # PRODUCT DOMAINS (e.g., twitter)
│       ├── shared/             # Product-specific utils
│       └── {app}/              # Individual apps (e.g., tweet, timeline)
│           ├── model/          # GraphQL schemas (.graphql)
│           ├── view/           # Web UI (React), CLI, Jobs
│           ├── controller/     # Business logic, resolvers
│           └── generated/      # Code generation output
│
├── go.mod                      # Go module definition
└── package.json                # Project-wide CLI scripts
```

---

## 🔌 Technology Layers

### Layer 1: Frontend (Presentation - Distributed)

```
┌─────────────────────────────────────────────┐
│  React Components (Mantine UI)              │
│  ├─ Colocated in services/{domain}/{app}/   │
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
│  GraphQL Schemas (Colocated)                │
│  ├─ services/{domain}/{app}/model/          │
│  ├─ scalars.graphql (shared)                │
│  └─ Federated single source of truth        │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  Code Generation                            │
│  ├─ grgn generate orchestrates all layers   │
│  ├─ Backend: gqlgen → Go types/resolvers    │
│  └─ Frontend: graphql-codegen → TS types    │
└─────────────────────────────────────────────┘
```

### Layer 3: Backend (MVC Pattern)

```
┌─────────────────────────────────────────────┐
│  Gin HTTP Server (cmd/server/)              │
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
│  Neo4j Driver (shared controller)           │
│  ├─ Connection pooling                      │
│  ├─ Session management                      │
│  └─ Cypher execution                        │
└─────────────────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────┐
│  Neo4j Fabric                               │
│  ├─ Nodes (entities)                        │
│  ├─ Relationships (connections)             │
│  └─ Multi-tenant isolation                  │
└─────────────────────────────────────────────┘
```

---

## 🔄 Code Generation Flow

```
┌──────────────────────┐
│  model/*.graphql     │
│  (Per-app schemas)   │
└──────────────────────┘
          │
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
│  │  web         │  │  services    │  │  neo4j       │ │
│  │ (Distributed)│  │ (Mod Monolith)│  │ (Fabric)     │ │
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
./services/*/view/web →  /app/src          (live reload)
./services           →   /app/services     (hot reload)
./pkg                →   /app/pkg
./cmd                →   /app/cmd
neo4j_data           →   /data             (persistence)
```

---

## 🔐 Authentication Flow

```
User
  │
  ├─ Login with Google
  │     │
  │     └─► OAuth Flow
  │           │
  │           └─► Core Auth Service (services/core/auth)
  │                 │
  │                 └─► JWT Token
  │                       │
  │                       └─► Store in HTTP-only cookie
  │
  └─ Subsequent Requests
        │
        └─► Cookie attached
              │
              └─► API Gateway validates JWT
                    │
                    ├─ Valid → Process request
                    │            │
                    │            └─► Access Neo4j Fabric with user context
                    │
                    └─ Invalid → Return 401
```

---

## 📊 Schema Design Workflow

```
Developer                                    grgn CLI
    │                                            │
    │  1. Edit model/*.graphql                   │
    │     (services/{domain}/{app}/model/)       │
    │                                            │
    │  2. Run grgn generate                      │
    ├────────────────────────────────────────────►│
    │                                            │
    │                            3. Read model/*.graphql
    │                            4. Generate:
    │                               - Go types
    │                               - Resolvers
    │                               - TypeScript types
    │                               - React Query hooks
    │                                            │
    │  5. Implement resolvers                    │
    │     (controller/*.go)                      │
    │                                            │
    │  6. Create migrations if needed            │
    │     (migrations/*.go)                      │
    │                                            │
    │  7. Review & commit                        │
    │ ◄────────────────────────────────────────────┤
    │                                            │
```

---

## 🧪 Testing Architecture

```
┌─────────────────────────────────────────────┐
│  Frontend Tests (Vitest)                    │
│  ├─ Component tests (Colocated)              │
│  ├─ Hook tests (services/**/view/web)       │
│  ├─ Integration tests                       │
│  └─ Mock GraphQL responses                  │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  Backend Tests (Go testing)                 │
│  ├─ Unit tests (services/**/controller)     │
│  ├─ Integration tests (resolvers)           │
│  ├─ Database tests (with test Neo4j)        │
│  └─ Table-driven tests                      │
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
      │      ├─ Run all domain tests
      │      ├─ Check architecture rules
      │      ├─ Build Docker images
      │      └─ Calculate coverage
      │
      └─► GitHub Actions (CD)
             │
             ├─ Build production images
             ├─ Push to container registry
             └─ Deploy to environment
```

---

## 📈 Scalability Considerations

### Horizontal Scaling

```
                    Load Balancer
                         │
        ┌────────────────┼────────────────┐
        ▼                ▼                ▼
    Service Node 1   Service Node 2   Service Node 3
        │                │                │
        └────────────────┼────────────────┘
                         ▼
                   Neo4j Fabric Cluster
```

### Caching Strategy

```
Frontend (Distributed)
    │
    └─► TanStack Query Cache (in-memory)
            │
            └─► HTTP Request
                    │
                    ▼
                Service Layer
                    │
                    ├─► Core Cache Service (services/core/shared)
                    │
                    └─► Neo4j Fabric
```

---

## 🎯 Key Design Principles

> See [mvc_design.md](mvc_design.md) Section 2 for detailed principles.

1. **MVC Pattern (Redefined)**
   - **Model**: Declarative GraphQL schemas colocated in each app.
   - **View**: Distributed React components, CLI tools, Background jobs.
   - **Controller**: Business logic handlers and policies.

2. **Internal SDK Pattern**
   - Product domains consume core services via interfaces.
   - Decoupled from external drivers.

3. **Domain Isolation (Configurable)**
   - Enforced by `grgn` CLI based on developer policy.
   - No illegal cross-product imports.

4. **Type Safety**
   - End-to-end type safety from Graph → Go → GraphQL → TypeScript.

5. **Configuration Locality**
   - Each app owns its configuration.

6. **Schema-First Development**
   - Single source of truth colocated with logic.

7. **Multi-Tenancy by Design**
   - Configurable isolation via property or dedicated database.

8. **CLI-Driven Development**
   - Standardized workflows via `grgn` CLI.

---

This architecture provides a solid foundation for building scalable, maintainable full-stack applications with graph database capabilities. For the complete specification, see **[mvc_design.md](mvc_design.md)**.
