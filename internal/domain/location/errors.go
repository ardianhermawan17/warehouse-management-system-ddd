package location

import "errors"

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrInvalidCode      = errors.New("invalid location code")
	ErrInvalidName      = errors.New("invalid location name")
	ErrInvalidCapacity  = errors.New("invalid capacity")
	ErrDuplicateCode    = errors.New("location code already exists")
	ErrCapacityExceeded = errors.New("location capacity exceeded")
)
