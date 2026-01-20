#!/usr/bin/env bash
# Run Go tests with coverage and generate HTML report

set -e

echo "Running Go tests with coverage..."
cd backend

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
go tool cover -func=coverage.out | tail -1

echo ""
echo "Coverage report generated: backend/coverage.html"
echo "To view: open backend/coverage.html"
