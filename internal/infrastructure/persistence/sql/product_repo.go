package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
)

// ProductRepository implements product.Repository
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create saves a new product
func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (sku_name, quantity)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, p.SKUName, p.Quantity).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	query := `
		SELECT id, sku_name, quantity
		FROM products
		WHERE id = $1
	`

	p := &product.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.SKUName, &p.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return p, nil
}

// GetBySKU retrieves a product by SKU name
func (r *ProductRepository) GetBySKU(ctx context.Context, skuName string) (*product.Product, error) {
	query := `
		SELECT id, sku_name, quantity
		FROM products
		WHERE sku_name = $1
	`

	p := &product.Product{}
	err := r.db.QueryRowContext(ctx, query, skuName).Scan(&p.ID, &p.SKUName, &p.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return p, nil
}

// List retrieves all products with pagination
func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	query := `
		SELECT id, sku_name, quantity
		FROM products
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		p := &product.Product{}
		if err := rows.Scan(&p.ID, &p.SKUName, &p.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	query := `
		UPDATE products
		SET sku_name = $1, quantity = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, p.SKUName, p.Quantity, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

// Count returns total number of products
func (r *ProductRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM products`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}
