package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewDB creates a new database connection
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Create tables
	schema := `
	-- Products table
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		sku_name VARCHAR(255) UNIQUE NOT NULL,
		quantity BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Locations table
	CREATE TABLE IF NOT EXISTS locations (
		id SERIAL PRIMARY KEY,
		code VARCHAR(100) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		capacity BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Stock movements table
	CREATE TABLE IF NOT EXISTS stock_movements (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL REFERENCES products(id),
		location_id INTEGER NOT NULL REFERENCES locations(id),
		type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
		quantity BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes
	CREATE INDEX IF NOT EXISTS idx_products_sku_name ON products(sku_name);
	CREATE INDEX IF NOT EXISTS idx_locations_code ON locations(code);
	CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements(product_id);
	CREATE INDEX IF NOT EXISTS idx_stock_movements_location_id ON stock_movements(location_id);
	CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements(created_at);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
