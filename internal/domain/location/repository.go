package location

import "context"

// Repository defines the contract for location persistence
type Repository interface {
	// Create saves a new location
	Create(ctx context.Context, location *Location) error

	// GetByID retrieves a location by ID
	GetByID(ctx context.Context, id int64) (*Location, error)

	// GetByCode retrieves a location by code
	GetByCode(ctx context.Context, code string) (*Location, error)

	// List retrieves all locations with pagination
	List(ctx context.Context, limit, offset int) ([]*Location, error)

	// Update updates an existing location
	Update(ctx context.Context, location *Location) error

	// Delete deletes a location
	Delete(ctx context.Context, id int64) error

	// Count returns total number of locations
	Count(ctx context.Context) (int64, error)
}
