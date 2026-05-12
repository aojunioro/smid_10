# SPEC_VISITAS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Visitas** representa o agendamento e a execução de visitas técnicas/comerciais aos leads. Uma visita é o momento em que o representante (vendedor externo, terapeuta) encontra o cliente. O sistema deve gerenciar o ciclo: agendamento, confirmação, atribuição de representante, execução com check-in GPS, registro do ocorrido (histórico), pós-visita (avaliação) e conversão em pedido.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Visita** | Agendamento da visita ao cliente | 0..N por lead |
| **VisitaCheckinEvent** | Evento GPS de check-in/check-out | 0..N por visita |
| **VisitaRepre** | Vínculo histórico representante × visita | 0..N |
| **VisitaCalendarioModel** | View/agregação para visão calendário | view |
| **PosVisita** | Avaliação pós-visita (notas, percepções) | 0..1 por visita |
| **PosVisitaProvider** | Provider de dados para PosVisita | service |
| **HistoricoRepre** | Registro do ocorrido na visita | 0..N por visita |
| **AgendamentoConfigHorario** | Slots de horários disponíveis | catálogo |
| **VisitaStatus** | Catálogo de status (mesmo de LeadStatus, tipo V) | catálogo |

### 1.3 Relacionamentos

```
Lead 1 ──── 0..N Visita
Visita 1 ──── 0..N VisitaCheckinEvent (GPS)
Visita 1 ──── 0..1 PosVisita (avaliação)
Visita 1 ──── 0..N HistoricoRepre (ocorridos)
Visita 1 ──── 0..1 PedidosRepre (venda concretizada)
Visita N ──── 1   SystemUser (login_recep / atendente)
Visita N ──── 0..1 SystemUser (login_repre / representante)
Visita N ──── 1   LeadStatus (stts_lead, tipo V)
```

### 1.4 Integrações com outros Módulos

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| LEADS | Bidirecional | Espelhamento de status quando categoria V; criação de visita altera lead |
| HISTORICOS | Saída | Visita gera registro em historicos.vis_id |
| PEDIDOS | Saída | Visita pode gerar pedido (venda) |
| KM_RODADO | Leitura | Cálculo de quilometragem para reembolso |
| GPS / Geolocalização | Externo | Check-in/out via API browser |
| PERMISSOES | Leitura | Roles (admin, gestor, atendente, representante) |

---

## 2. Glossário de Negócio

### 2.1 Termos do Domínio

| Termo | Definição |
|-------|-----------|
| **Visita** | Agendamento de encontro presencial com o lead |
| **Atendente** | Recepcionista que cadastra/agenda a visita (login_recep) |
| **Representante** | Profissional que realiza a visita (login_repre) |
| **Confirmação** | Validação prévia de que cliente confirma (confirm + login_conf + dt_confirm) |
| **Check-in** | Registro GPS do início da visita |
| **Check-out** | Registro GPS do fim da visita |
| **Pós-visita** | Avaliação pós-execução (notas e percepções) |
| **Histórico** | Registro do que ocorreu durante/após a visita |
| **Interesse** | Indicador de interesse capturado durante a visita |
| **Hist Feito / Pos Feito** | Flags indicando se histórico/pós-visita foram preenchidos |

### 2.2 Status de Visita (Categoria V)

| Status | Categoria | Cor | Descrição |
|--------|-----------|-----|-----------|
| **Pré-Agendado** | AGENDADO | #007bff | Visita agendada, aguardando confirmação |
| **Agendado** | AGENDADO | #007bff | Visita confirmada |
| **Andamento** | EM_ATENDIMENTO | #ffc107 | Visita em execução (com representante) |
| **Reagendado** | AGENDADO | #007bff | Visita remarcada |
| **Vendido** | VENDIDO | #28a745 | Visita gerou venda |
| **Entregue** | VENDIDO | #28a745 | Produto entregue |

### 2.3 Roles e Comportamentos

| Role | Constante | Descrição |
|------|-----------|-----------|
| **Admin** | `ROLE_ADMIN` | Acesso total, sem restrições |
| **Gestor (Supervisor)** | `ROLE_SUPERVISOR` | Acesso elevado, atribui representantes |
| **Atendente** | `ROLE_ATTENDANT` | Cadastra/edita em status iniciais |
| **Representante** | `ROLE_REPRESENTATIVE` | Executa visita, faz check-in, registra histórico |

---

## 3. Fluxos Principais

### 3.1 Agendamento da Visita

**Ator**: Atendente
**Pré-condições**: Lead cadastrado; sem visita ativa em status bloqueante

```
1. Atendente acessa o lead via LeadsForm
2. Aba Visitas (LeadVisitaDetalhes) é exibida
3. Atendente clica "Adicionar Visita"
4. Preenche dados:
   - dt_visita (data, obrigatório se status PRÉ-AG/AG/REAG)
   - hr_visita (hora, opcional inicialmente)
   - login_repre (representante, opcional inicialmente)
   - stts_lead (status, opcional - vide RN-014 do SPEC_LEADS)
5. Salva visita inline (TFieldList) ou via modal mobile
6. VisitaValidationService valida:
   - Endereço completo se status requer (PREAGENDADO/AGENDADO/REAGENDADO)
   - Não conflito de horário do representante
7. Persiste visita; lead pode entrar em Modo V se status válido
```

**Regras de Negócio**:
- RN-001: Não permitido criar nova visita se já existe em PREAGENDADO/AGENDADO/ANDAMENTO (BLOCK_NEW_VISIT)
- RN-002: Endereço completo obrigatório para PREAGENDADO/AGENDADO/REAGENDADO (REQUIRES_ADDRESS)
- RN-003: Status pode ficar vazio (visita opcional) - automação assistiva

### 3.2 Confirmação da Visita

```
1. Atendente liga para o cliente em data próxima à visita
2. Cliente confirma → atendente marca:
   - confirm = 'S'
   - login_conf (quem confirmou)
   - dt_confirm (timestamp)
3. Sistema atualiza status (Pré-Agendado → Agendado)
```

### 3.3 Atribuição de Representante (VisitaAtribuirForm)

**Ator**: Supervisor
**Pré-condições**: Visita em AGENDADO

```
1. Supervisor acessa lista de visitas
2. Abre VisitaAtribuirForm em cortina
3. Combo de representantes restrito ao grupo 4
4. Combo stts_lead carrega apenas tipo V (com fallback legado)
5. Atribui login_repre:
   - Vazio → categoria AGENDADO
   - Preenchido → categoria EM_ATENDIMENTO (Andamento)
6. onSave valida:
   - stts_lead pertence a categoria V
7. Persiste visita
8. LeadVisitStatusSyncService::syncLeadFromLatestVisit propaga ao lead
```

**Regras de Negócio**:
- RN-004: Apenas grupo 4 (Representantes) aparece no combo
- RN-005: stts_lead deve ser tipo V (validação backend)
- RN-006: Sem IDs hardcoded - usa categorias

### 3.4 Execução com GPS (Check-in/Check-out)

**Ator**: Representante
**Pré-condições**: Visita em ANDAMENTO; permissão GPS concedida

```
1. Representante abre HistoricoRepreForm (mobile)
2. GpsPermissionGate solicita permissão de geolocalização
3. SE permissão negada:
   3.1 Registra VisitaCheckinEvent com permission='N' e erro_codigo
   3.2 Bloqueia ações que dependem de GPS
4. SE permissão concedida:
   4.1 Captura lat/lng/accuracy_m
   4.2 Cria VisitaCheckinEvent (tipo=CHECK_IN, capturado_em)
   4.3 Permite registrar histórico
5. Ao fim da visita:
   5.1 Captura check-out GPS
   5.2 Cria VisitaCheckinEvent (tipo=CHECK_OUT)
6. KmCalculationService calcula distância para reembolso
```

**Regras de Negócio**:
- RN-007: Toda visita executada deve ter ao menos 1 check-in GPS
- RN-008: Eventos GPS são imutáveis (audit trail)
- RN-009: Falha de GPS é registrada com erro_codigo + erro_msg
- RN-010: Accuracy_m alta (>50m) gera flag de validação

### 3.5 Registro de Histórico (Ocorrido)

```
1. Representante (ou Atendente) registra HistoricoRepre:
   - vis_id (visita)
   - lead_id (preenchido a partir de visita.lead_id)
   - login (quem registrou)
   - motivo_id (HistoricoMotivo)
   - ocorr_id (HistoricoOcorrido)
   - hist (texto livre)
   - foto_hist (foto opcional)
2. Sistema atualiza visita.hist_feito = 'S'
3. Se ocorrido invalida visita (hist_motivo=2): visita perde validade
```

### 3.6 Pós-Visita (Avaliação)

```
1. Atendente liga para cliente após visita
2. Acessa PosVisitaForm em cortina
3. Preenche:
   - nota_repre (nota do representante)
   - nota_prod (nota do produto)
   - nota_empre (nota da empresa)
   - visitado (S/N - confirmou que visita ocorreu)
   - pontual (S/N)
   - jaleco (S/N - usou jaleco/uniforme)
   - adquiriu (S/N - comprou)
   - obs (observações)
4. Sistema atualiza visita.pos_feito = 'S'
5. PosVisitaList exibe agregações para gestão de qualidade
```

### 3.7 Calendário de Visitas

**VisitaCalendario** (desktop) e **VisitaCalendarioMobile**:

```
1. Visualização em calendário das visitas agendadas
2. Filtros por período, representante, atendente, unidade
3. Cores por status (LeadStatusHelper.getCorByCategoria)
4. Click no evento abre LeadsForm em modo edição
5. Mobile: visão em lista compacta com cards
```

### 3.8 Disponibilidade de Horários (Agendamento)

**AgendamentoConfigHorariosForm**: configura slots disponíveis

**AgendamentoDisponibilidadeView**: visualiza disponibilidade

```
1. Gestor configura horários disponíveis por unidade/representante
2. Sistema considera ao sugerir horários no agendamento
3. Bloqueia conflitos de agenda
```

---

## 4. Estados e Transições

### 4.1 Máquina de Estados

```
        ┌────────────┐
        │ Pré-Agend. │──┐
        └────────────┘  │ confirmação
              │         ▼
              │   ┌────────────┐
              │   │  Agendado  │
              │   └────────────┘
              │         │ atribui representante
              │         ▼
              │   ┌────────────┐
              │   │ Andamento  │
              │   └────────────┘
              │      │      │
   reagenda  │      │      │ vendeu
              ▼      ▼      ▼
        ┌──────────┐  ┌──────────┐
        │Reagendado│  │ Vendido  │
        └──────────┘  └──────────┘
                          │
                          ▼ entrega
                    ┌──────────┐
                    │ Entregue │
                    └──────────┘
```

### 4.2 Bloqueios por Status

| Conjunto | Status | Operação Bloqueada |
|----------|--------|--------------------|
| `BLOCK_NEW_VISIT` | PREAGENDADO, AGENDADO, ANDAMENTO | Criar nova visita para o lead |
| `BLOCK_EDIT_DELETE` | ANDAMENTO, REAGENDADO, VENDIDO, ENTREGUE | Editar/excluir visita (exceto admin/supervisor) |
| `ALLOW_ATTENDANT_EDIT` | PREAGENDADO, AGENDADO | Atendente pode editar/excluir |
| `REQUIRES_ADDRESS` | PREAGENDADO, AGENDADO, REAGENDADO | Endereço completo obrigatório |

### 4.3 Automações Assistivas (LeadVisitaGranularActionsTrait)

| Condição | Ação Sugerida |
|----------|---------------|
| dt_visita + hr_visita preenchidos | Sugere AGENDADO |
| login_repre preenchido | Sugere ANDAMENTO |
| login_repre removido | Sugere AGENDADO |
| dt_visita ou hr_visita limpos | Permite limpar status |

**Princípio**: Automações operam apenas em eventos reais (change/blur), nunca no bootstrap.

---

## 5. Visualizações

### 5.1 LeadVisitaDetalhes (Sub-controller)

**Tipo**: TFieldList inline dentro do LeadsForm
- Edição inline em tabela (desktop)
- Cards verticais (mobile, via mobile-premium.css)
- Reconstrução de linhas via `addDetail` (não sincronização frágil via JS)

### 5.2 LeadVisitaMobileActions

**Tipo**: Controller dedicado para ações mobile
- Edição via modal popup
- Encapsula transação e ciclo de vida
- Resolve "tela branca" e "crash de construtor"

### 5.3 VisitaList

**Tipo**: Listagem global de visitas
- Filtros: período, status, representante, atendente, unidade
- Busca rápida server-side
- Indicadores: hist_feito, pos_feito, confirm

### 5.4 VisitaCalendario (Desktop) / VisitaCalendarioMobile

**Tipo**: Visão em calendário
- Cores por categoria de status
- Click abre LeadsForm
- Filtros persistidos em sessão

### 5.5 PosVisitaList / PosVisitaGrid

**Tipo**: Listagem de avaliações pós-visita
- Agregações de notas
- Filtros por representante, período, unidade

---

## 6. Cadastros Auxiliares

### 6.1 Status de Visita

| Controller | Função |
|------------|--------|
| `VisitaStatusForm` / `VisitaStatusList` | Compartilhado com LeadStatus, filtrado por tipo V |

### 6.2 Configuração de Horários

| Controller | Função |
|------------|--------|
| `AgendamentoConfigHorariosForm` / `List` | Slots disponíveis |
| `AgendamentoDisponibilidadeView` | Visualização de disponibilidade |

### 6.3 Pós-Visita

| Controller | Função |
|------------|--------|
| `PosVisitaForm` | Cadastro/edição de avaliação |
| `PosVisitaList` | Listagem agregada |
| `PosVisitaGrid` | Grid embutido em outros forms |

---

## 7. Decisões Arquiteturais (ADRs)

### ADR-001: Status Compartilhado L/V no Mesmo Catálogo

**Decisão**: Não criar tabela separada para status de visita; usar `lead_status` com `lead_status_categoria.status_tipo='V'`.

**Consequências**:
- Catálogo único, evolução conjunta
- Filtro por tipo em queries
- Compatibilidade legada preservada

### ADR-002: Serviços de Domínio Separados (Status, Permission, Validation)

**Contexto**: Regras de status, permissão e validação eram espalhadas no controller.

**Decisão**: Três services dedicados:
- `VisitaStatusService` — constantes de status, normalização, classificação
- `VisitaPermissionService` — roles e permissões
- `VisitaValidationService` — validações (endereço completo, conflito, etc.)

**Consequências**:
- Reuso entre Form, List, Calendário
- Testabilidade
- Mudanças de regra centralizadas

### ADR-003: Dupla Visualização Desktop/Mobile

**Decisão**: TFieldList no desktop; cards verticais com modal no mobile.

**Consequências**:
- UX otimizada por dispositivo
- Controller dedicado (`LeadVisitaMobileActions`)
- CSS Bootstrap (`d-none d-md-block` / `d-block d-md-none`)

### ADR-004: Check-in GPS Imutável

**Decisão**: `VisitaCheckinEvent` é append-only, com `IDPOLICY=serial`.

**Consequências**:
- Audit trail completo
- Permissão registrada (S/N) com erro_codigo
- Suporte a accuracy_m para validação

### ADR-005: Status Opcional na Visita

**Contexto**: Operação de callcenter precisava deixar visita sem status.

**Decisão**: `stts_lead` é opcional, mesmo com data/hora/endereço preenchidos.

**Consequências**:
- Lead permanece em Modo L se visita não for válida
- Automações são assistivas, não coercitivas
- Validação no save permite vazio

### ADR-006: Representante Restrito ao Grupo 4

**Decisão**: Combo de representantes filtra apenas `usergroupids` contendo `4`.

**Consequências**:
- Segregação clara de papéis
- Validação consistente em todos os pontos
- Suporte a transição automática (presença → ANDAMENTO)

### ADR-007: KM Calculado por Service Dedicado

**Decisão**: `KmCalculationService` consome eventos GPS e calcula distâncias.

**Consequências**:
- Lógica de reembolso isolada
- Validação via `KmGpsValidationService`
- Configuração em `KmConfigService`

---

## 8. Referência à Implementação Atual

### 8.1 Controllers

| Controller | Localização |
|------------|-------------|
| `VisitaList` | `app/control/visitas/VisitaList.php` |
| `VisitaAtribuirForm` | `app/control/visitas/VisitaAtribuirForm.php` |
| `VisitaCalendario` | `app/control/visitas/VisitaCalendario.php` |
| `VisitaCalendarioMobile` | `app/control/visitas/VisitaCalendarioMobile.php` |
| `LeadVisitaDetalhes` | `app/control/leads/LeadVisitaDetalhes.php` |
| `LeadVisitaMobileActions` | `app/control/leads/LeadVisitaMobileActions.php` |
| `PosVisitaForm` / `PosVisitaList` / `PosVisitaGrid` | `app/control/visitas/pos-visita/` |
| `AgendamentoConfigHorariosForm` / `List` | `app/control/visitas/disponibilidade/` |
| `AgendamentoDisponibilidadeView` | `app/control/visitas/disponibilidade/` |
| `VisitaStatusForm` / `List` | `app/control/visitas/status/` |

### 8.2 Models

| Model | Tabela |
|-------|--------|
| `Visita` | `visitas` |
| `VisitaCheckinEvent` | `visita_checkin_event` |
| `VisitaRepre` | `visita_repre` |
| `VisitaCalendarioModel` | view |
| `PosVisita` | `pos_visita` |
| `PosVisitaProvider` | service |

### 8.3 Services

| Service | Responsabilidade |
|---------|------------------|
| `VisitaStatusService` | Regras de status, classificação, bloqueios |
| `VisitaPermissionService` | Roles e permissões |
| `VisitaValidationService` | Validações de endereço, conflitos |
| `LeadVisitStatusSyncService` | Espelhamento Lead × Visita (compartilhado) |
| `KmCalculationService` | Cálculo de KM |
| `KmGpsValidationService` | Validação GPS |
| `KmReembolsoDataService` | Dados de reembolso |
| `KmConfigService` | Configuração de KM |

### 8.4 Libs

| Lib | Responsabilidade |
|-----|------------------|
| `VisitaNavigationHelper` | Navegação protegida |
| `LeadVisitaGranularActionsTrait` | Motor granular UI |
| `LeadVisitaFieldListTrait` | Lifecycle field list |
| `LeadVisitaMobileTrait` | Comportamentos mobile |
| `GpsPermissionGate` | Gate de permissão GPS |

---

## 9. Considerações de Segurança

### 9.1 Permissões

| Role | Pode Criar | Pode Editar | Pode Excluir |
|------|-----------|-------------|--------------|
| Admin | Sempre | Sempre | Sempre |
| Gestor | Sempre | Sempre | Sempre |
| Atendente | Sim | Apenas em PREAG/AG | Apenas em PREAG/AG |
| Representante | Não | Apenas durante visita atribuída | Não |

### 9.2 Auditoria

- `Visita` usa `SystemChangeLogTrait` (registro automático)
- `VisitaCheckinEvent` é imutável (audit GPS)
- `system_change_log` permite recuperação de status anterior

---

## 10. Glossário Técnico

| Termo | Significado |
|-------|-------------|
| **TFieldList** | Componente de lista editável inline |
| **VisitaCheckinEvent** | Evento append-only de GPS |
| **Modo V** | Estado em que visita controla o lead |
| **Hist Feito** | Flag de histórico preenchido (boolean) |
| **Pos Feito** | Flag de pós-visita preenchida (boolean) |
| **REQUIRES_ADDRESS** | Conjunto de status que exigem endereço |
| **BLOCK_NEW_VISIT** | Conjunto que bloqueia nova visita |

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
