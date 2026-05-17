---
name: smid-handoff-update
description: Updates docs/handoff/CURRENT.md after a work session; optional journal entry. Use when closing a session, documenting progress, updating handoff, or preparing context for the next agent thread.
---

# SMID — Atualizar handoff (serie enxuta)

## O que editar

| Arquivo | Quando |
|---------|--------|
| **`docs/handoff/CURRENT.md`** | Sempre — unico documento vivo |
| `docs/handoff/journal/YYYY-MM-DD-<tema>.md` | Sessao longa ou decisao importante (~30 linhas max) |
| `docs/handoff/archive/` | **Nao editar** (historico congelado) |

Indice da serie: `docs/handoff/README.md`.

## Coletar evidencias

```bash
git log --oneline --no-merges <ultimo-sha-em-CURRENT>..HEAD
git diff --stat <ultimo-sha-em-CURRENT>..HEAD
```

## Atualizar CURRENT.md

| Secao | Acao |
|-------|------|
| Cabecalho | `Atualizado em`, `ultimo_commit` (SHA curto) |
| Resumo executivo | 1 paragrafo se mudou maturidade global |
| Proximo passo | Uma acao objetiva + referencia SPEC/skill |
| Pendencias | Checkboxes; remover concluido |
| Concluido (indice) | Nova linha na **tabela** (nao duplicar blocos de fase) |
| Riscos | So riscos ainda ativos |

**Nao** colar listas enormes de fases em CURRENT — detalhe vai para `journal/` ou permanece no archive.

## Journal (opcional)

Criar `docs/handoff/journal/YYYY-MM-DD-<tema>.md`:

- Entregue (bullets)
- Decisoes
- Proxima sessao (1 linha, ou "ver CURRENT")

## Principios

- pt-BR, factual, sem emoji.
- Meta: **CURRENT.md permanece abaixo de ~120 linhas**.
- Commit dedicado so se usuario pedir: `docs(handoff): update CURRENT after <topic>`

Alinhar `AGENTS.md` secao 3 se mudou maturidade global.
