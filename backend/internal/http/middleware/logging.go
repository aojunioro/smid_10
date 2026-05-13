package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger é um middleware que loga informações estruturadas de cada request:
// - request_id (injetado pelo middleware.RequestID do Echo)
// - método HTTP
// - path
// - status code
// - latência em milissegundos
// - IP do cliente
func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Extrai request_id do header injetado pelo middleware.RequestID
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = "unknown"
			}

			// Executa o handler
			err := next(c)

			// Calcula latência
			latency := time.Since(start).Milliseconds()

			// Log estruturado
			logger.Info("http_request",
				"request_id", requestID,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency_ms", latency,
				"client_ip", c.RealIP(),
			)

			return err
		}
	}
}
