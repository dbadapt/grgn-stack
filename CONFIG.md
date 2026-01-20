# Multi-Environment Configuration Guide

## Overview

GRGN Stack supports three environments: **development**, **staging**, and **production**.

## Environment Files

### Backend (.env files in project root)

- `.env.example` - Template with all available options
- `.env.development` - Development environment
- `.env.staging` - Staging environment
- `.env.production` - Production environment

All backend env vars are prefixed with `GRGN_STACK_`

### Frontend (.env files in web/ directory)

- `web/.env.example` - Template
- `web/.env.development` - Development environment
- `web/.env.staging` - Staging environment
- `web/.env.production` - Production environment

All frontend env vars must be prefixed with `VITE_` to be exposed to the client.

## Usage

### Development

```bash
# Copy example files
cp .env.example .env
cp web/.env.example web/.env

# Start with development compose
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

### Staging

```bash
# Use staging environment files
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up
```

### Production

```bash
# Use production environment files
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

## Backend Configuration

Configuration is loaded via the `pkg/config` package using Viper:

- Reads from .env files
- Can be overridden by environment variables
- Supports nested configuration with dot notation

Example usage in Go:

```go
import "github.com/yourusername/grgn-stack/pkg/config"

cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}

// Access configuration
dbURI := cfg.Database.Neo4jURI
if cfg.IsProduction() {
    // Production-specific logic
}
```

## Frontend Configuration

Environment variables are accessed via `import.meta.env`:

```typescript
import { env, isDevelopment } from './config/env';

// Access typed configuration
const apiUrl = env.apiUrl;

if (isDevelopment) {
  console.log('Running in development mode');
}
```

## Security Notes

1. **Never commit .env files with secrets** - Only .env.example files should be committed
2. **Use secrets management** in production (AWS Secrets Manager, Azure Key Vault, etc.)
3. **Rotate secrets regularly**
4. **Different secrets per environment**

## Adding New Configuration

### Backend

1. Add to `pkg/config/config.go` struct
2. Add default in `setDefaults()`
3. Add to `.env.example`
4. Update environment-specific .env files

### Frontend

1. Add VITE\_ prefixed variable to `web/.env.example`
2. Update `web/src/config/env.ts` interface
3. Update environment-specific .env files
