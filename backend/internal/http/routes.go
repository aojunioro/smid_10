package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/aojunioro/smid_10/backend/internal/config"
	"github.com/aojunioro/smid_10/backend/internal/db"
	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
	"github.com/aojunioro/smid_10/backend/internal/domain/common"
	"github.com/aojunioro/smid_10/backend/internal/domain/communication"
	"github.com/aojunioro/smid_10/backend/internal/domain/leads"
	"github.com/aojunioro/smid_10/backend/internal/domain/log"
	"github.com/aojunioro/smid_10/backend/internal/domain/tarefas"
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

	// Criar repositório de grupo
	groupRepo := admin.NewGroupRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de grupo
	groupService := admin.NewGroupService(groupRepo)

	// Criar handler de grupo
	groupHandler := handlers.NewGroupHandler(groupService)

	// Rotas de grupos (protegidas por JWT)
	groups := v1.Group("/groups")
	groups.GET("", groupHandler.ListGroups)
	groups.GET("/:id", groupHandler.GetGroup)
	groups.POST("", groupHandler.CreateGroup)
	groups.PUT("/:id", groupHandler.UpdateGroup)
	groups.DELETE("/:id", groupHandler.DeleteGroup)

	// Criar repositório de papel
	roleRepo := admin.NewRoleRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de papel
	roleService := admin.NewRoleService(roleRepo)

	// Criar handler de papel
	roleHandler := handlers.NewRoleHandler(roleService)

	// Rotas de papéis (protegidas por JWT)
	roles := v1.Group("/roles")
	roles.GET("", roleHandler.ListRoles)
	roles.GET("/:id", roleHandler.GetRole)
	roles.POST("", roleHandler.CreateRole)
	roles.PUT("/:id", roleHandler.UpdateRole)
	roles.DELETE("/:id", roleHandler.DeleteRole)

	// Criar repositório de programa
	programRepo := admin.NewProgramRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de programa
	programService := admin.NewProgramService(programRepo)

	// Criar handler de programa
	programHandler := handlers.NewProgramHandler(programService)

	// Rotas de programas (protegidas por JWT)
	programs := v1.Group("/programs")
	programs.GET("", programHandler.ListPrograms)
	programs.GET("/:id", programHandler.GetProgram)
	programs.POST("", programHandler.CreateProgram)
	programs.PUT("/:id", programHandler.UpdateProgram)
	programs.DELETE("/:id", programHandler.DeleteProgram)

	// Criar repositório de unidade
	unitRepo := admin.NewUnitRepository(permissionDB, common.DBAlias(db.AliasPermission))

	// Criar serviço de unidade
	unitService := admin.NewUnitService(unitRepo)

	// Criar handler de unidade
	unitHandler := handlers.NewUnitHandler(unitService)

	// Rotas de unidades (protegidas por JWT)
	units := v1.Group("/units")
	units.GET("", unitHandler.ListUnits)
	units.GET("/:id", unitHandler.GetUnit)
	units.POST("", unitHandler.CreateUnit)
	units.PUT("/:id", unitHandler.UpdateUnit)
	units.DELETE("/:id", unitHandler.DeleteUnit)

	// Obter pool de conexão do banco log
	logDB, err := pools.Get(db.AliasLog)
	if err != nil {
		panic(err)
	}

	// Criar repositórios de log
	accessLogRepo := log.NewAccessLogRepository(logDB, common.DBAlias(db.AliasLog))
	changeLogRepo := log.NewChangeLogRepository(logDB, common.DBAlias(db.AliasLog))
	sqlLogRepo := log.NewSqlLogRepository(logDB, common.DBAlias(db.AliasLog))

	// Criar serviços de log
	accessLogService := log.NewAccessLogService(accessLogRepo)
	changeLogService := log.NewChangeLogService(changeLogRepo)
	sqlLogService := log.NewSqlLogService(sqlLogRepo)

	// Criar handlers de log
	accessLogHandler := handlers.NewAccessLogHandler(accessLogService)
	changeLogHandler := handlers.NewChangeLogHandler(changeLogService)
	sqlLogHandler := handlers.NewSqlLogHandler(sqlLogService)

	// Rotas de logs (protegidas por JWT)
	logs := v1.Group("/logs")
	logs.GET("/access", accessLogHandler.ListAccessLogs)
	logs.GET("/change", changeLogHandler.ListChangeLogs)
	logs.GET("/sql", sqlLogHandler.ListSqlLogs)

	// Obter pool de conexão do banco communication
	commDB, err := pools.Get(db.AliasCommunication)
	if err != nil {
		panic(err)
	}

	// Criar repositórios de communication
	notificationRepo := communication.NewNotificationRepository(commDB, common.DBAlias(db.AliasCommunication))
	messageRepo := communication.NewMessageRepository(commDB, common.DBAlias(db.AliasCommunication))

	// Criar serviços de communication
	notificationService := communication.NewNotificationService(notificationRepo)
	messageService := communication.NewMessageService(messageRepo)

	// Criar handlers de communication
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	messageHandler := handlers.NewMessageHandler(messageService)

	// Rotas de comunicação (protegidas por JWT)
	notifications := v1.Group("/notifications")
	notifications.GET("", notificationHandler.ListNotifications)
	notifications.GET("/:id", notificationHandler.GetNotification)
	notifications.POST("", notificationHandler.CreateNotification)
	notifications.PUT("/:id", notificationHandler.UpdateNotification)

	messages := v1.Group("/messages")
	messages.GET("", messageHandler.ListMessages)
	messages.GET("/:id", messageHandler.GetMessage)
	messages.POST("", messageHandler.CreateMessage)
	messages.PUT("/:id", messageHandler.UpdateMessage)

	// Obter pool de conexão do banco smid
	smidDB, err := pools.Get(db.AliasSmid)
	if err != nil {
		panic(err)
	}

	// Criar repositório de tarefa
	tarefaRepo := tarefas.NewTarefaRepository(smidDB, common.DBAlias(db.AliasSmid))

	// Criar serviço de tarefa
	tarefaService := tarefas.NewTarefaService(tarefaRepo)

	// Criar handler de tarefa
	tarefaHandler := handlers.NewTarefaHandler(tarefaService)

	// Rotas de tarefas (protegidas por JWT)
	tarefas := v1.Group("/tarefas")
	tarefas.GET("", tarefaHandler.ListTarefas)
	tarefas.GET("/:id", tarefaHandler.GetTarefa)
	tarefas.POST("", tarefaHandler.CreateTarefa)
	tarefas.PUT("/:id", tarefaHandler.UpdateTarefa)

	// Criar repositório de lead
	leadRepo := leads.NewLeadRepository(smidDB, common.DBAlias(db.AliasSmid))

	// Criar serviço de lead
	leadService := leads.NewLeadService(leadRepo)

	// Criar handler de lead
	leadHandler := handlers.NewLeadHandler(leadService)

	// Rotas de leads (protegidas por JWT)
	leads := v1.Group("/leads")
	leads.GET("", leadHandler.ListLeads)
	leads.GET("/:id", leadHandler.GetLead)
	leads.POST("", leadHandler.CreateLead)
	leads.PUT("/:id", leadHandler.UpdateLead)
}
