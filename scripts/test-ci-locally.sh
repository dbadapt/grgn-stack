#!/bin/bash
# Test CI checks locally before pushing to GitHub
# This simulates what GitHub Actions will run

set -e
START_DIR=$(pwd)

echo -e "\n========================================"
echo "Testing CI Pipeline Locally"
echo "========================================\n"

FAILURES=()

# ============================================
# Backend Tests
# ============================================
echo -e "📦 Backend Tests (Go)"
echo "--------------------\n"

cd "$(dirname "$0")/../backend"

echo "✓ Running go fmt check..."
FMT_OUTPUT=$(gofmt -s -l .)
if [ -n "$FMT_OUTPUT" ]; then
    echo "❌ Code is not formatted. Run 'go fmt ./...'"
    echo "$FMT_OUTPUT"
    FAILURES+=("go fmt")
else
    echo "  ✓ go fmt passed"
fi

echo -e "\n✓ Running go vet..."
if go vet ./...; then
    echo "  ✓ go vet passed"
else
    echo "❌ go vet failed"
    FAILURES+=("go vet")
fi

echo -e "\n✓ Running go tests..."
if go test -v -race ./...; then
    echo "  ✓ go tests passed"
else
    echo "❌ go tests failed"
    FAILURES+=("go test")
fi

cd "$START_DIR"

# ============================================
# Frontend Tests
# ============================================
echo -e "\n\n📦 Frontend Tests (React/TypeScript)"
echo "------------------------------------\n"

cd "$(dirname "$0")/../web"

echo "✓ Installing dependencies..."
if npm ci --silent; then
    echo "  ✓ dependencies installed"
else
    echo "❌ npm ci failed"
    FAILURES+=("npm ci")
fi

echo -e "\n✓ Running ESLint..."
if npm run lint; then
    echo "  ✓ ESLint passed"
else
    echo "❌ ESLint failed"
    FAILURES+=("eslint")
fi

echo -e "\n✓ Running TypeScript check..."
if npx tsc --noEmit; then
    echo "  ✓ TypeScript check passed"
else
    echo "❌ TypeScript check failed"
    FAILURES+=("typescript")
fi

echo -e "\n✓ Running tests..."
if npm test; then
    echo "  ✓ Frontend tests passed"
else
    echo "❌ Frontend tests failed"
    FAILURES+=("frontend tests")
fi

echo -e "\n✓ Building production bundle..."
if npm run build; then
    echo "  ✓ Build successful"
else
    echo "❌ Build failed"
    FAILURES+=("build")
fi

cd "$START_DIR"

# ============================================
# GraphQL Code Generation
# ============================================
echo -e "\n\n📦 GraphQL Code Generation Check"
echo "--------------------------------\n"

echo "✓ Generating GraphQL code..."
if npm run generate; then
    echo "  ✓ GraphQL code generated"
else
    echo "❌ GraphQL generation failed"
    FAILURES+=("graphql generation")
fi

echo -e "\n✓ Checking for uncommitted changes..."
GIT_STATUS=$(git status --porcelain)
if [ -n "$GIT_STATUS" ]; then
    echo "❌ Generated code has uncommitted changes:"
    echo "$GIT_STATUS"
    echo -e "\nRun 'npm run generate' and commit the changes"
    FAILURES+=("graphql uncommitted")
else
    echo "  ✓ No uncommitted changes"
fi

# ============================================
# Summary
# ============================================
echo -e "\n\n========================================"
echo "Test Summary"
echo "========================================\n"

if [ ${#FAILURES[@]} -eq 0 ]; then
    echo "✅ All checks passed! Ready to push to GitHub."
    echo -e "\nNext steps:"
    echo "  1. git add ."
    echo "  2. git commit -m 'feat: add CI/CD pipeline'"
    echo "  3. git push"
    exit 0
else
    echo "❌ Some checks failed:"
    for failure in "${FAILURES[@]}"; do
        echo "  - $failure"
    done
    echo -e "\nPlease fix the issues before pushing to GitHub."
    exit 1
fi
