package handlers

import (
	"net/http"
	"strconv"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/commands"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/queries"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/stock"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// StockHandler handles stock movement endpoints
type StockHandler struct {
	recordCmd *commands.RecordStockMovementCommand
	listQuery *queries.ListStockMovementsQuery
	stockRepo stock.Repository
}

// NewStockHandler creates a new stock handler
func NewStockHandler(
	recordCmd *commands.RecordStockMovementCommand,
	listQuery *queries.ListStockMovementsQuery,
	stockRepo stock.Repository,
) *StockHandler {
	return &StockHandler{
		recordCmd: recordCmd,
		listQuery: listQuery,
		stockRepo: stockRepo,
	}
}

// RecordMovement records a stock movement
func (h *StockHandler) RecordMovement(c *gin.Context) {
	var req dto.RecordStockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	result, err := h.recordCmd.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("stock movement recorded successfully", result))
}

// GetMovement retrieves a stock movement by ID
func (h *StockHandler) GetMovement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid movement ID"))
		return
	}

	movement, err := h.stockRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("movement not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("movement retrieved successfully", &dto.StockMovementResponse{
		ID:         movement.ID,
		ProductID:  movement.ProductID,
		LocationID: movement.LocationID,
		Type:       string(movement.Type),
		Quantity:   movement.Quantity,
		CreatedAt:  movement.CreatedAt,
	}))
}

// ListMovements lists all stock movements
func (h *StockHandler) ListMovements(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to list movements"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("movements retrieved successfully", result))
}

// GetProductMovements retrieves all movements for a product
func (h *StockHandler) GetProductMovements(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid product ID"))
		return
	}

	movements, err := h.stockRepo.GetByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to get movements"))
		return
	}

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

	c.JSON(http.StatusOK, response.SuccessResponse("movements retrieved successfully", responses))
}

// GetLocationMovements retrieves all movements for a location
func (h *StockHandler) GetLocationMovements(c *gin.Context) {
	locationID, err := strconv.ParseInt(c.Param("location_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid location ID"))
		return
	}

	movements, err := h.stockRepo.GetByLocation(c.Request.Context(), locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to get movements"))
		return
	}

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

	c.JSON(http.StatusOK, response.SuccessResponse("movements retrieved successfully", responses))
}
