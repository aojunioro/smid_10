package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/suporte"
	"github.com/labstack/echo/v4"
)

type SuporteHandler struct {
	suporteService *suporte.SuporteService
}

func NewSuporteHandler(suporteService *suporte.SuporteService) *SuporteHandler {
	return &SuporteHandler{
		suporteService: suporteService,
	}
}

// CreateSuporte cria um novo chamado de suporte
func (h *SuporteHandler) CreateSuporte(c echo.Context) error {
	var req suporte.CreateSuporteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	suporte, err := h.suporteService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, suporte)
}

// GetSuporte busca um suporte por ID
func (h *SuporteHandler) GetSuporte(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	suporte, err := h.suporteService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if suporte == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Suporte não encontrado"})
	}

	return c.JSON(http.StatusOK, suporte)
}

// ListSuportes lista suportes com filtros opcionais
func (h *SuporteHandler) ListSuportes(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var pedID *int64
	if pedIDStr := c.QueryParam("ped_id"); pedIDStr != "" {
		if id, err := strconv.ParseInt(pedIDStr, 10, 64); err == nil {
			pedID = &id
		}
	}

	var login *string
	if loginStr := c.QueryParam("login"); loginStr != "" {
		login = &loginStr
	}

	var atribLogin *string
	if atribLoginStr := c.QueryParam("atrib_login"); atribLoginStr != "" {
		atribLogin = &atribLoginStr
	}

	var statusID *string
	if statusIDStr := c.QueryParam("status_id"); statusIDStr != "" {
		statusID = &statusIDStr
	}

	suportes, err := h.suporteService.List(c.Request().Context(), limit, offset, pedID, login, atribLogin, statusID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"suportes": suportes,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateSuporte atualiza um suporte existente
func (h *SuporteHandler) UpdateSuporte(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req suporte.UpdateSuporteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	suporte, err := h.suporteService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, suporte)
}

// DeleteSuporte exclui um suporte (soft delete)
func (h *SuporteHandler) DeleteSuporte(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.suporteService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
