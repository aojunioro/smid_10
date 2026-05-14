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
	logger.Info("SetupRouter iniciado")

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
	setupAuthRoutes(v1, cfg, pools, logger)

	// Rotas protegidas (requerem autenticação JWT)
	setupProtectedRoutes(v1, pools)
}

// setupAuthRoutes configura rotas de autenticação.
func setupAuthRoutes(v1 *echo.Group, cfg *config.Config, pools *db.Pools, logger *slog.Logger) {
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

	// Log para debug
	logger.Info("rotas de autenticação registradas", "route", "/api/v1/auth/login")
}

// setupProtectedRoutes configura rotas protegidas por autenticação JWT.
func setupProtectedRoutes(v1 *echo.Group, pools *db.Pools) {
	// Obter pool de conexão do banco permission
	permissionDB, err := pools.Get(db.AliasPermission)
	if err != nil {
		panic(err)
	}

	// Criar repositório de usuário
	userRepo := admin.NewUserRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de usuário
	userService := admin.NewUserService(userRepo)

	// Criar handler de usuário
	userHandler := handlers.NewUserHandler(userService)

	// Rotas de usuários (protegidas por JWT)
	users := v1.Group("/users")
	users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUser)
	users.POST("", userHandler.CreateUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
