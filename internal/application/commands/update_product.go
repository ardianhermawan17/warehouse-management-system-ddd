package commands

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
)

// UpdateProductCommand handles product update
type UpdateProductCommand struct {
	productRepo product.Repository
}

// NewUpdateProductCommand creates a new update product command
func NewUpdateProductCommand(productRepo product.Repository) *UpdateProductCommand {
	return &UpdateProductCommand{
		productRepo: productRepo,
	}
}

// Execute executes the update product command
func (c *UpdateProductCommand) Execute(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Get existing product
	prod, err := c.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.SKUName != "" {
		prod.SKUName = req.SKUName
	}
	if req.Quantity >= 0 {
		prod.Quantity = req.Quantity
	}

	// Save to repository
	if err := c.productRepo.Update(ctx, prod); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:       prod.ID,
		SKUName:  prod.SKUName,
		Quantity: prod.Quantity,
	}, nil
}
