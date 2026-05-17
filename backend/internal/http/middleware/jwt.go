package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
)

// ClaimsContextKey is the Echo context key for JWT claims.
const ClaimsContextKey = "claims"

// JWT validates Bearer tokens issued by AuthService.
func JWT(authService *admin.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Authorization header required",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "unauthorized",
					"message": "Invalid Authorization header format",
				})
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "invalid_token",
					"message": "Invalid or expired token",
				})
			}

			c.Set(ClaimsContextKey, claims)
			return next(c)
		}
	}
}

// GetClaims returns JWT claims stored by the JWT middleware.
func GetClaims(c echo.Context) (*admin.Claims, bool) {
	claims, ok := c.Get(ClaimsContextKey).(*admin.Claims)
	return claims, ok
}
