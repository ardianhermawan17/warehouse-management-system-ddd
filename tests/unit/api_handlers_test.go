package unit

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

// TestAuthLoginSuccess tests successful login
func TestAuthLoginSuccess(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

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

// TestAuthLoginInvalidCredentials tests login with invalid credentials
func TestAuthLoginInvalidCredentials(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Test with empty username (binding validation returns 400)
	loginReq := map[string]string{
		"username": "",
		"password": "testpass",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestAuthLoginMissingFields tests login with missing fields
func TestAuthLoginMissingFields(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Test with missing password
	loginReq := map[string]string{
		"username": "testuser",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestHealthEndpoint tests health check endpoint
func TestHealthCheckEndpoint(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

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

// ===================== PRODUCT TESTS =====================

// TestCreateProductSuccess tests successful product creation
func TestCreateProductSuccess(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Create product
	createReq := dto.CreateProductRequest{
		SKUName:  "SKU-001",
		Quantity: 100,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}
}

// TestCreateProductInvalidRequest tests product creation with invalid request
func TestCreateProductInvalidRequest(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Create product with missing SKU name
	createReq := map[string]interface{}{
		"quantity": 100,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestGetProductSuccess tests successful product retrieval
func TestGetProductSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test product
	prod, _ := product.NewProduct("SKU-001", 100)
	productRepo.Create(ctx, prod)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get product
	req := httptest.NewRequest("GET", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}
}

// TestGetProductNotFound tests getting non-existent product
func TestGetProductNotFound(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get non-existent product
	req := httptest.NewRequest("GET", "/api/v1/products/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestListProductsSuccess tests successful products listing
func TestListProductsSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test products
	prod1, _ := product.NewProduct("SKU-001", 100)
	prod2, _ := product.NewProduct("SKU-002", 200)
	productRepo.Create(ctx, prod1)
	productRepo.Create(ctx, prod2)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// List products
	req := httptest.NewRequest("GET", "/api/v1/products?limit=10&offset=0", nil)
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

// TestUpdateProductSuccess tests successful product update
func TestUpdateProductSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test product
	prod, _ := product.NewProduct("SKU-001", 100)
	productRepo.Create(ctx, prod)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Update product
	updateReq := dto.UpdateProductRequest{
		SKUName:  "SKU-001-UPDATED",
		Quantity: 150,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestDeleteProductSuccess tests successful product deletion
func TestDeleteProductSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test product
	prod, _ := product.NewProduct("SKU-001", 100)
	productRepo.Create(ctx, prod)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Delete product
	req := httptest.NewRequest("DELETE", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ===================== LOCATION TESTS =====================

// TestCreateLocationSuccess tests successful location creation
func TestCreateLocationSuccess(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Create location
	createReq := dto.LocationRequest{
		Code:     "LOC-A1",
		Name:     "Warehouse A",
		Capacity: 500,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/locations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}
}

// TestCreateLocationInvalidCapacity tests location creation with invalid capacity
func TestCreateLocationInvalidCapacity(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Create location with zero capacity
	createReq := dto.LocationRequest{
		Code:     "LOC-A1",
		Name:     "Warehouse A",
		Capacity: 0,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/locations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestGetLocationSuccess tests successful location retrieval
func TestGetLocationSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test location
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get location
	req := httptest.NewRequest("GET", "/api/v1/locations/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestListLocationsSuccess tests successful locations listing
func TestListLocationsSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test locations
	loc1, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	loc2, _ := location.NewLocation("LOC-B1", "Warehouse B", 300)
	locationRepo.Create(ctx, loc1)
	locationRepo.Create(ctx, loc2)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// List locations
	req := httptest.NewRequest("GET", "/api/v1/locations?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestUpdateLocationSuccess tests successful location update
func TestUpdateLocationSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test location
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Update location
	updateReq := dto.LocationRequest{
		Code:     "LOC-A1-UPDATED",
		Name:     "Warehouse A Updated",
		Capacity: 600,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/locations/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestDeleteLocationSuccess tests successful location deletion
func TestDeleteLocationSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test location
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Delete location
	req := httptest.NewRequest("DELETE", "/api/v1/locations/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ===================== STOCK MOVEMENT TESTS =====================

// TestRecordStockMovementInboundSuccess tests successful inbound stock movement
func TestRecordStockMovementInboundSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Record inbound movement
	movementReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "IN",
		Quantity:   50,
	}

	body, _ := json.Marshal(movementReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

// TestRecordStockMovementOutboundSuccess tests successful outbound stock movement
func TestRecordStockMovementOutboundSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Record outbound movement
	movementReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "OUT",
		Quantity:   50,
	}

	body, _ := json.Marshal(movementReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

// TestRecordStockMovementInsufficientStock tests outbound movement with insufficient stock
func TestRecordStockMovementInsufficientStock(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 50)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Record outbound movement with insufficient stock
	movementReq := dto.RecordStockMovementRequest{
		ProductID:  prod.ID,
		LocationID: loc.ID,
		Type:       "OUT",
		Quantity:   100, // More than available
	}

	body, _ := json.Marshal(movementReq)
	req := httptest.NewRequest("POST", "/api/v1/stock-movements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestGetStockMovementSuccess tests successful stock movement retrieval
func TestGetStockMovementSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	movement, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
	stockRepo.Create(ctx, movement)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get stock movement
	req := httptest.NewRequest("GET", "/api/v1/stock-movements/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestListStockMovementsSuccess tests successful stock movements listing
func TestListStockMovementsSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	move1, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
	move2, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeOUT, 20)
	stockRepo.Create(ctx, move1)
	stockRepo.Create(ctx, move2)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// List stock movements
	req := httptest.NewRequest("GET", "/api/v1/stock-movements?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestGetProductMovementsSuccess tests successful product movements retrieval
func TestGetProductMovementsSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	move1, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
	stockRepo.Create(ctx, move1)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get product movements
	req := httptest.NewRequest("GET", "/api/v1/stock-movements/product/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestGetLocationMovementsSuccess tests successful location movements retrieval
func TestGetLocationMovementsSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()
	// txManager removed (not needed for testing)

	// Create test data
	prod, _ := product.NewProduct("SKU-001", 100)
	loc, _ := location.NewLocation("LOC-A1", "Warehouse A", 500)
	productRepo.Create(ctx, prod)
	locationRepo.Create(ctx, loc)

	move1, _ := stock.NewStockMovement(prod.ID, loc.ID, stock.MovementTypeIN, 50)
	stockRepo.Create(ctx, move1)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, nil)

	// Get auth token
	token := getAuthToken(t, router)

	// Get location movements
	req := httptest.NewRequest("GET", "/api/v1/stock-movements/location/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ===================== HELPER FUNCTIONS =====================

// getAuthToken retrieves an authentication token from the login endpoint
func getAuthToken(t *testing.T, router *gin.Engine) string {
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

// ===================== MOCK TRANSACTION MANAGER =====================

// mockTransactionManager is a mock implementation of transaction manager
type mockTransactionManager struct{}

func (m *mockTransactionManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// Import gin package for the mock
var _ = (&gin.Engine{}).Group("")
