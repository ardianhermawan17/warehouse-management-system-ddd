package commands

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
)

// CreateProductCommand handles product creation
type CreateProductCommand struct {
	productRepo product.Repository
}

// NewCreateProductCommand creates a new create product command
func NewCreateProductCommand(productRepo product.Repository) *CreateProductCommand {
	return &CreateProductCommand{
		productRepo: productRepo,
	}
}

// Execute executes the create product command
func (c *CreateProductCommand) Execute(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Create product entity
	prod, err := product.NewProduct(req.SKUName, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := c.productRepo.Create(ctx, prod); err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:       prod.ID,
		SKUName:  prod.SKUName,
		Quantity: prod.Quantity,
	}, nil
}
