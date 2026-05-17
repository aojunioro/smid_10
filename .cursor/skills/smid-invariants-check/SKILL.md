---
name: smid-invariants-check
description: Reviews SMID 10 changes against SPEC_INDEX invariants INV-001 to INV-015 and AGENTS.md data rules before merge or session end. Use for PR review, pre-commit validation, or "check invariants" requests.
---

# SMID — Checagem de invariantes

Fonte: `docs/specs/SPEC_INDEX.md` secao 8.

## Checklist rapido

| ID | Verificar |
|----|-----------|
| INV-001 | Lead sempre com status valido |
| INV-002 | Status lead espelha visita valida |
| INV-003 | Exclusao visita restaura ultimo status tipo L |
| INV-004 | GPS valido antes de KM |
| INV-005 | Canal produto exclusivo S/N |
| INV-006 | Pedido televendas exige lead_id |
| INV-007 | Comissao cancelada fora do saldo |
| INV-008 | Suporte exige ped_id |
| INV-009 | Lancamento financeiro automatico idempotente |
| INV-010 | Historico com vis_id e lead_id |
| INV-011 | Relatorios por categoria, nao ID fixo |
| INV-012 | Duplicidade telefone (20min / 6h) |
| INV-013 | Notificacao tarefa so para responsavel |
| INV-014 | Permissoes no backend, nao so menu |
| INV-015 | Escritas auditaveis em change log |

## Queries (sempre que aplicavel)

- `excluido_em IS NULL` em listagens
- Filtro `unidd_id` coerente com JWT/sessao
- Transacoes nao cruzam aliases de banco

## Saida

Listar: **OK** | **NAO VERIFICADO** | **VIOLACAO** com arquivo/linha e sugestao de correcao.

Se violacao em codigo novo, corrigir antes de considerar a tarefa concluida.
