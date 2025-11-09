# WMS API - Warehouse Management System

A production-ready REST API for Warehouse Management System (WMS) built with Go, PostgreSQL, and Domain-Driven Design (DDD) principles.

## Cara running tanpa makefile

``bash
docker-compose -f build/docker/docker-compose.yml up -d
``
## Features

- **Domain-Driven Design**: Clean architecture with clear separation of concerns
- **Product Management**: Create, read, update, and delete products with SKU tracking
- **Location Management**: Manage storage locations with capacity constraints
- **Stock Movements**: Track inbound and outbound stock movements with audit trail
- **Business Rules Enforcement**:
  - Stock OUT cannot exceed available stock
  - Stock IN cannot exceed location capacity
  - Automatic product quantity updates on stock movements
- **JWT Authentication**: Secure endpoints with JWT tokens
- **PostgreSQL Database**: Reliable data persistence with manual SQL queries
- **Docker Ready**: Complete Docker setup for containerized deployment
- **Comprehensive Testing**: Unit and E2E tests included

## Project Structure

```
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Domain layer (pure business logic)
│   │   ├── product/
│   │   ├── location/
│   │   └── stock/
│   ├── application/                # Application layer (use cases)
│   │   ├── commands/
│   │   ├── queries/
│   │   └── dto/
│   ├── interfaces/                 # Interface layer (HTTP handlers)
│   │   └── http/
│   │       ├── handlers/
│   │       ├── middleware/
│   │       └── response/
│   └── infrastructure/             # Infrastructure layer (persistence, auth, config)
│       ├── persistence/sql/
│       ├── auth/
│       ├── config/
│       └── logging/
├── pkg/                            # Utility packages
├── scripts/                        # Helper scripts
├── build/docker/                   # Docker configuration
├── tests/                          # Test files
├── Makefile                        # Build automation
├── go.mod                          # Go module definition
└── README.md                       # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (optional, for containerized deployment)

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd coldstoreindo-ddd-go
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Copy the example environment file and update with your settings:

```bash
cp internal/infrastructure/config/env.example .env
```

Edit `.env` with your database credentials:

```env
SERVER_PORT=8080
ENVIRONMENT=development
DATABASE_DSN=postgres://user:password@localhost:5432/wms?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
```

### 4. Setup PostgreSQL

Create a PostgreSQL database:

```bash
createdb wms
```

Or use Docker:

```bash
docker run --name wms-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=wms -p 5432:5432 -d postgres:15-alpine
```

## Running the Application

### Local Development

```bash
# Build and run
make run

# Or run directly
go run ./cmd/api
```

The API will start on `http://localhost:8080`

### Using Docker Compose

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

## Database Migrations

Migrations run automatically on application startup. To manually run migrations:

```bash
make migrate
```

## Seeding Sample Data

```bash
make seed
```

This will create sample products and locations in the database.

## API Endpoints

### Authentication

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user",
  "password": "password"
}

Response:
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Products

All product endpoints require JWT authentication (add `Authorization: Bearer <token>` header)

#### Create Product
```
POST /api/v1/products
{
  "sku_name": "SKU-001",
  "quantity": 100
}
```

#### Get Product
```
GET /api/v1/products/:id
```

#### List Products
```
GET /api/v1/products?limit=10&offset=0
```

#### Update Product
```
PUT /api/v1/products/:id
{
  "sku_name": "SKU-001",
  "quantity": 150
}
```

#### Delete Product
```
DELETE /api/v1/products/:id
```

### Locations

#### Create Location
```
POST /api/v1/locations
{
  "code": "LOC-A1",
  "name": "Warehouse A - Shelf 1",
  "capacity": 500
}
```

#### Get Location
```
GET /api/v1/locations/:id
```

#### List Locations
```
GET /api/v1/locations?limit=10&offset=0
```

#### Update Location
```
PUT /api/v1/locations/:id
{
  "code": "LOC-A1",
  "name": "Warehouse A - Shelf 1",
  "capacity": 600
}
```

#### Delete Location
```
DELETE /api/v1/locations/:id
```

### Stock Movements

#### Record Stock Movement
```
POST /api/v1/stock-movements
{
  "product_id": 1,
  "location_id": 1,
  "type": "IN",
  "quantity": 50
}
```

Type can be "IN" (inbound) or "OUT" (outbound)

#### Get Movement
```
GET /api/v1/stock-movements/:id
```

#### List Movements
```
GET /api/v1/stock-movements?limit=10&offset=0
```

#### Get Product Movements
```
GET /api/v1/stock-movements/product/:product_id
```

#### Get Location Movements
```
GET /api/v1/stock-movements/location/:location_id
```

## Testing

### Run all tests
```bash
make test
```

### Run unit tests
```bash
go test -v ./tests/unit/...
```

### Run E2E tests
```bash
go test -v ./tests/e2e/...
```

## Database Schema

### Products Table
```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  sku_name VARCHAR(255) UNIQUE NOT NULL,
  quantity BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Locations Table
```sql
CREATE TABLE locations (
  id SERIAL PRIMARY KEY,
  code VARCHAR(100) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  capacity BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Stock Movements Table
```sql
CREATE TABLE stock_movements (
  id SERIAL PRIMARY KEY,
  product_id INTEGER NOT NULL REFERENCES products(id),
  location_id INTEGER NOT NULL REFERENCES locations(id),
  type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
  quantity BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Business Rules

1. **Stock OUT Validation**: Cannot record outbound movement if quantity exceeds available stock
2. **Stock IN Validation**: Cannot record inbound movement if total stock at location would exceed capacity
3. **Automatic Updates**: Product quantity automatically updates when stock movement is recorded
4. **Audit Trail**: All stock movements are logged with timestamp for audit purposes

## Architecture Highlights

### Domain Layer
- Pure business logic without external dependencies
- Entities: Product, Location, StockMovement
- Repositories: Interfaces for data persistence
- Services: Domain services for complex business rules

### Application Layer
- Commands: Handle state-changing operations (Create, Update, Record)
- Queries: Handle read operations (List)
- DTOs: Data transfer objects for API communication
- Ports: Abstractions for infrastructure dependencies

### Interface Layer
- HTTP Handlers: REST endpoint implementations
- Middleware: Authentication, logging, recovery
- Response: Standardized API response format

### Infrastructure Layer
- SQL Repositories: Manual SQL query implementations
- JWT Manager: Token generation and verification
- Configuration: Environment-based configuration
- Transaction Manager: Database transaction handling

## Error Handling

The API returns standardized error responses:

```json
{
  "success": false,
  "message": "error description"
}
```

Common HTTP status codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Security Considerations

1. **JWT Tokens**: All protected endpoints require valid JWT tokens
2. **Environment Variables**: Sensitive data stored in `.env` file (not committed)
3. **SQL Injection Prevention**: Using parameterized queries
4. **CORS**: Can be added via middleware if needed
5. **Rate Limiting**: Can be added via middleware if needed

## Performance Optimization

1. **Database Indexes**: Indexes on frequently queried columns
2. **Connection Pooling**: Configured database connection pool
3. **Pagination**: List endpoints support limit/offset pagination
4. **Query Optimization**: Efficient SQL queries with proper joins

## Deployment

### Docker Deployment

```bash
# Build and push image
docker build -f build/docker/Dockerfile -t wms-api:latest .

# Run with docker-compose
docker-compose -f build/docker/docker-compose.yml up -d
```

### Environment Variables for Production

```env
SERVER_PORT=8080
ENVIRONMENT=production
DATABASE_DSN=postgres://user:password@db-host:5432/wms?sslmode=require
JWT_SECRET=<strong-random-secret>
```

## Troubleshooting

### Database Connection Error
- Verify PostgreSQL is running
- Check DATABASE_DSN in .env file
- Ensure database exists

### JWT Token Error
- Verify JWT_SECRET is set in .env
- Check token format: `Authorization: Bearer <token>`
- Ensure token hasn't expired

### Port Already in Use
- Change SERVER_PORT in .env
- Or kill process using port 8080

## Contributing

1. Follow the existing code structure and patterns
2. Write tests for new features
3. Ensure all tests pass before submitting PR
4. Follow Go code style guidelines

## License

MIT License - See LICENSE file for details

## Support

For issues and questions, please create an issue in the repository.
