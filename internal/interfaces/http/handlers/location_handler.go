package handlers

import (
	"net/http"
	"strconv"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/application/dto"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/domain/location"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// LocationHandler handles location endpoints
type LocationHandler struct {
	locationRepo location.Repository
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(locationRepo location.Repository) *LocationHandler {
	return &LocationHandler{
		locationRepo: locationRepo,
	}
}

// CreateLocation creates a new location
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	loc, err := location.NewLocation(req.Code, req.Name, req.Capacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	if err := h.locationRepo.Create(c.Request.Context(), loc); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to create location"))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("location created successfully", &dto.LocationResponse{
		ID:       loc.ID,
		Code:     loc.Code,
		Name:     loc.Name,
		Capacity: loc.Capacity,
	}))
}

// GetLocation retrieves a location by ID
func (h *LocationHandler) GetLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid location ID"))
		return
	}

	loc, err := h.locationRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("location not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("location retrieved successfully", &dto.LocationResponse{
		ID:       loc.ID,
		Code:     loc.Code,
		Name:     loc.Name,
		Capacity: loc.Capacity,
	}))
}

// ListLocations lists all locations
func (h *LocationHandler) ListLocations(c *gin.Context) {
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

	locations, err := h.locationRepo.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to list locations"))
		return
	}

	total, err := h.locationRepo.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to count locations"))
		return
	}

	var responses []*dto.LocationResponse
	for _, loc := range locations {
		responses = append(responses, &dto.LocationResponse{
			ID:       loc.ID,
			Code:     loc.Code,
			Name:     loc.Name,
			Capacity: loc.Capacity,
		})
	}

	c.JSON(http.StatusOK, response.SuccessResponse("locations retrieved successfully", &dto.LocationListResponse{
		Data:   responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}))
}

// UpdateLocation updates a location
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid location ID"))
		return
	}

	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	loc, err := h.locationRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("location not found"))
		return
	}

	loc.Code = req.Code
	loc.Name = req.Name
	loc.Capacity = req.Capacity

	if err := h.locationRepo.Update(c.Request.Context(), loc); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to update location"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("location updated successfully", &dto.LocationResponse{
		ID:       loc.ID,
		Code:     loc.Code,
		Name:     loc.Name,
		Capacity: loc.Capacity,
	}))
}

// DeleteLocation deletes a location
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid location ID"))
		return
	}

	if err := h.locationRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("location not found"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("location deleted successfully", nil))
}
