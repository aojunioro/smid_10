# HANDOFF_BOOTSTRAP.md — Abertura do SMID 10 no Windsurf

## 1. Metadados

- **Tema**: bootstrap inicial do SMID 10
- **Data/Hora**: 2026-05-12
- **Thread de origem**: Cascade no workspace `smid_8_ajustes-gerais`
- **Responsável**: aojunioro
- **Repositório**: `git@github.com:aojunioro/smid_10.git`

---

## 2. Objetivo do Projeto

Reimplementar o SMID 8.x (Adianti/PHP) em **Go (backend) + Next.js/shadcn (frontend)** reaproveitando os bancos MySQL legados (`smid`, `permission`, `log`, `communication`) para coexistência transparente com o sistema atual até o cutover.

---

## 3. Decisões Já Fechadas

1. **Cenário 2** — Go (API REST) + Next.js (SPA separada) → `docs/adrs/0001-cenario-2-go-next.md`
2. **Reuso da base legada** — schema legado intocável durante a coexistência → `docs/adrs/0002-reuso-base-legado.md`
3. **Stack detalhada** — Echo, sqlc, JWT v5, Next 14, shadcn/ui, TailwindCSS, TanStack Query/Table, dnd-kit → `docs/adrs/0003-stack-detalhada.md`
4. **Postgres como destino pós-cutover** (Fase 8); preparação contínua em medidas baratas → `docs/adrs/0004-postgres-pos-cutover.md`
5. **Ambiente de testes em Docker Swarm na VPS compartilhada** (Traefik + `SmydiNet` + MariaDB 10.11) → `docs/adrs/0005-ambiente-testes-swarm-vps.md`
6. **Template inicial recomendado**: `shadcn-admin` (https://github.com/satnaing/shadcn-admin) para o frontend

---

## 4. Contexto Técnico Essencial

- **Origem**: derivado do workspace `smid_8_ajustes-gerais` (commit `a22f3d81` em 12/05/2026)
- **Worktree do SMID 8.x para referência de código PHP**: `/Users/antoniojr/ProjetoTrae/smid_8` (consultar somente para esclarecer ambiguidades — a fonte de verdade é o SPEC)
- **Fonte de verdade do projeto**: os 24 SPECs em `docs/specs/`
- **Restrições herdadas**: ver `docs/legacy-reference/SMID_PERSONALIZACOES.md`

---

## 5. Arquivos Importantes (ordem de leitura)

1. `README.md` — visão, 8 fases de evolução, stack, roadmap
2. `AGENTS.md` — regras de operação para agentes de IA
3. `docs/specs/SPEC_INDEX.md` — índice mestre dos 24 SPECs com checklist de reimplementação
4. `docs/specs/SPEC_DATABASE.md` — arquitetura multi-banco
5. `docs/specs/SPEC_REST_API.md` — contrato da API REST que o backend implementa
6. `docs/specs/SPEC_UX_UI.md` — padrões de UX que o frontend deve preservar
7. `docs/adrs/0001..0003` — decisões arquiteturais fechadas
8. `backend/README.md` — bootstrap do Go
9. `frontend/README.md` — bootstrap do Next

---

## 6. O Que Já Foi Concluído

### Bootstrap inicial
- Estrutura de pastas (`backend/`, `frontend/`, `docs/`, `deploy/`)
- 24 SPECs canônicos copiados para `docs/specs/`
- Schema base do legado copiado para `docs/legacy-schema/`
- Docs críticos de UX/UI copiados para `docs/legacy-reference/`
- **5 ADRs** escritos (0001-0005)
- `README.md`, `AGENTS.md` e `.gitignore` criados
- `backend/.env.example` com DSN dos 4 bancos
- Repositório git inicializado, commit inicial e push para `origin/main`

### Fase 0.1 — Backend bootstrap (concluída em 2026-05-12)
- `go.mod` inicializado (`github.com/aojunioro/smid_10/backend`, Go 1.26.3)
- `backend/cmd/server/main.go` — Echo v4, middlewares (`RequestID`, `Recover`, `CORS`), `GET /healthz`, graceful shutdown via SIGINT/SIGTERM, logger `slog` JSON
- `backend/internal/config/config.go` — `Load()` via `godotenv`, validação de DSNs/JWT obrigatórios, **suporte a Docker Secrets** via sufixo `_FILE`
- `backend/internal/db/pools.go` — 4 pools nomeados (smid/permission/log/communication) com `PingAll` e timeout por alias
- `go build ./...` e `go vet ./...` verdes

### Fase 0 — Ambiente de testes (concluída em 2026-05-12)
- `backend/Dockerfile` multi-stage (distroless, ≈ 15 MB)
- `backend/.dockerignore`
- `deploy/swarm-stack.yml` — stack `smid10` (MariaDB + API; web comentado)
- `deploy/mariadb-init/01-schemas.sql` — 4 schemas + usuário `smid10`
- `deploy/README.md` — runbook operacional completo
- `.github/workflows/backend-build.yml` — CI publica em `ghcr.io/aojunioro/smid10-api`
- `.gitignore` reforçado: `credenciais.md`, `secrets/`, `*.pem`, `*.key`, `deploy/.env*`

---

## 7. Pendências Abertas

### Deploy do ambiente de testes (pré-Fase 1) — concluído em 2026-05-13

1. ~~Trocar a senha root da VPS e configurar SSH por chave (segurança)~~
2. Criar registros DNS A:
   - `api.s10.smydi.com.br` → `216.144.235.25` (pendente)
   - `s10.smydi.com.br` → `216.144.235.25` (pendente)
3. ~~Push do código para `origin/main` para acionar o build no CI~~
4. ~~Na VPS: criar os 8 Swarm secrets (`deploy/README.md §3`)~~
5. ~~Na VPS: popular o volume `smid10_mariadb_init` com `01-schemas.sql` (`deploy/README.md §4`)~~
6. ~~Na VPS: `docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth`~~
7. ~~Validar `https://api.s10.smydi.com.br/healthz` retornando `status: "ok"`~~ (smoke-test interno OK, externo pendente DNS)

**Nota**: O volume `smid10_mariadb_data` já existia antes do deploy, então os scripts de init não foram executados automaticamente. Os bancos foram criados manualmente via SQL direto no container.

### Backend (Fase 0.2+)

1. Middleware de logging estruturado por requisição (request_id, latência, status)
2. Repositórios atrás de interfaces (preparação para ADR 0004 — Postgres pós-cutover)
3. Setup de testes com `testcontainers` para repositórios

### Frontend (Fase 0.1 do frontend)

1. Scaffold do `shadcn-admin` em `frontend/` (opção A do `frontend/README.md`)
2. Configurar `lib/api/client.ts` apontando para `https://api.s10.smydi.com.br` (staging) com fallback `localhost:8080`
3. Setar tema dark/light e branding
4. Descomentar o serviço `web` em `deploy/swarm-stack.yml` quando a imagem `smid10-web` estiver publicada

---

## 8. Próximo Passo Executável

> **Deploy do stack `smid10` na VPS**: seguir `deploy/README.md` seções §2 (DNS), §3 (secrets), §4 (volume init) e §6 (deploy). Validar `/healthz` interno e externo conforme §7. Em paralelo, iniciar o scaffold do frontend (`frontend/README.md`).

---

## 9. Riscos e Cuidados

1. **Não alterar schema legado** — qualquer DDL que toque tabelas existentes precisa de novo ADR
2. **Compatibilidade de hash de senha** — o legado usa o que o Adianti gerar; validar exatamente o mesmo algoritmo no Go (provavelmente bcrypt; confirmar no primeiro login)
3. **Encoding/collation** — algumas tabelas legadas têm mojibake (`latin1` com UTF-8 dentro); usar `charset=utf8mb4` no DSN e revisar caso a caso
4. **Sessões PHP do SMID 8.x são independentes** do JWT do SMID 10 — usuário pode estar logado em um e não no outro durante a coexistência
5. **Filtro de unidade**: toda query do SMID 10 deve respeitar `unidd_id` e `excluido_em IS NULL`
6. **RAM apertada na VPS** — ~1 GB livre; nunca remover `deploy.resources.limits.memory` do `swarm-stack.yml` ou outros stacks de produção podem cair (Chatwoot, Evolution, Typebot)
7. **Credenciais SSH** — o arquivo `credenciais.md` está gitignorado mas a senha root SSH passou pelo histórico do agente; rotacionar e migrar para chave SSH é recomendado

---

## 10. Prompt de Retomada (copiar e colar)

```md
Continuar o bootstrap do SMID 10.

Leia primeiro, nesta ordem:
- README.md
- AGENTS.md
- docs/handoff/HANDOFF_BOOTSTRAP.md
- docs/specs/SPEC_INDEX.md
- docs/adrs/0001-cenario-2-go-next.md
- docs/adrs/0002-reuso-base-legado.md
- docs/adrs/0003-stack-detalhada.md

Estado atual: bootstrap concluído (governança, SPECs, ADRs, READMEs).
Nenhum código Go ou TS implementado ainda.

Próximo passo: Fase 0.1 do backend.
- Inicializar go.mod
- Criar cmd/server/main.go com Echo + GET /healthz
- Criar internal/config/config.go (.env via godotenv)
- Criar internal/db/pools.go com pools para smid, permission, log, communication
- Estender /healthz para validar ping nos 4 pools

Respeitar as regras de AGENTS.md. Não alterar schema legado.
```

---

## 11. Comandos Úteis

```bash
# Clonar (em outra máquina)
git clone git@github.com:aojunioro/smid_10.git

# Verificar status
git status
git log --oneline -10

# Backend (após Fase 0.1)
cd backend
cp .env.example .env   # editar credenciais
go run ./cmd/server

# Frontend (após Fase 0.1)
cd frontend
pnpm install
pnpm dev
```

---

**Status**: Fase 0.1 do backend concluída; ambiente de deploy preparado; aguardando deploy efetivo na VPS e início do scaffold do frontend
