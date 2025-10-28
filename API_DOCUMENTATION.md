# WMS API - Complete API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All endpoints except `/auth/login` and `/health` require JWT authentication.

### Header Format
```
Authorization: Bearer <jwt_token>
```

### Example
```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." http://localhost:8080/api/v1/products
```

---

## Authentication Endpoints

### 1. Login

**Endpoint:** `POST /auth/login`

**Description:** Authenticate user and receive JWT token

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "invalid credentials"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

---

## Product Endpoints

### 1. Create Product

**Endpoint:** `POST /products`

**Authentication:** Required

**Description:** Create a new product

**Request Body:**
```json
{
  "sku_name": "string (required)",
  "quantity": "integer (required, >= 0)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "product created successfully",
  "data": {
    "id": 1,
    "sku_name": "SKU-001",
    "quantity": 100
  }
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "invalid request"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sku_name": "SKU-001",
    "quantity": 100
  }'
```

---

### 2. Get Product

**Endpoint:** `GET /products/:id`

**Authentication:** Required

**Description:** Retrieve a product by ID

**Path Parameters:**
- `id` (integer, required): Product ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "product retrieved successfully",
  "data": {
    "id": 1,
    "sku_name": "SKU-001",
    "quantity": 100
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "product not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <token>"
```

---

### 3. List Products

**Endpoint:** `GET /products`

**Authentication:** Required

**Description:** List all products with pagination

**Query Parameters:**
- `limit` (integer, optional, default: 10): Number of results per page
- `offset` (integer, optional, default: 0): Number of results to skip

**Response (200 OK):**
```json
{
  "success": true,
  "message": "products retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "sku_name": "SKU-001",
        "quantity": 100
      },
      {
        "id": 2,
        "sku_name": "SKU-002",
        "quantity": 50
      }
    ],
    "total": 2,
    "limit": 10,
    "offset": 0
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/products?limit=10&offset=0" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Update Product

**Endpoint:** `PUT /products/:id`

**Authentication:** Required

**Description:** Update an existing product

**Path Parameters:**
- `id` (integer, required): Product ID

**Request Body:**
```json
{
  "sku_name": "string (optional)",
  "quantity": "integer (optional, >= 0)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "product updated successfully",
  "data": {
    "id": 1,
    "sku_name": "SKU-001-UPDATED",
    "quantity": 150
  }
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sku_name": "SKU-001-UPDATED",
    "quantity": 150
  }'
```

---

### 5. Delete Product

**Endpoint:** `DELETE /products/:id`

**Authentication:** Required

**Description:** Delete a product

**Path Parameters:**
- `id` (integer, required): Product ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "product deleted successfully"
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "product not found"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <token>"
```

---

## Location Endpoints

### 1. Create Location

**Endpoint:** `POST /locations`

**Authentication:** Required

**Description:** Create a new storage location

**Request Body:**
```json
{
  "code": "string (required)",
  "name": "string (required)",
  "capacity": "integer (required, > 0)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "location created successfully",
  "data": {
    "id": 1,
    "code": "LOC-A1",
    "name": "Warehouse A - Shelf 1",
    "capacity": 500
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/locations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "LOC-A1",
    "name": "Warehouse A - Shelf 1",
    "capacity": 500
  }'
```

---

### 2. Get Location

**Endpoint:** `GET /locations/:id`

**Authentication:** Required

**Description:** Retrieve a location by ID

**Path Parameters:**
- `id` (integer, required): Location ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "location retrieved successfully",
  "data": {
    "id": 1,
    "code": "LOC-A1",
    "name": "Warehouse A - Shelf 1",
    "capacity": 500
  }
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/locations/1 \
  -H "Authorization: Bearer <token>"
```

---

### 3. List Locations

**Endpoint:** `GET /locations`

**Authentication:** Required

**Description:** List all storage locations with pagination

**Query Parameters:**
- `limit` (integer, optional, default: 10): Number of results per page
- `offset` (integer, optional, default: 0): Number of results to skip

**Response (200 OK):**
```json
{
  "success": true,
  "message": "locations retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "code": "LOC-A1",
        "name": "Warehouse A - Shelf 1",
        "capacity": 500
      }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/locations?limit=10&offset=0" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Update Location

**Endpoint:** `PUT /locations/:id`

**Authentication:** Required

**Description:** Update an existing location

**Path Parameters:**
- `id` (integer, required): Location ID

**Request Body:**
```json
{
  "code": "string (optional)",
  "name": "string (optional)",
  "capacity": "integer (optional, > 0)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "location updated successfully",
  "data": {
    "id": 1,
    "code": "LOC-A1",
    "name": "Warehouse A - Shelf 1 (Updated)",
    "capacity": 600
  }
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/locations/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 600
  }'
```

---

### 5. Delete Location

**Endpoint:** `DELETE /locations/:id`

**Authentication:** Required

**Description:** Delete a location

**Path Parameters:**
- `id` (integer, required): Location ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "location deleted successfully"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/locations/1 \
  -H "Authorization: Bearer <token>"
```

---

## Stock Movement Endpoints

### 1. Record Stock Movement

**Endpoint:** `POST /stock-movements`

**Authentication:** Required

**Description:** Record a stock movement (IN or OUT)

**Business Rules:**
- Stock OUT cannot exceed available product quantity
- Stock IN cannot exceed location capacity

**Request Body:**
```json
{
  "product_id": "integer (required, > 0)",
  "location_id": "integer (required, > 0)",
  "type": "string (required, 'IN' or 'OUT')",
  "quantity": "integer (required, > 0)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "stock movement recorded successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "location_id": 1,
    "type": "IN",
    "quantity": 50,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response (400 Bad Request - Insufficient Stock):**
```json
{
  "success": false,
  "message": "insufficient stock for outbound movement"
}
```

**Response (400 Bad Request - Capacity Exceeded):**
```json
{
  "success": false,
  "message": "location capacity exceeded"
}
```

**Example - Inbound Movement:**
```bash
curl -X POST http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "location_id": 1,
    "type": "IN",
    "quantity": 50
  }'
```

**Example - Outbound Movement:**
```bash
curl -X POST http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "location_id": 1,
    "type": "OUT",
    "quantity": 20
  }'
```

---

### 2. Get Stock Movement

**Endpoint:** `GET /stock-movements/:id`

**Authentication:** Required

**Description:** Retrieve a stock movement by ID

**Path Parameters:**
- `id` (integer, required): Movement ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "movement retrieved successfully",
  "data": {
    "id": 1,
    "product_id": 1,
    "location_id": 1,
    "type": "IN",
    "quantity": 50,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/stock-movements/1 \
  -H "Authorization: Bearer <token>"
```

---

### 3. List Stock Movements

**Endpoint:** `GET /stock-movements`

**Authentication:** Required

**Description:** List all stock movements with pagination

**Query Parameters:**
- `limit` (integer, optional, default: 10): Number of results per page
- `offset` (integer, optional, default: 0): Number of results to skip

**Response (200 OK):**
```json
{
  "success": true,
  "message": "movements retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "product_id": 1,
        "location_id": 1,
        "type": "IN",
        "quantity": 50,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
  }
}
```

**Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/stock-movements?limit=10&offset=0" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Product Movements

**Endpoint:** `GET /stock-movements/product/:product_id`

**Authentication:** Required

**Description:** Get all stock movements for a specific product

**Path Parameters:**
- `product_id` (integer, required): Product ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "movements retrieved successfully",
  "data": [
    {
      "id": 1,
      "product_id": 1,
      "location_id": 1,
      "type": "IN",
      "quantity": 50,
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "product_id": 1,
      "location_id": 1,
      "type": "OUT",
      "quantity": 20,
      "created_at": "2024-01-15T11:00:00Z"
    }
  ]
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/stock-movements/product/1 \
  -H "Authorization: Bearer <token>"
```

---

### 5. Get Location Movements

**Endpoint:** `GET /stock-movements/location/:location_id`

**Authentication:** Required

**Description:** Get all stock movements for a specific location

**Path Parameters:**
- `location_id` (integer, required): Location ID

**Response (200 OK):**
```json
{
  "success": true,
  "message": "movements retrieved successfully",
  "data": [
    {
      "id": 1,
      "product_id": 1,
      "location_id": 1,
      "type": "IN",
      "quantity": 50,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/stock-movements/location/1 \
  -H "Authorization: Bearer <token>"
```

---

## Health Check Endpoint

### Health Check

**Endpoint:** `GET /health`

**Authentication:** Not required

**Description:** Check API health status

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/health
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data or business rule violation |
| 401 | Unauthorized - Missing or invalid authentication token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Rate Limiting

Currently, there is no rate limiting implemented. Consider adding it for production use.

---

## Pagination

List endpoints support pagination using `limit` and `offset` query parameters:

- `limit`: Number of items per page (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

**Example:**
```bash
# Get 20 products, skip first 40
curl -X GET "http://localhost:8080/api/v1/products?limit=20&offset=40" \
  -H "Authorization: Bearer <token>"
```

---

## Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "message": "operation description",
  "data": {}
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "error description"
}
```

---

## Common Scenarios

### Scenario 1: Complete Stock Movement Workflow

1. **Login to get token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

2. **Create a product:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"sku_name": "SKU-001", "quantity": 100}'
```

3. **Create a location:**
```bash
curl -X POST http://localhost:8080/api/v1/locations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"code": "LOC-A1", "name": "Warehouse A", "capacity": 500}'
```

4. **Record inbound stock movement:**
```bash
curl -X POST http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "location_id": 1, "type": "IN", "quantity": 50}'
```

5. **Record outbound stock movement:**
```bash
curl -X POST http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "location_id": 1, "type": "OUT", "quantity": 20}'
```

6. **View product movements:**
```bash
curl -X GET http://localhost:8080/api/v1/stock-movements/product/1 \
  -H "Authorization: Bearer <token>"
```

---

## Troubleshooting

### 401 Unauthorized
- Ensure you're including the `Authorization: Bearer <token>` header
- Verify the token hasn't expired
- Check that the JWT_SECRET matches between login and protected endpoints

### 400 Bad Request
- Verify all required fields are present in the request body
- Check data types match the specification
- For stock movements, ensure business rules are satisfied

### 404 Not Found
- Verify the resource ID exists
- Check the endpoint path is correct

### 500 Internal Server Error
- Check server logs for detailed error information
- Verify database connection is working
- Ensure all required environment variables are set
