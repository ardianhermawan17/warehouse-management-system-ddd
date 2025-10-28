package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/config"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/persistence/sql"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := sql.NewDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := sql.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	productRepo := sql.NewProductRepository(db)
	locationRepo := sql.NewLocationRepository(db)

	ctx := context.Background()

	// Seed products
	products := []struct {
		sku      string
		quantity int64
	}{
		{"SKU-001", 100},
		{"SKU-002", 50},
		{"SKU-003", 200},
		{"SKU-004", 75},
		{"SKU-005", 150},
	}

	fmt.Println("Seeding products...")
	for _, p := range products {
		prod, err := product.NewProduct(p.sku, p.quantity)
		if err != nil {
			log.Printf("Failed to create product %s: %v", p.sku, err)
			continue
		}

		if err := productRepo.Create(ctx, prod); err != nil {
			log.Printf("Failed to save product %s: %v", p.sku, err)
			continue
		}

		fmt.Printf("Created product: %s (qty: %d)\n", p.sku, p.quantity)
	}

	// Seed locations
	locations := []struct {
		code     string
		name     string
		capacity int64
	}{
		{"LOC-A1", "Warehouse A - Shelf 1", 500},
		{"LOC-A2", "Warehouse A - Shelf 2", 500},
		{"LOC-B1", "Warehouse B - Shelf 1", 1000},
		{"LOC-B2", "Warehouse B - Shelf 2", 1000},
		{"LOC-C1", "Cold Storage - Zone 1", 300},
	}

	fmt.Println("\nSeeding locations...")
	for _, l := range locations {
		loc, err := location.NewLocation(l.code, l.name, l.capacity)
		if err != nil {
			log.Printf("Failed to create location %s: %v", l.code, err)
			continue
		}

		if err := locationRepo.Create(ctx, loc); err != nil {
			log.Printf("Failed to save location %s: %v", l.code, err)
			continue
		}

		fmt.Printf("Created location: %s - %s (capacity: %d)\n", l.code, l.name, l.capacity)
	}

	fmt.Println("\nDatabase seeding completed successfully!")
}
