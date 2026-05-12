---
description: Implementa um domínio (SPEC) do SMID 10 seguindo a estrutura canônica
---

Use esta workflow quando o usuário pedir para implementar um SPEC de domínio (leads, visitas, pedidos, etc.) no backend Go.

Pré-condição: o nome do SPEC deve ser conhecido (ex.: SPEC_LEADS). Se não, perguntar.

## Passos

1. Ler o SPEC do domínio em `docs/specs/SPEC_<DOMINIO>.md` integralmente e identificar:
   - Entidades e campos canônicos.
   - Regras de negócio (transições de status, validações, filtros).
   - Endpoints REST exigidos (verbo, path, request, response).
   - Invariantes específicas do domínio.

2. Reler obrigatoriamente `docs/specs/SPEC_INDEX.md` seção 8 (invariantes globais) e `AGENTS.md` seções 2, 4 e 6.

3. Validar invariantes universais que toda implementação respeita:
   - Toda query filtra por `unidd_id` quando a tabela tem esse campo.
   - Toda listagem filtra por `excluido_em IS NULL` (soft delete).
   - Status funcionais via `*_status_categoria`, nunca IDs hard-coded.
   - Nenhum DDL toca tabela legada (ADR 0002).

4. Criar a estrutura padrão em `backend/internal/domain/<dominio>/`:
   - `entity.go` — structs, sem dependência de HTTP nem SQL.
   - `repository.go` — interface dos métodos de persistência (ADR 0004).
   - `repository_mysql.go` — implementação usando pool via `db.Pools.Get(alias)`.
   - `service.go` — regras de negócio puras; recebe `Repository` por DI.
   - `errors.go` — sentinel errors (`ErrNotFound`, `ErrInvalidStatus`, ...).
   - `*_test.go` — unitários do serviço com mock do `Repository`.

5. Criar handlers REST em `backend/internal/http/handlers/<dominio>.go`:
   - Handlers finos: parse + validate + call service + serialize.
   - Validação via `validator/v10`.
   - Erros do domínio mapeados para HTTP com `errors.Is/As`; nunca expor `err.Error()` cru.
   - Registrar rotas em `backend/internal/http/routes.go`.

6. Regras de SQL (preparação Postgres — ADR 0004):
   - Evitar dialectismos MySQL: `GROUP_CONCAT`, `ON DUPLICATE KEY`, `LIMIT x, y`, `IFNULL`, backticks.
   - Preferir: `COALESCE`, `LIMIT n OFFSET m`, INSERT…ON CONFLICT-equivalente isolado.
   - Queries não-portáveis ficam em `internal/db/queries/<dominio>/*.sql` (sqlc).

7. Testes:
   - Unitários do serviço: cobertura ≥ 80% das regras do SPEC.
   - Integração (se aplicável): `testcontainers` com `mariadb:10.11`.
   // turbo
   - Rodar `go test -race ./internal/domain/<dominio>/...` e `go vet ./...`.

8. Validar contra o contrato em `docs/specs/SPEC_REST_API.md`: verbo, path, payload, status code.

9. Atualizar `docs/specs/SPEC_INDEX.md` marcando o checklist do domínio como concluído.

10. Commit Conventional Commits (inglês, sem emoji):
    - `feat(<dominio>): implement <feature> per SPEC_<DOMINIO>`
    - Mensagem cita o SPEC e os endpoints entregues.

## Checklist final

- [ ] SPEC lido integralmente
- [ ] Invariantes universais validados (`unidd_id`, `excluido_em`, categorias)
- [ ] Estrutura `internal/domain/<dominio>/` criada no padrão
- [ ] Repository atrás de interface (ADR 0004)
- [ ] SQL portável onde possível
- [ ] Handlers finos + validação + erros tipados
- [ ] Rotas registradas em `routes.go`
- [ ] Testes unitários verdes (≥ 80%)
- [ ] Contrato bate com `SPEC_REST_API.md`
- [ ] `SPEC_INDEX.md` atualizado
- [ ] Commit em Conventional Commits
