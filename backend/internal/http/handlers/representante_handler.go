package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/representantes"
	"github.com/labstack/echo/v4"
)

type RepresentanteHandler struct {
	despesaService *representantes.RepreDespesaExtraService
}

func NewRepresentanteHandler(despesaService *representantes.RepreDespesaExtraService) *RepresentanteHandler {
	return &RepresentanteHandler{
		despesaService: despesaService,
	}
}

// CreateDespesaExtra cria uma nova despesa extra do representante
func (h *RepresentanteHandler) CreateDespesaExtra(c echo.Context) error {
	var req representantes.CreateDespesaExtraRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	despesa, err := h.despesaService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, despesa)
}

// GetDespesaExtra busca uma despesa extra por ID
func (h *RepresentanteHandler) GetDespesaExtra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	despesa, err := h.despesaService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if despesa == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Despesa extra não encontrada"})
	}

	return c.JSON(http.StatusOK, despesa)
}

// ListDespesasExtras lista despesas extras com filtros opcionais
func (h *RepresentanteHandler) ListDespesasExtras(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var login *string
	if loginStr := c.QueryParam("login"); loginStr != "" {
		login = &loginStr
	}

	var status *string
	if statusStr := c.QueryParam("status"); statusStr != "" {
		status = &statusStr
	}

	despesas, err := h.despesaService.List(c.Request().Context(), limit, offset, login, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"despesas": despesas,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateDespesaExtra atualiza uma despesa extra existente
func (h *RepresentanteHandler) UpdateDespesaExtra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req representantes.UpdateDespesaExtraRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	despesa, err := h.despesaService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, despesa)
}

// DeleteDespesaExtra exclui uma despesa extra (soft delete)
func (h *RepresentanteHandler) DeleteDespesaExtra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.despesaService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
