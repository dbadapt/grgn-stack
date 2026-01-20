# Run database migrations (PowerShell)

param(
    [string]$Command = "up"
)

$ErrorActionPreference = "Stop"

# Build migration tool
Write-Host "Building migration tool..." -ForegroundColor Cyan
Push-Location backend
go build -o ..\bin\migrate.exe .\cmd\migrate
Pop-Location

# Run migrations
Write-Host "Running migrations..." -ForegroundColor Cyan
& ".\bin\migrate.exe" -command=$Command

Write-Host "✓ Done" -ForegroundColor Green
