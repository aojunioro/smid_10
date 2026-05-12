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
4. **Template inicial recomendado**: `shadcn-admin` (https://github.com/satnaing/shadcn-admin) para o frontend

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

- Estrutura de pastas (`backend/`, `frontend/`, `docs/`)
- 24 SPECs canônicos copiados para `docs/specs/`
- Schema base do legado copiado para `docs/legacy-schema/`
- Docs críticos de UX/UI copiados para `docs/legacy-reference/`
- 3 ADRs iniciais escritos
- `README.md`, `AGENTS.md` e `.gitignore` criados
- `backend/.env.example` com DSN dos 4 bancos
- Repositório git inicializado, commit inicial e push para `origin/main`

---

## 7. Pendências Abertas (Fase 0)

### Backend

1. Inicializar módulo: `cd backend && go mod init github.com/aojunioro/smid_10/backend`
2. Adicionar dependências (ver `backend/README.md` seção 3.2)
3. Criar `cmd/server/main.go` com Echo + `GET /healthz`
4. Criar `internal/config/config.go` lendo `.env` via `godotenv`
5. Criar `internal/db/pools.go` com 4 pools (`smid`, `permission`, `log`, `communication`)
6. Health check que valida ping nos 4 pools

### Frontend

1. Scaffold do `shadcn-admin` em `frontend/` (opção A do `frontend/README.md`)
2. Configurar `lib/api/client.ts` apontando para `http://localhost:8080`
3. Setar tema dark/light e branding

---

## 8. Próximo Passo Executável

> **Fase 0.1 do backend**: criar `cmd/server/main.go` minimal com Echo expondo `GET /healthz` retornando `200 OK` com `{"status":"ok"}`, lendo a porta de `.env`. Em seguida, implementar `internal/config/config.go` e `internal/db/pools.go` para conectar aos 4 bancos legados e estender o healthz para validar pings.

---

## 9. Riscos e Cuidados

1. **Não alterar schema legado** — qualquer DDL que toque tabelas existentes precisa de novo ADR
2. **Compatibilidade de hash de senha** — o legado usa o que o Adianti gerar; validar exatamente o mesmo algoritmo no Go (provavelmente bcrypt; confirmar no primeiro login)
3. **Encoding/collation** — algumas tabelas legadas têm mojibake (`latin1` com UTF-8 dentro); usar `charset=utf8mb4` no DSN e revisar caso a caso
4. **Sessões PHP do SMID 8.x são independentes** do JWT do SMID 10 — usuário pode estar logado em um e não no outro durante a coexistência
5. **Filtro de unidade**: toda query do SMID 10 deve respeitar `unidd_id` e `excluido_em IS NULL`

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

**Status**: bootstrap concluído, aguardando início da Fase 0.1
