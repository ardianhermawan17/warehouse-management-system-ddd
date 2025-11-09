package commands

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
)

// RecordStockMovementCommand handles stock movement recording
type RecordStockMovementCommand struct {
	stockService *stock.Service
	productRepo  product.Repository
	txManager    interface{} // TransactionManager
}

// NewRecordStockMovementCommand creates a new record stock movement command
func NewRecordStockMovementCommand(
	stockService *stock.Service,
	productRepo product.Repository,
	txManager interface{},
) *RecordStockMovementCommand {
	return &RecordStockMovementCommand{
		stockService: stockService,
		productRepo:  productRepo,
		txManager:    txManager,
	}
}

// Execute executes the record stock movement command
func (c *RecordStockMovementCommand) Execute(ctx context.Context, req *dto.RecordStockMovementRequest) (*dto.StockMovementResponse, error) {
	// Create stock movement entity
	movementType := stock.MovementType(req.Type)
	movement, err := stock.NewStockMovement(req.ProductID, req.LocationID, movementType, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Record movement with business rule validation
	if err := c.stockService.RecordMovement(ctx, movement); err != nil {
		return nil, err
	}

	// Update product quantity
	prod, err := c.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if movement.IsInbound() {
		if err := prod.IncreaseStock(req.Quantity); err != nil {
			return nil, err
		}
	} else {
		if err := prod.DecreaseStock(req.Quantity); err != nil {
			return nil, err
		}
	}

	if err := c.productRepo.Update(ctx, prod); err != nil {
		return nil, err
	}

	return &dto.StockMovementResponse{
		ID:         movement.ID,
		ProductID:  movement.ProductID,
		LocationID: movement.LocationID,
		Type:       string(movement.Type),
		Quantity:   movement.Quantity,
		CreatedAt:  movement.CreatedAt,
	}, nil
}
