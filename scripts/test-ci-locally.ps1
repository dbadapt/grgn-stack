#!/usr/bin/env pwsh
# Test CI checks locally before pushing to GitHub
# This simulates what GitHub Actions will run

$ErrorActionPreference = "Continue"
$startLocation = Get-Location

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Testing CI Pipeline Locally" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$failures = @()

# ============================================
# Backend Tests
# ============================================
Write-Host "📦 Backend Tests (Go)" -ForegroundColor Yellow
Write-Host "--------------------`n" -ForegroundColor Yellow

try {
    Set-Location "$PSScriptRoot\..\backend"

    Write-Host "✓ Running go fmt check..." -ForegroundColor Gray
    $fmtOutput = go fmt ./...
    if ($fmtOutput) {
        Write-Host "❌ Code is not formatted. Run 'go fmt ./...'" -ForegroundColor Red
        $failures += "go fmt"
    } else {
        Write-Host "  ✓ go fmt passed" -ForegroundColor Green
    }

    Write-Host "`n✓ Running go vet..." -ForegroundColor Gray
    go vet ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ go vet failed" -ForegroundColor Red
        $failures += "go vet"
    } else {
        Write-Host "  ✓ go vet passed" -ForegroundColor Green
    }

    Write-Host "`n✓ Running go tests..." -ForegroundColor Gray
    go test -v -race ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ go tests failed" -ForegroundColor Red
        $failures += "go test"
    } else {
        Write-Host "  ✓ go tests passed" -ForegroundColor Green
    }

} catch {
    Write-Host "❌ Backend tests failed: $_" -ForegroundColor Red
    $failures += "backend"
}

Set-Location $startLocation

# ============================================
# Frontend Tests
# ============================================
Write-Host "`n`n📦 Frontend Tests (React/TypeScript)" -ForegroundColor Yellow
Write-Host "------------------------------------`n" -ForegroundColor Yellow

try {
    Set-Location "$PSScriptRoot\..\web"

    Write-Host "✓ Installing dependencies..." -ForegroundColor Gray
    npm ci --silent
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ npm ci failed" -ForegroundColor Red
        $failures += "npm ci"
    } else {
        Write-Host "  ✓ dependencies installed" -ForegroundColor Green
    }

    Write-Host "`n✓ Running ESLint..." -ForegroundColor Gray
    npm run lint
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ ESLint failed" -ForegroundColor Red
        $failures += "eslint"
    } else {
        Write-Host "  ✓ ESLint passed" -ForegroundColor Green
    }

    Write-Host "`n✓ Running TypeScript check..." -ForegroundColor Gray
    npx tsc --noEmit
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ TypeScript check failed" -ForegroundColor Red
        $failures += "typescript"
    } else {
        Write-Host "  ✓ TypeScript check passed" -ForegroundColor Green
    }

    Write-Host "`n✓ Running tests..." -ForegroundColor Gray
    npm test
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Frontend tests failed" -ForegroundColor Red
        $failures += "frontend tests"
    } else {
        Write-Host "  ✓ Frontend tests passed" -ForegroundColor Green
    }

    Write-Host "`n✓ Building production bundle..." -ForegroundColor Gray
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Build failed" -ForegroundColor Red
        $failures += "build"
    } else {
        Write-Host "  ✓ Build successful" -ForegroundColor Green
    }

} catch {
    Write-Host "❌ Frontend tests failed: $_" -ForegroundColor Red
    $failures += "frontend"
}

Set-Location $startLocation

# ============================================
# GraphQL Code Generation
# ============================================
Write-Host "`n`n📦 GraphQL Code Generation Check" -ForegroundColor Yellow
Write-Host "--------------------------------`n" -ForegroundColor Yellow

try {
    Write-Host "✓ Generating GraphQL code..." -ForegroundColor Gray
    npm run generate
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ GraphQL generation failed" -ForegroundColor Red
        $failures += "graphql generation"
    } else {
        Write-Host "  ✓ GraphQL code generated" -ForegroundColor Green
    }

    Write-Host "`n✓ Checking for uncommitted changes..." -ForegroundColor Gray
    $gitStatus = git status --porcelain
    if ($gitStatus) {
        Write-Host "❌ Generated code has uncommitted changes:" -ForegroundColor Red
        Write-Host $gitStatus -ForegroundColor Yellow
        Write-Host "`nRun 'npm run generate' and commit the changes" -ForegroundColor Yellow
        $failures += "graphql uncommitted"
    } else {
        Write-Host "  ✓ No uncommitted changes" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ GraphQL check failed: $_" -ForegroundColor Red
    $failures += "graphql"
}

# ============================================
# Summary
# ============================================
Write-Host "`n`n========================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

if ($failures.Count -eq 0) {
    Write-Host "✅ All checks passed! Ready to push to GitHub." -ForegroundColor Green
    Write-Host "`nNext steps:" -ForegroundColor Cyan
    Write-Host "  1. git add ." -ForegroundColor Gray
    Write-Host "  2. git commit -m 'feat: add CI/CD pipeline'" -ForegroundColor Gray
    Write-Host "  3. git push" -ForegroundColor Gray
    exit 0
} else {
    Write-Host "❌ Some checks failed:" -ForegroundColor Red
    foreach ($failure in $failures) {
        Write-Host "  - $failure" -ForegroundColor Red
    }
    Write-Host "`nPlease fix the issues before pushing to GitHub." -ForegroundColor Yellow
    exit 1
}
