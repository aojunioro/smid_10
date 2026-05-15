package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/compras"
	"github.com/labstack/echo/v4"
)

type CompraHandler struct {
	compraService *compras.CompraService
}

func NewCompraHandler(compraService *compras.CompraService) *CompraHandler {
	return &CompraHandler{
		compraService: compraService,
	}
}

// CreateCompra cria uma nova compra
func (h *CompraHandler) CreateCompra(c echo.Context) error {
	var req compras.CreateCompraRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	compra, err := h.compraService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, compra)
}

// GetCompra busca uma compra por ID
func (h *CompraHandler) GetCompra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	compra, err := h.compraService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if compra == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Compra não encontrada"})
	}

	return c.JSON(http.StatusOK, compra)
}

// ListCompras lista compras com filtros opcionais
func (h *CompraHandler) ListCompras(c echo.Context) error {
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

	var fornecID *int64
	if fornecIDStr := c.QueryParam("fornec_id"); fornecIDStr != "" {
		if id, err := strconv.ParseInt(fornecIDStr, 10, 64); err == nil {
			fornecID = &id
		}
	}

	var statusID *int64
	if statusIDStr := c.QueryParam("status_id"); statusIDStr != "" {
		if id, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			statusID = &id
		}
	}

	var transpID *int64
	if transpIDStr := c.QueryParam("transp_id"); transpIDStr != "" {
		if id, err := strconv.ParseInt(transpIDStr, 10, 64); err == nil {
			transpID = &id
		}
	}

	var login *string
	if loginStr := c.QueryParam("login"); loginStr != "" {
		login = &loginStr
	}

	compras, err := h.compraService.List(c.Request().Context(), limit, offset, pedID, fornecID, statusID, transpID, login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"compras": compras,
		"limit":   limit,
		"offset":  offset,
	})
}

// UpdateCompra atualiza uma compra existente
func (h *CompraHandler) UpdateCompra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req compras.UpdateCompraRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	compra, err := h.compraService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, compra)
}

// DeleteCompra exclui uma compra (soft delete)
func (h *CompraHandler) DeleteCompra(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.compraService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
