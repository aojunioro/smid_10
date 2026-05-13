package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/aojunioro/smid_10/backend/internal/config"
	"github.com/aojunioro/smid_10/backend/internal/db"
	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/aojunioro/smid_10/backend/internal/domain/common"
	"github.com/aojunioro/smid_10/backend/internal/http/handlers"
	"github.com/aojunioro/smid_10/backend/internal/http/middleware"
)

// SetupRouter configura todas as rotas da aplicação.
func SetupRouter(e *echo.Echo, cfg *config.Config, pools *db.Pools, logger *slog.Logger) {
	// Middlewares globais
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Middleware de logging estruturado customizado
	e.Use(middleware.RequestLogger(logger))

	// Grupo de rotas da API v1
	v1 := e.Group("/api/v1")

	// Rotas de autenticação (públicas)
	setupAuthRoutes(v1, cfg, pools)

	// Rotas protegidas (requerem autenticação JWT)
	setupProtectedRoutes(v1, pools)
}

// setupAuthRoutes configura rotas de autenticação.
func setupAuthRoutes(v1 *echo.Group, cfg *config.Config, pools *db.Pools) {
	// Obter pool de conexão do banco permission
	permissionDB, err := pools.Get(db.AliasPermission)
	if err != nil {
		panic(err)
	}

	// Criar repositório de usuário
	userRepo := admin.NewUserRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de autenticação
	authService := admin.NewAuthService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
	)

	// Criar handler de autenticação
	authHandler := handlers.NewAuthHandler(authService)

	// Registrar rotas de autenticação
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)
}

// setupProtectedRoutes configura rotas protegidas por autenticação JWT.
func setupProtectedRoutes(v1 *echo.Group, pools *db.Pools) {
	// Rotas protegidas serão adicionadas aqui
	// Exemplo:
	// protected := v1.Group("")
	// protected.Use(echomiddleware.JWT(cfg.JWT.Secret))
	// protected.GET("/users", userHandler.List)
}
