# Handoff SMID 10 — serie de continuidade

Objetivo: **nao inflar a janela de contexto** do agente. Historico fica arquivado; cada sessao le so o necessario.

---

## O que ler (por situacao)

| Situacao | Arquivo | Tamanho alvo |
|----------|---------|----------------|
| Retomar trabalho / nova thread | **`CURRENT.md`** | ~80-120 linhas |
| Regras de operacao | `AGENTS.md` | secao 3 + harness |
| Detalhe de fase antiga | `archive/*.md` | sob demanda |
| Entrada de sessao pontual | `journal/` | 1 arquivo curto por sessao |

**Nao** carregar `archive/BOOTSTRAP_HISTORY.md` inteiro salvo pedido explicito de auditoria ou historico.

---

## Arquivos da serie

```
docs/handoff/
  README.md                 <- este indice
  CURRENT.md                <- UNICO obrigatorio entre sessoes
  archive/
    BOOTSTRAP_HISTORY.md    <- monolito original (fases 0.x concluidas)
  journal/
    YYYY-MM-DD-<tema>.md    <- opcional, 1 pagina por sessao produtiva
```

---

## Regras de manutencao

1. **`CURRENT.md` e o unico documento vivo** para pendencias, proximo passo e riscos.
2. Ao concluir uma fase grande, **resumir** em 3-5 bullets em `CURRENT.md` e mover detalhe para `archive/` ou `journal/`.
3. **Nao** acrescentar dezenas de `### Fase X.Y` em `CURRENT.md` — usar tabela-resumo + link para archive.
4. `archive/` e **append-only** (congelado); corrigir erro factual via nota em `journal/`, nao reescrever historico.
5. Atualizar `ultimo_commit` em `CURRENT.md` ao fechar sessao (SHA curto).

Skill do agente: `.cursor/skills/smid-handoff-update/`.

---

## Migracao (2026-05-17)

O arquivo `HANDOFF_BOOTSTRAP.md` virou ponte para esta serie. Conteudo historico em `archive/BOOTSTRAP_HISTORY.md`.
