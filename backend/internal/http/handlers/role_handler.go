package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	roleService *admin.RoleService
}

func NewRoleHandler(roleService *admin.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole handles POST /api/v1/roles
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req admin.CreateRoleRequest
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

	role, err := h.roleService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, role)
}

// GetRole handles GET /api/v1/roles/:id
func (h *RoleHandler) GetRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid role ID",
		})
	}

	role, err := h.roleService.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "role not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, role)
}

// ListRoles handles GET /api/v1/roles
func (h *RoleHandler) ListRoles(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	roles, err := h.roleService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": roles,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateRole handles PUT /api/v1/roles/:id
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid role ID",
		})
	}

	var req admin.UpdateRoleRequest
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

	role, err := h.roleService.Update(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "role not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, role)
}

// DeleteRole handles DELETE /api/v1/roles/:id
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid role ID",
		})
	}

	if err := h.roleService.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "role not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
