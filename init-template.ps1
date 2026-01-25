#!/usr/bin/env pwsh
# GRGN Stack Template Initialization Script
# Run this script to set up your new project from the template

param(
    [string]$ProjectName,
    [string]$GitHubUsername,
    [string]$RepoName,
    [switch]$SkipGitInit,
    [switch]$Help
)

$ErrorActionPreference = "Stop"

# Colors for output
function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Cyan }
function Write-Success { param($Message) Write-Host "[OK] $Message" -ForegroundColor Green }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red }

function Show-Help {
    Write-Host @"
GRGN Stack Template Initialization Script

USAGE:
    .\init-template.ps1 [OPTIONS]

OPTIONS:
    -ProjectName      Your project name (e.g., "MyApp")
    -GitHubUsername   Your GitHub username
    -RepoName         Repository name (defaults to project name lowercase)
    -SkipGitInit      Skip git initialization
    -Help             Show this help message

EXAMPLES:
    .\init-template.ps1
    .\init-template.ps1 -ProjectName "MyApp" -GitHubUsername "johndoe"
    .\init-template.ps1 -ProjectName "MyApp" -GitHubUsername "johndoe" -RepoName "my-app"

"@
}

function Get-UserInput {
    param(
        [string]$Prompt,
        [string]$Default = ""
    )

    if ($Default) {
        $userInput = Read-Host "$Prompt [$Default]"
        if ([string]::IsNullOrWhiteSpace($userInput)) { return $Default }
        return $userInput
    } else {
        do {
            $userInput = Read-Host $Prompt
        } while ([string]::IsNullOrWhiteSpace($userInput))
        return $userInput
    }
}

function Update-FileContent {
    param(
        [string]$FilePath,
        [string]$Find,
        [string]$Replace
    )

    if (Test-Path $FilePath) {
        $content = Get-Content $FilePath -Raw -ErrorAction SilentlyContinue
        if ($content) {
            $newContent = $content -replace [regex]::Escape($Find), $Replace
            Set-Content $FilePath $newContent -NoNewline
        }
    }
}

function Update-AllFiles {
    param(
        [string]$Find,
        [string]$Replace,
        [string[]]$Extensions = @("*.md", "*.go", "*.mod", "*.json", "*.yml", "*.yaml", "*.ts", "*.tsx", "*.graphql")
    )

    foreach ($ext in $Extensions) {
        Get-ChildItem -Path . -Recurse -Include $ext -File | ForEach-Object {
            Update-FileContent -FilePath $_.FullName -Find $Find -Replace $Replace
        }
    }
}

# Main script
if ($Help) {
    Show-Help
    exit 0
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Magenta
Write-Host "  GRGN Stack Template Initialization   " -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta
Write-Host ""

# Collect project information
if (-not $ProjectName) {
    $ProjectName = Get-UserInput "Enter your project name (e.g., MyApp)"
}

if (-not $GitHubUsername) {
    $GitHubUsername = Get-UserInput "Enter your GitHub username"
}

if (-not $RepoName) {
    $defaultRepo = $ProjectName.ToLower() -replace '\s+', '-'
    $RepoName = Get-UserInput "Enter repository name" -Default $defaultRepo
}

$EnvPrefix = ($ProjectName.ToUpper() -replace '\s+', '_' -replace '-', '_') + "_"

Write-Host ""
Write-Info "Configuration Summary:"
Write-Host "  Project Name:    $ProjectName"
Write-Host "  GitHub User:     $GitHubUsername"
Write-Host "  Repository:      $RepoName"
Write-Host "  Env Prefix:      $EnvPrefix"
Write-Host ""

$confirm = Read-Host "Proceed with these settings? (Y/n)"
if ($confirm -eq 'n' -or $confirm -eq 'N') {
    Write-Warn "Aborted by user"
    exit 1
}

Write-Host ""

# Step 1: Update Go module paths
Write-Info "Updating Go module paths..."
Update-AllFiles -Find "github.com/yourusername/grgn-stack" -Replace "github.com/$GitHubUsername/$RepoName"
Update-AllFiles -Find "github.com/dbadapt/grgn-stack" -Replace "github.com/$GitHubUsername/$RepoName"
Write-Success "Go module paths updated"

# Step 2: Update environment variable prefix
Write-Info "Updating environment variable prefix..."
Update-AllFiles -Find "GRGN_STACK_" -Replace $EnvPrefix
Write-Success "Environment prefix updated to $EnvPrefix"

# Step 3: Update project name in documentation
Write-Info "Updating project name in documentation..."
Update-AllFiles -Find "GRGN Stack" -Replace $ProjectName
Update-AllFiles -Find "grgn-stack" -Replace $RepoName
Write-Success "Project name updated"

# Step 4: Update README badges and links
Write-Info "Updating README badges..."
Update-FileContent -FilePath "README.md" -Find "YOUR_USERNAME" -Replace $GitHubUsername
Update-FileContent -FilePath "README.md" -Find "YOUR_REPO" -Replace $RepoName
# Also replace template owner's values
Update-FileContent -FilePath "README.md" -Find "dbadapt" -Replace $GitHubUsername
Write-Success "README updated"

# Step 5: Update other documentation
Write-Info "Updating documentation links..."
$docFiles = @(
    "CONTRIBUTING.md",
    "TESTING-CI.md",
    "CI-CD.md",
    "CONFIG.md"
)
foreach ($file in $docFiles) {
    if (Test-Path $file) {
        Update-FileContent -FilePath $file -Find "YOUR_USERNAME" -Replace $GitHubUsername
        Update-FileContent -FilePath $file -Find "YOUR_REPO" -Replace $RepoName
        # Also replace template owner's values
        Update-FileContent -FilePath $file -Find "dbadapt" -Replace $GitHubUsername
    }
}
Write-Success "Documentation updated"

# Step 6: Update LICENSE
Write-Info "Updating LICENSE..."
$currentYear = Get-Date -Format "yyyy"
$licenseOwner = Get-UserInput "Enter name/organization for LICENSE" -Default $GitHubUsername
Update-FileContent -FilePath "LICENSE" -Find "[YOUR NAME OR ORGANIZATION]" -Replace $licenseOwner
Update-FileContent -FilePath "LICENSE" -Find "[YEAR]" -Replace $currentYear
Write-Success "LICENSE updated"

# Step 7: Create .env files
Write-Info "Creating environment files..."
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    # Update prefix in .env
    (Get-Content ".env") -replace "GRGN_STACK_", $EnvPrefix | Set-Content ".env"
    Write-Success "Created .env from .env.example"
}

if (-not (Test-Path "web/.env")) {
    if (Test-Path "web/.env.example") {
        Copy-Item "web/.env.example" "web/.env"
        Write-Success "Created web/.env from web/.env.example"
    }
}

# Step 8: Initialize git (optional)
if (-not $SkipGitInit) {
    Write-Info "Initializing Git repository..."

    # Remove existing .git if present (template origin)
    if (Test-Path ".git") {
        Remove-Item -Recurse -Force ".git"
    }

    git init | Out-Null
    git add .
    git commit -m "Initial commit from GRGN Stack template" | Out-Null

    Write-Success "Git repository initialized"
    Write-Info "To add remote: git remote add origin https://github.com/$GitHubUsername/$RepoName.git"
}

# Step 9: Clean up template files
Write-Info "Cleaning up template-specific files..."

$templateFiles = @(
    "init-template.ps1",
    "init-template.sh",
    "TEMPLATE-SETUP.md",
    "USING-TEMPLATE.md"
)

$removeFiles = Read-Host "Remove template initialization files? (y/N)"
if ($removeFiles -eq 'y' -or $removeFiles -eq 'Y') {
    foreach ($file in $templateFiles) {
        if (Test-Path $file) {
            Remove-Item $file -Force
            Write-Info "Removed $file"
        }
    }

    # Amend the commit to remove template files
    if (-not $SkipGitInit) {
        git add .
        git commit --amend -m "Initial commit from GRGN Stack template" --no-edit | Out-Null
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  Initialization Complete!             " -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Review and edit .env and web/.env files"
Write-Host "  2. Install dependencies: npm install && cd web && npm install"
Write-Host "  3. Start development: docker-compose -f docker-compose.yml -f docker-compose.dev.yml up"
Write-Host ""
Write-Host "Documentation:"
Write-Host "  - README.md         - Project overview"
Write-Host "  - QUICK-REFERENCE.md - Command cheat sheet"
Write-Host "  - ARCHITECTURE.md   - System architecture"
Write-Host ""
Write-Success "Happy building with $ProjectName!"
Write-Host ""
