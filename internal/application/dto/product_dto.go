package dto

// CreateProductRequest is the DTO for creating a product
type CreateProductRequest struct {
	SKUName  string `json:"sku_name" binding:"required"`
	Quantity int64  `json:"quantity" binding:"required,min=0"`
}

// UpdateProductRequest is the DTO for updating a product
type UpdateProductRequest struct {
	SKUName  string `json:"sku_name"`
	Quantity int64  `json:"quantity" binding:"min=0"`
}

// ProductResponse is the DTO for product response
type ProductResponse struct {
	ID       int64  `json:"id"`
	SKUName  string `json:"sku_name"`
	Quantity int64  `json:"quantity"`
}

// ProductListResponse is the DTO for product list response
type ProductListResponse struct {
	Data   []*ProductResponse `json:"data"`
	Total  int64              `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}
