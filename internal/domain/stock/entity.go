package stock

import (
	"errors"
	"time"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeIN  MovementType = "IN"
	MovementTypeOUT MovementType = "OUT"
)

// StockMovement is the aggregate root for stock movement domain
type StockMovement struct {
	ID         int64
	ProductID  int64
	LocationID int64
	Type       MovementType
	Quantity   int64
	CreatedAt  time.Time
}

// NewStockMovement creates a new stock movement
func NewStockMovement(productID, locationID int64, movementType MovementType, quantity int64) (*StockMovement, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if locationID <= 0 {
		return nil, errors.New("invalid location ID")
	}
	if movementType != MovementTypeIN && movementType != MovementTypeOUT {
		return nil, errors.New("invalid movement type")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	return &StockMovement{
		ProductID:  productID,
		LocationID: locationID,
		Type:       movementType,
		Quantity:   quantity,
		CreatedAt:  time.Now(),
	}, nil
}

// IsInbound checks if movement is inbound
func (sm *StockMovement) IsInbound() bool {
	return sm.Type == MovementTypeIN
}

// IsOutbound checks if movement is outbound
func (sm *StockMovement) IsOutbound() bool {
	return sm.Type == MovementTypeOUT
}
