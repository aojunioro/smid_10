package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/financeiro"
	"github.com/labstack/echo/v4"
)

type FinanceiroHandler struct {
	financeiroService *financeiro.FinanceiroService
}

func NewFinanceiroHandler(financeiroService *financeiro.FinanceiroService) *FinanceiroHandler {
	return &FinanceiroHandler{
		financeiroService: financeiroService,
	}
}

// Contas a Pagar handlers
func (h *FinanceiroHandler) CreateContaPagar(c echo.Context) error {
	var req financeiro.CreateContaPagarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	conta, err := h.financeiroService.CreateContaPagar(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, conta)
}

func (h *FinanceiroHandler) GetContaPagar(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	conta, err := h.financeiroService.GetContaPagarByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if conta == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Conta a pagar não encontrada"})
	}

	return c.JSON(http.StatusOK, conta)
}

func (h *FinanceiroHandler) ListContasPagar(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var categoriaID *int64
	if categoriaIDStr := c.QueryParam("categoria_id"); categoriaIDStr != "" {
		if id, err := strconv.ParseInt(categoriaIDStr, 10, 64); err == nil {
			categoriaID = &id
		}
	}

	var pedidoID *int64
	if pedidoIDStr := c.QueryParam("pedido_id"); pedidoIDStr != "" {
		if id, err := strconv.ParseInt(pedidoIDStr, 10, 64); err == nil {
			pedidoID = &id
		}
	}

	var status *string
	if statusStr := c.QueryParam("status"); statusStr != "" {
		status = &statusStr
	}

	contas, err := h.financeiroService.ListContasPagar(c.Request().Context(), limit, offset, categoriaID, pedidoID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"contas_pagar": contas,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *FinanceiroHandler) UpdateContaPagar(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req financeiro.UpdateContaPagarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	conta, err := h.financeiroService.UpdateContaPagar(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, conta)
}

func (h *FinanceiroHandler) DeleteContaPagar(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.financeiroService.DeleteContaPagar(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Contas a Receber handlers
func (h *FinanceiroHandler) CreateContaReceber(c echo.Context) error {
	var req financeiro.CreateContaReceberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	conta, err := h.financeiroService.CreateContaReceber(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, conta)
}

func (h *FinanceiroHandler) GetContaReceber(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	conta, err := h.financeiroService.GetContaReceberByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if conta == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Conta a receber não encontrada"})
	}

	return c.JSON(http.StatusOK, conta)
}

func (h *FinanceiroHandler) ListContasReceber(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var categoriaID *int64
	if categoriaIDStr := c.QueryParam("categoria_id"); categoriaIDStr != "" {
		if id, err := strconv.ParseInt(categoriaIDStr, 10, 64); err == nil {
			categoriaID = &id
		}
	}

	var pedidoID *int64
	if pedidoIDStr := c.QueryParam("pedido_id"); pedidoIDStr != "" {
		if id, err := strconv.ParseInt(pedidoIDStr, 10, 64); err == nil {
			pedidoID = &id
		}
	}

	var status *string
	if statusStr := c.QueryParam("status"); statusStr != "" {
		status = &statusStr
	}

	contas, err := h.financeiroService.ListContasReceber(c.Request().Context(), limit, offset, categoriaID, pedidoID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"contas_receber": contas,
		"limit":          limit,
		"offset":         offset,
	})
}

func (h *FinanceiroHandler) UpdateContaReceber(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req financeiro.UpdateContaReceberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	conta, err := h.financeiroService.UpdateContaReceber(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, conta)
}

func (h *FinanceiroHandler) DeleteContaReceber(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.financeiroService.DeleteContaReceber(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
