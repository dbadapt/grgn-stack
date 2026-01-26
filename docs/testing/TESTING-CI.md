# Testing CI/CD Pipeline

## Overview

This guide explains how to test your CI/CD pipeline both locally and on GitHub.

## Local Testing

### Quick Test (Windows)

Run all CI checks locally before pushing:

```bash
npm run test:ci
```

Or directly:

```powershell
pwsh ./scripts/test-ci-locally.ps1
```

### Quick Test (Linux/Mac)

```bash
bash ./scripts/test-ci-locally.sh
```

### What Gets Tested

The local test script simulates what GitHub Actions will run:

1. **Backend (Go)**
   - `go fmt` formatting check
   - `go vet` static analysis
   - `go test -v -race` with race detection

2. **Frontend (TypeScript/React)**
   - ESLint code quality
   - TypeScript type checking
   - Vitest unit tests
   - Production build

3. **GraphQL**
   - Code generation
   - Uncommitted changes check

### Expected Output

```
========================================
Testing CI Pipeline Locally
========================================

📦 Backend Tests (Go)
--------------------
✓ Running go fmt check...
  ✓ go fmt passed
✓ Running go vet...
  ✓ go vet passed
✓ Running go tests...
  ✓ go tests passed

📦 Frontend Tests (React/TypeScript)
------------------------------------
✓ Running ESLint...
  ✓ ESLint passed
✓ Running TypeScript check...
  ✓ TypeScript check passed
✓ Running tests...
  ✓ Frontend tests passed
✓ Building production bundle...
  ✓ Build successful

📦 GraphQL Code Generation Check
--------------------------------
✓ Generating GraphQL code...
  ✓ GraphQL code generated
✓ Checking for uncommitted changes...
  ✓ No uncommitted changes

========================================
Test Summary
========================================

✅ All checks passed! Ready to push to GitHub.
```

## Testing on GitHub

### Step 1: Commit and Push

```bash
# Add all new files
git add .

# Commit with conventional commit message
git commit -m "feat: add CI/CD pipeline with GitHub Actions"

# Push to GitHub
git push origin develop
```

### Step 2: View Workflow Runs

1. Go to your GitHub repository
2. Click **Actions** tab
3. You'll see workflows running:
   - ✅ **CI** - Main test suite
   - ✅ **Coverage** - Code coverage tracking
   - ✅ **Docker Build** - Container image builds

### Step 3: Check Results

Click on any workflow run to see:

- Job summaries
- Test results
- Build logs
- Coverage reports

## Testing Individual Workflows

### Test CI Workflow Only

```bash
# Create a test branch
git checkout -b test/ci-pipeline

# Make a small change
echo "# Test" >> README.md

# Commit and push
git add README.md
git commit -m "test: verify CI workflow"
git push origin test/ci-pipeline

# Create a pull request to trigger CI
```

### Test Coverage Workflow

Coverage workflow runs automatically with CI. Check results at:

- **GitHub Actions** - Summary in workflow output
- **Codecov** (if configured) - https://codecov.io/gh/YOUR_USERNAME/grgn-stack

### Test Docker Build Workflow

Automatically runs when you push to `main`, `staging`, or `develop`:

```bash
# Push to develop branch
git push origin develop
```

To test manually:

1. Go to **Actions** tab
2. Select **Docker Build** workflow
3. Click **Run workflow**
4. Choose branch
5. Click **Run workflow** button

### Test Deploy Workflow

**Manual Deployment:**

1. Go to **Actions** tab
2. Select **Deploy** workflow
3. Click **Run workflow**
4. Choose environment (development/staging/production)
5. Click **Run workflow** button

**Automatic Deployment:**

- Push to `develop` → Deploys to development
- Push to `staging` → Deploys to staging
- Push to `main` → Deploys to production

## Common Issues & Fixes

### ❌ Tests Fail Locally

**Backend tests fail:**

```bash
cd backend
go test -v ./...
# Review and fix failing tests
```

**Frontend tests fail:**

```bash
cd web
npm test
# Review and fix failing tests
```

### ❌ GraphQL Code Out of Date

```bash
npm run generate
git add services/**/controller/generated services/**/view/web
git commit -m "chore: update generated GraphQL code"
```

### ❌ Linting Errors

**Auto-fix:**

```bash
cd web
npm run lint:fix
```

**Manual fix:**
Review ESLint output and fix issues manually.

### ❌ TypeScript Errors

```bash
cd web
npx tsc --noEmit
# Fix type errors shown in output
```

### ❌ Go Formatting

```bash
cd backend
go fmt ./...
git add .
git commit -m "style: format Go code"
```

## Workflow Configuration

### Required Secrets (for deployment)

Set in **Settings → Secrets and variables → Actions**:

```yaml
CODECOV_TOKEN         # Optional - for coverage tracking
DOCKER_REGISTRY       # Your Docker registry URL
DOCKER_USERNAME       # Docker registry username
DOCKER_PASSWORD       # Docker registry password/token
DEPLOY_HOST          # Deployment server hostname
DEPLOY_USER          # SSH username
DEPLOY_SSH_KEY       # SSH private key
```

### Environment Variables

Set per environment in **Settings → Environments**:

- `APP_URL` - Application URL for each environment

## Debugging Failed Workflows

### View Detailed Logs

1. Click on failed workflow run
2. Click on failed job (red X)
3. Expand failed step
4. Review error messages

### Re-run Failed Workflows

1. Open failed workflow run
2. Click **Re-run jobs** dropdown
3. Choose:
   - **Re-run all jobs** - Run everything again
   - **Re-run failed jobs** - Only retry failures

### Download Artifacts

Some workflows save artifacts (coverage reports, logs):

1. Scroll to bottom of workflow run
2. Click **Artifacts** section
3. Download files for local inspection

## Best Practices

### Before Pushing

✅ **Always run local tests:**

```bash
npm run test:ci
```

✅ **Check git status:**

```bash
git status
```

✅ **Review changes:**

```bash
git diff
```

### During Development

✅ **Use feature branches:**

```bash
git checkout -b feature/my-new-feature
```

✅ **Make small, focused commits:**

```bash
git commit -m "feat: add user authentication"
```

✅ **Create pull requests** for code review before merging to `main`

### After Pushing

✅ **Monitor GitHub Actions** - Watch for failures

✅ **Check coverage reports** - Ensure coverage doesn't drop

✅ **Review security alerts** - Address Trivy scan results

## Quick Reference

| Command            | Description                  |
| ------------------ | ---------------------------- |
| `npm run test:ci`  | Run all CI checks locally    |
| `npm test`         | Run backend + frontend tests |
| `npm run lint`     | Lint frontend code           |
| `npm run generate` | Generate GraphQL code        |
| `npm run coverage` | Generate coverage reports    |

## Next Steps

1. ✅ Test locally with `npm run test:ci`
2. ✅ Commit changes: `git add . && git commit -m "feat: add CI/CD"`
3. ✅ Push to GitHub: `git push`
4. ✅ Monitor Actions tab
5. ✅ Configure secrets for deployment
6. ✅ Set up branch protection rules

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Debugging Workflows](https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/about-monitoring-and-troubleshooting)
- [CI-CD.md](./CI-CD.md) - Complete pipeline documentation
