package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/historicos"
	"github.com/labstack/echo/v4"
)

type HistoricoHandler struct {
	historicoService *historicos.HistoricoService
}

func NewHistoricoHandler(historicoService *historicos.HistoricoService) *HistoricoHandler {
	return &HistoricoHandler{
		historicoService: historicoService,
	}
}

// CreateHistorico cria um novo histórico
func (h *HistoricoHandler) CreateHistorico(c echo.Context) error {
	var req historicos.CreateHistoricoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	historico, err := h.historicoService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, historico)
}

// GetHistorico busca um histórico por ID
func (h *HistoricoHandler) GetHistorico(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	historico, err := h.historicoService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if historico == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Histórico não encontrado"})
	}

	return c.JSON(http.StatusOK, historico)
}

// ListHistoricos lista históricos com filtros opcionais
func (h *HistoricoHandler) ListHistoricos(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var visID *int64
	if visIDStr := c.QueryParam("vis_id"); visIDStr != "" {
		if id, err := strconv.ParseInt(visIDStr, 10, 64); err == nil {
			visID = &id
		}
	}

	var leadID *int64
	if leadIDStr := c.QueryParam("lead_id"); leadIDStr != "" {
		if id, err := strconv.ParseInt(leadIDStr, 10, 64); err == nil {
			leadID = &id
		}
	}

	var login *string
	if loginStr := c.QueryParam("login"); loginStr != "" {
		login = &loginStr
	}

	var motivoID *int64
	if motivoIDStr := c.QueryParam("motivo_id"); motivoIDStr != "" {
		if id, err := strconv.ParseInt(motivoIDStr, 10, 64); err == nil {
			motivoID = &id
		}
	}

	var ocorridoID *int64
	if ocorridoIDStr := c.QueryParam("ocorrido_id"); ocorridoIDStr != "" {
		if id, err := strconv.ParseInt(ocorridoIDStr, 10, 64); err == nil {
			ocorridoID = &id
		}
	}

	historicos, err := h.historicoService.List(c.Request().Context(), limit, offset, visID, leadID, login, motivoID, ocorridoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"historicos": historicos,
		"limit":      limit,
		"offset":     offset,
	})
}

// UpdateHistorico atualiza um histórico existente
func (h *HistoricoHandler) UpdateHistorico(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req historicos.UpdateHistoricoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	historico, err := h.historicoService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, historico)
}

// DeleteHistorico exclui um histórico (soft delete)
func (h *HistoricoHandler) DeleteHistorico(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.historicoService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
