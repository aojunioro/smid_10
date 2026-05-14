package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/labstack/echo/v4"
)

type ProgramHandler struct {
	programService *admin.ProgramService
}

func NewProgramHandler(programService *admin.ProgramService) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
	}
}

// CreateProgram handles POST /api/v1/programs
func (h *ProgramHandler) CreateProgram(c echo.Context) error {
	var req admin.CreateProgramRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "name is required",
		})
	}

	program, err := h.programService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, program)
}

// GetProgram handles GET /api/v1/programs/:id
func (h *ProgramHandler) GetProgram(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid program ID",
		})
	}

	program, err := h.programService.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "program not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Program not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, program)
}

// ListPrograms handles GET /api/v1/programs
func (h *ProgramHandler) ListPrograms(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	programs, err := h.programService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"programs": programs,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateProgram handles PUT /api/v1/programs/:id
func (h *ProgramHandler) UpdateProgram(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid program ID",
		})
	}

	var req admin.UpdateProgramRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation_error",
			"message": "name is required",
		})
	}

	program, err := h.programService.Update(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "program not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Program not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, program)
}

// DeleteProgram handles DELETE /api/v1/programs/:id
func (h *ProgramHandler) DeleteProgram(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid program ID",
		})
	}

	if err := h.programService.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "program not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Program not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
