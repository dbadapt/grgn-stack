# Code Coverage Guide

## Overview

GRGN Stack uses comprehensive code coverage tracking for both Go backend and TypeScript frontend.

## Frontend Coverage (Vitest)

### Run Coverage

```bash
cd web
npm run test:coverage
```

### Configuration

Coverage is configured in `web/vite.config.ts`:

- **Provider**: V8 (fast, accurate)
- **Reporters**: text, json, html, lcov
- **Thresholds**: 70% for lines, functions, branches, statements
- **Exclusions**: test files, config files, node_modules

### View Reports

- **Terminal**: Coverage summary displayed after tests
- **HTML**: Open `web/coverage/index.html` in browser
- **CI Integration**: LCOV format for CI/CD pipelines

### Coverage Thresholds

Tests will fail if coverage drops below:

- Lines: 70%
- Functions: 70%
- Branches: 70%
- Statements: 70%

## Backend Coverage (Go)

### Run Coverage

```bash
# From project root
npm run coverage:go

# Or directly
cd backend
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### View Reports

- **Terminal**: Coverage summary in test output
- **HTML**: Open `backend/coverage.html` in browser
- **Function-level**: `go tool cover -func=coverage.out`

## Running All Coverage

From project root:

```bash
npm run coverage
```

This runs both Go and TypeScript coverage tests sequentially.

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Test with Coverage
  run: npm run coverage

- name: Upload Coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./web/coverage/lcov.info,./backend/coverage.out
```

## Best Practices

1. **Write tests first** - Aim for 80%+ coverage on critical paths
2. **Test behavior, not implementation** - Focus on what code does, not how
3. **Exclude test utilities** - Configuration already excludes test helpers
4. **Review coverage reports** - Look for untested edge cases
5. **Don't chase 100%** - Focus on meaningful test coverage

## Excluded from Coverage

### Frontend

- Test files (`*.test.ts`, `*.test.tsx`)
- Test setup (`src/test/`)
- Config files (vite.config.ts, eslint.config.js)
- Node modules

### Backend

- Vendor dependencies
- Generated code
- Test files (`*_test.go`)

## Coverage Files (Git Ignored)

- `web/coverage/`
- `backend/coverage.out`
- `backend/coverage.html`
- `*.lcov`
