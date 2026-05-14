package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/visitas"
	"github.com/labstack/echo/v4"
)

type VisitaHandler struct {
	visitaService *visitas.VisitaService
}

func NewVisitaHandler(visitaService *visitas.VisitaService) *VisitaHandler {
	return &VisitaHandler{
		visitaService: visitaService,
	}
}

// CreateVisita cria uma nova visita
func (h *VisitaHandler) CreateVisita(c echo.Context) error {
	var req visitas.CreateVisitaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	visita, err := h.visitaService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, visita)
}

// GetVisita busca uma visita por ID
func (h *VisitaHandler) GetVisita(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	visita, err := h.visitaService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if visita == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Visita não encontrada"})
	}

	return c.JSON(http.StatusOK, visita)
}

// ListVisitas lista visitas com filtros opcionais
func (h *VisitaHandler) ListVisitas(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var leadID *int64
	if leadIDStr := c.QueryParam("lead_id"); leadIDStr != "" {
		if id, err := strconv.ParseInt(leadIDStr, 10, 64); err == nil {
			leadID = &id
		}
	}

	var loginRepre *string
	if loginRepreStr := c.QueryParam("login_repre"); loginRepreStr != "" {
		loginRepre = &loginRepreStr
	}

	var loginRecep *string
	if loginRecepStr := c.QueryParam("login_recep"); loginRecepStr != "" {
		loginRecep = &loginRecepStr
	}

	var statusID *int64
	if statusIDStr := c.QueryParam("status_id"); statusIDStr != "" {
		if id, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			statusID = &id
		}
	}

	var dtVisita *string
	if dtVisitaStr := c.QueryParam("dt_visita"); dtVisitaStr != "" {
		dtVisita = &dtVisitaStr
	}

	visitas, err := h.visitaService.List(c.Request().Context(), limit, offset, leadID, loginRepre, loginRecep, statusID, dtVisita)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"visitas": visitas,
		"limit":   limit,
		"offset":  offset,
	})
}

// UpdateVisita atualiza uma visita existente
func (h *VisitaHandler) UpdateVisita(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req visitas.UpdateVisitaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	visita, err := h.visitaService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, visita)
}

// DeleteVisita exclui uma visita (soft delete)
func (h *VisitaHandler) DeleteVisita(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.visitaService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
