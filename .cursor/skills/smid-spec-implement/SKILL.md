---
name: smid-spec-implement
description: Implements a SMID 10 backend domain from docs/specs (leads, visitas, pedidos, etc.) with Go Echo handlers, domain services, MySQL repositories, and tests. Use when implementing SPEC_*, REST endpoints, internal/domain packages, or extending routes.go.
---

# SMID — Implementar dominio (backend)

## Pre-condicao

Nome do dominio conhecido (ex.: `LEADS` -> `docs/specs/SPEC_LEADS.md`). Se ambiguo, perguntar.

## Passos

1. Ler o SPEC integralmente: entidades, fluxos secao 3, maquina de estados secao 4, REST.
2. Ler `docs/specs/SPEC_INDEX.md` secao 8 (INV-001 a INV-015) e `AGENTS.md` ss 6 e 8.
3. Verificar o que ja existe em `backend/internal/domain/<dominio>/` e `routes.go` — estender, nao duplicar CRUD cego.

## Estrutura padrao

```
backend/internal/domain/<dominio>/
  entities.go
  repository.go          # interface
  repository_impl.go     # MySQL via pools.Get(db.Alias*)
  <dominio>_service.go
  errors.go
  *_test.go
backend/internal/http/handlers/<dominio>_handler.go
```

Registrar rotas em `backend/internal/http/routes.go`.

## Regras SQL e dados

- `excluido_em IS NULL` em listagens.
- `unidd_id` quando a tabela tiver o campo.
- Categorias de status, nunca IDs magicos.
- Sem DDL em tabelas legadas (ADR 0002).
- SQL portavel onde possivel (ADR 0004): `COALESCE`, `LIMIT n OFFSET m`.

## Handlers

- Validar body (validator ou manual).
- `errors.Is` para `ErrNotFound` -> 404, validacao -> 400.
- Respostas JSON alinhadas a `SPEC_REST_API.md`.

## Testes

```bash
cd backend && go test -race ./internal/domain/<dominio>/...
go vet ./...
```

Cobrir regras do SPEC no service (mocks do repository).

## Checklist final

- [ ] SPEC lido; invariantes aplicados
- [ ] Interface de repository (ADR 0004)
- [ ] Rotas registradas
- [ ] Testes verdes
- [ ] Sem segredos no codigo
- [ ] Atualizar handoff se sessao encerrar (`smid-handoff-update`)

Nao commitar sem pedido do usuario.
