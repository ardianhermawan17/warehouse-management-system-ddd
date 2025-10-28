package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
)

// StockRepository implements stock.Repository
type StockRepository struct {
	db *sql.DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

// Create saves a new stock movement
func (r *StockRepository) Create(ctx context.Context, m *stock.StockMovement) error {
	query := `
		INSERT INTO stock_movements (product_id, location_id, type, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, m.ProductID, m.LocationID, m.Type, m.Quantity).Scan(&m.ID)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	return nil
}

// GetByID retrieves a stock movement by ID
func (r *StockRepository) GetByID(ctx context.Context, id int64) (*stock.StockMovement, error) {
	query := `
		SELECT id, product_id, location_id, type, quantity, created_at
		FROM stock_movements
		WHERE id = $1
	`

	m := &stock.StockMovement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.ProductID, &m.LocationID, &m.Type, &m.Quantity, &m.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, stock.ErrMovementNotFound
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return m, nil
}

// GetByProduct retrieves all movements for a product
func (r *StockRepository) GetByProduct(ctx context.Context, productID int64) ([]*stock.StockMovement, error) {
	query := `
		SELECT id, product_id, location_id, type, quantity, created_at
		FROM stock_movements
		WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}
	defer rows.Close()

	var movements []*stock.StockMovement
	for rows.Next() {
		m := &stock.StockMovement{}
		if err := rows.Scan(&m.ID, &m.ProductID, &m.LocationID, &m.Type, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock movements: %w", err)
	}

	return movements, nil
}

// GetByLocation retrieves all movements for a location
func (r *StockRepository) GetByLocation(ctx context.Context, locationID int64) ([]*stock.StockMovement, error) {
	query := `
		SELECT id, product_id, location_id, type, quantity, created_at
		FROM stock_movements
		WHERE location_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}
	defer rows.Close()

	var movements []*stock.StockMovement
	for rows.Next() {
		m := &stock.StockMovement{}
		if err := rows.Scan(&m.ID, &m.ProductID, &m.LocationID, &m.Type, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock movements: %w", err)
	}

	return movements, nil
}

// List retrieves all stock movements with pagination
func (r *StockRepository) List(ctx context.Context, limit, offset int) ([]*stock.StockMovement, error) {
	query := `
		SELECT id, product_id, location_id, type, quantity, created_at
		FROM stock_movements
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements: %w", err)
	}
	defer rows.Close()

	var movements []*stock.StockMovement
	for rows.Next() {
		m := &stock.StockMovement{}
		if err := rows.Scan(&m.ID, &m.ProductID, &m.LocationID, &m.Type, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock movements: %w", err)
	}

	return movements, nil
}

// Count returns total number of stock movements
func (r *StockRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM stock_movements`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count stock movements: %w", err)
	}

	return count, nil
}
