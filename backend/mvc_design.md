# GRGN Stack Architecture Design

> **Version:** 2.0 (2026-01-25)
> **Status:** Design Specification
> **Stack:** Go + React + GraphQL + Neo4j

This document defines the architecture for a Multi-Domain Modular Monolith using the GRGN Stack. It establishes conventions for domain isolation, type-safe internal SDKs, and scalable file organization.

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Core Principles](#2-core-principles)
3. [File Layout Specification](#3-file-layout-specification)
4. [Domain Prefixing & Naming Conventions](#4-domain-prefixing--naming-conventions)
5. [MVC Redefined for 2026](#5-mvc-redefined-for-2026)
6. [Service Configuration System](#6-service-configuration-system)
7. [Internal SDK Pattern](#7-internal-sdk-pattern)
8. [Multi-Tenancy Architecture](#8-multi-tenancy-architecture)
9. [CLI Tool Specification (grgn)](#9-cli-tool-specification-grgn)
10. [Code Generation](#10-code-generation)
11. [Database Schema Management](#11-database-schema-management)
12. [Architecture Enforcement (Configurable)](#12-architecture-enforcement-configurable)
13. [Migration Path](#13-migration-path)
- [Appendix A: Quick Reference](#appendix-a-quick-reference)
- [Appendix B: Example Timeline Implementation](#appendix-b-example-timeline-implementation)
- [Appendix C: Architecture Policy Presets](#appendix-c-architecture-policy-presets)

---

## 1. System Overview

### 1.1 What is GRGN Stack?

GRGN (pronounced "Gur-gen") is a modular monolith architecture that treats **infrastructure (core)** and **product logic (domains)** as peers, while maintaining a strict hierarchy where product domains consume core services through well-defined interfaces.

### 1.2 Key Characteristics

- **Type-Safe Internal SDKs**: Domains don't call external services directly; they consume core controllers
- **Declarative Schema-First**: GraphQL schemas are the single source of truth
- **Graph-Native Data**: Neo4j enables natural relationship modeling
- **Multi-Tenant by Design**: Configurable isolation from property-level to database-level
- **CLI-Driven Development**: The `grgn` CLI automates scaffolding, validation, and deployment

### 1.3 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         GRGN Stack                                  │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────────────────────────┐    │
│  │   web/ (React)  │◄───│  internal/core/shared/view/         │    │
│  │   - Components  │    │  - Base components                  │    │
│  │   - Pages       │    │  - Theme/Design system              │    │
│  └────────┬────────┘    └─────────────────────────────────────┘    │
│           │ GraphQL                                                 │
│           ▼                                                         │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    GraphQL Gateway                           │   │
│  │           (Federated schemas from all domains)               │   │
│  └──────────────────────────┬──────────────────────────────────┘   │
│                             │                                       │
│  ┌──────────────────────────┴──────────────────────────────────┐   │
│  │                      internal/                               │   │
│  │  ┌─────────────────┐         ┌─────────────────────────┐    │   │
│  │  │      core/      │◄────────│      {product}/         │    │   │
│  │  │  ├─ shared/     │         │  ├─ shared/             │    │   │
│  │  │  ├─ auth/       │         │  ├─ {app1}/             │    │   │
│  │  │  ├─ tenant/     │         │  └─ {app2}/             │    │   │
│  │  │  └─ directory/  │         │                         │    │   │
│  │  └─────────────────┘         └─────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                             │                                       │
│  ┌──────────────────────────┴──────────────────────────────────┐   │
│  │                      pkg/ (Standalone)                       │   │
│  │           Importable by external Go projects                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                             │                                       │
│  ┌──────────────────────────┴──────────────────────────────────┐   │
│  │                    Neo4j (Fabric)                            │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │   │
│  │  │  core   │  │tenant_a │  │tenant_b │  │   shared        │ │   │
│  │  │   db    │  │   db    │  │   db    │  │   (fabric)      │ │   │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Core Principles

### 2.1 The Internal SDK Principle

Domains like `twitter` **never interact with raw external drivers** (Postmark, S3, Stripe). Instead, they consume controller logic inside `core/shared`. This ensures:

- **Single point of change**: Swapping Postmark → SendGrid requires changes in exactly one location
- **Consistent error handling**: All external service errors are wrapped consistently
- **Testability**: Mock the core interface, not external services
- **Audit logging**: All external calls flow through instrumented core controllers

### 2.2 Schema-First Development

GraphQL schemas (`.graphql` files) are the **single source of truth**:

1. Design entities in Arrows.app → Export JSON
2. Define GraphQL types from the visual model
3. Generate Go types, resolvers, and repository stubs
4. Implement business logic in controllers

### 2.3 Domain Isolation (Configurable)

Domain isolation policies are **defined by developers** in `service_config.yaml`. The stack provides sensible defaults but does not mandate a specific isolation strategy.

**Available Policies:**
- `strict` - No cross-domain imports allowed
- `relaxed` - Cross-domain imports allowed with explicit declarations
- `open` - No import restrictions (validation skipped)
- `custom` - Developer-defined import rules

The `grgn validate` command checks against **your chosen policy**, not a hardcoded ruleset.

### 2.4 Configuration Locality

Every app has its own `service_config.yaml` for:
- Feature flags
- Rate limits
- Validation rules
- Environment-specific overrides

No giant global config file.

---

## 3. File Layout Specification

> **CRITICAL**: This structure is mandatory. The `grgn` CLI validates conformance.

```
/backend/
├── main.go                          # Application entry point
├── internal/
│   ├── core/                        # INFRASTRUCTURE DOMAIN
│   │   ├── service_config.yaml      # Core domain configuration
│   │   │
│   │   ├── shared/                  # GLOBAL INFRASTRUCTURE
│   │   │   ├── service_config.yaml  # DB pools, AWS regions, etc.
│   │   │   ├── model/
│   │   │   │   ├── scalars.graphql  # DateTime, JSON, Email, UUID
│   │   │   │   └── common.graphql   # PageInfo, Error, interfaces
│   │   │   ├── view/
│   │   │   │   ├── web/             # Base React components
│   │   │   │   │   ├── theme/       # Design tokens, CSS variables
│   │   │   │   │   └── components/  # Button, Input, Modal, etc.
│   │   │   │   └── admin/           # Admin dashboard UI
│   │   │   ├── controller/
│   │   │   │   ├── aws_sdk.go       # S3, SQS, SES wrapper
│   │   │   │   ├── database.go      # Neo4j driver abstraction
│   │   │   │   ├── mailer.go        # Email service interface
│   │   │   │   ├── cache.go         # Redis/Memcached wrapper
│   │   │   │   └── interfaces.go    # Exported interfaces
│   │   │   └── generated/           # Code generation output
│   │   │       ├── graphql/         # gqlgen output
│   │   │       └── mappers/         # Type converters
│   │   │
│   │   ├── auth/                    # IDENTITY & ACCESS
│   │   │   ├── service_config.yaml  # Token expiry, MFA settings
│   │   │   ├── model/
│   │   │   │   ├── types.graphql    # CoreAuthUser, Session
│   │   │   │   ├── enums.graphql    # AUTH_PROVIDER, MFA_TYPE
│   │   │   │   └── inputs.graphql   # LoginInput, RegisterInput
│   │   │   ├── view/
│   │   │   │   ├── web/             # Login/Register UI
│   │   │   │   └── cli/             # Auth verification tools
│   │   │   ├── controller/
│   │   │   │   ├── login_handler.go # Strategy pattern for login
│   │   │   │   ├── mfa_logic.go     # TOTP/SMS verification
│   │   │   │   ├── session.go       # JWT/Session management
│   │   │   │   └── resolver.go      # GraphQL resolvers
│   │   │   └── generated/
│   │   │
│   │   ├── tenant/                  # MULTI-TENANCY
│   │   │   ├── service_config.yaml  # Isolation policies, tiers
│   │   │   ├── model/
│   │   │   │   ├── instance.graphql # CoreTenantInstance
│   │   │   │   └── subscription.graphql
│   │   │   ├── view/
│   │   │   │   ├── cli/             # Tenant provisioning
│   │   │   │   └── jobs/            # Cleanup, billing
│   │   │   ├── controller/
│   │   │   │   ├── provisioner.go   # DB isolation logic
│   │   │   │   ├── quota_enforcer.go
│   │   │   │   └── resolver.go
│   │   │   └── generated/
│   │   │
│   │   ├── directory/               # LDAP-STYLE SERVICES
│   │   │   ├── service_config.yaml
│   │   │   ├── model/
│   │   │   │   ├── types.graphql    # CoreDirectoryUser, Group, ACL
│   │   │   │   └── roles.graphql    # Role, Permission
│   │   │   ├── view/
│   │   │   ├── controller/
│   │   │   │   ├── user_service.go
│   │   │   │   ├── group_service.go
│   │   │   │   ├── acl_service.go
│   │   │   │   └── resolver.go
│   │   │   └── generated/
│   │   │
│   │   └── feature/                 # FEATURE FLAGS
│   │       ├── service_config.yaml
│   │       ├── model/
│   │       │   └── types.graphql    # CoreFeatureFlag
│   │       ├── controller/
│   │       │   ├── evaluator.go     # Flag evaluation logic
│   │       │   └── resolver.go
│   │       └── generated/
│   │
│   └── twitter/                     # PRODUCT DOMAIN (Example)
│       ├── service_config.yaml      # Twitter-wide settings
│       │
│       ├── shared/                  # TWITTER-SPECIFIC UTILS
│       │   ├── service_config.yaml  # Media CDN, hashtag rules
│       │   ├── model/
│       │   │   └── common.graphql   # Shared Twitter types
│       │   ├── view/
│       │   │   └── web/             # Twitter component library
│       │   └── controller/
│       │       ├── media_processor.go
│       │       └── hashtag_aggregator.go
│       │
│       ├── tweet/                   # ATOMIC UNIT
│       │   ├── service_config.yaml  # Char limits, media counts
│       │   ├── model/
│       │   │   ├── types.graphql    # TwitterTweet
│       │   │   ├── enums.graphql    # TWEET_STATUS
│       │   │   └── interactions.graphql
│       │   ├── view/
│       │   │   ├── web/             # Tweet components
│       │   │   └── jobs/            # Sentiment analysis
│       │   ├── controller/
│       │   │   ├── post_handler.go  # Injects CoreAuth
│       │   │   ├── deletion_policy.go
│       │   │   └── resolver.go
│       │   └── generated/
│       │
│       └── timeline/                # AGGREGATOR
│           ├── service_config.yaml  # Algorithm weights, TTLs
│           ├── model/
│           │   └── types.graphql    # TwitterTimeline
│           ├── view/
│           │   ├── mobile/          # Mobile API handlers
│           │   └── jobs/            # Cache warming
│           ├── controller/
│           │   ├── resolver.go      # Top-level resolver
│           │   ├── feed_algorithm.go
│           │   └── cache_manager.go
│           └── generated/
│
├── pkg/                             # STANDALONE PACKAGES
│   ├── config/                      # Configuration loader
│   │   └── config.go
│   ├── grgn/                        # Core interfaces (importable)
│   │   ├── auth.go                  # Auth interface
│   │   ├── tenant.go                # Tenant interface
│   │   ├── mailer.go                # Mailer interface
│   │   └── errors.go                # Standard errors
│   └── testing/                     # Test utilities
│       └── mocks/
│
└── cmd/
    ├── server/                      # HTTP server
    │   └── main.go
    ├── migrate/                     # Migration runner
    │   └── main.go
    └── worker/                      # Background job runner
        └── main.go

/web/                                # REACT APPLICATION
├── src/
│   ├── App.tsx
│   ├── domains/                     # Domain-specific UI
│   │   └── twitter/
│   │       ├── components/
│   │       └── pages/
│   └── graphql/
│       ├── generated.ts             # Codegen output
│       └── queries.graphql

/schema/                             # SCHEMA DEFINITIONS
├── schema.graphql                   # Combined GraphQL schema
└── graph-models/                    # Arrows.app JSON exports
    ├── core-model.json
    └── auth-model.json
```

---

## 4. Domain Prefixing & Naming Conventions

### 4.1 The Naming Problem

When multiple teams work on different domains, identifier collisions occur:
- Two `User` types in GraphQL
- Two `UserRepository` in Go
- Ambiguous imports

### 4.2 Domain Prefix Rules

Every domain declares a unique prefix in its `service_config.yaml`:

```yaml
# internal/core/service_config.yaml
domain:
  name: core
  prefix: Core
  
# internal/twitter/service_config.yaml  
domain:
  name: twitter
  prefix: Twitter
```

### 4.3 Naming Convention Matrix

| Artifact Type | Pattern | Core Example | Twitter Example |
|--------------|---------|--------------|-----------------|
| GraphQL Type | `{Prefix}{App}{Type}` | `CoreAuthUser` | `TwitterTweetPost` |
| GraphQL Enum | `{PREFIX}_{APP}_{ENUM}` | `CORE_AUTH_PROVIDER` | `TWITTER_TWEET_STATUS` |
| GraphQL Input | `{Prefix}{App}{Action}Input` | `CoreAuthLoginInput` | `TwitterTweetCreateInput` |
| Go Package | `{domain}/{app}` | `core/auth` | `twitter/tweet` |
| Go Type | `{Type}` (package provides namespace) | `auth.User` | `tweet.Post` |
| Go Interface | `I{Type}` | `auth.IUserService` | `tweet.IPostService` |
| Repository | `{Type}Repository` | `auth.UserRepository` | `tweet.PostRepository` |
| Controller | `{Type}Controller` | `auth.LoginController` | `tweet.PostController` |
| Database Label | `{Prefix}{App}{Type}` | `CoreAuthUser` | `TwitterTweetPost` |

### 4.4 CLI Validation

The `grgn` CLI validates naming:

```bash
$ grgn validate
✓ Checking domain prefixes...
✓ Checking GraphQL type names...
✓ Checking database labels...
✗ ERROR: Type 'User' in twitter/tweet/model/types.graphql missing prefix
         Expected: 'TwitterTweetUser' or reference to 'CoreAuthUser'
```

### 4.5 Cross-Domain References

Product domains reference core types without redeclaring:

```graphql
# twitter/timeline/model/types.graphql
type TwitterTimeline {
  owner: CoreAuthUser!          # Reference to core/auth type
  tweets: [TwitterTweetPost!]!  # Reference to twitter/tweet type
  lastUpdated: DateTime!        # Reference to core/shared scalar
}
```

---

## 5. MVC Redefined for 2026

### 5.1 Model Layer

**Definition**: Declarative `.graphql` files are the single source of truth.

**Structure**:
```
model/
├── types.graphql        # Entity definitions
├── enums.graphql        # Enumeration types
├── inputs.graphql       # Input types for mutations
├── interfaces.graphql   # Shared interfaces
└── directives.graphql   # Custom directives
```

**Rules**:
- One responsibility per file (prevents "schema bloat")
- All types must follow naming conventions
- Relationships defined via field references
- Validation rules as directives

**Example**:
```graphql
# twitter/tweet/model/types.graphql
"""
A tweet represents a single post in the Twitter domain.
"""
type TwitterTweetPost @requiresAuth {
  id: ID!
  author: CoreAuthUser!
  content: String! @length(max: 280)
  status: TWITTER_TWEET_STATUS!
  media: [TwitterTweetMedia!]!
  metrics: TwitterTweetMetrics!
  createdAt: DateTime!
  updatedAt: DateTime!
}
```

### 5.2 View Layer

**Definition**: The "Consumer" of the domain. This is NOT just HTML.

**View Types**:

| Type | Location | Purpose |
|------|----------|---------|
| Web UI | `view/web/` | React components for browser |
| Mobile API | `view/mobile/` | REST/GraphQL handlers for native apps |
| CLI Tool | `view/cli/` | Admin/operator command-line tools |
| Background Job | `view/jobs/` | Scheduled tasks, workers, CRON |
| Admin Dashboard | `view/admin/` | Internal admin interfaces |

**Example Structure**:
```
tweet/view/
├── web/
│   ├── TweetCard.tsx
│   ├── TweetComposer.tsx
│   └── TweetThread.tsx
├── jobs/
│   ├── sentiment_analyzer.go    # Runs every hour
│   ├── trending_calculator.go   # Runs every 5 min
│   └── cleanup_deleted.go       # Runs daily
└── cli/
    └── tweet_admin.go           # grgn tweet:admin commands
```

**Design Principle**: When creating a new module, always ask:
- "What is the Web view for this?"
- "What is the CLI view for this?"
- "What background jobs does this need?"

### 5.3 Controller Layer

**Definition**: The "Imperative Logic" that orchestrates business operations.

**Structure**:
```
controller/
├── resolver.go           # GraphQL resolver (entry point)
├── {feature}_handler.go  # Feature-specific logic
├── {policy}_policy.go    # Business rules
├── interfaces.go         # Exported interfaces for SDK
└── generated/            # Machine-generated code
    ├── graphql/          # gqlgen output
    ├── repository/       # Repository stubs
    └── mappers/          # Type converters
```

**Separation of Concerns**:
- `resolver.go`: Thin layer, delegates to handlers
- `*_handler.go`: Implements use cases
- `*_policy.go`: Encapsulates business rules
- `interfaces.go`: Defines what other domains can use

**Example**:
```go
// twitter/tweet/controller/post_handler.go
package tweet

import (
    "context"
    "github.com/yourorg/grgn-stack/internal/core/auth"
    "github.com/yourorg/grgn-stack/internal/core/shared"
)

type PostHandler struct {
    authService   auth.IUserService       // From core/auth
    mediaService  shared.IMediaProcessor  // From core/shared
    repository    IPostRepository         // Local interface
    config        *TweetConfig            // From service_config.yaml
}

func NewPostHandler(
    authService auth.IUserService,
    mediaService shared.IMediaProcessor,
    repo IPostRepository,
    config *TweetConfig,
) *PostHandler {
    return &PostHandler{
        authService:  authService,
        mediaService: mediaService,
        repository:   repo,
        config:       config,
    }
}

func (h *PostHandler) CreatePost(ctx context.Context, input CreatePostInput) (*Post, error) {
    // 1. Validate via core auth
    user, err := h.authService.GetCurrentUser(ctx)
    if err != nil {
        return nil, err
    }
    
    // 2. Apply business policy
    if err := h.validateContent(input.Content); err != nil {
        return nil, err
    }
    
    // 3. Process media via core shared
    mediaURLs, err := h.mediaService.ProcessMedia(ctx, input.Media)
    if err != nil {
        return nil, err
    }
    
    // 4. Persist via repository
    return h.repository.Create(ctx, user.ID, input.Content, mediaURLs)
}
```

---

## 6. Service Configuration System

### 6.1 Configuration Hierarchy

```
service_config.yaml (root default)
    └── core/service_config.yaml
        └── core/auth/service_config.yaml
            └── environment overrides (.env)
```

Configuration merges top-down with lower levels overriding higher levels.

### 6.2 Schema Definition

```yaml
# service_config.yaml (Root - distributed with stack)
version: "1.0"
domain:
  name: grgn
  prefix: Grgn

defaults:
  pagination:
    defaultPageSize: 20
    maxPageSize: 100
  rateLimit:
    requestsPerMinute: 100
    burstSize: 20

# ============================================
# ARCHITECTURE POLICY (Developer's Choice)
# ============================================
# Choose a preset or define custom rules
architecture:
  # Presets: strict | relaxed | open | custom
  isolation: relaxed
  
  # Naming validation: required | recommended | disabled
  naming: recommended
  
  # Behavior on violations: error | warn | ignore
  onViolation: warn
  
# See Section 11 for full policy options and examples
# ============================================

# core/auth/service_config.yaml
extends: ../service_config.yaml

domain:
  name: core
  prefix: Core

app:
  name: auth
  
config:
  token:
    accessExpiry: 15m
    refreshExpiry: 7d
    issuer: "grgn-auth"
  
  mfa:
    enabled: true
    providers:
      - totp
      - sms
    totpIssuer: "GRGN Stack"
  
  oauth:
    google:
      enabled: ${GRGN_OAUTH_GOOGLE_ENABLED:false}
      clientId: ${GRGN_OAUTH_GOOGLE_CLIENT_ID}
    apple:
      enabled: ${GRGN_OAUTH_APPLE_ENABLED:false}
      clientId: ${GRGN_OAUTH_APPLE_CLIENT_ID}

  validation:
    password:
      minLength: 8
      requireUppercase: true
      requireNumber: true
      requireSpecial: true

# Environment variable interpolation syntax:
# ${VAR_NAME}           - Required, fails if missing
# ${VAR_NAME:default}   - Optional with default value
```

### 6.3 Go Integration (uber-go/config)

```go
package config

import (
    "go.uber.org/config"
)

type AuthConfig struct {
    Token TokenConfig `yaml:"token"`
    MFA   MFAConfig   `yaml:"mfa"`
    OAuth OAuthConfig `yaml:"oauth"`
}

type TokenConfig struct {
    AccessExpiry  time.Duration `yaml:"accessExpiry"`
    RefreshExpiry time.Duration `yaml:"refreshExpiry"`
    Issuer        string        `yaml:"issuer"`
}

func LoadAuthConfig() (*AuthConfig, error) {
    // Load with inheritance
    provider, err := config.NewYAML(
        config.File("service_config.yaml"),
        config.File("core/service_config.yaml"),
        config.File("core/auth/service_config.yaml"),
        config.Expand(os.LookupEnv), // Environment variable expansion
    )
    if err != nil {
        return nil, err
    }
    
    var cfg AuthConfig
    if err := provider.Get("config").Populate(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

### 6.4 Configuration Validation

The `grgn` CLI validates configuration:

```bash
$ grgn config:validate
✓ Loading root service_config.yaml
✓ Loading core/service_config.yaml
✓ Loading core/auth/service_config.yaml
✓ Validating schema compliance
✓ Checking environment variables
  ⚠ GRGN_OAUTH_GOOGLE_CLIENT_ID not set (using default: disabled)
  ⚠ GRGN_OAUTH_APPLE_CLIENT_ID not set (using default: disabled)
✓ Configuration valid
```

---

## 7. Internal SDK Pattern

### 7.1 Concept

Product domains consume core services through **well-defined interfaces**, never raw implementations:

```
┌──────────────────┐     Interface      ┌─────────────────────┐
│  twitter/tweet   │ ─────────────────► │  core/shared        │
│  PostHandler     │                    │  IMailer            │
└──────────────────┘                    └──────────┬──────────┘
                                                   │
                                        ┌──────────┴──────────┐
                                        │  Implementation     │
                                        │  SendGridMailer     │
                                        │  PostmarkMailer     │
                                        └─────────────────────┘
```

### 7.2 Interface Definition

```go
// pkg/grgn/mailer.go - Standalone, importable
package grgn

import "context"

// IMailer defines the contract for email services.
// Product domains depend on this interface, not implementations.
type IMailer interface {
    // SendEmail sends a single email
    SendEmail(ctx context.Context, req EmailRequest) error
    
    // SendTemplate sends a templated email
    SendTemplate(ctx context.Context, template string, data any, recipients []string) error
    
    // SendBulk sends emails in batch
    SendBulk(ctx context.Context, requests []EmailRequest) ([]EmailResult, error)
}

type EmailRequest struct {
    To      []string
    Subject string
    Body    string
    IsHTML  bool
}

type EmailResult struct {
    To     string
    Status EmailStatus
    Error  error
}
```

### 7.3 Implementation

```go
// internal/core/shared/controller/mailer.go
package shared

import (
    "github.com/yourorg/grgn-stack/pkg/grgn"
    "github.com/sendgrid/sendgrid-go"
)

// SendGridMailer implements grgn.IMailer using SendGrid
type SendGridMailer struct {
    client *sendgrid.Client
    config *MailerConfig
}

func NewSendGridMailer(config *MailerConfig) grgn.IMailer {
    return &SendGridMailer{
        client: sendgrid.NewSendClient(config.APIKey),
        config: config,
    }
}

func (m *SendGridMailer) SendEmail(ctx context.Context, req grgn.EmailRequest) error {
    // Implementation using SendGrid
}
```

### 7.4 Consumer Usage

```go
// internal/twitter/tweet/controller/post_handler.go
package tweet

import "github.com/yourorg/grgn-stack/pkg/grgn"

type PostHandler struct {
    mailer grgn.IMailer // Interface, not SendGridMailer
}

func (h *PostHandler) NotifyMentions(ctx context.Context, post *Post) error {
    mentions := extractMentions(post.Content)
    
    for _, mention := range mentions {
        err := h.mailer.SendTemplate(ctx, "mention_notification", map[string]any{
            "mentioner": post.Author.Name,
            "content":   post.Content,
        }, []string{mention.Email})
        
        if err != nil {
            // Log but don't fail - mentions are best-effort
            log.Warn("failed to notify mention", "user", mention.ID, "error", err)
        }
    }
    
    return nil
}
```

### 7.5 Benefits

1. **Single Point of Change**: Swap SendGrid → Postmark in one file
2. **Testability**: Mock `IMailer` in tests, no external calls
3. **Consistency**: All email errors handled uniformly
4. **Observability**: Add logging/metrics in one place

---

## 8. Multi-Tenancy Architecture

### 8.1 Isolation Strategies

GRGN supports **configurable tenant isolation** based on tenant requirements:

| Strategy | Use Case | Implementation |
|----------|----------|----------------|
| **Property-based** | Small tenants, shared infrastructure | `tenant_id` property on all nodes |
| **Database-based** | Enterprise tenants, compliance requirements | Separate Neo4j database per tenant |
| **Hybrid** | Mixed tenant tiers | Small tenants share DB, enterprise get dedicated |

### 8.2 Property-Based Isolation

```graphql
# core/tenant/model/types.graphql
type CoreTenantInstance {
  id: ID!
  name: String!
  tier: CORE_TENANT_TIER!
  isolationMode: CORE_TENANT_ISOLATION!
  databaseName: String  # Only for DATABASE isolation
  createdAt: DateTime!
}

enum CORE_TENANT_ISOLATION {
  SHARED    # Property-based, tenant_id on nodes
  DEDICATED # Separate database
}
```

**Neo4j Implementation**:
```cypher
// All tenant-specific nodes have tenant_id
CREATE (t:TwitterTweetPost {
  id: randomUUID(),
  tenant_id: $tenantId,  // Always present
  content: $content,
  // ...
})

// Shared nodes (core/shared) have no tenant_id
CREATE (c:CoreSharedConfig {
  id: randomUUID(),
  key: "feature_flags",
  // No tenant_id - visible to all
})
```

**Query Enforcement**:
```go
// internal/core/tenant/controller/query_builder.go
func (qb *QueryBuilder) ForTenant(ctx context.Context, query string) string {
    tenant := TenantFromContext(ctx)
    if tenant == nil || tenant.IsolationMode == "DEDICATED" {
        return query // No modification needed for dedicated DB
    }
    
    // Inject tenant filter for shared mode
    return fmt.Sprintf(`
        WITH $tenantId as tid
        %s
        WHERE n.tenant_id = tid OR n.tenant_id IS NULL
    `, query)
}
```

### 8.3 Database-Based Isolation (Neo4j Fabric)

```yaml
# neo4j.conf
fabric.database.name=fabric
fabric.graph.0.name=core
fabric.graph.0.uri=neo4j://core-db:7687
fabric.graph.1.name=tenant_acme
fabric.graph.1.uri=neo4j://tenant-acme-db:7687
fabric.graph.2.name=tenant_globex
fabric.graph.2.uri=neo4j://tenant-globex-db:7687
```

**Fabric Queries**:
```cypher
// Query across tenant + core
USE fabric.core
MATCH (u:CoreAuthUser {id: $userId})
USE fabric.tenant_acme
MATCH (t:TwitterTweetPost {author_id: $userId})
RETURN u, t
```

### 8.4 Tenant Provisioning

```go
// internal/core/tenant/controller/provisioner.go
func (p *Provisioner) CreateTenant(ctx context.Context, input CreateTenantInput) (*Tenant, error) {
    tenant := &Tenant{
        ID:            uuid.New().String(),
        Name:          input.Name,
        Tier:          input.Tier,
        IsolationMode: determineIsolation(input.Tier),
    }
    
    if tenant.IsolationMode == IsolationDedicated {
        // Provision dedicated database
        dbName, err := p.provisionDatabase(ctx, tenant.ID)
        if err != nil {
            return nil, err
        }
        tenant.DatabaseName = dbName
        
        // Run migrations on new database
        if err := p.runMigrations(ctx, dbName); err != nil {
            return nil, err
        }
    }
    
    return p.repository.Create(ctx, tenant)
}

func determineIsolation(tier TenantTier) IsolationMode {
    switch tier {
    case TierEnterprise, TierCompliance:
        return IsolationDedicated
    default:
        return IsolationShared
    }
}
```

### 8.5 Tenant Context Middleware

```go
// internal/core/tenant/controller/middleware.go
func TenantMiddleware(resolver TenantResolver) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract tenant from subdomain, header, or JWT
        tenantID := extractTenantID(c)
        
        tenant, err := resolver.Resolve(c.Request.Context(), tenantID)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid tenant"})
            return
        }
        
        // Add tenant to context
        ctx := WithTenant(c.Request.Context(), tenant)
        c.Request = c.Request.WithContext(ctx)
        
        // Configure database session for tenant
        if tenant.IsolationMode == IsolationDedicated {
            ctx = WithDatabase(ctx, tenant.DatabaseName)
            c.Request = c.Request.WithContext(ctx)
        }
        
        c.Next()
    }
}
```

---

## 9. CLI Tool Specification (grgn)

### 9.1 Overview

The `grgn` CLI is the primary development tool for the GRGN stack, inspired by Laravel Artisan, AdonisJS Ace, and Django Admin.

```bash
$ grgn --help
GRGN Stack CLI v1.0.0

Usage:
  grgn <command> [options]

Available Commands:
  init          Initialize a new GRGN project
  make          Generate code (model, controller, view, migration)
  migrate       Run database migrations
  validate      Validate project structure and naming
  config        Configuration management
  serve         Start development server
  deploy        Deploy to cloud providers
  domain        Domain management commands

Flags:
  -h, --help      Show help
  -v, --version   Show version
  --verbose       Verbose output
```

### 9.2 Command Reference

#### `grgn init`
```bash
$ grgn init my-project
$ grgn init my-project --domain=acme --template=saas
```
Creates a new GRGN project with the specified domain name.

#### `grgn make`
```bash
# Generate a complete app with model, controller, and view
$ grgn make:app twitter/tweet

# Generate individual components
$ grgn make:model twitter/tweet Post
$ grgn make:controller twitter/tweet PostHandler
$ grgn make:view twitter/tweet/web TweetCard
$ grgn make:migration twitter/tweet add_media_support
$ grgn make:job twitter/tweet SentimentAnalyzer
$ grgn make:resolver twitter/tweet

# Scaffold from GraphQL model (creates controller + view + migration)
$ grgn make:scaffold twitter/tweet --from=model/types.graphql
```

#### `grgn migrate`

See [Section 11: Database Schema Management](#11-database-schema-management) for full details.

```bash
$ grgn migrate              # Run all pending migrations
$ grgn migrate:down         # Rollback last migration
$ grgn migrate:down --steps=3  # Rollback N migrations
$ grgn migrate:status       # Show migration status
$ grgn migrate:validate     # Check migration files for issues

# Create new migrations
$ grgn migrate:create twitter/tweet add_reactions        # Cypher migration
$ grgn migrate:create twitter/tweet data_backfill --go   # Go migration

# Domain/app filtering
$ grgn migrate --domain=twitter
$ grgn migrate --app=twitter/tweet

# Tenant-specific migrations
$ grgn migrate --tenant=acme
$ grgn migrate --parallel   # All tenants in parallel

# Safety options
$ grgn migrate --dry-run    # Show what would run
```

#### `grgn validate`
```bash
$ grgn validate              # Validate against YOUR policy in service_config.yaml
$ grgn validate:naming       # Check naming conventions (if enabled)
$ grgn validate:imports      # Check import restrictions (if enabled)
$ grgn validate:config       # Validate service configs
$ grgn validate:schema       # Validate GraphQL schemas

# Policy management
$ grgn validate:policy       # Show current validation policy
$ grgn validate:dry-run      # Preview what would be validated
$ grgn validate:infer        # Generate policy from existing codebase

# Options
$ grgn validate --fix        # Auto-resolve issues where possible
$ grgn validate --ci         # CI mode: exit non-zero on warnings
$ grgn validate --changed-only  # Only validate changed files
$ grgn validate --skip-all   # Skip all validation (emergency override)
```

#### `grgn config`
```bash
$ grgn config:show           # Display merged configuration
$ grgn config:show auth      # Show auth-specific config
$ grgn config:validate       # Validate all configs
$ grgn config:env            # Show required environment variables
```

#### `grgn domain`
```bash
$ grgn domain:list           # List all domains
$ grgn domain:rename old new # Rename domain (fresh projects only)
$ grgn domain:add products   # Add new product domain
```

#### `grgn deploy`
```bash
$ grgn deploy:docker         # Build Docker images
$ grgn deploy:aws            # Deploy to AWS (ECS/EKS)
$ grgn deploy:gcp            # Deploy to GCP (Cloud Run/GKE)
$ grgn deploy:azure          # Deploy to Azure (AKS)
```

### 9.3 Implementation Location

```
/cmd/grgn/                   # CLI entry point
├── main.go
├── commands/
│   ├── init.go
│   ├── make.go
│   ├── migrate.go
│   ├── validate.go
│   ├── config.go
│   ├── domain.go
│   └── deploy.go
├── templates/               # Code generation templates
│   ├── model.graphql.tmpl
│   ├── controller.go.tmpl
│   ├── resolver.go.tmpl
│   ├── migration.go.tmpl
│   └── view/
│       ├── component.tsx.tmpl
│       └── job.go.tmpl
└── validators/              # Validation logic
    ├── naming.go
    ├── imports.go
    └── config.go
```

### 9.4 Example Workflow

```bash
# 1. Create new app
$ grgn make:app twitter/notifications

# 2. Define model (manually edit or use arrows.app)
$ code internal/twitter/notifications/model/types.graphql

# 3. Generate from model
$ grgn make:scaffold twitter/notifications

# Output:
# ✓ Created controller/resolver.go
# ✓ Created controller/notification_handler.go
# ✓ Created controller/generated/graphql/
# ✓ Created migration 003_notifications.go
# ✓ Created view/web/NotificationList.tsx
# ✓ Created view/jobs/notification_sender.go
# ✓ Updated schema/schema.graphql

# 4. Run migrations
$ grgn migrate

# 5. Validate
$ grgn validate
```

---

## 10. Code Generation

### 10.1 Generated Directory Structure

Each app's `generated/` directory contains machine-produced code:

```
controller/generated/
├── graphql/              # gqlgen output
│   ├── generated.go      # Schema execution code
│   ├── models_gen.go     # Go structs from GraphQL
│   └── resolver.go       # Resolver interface
├── repository/           # Repository stubs
│   ├── interfaces.go     # Repository interfaces
│   └── neo4j_impl.go     # Neo4j implementation stubs
└── mappers/              # Type converters
    ├── graphql_to_model.go
    └── model_to_graphql.go
```

### 10.2 Generation Flow

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│ model/*.graphql │ ──► │   grgn generate  │ ──► │   generated/    │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                              │
                              ├── gqlgen (GraphQL → Go)
                              ├── Repository generator
                              └── Mapper generator
```

### 10.3 gqlgen Configuration

```yaml
# gqlgen.yml (per-app)
schema:
  - model/*.graphql
  - ../../core/shared/model/*.graphql  # Include shared scalars

exec:
  filename: controller/generated/graphql/generated.go
  package: graphql

model:
  filename: controller/generated/graphql/models_gen.go
  package: graphql

resolver:
  layout: follow-schema
  dir: controller/generated/graphql
  package: graphql

autobind:
  - github.com/yourorg/grgn-stack/internal/twitter/tweet/controller

models:
  DateTime:
    model: github.com/yourorg/grgn-stack/pkg/grgn.DateTime
  UUID:
    model: github.com/yourorg/grgn-stack/pkg/grgn.UUID
```

### 10.4 Repository Generation

From GraphQL types, generate repository interfaces:

```graphql
# Input
type TwitterTweetPost {
  id: ID!
  author: CoreAuthUser!
  content: String!
  createdAt: DateTime!
}
```

```go
// Output: controller/generated/repository/interfaces.go
package repository

type IPostRepository interface {
    Create(ctx context.Context, post *Post) (*Post, error)
    FindByID(ctx context.Context, id string) (*Post, error)
    FindByAuthor(ctx context.Context, authorID string, pagination Pagination) ([]*Post, error)
    Update(ctx context.Context, id string, updates PostUpdates) (*Post, error)
    Delete(ctx context.Context, id string) error
}
```

---

## 11. Database Schema Management

GRGN uses **golang-migrate** for database schema migrations, wrapped by the `grgn` CLI for ease of use. This provides a standard, well-tested migration workflow adapted for Neo4j and the modular monolith structure.

### 11.1 Overview

**Key Characteristics:**
- **Sequential versioned migrations**: Standard `{version}_{description}.{direction}.{ext}` naming
- **Dual format support**: Cypher files for simple migrations, Go files for complex logic
- **Hybrid location strategy**: Core migrations centralized, product domain migrations colocated
- **Domain priority ordering**: Core always migrates first, then product domains alphabetically
- **Reversibility required**: Every `up` migration must have a corresponding `down` migration
- **Multi-tenant aware**: Configurable per-tenant based on isolation mode
- **Embedded library**: golang-migrate embedded in `grgn` CLI (no external dependencies)

### 11.2 Migration File Structure

```
/backend/
├── migrations/                      # CENTRAL - Core domain migrations
│   ├── 000001_initial_schema.up.cypher
│   ├── 000001_initial_schema.down.cypher
│   ├── 000002_auth_providers.up.cypher
│   ├── 000002_auth_providers.down.cypher
│   ├── 000003_tenant_support.up.go      # Complex migration in Go
│   ├── 000003_tenant_support.down.go
│   └── migrations.go                # Go migration registry
│
├── internal/
│   ├── core/                        # Core migrations in central /migrations/
│   │   ├── auth/
│   │   ├── tenant/
│   │   └── ...
│   │
│   └── twitter/                     # Product domain migrations colocated
│       ├── migrations/              # Twitter-wide migrations
│       │   ├── 000001_twitter_init.up.cypher
│       │   └── 000001_twitter_init.down.cypher
│       │
│       ├── tweet/
│       │   └── migrations/          # Tweet-specific migrations
│       │       ├── 000001_tweet_schema.up.cypher
│       │       ├── 000001_tweet_schema.down.cypher
│       │       ├── 000002_add_media.up.cypher
│       │       └── 000002_add_media.down.cypher
│       │
│       └── timeline/
│           └── migrations/
│               ├── 000001_timeline_cache.up.cypher
│               └── 000001_timeline_cache.down.cypher
```

### 11.3 Migration File Formats

#### Cypher Files (Simple Migrations)

```cypher
-- 000001_initial_schema.up.cypher
-- Create User node constraints and indexes

CREATE CONSTRAINT user_id_unique IF NOT EXISTS
FOR (u:User) REQUIRE u.id IS UNIQUE;

CREATE CONSTRAINT user_email_unique IF NOT EXISTS
FOR (u:User) REQUIRE u.email IS UNIQUE;

CREATE INDEX user_created_at IF NOT EXISTS
FOR (u:User) ON (u.createdAt);
```

```cypher
-- 000001_initial_schema.down.cypher
-- Reverse the initial schema

DROP INDEX user_created_at IF EXISTS;
DROP CONSTRAINT user_email_unique IF EXISTS;
DROP CONSTRAINT user_id_unique IF EXISTS;
```

#### Go Files (Complex Migrations)

Use Go files when you need:
- Conditional logic based on existing data
- Data transformations
- External API calls during migration
- Complex multi-step operations

```go
// 000003_tenant_support.up.go
package migrations

import (
    "context"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
    Register(Migration{
        Version:     3,
        Description: "Add tenant support to existing nodes",
        Up: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
            // Add tenant_id to existing users (default tenant)
            if _, err := tx.Run(ctx, `
                MATCH (u:User) WHERE u.tenant_id IS NULL
                SET u.tenant_id = 'default'
            `, nil); err != nil {
                return err
            }
            
            // Create tenant constraint
            if _, err := tx.Run(ctx, `
                CREATE CONSTRAINT tenant_id_unique IF NOT EXISTS
                FOR (t:Tenant) REQUIRE t.id IS UNIQUE
            `, nil); err != nil {
                return err
            }
            
            // Create default tenant
            _, err := tx.Run(ctx, `
                MERGE (t:Tenant {id: 'default', name: 'Default', createdAt: datetime()})
            `, nil)
            return err
        },
        Down: func(ctx context.Context, tx neo4j.ManagedTransaction) error {
            // Remove default tenant
            if _, err := tx.Run(ctx, `
                MATCH (t:Tenant {id: 'default'}) DELETE t
            `, nil); err != nil {
                return err
            }
            
            // Remove tenant_id from users
            if _, err := tx.Run(ctx, `
                MATCH (u:User) REMOVE u.tenant_id
            `, nil); err != nil {
                return err
            }
            
            // Drop constraint
            _, err := tx.Run(ctx, `
                DROP CONSTRAINT tenant_id_unique IF EXISTS
            `, nil)
            return err
        },
    })
}
```

### 11.4 Migration Execution Order

Migrations execute in **domain priority order**:

```
1. /backend/migrations/          (core - always first)
2. /internal/core/*/migrations/  (core apps - alphabetical)
3. /internal/{product}/migrations/ (product domains - alphabetical)
4. /internal/{product}/*/migrations/ (product apps - alphabetical)
```

**Example execution order:**
```
[core]     migrations/000001_initial_schema
[core]     migrations/000002_auth_providers
[core]     migrations/000003_tenant_support
[core]     internal/core/auth/migrations/000001_mfa_tables
[core]     internal/core/tenant/migrations/000001_billing
[product]  internal/commerce/migrations/000001_commerce_init
[product]  internal/commerce/cart/migrations/000001_cart_schema
[product]  internal/twitter/migrations/000001_twitter_init
[product]  internal/twitter/tweet/migrations/000001_tweet_schema
[product]  internal/twitter/tweet/migrations/000002_add_media
[product]  internal/twitter/timeline/migrations/000001_timeline_cache
```

### 11.5 CLI Commands

The `grgn` CLI wraps golang-migrate for a streamlined experience:

```bash
# Run all pending migrations
$ grgn migrate
Running migrations...
  [core] 000001_initial_schema ✓
  [core] 000002_auth_providers ✓
  [twitter/tweet] 000001_tweet_schema ✓
Migrations complete: 3 applied, 0 failed

# Migrate specific domain only
$ grgn migrate --domain=twitter
$ grgn migrate --domain=core

# Migrate specific app only
$ grgn migrate --app=twitter/tweet

# Rollback last migration
$ grgn migrate:down
Rolling back: [twitter/tweet] 000001_tweet_schema
Rollback complete

# Rollback N migrations
$ grgn migrate:down --steps=3

# Rollback to specific version
$ grgn migrate:to --version=000002

# Show migration status
$ grgn migrate:status
Migration Status:
  [core] migrations/
    ✓ 000001_initial_schema (applied: 2026-01-20 10:30:00)
    ✓ 000002_auth_providers (applied: 2026-01-20 10:30:01)
    ○ 000003_tenant_support (pending)
  [twitter/tweet] migrations/
    ✓ 000001_tweet_schema (applied: 2026-01-21 14:00:00)
    ○ 000002_add_media (pending)

# Create new migration
$ grgn migrate:create twitter/tweet add_reactions
Created:
  internal/twitter/tweet/migrations/000003_add_reactions.up.cypher
  internal/twitter/tweet/migrations/000003_add_reactions.down.cypher

# Create Go migration (for complex logic)
$ grgn migrate:create twitter/tweet data_backfill --go
Created:
  internal/twitter/tweet/migrations/000003_data_backfill.up.go
  internal/twitter/tweet/migrations/000003_data_backfill.down.go

# Validate migrations (check for issues)
$ grgn migrate:validate
Validating migrations...
  ✓ All up migrations have corresponding down migrations
  ✓ No gaps in version numbers
  ✓ No duplicate versions across domains
  ⚠ Warning: 000003_data_backfill.down.go is empty

# Dry run (show what would execute)
$ grgn migrate --dry-run
Would apply:
  [core] 000003_tenant_support
  [twitter/tweet] 000002_add_media

# Force migration (skip version check - use with caution)
$ grgn migrate:force --version=000003
```

### 11.6 Multi-Tenant Migration Strategy

Migration behavior depends on tenant isolation mode:

```yaml
# service_config.yaml
migrations:
  # How to handle multi-tenant migrations
  tenantStrategy: auto  # auto | shared | per-tenant
  
  # auto: Determines based on tenant.isolationMode
  #   - SHARED tenants: Run once (shared schema)
  #   - DEDICATED tenants: Run per-database
```

#### Shared Tenants (Property-Based Isolation)

```bash
# Migrations run once against the shared database
$ grgn migrate
Running migrations on shared database...
  [core] 000001_initial_schema ✓
Done.
```

#### Dedicated Tenants (Database-Per-Tenant)

```bash
# Migrations run against each tenant database
$ grgn migrate
Running migrations...
  Database: core
    [core] 000001_initial_schema ✓
  Database: tenant_acme
    [core] 000001_initial_schema ✓
    [twitter/tweet] 000001_tweet_schema ✓
  Database: tenant_globex
    [core] 000001_initial_schema ✓
    [twitter/tweet] 000001_tweet_schema ✓
Done: 3 databases migrated

# Migrate specific tenant only
$ grgn migrate --tenant=acme

# Migrate all tenants in parallel (faster)
$ grgn migrate --parallel
```

### 11.7 Migration Tracking

Migrations are tracked in Neo4j using a `_Migration` node:

```cypher
(:_Migration {
  version: 1,
  description: "initial_schema",
  source: "core",           // Which domain/app
  appliedAt: datetime(),
  checksum: "abc123...",    // File hash for drift detection
  executionTime: 1234       // Milliseconds
})
```

Query migration history:
```cypher
MATCH (m:_Migration) 
RETURN m.version, m.description, m.source, m.appliedAt 
ORDER BY m.appliedAt
```

### 11.8 CI/CD Integration

```yaml
# .github/workflows/ci.yml
jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate Migrations
        run: grgn migrate:validate
      
      - name: Dry Run
        run: grgn migrate --dry-run
        
      - name: Apply Migrations (staging)
        if: github.ref == 'refs/heads/staging'
        run: grgn migrate
        env:
          GRGN_DATABASE_URI: ${{ secrets.STAGING_NEO4J_URI }}
          
      - name: Apply Migrations (production)
        if: github.ref == 'refs/heads/main'
        run: |
          grgn migrate --dry-run  # Safety check
          grgn migrate
        env:
          GRGN_DATABASE_URI: ${{ secrets.PROD_NEO4J_URI }}
```

### 11.9 Best Practices

#### DO:
- **Keep migrations small and focused** - One logical change per migration
- **Always write down migrations** - Required for rollback capability
- **Test migrations locally** - Run up/down/up cycle before committing
- **Use descriptive names** - `000005_add_user_preferences` not `000005_update`
- **Include comments** - Explain why, not just what
- **Use IF EXISTS/IF NOT EXISTS** - Make migrations idempotent where possible

#### DON'T:
- **Modify existing migrations** - Create new ones instead
- **Delete migration files** - Keep history intact
- **Skip version numbers** - Keep sequence continuous
- **Mix schema and data** - Separate constraint changes from data migrations
- **Assume empty database** - Migrations may run on existing data

#### Example: Safe Column Addition

```cypher
-- 000005_add_user_preferences.up.cypher
-- Add preferences to User nodes

-- Add property (safe - Neo4j is schemaless)
-- New users will have this set by application code

-- Create index for queries on preferences
CREATE INDEX user_preferences_theme IF NOT EXISTS
FOR (u:User) ON (u.preferences_theme);

-- Backfill existing users with defaults
MATCH (u:User) WHERE u.preferences_theme IS NULL
SET u.preferences_theme = 'system',
    u.preferences_language = 'en';
```

```cypher
-- 000005_add_user_preferences.down.cypher
-- Remove preferences from User nodes

-- Drop index first
DROP INDEX user_preferences_theme IF EXISTS;

-- Remove properties
MATCH (u:User)
REMOVE u.preferences_theme, u.preferences_language;
```

### 11.10 Golang-Migrate Integration Details

The `grgn` CLI embeds golang-migrate as a library:

```go
// pkg/migrate/migrator.go
package migrate

import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/neo4j"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
    db     *database.Neo4jDB
    config *MigrateConfig
}

func (m *Migrator) Up(ctx context.Context) error {
    // Collect migrations from all sources in priority order
    sources := m.collectMigrationSources()
    
    for _, source := range sources {
        migrator, err := migrate.New(
            source.Path,
            m.db.GetDriver(),
        )
        if err != nil {
            return fmt.Errorf("[%s] failed to create migrator: %w", source.Domain, err)
        }
        
        if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
            return fmt.Errorf("[%s] migration failed: %w", source.Domain, err)
        }
        
        log.Printf("[%s] migrations applied successfully", source.Domain)
    }
    
    return nil
}

func (m *Migrator) collectMigrationSources() []MigrationSource {
    var sources []MigrationSource
    
    // 1. Central core migrations (always first)
    sources = append(sources, MigrationSource{
        Domain: "core",
        Path:   "file://migrations",
    })
    
    // 2. Core app migrations (alphabetical)
    coreApps := m.findMigrationDirs("internal/core/*/migrations")
    sort.Strings(coreApps)
    for _, app := range coreApps {
        sources = append(sources, MigrationSource{
            Domain: filepath.Base(filepath.Dir(app)),
            Path:   "file://" + app,
        })
    }
    
    // 3. Product domain migrations (alphabetical)
    // ... similar pattern
    
    return sources
}
```

---

## 12. Architecture Enforcement (Configurable)

Architecture rules are **defined by developers**, not mandated by the stack. The `grgn` CLI validates against your chosen policies.

### 12.1 Isolation Policy Configuration

Define your isolation strategy in the root `service_config.yaml`:

```yaml
# service_config.yaml (project root)
version: "1.0"
domain:
  name: grgn
  prefix: Grgn

# Architecture validation settings
architecture:
  # Isolation policy: strict | relaxed | open | custom
  isolation: relaxed
  
  # Naming validation: required | recommended | disabled
  naming: recommended
  
  # What to do on violations: error | warn | ignore
  onViolation: warn
  
  # Custom rules (only used when isolation: custom)
  rules:
    # Define explicit allow/deny patterns
    imports:
      # Pattern: which packages can import what
      - from: "internal/twitter/*"
        allow:
          - "internal/core/*"
          - "internal/commerce/*"  # Allow cross-product import
          - "pkg/*"
        deny:
          - "internal/twitter/*/controller/generated/*"  # No importing generated code
      
      - from: "internal/core/*"
        allow:
          - "pkg/*"
        # Implicitly denies internal/* imports
    
    # Naming patterns (regex)
    naming:
      graphqlTypes: "^[A-Z][a-zA-Z]+$"        # Just PascalCase, no prefix required
      graphqlEnums: "^[A-Z][A-Z0-9_]+$"       # UPPER_SNAKE
      goPackages: "^[a-z][a-z0-9]*$"          # lowercase
```

### 12.2 Isolation Policies Explained

#### `strict` - Maximum Isolation
```yaml
architecture:
  isolation: strict
```
- Product domains cannot import from other product domains
- All domains can only import from `core/` and `pkg/`
- Core cannot import from any product domain
- Best for: Large teams, compliance requirements

#### `relaxed` - Declared Dependencies (Default)
```yaml
architecture:
  isolation: relaxed
```
- Cross-domain imports allowed if declared in `service_config.yaml`
- Each domain declares its dependencies explicitly
- Undeclared dependencies trigger warnings/errors
- Best for: Most projects, balanced flexibility

```yaml
# internal/twitter/service_config.yaml
domain:
  name: twitter
  prefix: Twitter
  
  # Explicit dependencies
  dependencies:
    - core/auth      # Can import auth
    - core/shared    # Can import shared infra
    - commerce/cart  # Can import cart (cross-product)
```

#### `open` - No Restrictions
```yaml
architecture:
  isolation: open
```
- No import validation performed
- Naming validation still applies (if enabled)
- Best for: Small teams, rapid prototyping, monolith-first

#### `custom` - Full Control
```yaml
architecture:
  isolation: custom
  rules:
    imports:
      # Your rules here
```
- Define exact patterns for allowed/denied imports
- Maximum flexibility for complex architectures
- Best for: Enterprise with specific requirements

### 12.3 Naming Policy Configuration

Naming rules are also configurable:

```yaml
architecture:
  naming: recommended  # required | recommended | disabled
  
  # Override specific rules
  namingRules:
    # Require domain prefix on GraphQL types?
    requireDomainPrefix: false  # Default: true when naming: required
    
    # Prefix separator style
    prefixStyle: none  # none | pascal | underscore
    # none: "UserProfile" (package provides namespace)
    # pascal: "TwitterUserProfile"  
    # underscore: "Twitter_UserProfile"
    
    # Database labels
    databaseLabels: match_graphql  # match_graphql | custom | none
```

### 12.4 Interactive Validation

The CLI provides options during validation:

```bash
$ grgn validate:imports

Checking import rules (policy: relaxed)...

⚠ WARNING: internal/twitter/timeline imports internal/commerce/pricing
  This cross-domain import is not declared in twitter/service_config.yaml
  
  Options:
    [1] Add to dependencies (update service_config.yaml)
    [2] Ignore this warning
    [3] Fail validation
    [4] Add to global exceptions
  
  Choice (1-4): 1
  
  ✓ Added commerce/pricing to twitter dependencies
  ✓ Updated internal/twitter/service_config.yaml

Validation complete: 0 errors, 1 warning (resolved)
```

### 12.5 CI Integration (Respects Your Policy)

```yaml
# .github/workflows/ci.yml
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install grgn CLI
        run: go install ./cmd/grgn
      
      - name: Validate Architecture
        run: |
          # Validates against YOUR service_config.yaml settings
          grgn validate
          
      # Or run specific validations
      - name: Check Naming (if enabled)
        run: grgn validate:naming --ci  # --ci flag exits non-zero on warnings
        
      - name: Check Imports (if enabled)  
        run: grgn validate:imports --ci
```

### 12.6 Pre-commit Hook (Optional)

```bash
#!/bin/bash
# .git/hooks/pre-commit (install with: grgn hooks:install)

# Read policy from config
POLICY=$(grgn config:get architecture.isolation)

if [ "$POLICY" != "open" ]; then
    grgn validate --changed-only --ci
    
    if [ $? -ne 0 ]; then
        echo ""
        echo "Validation failed. Options:"
        echo "  1. Fix the issues"
        echo "  2. Run: grgn validate --fix (auto-resolve where possible)"
        echo "  3. Update architecture policy in service_config.yaml"
        exit 1
    fi
fi
```

### 12.7 Validation Summary Commands

```bash
# Show current policy
$ grgn validate:policy
Architecture Policy:
  Isolation: relaxed
  Naming: recommended  
  On Violation: warn

# Show what would be validated
$ grgn validate:dry-run
Would validate:
  ✓ Import rules (relaxed policy, 3 domains)
  ✓ Naming conventions (recommended, PascalCase types)
  ○ Schema consistency (disabled)
  
# Skip validation entirely for a session
$ grgn validate --skip-all

# Generate policy from current codebase (useful for adopting on existing projects)
$ grgn validate:infer
Analyzing codebase...
  Found 5 cross-domain imports
  Found 2 naming pattern variants
  
Generated policy saved to service_config.yaml.suggested
Review and merge into your service_config.yaml
```

### 12.8 Escape Hatches

For legitimate exceptions:

```go
// Use directive comments to suppress warnings
//grgn:allow-import internal/legacy/oldcode
import "github.com/yourorg/grgn-stack/internal/legacy/oldcode"
```

```yaml
# Or declare global exceptions in config
architecture:
  exceptions:
    imports:
      - pattern: "internal/legacy/*"
        reason: "Legacy code being migrated"
        expires: "2026-06-01"  # Optional expiration
```

---

## 13. Migration Path

### 13.1 From Current Structure

The existing codebase has:
- `internal/graphql/` - Flat GraphQL setup
- `internal/repository/` - Flat repository
- `internal/database/` - Database connection
- `pkg/config/` - Configuration

### 13.2 Migration Steps

1. **Phase 1: Core Foundation**
   - Create `internal/core/shared/` structure
   - Move database, config to core/shared
   - Extract interfaces to `pkg/grgn/`

2. **Phase 2: Auth Domain**
   - Create `internal/core/auth/`
   - Move user-related code
   - Add service_config.yaml

3. **Phase 3: Tenant Domain**
   - Create `internal/core/tenant/`
   - Implement tenant isolation
   - Add provisioning logic

4. **Phase 4: CLI Tool**
   - Create `cmd/grgn/`
   - Implement validation commands
   - Add code generation

5. **Phase 5: Product Domains**
   - Create sample product domain
   - Document patterns
   - Validate architecture

### 13.3 Backwards Compatibility

During migration:
- Existing endpoints continue to work
- New code follows new patterns
- Gradual migration via facade pattern

---

## Appendix A: Quick Reference

### A.1 File Naming

| Type | Convention | Example |
|------|------------|---------|
| GraphQL Schema | `{purpose}.graphql` | `types.graphql`, `enums.graphql` |
| Go Controller | `{feature}_handler.go` | `post_handler.go` |
| Go Policy | `{rule}_policy.go` | `deletion_policy.go` |
| Go Resolver | `resolver.go` | `resolver.go` |
| React Component | `{Name}.tsx` | `TweetCard.tsx` |
| Background Job | `{name}_job.go` | `sentiment_analyzer.go` |
| Migration (Cypher) | `{NNNNNN}_{desc}.{dir}.cypher` | `000001_init.up.cypher` |
| Migration (Go) | `{NNNNNN}_{desc}.{dir}.go` | `000003_backfill.up.go` |
| Config | `service_config.yaml` | `service_config.yaml` |

### A.2 Import Rules (Depends on Your Policy)

```go
// Always allowed (all policies)
import "github.com/yourorg/grgn-stack/internal/core/auth"
import "github.com/yourorg/grgn-stack/internal/core/shared"
import "github.com/yourorg/grgn-stack/pkg/grgn"

// Cross-product imports - depends on your architecture.isolation setting:

// If isolation: strict → ❌ Not Allowed
// If isolation: relaxed → ⚠ Allowed if declared in dependencies
// If isolation: open → ✅ Allowed
// If isolation: custom → Depends on your rules
import "github.com/yourorg/grgn-stack/internal/commerce/cart"
```

**To declare a dependency (relaxed mode):**
```yaml
# internal/twitter/service_config.yaml
domain:
  name: twitter
  dependencies:
    - core/auth
    - core/shared
    - commerce/cart  # Now allowed!
```

### A.3 GraphQL Naming

```graphql
# Type: {Domain}{App}{Type}
type TwitterTweetPost { ... }
type CoreAuthUser { ... }

# Enum: {DOMAIN}_{APP}_{ENUM}
enum TWITTER_TWEET_STATUS { DRAFT PUBLISHED DELETED }
enum CORE_AUTH_PROVIDER { GOOGLE APPLE SAML LOCAL }

# Input: {Domain}{App}{Action}Input
input TwitterTweetCreateInput { ... }
input CoreAuthLoginInput { ... }

# Query: {domain}{App}{Action}
type Query {
  twitterTweetGet(id: ID!): TwitterTweetPost
  coreAuthMe: CoreAuthUser
}

# Mutation: {domain}{App}{Action}
type Mutation {
  twitterTweetCreate(input: TwitterTweetCreateInput!): TwitterTweetPost
  coreAuthLogin(input: CoreAuthLoginInput!): CoreAuthSession
}
```

---

## Appendix B: Example Timeline Implementation

Complete example showing MVC pattern in practice:

```graphql
# internal/twitter/timeline/model/types.graphql
type TwitterTimeline {
  owner: CoreAuthUser!
  tweets: [TwitterTweetPost!]!
  cursor: String
  hasMore: Boolean!
  lastUpdated: DateTime!
}

extend type Query {
  twitterTimelineHome(cursor: String, limit: Int = 20): TwitterTimeline!
  twitterTimelineUser(userId: ID!, cursor: String, limit: Int = 20): TwitterTimeline!
}
```

```go
// internal/twitter/timeline/controller/resolver.go
package timeline

type Resolver struct {
    feedAlgorithm *FeedAlgorithm
    cacheManager  *CacheManager
}

func (r *Resolver) TwitterTimelineHome(ctx context.Context, cursor *string, limit int) (*Timeline, error) {
    user := auth.UserFromContext(ctx)
    
    // Try cache first
    if cached, ok := r.cacheManager.Get(ctx, user.ID, "home"); ok {
        return cached, nil
    }
    
    // Generate fresh timeline
    timeline, err := r.feedAlgorithm.GenerateHomeFeed(ctx, user.ID, cursor, limit)
    if err != nil {
        return nil, err
    }
    
    // Cache for next request
    r.cacheManager.Set(ctx, user.ID, "home", timeline, 5*time.Minute)
    
    return timeline, nil
}
```

```go
// internal/twitter/timeline/controller/feed_algorithm.go
package timeline

type FeedAlgorithm struct {
    db           database.Neo4jDB
    tweetService tweet.IPostService // Internal SDK
}

func (f *FeedAlgorithm) GenerateHomeFeed(ctx context.Context, userID string, cursor *string, limit int) (*Timeline, error) {
    // Neo4j graph traversal
    query := `
        MATCH (u:CoreAuthUser {id: $userId})-[:FOLLOWS]->(following:CoreAuthUser)
        MATCH (following)-[:POSTED]->(t:TwitterTweetPost)
        WHERE t.status = 'PUBLISHED'
        RETURN t
        ORDER BY t.createdAt DESC
        LIMIT $limit
    `
    
    tweets, err := f.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        // Execute query and map results
    })
    
    return &Timeline{
        OwnerID:     userID,
        Tweets:      tweets,
        LastUpdated: time.Now(),
    }, nil
}
```

---

## Appendix C: Architecture Policy Presets

Ready-to-use policy configurations for common scenarios.

### C.1 Startup / Rapid Prototyping

```yaml
# service_config.yaml - Move fast, minimal rules
architecture:
  isolation: open
  naming: disabled
  onViolation: ignore
```

### C.2 Small Team (Recommended Default)

```yaml
# service_config.yaml - Balanced flexibility and guardrails
architecture:
  isolation: relaxed
  naming: recommended
  onViolation: warn
```

### C.3 Large Team / Multiple Squads

```yaml
# service_config.yaml - Clear boundaries between teams
architecture:
  isolation: strict
  naming: required
  onViolation: error
  
  # Optional: Define team ownership
  teams:
    platform:
      owns: ["core/*"]
    growth:
      owns: ["twitter/*", "notifications/*"]
    commerce:
      owns: ["commerce/*", "payments/*"]
```

### C.4 Enterprise / Compliance

```yaml
# service_config.yaml - Maximum control and auditability
architecture:
  isolation: custom
  naming: required
  onViolation: error
  
  rules:
    imports:
      # PCI-DSS: Payment code is isolated
      - from: "internal/payments/*"
        allow: ["pkg/*"]
        deny: ["internal/*"]  # No internal dependencies
      
      # HIPAA: Health data isolated  
      - from: "internal/health/*"
        allow: ["internal/core/auth", "pkg/*"]
        deny: ["internal/*"]
      
      # Everything else follows relaxed rules
      - from: "internal/*"
        allow: ["internal/core/*", "pkg/*"]
    
  # Audit trail
  audit:
    logViolations: true
    logResolutions: true
    
  # Exception governance
  exceptions:
    requireApproval: true
    maxDuration: 30d
```

### C.5 Migrating Legacy Codebase

```yaml
# service_config.yaml - Start permissive, tighten over time
architecture:
  isolation: relaxed
  naming: recommended
  onViolation: warn
  
  # Grandfather existing violations
  exceptions:
    imports:
      - pattern: "internal/legacy/*"
        reason: "Pre-migration code"
        expires: "2026-12-31"
      - pattern: "internal/*/old_*"
        reason: "Deprecated modules"
        expires: "2026-06-30"
  
  # Track progress
  migration:
    trackViolations: true
    targetPolicy: strict  # Where we're headed
    targetDate: "2027-01-01"
```

### C.6 Monorepo with Multiple Products

```yaml
# service_config.yaml - Product boundaries with shared core
architecture:
  isolation: custom
  naming: required
  onViolation: error
  
  rules:
    imports:
      # Product A can only use core and its own code
      - from: "internal/product-a/*"
        allow: ["internal/core/*", "internal/product-a/*", "pkg/*"]
      
      # Product B can only use core and its own code
      - from: "internal/product-b/*"
        allow: ["internal/core/*", "internal/product-b/*", "pkg/*"]
      
      # Shared libraries available to all
      - from: "internal/shared-libs/*"
        allow: ["pkg/*"]
      
      # Products can use shared libs
      - from: "internal/product-*/*"
        allow: ["internal/shared-libs/*"]
```

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-20 | - | Initial notes and concepts |
| 2.0 | 2026-01-25 | - | Formalized architecture specification |
| 2.1 | 2026-01-25 | - | Made domain isolation configurable (developer choice) |
| 2.2 | 2026-01-25 | - | Added Database Schema Management (golang-migrate, Section 11) |
