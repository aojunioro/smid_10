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
2. ~~Criar registros DNS A~~:
   - `api.s10.smydi.com.br` → `216.144.235.25` (concluído, DNS only na Cloudflare)
   - `s10.smydi.com.br` → `216.144.235.25` (pendente - frontend ainda não deployado)
3. ~~Push do código para `origin/main` para acionar o build no CI~~
4. ~~Na VPS: criar os 8 Swarm secrets (`deploy/README.md §3`)~~
5. ~~Na VPS: popular o volume `smid10_mariadb_init` com `01-schemas.sql` (`deploy/README.md §4`)~~
6. ~~Na VPS: `docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth`~~
7. ~~Validar `https://api.s10.smydi.com.br/healthz` retornando `status: "ok"`~~ (smoke-test interno e externo OK)

**Nota**: O volume `smid10_mariadb_data` já existia antes do deploy, então os scripts de init não foram executados automaticamente. Os bancos foram criados manualmente via SQL direto no container. O certificado Let's Encrypt foi gerado após mudar o DNS de api.s10.smydi.com.br para "DNS only" na Cloudflare (HTTP challenge não funciona através do proxy).

### Backend (Fase 0.2+)

1. ~~Middleware de logging estruturado por requisição (request_id, latência, status)~~ — concluído
2. ~~Repositórios atrás de interfaces (preparação para ADR 0004 — Postgres pós-cutover)~~ — concluído
3. ~~Setup de testes com `testcontainers` para repositórios~~ — concluído

### Backend (Fase 0.3 — Plataforma Admin) — concluído em 2026-05-13

1. ~~Ler SPEC_ADMIN e planejar implementação~~ — concluído
2. ~~Criar entidades SystemUser, SystemGroup, SystemRole~~ — concluído
3. ~~Criar repositórios para Admin (banco permission)~~ — concluído
4. ~~Criar script de migração para tabelas admin~~ — concluído
5. ~~Redefinir senhas do MariaDB e aplicar migração~~ — concluído
6. ~~Atualizar secrets do Swarm com novas senhas e deploy~~ — concluído
7. ~~Implementar serviço de autenticação (login, JWT)~~ — concluído
8. ~~Criar handler REST para login e rotas~~ — concluído
9. ~~Corrigir erros de compilação e build do backend~~ — concluído
10. ~~Deploy da nova imagem na VPS~~ — concluído
11. ~~Corrigir configuração CORS e testar login~~ — concluído
12. ~~Adicionar colunas faltantes na tabela system_users~~ — concluído
13. ~~Corrigir SystemUser para usar ponteiros em campos nullable~~ — concluído
14. ~~Testar endpoint de login com sucesso~~ — concluído (POST /api/v1/auth/login retornando JWT)
15. ~~Implementar CRUD de usuários (Create, Read, Update, Delete)~~ — concluído (POST /api/v1/users, GET /api/v1/users, GET /api/v1/users/:id, PUT /api/v1/users/:id, DELETE /api/v1/users/:id)
16. ~~Implementar CRUD de grupos (SystemGroup)~~ — concluído (POST /api/v1/groups, GET /api/v1/groups, GET /api/v1/groups/:id, PUT /api/v1/groups/:id, DELETE /api/v1/groups/:id)
17. ~~Implementar CRUD de papéis (SystemRole)~~ — concluído (POST /api/v1/roles, GET /api/v1/roles, GET /api/v1/roles/:id, PUT /api/v1/roles/:id, DELETE /api/v1/roles/:id)
18. ~~Implementar CRUD de programas (SystemProgram)~~ — concluído (POST /api/v1/programs, GET /api/v1/programs, GET /api/v1/programs/:id, PUT /api/v1/programs/:id, DELETE /api/v1/programs/:id)
19. ~~Implementar CRUD de unidades (SystemUnit)~~ — concluído (POST /api/v1/units, GET /api/v1/units, GET /api/v1/units/:id, PUT /api/v1/units/:id, DELETE /api/v1/units/:id)
20. ~~Criar tabelas faltantes no banco de dados~~ — concluído (system_groups, system_roles, system_programs, system_units)
21. ~~Reverter CI para Docker Hub devido a problema com GHCR~~ — concluído

**Nota**: O login está funcionando corretamente em https://api.s10.smydi.com.br/api/v1/auth/login. As credenciais de teste são login: admin, senha: Admin123!. Todos os endpoints CRUD de usuários, grupos, papéis, programas e unidades estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.4 — Plataformas Transversais) — concluído em 2026-05-14

1. ~~Ler SPEC_LOG e planejar implementação~~ — concluído
2. ~~Criar banco de dados log e tabelas (system_access_log, system_change_log, system_sql_log, system_request_log, system_access_notification_log, system_schedule_log, system_sql_changes)~~ — concluído
3. ~~Implementar leitura de SystemAccessLog~~ — concluído (GET /api/v1/logs/access)
4. ~~Implementar leitura de SystemChangeLog~~ — concluído (GET /api/v1/logs/change)
5. ~~Implementar leitura de SystemSqlLog~~ — concluído (GET /api/v1/logs/sql)
6. ~~Deploy e testar endpoints de logs~~ — concluído

7. ~~Ler SPEC_COMMUNICATION e planejar implementação~~ — concluído
8. ~~Criar banco de dados communication e tabelas (system_notification, system_message)~~ — concluído
9. ~~Implementar CRUD de SystemNotification~~ — concluído (POST /api/v1/notifications, GET /api/v1/notifications, GET /api/v1/notifications/:id, PUT /api/v1/notifications/:id)
10. ~~Implementar CRUD de SystemMessage~~ — concluído (POST /api/v1/messages, GET /api/v1/messages, GET /api/v1/messages/:id, PUT /api/v1/messages/:id)
11. ~~Deploy e testar endpoints de comunicação~~ — concluído

12. ~~Ler SPEC_TAREFAS e planejar implementação~~ — concluído
13. ~~Criar banco de dados smid e tabela tarefas~~ — concluído
14. ~~Implementar CRUD de Tarefa~~ — concluído (POST /api/v1/tarefas, GET /api/v1/tarefas, GET /api/v1/tarefas/:id, PUT /api/v1/tarefas/:id)
15. ~~Deploy e testar endpoints de tarefas~~ — concluído

**Nota**: Todos os endpoints de logs, comunicação e tarefas estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.5 — Plataformas de Negócio) — concluído em 2026-05-14

1. ~~Ler SPEC_LEADS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (leads, lead_status, lead_meio, midias, endereco)~~ — concluído
4. ~~Atualizar entidades Lead para corresponder à tabela do banco~~ — concluído
5. ~~Implementar repositório LeadRepository~~ — concluído
6. ~~Implementar serviço LeadService~~ — concluído
7. ~~Implementar handler LeadHandler~~ — concluído
8. ~~Adicionar rotas de leads no routes.go~~ — concluído
9. ~~Corrigir placeholders SQL (PostgreSQL $1 para MySQL ?)~~ — concluído
10. ~~Corrigir erro de INSERT simplificando para colunas obrigatórias~~ — concluído
11. ~~Corrigir erro de Scan usando COALESCE para campos de data NULL~~ — concluído
12. ~~Deploy e testar endpoints de leads~~ — concluído (POST /api/v1/leads, GET /api/v1/leads, GET /api/v1/leads/:id, PUT /api/v1/leads/:id)
13. ~~Configurar nginx proxy para API no frontend~~ — concluído

**Nota**: Todos os endpoints de leads estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/. O nginx do frontend foi configurado para fazer proxy das requisições /api para o backend.

### Backend (Fase 0.6 — Plataformas de Negócio) — concluído em 2026-05-14

1. ~~Ler SPEC_VISITAS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (visitas, visita_checkin_event, pos_visita)~~ — concluído
4. ~~Atualizar entidades Visita, VisitaCheckinEvent, PosVisita~~ — concluído
5. ~~Implementar repositório VisitaRepository~~ — concluído
6. ~~Implementar serviço VisitaService~~ — concluído
7. ~~Implementar handler VisitaHandler~~ — concluído
8. ~~Adicionar rotas de visitas no routes.go~~ — concluído
9. ~~Corrigir erros de SQL usando sql.NullTime~~ — concluído
10. ~~Deploy e testar endpoints de visitas~~ — concluído (POST /api/v1/visitas, GET /api/v1/visitas, GET /api/v1/visitas/:id, PUT /api/v1/visitas/:id, DELETE /api/v1/visitas/:id)

**Nota**: Todos os endpoints de visitas estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.6 — Plataformas de Negócio) — concluído em 2026-05-14

1. ~~Ler SPEC_HISTORICOS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (historicos, hist_motivo, hist_ocorrido)~~ — concluído
4. ~~Atualizar entidades HistoricoRepre, HistoricoMotivo, HistoricoOcorrido~~ — concluído
5. ~~Implementar repositório HistoricoRepository~~ — concluído
6. ~~Implementar serviço HistoricoService~~ — concluído
7. ~~Implementar handler HistoricoHandler~~ — concluído
8. ~~Adicionar rotas de históricos no routes.go~~ — concluído
9. ~~Deploy e testar endpoints de históricos~~ — concluído (POST /api/v1/historicos, GET /api/v1/historicos, GET /api/v1/historicos/:id, PUT /api/v1/historicos/:id, DELETE /api/v1/historicos/:id)

**Nota**: Todos os endpoints de históricos estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.6 — Plataformas de Negócio) — concluído em 2026-05-15

1. ~~Ler SPEC_PEDIDOS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (pedidos, ped_status_categoria, ped_status, ped_fpgto, ped_cpgto, ped_canal, ped_prod_item)~~ — concluído
4. ~~Atualizar entidades Pedido, PedProdItem, PedStatus, PedStatusCategoria, PedFormaPagamento, PedCondicaoPagamento, PedCanal~~ — concluído
5. ~~Implementar repositório PedidoRepository~~ — concluído
6. ~~Implementar serviço PedidoService~~ — concluído
7. ~~Implementar handler PedidoHandler~~ — concluído
8. ~~Adicionar rotas de pedidos no routes.go~~ — concluído
9. ~~Deploy e testar endpoints de pedidos~~ — concluído (POST /api/v1/pedidos, GET /api/v1/pedidos, GET /api/v1/pedidos/:id, PUT /api/v1/pedidos/:id, DELETE /api/v1/pedidos/:id)

**Nota**: Todos os endpoints de pedidos estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.6 — Plataformas de Negócio) — concluído em 2026-05-15

1. ~~Ler SPEC_PRODUTOS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (produtos, prod_categ, prod_medidas, prod_modelos)~~ — concluído
4. ~~Atualizar entidades Produto, ProdCateg, ProdMedidas, ProdModelos~~ — concluído
5. ~~Implementar repositório ProdutoRepository~~ — concluído
6. ~~Implementar serviço ProdutoService~~ — concluído
7. ~~Implementar handler ProdutoHandler~~ — concluído
8. ~~Adicionar rotas de produtos no routes.go~~ — concluído
9. ~~Deploy e testar endpoints de produtos~~ — concluído (POST /api/v1/produtos, GET /api/v1/produtos, GET /api/v1/produtos/:id, PUT /api/v1/produtos/:id, DELETE /api/v1/produtos/:id)

**Nota**: Todos os endpoints de produtos estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.7 — Televendas) — concluído em 2026-05-15

1. ~~Ler SPEC_TELEVENDAS e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (televendas_contatos, televendas_status, tele_coeficiente, orcam, orcam_status, orcam_prod_item)~~ — concluído
4. ~~Atualizar entidades TelevendasContato, TelevendasStatus, TeleCoeficiente, Orcam, OrcamStatus, OrcamProdItem~~ — concluído
5. ~~Implementar repositório TelevendasContatoRepository~~ — concluído
6. ~~Implementar serviço TelevendasContatoService~~ — concluído
7. ~~Implementar handler TelevendasHandler~~ — concluído
8. ~~Adicionar rotas de televendas no routes.go~~ — concluído
9. ~~Deploy e testar endpoints de televendas~~ — concluído (POST /api/v1/televendas/contatos, GET /api/v1/televendas/contatos, GET /api/v1/televendas/contatos/:id, PUT /api/v1/televendas/contatos/:id, DELETE /api/v1/televendas/contatos/:id)

**Nota**: Todos os endpoints de televendas estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Backend (Fase 0.8 — Representantes) — concluído em 2026-05-15

1. ~~Ler SPEC_REPRESENTANTES e planejar implementação~~ — concluído
2. ~~Verificar tabelas existentes no banco smid~~ — concluído
3. ~~Criar tabelas faltantes (repre_despesa_categ, repre_despesa_extra)~~ — concluído
4. ~~Atualizar entidades RepreDespesaExtra, RepreDespesaCateg~~ — concluído
5. ~~Implementar repositório RepreDespesaExtraRepository~~ — concluído
6. ~~Implementar serviço RepreDespesaExtraService~~ — concluído
7. ~~Implementar handler RepresentanteHandler~~ — concluído
8. ~~Adicionar rotas de representantes no routes.go~~ — concluído
9. ~~Deploy e testar endpoints de representantes~~ — concluído (POST /api/v1/representantes/despesas-extras, GET /api/v1/representantes/despesas-extras, GET /api/v1/representantes/despesas-extras/:id, PUT /api/v1/representantes/despesas-extras/:id, DELETE /api/v1/representantes/despesas-extras/:id)

**Nota**: Todos os endpoints de representantes estão funcionando corretamente em https://api.s10.smydi.com.br/api/v1/.

### Frontend (Fase 0.1 do frontend)

1. ~~Scaffold do `shadcn-admin` em `frontend/` (opção A do `frontend/README.md`)~~ — concluído
2. ~~Configurar `lib/api/` apontando para `https://api.s10.smydi.com.br`~~ — concluído
3. ~~Tema dark/light já configurado no shadcn-admin~~ — concluído
4. ~~Criar Dockerfile e configurar CI para build da imagem~~ — concluído
5. ~~Descomentar serviço web em swarm-stack.yml~~ — concluído
6. ~~Mudar CI de GHCR para Docker Hub (GHCR tinha problemas de visibilidade)~~ — concluído
7. ~~Deploy do frontend na VPS~~ — concluído
8. ~~Configurar DNS para `s10.smydi.com.br` → `216.144.235.25` (DNS only na Cloudflare)~~ — concluído
9. ~~Smoke-test do frontend: curl https://s10.smydi.com.br~~ — concluído

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

**Status**: Fase 0.1 do backend concluída; Fase 0.2 do backend concluída; Fase 0.3 (Plataforma Admin) concluída com login funcionando; ambiente de deploy operacional em VPS; frontend Fase 0.1 concluída
