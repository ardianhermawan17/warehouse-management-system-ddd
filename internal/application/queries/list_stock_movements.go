package queries

import (
	"context"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
)

// ListStockMovementsQuery handles stock movement listing
type ListStockMovementsQuery struct {
	stockRepo stock.Repository
}

// NewListStockMovementsQuery creates a new list stock movements query
func NewListStockMovementsQuery(stockRepo stock.Repository) *ListStockMovementsQuery {
	return &ListStockMovementsQuery{
		stockRepo: stockRepo,
	}
}

// Execute executes the list stock movements query
func (q *ListStockMovementsQuery) Execute(ctx context.Context, limit, offset int) (*dto.StockMovementListResponse, error) {
	// Get movements
	movements, err := q.stockRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := q.stockRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var responses []*dto.StockMovementResponse
	for _, m := range movements {
		responses = append(responses, &dto.StockMovementResponse{
			ID:         m.ID,
			ProductID:  m.ProductID,
			LocationID: m.LocationID,
			Type:       string(m.Type),
			Quantity:   m.Quantity,
			CreatedAt:  m.CreatedAt,
		})
	}

	return &dto.StockMovementListResponse{
		Data:   responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
