# AGENTS.md — SMID 10

Guia canonico para agentes de IA e desenvolvedores no SMID 10. Em conflito com qualquer outro documento do repositorio, **este prevalece**.

---

## 1. Contexto rapido

| Item | Valor |
|------|--------|
| Legado | SMID 8.x (Adianti/PHP) |
| SMID 10 | **Go (API REST)** + **Vite/React (SPA)** + shadcn/ui |
| Bancos | 4 aliases MySQL: `smid`, `permission`, `log`, `communication` |
| Coexistencia | SMID 8.x e SMID 10 ate cutover; schema legado **intocavel** |
| Fonte de verdade | `docs/specs/` (24 SPECs) + ADRs em `docs/adrs/` |
| Continuidade entre sessoes | `docs/handoff/CURRENT.md` (indice: `docs/handoff/README.md`) |
| Referencia PHP (somente leitura) | worktree `smid_8` — ver handoff; nunca copiar controllers verbatim |

**URLs de staging (quando deploy ativo):**

- API: `https://api.s10.smydi.com.br`
- App: `https://s10.smydi.com.br`
- Health: `GET /healthz` (4 pools)

**Nota de stack:** ADR 0003 cita Next.js; o frontend atual e **Vite + TanStack Router** (`frontend/package.json`). Novas telas seguem a estrutura real em `frontend/src/`, nao `app/`. Mudanca para Next exige ADR.

---

## 2. Harness do agente (Cursor)

Usar estes artefatos antes de codar:

| Artefato | Caminho | Uso |
|----------|---------|-----|
| Skills (projeto) | `.cursor/skills/*/SKILL.md` | Workflows: SPEC, frontend, handoff, deploy, invariantes |
| Regras (projeto) | `.cursor/rules/*.mdc` | Padroes por tipo de arquivo (Go, TS, SPECs) |
| MCP local | `.cursor/mcp.json` | Ferramentas `smid-tools` + docs |
| MCP opcionais | `.cursor/mcp.json.example` | GitHub, MySQL (requer tokens locais) |
| Workflows Windsurf | `.windsurf/workflows/*.md` | Equivalentes legados; preferir skills no Cursor |

**Leitura obrigatoria por tipo de tarefa:**

| Tarefa | Ler primeiro |
|--------|----------------|
| Qualquer feature | `docs/specs/SPEC_INDEX.md` + SPEC do dominio |
| Backend dominio | Skill `smid-spec-implement` + `AGENTS.md` ss 4 e 6 |
| Frontend | Skill `smid-frontend-feature` + `docs/specs/SPEC_UX_UI.md` |
| Deploy VPS | Skill `smid-deploy-vps` + `deploy/README.md` |
| Fim de sessao | Skill `smid-handoff-update` |
| Seguranca/auth | `docs/specs/SPEC_ADMIN.md` + `SPEC_REST_API.md` |

---

## 3. Estado atual do codigo (espelho de `docs/handoff/CURRENT.md`)

Resumo para nao supor maturidade inexistente:

| Camada | Estado |
|--------|--------|
| Backend | CRUD REST amplo por dominio; **middleware JWT nas rotas protegidas ainda pendente** (`ValidateToken` existe, nao aplicado em `routes.go`) |
| Regras de negocio | Maioria dos dominios em CRUD raso; invariantes INV-001+ nao cobertos por testes |
| Dominios sem API | `SPEC_METAS`, `SPEC_RELATORIOS`, `SPEC_INTEGRACOES_JOBS` |
| Frontend | Template shadcn-admin; **login mock**; `lib/api/client.ts` sem uso nas features |
| Testes Go | Quase so `leads/repository_test.go` |
| sqlc | Previsto no ADR; repositorios atuais usam SQL manual |

**Prioridade de produto recomendada:** (1) JWT + login real, (2) vertical slice Leads UI+API com regras do SPEC, (3) sync Lead-Visita, (4) demais dominios.

---

## 4. Principios nao negociaveis

1. **SPECs primeiro** — `docs/specs/SPEC_INDEX.md` antes de implementar
2. **Schema legado intocavel** — novas tabelas `s10_*` ou schema dedicado; sem ALTER em tabelas legadas (ADR 0002)
3. **Multi-banco** — transacoes por alias; nunca cruzar TX entre `smid` e `permission`
4. **Invariantes** — `SPEC_INDEX.md` secao 8 (INV-001 a INV-015)
5. **Mobile premium** — breakpoints 320, 375, 414, 428, 768 e desktop
6. **REST + JWT** — contrato em `SPEC_REST_API.md`
7. **pt-BR** — UI, docs de produto, mensagens ao usuario
8. **Ingles** — codigo (nomes, comentarios tecnicos curtos)
9. **Sem emojis** — codigo, docs, commits
10. **Idempotencia** — webhooks e jobs por chave externa
11. **Git** — nao commitar/push sem instrucao explicita do usuario

---

## 5. Workflow padrao

```
ANALISAR -> PLANEJAR -> IMPLEMENTAR -> TESTAR -> VALIDAR -> DOCUMENTAR
```

1. **ANALISAR** — SPEC do dominio + ADRs + `docs/handoff/CURRENT.md` (nao carregar archive inteiro)
2. **PLANEJAR** — endpoints, telas, invariantes afetados
3. **IMPLEMENTAR** — backend antes do frontend (salvo tarefa so UI)
4. **TESTAR** — `go test`, `pnpm test`; mobile nos breakpoints
5. **VALIDAR** — skill `smid-invariants-check` ou checklist secao 12
6. **DOCUMENTAR** — SPEC/ADR se regra mudou; handoff se sessao encerrada

---

## 6. Backend (Go)

### 6.1 Estrutura real

```
backend/
  cmd/server/main.go
  internal/
    config/           # .env, Docker secrets *_FILE
    db/pools.go       # AliasSmid, AliasPermission, AliasLog, AliasCommunication
    domain/<dominio>/ # entities, repository, service, *_impl.go
    http/
      handlers/
      middleware/     # logging (auth JWT: criar/aplicar aqui)
      routes.go
```

JWT hoje: `internal/domain/admin/auth_service.go` (`Login`, `ValidateToken`).

### 6.2 Camadas

- **Handlers finos** — bind, validate, service, JSON; erros sanitizados
- **Services** — regras de negocio; sem import de Echo
- **Repositories** — interface + `*_impl.go`; pool via `db.Pools.Get(alias)`
- **Context** — propagar `context.Context`
- **Logging** — `slog`; nunca senha/token/PII

### 6.3 Convencoes

| Item | Convencao |
|------|-----------|
| Pacote | `lowercase` |
| Tipos exportados | `PascalCase` |
| Arquivos | `snake_case.go` |
| Erros | `ErrXxx` + `fmt.Errorf("%w", err)` |
| Testes | `*_test.go`, `testify/assert` |
| Commits | Conventional Commits em **ingles** |

### 6.4 Comandos uteis

```bash
cd backend && cp .env.example .env
go build ./...
go vet ./...
go test -race ./internal/domain/<dominio>/...
go run ./cmd/server
```

---

## 7. Frontend (Vite + React + shadcn/ui)

### 7.1 Estrutura real

```
frontend/src/
  routes/              # TanStack Router (file-based)
    _authenticated/    # layout autenticado
    (auth)/            # sign-in, etc.
  features/            # modulos por dominio (criar smid aqui)
  components/
    ui/                # shadcn gerados
    layout/            # sidebar, shell
  lib/api/client.ts    # ApiClient + apiClient singleton
  stores/auth-store.ts # token (hoje mock no sign-in)
```

Componentes SMID reutilizaveis: criar `frontend/src/components/smid/` (SidePanel, DataTable, etc.) conforme `SPEC_UX_UI.md`.

### 7.2 Regras

- shadcn/ui primeiro; custom em `components/smid/`
- TanStack Query para API; `apiClient.setAuthToken` apos login
- react-hook-form + zod em formularios
- Mobile-first Tailwind (`base` -> `md:` -> `lg:`)
- Cortina lateral: `Sheet` shadcn em wrapper `SidePanel`
- Tema: provider em `context/theme-provider.tsx`
- Icones: `lucide-react` (sem emojis)

### 7.3 Comandos

```bash
cd frontend && pnpm install
pnpm dev          # Vite, porta padrao do projeto
pnpm test
pnpm build
```

Env: `VITE_API_BASE_URL` (ver `frontend/.env.example`).

---

## 8. Banco de dados

- Listagens: `excluido_em IS NULL`
- Escopo: `unidd_id` quando a tabela tiver o campo (claims JWT futuros)
- Status: `lead_status_categoria` / `ped_status_categoria` — **nunca IDs fixos** (INV-011)
- SQL: parametrizado; preferir sqlc em codigo novo (ADR 0004 — evitar dialectismos MySQL onde possivel)
- Migrations SMID 10: `golang-migrate` apenas para deltas `s10_*`

Aliases em Go: `db.AliasSmid`, `db.AliasPermission`, `db.AliasLog`, `db.AliasCommunication`.

---

## 9. Seguranca

| Item | Regra |
|------|--------|
| Senhas | bcrypt compativel com legado |
| JWT | HS256, 8h; aplicar middleware em `/api/v1/*` exceto `/auth/login` |
| Validacao | 100% no backend |
| Auditoria | escritas relevantes -> `SystemChangeLog` (SPEC_LOG) |
| Secrets | `.env` / Swarm secrets; nunca no git |
| CORS | restrito em producao |

---

## 10. Mapa de dominios (backend)

| Status API | Dominios |
|------------|----------|
| Implementado (CRUD) | admin, log (leitura), communication, tarefas, leads, visitas, historicos, pedidos, produtos, televendas (contatos), representantes (despesas), suporte, compras, financeiro, comissoes, km |
| Pendente | metas, relatorios, integracoes_jobs, televendas (orcam/fila), regras transversais |

Rotas: `backend/internal/http/routes.go`.

---

## 11. Workflow Git

- `main` = integracao; trabalho em `feature/<modulo>-<descricao>` ou `fix/<descricao>`
- Nao commitar sem pedido explicito do usuario
- Nao push/merge/checkout destrutivo sem autorizacao

---

## 12. Checklist antes de concluir

- [ ] SPEC(s) e invariantes lidos
- [ ] `unidd_id` e `excluido_em` nas queries aplicaveis
- [ ] Testes passando (`go test` / `pnpm test`)
- [ ] Mobile nos breakpoints obrigatorios
- [ ] Sem credenciais ou PII em logs
- [ ] SPEC/ADR atualizados se regra mudou
- [ ] Handoff atualizado se encerrar sessao longa (skill `smid-handoff-update`)

---

## 13. Quando atualizar este arquivo

- Nova convencao transversal
- Mudanca de stack ou principio (com ADR)
- Mudanca material no harness (skills, MCP, estrutura de pastas)

---

**Versao**: 0.2.0  
**Status**: documento vivo — alinhado ao harness em `.cursor/`
