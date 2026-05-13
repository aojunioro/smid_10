# Deploy — Ambiente de Testes SMID 10 (VPS Docker Swarm)

Runbook para subir o stack `smid10` na VPS compartilhada com Chatwoot, Evolution, Portainer e Traefik.

---

## 1. Arquitetura

```
Internet
  │
  ▼
Traefik (já existente)  ── certresolver=letsencryptresolver, entrypoint=websecure
  │
  ├──  api.s10.smydi.com.br  ──►  smid10_api (Go, porta interna 8080)
  │                                     │
  │                                     ▼
  │                              smid10_mariadb (porta interna 3306, sem expose)
  │
  └──  s10.smydi.com.br      ──►  smid10_web (Next.js — habilitar quando frontend pronto)
```

Rede overlay compartilhada: **`SmydiNet`** (já existe na VPS).

---

## 2. DNS

Antes do primeiro deploy, criar registros A apontando para a VPS (`216.144.235.25`):

```
api.s10.smydi.com.br   A   216.144.235.25
s10.smydi.com.br       A   216.144.235.25
```

O Traefik resolve o desafio HTTP-01 do Let's Encrypt automaticamente.

---

## 3. Secrets

Todos os segredos vivem como **Docker Swarm secrets**, nunca em variáveis de ambiente do compose nem no repositório.

### 3.1 Criar (uma única vez por VPS)

Gere senhas fortes e injete via stdin. Substitua `<...>` pelos valores reais:

```bash
# Senha root do MariaDB (use openssl rand -base64 32)
openssl rand -base64 32 | docker secret create smid10_mariadb_root_password -

# Senha do usuário de aplicação smid10
openssl rand -base64 32 | docker secret create smid10_mariadb_user_password -

# DSNs Go — uma por banco. O host é o nome do serviço Swarm: `mariadb`.
# Formato: user:password@tcp(host:3306)/dbname?parseTime=true&loc=Local&charset=utf8mb4
USER_PWD=$(docker secret inspect smid10_mariadb_user_password --pretty 2>/dev/null || echo) # apenas referência
# Use a MESMA senha que você acabou de gerar acima.
read -s -p 'Senha do usuário smid10 (a que acabou de gerar): ' PWD; echo
printf 'smid10:%s@tcp(mariadb:3306)/smid?parseTime=true&loc=Local&charset=utf8mb4'          "$PWD" | docker secret create smid10_dsn_smid -
printf 'smid10:%s@tcp(mariadb:3306)/permission?parseTime=true&loc=Local&charset=utf8mb4'    "$PWD" | docker secret create smid10_dsn_permission -
printf 'smid10:%s@tcp(mariadb:3306)/log?parseTime=true&loc=Local&charset=utf8mb4'           "$PWD" | docker secret create smid10_dsn_log -
printf 'smid10:%s@tcp(mariadb:3306)/communication?parseTime=true&loc=Local&charset=utf8mb4' "$PWD" | docker secret create smid10_dsn_communication -
unset PWD

# JWT
openssl rand -base64 48 | docker secret create smid10_jwt_secret -
```

Confira:

```bash
docker secret ls | grep smid10
```

### 3.2 Rotacionar um secret

Swarm secrets são imutáveis. Para rotacionar:

```bash
docker secret create smid10_jwt_secret_v2 - <<< 'novo-valor'
# Edite swarm-stack.yml apontando para `smid10_jwt_secret_v2`
docker stack deploy -c deploy/swarm-stack.yml smid10
docker secret rm smid10_jwt_secret   # após confirmar funcionamento
```

---

## 4. Volume de Inicialização do MariaDB

O MariaDB lê arquivos em `/docker-entrypoint-initdb.d/` apenas na primeira inicialização (volume vazio). Use um volume nomeado populado uma única vez:

```bash
# Cria o volume
docker volume create smid10_mariadb_init

# Copia o SQL de inicialização para dentro do volume
docker run --rm \
  -v smid10_mariadb_init:/dst \
  -v $(pwd)/deploy/mariadb-init:/src:ro \
  alpine sh -c 'cp /src/*.sql /dst/ && chmod 644 /dst/*.sql'

# Confere
docker run --rm -v smid10_mariadb_init:/v alpine ls -l /v
```

Para reaplicar (recriar bancos), apague o volume de dados e o stack sobe limpo:

```bash
docker stack rm smid10
docker volume rm smid10_smid10_mariadb_data    # cuidado: apaga TODOS os dados
docker stack deploy -c deploy/swarm-stack.yml smid10
```

---

## 5. Publicar a Imagem da API

Na primeira vez (e em cada push em `main`), o GitHub Actions builda e publica em `docker.io/aojunioro/smid10-api:dev` automaticamente via `.github/workflows/backend-build.yml`.

Build local (para testar antes do CI):

```bash
docker buildx build --platform linux/amd64 \
  -t aojunioro/smid10-api:dev \
  -f backend/Dockerfile \
  --push \
  backend
```

> Para `--push` funcionar localmente: `echo $DOCKER_PASSWORD | docker login docker.io -u $DOCKER_USERNAME --password-stdin`

A imagem é pública (ou privada, conforme configuração no Docker Hub). Se privada, criar pull secret no Swarm:

```bash
docker login docker.io
docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth
```

---

## 6. Deploy

```bash
# Na VPS, dentro de uma checkout do repositório:
docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth
```

Acompanhar:

```bash
docker stack services smid10
docker service ps smid10_mariadb   --no-trunc
docker service ps smid10_api       --no-trunc
docker service logs -f smid10_api
```

---

## 7. Verificação

### 7.1 Direto na rede interna (sem Traefik)

```bash
docker run --rm --network SmydiNet curlimages/curl:latest \
  -sS http://smid10_api:8080/healthz | jq
```

Esperado: status `ok`, todos os 4 bancos respondendo.

### 7.2 Através do Traefik (público)

```bash
curl -sS https://api.s10.smydi.com.br/healthz | jq
```

Primeira requisição pode levar 10–30 s aguardando emissão do certificado Let's Encrypt.

---

## 8. Rollback

```bash
# Voltar para a tag anterior
SMID10_API_TAG=sha-<curto> docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth
```

Ou via `update --rollback`:

```bash
docker service update --rollback smid10_api
```

---

## 9. Recursos e Convivência com Outros Stacks

Os limites em `swarm-stack.yml` foram dimensionados para a VPS atual (3.3 GB RAM, ~1 GB livre):

| Serviço | Reserva | Limite |
|---|---|---|
| mariadb | 192 MB | 512 MB |
| api | 32 MB | 128 MB |
| web (futuro) | 96 MB | 256 MB |

**Total no pior caso**: ~896 MB. Monitore com:

```bash
docker stats --no-stream
```

Se a VPS começar a fazer swap pesado, reduzir `api.replicas` para 0 temporariamente ou aumentar a VPS.

---

## 10. Remoção Completa

```bash
docker stack rm smid10
docker volume rm smid10_smid10_mariadb_data smid10_smid10_mariadb_init
docker secret ls | awk '/smid10_/ {print $2}' | xargs -r docker secret rm
```

---

## 11. Troubleshooting Rápido

| Sintoma | Causa provável | Ação |
|---|---|---|
| `/healthz` retorna 503 com `Access denied` | Senha do secret diverge do usuário no MariaDB | Apague o volume de dados e refaça secrets em ordem |
| Traefik não roteia | Label `traefik.docker.network` errada ou serviço fora da rede `SmydiNet` | Confirme `docker service inspect smid10_api` |
| Certificado não emite | DNS não propagou ou porta 80 bloqueada | `dig api.s10.smydi.com.br`; checar firewall |
| OOM kill da API | Limite muito baixo OU vazamento | Subir `limits.memory` para 192M; checar `docker stats` |
| MariaDB não sobe | Init scripts com erro de sintaxe | `docker service logs smid10_mariadb` |
