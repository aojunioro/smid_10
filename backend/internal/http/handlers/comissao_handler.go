package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/comissoes"
	"github.com/labstack/echo/v4"
)

type ComissaoHandler struct {
	comissaoService *comissoes.ComissaoService
}

func NewComissaoHandler(comissaoService *comissoes.ComissaoService) *ComissaoHandler {
	return &ComissaoHandler{
		comissaoService: comissaoService,
	}
}

// Comissao handlers
func (h *ComissaoHandler) CreateComissao(c echo.Context) error {
	var req comissoes.CreateComissaoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	comissao, err := h.comissaoService.CreateComissao(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, comissao)
}

func (h *ComissaoHandler) GetComissao(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	comissao, err := h.comissaoService.GetComissaoByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if comissao == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Comissão não encontrada"})
	}

	return c.JSON(http.StatusOK, comissao)
}

func (h *ComissaoHandler) ListComissoes(c echo.Context) error {
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

	var sttComis *string
	if sttComisStr := c.QueryParam("stt_comis"); sttComisStr != "" {
		sttComis = &sttComisStr
	}

	comissoes, err := h.comissaoService.ListComissoes(c.Request().Context(), limit, offset, pedID, sttComis)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"comissoes": comissoes,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *ComissaoHandler) UpdateComissao(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req comissoes.UpdateComissaoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	comissao, err := h.comissaoService.UpdateComissao(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, comissao)
}

func (h *ComissaoHandler) DeleteComissao(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.comissaoService.DeleteComissao(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ComissItem handlers
func (h *ComissaoHandler) CreateComissItem(c echo.Context) error {
	var req comissoes.CreateComissItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	item, err := h.comissaoService.CreateComissItem(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *ComissaoHandler) GetComissItem(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	item, err := h.comissaoService.GetComissItemByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Item de comissão não encontrado"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ComissaoHandler) ListComissItems(c echo.Context) error {
	comisID, err := strconv.ParseInt(c.Param("comis_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_comis_id", "message": "ID de comissão inválido"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	itens, err := h.comissaoService.ListComissItems(c.Request().Context(), comisID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"itens":  itens,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *ComissaoHandler) UpdateComissItem(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req comissoes.UpdateComissItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	item, err := h.comissaoService.UpdateComissItem(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ComissaoHandler) DeleteComissItem(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.comissaoService.DeleteComissItem(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
