package main

import (
	"fmt"
	"log"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/config"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/persistence/sql"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http"
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
	stockRepo := sql.NewStockRepository(db)
	txManager := sql.NewTransactionManager(db)

	// Setup HTTP server
	router := http.SetupRouter(cfg, productRepo, locationRepo, stockRepo, txManager)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
