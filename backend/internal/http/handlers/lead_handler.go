package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/leads"
	"github.com/labstack/echo/v4"
)

type LeadHandler struct {
	leadService *leads.LeadService
}

func NewLeadHandler(leadService *leads.LeadService) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
	}
}

// CreateLead handles POST /api/v1/leads
func (h *LeadHandler) CreateLead(c echo.Context) error {
	var req leads.CreateLeadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Fone1 == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "fone1 is required",
		})
	}

	if req.Nome == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "nome is required",
		})
	}

	if req.StatusID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "status_id is required",
		})
	}

	lead, err := h.leadService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, lead)
}

// GetLead handles GET /api/v1/leads/:id
func (h *LeadHandler) GetLead(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid lead ID",
		})
	}

	lead, err := h.leadService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	if lead == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "not_found",
			"message": "Lead not found",
		})
	}

	return c.JSON(http.StatusOK, lead)
}

// ListLeads handles GET /api/v1/leads
func (h *LeadHandler) ListLeads(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	leads, err := h.leadService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"leads": leads,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateLead handles PUT /api/v1/leads/:id
func (h *LeadHandler) UpdateLead(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid lead ID",
		})
	}

	var req leads.UpdateLeadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	lead, err := h.leadService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, lead)
}
