package stock

import "context"

// Repository defines the contract for stock movement persistence
type Repository interface {
	// Create saves a new stock movement
	Create(ctx context.Context, movement *StockMovement) error

	// GetByID retrieves a stock movement by ID
	GetByID(ctx context.Context, id int64) (*StockMovement, error)

	// GetByProduct retrieves all movements for a product
	GetByProduct(ctx context.Context, productID int64) ([]*StockMovement, error)

	// GetByLocation retrieves all movements for a location
	GetByLocation(ctx context.Context, locationID int64) ([]*StockMovement, error)

	// List retrieves all stock movements with pagination
	List(ctx context.Context, limit, offset int) ([]*StockMovement, error)

	// Count returns total number of stock movements
	Count(ctx context.Context) (int64, error)
}
