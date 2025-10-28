#!/bin/bash

# Migration script to run database migrations

set -e

echo "Running database migrations..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Run migrations
go run ./cmd/api

echo "Migrations completed successfully!"
