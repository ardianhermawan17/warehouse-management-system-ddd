package application

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
)

// Ports defines all external dependencies for application layer

// ProductRepository is the port for product persistence
type ProductRepository = product.Repository

// LocationRepository is the port for location persistence
type LocationRepository = location.Repository

// StockRepository is the port for stock movement persistence
type StockRepository = stock.Repository

// TransactionManager defines transaction management contract
type TransactionManager interface {
	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (context.Context, error)

	// CommitTx commits the transaction
	CommitTx(ctx context.Context) error

	// RollbackTx rolls back the transaction
	RollbackTx(ctx context.Context) error

	// WithTx executes a function within a transaction
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// Clock provides current time
type Clock interface {
	Now() interface{}
}
