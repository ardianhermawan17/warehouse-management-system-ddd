package product

import "context"

// Repository defines the contract for product persistence
type Repository interface {
	// Create saves a new product
	Create(ctx context.Context, product *Product) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id int64) (*Product, error)

	// GetBySKU retrieves a product by SKU name
	GetBySKU(ctx context.Context, skuName string) (*Product, error)

	// List retrieves all products with pagination
	List(ctx context.Context, limit, offset int) ([]*Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete deletes a product
	Delete(ctx context.Context, id int64) error

	// Count returns total number of products
	Count(ctx context.Context) (int64, error)
}
