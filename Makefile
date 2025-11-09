.PHONY: help build run test clean docker-build docker-up docker-down migrate seed

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make migrate        - Run database migrations"
	@echo "  make seed           - Seed database with sample data"

build:
	@echo "Building application..."
	go build -o bin/wms-api ./cmd/api

run: build
	@echo "Running application..."
	./bin/wms-api

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

docker-build:
	@echo "Building Docker image..."
	docker build -f build/docker/Dockerfile -t wms-api:latest .

docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f build/docker/docker-compose.yml up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f build/docker/docker-compose.yml down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose -f build/docker/docker-compose.yml logs -f

migrate:
	@echo "Running migrations..."
	go run ./cmd/api

seed:
	@echo "Seeding database..."
	go run ./scripts/seed.go

.DEFAULT_GOAL := help
