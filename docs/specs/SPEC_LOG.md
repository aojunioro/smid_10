# SPEC_LOG.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Log** registra acessos, requisições, SQLs, alterações de dados, agendamentos e eventos de integrações. Ele fornece rastreabilidade técnica e auditoria funcional para suporte, compliance e diagnóstico.

### 1.2 Entidades Principais

| Entidade | Descrição |
|----------|-----------|
| **SystemAccessLog** | Registro de acesso/login |
| **SystemRequestLog** | Registro de requisições HTTP/controladores |
| **SystemSqlLog** | Registro de SQL executado |
| **SystemChangeLog** | Auditoria de alteração de registros |
| **SystemScheduleLog** | Log de tarefas agendadas |
| **SystemAccessNotificationLog** | Notificações de acesso |
| **EvolutionWebhookLog/Event** | Logs de webhook Evolution |
| **Integration Sync Logs** | Logs de 3C+, GoTo, WhatsApp e notificações |

---

## 2. Glossário

| Termo | Definição |
|-------|-----------|
| **Audit Log** | Registro de alteração de dados de negócio |
| **Request Log** | Registro de requisição e parâmetros técnicos |
| **SQL Log** | Registro de instruções SQL para diagnóstico |
| **Access Log** | Registro de autenticação/acesso |
| **Schedule Log** | Execução de rotinas agendadas |
| **Payload Log** | Corpo bruto ou tratado de integração externa |

---

## 3. Fluxos Principais

### 3.1 Auditoria por SystemChangeLogTrait

```
1. Model usa SystemChangeLogTrait
2. Ao salvar/excluir, trait detecta alterações
3. SystemChangeLogService persiste:
   - tabela alvo
   - chave do registro
   - coluna alterada
   - valor anterior/novo
   - usuário
   - data/hora
4. Painéis de auditoria exibem histórico por entidade
```

### 3.2 Log de Acesso

```
1. Usuário autentica ou tenta autenticar
2. SystemAccessLogService registra login, IP, user-agent e resultado
3. SystemAccessLogList permite consulta administrativa
```

### 3.3 Log de Requisição

```
1. Request entra no framework
2. SystemRequestLogService registra classe, método, parâmetros e tempo
3. SystemRequestLogList/View exibe detalhes
```

### 3.4 Log SQL

```
1. Operações SQL relevantes são interceptadas/registradas
2. SystemSqlLogService persiste instrução, contexto e usuário
3. SystemSqlLogList permite auditoria técnica
```

### 3.5 Logs de Integrações

```
1. Webhook ou sync externo é recebido
2. Payload bruto e status de processamento são salvos
3. Listas específicas permitem inspeção e troubleshooting
4. Falhas ficam rastreáveis por evento/instância
```

---

## 4. Visualizações

| Tela | Função |
|------|--------|
| `SystemLogDashboard` | Painel geral de logs |
| `SystemAccessLogList` | Acessos |
| `SystemRequestLogList/View` | Requisições |
| `SystemSqlLogList` | SQL |
| `SystemChangeLogView` | Alterações |
| `SystemScheduleLogList` | Agendamentos |
| `SystemSessionVarsView` | Sessão atual |
| `EvolutionWebhook*` | Webhooks Evolution |
| `EvolutionInstanceLogList` | Logs por instância |

---

## 5. Decisões Arquiteturais

### ADR-001: Auditoria por Trait

**Decisão**: Models críticos usam `SystemChangeLogTrait`.

**Consequências**:
- Baixo acoplamento nos controllers
- Auditoria consistente por entidade
- Requer models bem mapeados

### ADR-002: Logs Técnicos em Domínio Separado

**Decisão**: Logs técnicos ficam em `app/model/log` e `app/control/log`.

**Consequências**:
- Separação de responsabilidade
- Consultas administrativas centralizadas

### ADR-003: Não Customizar Log Upstream sem Necessidade

**Decisão**: Controllers de log devem ficar alinhados ao template limpo sempre que possível.

**Consequências**:
- Upgrade do Adianti/adminbs5 mais seguro
- Customizações devem ser pontuais e documentadas

### ADR-004: Payloads de Integração Preservados

**Decisão**: Logs de webhooks/sync mantêm payload/status.

**Consequências**:
- Reprocessamento e diagnóstico viáveis
- Atenção a dados sensíveis

---

## 6. Referência à Implementação Atual

### 6.1 Controllers

`app/control/log/SystemAccessLogList.php`, `SystemChangeLogView.php`, `SystemLogDashboard.php`, `SystemRequestLogList.php`, `SystemRequestLogView.php`, `SystemScheduleLogList.php`, `SystemSessionVarsView.php`, `SystemSqlLogList.php`, `EvolutionWebhookEventList.php`, `EvolutionWebhookLogList.php`, `EvolutionLogPayloadView.php`, `EvolutionInstanceLogList.php`.

### 6.2 Models/Services

| Artefato | Função |
|----------|--------|
| `SystemAccessLog` | Acessos |
| `SystemRequestLog` | Requisições |
| `SystemSqlLog` | SQL |
| `SystemChangeLog` | Alterações |
| `SystemScheduleLog` | Rotinas |
| `SystemChangeLogTrait` | Instrumentação de models |
| `System*LogService` | Persistência/consulta de logs |

---

## 7. Segurança

- Logs podem conter dados pessoais e payloads sensíveis
- Acesso restrito a Admin/Suporte técnico autorizado
- Mascarar senhas, tokens e credenciais
- Retenção de logs deve ser governada por política operacional

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
