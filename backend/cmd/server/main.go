package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/aojunioro/smid_10/backend/internal/config"
	"github.com/aojunioro/smid_10/backend/internal/db"
)

func main() {
	logger := newLogger()
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("falha ao carregar configuração", "err", err)
		os.Exit(1)
	}

	pools, err := db.Open(cfg.DB)
	if err != nil {
		logger.Error("falha ao abrir pools de banco", "err", err)
		os.Exit(1)
	}
	defer pools.Close()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
	}))

	e.GET("/healthz", healthHandler(pools))

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("servidor iniciado", "addr", srv.Addr, "env", cfg.App.Env)
		if err := e.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("servidor falhou", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("encerrando servidor")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("erro no shutdown", "err", err)
	}
}

// healthHandler retorna 200 se todos os pools responderem ao ping, ou 503
// caso qualquer um falhe. A resposta lista o status individual de cada alias.
func healthHandler(pools *db.Pools) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		results := pools.PingAll(ctx)
		databases := make(map[string]string, len(results))
		overallOK := true
		for alias, err := range results {
			if err != nil {
				overallOK = false
				databases[string(alias)] = "down: " + err.Error()
				continue
			}
			databases[string(alias)] = "ok"
		}

		status := "ok"
		httpCode := http.StatusOK
		if !overallOK {
			status = "degraded"
			httpCode = http.StatusServiceUnavailable
		}

		return c.JSON(httpCode, map[string]any{
			"status":    status,
			"databases": databases,
			"time":      time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	opts := &slog.HandlerOptions{Level: level}
	if os.Getenv("LOG_FORMAT") == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
