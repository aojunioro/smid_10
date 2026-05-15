package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/km"
	"github.com/labstack/echo/v4"
)

type KmHandler struct {
	kmService *km.KmService
}

func NewKmHandler(kmService *km.KmService) *KmHandler {
	return &KmHandler{
		kmService: kmService,
	}
}

// KmConfig handlers
func (h *KmHandler) CreateKmConfig(c echo.Context) error {
	var req km.CreateKmConfigRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	config, err := h.kmService.CreateKmConfig(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, config)
}

func (h *KmHandler) GetKmConfig(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	config, err := h.kmService.GetKmConfigByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if config == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Configuração de KM não encontrada"})
	}

	return c.JSON(http.StatusOK, config)
}

func (h *KmHandler) ListKmConfigs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	configs, err := h.kmService.ListKmConfigs(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"configs": configs,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *KmHandler) UpdateKmConfig(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req km.UpdateKmConfigRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	config, err := h.kmService.UpdateKmConfig(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, config)
}

func (h *KmHandler) DeleteKmConfig(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.kmService.DeleteKmConfig(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// KmValorKmVigencia handlers
func (h *KmHandler) CreateKmValorKmVigencia(c echo.Context) error {
	var req km.CreateKmValorKmVigenciaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	valorKm, err := h.kmService.CreateKmValorKmVigencia(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, valorKm)
}

func (h *KmHandler) GetKmValorKmVigencia(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	valorKm, err := h.kmService.GetKmValorKmVigenciaByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if valorKm == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Valor de KM por vigência não encontrado"})
	}

	return c.JSON(http.StatusOK, valorKm)
}

func (h *KmHandler) ListKmValorKmVigencia(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	valores, err := h.kmService.ListKmValorKmVigencia(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"valores": valores,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *KmHandler) UpdateKmValorKmVigencia(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req km.UpdateKmValorKmVigenciaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	valorKm, err := h.kmService.UpdateKmValorKmVigencia(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, valorKm)
}

func (h *KmHandler) DeleteKmValorKmVigencia(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.kmService.DeleteKmValorKmVigencia(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// KmReembolsoLote handlers
func (h *KmHandler) CreateKmReembolsoLote(c echo.Context) error {
	var req km.CreateKmReembolsoLoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	lote, err := h.kmService.CreateKmReembolsoLote(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, lote)
}

func (h *KmHandler) GetKmReembolsoLote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	lote, err := h.kmService.GetKmReembolsoLoteByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if lote == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Lote de reembolso não encontrado"})
	}

	return c.JSON(http.StatusOK, lote)
}

func (h *KmHandler) ListKmReembolsoLotes(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var loginRepre *string
	if loginRepreStr := c.QueryParam("login_repre"); loginRepreStr != "" {
		loginRepre = &loginRepreStr
	}

	var statusPagamento *string
	if statusPagamentoStr := c.QueryParam("status_pagamento"); statusPagamentoStr != "" {
		statusPagamento = &statusPagamentoStr
	}

	lotes, err := h.kmService.ListKmReembolsoLotes(c.Request().Context(), limit, offset, loginRepre, statusPagamento)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"lotes":  lotes,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *KmHandler) UpdateKmReembolsoLote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req km.UpdateKmReembolsoLoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	lote, err := h.kmService.UpdateKmReembolsoLote(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, lote)
}

func (h *KmHandler) DeleteKmReembolsoLote(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.kmService.DeleteKmReembolsoLote(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
