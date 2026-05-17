---
name: smid-deploy-vps
description: Deploys or redeploys the smid10 Docker Swarm stack on the VPS, runs smoke tests on api.s10.smydi.com.br, and follows deploy/README.md. Use for deploy, redeploy API, staging smoke-test, or swarm-stack.yml changes.
---

# SMID — Deploy VPS (Swarm)

Referencia canonica: `deploy/README.md`.

## Pre-checagem

- Confirmar com usuario antes de SSH/deploy em producao compartilhada.
- Stack name: `smid10`.
- Nunca colocar senhas no compose; usar Swarm secrets.

## Redeploy API (stack existente)

```bash
docker service update --image <registry>/smid10-api:<tag> smid10_api --with-registry-auth
docker service ps smid10_api --no-trunc | head -10
```

## Smoke-test obrigatorio

1. Interno: `curl -s http://smid10_api:8080/healthz` (na rede overlay).
2. Externo: `curl -s https://api.s10.smydi.com.br/healthz` — esperado `status: ok` nos 4 bancos.

## Primeiro deploy

Seguir `deploy/README.md` secoes DNS, secrets, volume init MariaDB, `docker stack deploy`.

## Apos deploy

- Atualizar handoff secao deploy se relevante.
- Nao commitar credenciais.

Imagem CI: ver `.github/workflows/backend-build.yml` e registry atual no README/deploy.
