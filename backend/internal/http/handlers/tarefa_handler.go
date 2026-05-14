package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/tarefas"
	"github.com/labstack/echo/v4"
)

type TarefaHandler struct {
	tarefaService *tarefas.TarefaService
}

func NewTarefaHandler(tarefaService *tarefas.TarefaService) *TarefaHandler {
	return &TarefaHandler{
		tarefaService: tarefaService,
	}
}

// CreateTarefa handles POST /api/v1/tarefas
func (h *TarefaHandler) CreateTarefa(c echo.Context) error {
	var req tarefas.CreateTarefaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Tarefa == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "tarefa is required",
		})
	}

	if req.DtTarefa == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "dt_tarefa is required",
		})
	}

	if req.HrTarefa == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "hr_tarefa is required",
		})
	}

	if req.Login == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "login is required",
		})
	}

	tarefa, err := h.tarefaService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, tarefa)
}

// GetTarefa handles GET /api/v1/tarefas/:id
func (h *TarefaHandler) GetTarefa(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid tarefa ID",
		})
	}

	tarefa, err := h.tarefaService.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "tarefa not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Tarefa not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tarefa)
}

// ListTarefas handles GET /api/v1/tarefas
func (h *TarefaHandler) ListTarefas(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	tarefas, err := h.tarefaService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tarefas": tarefas,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateTarefa handles PUT /api/v1/tarefas/:id
func (h *TarefaHandler) UpdateTarefa(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid tarefa ID",
		})
	}

	var req tarefas.UpdateTarefaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	tarefa, err := h.tarefaService.Update(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "tarefa not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Tarefa not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tarefa)
}
