# SPEC_REST_API.md — Superfície REST do SMID

> Este SPEC descreve o contrato REST para integrações externas, aplicações mobile e consumidores headless. É agnóstico de linguagem — implemente os endpoints com o verbo, payload e regra de autorização indicados.

---

## 1. Propósito

Expor recursos do SMID para:
- Apps mobile nativos
- Integrações B2B (parceiros, ERPs, BI)
- Webhooks de integrações externas (Evolution, 3C+, GoTo)
- Automação e scripts internos

---

## 2. Autenticação

### 2.1 JWT (padrão)

| Item | Valor |
|------|-------|
| Algoritmo | HS256 |
| Header | `Authorization: Bearer <token>` |
| Claims obrigatórios | `sub` (login), `iat`, `exp`, `system_unit_id`, `groups` |
| Expiração padrão | 8 horas |
| Renovação | Endpoint `/auth/refresh` |

### 2.2 REST Key (integrações confiáveis)

Serviços tipo `AdiantiRecordService` usam `REST_KEY` dedicada por recurso. Deve ser rotacionável e armazenada criptografada.

### 2.3 Webhook Signature

Webhooks externos (Evolution, 3C+) devem validar assinatura HMAC quando suportado pelo provedor. Fallback: validar IP + nome da instância.

---

## 3. Convenções

### 3.1 Formato

- Requisição: `application/json; charset=utf-8`
- Resposta: `application/json; charset=utf-8`
- Datas: ISO 8601 (`YYYY-MM-DDTHH:MM:SS-03:00`)
- Valores monetários: `number` em reais com 2 casas decimais

### 3.2 Respostas Padrão

**Sucesso (200)**:
```json
{ "status": "ok", "data": { ... } }
```

**Erro de validação (400)**:
```json
{ "status": "error", "code": "VALIDATION", "message": "...", "errors": [ ... ] }
```

**Não autenticado (401)**:
```json
{ "status": "error", "code": "UNAUTHENTICATED" }
```

**Não autorizado (403)**:
```json
{ "status": "error", "code": "FORBIDDEN", "message": "..." }
```

**Não encontrado (404)**:
```json
{ "status": "error", "code": "NOT_FOUND" }
```

**Servidor (500)**:
```json
{ "status": "error", "code": "INTERNAL", "message": "..." }
```

### 3.3 Paginação

Query params: `page`, `per_page` (default 20, max 200).

Resposta com meta:
```json
{
  "status": "ok",
  "data": [...],
  "meta": { "page": 1, "per_page": 20, "total": 120 }
}
```

### 3.4 Escopo de Unidade

Toda listagem filtra automaticamente por `system_unit_ids` do usuário autenticado. Registros sem unidade (`unidd_id IS NULL`) são visíveis a todos.

---

## 4. Endpoints por Domínio

> Endpoints abaixo são a forma **canônica recomendada**. A implementação atual em PHP expõe recursos via `AdiantiRecordService` e controllers específicos. Ao reimplementar, siga este contrato.

### 4.1 Autenticação

| Verbo | Endpoint | Função |
|-------|----------|--------|
| POST | `/auth/login` | Autentica login+senha, retorna JWT |
| POST | `/auth/refresh` | Renova JWT |
| POST | `/auth/logout` | Invalida sessão/JWT |
| GET  | `/auth/me` | Dados do usuário autenticado |

### 4.2 Leads

| Verbo | Endpoint | Função |
|-------|----------|--------|
| GET | `/leads` | Lista com paginação e filtros |
| GET | `/leads/{id}` | Detalhe |
| POST | `/leads` | Cria lead (respeita duplicidade técnica/negócio) |
| PUT | `/leads/{id}` | Atualiza |
| DELETE | `/leads/{id}` | Soft delete |
| GET | `/leads/search-by-phone` | Busca por telefone (ver `SPEC_PHONE_SEARCH_BOX`) |

### 4.3 Visitas

| Verbo | Endpoint | Função |
|-------|----------|--------|
| GET | `/visitas` | Lista |
| POST | `/visitas` | Agenda visita (sincroniza status do lead) |
| PUT | `/visitas/{id}` | Atualiza |
| POST | `/visitas/{id}/checkin` | Registra check-in com GPS |
| POST | `/visitas/{id}/historico` | Cria histórico da visita |

### 4.4 Pedidos

| Verbo | Endpoint | Função |
|-------|----------|--------|
| GET | `/pedidos` | Lista |
| POST | `/pedidos` | Cria (dispara comissão e financeiro) |
| PUT | `/pedidos/{id}` | Atualiza |
| PUT | `/pedidos/{id}/status` | Transição de status |

### 4.5 Televendas

| Verbo | Endpoint | Função |
|-------|----------|--------|
| GET | `/televendas/fila` | Fila consolidada |
| POST | `/televendas/contatos` | Registra contato |
| POST | `/televendas/orcamentos` | Cria orçamento |
| POST | `/televendas/orcamentos/{id}/converter` | Converte em pedido |

### 4.6 KM Rodado

| Verbo | Endpoint | Função |
|-------|----------|--------|
| GET | `/km/lotes` | Lista lotes |
| POST | `/km/lotes` | Cria lote de reembolso |
| POST | `/km/lotes/{id}/calcular` | Dispara cálculo GPS |
| POST | `/km/lotes/{id}/aprovar` | Aprova lote |

### 4.7 Integrações

| Verbo | Endpoint | Função |
|-------|----------|--------|
| POST | `/webhooks/evolution` | Recebe evento WhatsApp |
| POST | `/webhooks/3cplus` | Recebe evento de chamada |
| POST | `/webhooks/goto` | Recebe evento GoTo Connect |

Webhooks não autenticam por JWT. Validação é por assinatura ou chave de instância.

---

## 5. Regras Transversais

### 5.1 Idempotência

- POST de webhooks deve ser idempotente por chave externa (`call_id`, `message_id`, `event_id`)
- POST de lançamentos financeiros derivados deve usar chave de origem (`ped_id`, `comissao_id`)

### 5.2 Auditoria

Toda escrita deve gerar entrada em `SystemChangeLog` via `SystemChangeLogTrait` no model consumidor.

### 5.3 Validação

- Backend sempre valida payload completo
- Nunca confiar em validação apenas do cliente
- Validar permissão de unidade em TODAS as operações

### 5.4 Taxa de Requisições

Recomendação: rate limit por token/IP em endpoints de busca e webhooks.

---

## 6. Versionamento

- Prefixo opcional `/v1/` recomendado para novas implementações
- Mudanças incompatíveis exigem nova versão
- Adição de campo em resposta é não-breaking
- Remoção ou renomeação é breaking

---

## 7. Segurança

| Risco | Mitigação |
|-------|-----------|
| JWT vazado | Expiração curta + rotação + blacklist em logout |
| CSRF em sessão cookie | Token anti-CSRF quando usar cookie |
| Injection | Sempre parametrizar queries |
| Mass assignment | Whitelist explícita de atributos gravados |
| IDOR | Validar `unidd_id` e ownership em cada endpoint |
| Rate abuse | Throttling por token/IP |
| Exposição de PII | Não retornar CPF/senha/token em listagens |

---

## 8. Referência à Implementação Atual

| Artefato | Localização |
|----------|-------------|
| Exemplo de REST record service | `app/service/rest/RecordServiceExample.php` |
| Autenticação REST | `app/service/auth/ApplicationAuthenticationRestService.php` |
| Autenticação app | `app/service/auth/ApplicationAuthenticationService.php` |
| Sessão | `app/service/auth/ApplicationSessionService.php` |
| JWT | Biblioteca `firebase/php-jwt` |

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
