#!/bin/bash

# Development script to run the application locally

set -e

echo "Starting WMS API in development mode..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Run the application
go run ./cmd/api
