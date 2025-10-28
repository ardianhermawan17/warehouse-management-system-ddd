package http

import (
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/commands"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/queries"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/auth"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/config"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/persistence/sql"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/handlers"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the HTTP router
func SetupRouter(
	cfg *config.Config,
	productRepo product.Repository,
	locationRepo location.Repository,
	stockRepo stock.Repository,
	txManager *sql.TransactionManager,
) *gin.Engine {
	router := gin.Default()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Public routes
	authHandler := handlers.NewAuthHandler(jwtManager)
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Product routes
		productHandler := setupProductHandler(productRepo)
		protected.POST("/products", productHandler.CreateProduct)
		protected.GET("/products", productHandler.ListProducts)
		protected.GET("/products/:id", productHandler.GetProduct)
		protected.PUT("/products/:id", productHandler.UpdateProduct)
		protected.DELETE("/products/:id", productHandler.DeleteProduct)

		// Location routes
		locationHandler := handlers.NewLocationHandler(locationRepo)
		protected.POST("/locations", locationHandler.CreateLocation)
		protected.GET("/locations", locationHandler.ListLocations)
		protected.GET("/locations/:id", locationHandler.GetLocation)
		protected.PUT("/locations/:id", locationHandler.UpdateLocation)
		protected.DELETE("/locations/:id", locationHandler.DeleteLocation)

		// Stock movement routes
		stockHandler := setupStockHandler(productRepo, locationRepo, stockRepo)
		protected.POST("/stock-movements", stockHandler.RecordMovement)
		protected.GET("/stock-movements", stockHandler.ListMovements)
		protected.GET("/stock-movements/:id", stockHandler.GetMovement)
		protected.GET("/stock-movements/product/:product_id", stockHandler.GetProductMovements)
		protected.GET("/stock-movements/location/:location_id", stockHandler.GetLocationMovements)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

// setupProductHandler sets up product handler with all dependencies
func setupProductHandler(productRepo product.Repository) *handlers.ProductHandler {
	createCmd := commands.NewCreateProductCommand(productRepo)
	updateCmd := commands.NewUpdateProductCommand(productRepo)
	listQuery := queries.NewListProductsQuery(productRepo)

	return handlers.NewProductHandler(createCmd, updateCmd, listQuery, productRepo)
}

// setupStockHandler sets up stock handler with all dependencies
func setupStockHandler(
	productRepo product.Repository,
	locationRepo location.Repository,
	stockRepo stock.Repository,
) *handlers.StockHandler {
	stockService := stock.NewService(productRepo, locationRepo, stockRepo)
	recordCmd := commands.NewRecordStockMovementCommand(stockService, productRepo, nil)
	listQuery := queries.NewListStockMovementsQuery(stockRepo)

	return handlers.NewStockHandler(recordCmd, listQuery, stockRepo)
}
