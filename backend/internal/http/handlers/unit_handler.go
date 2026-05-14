package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/labstack/echo/v4"
)

type UnitHandler struct {
	unitService *admin.UnitService
}

func NewUnitHandler(unitService *admin.UnitService) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
	}
}

// CreateUnit handles POST /api/v1/units
func (h *UnitHandler) CreateUnit(c echo.Context) error {
	var req admin.CreateUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "name is required",
		})
	}

	unit, err := h.unitService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, unit)
}

// GetUnit handles GET /api/v1/units/:id
func (h *UnitHandler) GetUnit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid unit ID",
		})
	}

	unit, err := h.unitService.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "unit not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Unit not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, unit)
}

// ListUnits handles GET /api/v1/units
func (h *UnitHandler) ListUnits(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	units, err := h.unitService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"units": units,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateUnit handles PUT /api/v1/units/:id
func (h *UnitHandler) UpdateUnit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid unit ID",
		})
	}

	var req admin.UpdateUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "name is required",
		})
	}

	unit, err := h.unitService.Update(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "unit not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Unit not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, unit)
}

// DeleteUnit handles DELETE /api/v1/units/:id
func (h *UnitHandler) DeleteUnit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid unit ID",
		})
	}

	if err := h.unitService.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "unit not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Unit not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
