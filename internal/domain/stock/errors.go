package stock

import "errors"

var (
	ErrMovementNotFound    = errors.New("stock movement not found")
	ErrInvalidProductID    = errors.New("invalid product ID")
	ErrInvalidLocationID   = errors.New("invalid location ID")
	ErrInvalidMovementType = errors.New("invalid movement type")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInsufficientStock   = errors.New("insufficient stock for outbound movement")
	ErrCapacityExceeded    = errors.New("location capacity exceeded for inbound movement")
)
