package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/produtos"
	"github.com/labstack/echo/v4"
)

type ProdutoHandler struct {
	produtoService *produtos.ProdutoService
}

func NewProdutoHandler(produtoService *produtos.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		produtoService: produtoService,
	}
}

// CreateProduto cria um novo produto
func (h *ProdutoHandler) CreateProduto(c echo.Context) error {
	var req produtos.CreateProdutoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	produto, err := h.produtoService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, produto)
}

// GetProduto busca um produto por ID
func (h *ProdutoHandler) GetProduto(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	produto, err := h.produtoService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}
	if produto == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not_found", "message": "Produto não encontrado"})
	}

	return c.JSON(http.StatusOK, produto)
}

// ListProdutos lista produtos com filtros opcionais
func (h *ProdutoHandler) ListProdutos(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	var categID *int64
	if categIDStr := c.QueryParam("categ_id"); categIDStr != "" {
		if id, err := strconv.ParseInt(categIDStr, 10, 64); err == nil {
			categID = &id
		}
	}

	var medID *int64
	if medIDStr := c.QueryParam("med_id"); medIDStr != "" {
		if id, err := strconv.ParseInt(medIDStr, 10, 64); err == nil {
			medID = &id
		}
	}

	var ativo *string
	if ativoStr := c.QueryParam("ativo"); ativoStr != "" {
		ativo = &ativoStr
	}

	var televendas *string
	if televendasStr := c.QueryParam("televendas"); televendasStr != "" {
		televendas = &televendasStr
	}

	produtos, err := h.produtoService.List(c.Request().Context(), limit, offset, categID, medID, ativo, televendas)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"produtos": produtos,
		"limit":   limit,
		"offset":  offset,
	})
}

// UpdateProduto atualiza um produto existente
func (h *ProdutoHandler) UpdateProduto(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	var req produtos.UpdateProdutoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": err.Error()})
	}

	produto, err := h.produtoService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.JSON(http.StatusOK, produto)
}

// DeleteProduto exclui um produto (soft delete)
func (h *ProdutoHandler) DeleteProduto(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_id", "message": "ID inválido"})
	}

	if err := h.produtoService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal_error", "message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
