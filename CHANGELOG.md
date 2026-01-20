# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial GRGN Stack template release
- Automated initialization scripts (Windows & Linux/Mac)
- Comprehensive documentation suite
- GitHub template files
- Docker Compose multi-environment support
- GraphQL code generation (backend & frontend)
- Neo4j database with migrations
- Testing infrastructure (Go & Vitest)
- CI/CD pipelines (GitHub Actions)

### Changed

### Deprecated

### Removed

### Fixed

- **Go 1.24.0 Required**: Upgraded from Go 1.22 to Go 1.24.0 due to transitive dependency requirements:
  - `gin-contrib/sse@v1.1.0` requires Go 1.23+
  - `golang.org/x/crypto`, `golang.org/x/net`, and other `golang.org/x/*` packages require Go 1.24+
  - Updated: `go.work`, `backend/go.mod`, `pkg/go.mod`, `backend/Dockerfile`, `.github/workflows/ci.yml`
  - Run `go mod tidy` after any Go version changes to sync dependencies

### Security

## [1.0.0] - YYYY-MM-DD

### Added

- First stable release of GRGN Stack Template

---

## Template for New Releases

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added

- New features

### Changed

- Changes in existing functionality

### Deprecated

- Soon-to-be removed features

### Removed

- Removed features

### Fixed

- Bug fixes

### Security

- Security improvements
```
