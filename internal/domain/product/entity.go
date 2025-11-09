package product

import "errors"

// Product is the aggregate root for product domain
type Product struct {
	ID       int64
	SKUName  string
	Quantity int64
}

// NewProduct creates a new product
func NewProduct(skuName string, quantity int64) (*Product, error) {
	if skuName == "" {
		return nil, errors.New("SKU name cannot be empty")
	}
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	return &Product{
		SKUName:  skuName,
		Quantity: quantity,
	}, nil
}

// IncreaseStock increases the product quantity
func (p *Product) IncreaseStock(quantity int64) error {
	if quantity <= 0 {
		return errors.New("quantity to increase must be positive")
	}
	p.Quantity += quantity
	return nil
}

// DecreaseStock decreases the product quantity
func (p *Product) DecreaseStock(quantity int64) error {
	if quantity <= 0 {
		return errors.New("quantity to decrease must be positive")
	}
	if p.Quantity < quantity {
		return errors.New("insufficient stock")
	}
	p.Quantity -= quantity
	return nil
}
