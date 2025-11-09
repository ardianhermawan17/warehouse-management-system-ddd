package queries

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
)

// ListProductsQuery handles product listing
type ListProductsQuery struct {
	productRepo product.Repository
}

// NewListProductsQuery creates a new list products query
func NewListProductsQuery(productRepo product.Repository) *ListProductsQuery {
	return &ListProductsQuery{
		productRepo: productRepo,
	}
}

// Execute executes the list products query
func (q *ListProductsQuery) Execute(ctx context.Context, limit, offset int) (*dto.ProductListResponse, error) {
	// Get products
	products, err := q.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := q.productRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var responses []*dto.ProductResponse
	for _, p := range products {
		responses = append(responses, &dto.ProductResponse{
			ID:       p.ID,
			SKUName:  p.SKUName,
			Quantity: p.Quantity,
		})
	}

	return &dto.ProductListResponse{
		Data:   responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
