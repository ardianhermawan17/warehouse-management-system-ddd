package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/config"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/persistence/sql"
	httpinterface "github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http"
)

// TestLoginEndpoint tests the login endpoint
func TestLoginEndpoint(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	db, err := sql.NewDB("postgres://user:password@localhost:5432/wms_test?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping E2E test: database not available - %v", err)
	}
	defer db.Close()

	productRepo := sql.NewProductRepository(db)
	locationRepo := sql.NewLocationRepository(db)
	stockRepo := sql.NewStockRepository(db)
	txManager := sql.NewTransactionManager(db)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

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

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	db, err := sql.NewDB("postgres://user:password@localhost:5432/wms_test?sslmode=disable")
	if err != nil {
		t.Skipf("Skipping E2E test: database not available - %v", err)
	}
	defer db.Close()

	productRepo := sql.NewProductRepository(db)
	locationRepo := sql.NewLocationRepository(db)
	stockRepo := sql.NewStockRepository(db)
	txManager := sql.NewTransactionManager(db)

	router := httpinterface.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

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
