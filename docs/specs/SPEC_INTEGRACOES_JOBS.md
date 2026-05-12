# SPEC_INTEGRACOES_JOBS.md

## 1. Visão Geral

### 1.1 Propósito

Este domínio documenta todas as integrações externas e jobs agendados do SMID. Ele cobre canais de entrada de leads (WhatsApp/Evolution), integrações telefônicas (3C+, GoTo Connect), autenticação alternativa (LDAP), jobs do sistema e o mecanismo de agendamento interno.

### 1.2 Integrações Existentes

| Integração | Direção | Protocolo | Propósito |
|-----------|---------|-----------|-----------|
| **Evolution / WhatsApp** | Entrada | Webhook HTTP | Ingresso de leads via mensagem WhatsApp |
| **3C Plus** | Entrada/Saída | REST API | Relatórios de chamadas, leads receptivos |
| **GoTo Connect** | Entrada | OAuth2 + REST | Chamadas e sincronização de leads |
| **WhatsApp Notify** | Saída | REST | Notificação de novos leads para equipe |
| **GPS/Maps** | Saída | REST API | Geocoding e distâncias para KM rodado |
| **LDAP** | Entrada | LDAP/LDAPS | Autenticação alternativa de usuários |
| **reCAPTCHA** | Entrada | REST | Validação anti-bot no login |

---

## 2. Glossário

| Termo | Definição |
|-------|-----------|
| **Webhook** | Chamada HTTP de entrada acionada por evento externo |
| **Ingresso** | Processo de criação ou atualização de lead a partir de fonte externa |
| **Sync Service** | Serviço que reconcilia dados entre SMID e sistema externo |
| **Receptive Rule** | Regra que qualifica chamada receptiva como lead |
| **OAuth2** | Protocolo de autorização delegada usado pelo GoTo |
| **Job** | Tarefa executada periodicamente sem interação humana |
| **Schedule** | Configuração de job recorrente via `SystemSchedule` |
| **Payload** | Corpo bruto do evento/integração, preservado para auditoria |

---

## 3. Evolution / WhatsApp

### 3.1 Arquitetura

```
WhatsApp ──► Evolution Server ──► Webhook POST ──► SMID
                                                     │
                                          EvolutionWebhookService
                                                     │
                                          EvolutionLeadIngressService
                                                     │
                                  ┌──────────────────┴───────────────────┐
                                  ▼                                       ▼
                         Lead criado/duplicado              LeadWhatsappNotificationService
```

### 3.2 Fluxo de Ingresso

```
1. Evolution recebe mensagem WhatsApp
2. POST para endpoint SMID com payload
3. EvolutionWebhookService valida instância e persiste EvolutionWebhookLog
4. EvolutionLeadIngressService:
   a. Normaliza telefone
   b. Verifica duplicidade técnica (20 minutos)
   c. Verifica duplicidade de negócio (6 horas)
   d. Cria Lead ou registra LeadDuplicado
   e. Define status inicial e mídia
5. LeadWhatsappNotificationService notifica equipe
6. EvolutionWebhookEvent registra evento processado
```

### 3.3 Gestão de Instâncias

```
- EvolutionInstanceForm/List: cadastro de instâncias WhatsApp
- EvolutionInstanceConnectModal: QR Code para conexão
- EvolutionInstanceCardsView: visão de status das instâncias
- EvolutionInstanceWizardForm: wizard de configuração
- EvolutionInstanceLog: logs por instância
```

### 3.4 Modelos

| Model | Tabela |
|-------|--------|
| `EvolutionWebhookLog` | Log bruto do webhook |
| `EvolutionWebhookEvent` | Evento processado |
| `EvolutionInstanceLog` | Log por instância |

### 3.5 Services

| Service | Responsabilidade |
|---------|-----------------|
| `EvolutionApiClient` | Comunicação com API Evolution |
| `EvolutionConfigService` | Leitura de configuração |
| `EvolutionCryptoService` | Criptografia de credenciais |
| `EvolutionInstanceService` | Ciclo de vida das instâncias |
| `EvolutionWebhookService` | Recepção e validação de webhooks |
| `EvolutionLeadIngressService` | Criação de leads a partir de eventos |
| `LeadWhatsappNotificationService` | Notificações de entrada para equipe |

---

## 4. 3C Plus

### 4.1 Propósito

Integração com plataforma de telefonia 3C+. Permite sincronizar relatórios de chamadas com leads e contatos do Televendas.

### 4.2 Fluxo

```
1. ThreeCPlusAuthService obtém token de acesso
2. ThreeCPlusCallReportService busca chamadas por período
3. ThreeCPlusLeadSyncService processa chamadas:
   a. Identifica lead por telefone
   b. Aplica ThreeCPlusReceptiveRule se chamada receptiva
   c. Cria/atualiza TelevendasContato
   d. ThreeCPlusProcessedCall evita reprocessamento
4. ThreeCPlusSyncLog registra execução
5. SyncProgramForm permite sincronização manual
```

### 4.3 Modelos

| Model | Tabela |
|-------|--------|
| `ThreeCPlusConfig` | Configuração |
| `ThreeCPlusProcessedCall` | Chamadas já processadas |
| `ThreeCPlusReceptiveRule` | Regras de qualificação receptiva |
| `ThreeCPlusSyncLog` | Log de sincronizações |

### 4.4 Administração

| Controller | Função |
|------------|--------|
| `ThreeCPlusConfigForm` | Configuração de credenciais |
| `ThreeCPlusReceptiveRuleForm/List` | Regras receptivas |

---

## 5. GoTo Connect

### 5.1 Propósito

Integração OAuth2 com GoTo Connect para importar relatórios de chamadas e sincronizar leads.

### 5.2 Fluxo

```
1. GotoConnectOAuthService executa fluxo OAuth2
2. Token armazenado e renovado automaticamente
3. GotoConnectCallReportService coleta chamadas
4. GotoLeadSyncService identifica e cria/atualiza leads
5. GotoConnectSyncLog registra execução
```

### 5.3 Modelos/Services

| Artefato | Responsabilidade |
|----------|-----------------|
| `GotoConnectConfig` | Credenciais OAuth e configuração |
| `GotoConnectSyncLog` | Log de sincronizações |
| `GotoConnectOAuthService` | Fluxo OAuth2 |
| `GotoConnectConfigService` | Leitura de configuração |
| `GotoConnectCryptoService` | Criptografia de secrets |
| `GotoConnectCallReportService` | Coleta de chamadas |
| `GotoLeadSyncService` | Criação/atualização de leads |

### 5.4 Administração

`GotoConnectConfigForm` gerencia credenciais no painel Admin.

---

## 6. Notificações WhatsApp de Leads

### 6.1 Fluxo

```
1. Novo lead criado por ingresso ou manual
2. LeadWhatsappNotificationService verifica configurações
3. Identifica destinatários conforme LeadWhatsappNotificationConfig
4. Envia mensagem via Evolution API
5. LeadWhatsappNotificationLog registra envio
```

### 6.2 Tipos

| Tipo | Gatilho |
|------|---------|
| `NEW_LEAD` | Lead novo criado |
| `DUPLICATE` | Lead identificado como duplicado |

---

## 7. Autenticação LDAP

### 7.1 Fluxo

```
1. Usuário faz login com credenciais corporativas
2. LdapAuthenticationService valida contra diretório LDAP
3. Se válido, retorna SystemUser correspondente
4. ApplicationAuthenticationService conclui carregamento de sessão
```

### 7.2 Configuração

Definido em `app/config/application.ini`:
```
[permission]
auth_service = LdapAuthenticationService
```

---

## 8. Jobs e Agendamentos

### 8.1 Arquitetura

```
SystemSchedule (banco communication)
        │
        ▼
SystemScheduleService::run() ─── disparado por CLI/cron externo
        │
        ├── registra SystemScheduleLog
        └── executa classe/método configurado
```

### 8.2 Tipos de Schedule

| Tipo | Código | Descrição |
|------|--------|-----------|
| Mensal | `M` | Executa no dia do mês especificado |
| Semanal | `W` | Executa no dia da semana especificado |
| Diário | `D` | Executa no horário especificado |
| Fixo | `F` | Executa sempre (sem filtro de data) |

### 8.3 Jobs Conhecidos do Sistema

| Job | Classe | Função |
|-----|--------|--------|
| Agendamento horário padrão | `AgendamentoHorarioDefaultService` | Define horário padrão de agendamento |
| Sincronização de programas | `SyncProgramForm` | Sincroniza programas/permissões |
| Notificações financeiras | `FinNotificacaoService` | Dispara alertas de vencimento |
| Serviço de exemplo | `MyExampleService` / `MyExampleJob` | Template para novos jobs |

### 8.4 Notificações Financeiras por Job

```
1. FinNotificacaoService verifica contas próximas ao vencimento
2. Lê FinNotificacaoConfig para destinatários e janelas
3. Envia alertas via canal configurado
4. FinNotificacaoLog registra envio
```

---

## 9. Segurança de Integrações

| Risco | Mitigação |
|-------|-----------|
| Exposição de chave API | Credenciais criptografadas por serviços `*CryptoService` |
| Payload malicioso | Validação de instância/origem antes de processar |
| Reprocessamento de evento | `ProcessedCall`, `WebhookEvent` evitam duplicidade |
| Credencial OAuth expirada | Renovação automática de token |
| LDAP injection | Escapar caracteres em bind |
| reCAPTCHA bypass | Validação server-side obrigatória |

---

## 10. Referência à Implementação Atual

### 10.1 Controllers Admin de Integrações

| Controller | Localização |
|------------|-------------|
| `EvolutionConfigForm` | `app/control/admin/` |
| `EvolutionInstance*` | `app/control/admin/` |
| `ThreeCPlusConfigForm` | `app/control/admin/` |
| `ThreeCPlusReceptiveRuleForm/List` | `app/control/admin/` |
| `GotoConnectConfigForm` | `app/control/admin/` |
| `WhatsNotificacoesForm` | `app/control/admin/` |
| `SyncProgramForm` | `app/control/admin/` |

### 10.2 Logs de Integrações

| Controller | Localização |
|------------|-------------|
| `EvolutionWebhookEventList` | `app/control/log/` |
| `EvolutionWebhookLogList` | `app/control/log/` |
| `EvolutionLogPayloadView` | `app/control/log/` |
| `EvolutionInstanceLogList` | `app/control/log/` |

### 10.3 Services

Localizados em `app/service/integration/` — todos prefixados por integração.

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
