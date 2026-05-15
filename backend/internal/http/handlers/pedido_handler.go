package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/pedidos"
	"github.com/labstack/echo/v4"
)

type PedidoHandler struct {
	pedidoService *pedidos.PedidoService
}

func NewPedidoHandler(pedidoService *pedidos.PedidoService) *PedidoHandler {
	return &PedidoHandler{
		pedidoService: pedidoService,
	}
}

// CreatePedido cria um novo pedido
func (h *PedidoHandler) CreatePedido(c echo.Context) error {
	var req pedidos.CreatePedidoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	pedido, err := h.pedidoService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, pedido)
}

// GetPedido busca um pedido por ID
func (h *PedidoHandler) GetPedido(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	pedido, err := h.pedidoService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if pedido == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Pedido não encontrado"})
	}

	return c.JSON(http.StatusOK, pedido)
}

// ListPedidos lista pedidos com filtros opcionais
func (h *PedidoHandler) ListPedidos(c echo.Context) error {
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

	var statusID *int64
	if statusIDStr := c.QueryParam("status_id"); statusIDStr != "" {
		if id, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			statusID = &id
		}
	}

	var loginRepre *string
	if loginRepreStr := c.QueryParam("login_repre"); loginRepreStr != "" {
		loginRepre = &loginRepreStr
	}

	var dtPed *string
	if dtPedStr := c.QueryParam("dt_ped"); dtPedStr != "" {
		dtPed = &dtPedStr
	}

	pedidos, err := h.pedidoService.List(c.Request().Context(), limit, offset, leadID, statusID, loginRepre, dtPed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"pedidos": pedidos,
		"limit":   limit,
		"offset":  offset,
	})
}

// UpdatePedido atualiza um pedido existente
func (h *PedidoHandler) UpdatePedido(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req pedidos.UpdatePedidoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	pedido, err := h.pedidoService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, pedido)
}

// DeletePedido exclui um pedido (soft delete)
func (h *PedidoHandler) DeletePedido(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.pedidoService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
