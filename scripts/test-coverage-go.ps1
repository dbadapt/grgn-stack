# Run Go tests with coverage and generate HTML report

Write-Host "Running Go tests with coverage..." -ForegroundColor Cyan
Set-Location backend

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
$lastLine = go tool cover -func=coverage.out | Select-Object -Last 1
Write-Host $lastLine -ForegroundColor Green

Write-Host ""
Write-Host "Coverage report generated: backend/coverage.html" -ForegroundColor Green
Write-Host "To view: Start-Process backend/coverage.html" -ForegroundColor Yellow
