package handlers

import (
	"net/http"
	"strconv"

	"github.com/aojunioro/smid_10/backend/internal/domain/log"
	"github.com/labstack/echo/v4"
)

type AccessLogHandler struct {
	accessLogService *log.AccessLogService
}

func NewAccessLogHandler(accessLogService *log.AccessLogService) *AccessLogHandler {
	return &AccessLogHandler{
		accessLogService: accessLogService,
	}
}

// ListAccessLogs handles GET /api/v1/logs/access
func (h *AccessLogHandler) ListAccessLogs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.accessLogService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_logs": logs,
		"limit": limit,
		"offset": offset,
	})
}

type ChangeLogHandler struct {
	changeLogService *log.ChangeLogService
}

func NewChangeLogHandler(changeLogService *log.ChangeLogService) *ChangeLogHandler {
	return &ChangeLogHandler{
		changeLogService: changeLogService,
	}
}

// ListChangeLogs handles GET /api/v1/logs/change
func (h *ChangeLogHandler) ListChangeLogs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.changeLogService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"change_logs": logs,
		"limit": limit,
		"offset": offset,
	})
}

type SqlLogHandler struct {
	sqlLogService *log.SqlLogService
}

func NewSqlLogHandler(sqlLogService *log.SqlLogService) *SqlLogHandler {
	return &SqlLogHandler{
		sqlLogService: sqlLogService,
	}
}

// ListSqlLogs handles GET /api/v1/logs/sql
func (h *SqlLogHandler) ListSqlLogs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.sqlLogService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sql_logs": logs,
		"limit": limit,
		"offset": offset,
	})
}

type RequestLogHandler struct {
	requestLogService *log.RequestLogService
}

func NewRequestLogHandler(requestLogService *log.RequestLogService) *RequestLogHandler {
	return &RequestLogHandler{
		requestLogService: requestLogService,
	}
}

// ListRequestLogs handles GET /api/v1/logs/request
func (h *RequestLogHandler) ListRequestLogs(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.requestLogService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"request_logs": logs,
		"limit": limit,
		"offset": offset,
	})
}
