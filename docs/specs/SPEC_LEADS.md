# SPEC_LEADS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Leads** representa o núcleo comercial do SMID — o ponto de entrada do funil de vendas. Um lead é um potencial cliente que entrou em contato com a empresa por telefone, WhatsApp, formulário web, callcenter ou outro canal. O sistema deve capturar, qualificar, distribuir, acompanhar e converter este lead até a venda ou descarte definitivo.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Lead** | Registro principal do potencial cliente | 1 por cliente em negociação |
| **Visita** | Agendamento/visita técnica vinculada ao lead | 0..N |
| **Pedido** | Venda concretizada vinculada ao lead | 0..N |
| **LeadStatus** | Catálogo de situações no funil | catálogo |
| **LeadStatusCategoria** | Categorização funcional L/V dos status | catálogo |
| **LeadEventTent** | Histórico de tentativas de contato | 0..N por lead |
| **LeadHistorico** | Linha do tempo agregada do lead | 0..N por lead |
| **LeadDuplicado** | Registro de tentativa duplicada (janela 6h) | 0..N por lead |
| **LeadMeio** | Canais de origem (catálogo) | catálogo |
| **Midia** / **MidiaTipo** | Mídias e tipos (catálogos) | catálogo |
| **LeadMotivoPend** | Motivos de pendência (catálogo) | catálogo |
| **LeadMotivoPerd** | Motivos de perda (catálogo) | catálogo |
| **Equipe** | Agrupamento de atendentes/representantes | catálogo |
| **LeadWhatsappNotificationConfig** | Configuração de notificação por usuário | 0..N |
| **LeadWhatsappNotificationLog** | Log/histórico de notificações enviadas | 0..N |

### 1.3 Relacionamentos

```
Lead 1 ──── 0..N Visita
Lead 1 ──── 0..N Pedido
Lead 1 ──── 0..N LeadEventTent (tentativas)
Lead 1 ──── 0..N LeadDuplicado
Lead 1 ──── 0..N LeadHistorico
Lead N ──── 1   Midia (origem comercial)
Lead N ──── 1   LeadMeio (canal técnico)
Lead N ──── 1   LeadStatus → LeadStatusCategoria (L/V)
Lead N ──── 1   SystemUnit (unidade de atendimento)
Lead N ──── 1   SystemUser (atendente / login_recep)
Lead N ──── 0..1 SystemUser (supervisor / login_super)
Equipe N ──── N SystemUser (membros)
```

### 1.4 Integrações com outros Módulos

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| VISITAS | Bidirecional | Espelhamento de status; visitas controlam leads quando em categoria V |
| PEDIDOS | Leitura | Pedidos vinculados a leads; dashboard de evolução |
| COMUNICACAO | Entrada | Webhooks Evolution (WhatsApp) criam leads automaticamente |
| TELEVENDAS | Entrada | Integrações 3C+ e GoTo importam leads de callcenter |
| FINANCEIRO | Leitura | Análise de conversão e comissões |
| LOG / AUDITORIA | Leitura | Recuperação de status anterior via `system_change_log` |
| PERMISSOES | Leitura | Definição de grupos (atendente, supervisor, representante) |

---

## 2. Glossário de Negócio

### 2.1 Termos do Domínio

| Termo | Definição | Sinônimos |
|-------|-----------|-----------|
| **Lead** | Registro de um potencial cliente que entrou em contato | Prospect, Interessado |
| **Atendente** | Usuário (grupo 2) que recebe e qualifica o lead | Recepcionista, Operador |
| **Representante** | Usuário (grupo 4) que realiza a visita técnica | Terapeuta, Vendedor externo |
| **Supervisão** | Usuário (grupo 5) que acompanha e distribui leads | Supervisor, Coordenador |
| **Tentativa** | Registro de tentativa de contato com o lead | Tentativa de ligação |
| **Contato OK** | Flag indicando contato bem-sucedido | Falado, Contactado |
| **Meio** | Canal técnico de entrada (telefone, WhatsApp, etc.) | Canal, Origem técnica |
| **Mídia** | Campanha ou veículo publicitário (Google, Instagram, etc.) | Campanha, Fonte |
| **Equipe** | Grupo de atendentes/representantes sob uma supervisão | Time |
| **Cortina** | Side panel lateral do SMID para edições/filtros | Side Panel |
| **Modo Lead (L)** | Aba Lead em modo editável (combo de status) | Modo livre |
| **Modo Visita (V)** | Aba Lead em modo readonly (badge sincronizado) | Modo bloqueado |

### 2.2 Categorias de Status (Matriz Funcional L/V)

A categorização separa status por tipo funcional via `lead_status_categoria.status_tipo`:

| Tipo | Código | Cor padrão | Descrição |
|------|--------|------------|-----------|
| **L** | NOVO | #6c757d | Lead recém-criado |
| **L** | CONTACTO | #17a2b8 | Lead contactado com sucesso |
| **L** | AGUARDANDO | #fd7e14 | Aguardando retorno do cliente |
| **L** | PERDIDO | #dc3545 | Lead perdido (não converteu) |
| **L** | INVALIDO | #6c757d | Lead inválido (dados errados/duplicado) |
| **V** | AGENDADO | #007bff | Visita agendada |
| **V** | EM_ATENDIMENTO | #ffc107 | Visita em andamento (com representante) |
| **V** | VENDIDO | #28a745 | Venda concluída |

**Mapeamento HAIFLEX (referência histórica):**

| ID | Nome | Categoria | Tipo |
|----|------|-----------|------|
| 1 | Pendente | NOVO | L |
| 5 | Retorno | CONTACTO | L |
| 10 | Falado | CONTACTO | L |
| 2 | Pré-Agendado | AGENDADO | V |
| 3 | Agendado | AGENDADO | V |
| 4 | Andamento | EM_ATENDIMENTO | V |
| 11 | Reagendado | AGENDADO | V |
| 7 | Vendido | VENDIDO | V |
| 8 | Entregue | VENDIDO | V |
| 6 | Perdida | PERDIDO | L |
| 9 | Repicada | PERDIDO | L |

### 2.3 Janelas Temporais e Indicadores

| Janela | Tempo | Propósito |
|--------|-------|-----------|
| Duplicidade técnica | 20 minutos | Hash de mensagem + telefone (Evolution) |
| Janela de retorno | 6 horas | Tempo mínimo para registrar duplicado de negócio |
| Indicador de inatividade | 3 horas | Badge amarelo (atenção) |
| Indicador de urgência | 5+ horas | Badge vermelho (crítico) |
| TTL de lookups | 900 segundos (15min) | Cache de catálogos auxiliares |
| Limite Kanban | 400 cards | Performance da visão Kanban |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Lead (Entrada)

#### 3.1.1 Via WhatsApp (Evolution Webhook)

**Ator**: Sistema (assíncrono)
**Pré-condições**: Webhook configurado e instância Evolution ativa

```
1. Evolution envia webhook → SMID grava EvolutionWebhookEvent (audit trail)
2. EvolutionLeadIngressService::processMessageUpsert() é chamado
3. Verifica se é mensagem fromMe (própria) → IGNORA
4. Extrai telefone do remoteJid e normaliza para formato BR
5. Adquire lock distribuído por telefone (anti-race-condition)
6. Verifica hash de mensagem na janela de 20min → IGNORA se duplicado técnico
7. Busca lead existente por telefone (com candidatos legados)
8. SE lead NÃO existe:
   8.1 Cria Lead com defaults (meio=WhatsApp, status=NOVO/Pendente)
   8.2 Captura foto de perfil do WhatsApp via API Evolution
   8.3 Armazena URL em whatsapp_profile_photo_url
   8.4 Dispara LeadWhatsappNotificationService::notifyNewLead()
   8.5 Retorna LEAD_CREATED
9. SE lead existe E criado há menos de 6h: IGNORA (janela retorno)
10. SE lead existe E criado há mais de 6h:
    10.1 Verifica se já existe duplicado recente
    10.2 SE não existe: cria LeadDuplicado + dispara notificação
    10.3 Retorna DUPLICATE_CREATED
```

**Regras de Negócio**:
- RN-001: Mensagens fromMe são ignoradas
- RN-002: Telefone normalizado em formato brasileiro (10/11 dígitos)
- RN-003: Lock distribuído evita criação dupla em paralelo
- RN-004: Janela técnica (20min/hash) ≠ janela de negócio (6h/telefone)
- RN-005: Foto de perfil capturada quando disponível na API

#### 3.1.2 Via Callcenter (3C+ / GoTo)

**Ator**: Sistema (sincronização periódica) ou Operador (popup)

```
1. Integração recebe evento de chamada (recebida/perdida)
2. Extrai número do caller
3. ThreeCPlusLeadSyncService / GotoLeadSyncService:
   3.1 Busca lead existente
   3.2 SE não existe: cria com meio=Telefone
   3.3 Registra ThreeCPlusProcessedCall (auditoria)
   3.4 Vincula gravação se disponível
4. Popup de atendimento exibe dados do lead ao operador
5. ThreeCPlusSyncPeriodForm permite sincronização manual de período
```

#### 3.1.3 Via Cadastro Manual (LeadsForm)

**Ator**: Atendente
**Pré-condições**: Usuário autenticado com permissão

```
1. Atendente acessa Leads > Novo Lead (cortina lateral)
2. Preenche fone1 (obrigatório) — máscara aplicada
3. onCheckPhone (exit action):
   3.1 Verifica duplicidade por telefone
   3.2 SE existe: alerta com link para abrir lead existente
4. Preenche dados do cliente (nome, profissão, idade, patologia)
5. Opcionalmente preenche dados do acompanhante/familiar
6. Seleciona Meio (obrigatório) — combo de canais técnicos ativos
7. Seleciona Mídia (obrigatório) — combo de campanhas ativas
8. Seleciona Unidade (obrigatório) — restrita às unidades do usuário
9. Seleciona Atendente (obrigatório) — pré-preenchido com login atual
10. Opcionalmente seleciona Supervisão
11. Aba Endereço: CEP com onLookupCep (auto-complete ViaCEP)
12. Salva: status inicial sempre NOVO (resolvido por findDefaultLeadStatusId)
```

**Regras de Negócio**:
- RN-006: Telefone deve ter 10 ou 11 dígitos
- RN-007: Campos texto têm limites: nome (100), profissão (100), patologia (255), obs (300), referencia (500)
- RN-008: Status inicial é resolvido dinamicamente do catálogo L (não hardcoded)
- RN-009: Lead é "guardado" em sessão (CURRENT_ID_SESSION_KEY) com guard token

### 3.2 Qualificação e Distribuição

#### 3.2.1 Edição Individual

```
1. Atendente abre lead via LeadsList ou Kanban
2. LeadsForm carrega em cortina com 3 abas (Lead, Cliente, Endereço)
3. Subcontrollers carregam em paralelo:
   - LeadVisitaDetalhes (visitas)
   - HistoricosController (histórico)
   - PosVisitaController (pós-visita)
   - PedidosController (pedidos)
4. Atendente atualiza campos:
   - Marca Contato OK (incrementa tentativas)
   - Seleciona Motivo de Pendência (combo restrito a categorias L)
   - Seleciona Motivo de Perda (combo restrito a categorias L)
   - Atualiza Status (combo no Modo L; readonly badge no Modo V)
5. onSave persiste:
   5.1 Validação backend (incluindo TMaxLengthValidator)
   5.2 LeadVisitStatusSyncService::resolveLeadSaveContext()
   5.3 Decide modo (L/V) com base na visita ativa válida
   5.4 Persiste lead + endereço + tentativa se aplicável
   5.5 Reaplica modo visual (combo vs badge)
```

#### 3.2.2 Distribuição em Lote (LeadsDistribuir)

**Ator**: Supervisor
**Pré-condições**: Usuário no grupo Supervisão (5)

```
1. Supervisor acessa Leads > Distribuir
2. Sistema exibe leads em batch (paginação 20/página)
3. Filtros disponíveis: status visual, DDD, unidade, atendente
4. Para cada linha, supervisor pode alterar:
   - Atendente (TCombo dinâmico)
   - Mídia (TCombo)
   - Unidade (TCombo)
   - Meio (TCombo)
5. Indicadores visuais:
   - Distribution Status: completo/parcial/pendente
   - Cor baseada nos 4 campos críticos preenchidos
6. Salvar Todos persiste alterações em massa em transação única
7. LeadDistributionStats agrega contadores por status visual
```

**Regras de Negócio**:
- RN-010: Apenas grupo 5 (Supervisão) tem acesso à tela
- RN-011: Atendentes só podem ver leads atribuídos a eles
- RN-012: Distribution Status é calculado em tempo real (sem persistência)

### 3.3 Evolução do Lead (Máquina de Estados)

#### 3.3.1 Ciclo L (Sem Visita)

```
NOVO (Pendente) ──tentativa──▶ CONTACTO (Falado/Retorno)
                                       │
                ┌─────cliente aceita──▶ AGENDADO (entra em V)
                │
                ├─────cliente recusa──▶ PERDIDO
                │
                └─────dados errados──▶ INVALIDO
```

#### 3.3.2 Ciclo V (Com Visita Válida)

```
AGENDADO (Pré-Agendado → Agendado)
    │
    ├──representante atribuído──▶ EM_ATENDIMENTO (Andamento)
    │                                  │
    │                                  ├──visita realizada──▶ VENDIDO
    │                                  │
    │                                  └──não realizada──▶ REAGENDADO ──▶ AGENDADO
    │
    └──visita cancelada──▶ volta para último L (ex: PERDIDO/RETORNO)
```

### 3.4 Sincronização Lead × Visita (Espelhamento Bidirecional)

**Conceito Central**: O status do lead e da visita são espelhados **apenas** quando a visita é "válida para sincronismo".

**Visita Válida para Sincronismo**:
- Possui `stts_lead` preenchido E
- Status pertence a categoria `AGENDADO` ou `VENDIDO`

**Fluxos**:
```
Alteração em Lead.status_id (Modo L editável):
  → SE existe visita ativa válida: propaga para visita.stts_lead
  → SE visita não é válida: lead permanece independente

Alteração em Visita.stts_lead (categoria V):
  → Propaga para lead.status_id
  → Lead entra em Modo V (badge readonly)

Exclusão da última visita:
  → Recupera último status L do system_change_log
  → Filtra logs com class_name in ['LeadVisitaDetalhes', 'VisitaAtribuirForm']
  → Restaura lead para esse status
  → Lead volta para Modo L editável

Carregamento de tela (load):
  → NUNCA persiste alterações automáticas
  → Apenas resolve estado visual
  → Bloqueia 'definir_valor' / 'limpar' do motor granular no bootstrap
```

**Regras de Negócio**:
- RN-013: Sincronização ocorre só para categorias AGENDADO ou VENDIDO
- RN-014: Modo visual reflete autoridade real, não apenas existência de visita
- RN-015: Nenhuma escrita em banco apenas por carregar tela
- RN-016: Lead nunca fica sem status (fallback dinâmico no catálogo L)
- RN-017: Recuperação de último L usa system_change_log com filtro de contexto

### 3.5 Auditoria de Lead (LeadAuditLogPanel)

```
1. Atendente/Supervisor abre painel de auditoria do lead
2. Sistema consulta system_change_log filtrando:
   - target_record_id = lead.id
   - target_table in ['leads', 'endereco', 'visitas']
3. Exibe DataGrid com:
   - Data/hora da alteração
   - Usuário responsável
   - Tabela alterada (rotulada: Lead/Endereço/Visita)
   - Coluna alterada (com lookup para nomes amigáveis)
   - Valor anterior → valor novo (resolvendo IDs para nomes)
   - Origem (controller que originou: Form/Distribuir/Atribuir)
4. columnLookupMap traduz IDs para nomes:
   - status_id → status (LeadStatus)
   - midia_id → midia (Midia)
   - unidd_id → unidade (SystemUnit)
   - meio_id → meio (LeadMeio)
   - mot_pend_id → motivo_pend (LeadMotivoPend)
   - mot_perd_id → motivo_perd (LeadMotivoPerd)
```

**Regras de Negócio**:
- RN-018: Auditoria é somente leitura (não permite edição/reversão)
- RN-019: SystemChangeLogTrait registra automaticamente em afterStore/afterDelete

### 3.6 Notificação de Novo Lead (WhatsApp)

```
1. Lead criado via Evolution dispara LeadWhatsappNotificationService::notifyNewLead()
2. Service busca recipientes configurados:
   - LeadWhatsappNotificationConfig por unidade/grupo
   - Filtra usuários ativos com permissão
3. Para cada recipiente:
   3.1 Cria LeadWhatsappNotificationLog (visualized = N)
   3.2 Mensagem padrão: "Novo Lead inserido no sistema"
4. Sistema exibe badge no menu com countPendingForUser
5. Usuário visualiza notificação → marca como lida (visualized = S)
```

**Variantes**:
- TYPE_NEW_LEAD: novo lead criado
- TYPE_DUPLICATE: lead duplicado registrado

---

## 4. Estados e Transições

### 4.1 Máquina de Estados Completa

```
┌─────────┐    ┌───────────┐    ┌────────────┐    ┌──────────────┐
│  NOVO   │───▶│ CONTACTO  │───▶│  AGENDADO  │───▶│ EM_ATENDIM.  │
│Pendente │    │Falado/Ret │    │Pré-Ag/Ag.  │    │  Andamento   │
└─────────┘    └───────────┘    └────────────┘    └──────────────┘
     │              │ │                │                  │
     │              │ ▼                │                  │
     │              │ ┌──────────────┐ │                  │
     │              │ │  AGUARDANDO  │ │                  │
     │              │ └──────────────┘ │                  │
     │              │       │          │                  │
     ▼              ▼       ▼          ▼                  ▼
┌─────────┐   ┌─────────────────────────────┐    ┌──────────────┐
│INVALIDO │   │            PERDIDO          │    │   VENDIDO    │
└─────────┘   │     (Perdida / Repicada)    │    │Vendido/Entreg│
              └─────────────────────────────┘    └──────────────┘
```

### 4.2 Transições Permitidas

| De | Para | Gatilho | Regra |
|----|------|---------|-------|
| NOVO | CONTACTO | Marcar Contato OK | Incrementa tentativas |
| NOVO/CONTACTO | AGENDADO | Criar visita válida (V) | Lead entra em Modo V |
| QUALQUER L | PERDIDO | Selecionar motivo_perd_id | Manual |
| QUALQUER L | INVALIDO | Marcar como inválido | Manual |
| CONTACTO | AGUARDANDO | Cliente pediu retorno | Manual |
| AGUARDANDO | CONTACTO | Retorno realizado | Manual |
| AGENDADO | EM_ATENDIMENTO | Atribuir representante | VisitaAtribuirForm |
| EM_ATENDIMENTO | VENDIDO | Confirmar venda na visita | LeadVisitaDetalhes |
| EM_ATENDIMENTO | AGENDADO | Remover representante | VisitaAtribuirForm |
| QUALQUER V | último L | Excluir última visita | Recupera de system_change_log |

### 4.3 Automações Assistivas (Não Coercitivas)

| Condição | Ação Sugerida | Origem |
|----------|---------------|--------|
| Visita ganha dt_visita+hr_visita | Sugere AGENDADO | LeadVisitaGranularActions |
| Visita ganha login_repre | Sugere ANDAMENTO | LeadVisitaGranularActions |
| Visita perde login_repre | Sugere AGENDADO | LeadVisitaGranularActions |
| Visita limpa data/hora | Permite limpar status | LeadVisitaGranularActions |

**Princípio**: O usuário sempre pode sobrescrever a sugestão. O sistema **nunca** força status no carregamento de tela.

---

## 5. Visualizações Alternativas

### 5.1 LeadsList (Listagem Padrão)

**Tipo**: DataGrid paginado com filtros avançados em cortina

**Funcionalidades**:
- Busca rápida (server-side) com persistência em sessão
- Filtros avançados: status, mídia, unidade, atendente, motivos, datas, contato_ok
- Botões de status visual coloridos (filtro rápido por categoria)
- Timeline de tentativas inline (transformer)
- Indicadores de tempo (cores conforme janelas 3h/5h)
- Foto de perfil WhatsApp no card
- Exportação Excel (TTableWriterXLS)
- Paginação mobile (SMIDMobilePagination)

### 5.2 LeadsKanbanView (Funil Kanban)

**Tipo**: Visão Kanban com colunas dinâmicas por status

**Funcionalidades**:
- Reusa filtros de LeadsList via sessão (LeadVisionFilterService)
- Colunas dinâmicas por status efetivo dos leads
- Cards com: nome, telefone, atendente, tentativas, última visita
- Drag-and-drop entre colunas (onUpdateItemDrop)
- Limite de segurança: MAX_CARDS = 400
- Acessível via ícone discreto no LeadsList (theme.js)
- Abre em aba interna em desktop (allow_page_tabs)

### 5.3 LeadVisitNextDayDashboard (Visitas Amanhã)

**Tipo**: Dashboard operacional para visualizar visitas do dia seguinte

**Funcionalidades**:
- Filtros por data, atendente, representante, unidade
- Gráficos de distribuição (TBarChart)
- Agrupamento por atendente/representante/região
- Indicadores de carga operacional

### 5.4 LeadAuditLogPanel (Auditoria)

**Tipo**: DataGrid em cortina com histórico de alterações

**Funcionalidades**:
- Filtra system_change_log por lead.id
- Resolve IDs para nomes via lookups
- Mostra origem (controller) e usuário
- Não permite edição (somente leitura)

### 5.5 LeadDuplicadoList (Duplicados)

**Tipo**: Listagem dedicada de leads duplicados

**Funcionalidades**:
- Visualização de lead_duplicados (registros da janela de 6h)
- Vincula ao lead original
- Ações: abrir lead original, marcar como tratado

### 5.6 DashboardEvolucao (Lead Time)

**Tipo**: Dashboard analítico de tempo de conversão

**Métricas**:
- T0: criado_em do lead
- T1: criado_em da visita (agendamento)
- T2: dt_visita + hr_visita (efetivação)
- T3: dt_ped do pedido (venda)

**Visualizações**:
- Funil temporal entre etapas (média/mediana)
- Distribuição por faixa de tempo
- Comparativo por mídia/unidade/atendente
- Grade detalhada anonimizada

---

## 6. Cadastros Auxiliares

### 6.1 Status e Categorias

| Controller | Função |
|------------|--------|
| `LeadStatusForm` / `LeadStatusList` | CRUD de status com vínculo a categoria |
| `LeadCategStatusForm` / `LeadCategStatusList` | CRUD de categorias funcionais |

**Regras**:
- Status protegidos (sistema) não podem ter nome alterado
- Categoria de sistema (sistema='S') não pode ser deletada
- Deve existir pelo menos um status com `stt_inicial = S`
- `LeadStatusHelper::clearCache()` deve ser chamado após alterações

### 6.2 Mídias (Origens Comerciais)

| Controller | Função |
|------------|--------|
| `MidiaForm` / `MidiaList` | CRUD de mídias (Google, Instagram, etc.) |
| `MidiaTipoForm` / `MidiaTipoList` | CRUD de tipos de mídia (orgânico, pago, etc.) |

### 6.3 Meios (Canais Técnicos)

| Controller | Função |
|------------|--------|
| `LeadMeioForm` / `LeadMeioList` | CRUD de meios (Telefone, WhatsApp, etc.) |

### 6.4 Motivos

| Controller | Função |
|------------|--------|
| `LeadMotivoPendForm` / `LeadMotivoPendList` | Motivos de pendência (com cor) |
| `LeadMotivoPerdForm` / `LeadMotivoPerdList` | Motivos de perda (com cor) |

### 6.5 Equipes

| Controller | Função |
|------------|--------|
| `EquipeForm` / `EquipeList` | CRUD de equipes (atendentes + supervisão + cor) |

**Regras**:
- Equipe vincula supervisor a múltiplos atendentes
- Cor da equipe é usada em indicadores visuais
- Lead herda equipe via atendente atribuído

---

## 7. Sistema de Notificações WhatsApp

### 7.1 Componentes

| Arquivo | Responsabilidade |
|---------|------------------|
| `LeadWhatsappNotificationService` | Lógica de notificação |
| `LeadWhatsappNotificationConfig` | Configuração por usuário/unidade |
| `LeadWhatsappNotificationLog` | Auditoria/histórico de notificações |

### 7.2 Eventos Disparadores

| Evento | Tipo | Origem |
|--------|------|--------|
| Novo lead Evolution | `NEW_LEAD` | EvolutionLeadIngressService |
| Lead duplicado | `DUPLICATE_LEAD` | EvolutionLeadIngressService |

### 7.3 Fluxo

```
1. Evento dispara notifyNewLead/notifyDuplicateLead
2. Service identifica recipientes:
   - LeadWhatsappNotificationConfig ativa
   - Filtros por unidade/grupo aplicados
3. Para cada recipiente:
   3.1 Cria LeadWhatsappNotificationLog (visualized = N)
   3.2 Vincula lead_id e duplicate_id (se aplicável)
4. Badge de notificação no menu mostra countPendingForUser
5. Visualização marca log como visualized = S
```

---

## 8. Motor de Regras Granulares (UI Layer)

### 8.1 Conceito

O motor de regras granulares (`LeadVisitaGranularActionsTrait`, `ConditionalRuleTrait`) é um sistema de comportamento UI que aplica regras como:
- visibilidade de campos
- bloqueio/readonly
- definir valor sugerido
- limpar campo

**Importante**: Essas regras operam apenas na camada de UI. **Não substituem** o serviço central de sincronização (`LeadVisitStatusSyncService`).

### 8.2 Tipos de Ação

| Ação | Quando aplicar | Bloqueada no Bootstrap? |
|------|----------------|------------------------|
| `visibilidade` | Sempre | Não |
| `bloquear` | Sempre | Não |
| `readonly` | Sempre | Não |
| `definir_valor` | Apenas em change/blur real | **Sim** |
| `limpar` | Apenas em change/blur real | **Sim** |

### 8.3 Bloqueio de Mutação no Carregamento

Para evitar regressão de status ao abrir tela, ações que alteram dados (`definir_valor`, `limpar`) são **bloqueadas no bootstrap**. Apenas eventos reais do usuário (change/blur) podem disparar essas mutações.

### 8.4 Mobile vs Desktop

`LeadVisitaMobileTrait` replica o comportamento granular para o fluxo mobile. As regras devem se comportar igualmente em ambos.

---

## 9. Decisões Arquiteturais (ADRs)

### ADR-001: Categorização Funcional de Status (L/V)

**Contexto**: Clientes precisam customizar nomes/quantidade de status sem quebrar relatórios.

**Decisão**: Separar status por tipo funcional via `lead_status_categoria.status_tipo` (`L` ou `V`).

**Consequências**:
- Relatórios usam categorias, não IDs hardcoded
- Múltiplos status por categoria são suportados
- Complexidade adicional no espelhamento L/V

### ADR-002: Espelhamento Bidirecional Lead × Visita

**Contexto**: Status precisa estar sincronizado, mas nem toda visita deve controlar o lead.

**Decisão**: Sincronização só ocorre quando visita tem categoria AGENDADO ou VENDIDO.

**Consequências**:
- Lead sem visita válida permanece independente
- Modo visual (combo vs badge) reflete autoridade real
- Requer serviço central (`LeadVisitStatusSyncService`)

### ADR-003: Duplicidade em Janelas Temporais Distintas

**Contexto**: Múltiplos canais podem criar o mesmo lead simultaneamente.

**Decisão**: Duas janelas — técnica (20min/hash) e de negócio (6h/telefone).

**Consequências**:
- Proteção contra floods técnicos
- Permite recontato após período razoável
- Requer lock distribuído no ingresso

### ADR-004: Arquitetura de Ingresso via Webhook Assíncrono

**Contexto**: Volume alto de mensagens WhatsApp requer processamento confiável.

**Decisão**: Webhook → grava EvolutionWebhookEvent → processamento → criação lead.

**Consequências**:
- Resiliência contra falhas
- Audit trail completo
- Possibilidade de replay
- Latência de até alguns segundos

### ADR-005: LeadStatusHelper como API Canônica de Categorias

**Contexto**: Queries com IDs hardcoded quebravam ao customizar status.

**Decisão**: Toda query/lógica de categoria usa `LeadStatusHelper::getXxxIds()` ou `getXxxSQL()`.

**Consequências**:
- Cache em memória por requisição
- `clearCache()` obrigatório após CRUD de status
- Robustez contra customização

### ADR-006: Recuperação de Status via system_change_log

**Contexto**: Excluir última visita precisa restaurar último L conhecido.

**Decisão**: Buscar logs com `class_name in ['LeadVisitaDetalhes', 'VisitaAtribuirForm']` para identificar transição L→V e recuperar L anterior.

**Consequências**:
- Não requer coluna histórica no Lead
- Audit log já existe (SystemChangeLogTrait)
- Falha graciosa se log não existir (fallback ao default L)

### ADR-007: Visão Kanban Sem Alterar LeadsList

**Contexto**: Adicionar Kanban sem regressão na listagem padrão.

**Decisão**: Nova página (`LeadsKanbanView`) + serviço de filtros compartilhado (`LeadVisionFilterService`) + injeção de ícone via `theme.js`.

**Consequências**:
- Zero alteração em `LeadsList.php`
- Reuso 100% dos filtros via sessão
- Padrão replicável para futuras visões

### ADR-008: Motor Granular ≠ Serviço de Domínio

**Contexto**: Regras de UI tentavam fazer persistência cross-table.

**Decisão**: Motor granular faz apenas UI; persistência sempre via service central.

**Consequências**:
- Regras de UI são puramente declarativas
- Sincronização L/V centralizada em `LeadVisitStatusSyncService`
- Bloqueio de `definir_valor`/`limpar` no bootstrap evita regressão

### ADR-009: Cache de Lookups com TTL

**Contexto**: Combos auxiliares (mídia, meio, motivos) consultados a cada renderização.

**Decisão**: `LookupCacheService` com TTL de 900s (15min).

**Consequências**:
- Performance significativamente melhor
- Inconsistência tolerável (15min) em catálogos
- Cache invalidado após CRUD nos cadastros auxiliares

---

## 10. Referência à Implementação Atual

### 10.1 Controllers

| Controller | Responsabilidade | Localização |
|------------|------------------|-------------|
| `LeadsForm` | CRUD principal, sincronização status, 3 abas + subcontrollers | `app/control/leads/LeadsForm.php` |
| `LeadsList` | Listagem com filtros, busca rápida, exportação | `app/control/leads/LeadsList.php` |
| `LeadsKanbanView` | Visão Kanban (funil) | `app/control/leads/LeadsKanbanView.php` |
| `LeadsDistribuir` | Edição em lote (distribuição) | `app/control/leads/LeadsDistribuir.php` |
| `LeadVisitaDetalhes` | Subcontroller de visitas vinculadas | `app/control/leads/LeadVisitaDetalhes.php` |
| `LeadVisitaMobileActions` | Ações mobile específicas | `app/control/leads/LeadVisitaMobileActions.php` |
| `LeadAuditLogPanel` | Painel de auditoria do lead | `app/control/leads/LeadAuditLogPanel.php` |
| `LeadDuplicadoList` | Listagem de duplicados | `app/control/leads/LeadDuplicadoList.php` |
| `LeadVisitNextDayDashboard` | Dashboard de visitas D+1 | `app/control/leads/LeadVisitNextDayDashboard.php` |
| `LeadListResultSearch` | Busca específica de resultados | `app/control/leads/LeadListResultSearch.php` |
| `DashboardEvolucao` | Dashboard analítico de lead time | `app/control/leads/DashboardEvolucao.php` |
| `ThreeCPlusSyncPeriodForm` | Sincronização periódica 3C+ | `app/control/leads/ThreeCPlusSyncPeriodForm.php` |

**Cadastros auxiliares**:
- `LeadStatusForm`/`LeadStatusList` — `app/control/leads/status/`
- `LeadCategStatusForm`/`LeadCategStatusList` — `app/control/leads/status/categorias/`
- `LeadMeioForm`/`LeadMeioList` — `app/control/leads/meios/`
- `MidiaForm`/`MidiaList`/`MidiaTipoForm`/`MidiaTipoList` — `app/control/leads/midias/`
- `LeadMotivoPendForm`/`LeadMotivoPendList` — `app/control/leads/motivos/`
- `LeadMotivoPerdForm`/`LeadMotivoPerdList` — `app/control/leads/motivos/`
- `EquipeForm`/`EquipeList` — `app/control/leads/equipes/`

### 10.2 Models

| Model | Tabela | Responsabilidade |
|-------|--------|------------------|
| `Lead` | `leads` | Entidade principal, relacionamentos com SystemChangeLogTrait |
| `LeadStatus` | `lead_status` | Catálogo de status com vínculo a categoria |
| `LeadStatusCategoria` | `lead_status_categoria` | Categorias L/V |
| `LeadMeio` | `lead_meio` | Canais técnicos |
| `LeadMotivoPend` | `lead_motivo_pend` | Motivos de pendência |
| `LeadMotivoPerd` | `lead_motivo_perd` | Motivos de perda |
| `LeadEventTent` | `lead_event_tent` | Histórico de tentativas |
| `LeadHistorico` | `lead_historico` | Linha do tempo agregada |
| `LeadDuplicado` | `lead_duplicados` | Registro de duplicados |
| `LeadDistributionStats` | view/agregação | Estatísticas para Distribuir |
| `LeadWhatsappNotificationConfig` | `lead_whatsapp_notification_config` | Config de notificação |
| `LeadWhatsappNotificationLog` | `lead_whatsapp_notification_log` | Auditoria de notificações |

### 10.3 Services

| Service | Responsabilidade |
|---------|------------------|
| `LeadVisitStatusSyncService` | Sincronização L/V centralizada (canônica) |
| `EvolutionLeadIngressService` | Processamento de webhooks WhatsApp |
| `LeadWhatsappNotificationService` | Notificações de novos leads |
| `ThreeCPlusLeadSyncService` | Integração 3C+ |
| `GotoLeadSyncService` | Integração GoTo |
| `LeadVisitNextDayDashboardService` | Agregações para dashboard D+1 |

### 10.4 Libs e Helpers

| Lib | Responsabilidade |
|-----|------------------|
| `LeadStatusHelper` | API canônica de categorias (`getContactoIds()`, `getVisitasValidasSQL()`, etc.) |
| `LeadVisionFilterService` | Filtros compartilhados entre visões (List/Kanban) |
| `LeadNavigationHelper` | Navegação protegida (volta de contexto) |
| `LeadVisitaGranularActionsTrait` | Motor de regras granulares (UI) |
| `LeadVisitaFieldListTrait` | Lifecycle do field list de visitas |
| `LeadVisitaMobileTrait` | Comportamentos mobile específicos |
| `SMIDLeadProfilePhotoTrait` | Lifecycle da foto de perfil WhatsApp |
| `LookupCacheService` | Cache com TTL de catálogos |

### 10.5 Mapeamento de Fluxos

| Fluxo | Método Principal | Localização |
|-------|------------------|-------------|
| Cadastro manual | `onSave()` | `LeadsForm.php:~1300` |
| Cadastro Evolution | `processMessageUpsert()` | `EvolutionLeadIngressService` |
| Verificação duplicidade | `onCheckPhone()` | `LeadsForm.php:~700` |
| Alteração tentativas | `onChangeTentativas()` | `LeadsForm.php:~850` |
| Consulta CEP | `onLookupCep()` | `LeadsForm.php:~900` |
| Sincronização status (load) | `syncLeadStatusOnLoad()` | `LeadsForm.php` |
| Resolução save context | `resolveLeadSaveContext()` | `LeadVisitStatusSyncService` |
| Distribuição em lote | `onSaveCollection()` | `LeadsDistribuir.php` |
| Visão Kanban | `onReload()` | `LeadsKanbanView.php` |
| Auditoria | `onLoad()` | `LeadAuditLogPanel.php` |
| Notificação | `notifyNewLead()` | `LeadWhatsappNotificationService` |

### 10.6 Scripts SQL

Diretório: `app/database/scripts/3 - leads/` e `11 - categorizacao de leads/`

Principais scripts:
- `04_add_status_tipo_lead_status_categoria.sql` — adiciona tipo L/V
- `15_sync_status_first_migration_safe.sql` — saneamento seguro de status
- `08_legacy_adjust_visao_repre_e_stts_lead.sql` — ajustes legados
- `09_reconcile_status_from_legacy.sql` — reconciliação

---

## 11. Considerações de Segurança

### 11.1 Permissões por Grupo

| Grupo | ID | Permissões no Domínio Leads |
|-------|----|-----------------------------|
| **Admin** | 1 | Todas as operações |
| **Atendente** | 2 | Criar, editar leads próprios; ver lista |
| **Representante** | 4 | Ver leads de visitas atribuídas |
| **Supervisão** | 5 | Distribuir leads; alterar atendente; ver todos |

### 11.2 Validações de Acesso

- Lead só pode ser editado por atendente atribuído ou supervisor
- Status readonly quando visita válida existe (independente de permissão)
- Distribuição em lote restrita a supervisores (grupo 5)
- LeadsForm registra alterações via `SystemChangeLogTrait`
- Sessão protegida com guard token (`CURRENT_ID_TOKEN_SESSION_KEY`)

### 11.3 Campos Sensíveis

- Telefone (PII) — exibido com máscara, armazenado sem máscara
- Foto de perfil WhatsApp — URL externa armazenada, exibida com fallback
- Patologia/Observações — informação sensível de saúde, limites de tamanho

---

## 12. Glossário Técnico

| Termo Técnico | Significado |
|---------------|-------------|
| **TCombo** | Componente Adianti de seleção dropdown |
| **TEntry** | Campo de texto simples |
| **TText** | Área de texto multi-linha |
| **TMultiCombo** | Combo com seleção múltipla |
| **TDateRangeField** | Campo customizado SMID de período com presets |
| **TSidePanel / Cortina** | Painel lateral deslizante |
| **Badge** | Indicador visual colorido com texto |
| **Subcontroller** | Controller aninhado (visitas dentro do lead) |
| **SMIDStandardForm** | Classe base para formulários SMID padrão |
| **SMIDPage** | Classe base para páginas SMID |
| **TRecord** | Active Record base do Adianti |
| **TRepository** | Repositório de queries do Adianti |
| **TCriteria/TFilter** | Filtros do Adianti |
| **TTransaction** | Gestão de transação por banco |
| **SystemChangeLogTrait** | Auditoria automática em models |
| **Modo L** | Aba Lead em modo editável (combo de status) |
| **Modo V** | Aba Lead em modo readonly (badge sincronizado) |

---

**Versão**: 2.0  
**Data**: 2026-05-11  
**Autor**: SMID Architecture Team  
**Status**: Draft consolidado para revisão  
**Histórico**:
- v1.0 (2026-05-11): Estrutura inicial com 10 seções
- v2.0 (2026-05-11): Unificado — adicionados 8 controllers, 3 models, 6 libs, 5 ADRs, seções Visualizações Alternativas, Cadastros Auxiliares, Notificações e Motor Granular
