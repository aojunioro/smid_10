---
description: Atualiza HANDOFF_BOOTSTRAP.md ao final de uma sessão
---

Use ao final de qualquer sessão produtiva quando o usuário pedir para **documentar a sessão**, **fechar a sessão**, ou **atualizar o handoff**.

Documento alvo: `docs/handoff/HANDOFF_BOOTSTRAP.md`.

## Passos

1. **Coletar mudanças da sessão**:
   ```
   git log --oneline --no-merges <ultimo-sha-no-handoff>..HEAD
   git diff --stat <ultimo-sha-no-handoff>..HEAD
   ```
   O `<ultimo-sha-no-handoff>` é deduzido do estado atual de §6 do handoff.

2. **Reler o handoff inteiro** antes de editar. Atentar especialmente para:
   - §3 (Decisões fechadas)
   - §6 (O que já foi concluído)
   - §7 (Pendências)
   - §8 (Próximo passo executável)
   - §9 (Riscos)
   - Status final

3. **Atualizar §3** se houve nova ADR; cada ADR uma linha numerada com caminho.

4. **Atualizar §6** adicionando um sub-bloco `### Fase X.Y — <titulo> (concluida em <data>)` com bullets factuais. Não apagar blocos antigos.

5. **Atualizar §7** removendo o concluído e adicionando novas pendências. Agrupar por área (Backend, Frontend, Deploy, Docs).

6. **Atualizar §8** com o próximo passo executável objetivo (uma frase, referência a runbook/SPEC).

7. **Atualizar §9** se surgiu risco novo (dependência externa, recurso apertado, decisão reversível).

8. **Atualizar o status final** refletindo o estado atual em uma frase.

9. **Não criar** `progress.txt` paralelos se o handoff já cobre. Centralizar no handoff.

10. **Commit dedicado** (não misturar com feature):
    ```
    docs(handoff): atualiza estado após <tema da sessão>
    ```

## Princípios

- Handoff é o **único** documento de continuidade entre threads. Se não está lá, não existe para a próxima sessão.
- Linguagem em pt-BR, factual, sem emoji (AGENTS.md §2.9).
- Cada bullet em §6 deve ser auditável: aponta para arquivo, SHA de commit ou ADR.
- Pendências sem responsável + data de revisão devem ser deslocadas para um "Backlog" se ainda não há plano de execução.
