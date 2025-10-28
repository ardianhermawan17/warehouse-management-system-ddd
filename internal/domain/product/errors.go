package product

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidSKU        = errors.New("invalid SKU name")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrDuplicateSKU      = errors.New("SKU already exists")
)
