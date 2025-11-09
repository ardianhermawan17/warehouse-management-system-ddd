package handlers

import (
	"net/http"
	"strconv"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/commands"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/queries"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/product"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	createCmd   *commands.CreateProductCommand
	updateCmd   *commands.UpdateProductCommand
	listQuery   *queries.ListProductsQuery
	productRepo product.Repository
}

// NewProductHandler creates a new product handler
func NewProductHandler(
	createCmd *commands.CreateProductCommand,
	updateCmd *commands.UpdateProductCommand,
	listQuery *queries.ListProductsQuery,
	productRepo product.Repository,
) *ProductHandler {
	return &ProductHandler{
		createCmd:   createCmd,
		updateCmd:   updateCmd,
		listQuery:   listQuery,
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	result, err := h.createCmd.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("product created successfully", result))
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid product ID"))
		return
	}

	prod, err := h.productRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("product not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("product retrieved successfully", &dto.ProductResponse{
		ID:       prod.ID,
		SKUName:  prod.SKUName,
		Quantity: prod.Quantity,
	}))
}

// ListProducts lists all products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	result, err := h.listQuery.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to list products"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("products retrieved successfully", result))
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid product ID"))
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	result, err := h.updateCmd.Execute(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("product updated successfully", result))
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid product ID"))
		return
	}

	if err := h.productRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("product not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("product deleted successfully", nil))
}
