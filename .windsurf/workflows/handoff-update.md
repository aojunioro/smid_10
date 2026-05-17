---
description: Atualiza a serie de handoff (CURRENT.md) ao final de uma sessao
---

Use ao fechar sessao produtiva ou quando o usuario pedir atualizar handoff.

Documentacao da serie: `docs/handoff/README.md`.
Alvo vivo: **`docs/handoff/CURRENT.md`** (nao inflar; meta ~120 linhas).

## Passos

1. Coletar `git log` e `git diff` desde o `ultimo_commit` em CURRENT.
2. Atualizar cabecalho (data, SHA).
3. Ajustar proximo passo, pendencias (checkboxes), tabela "Concluido (indice)".
4. Se sessao longa: criar `docs/handoff/journal/YYYY-MM-DD-<tema>.md` (~30 linhas).
5. **Nao** editar `archive/`; **nao** reexpandir CURRENT com blocos `### Fase X.Y`.

Skill Cursor: `.cursor/skills/smid-handoff-update/`.

Commit dedicado so se usuario pedir: `docs(handoff): update CURRENT after <tema>`.
