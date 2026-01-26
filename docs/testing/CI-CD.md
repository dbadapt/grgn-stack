# CI/CD Pipeline Documentation

## Overview

GRGN Stack uses **GitHub Actions** for continuous integration and deployment. The pipeline automatically tests, builds, and deploys your application across multiple environments.

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**

- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`

**Jobs:**

#### Backend Job

- Sets up Go 1.24
- Runs `go fmt` check
- Runs `go vet` linting
- Executes tests with race detection
- Spins up Neo4j service for integration tests
- Uploads coverage artifacts

#### Frontend Job

- Sets up Node.js 18
- Runs ESLint
- Performs TypeScript type checking
- Executes Vitest tests
- Builds production bundle

#### GraphQL Job

- Generates backend (gqlgen) and frontend (graphql-codegen) code
- Verifies generated code is up to date
- Fails if uncommitted changes detected

**Status:** Runs on every push/PR for immediate feedback

---

### 2. Coverage Workflow (`.github/workflows/coverage.yml`)

**Triggers:**

- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Features:**

- Generates coverage reports for Go and TypeScript
- Uploads to Codecov for tracking
- Enforces 70% coverage threshold
- Creates coverage summary in GitHub UI
- Stores coverage artifacts for 30 days

**Required Secret:**

- `CODECOV_TOKEN` - Get from [codecov.io](https://codecov.io)

---

### 3. Docker Build Workflow (`.github/workflows/docker.yml`)

**Triggers:**

- Push to `main`, `staging`, or `develop`
- Tags matching `v*.*.*` pattern
- Pull requests to `main`

**Features:**

- Builds multi-platform images (amd64, arm64)
- Pushes to GitHub Container Registry (ghcr.io)
- Uses build cache for faster builds
- Runs Trivy security scanning
- Uploads security results to GitHub Security tab

**Image Tags Generated:**

- Branch name (e.g., `main`, `develop`)
- PR number (e.g., `pr-123`)
- Semantic version (e.g., `v1.2.3`)
- Commit SHA (e.g., `main-abc1234`)
- `latest` (for default branch only)

**Permissions:** Automatically uses `GITHUB_TOKEN` for authentication

---

### 4. Deploy Workflow (`.github/workflows/deploy.yml`)

**Status:** ⚠️ **Currently Disabled** - Manual-only until deployment secrets are configured

**Triggers:**

- ~~Push to `main` → Production~~ (Disabled)
- ~~Push to `staging` → Staging~~ (Disabled)
- ~~Push to `develop` → Development~~ (Disabled)
- Manual dispatch (workflow_dispatch) - **Active**

**Note:** Automatic deployment is currently disabled to prevent failures from missing secrets. The workflow can only be triggered manually from the Actions tab. To re-enable automatic deployment, uncomment the `push:` section in `.github/workflows/deploy.yml` after configuring the required secrets below.

**Deployment Process:**

1. Determines target environment
2. Builds Docker images
3. Tags with commit SHA and environment
4. Pushes to container registry
5. SSHs into deployment server
6. Pulls latest code and images
7. Runs `docker-compose` with environment-specific config
8. Verifies deployment

**Required Secrets:**

```yaml
DOCKER_REGISTRY: your-registry.com
DOCKER_USERNAME: your-username
DOCKER_PASSWORD: your-token
DEPLOY_HOST: server.example.com
DEPLOY_USER: deploy
DEPLOY_SSH_KEY: |
  -----BEGIN OPENSSH PRIVATE KEY-----
  ...
  -----END OPENSSH PRIVATE KEY-----
```

**Required Variables (per environment):**

- `APP_URL` - Application URL for verification

---

## Setup Instructions

### 1. Initial GitHub Setup

```bash
# Enable GitHub Actions in your repository
# Settings → Actions → General → Allow all actions
```

### 2. Configure Secrets

Go to **Settings → Secrets and variables → Actions**:

**Repository Secrets:**

- `CODECOV_TOKEN` (optional, for coverage tracking)
- `DOCKER_REGISTRY`
- `DOCKER_USERNAME`
- `DOCKER_PASSWORD`
- `DEPLOY_HOST`
- `DEPLOY_USER`
- `DEPLOY_SSH_KEY`

**Environment Secrets (per environment):**
Create environments: `development`, `staging`, `production`

- Add environment-specific secrets as needed

### 3. Configure Environment Variables

**Settings → Secrets and variables → Actions → Variables**:

- `APP_URL` for each environment

### 4. Enable Container Registry

GitHub Container Registry is automatically enabled. Images will be pushed to:

```
ghcr.io/<your-username>/grgn-stack-backend
ghcr.io/<your-username>/grgn-stack-web
```

### 5. Setup Codecov (Optional)

1. Visit [codecov.io](https://codecov.io)
2. Sign in with GitHub
3. Add your repository
4. Copy the token
5. Add as `CODECOV_TOKEN` secret

---

## Branch Strategy

```
main (production)
├── staging
│   └── develop
│       └── feature/*
```

**Branch Policies:**

- `main` → Production deployment (requires PR review)
- `staging` → Staging deployment (for QA)
- `develop` → Development deployment (active development)
- `feature/*` → CI only (no deployment)

---

## Workflow Status Badges

Add to README.md:

```markdown
![CI](https://github.com/<username>/grgn-stack/workflows/CI/badge.svg)
![Coverage](https://github.com/<username>/grgn-stack/workflows/Coverage/badge.svg)
![Docker](https://github.com/<username>/grgn-stack/workflows/Docker%20Build/badge.svg)
[![codecov](https://codecov.io/gh/<username>/grgn-stack/branch/main/graph/badge.svg)](https://codecov.io/gh/<username>/grgn-stack)
```

---

## Manual Deployment

You can trigger manual deployments:

1. Go to **Actions** tab
2. Select **Deploy** workflow
3. Click **Run workflow**
4. Choose target environment
5. Click **Run workflow**

---

## Troubleshooting

### Build Failures

**Go tests failing:**

```bash
# Run locally with same settings
cd backend
go test -v -race ./...
```

**Frontend tests failing:**

```bash
# Run locally
cd web
npm test
```

**GraphQL codegen not up to date:**

```bash
npm run generate
git add .
git commit -m "chore: update generated GraphQL code"
```

### Deployment Failures

**Check logs:**

1. Go to Actions tab
2. Click failed workflow run
3. Expand failed job
4. Review error messages

**Common issues:**

- SSH key not configured correctly
- Docker registry authentication failed
- Server disk space full
- Environment variables missing

**Manual rollback:**

```bash
# SSH into server
ssh deploy@server.example.com
cd /opt/grgn-stack
git checkout <previous-commit>
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## Performance Optimization

**Cache Strategy:**

- Go modules cached by branch
- NPM packages cached by `package-lock.json`
- Docker layers cached with GitHub Actions cache

**Speed Tips:**

- Keep dependencies minimal
- Use `npm ci` instead of `npm install`
- Leverage multi-stage Docker builds
- Use `.dockerignore` to reduce context size

---

## Security

**Best Practices:**

- Never commit secrets to repository
- Use environment-specific secrets
- Enable branch protection rules
- Require status checks before merging
- Enable Trivy security scanning
- Review security alerts in Security tab
- Rotate secrets regularly

**Branch Protection:**
Settings → Branches → Add rule for `main`:

- ✅ Require pull request reviews
- ✅ Require status checks (CI, Coverage)
- ✅ Require branches to be up to date
- ✅ Require linear history

---

## Cost Optimization

GitHub Actions is free for public repositories and includes:

- 2,000 minutes/month for private repositories (free tier)
- Unlimited for public repositories

**Tips to reduce usage:**

- Use `cache` actions effectively
- Skip workflows for documentation changes:
  ```yaml
  paths-ignore:
    - '**.md'
    - 'docs/**'
  ```
- Use `concurrency` to cancel outdated runs

---

## Next Steps

1. **Set up deployment server** - Configure a server with Docker
2. **Configure domain names** - Point DNS to deployment servers
3. **Set up SSL/TLS** - Use Let's Encrypt with Traefik or Nginx
4. **Enable monitoring** - Add health check endpoints
5. **Set up alerts** - Configure GitHub Actions notifications

---

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Codecov Documentation](https://docs.codecov.com)
- [Trivy Security Scanner](https://aquasecurity.github.io/trivy)
