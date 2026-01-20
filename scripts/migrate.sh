#!/bin/bash
# Run database migrations

set -e

# Build migration tool
echo "Building migration tool..."
cd backend
go build -o ../bin/migrate ./cmd/migrate
cd ..

# Run migrations
echo "Running migrations..."
./bin/migrate -command="${1:-up}"

echo "✓ Done"
