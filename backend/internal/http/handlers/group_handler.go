package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	groupService *admin.GroupService
}

func NewGroupHandler(groupService *admin.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup handles POST /api/v1/groups
func (h *GroupHandler) CreateGroup(c echo.Context) error {
	var req admin.CreateGroupRequest
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

	group, err := h.groupService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, group)
}

// GetGroup handles GET /api/v1/groups/:id
func (h *GroupHandler) GetGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid group ID",
		})
	}

	group, err := h.groupService.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "group not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Group not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, group)
}

// ListGroups handles GET /api/v1/groups
func (h *GroupHandler) ListGroups(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := h.groupService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"groups": groups,
		"limit": limit,
		"offset": offset,
	})
}

// UpdateGroup handles PUT /api/v1/groups/:id
func (h *GroupHandler) UpdateGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid group ID",
		})
	}

	var req admin.UpdateGroupRequest
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

	group, err := h.groupService.Update(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "group not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Group not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, group)
}

// DeleteGroup handles DELETE /api/v1/groups/:id
func (h *GroupHandler) DeleteGroup(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_id",
			"message": "Invalid group ID",
		})
	}

	if err := h.groupService.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "group not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "not_found",
				"message": "Group not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
