package stock

import (
	"context"
	"errors"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
)

// Service contains business rules for stock movements
type Service struct {
	productRepo  product.Repository
	locationRepo location.Repository
	stockRepo    Repository
}

// NewService creates a new stock service
func NewService(
	productRepo product.Repository,
	locationRepo location.Repository,
	stockRepo Repository,
) *Service {
	return &Service{
		productRepo:  productRepo,
		locationRepo: locationRepo,
		stockRepo:    stockRepo,
	}
}

// RecordMovement records a stock movement with business rule validation
func (s *Service) RecordMovement(ctx context.Context, movement *StockMovement) error {
	// Validate product exists
	prod, err := s.productRepo.GetByID(ctx, movement.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Validate location exists
	loc, err := s.locationRepo.GetByID(ctx, movement.LocationID)
	if err != nil {
		return errors.New("location not found")
	}

	// Apply business rules based on movement type
	if movement.IsOutbound() {
		// Stock OUT cannot exceed available stock
		if prod.Quantity < movement.Quantity {
			return errors.New("insufficient stock for outbound movement")
		}
	} else if movement.IsInbound() {
		// Stock IN cannot exceed location capacity
		currentStock, err := s.getLocationStock(ctx, movement.LocationID)
		if err != nil {
			return err
		}

		if !loc.CanAccommodate(currentStock, movement.Quantity) {
			return errors.New("location capacity exceeded")
		}
	}

	// Record the movement
	return s.stockRepo.Create(ctx, movement)
}

// getLocationStock gets total stock at a location
func (s *Service) getLocationStock(ctx context.Context, locationID int64) (int64, error) {
	movements, err := s.stockRepo.GetByLocation(ctx, locationID)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, m := range movements {
		if m.IsInbound() {
			total += m.Quantity
		} else {
			total -= m.Quantity
		}
	}

	return total, nil
}
