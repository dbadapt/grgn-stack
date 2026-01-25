# Contributing to GRGN Stack

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Assume good intentions

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](../../issues)
2. If not, create a new issue using the bug report template
3. Provide detailed information:
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details
   - Relevant logs or screenshots

### Suggesting Features

1. Check if the feature has been suggested in [Issues](../../issues)
2. Create a new issue using the feature request template
3. Explain the problem it solves and proposed solution
4. Be open to discussion and feedback

### Submitting Pull Requests

1. **Fork the repository** and create a branch from `main`

   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**
   - Follow the existing code style
   - Write clear, concise commit messages
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**

   ```bash
   # Run all tests
   npm test

   # Run specific tests
   npm run test:backend    # Backend
   npm run test:frontend   # Frontend
   ```

4. **Ensure code quality**

   ```bash
   # Format Go code
   cd backend && go fmt ./... && cd ..

   # Lint TypeScript
   cd web && npm run lint && cd ..
   ```

5. **Push to your fork**

   ```bash
   git push origin feature/my-feature
   ```

6. **Open a Pull Request**
   - Use the PR template
   - Link related issues
   - Provide clear description of changes
   - Ensure CI checks pass

## Development Setup

### Prerequisites

- Docker & Docker Compose
- Node.js 18+
- Go 1.24+ (see [Go Version Management](#go-version-management) below)

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/dbadapt/grgn-stack.git
cd YOUR_REPO

# Install dependencies
npm install
cd web && npm install && cd ..

# Setup environment
cp .env.example .env
cp web/.env.example web/.env

# Start development environment
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## Project Structure

```
cmd/             # Entry points (server, migrate, worker)
pkg/             # Standalone packages (config)
services/        # Modular monolith domains
  core/          # Infrastructure (shared, auth, tenant)
  {product}/     # Product domains with apps
web/             # React frontend (TypeScript)
scripts/         # Utility scripts
.github/         # CI/CD workflows
```

## Coding Standards

### Go (Backend)

- Follow standard Go conventions
- Use `go fmt` for formatting
- Run `go vet` before committing
- Write table-driven tests
- Document exported functions

### TypeScript (Frontend)

- Use TypeScript strict mode
- Follow ESLint rules
- Use functional components with hooks
- Write tests with Vitest and Testing Library

### GraphQL

- Use clear, descriptive names
- Document queries and mutations
- Keep schema changes backward compatible

## Testing Guidelines

### Go Version Management

> ⚠️ **Important**: This project requires **Go 1.24.0 or higher** due to transitive dependency requirements.

#### Why Go 1.24?

Several dependencies in the Go ecosystem have updated their minimum Go version requirements:

| Package                  | Requires |
| ------------------------ | -------- |
| `gin-contrib/sse@v1.1.0` | Go 1.23+ |
| `golang.org/x/crypto`    | Go 1.24+ |
| `golang.org/x/net`       | Go 1.24+ |
| `golang.org/x/sys`       | Go 1.24+ |
| `golang.org/x/text`      | Go 1.24+ |

#### Files That Must Stay in Sync

When changing Go versions, update **all** of these files:

1. `go.work` - Workspace Go version
2. `backend/go.mod` - Backend module version
3. `pkg/go.mod` - Shared package module version
4. `backend/Dockerfile` - Docker image base
5. `.github/workflows/ci.yml` - CI Go version

#### After Changing Go Versions

```bash
# Always run go mod tidy after version changes
cd backend && go mod tidy && cd ..
cd pkg && go mod tidy && cd ..

# Verify local build
cd backend && go build ./... && cd ..

# Rebuild Docker containers
docker-compose down
docker-compose up --build
```

#### Troubleshooting Version Errors

If you see errors like:

```
module X requires go >= 1.XX (running go 1.YY; GOTOOLCHAIN=local)
```

This means a dependency needs a newer Go version. Check the dependency chain:

```bash
cd backend
go mod graph | Select-String "package-name"  # Windows
go mod graph | grep "package-name"           # Linux/Mac
```

### Backend Tests

```bash
cd backend
go test -v ./...
```

- Test all public functions
- Use mocks for external dependencies
- Aim for >80% coverage

### Frontend Tests

```bash
cd web
npm test
```

- Test user interactions
- Mock GraphQL queries
- Test error states

## Documentation

- Update README.md for user-facing changes
- Update relevant .md files for architectural changes
- Add inline comments for complex logic
- Update CHANGELOG.md (if exists)

## Commit Message Guidelines

Use clear, descriptive commit messages:

```
type(scope): Brief description

Detailed explanation if needed

Fixes #123
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Maintenance tasks

**Examples:**

```
feat(auth): Add Google OAuth integration
fix(graphql): Resolve user query race condition
docs(readme): Update installation instructions
```

## Review Process

1. All PRs require at least one review
2. CI checks must pass
3. Tests must pass
4. Code coverage should not decrease significantly
5. Maintainers may request changes

## Questions?

- Open a [Discussion](../../discussions) for general questions
- Open an [Issue](../../issues) for bugs or features
- Review existing documentation

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

**Thank you for contributing to GRGN Stack! 🚀**
