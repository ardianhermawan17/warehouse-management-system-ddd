# Quick Start Guide - WMS API

Get the WMS API up and running in 5 minutes!

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Git

## Step 1: Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd coldstoreindo-ddd-go

# Download dependencies
go mod download
```

## Step 2: Configure Database

### Option A: Local PostgreSQL

```bash
# Create database
createdb wms

# Update .env file
cp internal/infrastructure/config/env.example .env
```

Edit `.env`:
```env
DATABASE_DSN=postgres://user:password@localhost:5432/wms?sslmode=disable
JWT_SECRET=your-secret-key
```

### Option B: Docker PostgreSQL

```bash
docker run --name wms-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=wms \
  -p 5432:5432 \
  -d postgres:15-alpine
```

Update `.env`:
```env
DATABASE_DSN=postgres://postgres:password@localhost:5432/wms?sslmode=disable
```

## Step 3: Run the Application

```bash
# Build and run
make run

# Or run directly
go run ./cmd/api
```

You should see:
```
Starting server on :8080
```

## Step 4: Seed Sample Data (Optional)

In another terminal:
```bash
make seed
```

This creates sample products and locations.

## Step 5: Test the API

### 1. Login to get JWT token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

Save the token for next requests.

### 2. Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "sku_name": "PROD-001",
    "quantity": 100
  }'
```

### 3. Create a Location

```bash
curl -X POST http://localhost:8080/api/v1/locations \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SHELF-A1",
    "name": "Warehouse A - Shelf 1",
    "capacity": 500
  }'
```

### 4. Record Stock Movement

```bash
curl -X POST http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "location_id": 1,
    "type": "IN",
    "quantity": 50
  }'
```

### 5. List Products

```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

## Using Docker Compose

For a complete setup with database:

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

The API will be available at `http://localhost:8080`

## Project Structure Overview

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/           # Business logic (Product, Location, Stock)
│   ├── application/      # Use cases (Commands, Queries)
│   ├── interfaces/       # HTTP handlers and middleware
│   └── infrastructure/   # Database, auth, config
├── tests/                # Unit and E2E tests
├── build/docker/         # Docker configuration
└── scripts/              # Helper scripts
```

## Common Commands

```bash
# Build
make build

# Run
make run

# Test
make test

# Clean
make clean

# Docker
make docker-build
make docker-up
make docker-down
make docker-logs

# Database
make migrate
make seed
```

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login and get JWT token |
| POST | `/products` | Create product |
| GET | `/products` | List products |
| GET | `/products/:id` | Get product |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |
| POST | `/locations` | Create location |
| GET | `/locations` | List locations |
| GET | `/locations/:id` | Get location |
| PUT | `/locations/:id` | Update location |
| DELETE | `/locations/:id` | Delete location |
| POST | `/stock-movements` | Record stock movement |
| GET | `/stock-movements` | List movements |
| GET | `/stock-movements/:id` | Get movement |
| GET | `/stock-movements/product/:id` | Get product movements |
| GET | `/stock-movements/location/:id` | Get location movements |

## Troubleshooting

### Database Connection Error
```
Failed to connect to database
```
- Check PostgreSQL is running
- Verify DATABASE_DSN in .env
- Ensure database exists

### Port Already in Use
```
listen tcp :8080: bind: address already in use
```
- Change SERVER_PORT in .env
- Or kill process: `lsof -ti:8080 | xargs kill -9`

### JWT Token Error
```
invalid or expired token
```
- Ensure token is in `Authorization: Bearer <token>` format
- Get a new token from `/auth/login`

## Next Steps

1. Read [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for complete API reference
2. Read [README.md](README.md) for detailed documentation
3. Check [tests/](tests/) for test examples
4. Explore [internal/](internal/) for code structure

## Support

For issues, check:
- Application logs in terminal
- Database logs: `docker logs wms-postgres`
- API documentation: [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

Happy coding! 🚀
