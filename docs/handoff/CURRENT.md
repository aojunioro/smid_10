# SMID 10 — Estado atual

> Handoff operacional. Ler **somente este arquivo** para retomar. Indice da serie: `docs/handoff/README.md`.

| Campo | Valor |
|-------|--------|
| Atualizado em | 2026-05-17 |
| Ultimo commit | `40c4674` (docs handoff SPEC_LOG) |
| Fase README | 0 concluida; **Fase 1 (plataforma) em andamento** |

---

## Leitura minima (ordem)

1. Este arquivo
2. `AGENTS.md` (estado do codigo + harness `.cursor/`)
3. `docs/specs/SPEC_<DOMINIO>.md` da tarefa

---

## Resumo executivo

- **Infra**: API `api.s10.smydi.com.br`, app `s10.smydi.com.br`, stack Swarm `smid10`, healthz OK nos 4 bancos.
- **Backend**: CRUD REST para a maioria dos dominios; regras de negocio do SPEC ainda rasas; **JWT nao aplicado** nas rotas protegidas.
- **Frontend**: Vite + shadcn-admin deployado; **login mock**; sem telas SMID (leads, visitas, etc.).
- **Harness**: AGENTS v0.2, skills/rules/MCP em `.cursor/` (2026-05-17).
- **Sem API ainda**: metas, relatorios, integracoes_jobs.

---

## Proximo passo executavel

**Plataforma + primeiro vertical slice**

1. Middleware JWT em `/api/v1/*` (exceto `/auth/login`, `/healthz`).
2. Login real no frontend (`POST /api/v1/auth/login`) + `apiClient.setAuthToken`.
3. Tela **Leads** (lista + cortina) consumindo API, com regras basicas do `SPEC_LEADS` (status, soft delete, `unidd_id`).

Referencia: `AGENTS.md` secao 3; skill `smid-frontend-feature`.

---

## Pendencias abertas

### Plataforma (prioridade alta)

- [ ] Middleware JWT em `backend/internal/http/middleware/`
- [ ] Refresh token (SPEC_REST_API) — apos JWT base
- [ ] Escopo `unidd_id` via claims em listagens
- [ ] Auditoria de escrita (`SystemChangeLog`) nas mutacoes relevantes

### Backend (negocio)

- [ ] Regras Lead ↔ Visita (INV-002, INV-003)
- [ ] Duplicidade de telefone (INV-012)
- [ ] Dominios: `SPEC_METAS`, `SPEC_RELATORIOS`, `SPEC_INTEGRACOES_JOBS`
- [ ] Televendas: orcamentos/fila (alem de contatos)
- [ ] Cobertura de testes alem de `leads/repository_test.go`

### Frontend

- [ ] Substituir mock em `user-auth-form.tsx`
- [ ] Menu/sidebar SMID (remover demos do template)
- [ ] Componentes `components/smid/` (SidePanel, etc. per SPEC_UX_UI)

### Docs / harness

- [ ] ADR ou nota se mantiver Vite vs Next (drift com ADR 0003)
- [ ] Entradas futuras em `journal/` ao fechar sessoes longas

---

## Concluido (indice — detalhe no archive)

| Area | Ate | Referencia |
|------|-----|------------|
| Bootstrap repo, SPECs, ADRs | 2026-05-12 | `archive/BOOTSTRAP_HISTORY.md` |
| Backend healthz, pools, deploy VPS | 2026-05-13 | idem |
| Admin + auth login JWT emitido | 2026-05-13 | idem |
| Log, communication, tarefas | 2026-05-14 | idem |
| Leads, visitas, historicos, pedidos, produtos | 2026-05-15 | idem |
| Televendas contatos, representantes despesas, suporte, compras | 2026-05-15 | idem |
| Financeiro, comissoes, KM, request log | 2026-05-15 | idem |
| Frontend scaffold + deploy | 2026-05-14/15 | idem |
| Harness Cursor (skills, rules, MCP) | 2026-05-17 | `journal/2026-05-17-harness.md` |

---

## Riscos ativos

1. Schema legado intocavel — DDL em tabela existente exige ADR.
2. Rotas API sem JWT — tratar como bloqueador de staging publico.
3. Sessoes PHP 8.x e JWT 10 independentes (coexistencia).
4. Queries: `excluido_em IS NULL` e `unidd_id` sempre que aplicavel.
5. VPS com RAM limitada — nao remover limits do `swarm-stack.yml`.
6. Credenciais: nunca commitar; rotacionar SSH se exposto em sessao antiga.

---

## Prompt de retomada (copiar)

```md
Continuar SMID 10.

Ler: docs/handoff/CURRENT.md, AGENTS.md, SPEC do dominio.
Nao carregar archive salvo necessidade.

Proximo passo: JWT + login real + tela Leads.
Respeitar invariantes SPEC_INDEX secao 8. Nao alterar schema legado.
```

---

## Comandos uteis

```bash
cd backend && go run ./cmd/server
cd frontend && pnpm dev
curl -s https://api.s10.smydi.com.br/healthz
```

Histórico completo das fases 0.x: `docs/handoff/archive/BOOTSTRAP_HISTORY.md`.
