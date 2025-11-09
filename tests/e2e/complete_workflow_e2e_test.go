package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/config"
	httpinterface "github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http"
	"github.com/gin-gonic/gin"
)

// Mock implementations for E2E tests
// These are imported from unit tests but redefined here for clarity

// MockProductRepository is a mock implementation of product.Repository
type MockProductRepository struct {
	products map[int64]*product.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[int64]*product.Product),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, p *product.Product) error {
	p.ID = int64(len(m.products) + 1)
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, product.ErrProductNotFound
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, skuName string) (*product.Product, error) {
	for _, p := range m.products {
		if p.SKUName == skuName {
			return p, nil
		}
	}
	return nil, product.ErrProductNotFound
}

func (m *MockProductRepository) List(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	var result []*product.Product
	for _, p := range m.products {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockProductRepository) Update(ctx context.Context, p *product.Product) error {
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id int64) error {
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.products)), nil
}

// MockLocationRepository is a mock implementation of location.Repository
type MockLocationRepository struct {
	locations map[int64]*location.Location
}

func NewMockLocationRepository() *MockLocationRepository {
	return &MockLocationRepository{
		locations: make(map[int64]*location.Location),
	}
}

func (m *MockLocationRepository) Create(ctx context.Context, l *location.Location) error {
	l.ID = int64(len(m.locations) + 1)
	m.locations[l.ID] = l
	return nil
}

func (m *MockLocationRepository) GetByID(ctx context.Context, id int64) (*location.Location, error) {
	if l, ok := m.locations[id]; ok {
		return l, nil
	}
	return nil, location.ErrLocationNotFound
}

func (m *MockLocationRepository) GetByCode(ctx context.Context, code string) (*location.Location, error) {
	for _, l := range m.locations {
		if l.Code == code {
			return l, nil
		}
	}
	return nil, location.ErrLocationNotFound
}

func (m *MockLocationRepository) List(ctx context.Context, limit, offset int) ([]*location.Location, error) {
	var result []*location.Location
	for _, l := range m.locations {
		result = append(result, l)
	}
	return result, nil
}

func (m *MockLocationRepository) Update(ctx context.Context, l *location.Location) error {
	m.locations[l.ID] = l
	return nil
}

func (m *MockLocationRepository) Delete(ctx context.Context, id int64) error {
	delete(m.locations, id)
	return nil
}

func (m *MockLocationRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.locations)), nil
}

// MockStockRepository is a mock implementation of stock.Repository
type MockStockRepository struct {
	movements map[int64]*stock.StockMovement
}

func NewMockStockRepository() *MockStockRepository {
	return &MockStockRepository{
		movements: make(map[int64]*stock.StockMovement),
	}
}

func (m *MockStockRepository) Create(ctx context.Context, sm *stock.StockMovement) error {
	sm.ID = int64(len(m.movements) + 1)
	m.movements[sm.ID] = sm
	return nil
}

func (m *MockStockRepository) GetByID(ctx context.Context, id int64) (*stock.StockMovement, error) {
	if sm, ok := m.movements[id]; ok {
		return sm, nil
	}
	return nil, stock.ErrMovementNotFound
}

func (m *MockStockRepository) GetByProduct(ctx context.Context, productID int64) ([]*stock.StockMovement, error) {
	var result []*stock.StockMovement
	for _, sm := range m.movements {
		if sm.ProductID == productID {
			result = append(result, sm)
		}
	}
	return result, nil
}

func (m *MockStockRepository) GetByLocation(ctx context.Context, locationID int64) ([]*stock.StockMovement, error) {
	var result []*stock.StockMovement
	for _, sm := range m.movements {
		if sm.LocationID == locationID {
			result = append(result, sm)
		}
	}
	return result, nil
}

func (m *MockStockRepository) List(ctx context.Context, limit, offset int) ([]*stock.StockMovement, error) {
	var result []*stock.StockMovement
	for _, sm := range m.movements {
		result = append(result, sm)
	}
	return result, nil
}

func (m *MockStockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.movements)), nil
}

// MockTransactionManager is a mock implementation of transaction manager
type MockTransactionManager struct{}

func NewMockTransactionManager() *MockTransactionManager {
	return &MockTransactionManager{}
}

func (m *MockTransactionManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// ===================== AUTH E2E TESTS =====================

// TestE2EAuthLoginFlow tests complete authentication login flow
func TestE2EAuthLoginFlow(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

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

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}

	if response["data"] == nil {
		t.Error("Expected token in response data")
	}
}

// TestE2EAuthLoginWithToken tests using token for subsequent requests
func TestE2EAuthLoginWithToken(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	productRepo.Create(ctx, prod)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get token
	token := getTokenE2E(t, router)

	// Use token to access protected endpoint
	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ===================== HEALTH CHECK E2E TEST =====================

// TestE2EHealthCheckNoAuth tests health endpoint without auth
func TestE2EHealthCheckNoAuth(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Error("Expected status to be 'ok'")
	}
}

// ===================== PRODUCT E2E TESTS =====================

// TestE2ECompleteProductWorkflow tests complete product workflow
func TestE2ECompleteProductWorkflow(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// 1. Create Product
	createReq := dto.CreateProductRequest{
		SKUName:  "SKU-PROD-001",
		Quantity: 100,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Step 1: Create - Expected status 201, got %d", w.Code)
	}

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// 2. Get Product
	req = httptest.NewRequest("GET", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 2: Get - Expected status 200, got %d", w.Code)
	}

	// 3. List Products
	req = httptest.NewRequest("GET", "/api/v1/products?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 3: List - Expected status 200, got %d", w.Code)
	}

	// 4. Update Product
	updateReq := dto.UpdateProductRequest{
		SKUName:  "SKU-PROD-001-UPDATED",
		Quantity: 150,
	}

	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PUT", "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 4: Update - Expected status 200, got %d", w.Code)
	}

	// 5. Delete Product
	req = httptest.NewRequest("DELETE", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 5: Delete - Expected status 200, got %d", w.Code)
	}

	// 6. Verify deletion - Get should return 404
	req = httptest.NewRequest("GET", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Step 6: Verify deletion - Expected status 404, got %d", w.Code)
	}
}

// TestE2EProductPagination tests product pagination
func TestE2EProductPagination(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create multiple products
	for i := 1; i <= 5; i++ {
		prod, _ := product.NewProduct("SKU-"+string(rune(i)), int64(100*i))
		productRepo.Create(ctx, prod)
	}

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// List with limit
	req := httptest.NewRequest("GET", "/api/v1/products?limit=2&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}
}

// ===================== LOCATION E2E TESTS =====================

// TestE2ECompleteLocationWorkflow tests complete location workflow
func TestE2ECompleteLocationWorkflow(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// 1. Create Location
	createReq := dto.LocationRequest{
		Code:     "LOC-001",
		Name:     "Main Warehouse",
		Capacity: 1000,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/locations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Step 1: Create - Expected status 201, got %d", w.Code)
	}

	// 2. Get Location
	req = httptest.NewRequest("GET", "/api/v1/locations/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 2: Get - Expected status 200, got %d", w.Code)
	}

	// 3. List Locations
	req = httptest.NewRequest("GET", "/api/v1/locations?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 3: List - Expected status 200, got %d", w.Code)
	}

	// 4. Update Location
	updateReq := dto.LocationRequest{
		Code:     "LOC-001-UPDATED",
		Name:     "Main Warehouse Updated",
		Capacity: 1500,
	}

	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PUT", "/api/v1/locations/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 4: Update - Expected status 200, got %d", w.Code)
	}

	// 5. Delete Location
	req = httptest.NewRequest("DELETE", "/api/v1/locations/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 5: Delete - Expected status 200, got %d", w.Code)
	}
}

// ===================== STOCK MOVEMENT E2E TESTS =====================

// TestE2ECompleteStockMovementWorkflow tests complete stock movement workflow
func TestE2ECompleteStockMovementWorkflow(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// 1. Record Inbound Movement
	inboundReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "IN",
		Quantity:   50,
	}

	body, _ := json.Marshal(inboundReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Step 1: Inbound - Expected status 201, got %d", w.Code)
	}

	// 2. Record Outbound Movement
	outboundReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "OUT",
		Quantity:   20,
	}

	body, _ = json.Marshal(outboundReq)
	req = httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Step 2: Outbound - Expected status 201, got %d", w.Code)
	}

	// 3. Get Stock Movement
	req = httptest.NewRequest("GET", "/api/v1/stock-movements/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 3: Get - Expected status 200, got %d", w.Code)
	}

	// 4. List Stock Movements
	req = httptest.NewRequest("GET", "/api/v1/stock-movements?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 4: List - Expected status 200, got %d", w.Code)
	}

	// 5. Get Product Movements
	req = httptest.NewRequest("GET", "/api/v1/stock-movements/product/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 5: Product Movements - Expected status 200, got %d", w.Code)
	}

	// 6. Get Location Movements
	req = httptest.NewRequest("GET", "/api/v1/stock-movements/location/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Step 6: Location Movements - Expected status 200, got %d", w.Code)
	}
}

// TestE2EStockMovementInsufficientStock tests stock movement with insufficient stock
func TestE2EStockMovementInsufficientStock(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create product with limited quantity
	prod, _ := product.NewProduct("SKU-001", 30)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// Try to move out more than available
	outboundReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "OUT",
		Quantity:   50, // More than available (30)
	}

	body, _ := json.Marshal(outboundReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for insufficient stock, got %d", w.Code)
	}
}

// TestE2EStockMovementCapacityExceeded tests stock movement when capacity exceeded
func TestE2EStockMovementCapacityExceeded(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create location with limited capacity
	prod, _ := product.NewProduct("SKU-001", 200)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 100) // Limited capacity
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// Try to move in more than capacity allows
	inboundReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "IN",
		Quantity:   150, // More than capacity (100)
	}

	body, _ := json.Marshal(inboundReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for capacity exceeded, got %d", w.Code)
	}
}

// TestE2EEndToEndInventoryManagement tests complete inventory management scenario
func TestE2EEndToEndInventoryManagement(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// 1. Create Product
	productReq := dto.CreateProductRequest{
		SKUName:  "LAPTOP-001",
		Quantity: 50,
	}
	body, _ := json.Marshal(productReq)
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create product")
	}

	// 2. Create Locations
	locations := []string{"LOC-MAIN", "LOC-BACKUP"}
	for _, code := range locations {
		locReq := dto.LocationRequest{
			Code:     code,
			Name:     "Warehouse " + code,
			Capacity: 500,
		}
		body, _ := json.Marshal(locReq)
		req := httptest.NewRequest("POST", "/api/v1/locations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create location %s", code)
		}
	}

	// 3. Record inbound movements
	movementReq := dto.RecordStockMovementRequest{
		ProductID:  1,
		LocationID: 1,
		Type:       "IN",
		Quantity:   30,
	}
	body, _ = json.Marshal(movementReq)
	req = httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to record inbound movement")
	}

	// 4. Record outbound movement
	outboundReq := dto.RecordStockMovementRequest{
		ProductID:  1,
		LocationID: 1,
		Type:       "OUT",
		Quantity:   10,
	}
	body, _ = json.Marshal(outboundReq)
	req = httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to record outbound movement")
	}

	// 5. Verify movements
	req = httptest.NewRequest("GET", "/api/v1/stock-movements/product/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get product movements")
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected successful response")
	}
}

// ===================== ERROR HANDLING E2E TESTS =====================

// TestE2EUnauthorizedAccessWithoutToken tests accessing protected endpoint without token
func TestE2EUnauthorizedAccessWithoutToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Try to access protected endpoint without token
	req := httptest.NewRequest("GET", "/api/v1/products", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestE2EUnauthorizedAccessWithInvalidToken tests accessing protected endpoint with invalid token
func TestE2EUnauthorizedAccessWithInvalidToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Try to access protected endpoint with invalid token
	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestE2EInvalidProductID tests accessing product with invalid ID
func TestE2EInvalidProductID(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	token := getTokenE2E(t, router)

	// Try to access product with invalid ID
	req := httptest.NewRequest("GET", "/api/v1/products/invalid-id", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// ===================== HELPER FUNCTIONS =====================

// getTokenE2E retrieves an authentication token for E2E tests
func getTokenE2E(t *testing.T, router *gin.Engine) string {
	loginReq := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	return data["token"].(string)
}
