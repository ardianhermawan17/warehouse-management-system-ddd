# Testing Guide - WMS API

Comprehensive guide for testing the WMS API application.

## Testing Strategy

The project uses a multi-layered testing approach:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test components working together
3. **E2E Tests**: Test complete workflows through the API

## Running Tests

### Run All Tests
```bash
make test
```

### Run Specific Test Package
```bash
go test -v ./tests/unit/...
go test -v ./tests/e2e/...
```

### Run Specific Test
```bash
go test -v ./tests/unit -run TestStockOutValidation
```

### Run with Coverage
```bash
go test -v -cover ./...
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Unit Tests

### Location: `tests/unit/stock_service_test.go`

Unit tests focus on testing business logic in isolation using mock repositories.

#### Mock Repositories

The test file includes mock implementations of all repositories:

```go
type MockProductRepository struct {
    products map[int64]*product.Product
}

type MockLocationRepository struct {
    locations map[int64]*location.Location
}

type MockStockRepository struct {
    movements map[int64]*stock.StockMovement
}
```

#### Test Cases

##### 1. Stock OUT Validation
```go
func TestStockOutValidation(t *testing.T) {
    // Setup
    ctx := context.Background()
    productRepo := NewMockProductRepository()
    locationRepo := NewMockLocationRepository()
    stockRepo := NewMockStockRepository()

    // Create test data
    prod, _ := product.NewProduct("SKU-001", 100)
    productRepo.Create(ctx, prod)

    loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
    locationRepo.Create(ctx, loc)

    service := stock.NewService(productRepo, locationRepo, stockRepo)

    // Test: Stock OUT cannot exceed available stock
    movement, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeOUT, 150)
    err := service.RecordMovement(ctx, movement)

    if err == nil {
        t.Error("Expected error for insufficient stock, got nil")
    }
}
```

**What it tests:**
- Validates that outbound movements cannot exceed available stock
- Ensures business rule is enforced at domain level

##### 2. Stock IN Capacity Validation
```go
func TestStockInCapacityValidation(t *testing.T) {
    // Setup
    ctx := context.Background()
    productRepo := NewMockProductRepository()
    locationRepo := NewMockLocationRepository()
    stockRepo := NewMockStockRepository()

    // Create test data with limited capacity
    prod, _ := product.NewProduct("SKU-001", 100)
    productRepo.Create(ctx, prod)

    loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 100)
    locationRepo.Create(ctx, loc)

    service := stock.NewService(productRepo, locationRepo, stockRepo)

    // Test: Stock IN cannot exceed location capacity
    movement, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 150)
    err := service.RecordMovement(ctx, movement)

    if err == nil {
        t.Error("Expected error for capacity exceeded, got nil")
    }
}
```

**What it tests:**
- Validates that inbound movements cannot exceed location capacity
- Ensures capacity constraints are enforced

##### 3. Successful Stock Movement
```go
func TestSuccessfulStockMovement(t *testing.T) {
    // Setup
    ctx := context.Background()
    productRepo := NewMockProductRepository()
    locationRepo := NewMockLocationRepository()
    stockRepo := NewMockStockRepository()

    // Create test data
    prod, _ := product.NewProduct("SKU-001", 100)
    productRepo.Create(ctx, prod)

    loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
    locationRepo.Create(ctx, loc)

    service := stock.NewService(productRepo, locationRepo, stockRepo)

    // Test: Successful stock IN
    movement, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
    err := service.RecordMovement(ctx, movement)

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Verify movement was recorded
    movements, _ := stockRepo.GetByProduct(ctx, prod.ID)
    if len(movements) != 1 {
        t.Errorf("Expected 1 movement, got %d", len(movements))
    }
}
```

**What it tests:**
- Validates successful stock movement recording
- Verifies data is persisted correctly

## E2E Tests

### Location: `tests/e2e/api_e2e_test.go`

E2E tests verify complete workflows through the HTTP API.

#### Prerequisites for E2E Tests
- PostgreSQL running
- Database `wms_test` created
- Application can connect to database

#### Test Cases

##### 1. Login Endpoint
```go
func TestLoginEndpoint(t *testing.T) {
    // Setup router
    cfg := &config.Config{
        JWTSecret: "test-secret",
    }
    db, _ := sql.NewDB("postgres://user:password@localhost:5432/wms_test?sslmode=disable")
    router := http.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

    // Test login
    loginReq := map[string]string{
        "username": "testuser",
        "password": "testpass",
    }

    body, _ := json.Marshal(loginReq)
    req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Verify response
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

**What it tests:**
- Login endpoint returns 200 OK
- Response contains JWT token
- Token can be used for authenticated requests

##### 2. Health Endpoint
```go
func TestHealthEndpoint(t *testing.T) {
    // Setup
    router := http.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

    // Test health check
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Verify response
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

**What it tests:**
- Health endpoint is accessible
- Returns correct status code
- No authentication required

## Writing New Tests

### Unit Test Template

```go
func TestNewFeature(t *testing.T) {
    // 1. Setup
    ctx := context.Background()
    repo := NewMockRepository()

    // 2. Create test data
    entity, _ := domain.NewEntity(...)
    repo.Create(ctx, entity)

    // 3. Execute
    result, err := repo.GetByID(ctx, entity.ID)

    // 4. Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if result.ID != entity.ID {
        t.Errorf("Expected ID %d, got %d", entity.ID, result.ID)
    }
}
```

### E2E Test Template

```go
func TestNewEndpoint(t *testing.T) {
    // 1. Setup
    router := http.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

    // 2. Create request
    req := httptest.NewRequest("GET", "/api/v1/endpoint", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    // 3. Execute
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 4. Assert
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## Testing Best Practices

### 1. Test Isolation
- Each test should be independent
- Use fresh mock repositories for each test
- Don't rely on test execution order

### 2. Clear Test Names
```go
// Good
func TestStockOutValidation(t *testing.T) {}

// Bad
func TestStock(t *testing.T) {}
```

### 3. Arrange-Act-Assert Pattern
```go
func TestExample(t *testing.T) {
    // Arrange: Setup test data
    repo := NewMockRepository()
    entity, _ := domain.NewEntity(...)

    // Act: Execute the code
    result, err := repo.Create(context.Background(), entity)

    // Assert: Verify results
    if err != nil {
        t.Error("Unexpected error")
    }
}
```

### 4. Test Error Cases
```go
func TestErrorCase(t *testing.T) {
    repo := NewMockRepository()

    // Test error condition
    _, err := repo.GetByID(context.Background(), 999)

    if err == nil {
        t.Error("Expected error for non-existent ID")
    }
}
```

### 5. Use Table-Driven Tests
```go
func TestMultipleCases(t *testing.T) {
    tests := []struct {
        name    string
        input   int64
        want    error
    }{
        {"valid", 1, nil},
        {"invalid", -1, errors.New("invalid")},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := validate(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Testing Scenarios

### Scenario 1: Complete Stock Movement Workflow

```go
func TestCompleteStockWorkflow(t *testing.T) {
    ctx := context.Background()
    productRepo := NewMockProductRepository()
    locationRepo := NewMockLocationRepository()
    stockRepo := NewMockStockRepository()

    // 1. Create product
    prod, _ := product.NewProduct("SKU-001", 100)
    productRepo.Create(ctx, prod)

    // 2. Create location
    loc, _ := location.NewLocation("LOC-A1", "Warehouse", 500)
    locationRepo.Create(ctx, loc)

    // 3. Record inbound movement
    inbound, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
    service := stock.NewService(productRepo, locationRepo, stockRepo)
    service.RecordMovement(ctx, inbound)

    // 4. Record outbound movement
    outbound, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeOUT, 20)
    service.RecordMovement(ctx, outbound)

    // 5. Verify movements
    movements, _ := stockRepo.GetByProduct(ctx, prod.ID)
    if len(movements) != 2 {
        t.Errorf("Expected 2 movements, got %d", len(movements))
    }
}
```

### Scenario 2: API Integration Test

```go
func TestProductAPIWorkflow(t *testing.T) {
    router := http.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

    // 1. Login
    loginReq := map[string]string{"username": "admin", "password": "pass"}
    body, _ := json.Marshal(loginReq)
    req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    var loginResp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &loginResp)
    token := loginResp["data"].(map[string]interface{})["token"].(string)

    // 2. Create product
    createReq := map[string]interface{}{"sku_name": "SKU-001", "quantity": 100}
    body, _ = json.Marshal(createReq)
    req = httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("Expected 201, got %d", w.Code)
    }

    // 3. List products
    req = httptest.NewRequest("GET", "/api/v1/products", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Debugging Tests

### Run with Verbose Output
```bash
go test -v ./tests/unit/...
```

### Run with Print Statements
```go
func TestDebug(t *testing.T) {
    t.Log("Debug message")
    t.Logf("Formatted: %v", value)
}
```

### Run Single Test
```bash
go test -v ./tests/unit -run TestStockOutValidation
```

### Run with Timeout
```bash
go test -timeout 30s ./...
```

## Coverage Analysis

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./...
```

### View Coverage in Browser
```bash
go tool cover -html=coverage.out
```

### Check Coverage by Package
```bash
go test -cover ./internal/domain/...
go test -cover ./internal/application/...
go test -cover ./internal/infrastructure/...
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: wms_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go test -v ./...
```

## Performance Testing

### Benchmark Tests
```go
func BenchmarkProductCreation(b *testing.B) {
    repo := NewMockProductRepository()
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        prod, _ := product.NewProduct("SKU-001", 100)
        repo.Create(ctx, prod)
    }
}
```

### Run Benchmarks
```bash
go test -bench=. -benchmem ./tests/unit/...
```

## Troubleshooting Tests

### Test Fails with "database connection refused"
- Ensure PostgreSQL is running
- Check DATABASE_DSN in test setup
- Create test database: `createdb wms_test`

### Test Fails with "port already in use"
- Kill process using port: `lsof -ti:8080 | xargs kill -9`
- Or use different port in test

### Mock Repository Not Working
- Ensure mock implements full interface
- Check all methods are implemented
- Verify mock data is initialized

## Test Maintenance

### Keep Tests Updated
- Update tests when business logic changes
- Add tests for new features
- Remove tests for deprecated features

### Refactor Tests
- Extract common setup into helper functions
- Use table-driven tests for multiple cases
- Keep tests DRY (Don't Repeat Yourself)

### Document Complex Tests
```go
// TestComplexScenario verifies that when a product has insufficient stock,
// an outbound movement is rejected with the appropriate error message.
// This ensures the business rule is enforced at the domain level.
func TestComplexScenario(t *testing.T) {
    // ...
}
```

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Testing Best Practices](https://golang.org/doc/effective_go#testing)

---

**Happy Testing!** 🧪
