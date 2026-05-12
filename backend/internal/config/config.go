package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// resolveSecret retorna o valor de uma variável de ambiente, com suporte ao
// padrão Docker `<NAME>_FILE` que aponta para um arquivo contendo o segredo
// (ex.: montado a partir de um Docker Swarm secret em `/run/secrets/...`).
// Se ambos `NAME` e `NAME_FILE` estiverem definidos, `NAME_FILE` tem precedência.
func resolveSecret(key string) string {
	if path, ok := os.LookupEnv(key + "_FILE"); ok && path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimRight(string(data), "\r\n")
		}
	}
	return os.Getenv(key)
}

// Config agrega toda a configuração de runtime carregada de variáveis de
// ambiente (com fallback opcional para um arquivo .env).
type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
	CORS CORSConfig
	Log LogConfig
}

// AppConfig contém parâmetros gerais do servidor HTTP.
type AppConfig struct {
	Port string
	Env  string
}

// DBConfig agrupa as DSNs dos quatro bancos legados e os parâmetros do pool.
type DBConfig struct {
	SmidDSN          string
	PermissionDSN    string
	LogDSN           string
	CommunicationDSN string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig contém os parâmetros de emissão e validação de tokens.
type JWTConfig struct {
	Secret               string
	Expiration           time.Duration
	RefreshExpiration    time.Duration
}

// CORSConfig contém as origens permitidas.
type CORSConfig struct {
	AllowedOrigins []string
}

// LogConfig contém os parâmetros do logger estruturado.
type LogConfig struct {
	Level  string
	Format string
}

// Load carrega configuração a partir do ambiente. Se um arquivo .env existir
// no diretório de execução ele é carregado primeiro (sem sobrescrever valores
// já definidos no ambiente).
func Load() (*Config, error) {
	_ = godotenv.Load() // arquivo .env é opcional em produção

	cfg := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			SmidDSN:          resolveSecret("DB_SMID_DSN"),
			PermissionDSN:    resolveSecret("DB_PERMISSION_DSN"),
			LogDSN:           resolveSecret("DB_LOG_DSN"),
			CommunicationDSN: resolveSecret("DB_COMMUNICATION_DSN"),
			MaxOpenConns:     getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:     getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime:  time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		},
		JWT: JWTConfig{
			Secret:            resolveSecret("JWT_SECRET"),
			Expiration:        time.Duration(getEnvInt("JWT_EXPIRATION_HOURS", 8)) * time.Hour,
			RefreshExpiration: time.Duration(getEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 168)) * time.Hour,
		},
		CORS: CORSConfig{
			AllowedOrigins: splitAndTrim(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	missing := []string{}
	if c.DB.SmidDSN == "" {
		missing = append(missing, "DB_SMID_DSN")
	}
	if c.DB.PermissionDSN == "" {
		missing = append(missing, "DB_PERMISSION_DSN")
	}
	if c.DB.LogDSN == "" {
		missing = append(missing, "DB_LOG_DSN")
	}
	if c.DB.CommunicationDSN == "" {
		missing = append(missing, "DB_COMMUNICATION_DSN")
	}
	if c.JWT.Secret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return fmt.Errorf("config: variáveis de ambiente obrigatórias ausentes: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
