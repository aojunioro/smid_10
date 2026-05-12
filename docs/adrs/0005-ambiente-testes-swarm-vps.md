# ADR 0005 — Ambiente de Testes em Docker Swarm na VPS Compartilhada

**Status**: Aceito
**Data**: 2026-05-12
**Decisores**: aojunioro
**Relaciona-se com**: ADR 0001 (Cenário Go+Next), ADR 0002 (Reuso da Base Legada), ADR 0003 (Stack Detalhada), ADR 0004 (Postgres Pós-Cutover)

---

## Contexto

O SMID 10 precisa de um ambiente de execução acessível por URL pública desde a Fase 0.1, por três motivos:

1. O Mac do desenvolvedor tem **pouco espaço em disco** disponível para rodar containers locais.
2. O healthz da API só sai de `degraded` para `ok` com **MySQL real** alcançável.
3. O frontend Next.js e o backend Go consumindo `https://api.s10...` desde o início **eliminam refactor** de CORS, JWT issuer e URLs base mais tarde.

Existe uma VPS Debian 11 (`216.144.235.25`) já operando em produção com vários stacks ativos: **Portainer, Chatwoot, Evolution API, Typebot, MinIO, Postgres 14, pgvector, sites estáticos, Traefik v3.4** — nenhum deles pode parar.

---

## Decisão

Subir um stack **`smid10`** no Docker Swarm já existente da VPS, plugado na rede overlay **`SmydiNet`** que serve o Traefik, com isolamento por limites de recurso para não impactar os stacks de produção.

### Inventário relevante da VPS (capturado em 2026-05-12)

- **OS**: Debian GNU/Linux 11 (bullseye), kernel 5.10
- **Docker**: 28.5.1, Swarm ativo, **single-node manager**
- **Recursos**: 3.3 GB RAM total, ~1 GB livre, 34 GB disco livre
- **Rede overlay compartilhada**: `SmydiNet`
- **Traefik v3.4.0** com:
  - `--providers.swarm=true`
  - `--providers.docker.network=SmydiNet`
  - `--providers.docker.exposedbydefault=false`
  - Entrypoints: `web` (80) → redireciona para `websecure` (443)
  - Certresolver: **`letsencryptresolver`** (HTTP-01), email `aojunioro@gmail.com`
- **Padrão de domínio dos serviços internos**: `*.smydi.com.br`

### Topologia adotada

```
Internet
   │
   ▼
Traefik (já existente)
   │
   ├──  api.s10.smydi.com.br  ──►  smid10_api      (Go, distroless, porta interna 8080)
   │                                     │
   │                                     ▼
   │                              smid10_mariadb   (MariaDB 10.11, sem expose)
   │
   └──  s10.smydi.com.br      ──►  smid10_web      (Next.js — habilitar quando frontend pronto)
```

### Componentes do stack

| Serviço | Imagem | Exposição | Reserva | Limite |
|---|---|---|---|---|
| `mariadb` | `mariadb:10.11` | só rede overlay | 192 MB | 512 MB |
| `api` | `ghcr.io/aojunioro/smid10-api:${SMID10_API_TAG:-dev}` | Traefik (HTTPS, `api.s10.smydi.com.br`) | 32 MB | 128 MB |
| `web` (futuro) | `ghcr.io/aojunioro/smid10-web:${SMID10_WEB_TAG:-dev}` | Traefik (HTTPS, `s10.smydi.com.br`) | 96 MB | 256 MB |

### Decisões específicas

1. **MariaDB 10.11 no lugar de MySQL 8** para os testes:
   - Drop-in compatível com o dialeto Adianti/SMID 8.x.
   - ~40% menos RAM que MySQL 8 (crítico no orçamento de ~1 GB livre).
   - Driver Go `go-sql-driver/mysql` funciona idêntico — zero mudança no código.
2. **Sem porta MySQL exposta**: a API conecta via DNS interno do Swarm (`mariadb:3306`).
3. **Subdomínios**: `s10.smydi.com.br` (web) e `api.s10.smydi.com.br` (API), seguindo o padrão dos serviços existentes (`type.`, `chat.`, `whats.`, `evo.`, `portainer.`).
4. **Segredos** como **Docker Swarm secrets**, montados em `/run/secrets/*`; o `internal/config/config.go` foi estendido para resolver `<NAME>_FILE` automaticamente (padrão Docker oficial).
5. **Registry**: **GitHub Container Registry** (`ghcr.io/aojunioro/smid10-api`), grátis para o repositório privado, autenticação no CI via `GITHUB_TOKEN` sem credencial adicional.
6. **CI**: GitHub Actions builda e publica a imagem em cada push em `main` que toque `backend/` — `.github/workflows/backend-build.yml`.
7. **Deploy**: **manual via SSH** nesta fase inicial (`docker stack deploy -c deploy/swarm-stack.yml smid10 --with-registry-auth`). Reavaliar para deploy automático quando o ritmo de iteração exigir.
8. **Limites de memória explícitos** em todos os serviços para **proteger Chatwoot, Evolution e demais stacks de produção** de OOM kill cruzado.

---

## Implementação Concluída

### Arquivos criados

- `backend/Dockerfile` — multi-stage `golang:1.26-alpine` → `distroless/static-debian12:nonroot`, binário estático `CGO_ENABLED=0`, imagem final ≈ 15 MB.
- `backend/.dockerignore` — exclui `.git`, `.env*`, `vendor/`, etc.
- `deploy/swarm-stack.yml` — stack completo: MariaDB + API com labels Traefik, secrets externos, volumes nomeados, `update_config` rolling com rollback.
- `deploy/mariadb-init/01-schemas.sql` — cria os 4 schemas legados (`smid`, `permission`, `log`, `communication`) em `utf8mb4` e concede grants restritos ao usuário `smid10`.
- `deploy/README.md` — runbook operacional: DNS, criação de secrets, init do volume MariaDB, deploy, verificação interna (curl pela overlay) e externa (via Traefik), rollback, troubleshooting.
- `.github/workflows/backend-build.yml` — pipeline `docker buildx` com cache GHA, push para `ghcr.io`.

### Alterações no código

- `backend/internal/config/config.go` — função `resolveSecret(key)` lê `<KEY>_FILE` quando presente, com fallback para a variável direta. Aplicada a `DB_*_DSN` e `JWT_SECRET`. Testes existentes (`go build`, `go vet`) permanecem verdes.

### Higiene de segurança

- `.gitignore` reforçado: `credenciais.md`, `credentials.md`, `secrets/`, `*.pem`, `*.key`, `deploy/.env*`.
- Arquivo local de credenciais SSH **nunca rastreado** pelo git (verificado via `git ls-files`).

---

## Consequências

### Positivas

- Ambiente de testes público desde a Fase 0.1 sem consumir disco do Mac.
- Reuso da infra Traefik+Swarm já operacional: zero custo de TLS, roteamento e renovação de certificados.
- API e web já apontam para URLs definitivas do ambiente de staging — CORS, base URL do frontend e issuer do JWT corretos desde o dia 1.
- Padrão `_FILE` para secrets é compatível com Kubernetes, Nomad, ECS e Docker Swarm — facilita futura migração de orquestrador.
- CI publicando em GHCR garante imagens versionadas por SHA e tag `dev` para o ambiente.

### Negativas

- Orçamento de RAM apertado: ~896 MB no pior caso do stack `smid10`, contra ~1 GB livre na VPS. Requer monitoramento ativo de `docker stats` durante uso.
- Acoplamento operacional: um bug grave no SMID 10 que ignore os `limits.memory` pode pressionar Chatwoot/Typebot. Mitigado pelos limites declarados, mas exige disciplina.
- Deploy manual via SSH gera fricção; será automatizado em ADR futura quando o pipeline estabilizar.

### Neutras

- Não alteramos nada nos stacks existentes (`traefik`, `portainer`, `chatwoot-nestor`, `evolution`, `typebot`, `minio`, `postgres`, `pgvector`, sites).
- O `postgres:14` existente fica intocado; pode virar candidato natural para o destino do **ADR 0004** (cutover) quando chegar a hora.

---

## Alternativas Consideradas

| Alternativa | Por que não agora |
|---|---|
| **Rodar tudo local no Mac via `docker compose`** | Pouco espaço em disco; obriga o frontend a apontar para `localhost` (refactor depois). |
| **`docker compose` puro na VPS, sem Swarm** | Vai contra o ambiente: Traefik está configurado com `--providers.swarm=true` e Portainer gerencia stacks Swarm. |
| **Stack separado em VPS dedicada** | Custo financeiro adicional sem ganho técnico nesta fase. |
| **MySQL 8 oficial** | Consome ~40% mais RAM que MariaDB 10.11 com a mesma compatibilidade prática para o legado Adianti. Reavaliar se aparecer dependência específica do MySQL 8 (raro). |
| **Portainer Stack from Git + webhook** | Boa opção, mas adiciona acoplamento com a UI do Portainer. Deixado como evolução natural quando o deploy manual ficar tedioso. |
| **GitHub Actions com SSH para deploy automático** | Exige guardar chave SSH como secret no GitHub e abrir vetor de ataque. Reavaliar quando o time crescer. |
| **Registry próprio (`registry:2`) na VPS** | Mais 50 MB de RAM e mais um ponto de falha; GHCR resolve sem custo. |

---

## Limites

| O que **PODE** fazer agora | O que **NÃO PODE** fazer agora |
|---|---|
| Atualizar `swarm-stack.yml` e redeployar | Remover `limits.memory` dos serviços `smid10_*` |
| Subir novos serviços auxiliares no stack `smid10` | Tocar nos stacks `traefik`, `portainer`, `chatwoot-nestor`, `evolution`, etc. |
| Trocar a imagem da API por uma nova tag `dev` ou `sha-*` | Publicar segredos no `swarm-stack.yml` ou em `.env` rastreado |
| Adicionar middlewares Traefik específicos do SMID 10 | Mudar a configuração global do serviço `traefik` |
| Criar/rotacionar Swarm secrets do prefixo `smid10_*` | Reaproveitar secrets de outros stacks |
| Apagar o volume `smid10_mariadb_data` para reset | Apagar volumes de outros stacks |

---

## Critérios para Revisão

Reavaliar esta ADR quando **qualquer um** for verdadeiro:

- A VPS atingir uso de RAM > 85% sustentado por mais de 24 h.
- O deploy manual via SSH passar a ser executado mais de 3× por semana.
- Surgir necessidade de um ambiente de produção separado (e não apenas staging).
- Migração definitiva para Postgres (gatilho do ADR 0004) — o stack vai precisar ser revisitado para trocar `mariadb` por `postgres`.

---

## Referências

- `deploy/README.md` — runbook operacional completo
- `deploy/swarm-stack.yml` — definição declarativa do stack
- `.github/workflows/backend-build.yml` — pipeline de build e push
- ADR 0002 — exigência de compatibilidade com MySQL legado
- ADR 0004 — destino pós-cutover
