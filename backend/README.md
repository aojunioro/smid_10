# Backend SMID 10 (Go)

API REST do SMID 10 escrita em Go.

---

## 1. Stack

Ver `docs/adrs/0003-stack-detalhada.md`.

Resumo:
- Go 1.22+
- Echo v4
- sqlc + go-sql-driver/mysql
- golang-jwt/jwt v5
- log/slog (stdlib)
- bcrypt

---

## 2. Estrutura

```
backend/
├── cmd/
│   └── server/
│       └── main.go              ← entrypoint
├── internal/
│   ├── config/                  ← carregamento .env + struct Config
│   │   └── config.go
│   ├── db/                      ← pools por alias (smid, permission, log, communication)
│   │   ├── pools.go
│   │   └── queries/             ← arquivos .sql para sqlc
│   ├── auth/                    ← JWT, bcrypt, claims
│   │   ├── jwt.go
│   │   └── password.go
│   ├── domain/                  ← um pacote por SPEC
│   │   ├── leads/
│   │   ├── visitas/
│   │   ├── pedidos/
│   │   └── ...
│   ├── http/
│   │   ├── handlers/            ← handlers REST por domínio
│   │   ├── middleware/          ← auth, cors, logging, recovery
│   │   └── routes.go            ← registro de rotas
│   └── observability/           ← logging, métricas (futuro)
├── pkg/                         ← código reutilizável fora do internal
├── migrations/                  ← golang-migrate (deltas SMID 10)
├── sqlc.yaml                    ← config do sqlc
├── go.mod
├── go.sum
└── .env.example
```

---

## 3. Bootstrap (Fase 0)

### 3.1 Inicializar módulo Go

```bash
cd backend
go mod init github.com/aojunioro/smid_10/backend
```

### 3.2 Dependências iniciais

```bash
go get github.com/labstack/echo/v4
go get github.com/go-sql-driver/mysql
go get github.com/golang-jwt/jwt/v5
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc
go get -tool github.com/golang-migrate/migrate/v4/cmd/migrate
```

### 3.3 Variáveis de ambiente

Copiar `.env.example` para `.env` e ajustar:

```bash
APP_PORT=8080
APP_ENV=development

DB_SMID_DSN=user:pass@tcp(127.0.0.1:3306)/smid?parseTime=true&loc=Local
DB_PERMISSION_DSN=user:pass@tcp(127.0.0.1:3306)/permission?parseTime=true&loc=Local
DB_LOG_DSN=user:pass@tcp(127.0.0.1:3306)/log?parseTime=true&loc=Local
DB_COMMUNICATION_DSN=user:pass@tcp(127.0.0.1:3306)/communication?parseTime=true&loc=Local

JWT_SECRET=change-me
JWT_EXPIRATION_HOURS=8

CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 3.4 Rodar

```bash
go run ./cmd/server
```

---

## 4. Convenções

Ver `AGENTS.md` na raiz do projeto, seção 4.

Resumo:
- Um pacote por domínio
- Handlers finos, serviços de domínio gordos
- Transações explícitas por alias
- Erros tipados; nunca expor `error.Error()` cru
- `context.Context` propagado em todas as camadas
- Testes unitários por pacote

---

## 5. Roadmap do Backend

| Fase | Entregável |
|------|-----------|
| 0.1 | `main.go` com Echo + `GET /healthz` |
| 0.2 | Pools de conexão MySQL para os 4 aliases |
| 0.3 | Middleware de logging estruturado |
| 1.1 | `POST /auth/login` validando contra `permission.system_users` |
| 1.2 | JWT emitido + `GET /auth/me` |
| 1.3 | Middleware de autenticação |
| 1.4 | Multi-unidade na sessão |
| 2.1 | `GET /leads` e `GET /leads/{id}` |
| ... | (ver SPECs por domínio) |

---

## 6. Testes

```bash
go test ./...                   # todos
go test -race ./...             # com race detector
go test -cover ./...            # cobertura
```

---

## 7. Build

```bash
go build -o bin/smid10-server ./cmd/server
```

Binário único, sem dependências externas além das libs C do `go-sql-driver` (CGO desabilitável).

Para build estático:

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/smid10-server ./cmd/server
```
