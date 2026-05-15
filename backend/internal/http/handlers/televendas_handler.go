package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/televendas"
	"github.com/labstack/echo/v4"
)

type TelevendasHandler struct {
	contatoService *televendas.TelevendasContatoService
}

func NewTelevendasHandler(contatoService *televendas.TelevendasContatoService) *TelevendasHandler {
	return &TelevendasHandler{
		contatoService: contatoService,
	}
}

// CreateContato cria um novo contato de televendas
func (h *TelevendasHandler) CreateContato(c echo.Context) error {
	var req televendas.CreateContatoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	contato, err := h.contatoService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, contato)
}

// GetContato busca um contato por ID
func (h *TelevendasHandler) GetContato(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	contato, err := h.contatoService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if contato == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Contato não encontrado"})
	}

	return c.JSON(http.StatusOK, contato)
}

// ListContatos lista contatos com filtros opcionais
func (h *TelevendasHandler) ListContatos(c echo.Context) error {
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

	var login *string
	if loginStr := c.QueryParam("login"); loginStr != "" {
		login = &loginStr
	}

	var statusID *string
	if statusIDStr := c.QueryParam("status_id"); statusIDStr != "" {
		statusID = &statusIDStr
	}

	contatos, err := h.contatoService.List(c.Request().Context(), limit, offset, leadID, login, statusID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"contatos": contatos,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateContato atualiza um contato existente
func (h *TelevendasHandler) UpdateContato(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req televendas.UpdateContatoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	contato, err := h.contatoService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, contato)
}

// DeleteContato exclui um contato (soft delete)
func (h *TelevendasHandler) DeleteContato(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.contatoService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
