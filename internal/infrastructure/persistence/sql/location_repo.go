package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
)

// LocationRepository implements location.Repository
type LocationRepository struct {
	db *sql.DB
}

// NewLocationRepository creates a new location repository
func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// Create saves a new location
func (r *LocationRepository) Create(ctx context.Context, l *location.Location) error {
	query := `
		INSERT INTO locations (code, name, capacity)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, l.Code, l.Name, l.Capacity).Scan(&l.ID)
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

// GetByID retrieves a location by ID
func (r *LocationRepository) GetByID(ctx context.Context, id int64) (*location.Location, error) {
	query := `
		SELECT id, code, name, capacity
		FROM locations
		WHERE id = $1
	`

	l := &location.Location{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&l.ID, &l.Code, &l.Name, &l.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, location.ErrLocationNotFound
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return l, nil
}

// GetByCode retrieves a location by code
func (r *LocationRepository) GetByCode(ctx context.Context, code string) (*location.Location, error) {
	query := `
		SELECT id, code, name, capacity
		FROM locations
		WHERE code = $1
	`

	l := &location.Location{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(&l.ID, &l.Code, &l.Name, &l.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, location.ErrLocationNotFound
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return l, nil
}

// List retrieves all locations with pagination
func (r *LocationRepository) List(ctx context.Context, limit, offset int) ([]*location.Location, error) {
	query := `
		SELECT id, code, name, capacity
		FROM locations
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}
	defer rows.Close()

	var locations []*location.Location
	for rows.Next() {
		l := &location.Location{}
		if err := rows.Scan(&l.ID, &l.Code, &l.Name, &l.Capacity); err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating locations: %w", err)
	}

	return locations, nil
}

// Update updates an existing location
func (r *LocationRepository) Update(ctx context.Context, l *location.Location) error {
	query := `
		UPDATE locations
		SET code = $1, name = $2, capacity = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, l.Code, l.Name, l.Capacity, l.ID)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return location.ErrLocationNotFound
	}

	return nil
}

// Delete deletes a location
func (r *LocationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM locations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return location.ErrLocationNotFound
	}

	return nil
}

// Count returns total number of locations
func (r *LocationRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM locations`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count locations: %w", err)
	}

	return count, nil
}
