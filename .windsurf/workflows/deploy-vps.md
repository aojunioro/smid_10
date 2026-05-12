---
description: Deploy manual do stack smid10 na VPS Docker Swarm com smoke-test
---

Use esta workflow quando o usuário pedir para **fazer deploy** do SMID 10 na VPS, redeployar a API após nova imagem, ou rodar um smoke-test no ambiente de staging.

Documento de referência canônico: `deploy/README.md`.

## Pré-checagem (sempre)

1. **Confirmar que o usuário não tem outras stacks pausando**. Listar com:
   ```
   sshpass -e ssh root@216.144.235.25 'docker stack ls'
   ```
2. Confirmar que a tag a deployar existe no GHCR (default: `dev`):
   ```
   sshpass -e ssh root@216.144.235.25 'docker manifest inspect ghcr.io/aojunioro/smid10-api:${TAG:-dev} >/dev/null && echo OK || echo MISSING'
   ```

## Primeiro deploy (stack ainda não existe)

Se `docker stack ls | grep smid10` retorna vazio, executar **em ordem**:

1. **DNS**: verificar que `api.s10.smydi.com.br` e `s10.smydi.com.br` resolvem para `216.144.235.25` (`dig +short`).
2. **Secrets**: seguir `deploy/README.md §3.1` literalmente. Pedir ao usuário as senhas (ou gerar e mostrar uma vez) antes de criar. **Nunca** colocar segredos em variáveis de ambiente do compose.
3. **Volume de init MariaDB**: seguir `deploy/README.md §4` para popular `smid10_mariadb_init` com `01-schemas.sql`.
4. **Pull com auth**: se a imagem GHCR estiver privada:
   ```
   sshpass -e ssh root@216.144.235.25 'docker login ghcr.io'
   ```
5. **Deploy**:
   ```
   sshpass -e ssh root@216.144.235.25 'cd /srv/smid10 && docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth'
   ```

## Redeploy (stack já existe, só nova imagem)

```
sshpass -e ssh root@216.144.235.25 'docker service update --image ghcr.io/aojunioro/smid10-api:${TAG:-dev} smid10_api --with-registry-auth'
```

Acompanhar:
```
sshpass -e ssh root@216.144.235.25 'docker service ps smid10_api --no-trunc | head -10'
```

## Smoke-test obrigatório após cada deploy

1. **Interno** (pela rede overlay, isola Traefik):
   ```
   sshpass -e ssh root@216.144.235.25 'docker run --rm --network SmydiNet curlimages/curl -sS http://smid10_api:8080/healthz'
   ```
   Esperado: HTTP 200, `status: "ok"`, **todos** os 4 bancos respondendo.

2. **Externo** (pelo Traefik com TLS):
   ```
   curl -sS https://api.s10.smydi.com.br/healthz | jq
   ```

3. **Recursos**:
   ```
   sshpass -e ssh root@216.144.235.25 'docker stats --no-stream | grep smid10'
   ```
   Confirmar que `mem_usage` está dentro dos `limits` declarados (ver `swarm-stack.yml`).

## Em caso de falha

- `status: "degraded"` com `Access denied`: secrets divergem do init do MariaDB. Refazer secrets na ordem do `deploy/README.md §3` ou apagar volume `smid10_smid10_mariadb_data` (perde dados de teste).
- Certificado Let's Encrypt não emite: confirmar DNS + porta 80 da VPS aberta + log do Traefik (`docker service logs traefik_traefik | tail -50`).
- OOM kill: subir `limits.memory` da API em `swarm-stack.yml` para 192M e redeployar.
- Outros stacks (Chatwoot, Evolution) pesando: `docker stats` para confirmar; se necessário reduzir `smid10_api` replicas para 0 temporariamente.

## Rollback

```
sshpass -e ssh root@216.144.235.25 'docker service update --rollback smid10_api'
```

## Sair sem deixar lixo

Nunca commitar `.env` ou conteúdo de `/tmp/.smid10_creds` no repo. Confirmar `git status` limpo após a sessão.
