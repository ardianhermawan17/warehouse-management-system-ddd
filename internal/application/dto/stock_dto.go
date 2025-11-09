package dto

import "time"

// RecordStockMovementRequest is the DTO for recording stock movement
type RecordStockMovementRequest struct {
	ProductID  int64  `json:"product_id" binding:"required,min=1"`
	LocationID int64  `json:"location_id" binding:"required,min=1"`
	Type       string `json:"type" binding:"required,oneof=IN OUT"`
	Quantity   int64  `json:"quantity" binding:"required,min=1"`
}

// StockMovementResponse is the DTO for stock movement response
type StockMovementResponse struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	LocationID int64     `json:"location_id"`
	Type       string    `json:"type"`
	Quantity   int64     `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
}

// StockMovementListResponse is the DTO for stock movement list response
type StockMovementListResponse struct {
	Data   []*StockMovementResponse `json:"data"`
	Total  int64                    `json:"total"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

// LocationStockResponse is the DTO for location stock response
type LocationStockResponse struct {
	LocationID int64 `json:"location_id"`
	ProductID  int64 `json:"product_id"`
	Quantity   int64 `json:"quantity"`
}

// LocationRequest is the DTO for creating/updating location
type LocationRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Capacity int64  `json:"capacity" binding:"required,min=1"`
}

// LocationResponse is the DTO for location response
type LocationResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Capacity int64  `json:"capacity"`
}

// LocationListResponse is the DTO for location list response
type LocationListResponse struct {
	Data   []*LocationResponse `json:"data"`
	Total  int64               `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}
