package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
)

// AuthHandler gerencia endpoints de autenticação.
type AuthHandler struct {
	authService *admin.AuthService
}

// NewAuthHandler cria uma nova instância de AuthHandler.
func NewAuthHandler(authService *admin.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest representa o corpo da requisição de login.
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa a resposta de login bem-sucedido.
type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt string      `json:"expires_at"`
	User      admin.SystemUser `json:"user"`
}

// Login autentica um usuário e retorna um token JWT.
// @Summary Login
// @Description Autentica um usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciais de login"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validar campos obrigatórios
	if req.Login == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Login and password are required",
		})
	}

	// Chamar serviço de autenticação
	loginReq := admin.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	}

	resp, err := h.authService.Login(c.Request().Context(), loginReq)
	if err != nil {
		if err == admin.ErrInvalidCredentials || err == admin.ErrUserNotFound {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid login or password",
			})
		}
		if err == admin.ErrUserInactive {
			return c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "user_inactive",
				Message: "User account is inactive",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to authenticate user",
		})
	}

	// Converter para formato de resposta
	loginResp := LoginResponse{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:      resp.User,
	}

	return c.JSON(http.StatusOK, loginResp)
}

// ErrorResponse representa uma resposta de erro.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
