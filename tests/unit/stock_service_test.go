package unit

import (
	"context"
	"testing"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
)

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

// Test cases
func TestStockOutValidation(t *testing.T) {
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

func TestStockInCapacityValidation(t *testing.T) {
	ctx := context.Background()
	productRepo := NewMockProductRepository()
	locationRepo := NewMockLocationRepository()
	stockRepo := NewMockStockRepository()

	// Create test data
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

func TestSuccessfulStockMovement(t *testing.T) {
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
