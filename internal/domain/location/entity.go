package location

import "errors"

// Location is the aggregate root for location domain
type Location struct {
	ID       int64
	Code     string
	Name     string
	Capacity int64
}

// NewLocation creates a new location
func NewLocation(code, name string, capacity int64) (*Location, error) {
	if code == "" {
		return nil, errors.New("location code cannot be empty")
	}
	if name == "" {
		return nil, errors.New("location name cannot be empty")
	}
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}

	return &Location{
		Code:     code,
		Name:     name,
		Capacity: capacity,
	}, nil
}

// CanAccommodate checks if location can accommodate given quantity
func (l *Location) CanAccommodate(currentStock, incomingQuantity int64) bool {
	return currentStock+incomingQuantity <= l.Capacity
}
