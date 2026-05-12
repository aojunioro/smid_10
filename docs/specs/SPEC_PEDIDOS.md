# SPEC_PEDIDOS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Pedidos** representa a venda concretizada — o resultado positivo do funil comercial. Um pedido é gerado a partir de um lead que passou por uma visita e converteu em compra. O sistema gerencia: cadastro, itens/produtos, condições e formas de pagamento, status do pedido, integração com financeiro, entrega, gestão por canal, distribuição comissional e auditoria.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **PedidosRepre** | Pedido principal (venda) | 0..N por lead |
| **PedProdItem** | Itens do pedido (produtos vendidos) | 0..N por pedido |
| **PedStatus** | Catálogo de status do pedido | catálogo |
| **PedStatusCategoria** | Categorias funcionais (NOVO, APROVADO, ENTREGUE, etc.) | catálogo |
| **PedFormaPagamento** | Formas de pagamento (dinheiro, cartão, pix) | catálogo |
| **PedCondicaoPagamento** | Condições de pagamento (à vista, 12x, etc.) | catálogo |
| **PedCanal** | Canais de venda (próprio, parceiro, etc.) | catálogo |

### 1.3 Relacionamentos

```
Lead 1 ──── 0..N PedidosRepre
PedidosRepre 1 ──── 0..N PedProdItem
PedidosRepre N ──── 1   PedStatus → PedStatusCategoria
PedidosRepre N ──── 1   PedFormaPagamento
PedidosRepre N ──── 1   PedCondicaoPagamento
PedidosRepre N ──── 1   PedCanal
PedidosRepre N ──── 1   SystemUser (login_repre)
PedidosRepre N ──── 1   SystemUser (login - atendente que registrou)
PedidosRepre 0..1 ──── Visita (visita que gerou a venda, via lead+data)
```

### 1.4 Integrações com outros Módulos

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| LEADS | Saída | Lead vinculado via `lead_id`; pedido confirma conversão |
| VISITAS | Saída | Pedido pode disparar status VENDIDO na última visita |
| FINANCEIRO | Bidirecional | Lançamentos automáticos; comissões; conciliação |
| TELEVENDAS | Entrada | Orçamento pré-preenchido (PedidosTeleOrcamentoPreloadHelper) |
| RELATÓRIOS | Leitura | Dashboard de evolução de vendas (T3: dt_ped) |
| AUDITORIA | Saída | system_change_log via SystemChangeLogTrait |

---

## 2. Glossário de Negócio

### 2.1 Termos do Domínio

| Termo | Definição |
|-------|-----------|
| **Pedido** | Venda registrada no sistema |
| **n_ped** | Número do pedido (identificador comercial) |
| **dt_ped** | Data comercial do pedido (autoridade para relatórios) |
| **dt_prev** | Data prevista de entrega |
| **dt_quit** | Data de quitação financeira |
| **dt_entr** | Data efetiva de entrega |
| **total_ped** | Valor bruto do pedido |
| **valor_liquido** | Valor líquido após taxas |
| **taxa_financeira** | Taxa de financiamento aplicada |
| **entrada_ped** | Valor de entrada paga |
| **obs_ped** | Observação operacional |
| **obs_ped_ger** | Observação de gerência (sigilosa) |
| **img_ped** | Imagem/foto do pedido (recibo, contrato) |

### 2.2 Categorias de Status (Matriz Funcional)

| Código | Nome | Descrição |
|--------|------|-----------|
| `NOVO` | Novo Pedido | Pedido recém-criado |
| `PROCESSANDO` | Em Processamento | Análise ou processamento interno |
| `APROVADO` | Aprovado | Aprovado para produção/entrega |
| `ENTREGA` | Em Entrega | Em transporte |
| `ENTREGUE` | Entregue | Cliente recebeu |
| `AGUARDANDO` | Aguardando | Aguardando ação externa (pagamento, doc) |
| `CANCELADO` | Cancelado | Pedido cancelado |

### 2.3 Roles e Comportamentos

| Role | Acesso |
|------|--------|
| **Admin** | Total, sem restrições |
| **Gestor** | Total operacional, vê obs_ped_ger |
| **Atendente / Recepção** | Cria/edita pedidos próprios, sem obs_ped_ger |
| **Representante** | Cria/edita pedidos próprios (apenas até status NOVO/PROCESSANDO) |
| **Financeiro** | Visão dedicada (PedidosGridFinan) |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Pedido (Manual)

**Ator**: Atendente ou Representante
**Pré-condições**: Lead existente; visita normalmente realizada

```
1. Atendente abre lead via LeadsForm
2. Aba Pedidos exibe PedidosController (subcontroller)
3. Clica "Adicionar Pedido"
4. PedidosForm em cortina:
   4.1 lead_id pré-preenchido
   4.2 login_repre da última visita pré-preenchido
   4.3 login do usuário logado pré-preenchido
5. Preenche:
   - n_ped (número do pedido)
   - dt_ped (data comercial, obrigatória)
   - status_id (status inicial: NOVO ou PROCESSANDO)
   - canal_id, fpgto_id, cpgto_id
   - total_ped, entrada_ped (cálculos automáticos)
   - dt_prev, dt_entr, dt_quit (conforme fluxo)
   - obs_ped (operacional)
   - obs_ped_ger (apenas gerência)
   - img_ped (foto opcional)
6. Adiciona itens (PedProdItem):
   - produto, quantidade, valor unitário, descontos
7. Salva
8. SystemChangeLogTrait registra criação
9. Sincroniza visita.stts_lead = VENDIDO se aplicável
```

**Regras de Negócio**:
- RN-001: lead_id e dt_ped são obrigatórios
- RN-002: Status inicial deve ser categoria NOVO ou PROCESSANDO
- RN-003: total_ped >= entrada_ped
- RN-004: obs_ped_ger só visível/editável para gestor+

### 3.2 Cadastro via Televendas (Pré-preenchimento)

**Ator**: Operador de Televendas
**Pré-condições**: Atendimento em curso

```
1. Sistema invoca PedidosTeleOrcamentoPreloadHelper
2. Helper extrai dados do contexto:
   - Lead atual
   - Última visita (login_repre, dt_visita)
   - Histórico de pedidos do lead (template)
3. PedidosForm é aberto com campos pré-preenchidos
4. Operador ajusta e salva
```

### 3.3 Edição de Pedido

```
1. Atendente abre pedido via PedidosList ou Lead
2. PedidosForm em cortina
3. Restrições por status:
   - NOVO/PROCESSANDO: edição livre (atendente)
   - APROVADO+: edição restrita ao gestor
   - ENTREGUE/CANCELADO: somente leitura (exceto admin)
4. Atualiza status conforme fluxo:
   - dt_quit preenchida → status APROVADO
   - dt_entr preenchida → status ENTREGUE
5. Salva e registra alteração no audit log
```

### 3.4 Detalhamento de Produtos (PedProdDetalhes)

```
1. Atendente clica em pedido na lista
2. PedProdDetalhes exibe itens em cortina/modal
3. Edição inline de itens:
   - Adicionar PedProdItem
   - Alterar quantidade/valor
   - Remover item
4. Total recalculado automaticamente
5. Salva atualização do pedido + itens
```

### 3.5 Visualizações Alternativas

#### 3.5.1 PedidosList (Listagem Padrão)

```
- Filtros: período (dt_ped, dt_prev, dt_quit, dt_entr, dt_visita), status, mídia,
  unidade, atendente, representante, valores (vlr_min/max), forma/condição
- Busca rápida (n_ped, cpf, nome)
- Paginação server-side
- Exportação Excel
```

#### 3.5.2 PedidosKanbanView

```
- Visão em funil por status
- Reusa filtros via PedidoVisionFilterService
- Drag-and-drop entre status
- Limite de cards configurável
```

#### 3.5.3 PedidosGrid (Embutido)

Grid embutido em LeadsForm, lista pedidos do lead.

#### 3.5.4 PedidosGridFinan

Visão financeira dedicada com:
- Status de pagamento
- Pendências
- Comissões
- Conciliação

#### 3.5.5 PedidoAuditLogPanel

Painel de auditoria via system_change_log:
- Filtra por target_table='pedidos' e record_id
- Resolve IDs para nomes (status, fpgto, cpgto, canal)
- Mostra origem (controller) e usuário

#### 3.5.6 PedidosMobileActions

Controller dedicado para ações mobile (similar a LeadVisitaMobileActions).

### 3.6 Lançamento Financeiro

**Referência**: `docs/financeiro/lancamentoPedidos/LANCAMENTO_PEDIDOS.md`

```
1. Pedido aprovado dispara lançamento financeiro
2. Sistema cria FinLancamento vinculado ao pedido
3. Calcula valor_liquido = total_ped - taxa_financeira
4. Comissão calculada conforme regras de negócio
5. Conciliação manual ou automática
```

---

## 4. Estados e Transições

### 4.1 Máquina de Estados

```
┌─────────┐    ┌──────────────┐    ┌──────────┐    ┌─────────┐    ┌──────────┐
│  NOVO   │───▶│ PROCESSANDO  │───▶│ APROVADO │───▶│ ENTREGA │───▶│ ENTREGUE │
└─────────┘    └──────────────┘    └──────────┘    └─────────┘    └──────────┘
     │                │                  │              │              │
     ▼                ▼                  ▼              ▼              │
┌─────────┐    ┌──────────────┐    ┌──────────────────────────┐       │
│CANCELADO│    │  AGUARDANDO  │    │       CANCELADO          │       │
└─────────┘    └──────────────┘    └──────────────────────────┘       │
                      │                                                │
                      └────────retoma processamento───────────────────┘
```

### 4.2 Transições Permitidas

| De | Para | Gatilho |
|----|------|---------|
| NOVO | PROCESSANDO | Início análise |
| NOVO/PROCESSANDO | APROVADO | dt_quit preenchida ou aprovação manual |
| NOVO/PROCESSANDO | AGUARDANDO | Falta documento/pagamento |
| AGUARDANDO | PROCESSANDO | Pendência resolvida |
| APROVADO | ENTREGA | Início do transporte |
| ENTREGA | ENTREGUE | dt_entr preenchida |
| QUALQUER | CANCELADO | Cancelamento manual (gestor+) |

### 4.3 Bloqueios por Status

| Status | Bloqueio |
|--------|----------|
| ENTREGUE | Edição (exceto admin) |
| CANCELADO | Edição (exceto admin) |
| APROVADO+ | Edição restrita ao gestor |

---

## 5. Cadastros Auxiliares

### 5.1 Status de Pedidos

| Controller | Função |
|------------|--------|
| `PedStatusForm` / `PedStatusList` | CRUD de status com vínculo a categoria |
| `PedCategStatusForm` / `PedCategStatusList` | CRUD de categorias funcionais |

**Regras**:
- Status protegidos (sistema='S') não podem ser deletados
- `PedStatusHelper` fornece API canônica (similar a LeadStatusHelper)

### 5.2 Formas de Pagamento

| Controller | Função |
|------------|--------|
| `PedFormaPagamentoForm` / `PedFormaPagamentoList` | CRUD (dinheiro, cartão, pix, etc.) |

### 5.3 Condições de Pagamento

| Controller | Função |
|------------|--------|
| `PedCondicaoPagamentoForm` / `PedCondicaoPagamentoList` | CRUD (à vista, 12x, etc.) |

### 5.4 Canais de Venda

| Controller | Função |
|------------|--------|
| `PedCanalForm` / `PedCanalList` | CRUD (próprio, parceiro, etc.) |

---

## 6. Decisões Arquiteturais (ADRs)

### ADR-001: Categorização Funcional de Status

**Decisão**: Mesmo padrão do domínio Leads — `ped_status_categoria` separa status por categoria funcional.

**Consequências**:
- Relatórios usam `PedStatusHelper::getEntregaSQL()` em vez de IDs hardcoded
- Cliente customiza nomes sem quebrar relatórios
- Múltiplos status por categoria suportados

### ADR-002: dt_ped como Autoridade Comercial

**Contexto**: Pedidos podem ser cadastrados depois da venda real.

**Decisão**: `dt_ped` é a data comercial usada em todos os relatórios. `criado_em` é apenas auditoria.

**Consequências**:
- Dashboard de Evolução usa dt_ped (T3)
- Filtros principais por dt_ped
- Atendente pode editar dt_ped (com permissão)

### ADR-003: PedidosListModel para Queries Complexas

**Contexto**: Listagem de pedidos requer múltiplos JOINs e filtros.

**Decisão**: Model dedicado `PedidosListModel` com método estático `getPedidos()`.

**Consequências**:
- Queries SQL otimizadas
- Reuso entre PedidosList, PedidosKanbanView, PedidosGridFinan
- Fora do padrão TRecord, mas justificado por performance

### ADR-004: Pré-preenchimento via Helpers

**Decisão**: Helpers dedicados:
- `PedidosTeleOrcamentoPreloadHelper` — contexto de televendas
- `PedidosTeleNavigationHelper` — navegação contextual
- `PedidosRepreNavigationHelper` — navegação por representante

**Consequências**:
- Lógica de pré-preenchimento isolada
- Fluxos contextuais sem duplicação
- Padrão replicável

### ADR-005: Visão Financeira Dedicada (PedidosGridFinan)

**Decisão**: Grid separado focado em métricas financeiras.

**Consequências**:
- Permissões mais restritivas
- Foco em pagamento, comissão, conciliação
- Sem mistura com operação comercial

### ADR-006: img_ped como BLOB

**Decisão**: Imagem do pedido armazenada como BLOB no banco.

**Consequências**:
- Backup integrado
- Sem dependência de filesystem externo
- Cuidado com tamanho do banco (compressão recomendada)

### ADR-007: obs_ped vs obs_ped_ger

**Decisão**: Duas observações distintas — operacional pública e gerencial restrita.

**Consequências**:
- Atendente vê apenas obs_ped
- Gestor vê ambas
- Audit log separa alterações

### ADR-008: Auditoria via PedidoAuditLogPanel

**Decisão**: Painel dedicado consumindo system_change_log.

**Consequências**:
- Compliance e rastreabilidade
- Lookups resolvem IDs para nomes
- Origem (controller) registrada

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `PedidosForm` | `app/control/pedidos/PedidosForm.php` |
| `PedidosList` | `app/control/pedidos/PedidosList.php` |
| `PedidosGrid` | `app/control/pedidos/PedidosGrid.php` |
| `PedidosGridFinan` | `app/control/pedidos/PedidosGridFinan.php` |
| `PedidosKanbanView` | `app/control/pedidos/PedidosKanbanView.php` |
| `PedidosMobileActions` | `app/control/pedidos/PedidosMobileActions.php` |
| `PedProdDetalhes` | `app/control/pedidos/PedProdDetalhes.php` |
| `PedidoAuditLogPanel` | `app/control/pedidos/PedidoAuditLogPanel.php` |

**Cadastros auxiliares**:
- `PedStatusForm/List` — `app/control/pedidos/status/`
- `PedCategStatusForm/List` — `app/control/pedidos/status/categorias/`
- `PedCanalForm/List` — `app/control/pedidos/canal/`
- `PedCondicaoPagamentoForm/List` — `app/control/pedidos/cpgto/`
- `PedFormaPagamentoForm/List` — `app/control/pedidos/fpgto/`

### 7.2 Models

| Model | Tabela | Auditoria |
|-------|--------|-----------|
| `PedidosRepre` | `pedidos` | SystemChangeLogTrait |
| `PedProdItem` | `ped_prod_item` | — |
| `PedStatus` | `ped_status` | — |
| `PedStatusCategoria` | `ped_status_categoria` | — |
| `PedFormaPagamento` | `ped_fpgto` | — |
| `PedCondicaoPagamento` | `ped_cpgto` | — |
| `PedCanal` | `ped_canal` | — |
| `PedidosListModel` | (queries) | — |

### 7.3 Libs e Helpers

| Lib | Responsabilidade |
|-----|------------------|
| `PedStatusHelper` | API canônica de categorias |
| `PedidoVisionFilterService` | Filtros compartilhados (List/Kanban) |
| `PedidosRepreNavigationHelper` | Navegação por representante |
| `PedidosTeleNavigationHelper` | Navegação televendas |
| `PedidosTeleOrcamentoPreloadHelper` | Pré-preenchimento televendas |

### 7.4 Mapeamento de Fluxos

| Fluxo | Método Principal |
|-------|------------------|
| Cadastro | `PedidosForm::onSave()` |
| Listagem | `PedidosListModel::getPedidos()` |
| Itens | `PedProdDetalhes::onSave()` |
| Kanban | `PedidosKanbanView::onReload()` |
| Auditoria | `PedidoAuditLogPanel::onLoad()` |
| Pré-preenchimento Tele | `PedidosTeleOrcamentoPreloadHelper::preload()` |

### 7.5 Scripts SQL

Diretório: `app/database/scripts/12 - categorizacao de pedidos/`

Principais:
- Criação de `ped_status_categoria`
- Migração de status para categorias
- Criação de `v_pedidos_validos`
- Criação de `v_evolucao_vendas` (lead time)

---

## 8. Considerações de Segurança

### 8.1 Permissões

| Role | Cria | Edita | Exclui | Vê obs_ped_ger |
|------|------|-------|--------|----------------|
| Admin | Sempre | Sempre | Sempre | Sim |
| Gestor | Sempre | Sempre | Conforme | Sim |
| Atendente | Sim | Status iniciais | Recém-criados | Não |
| Representante | Sim | Próprios em NOVO | Não | Não |
| Financeiro | Não | Apenas campos finan. | Não | Não |

### 8.2 Auditoria

- `PedidosRepre` usa SystemChangeLogTrait (registro automático)
- PedidoAuditLogPanel exibe histórico
- Alterações em status, valores e datas são registradas

### 8.3 Campos Sensíveis

- **CPF**: PII, exibido com permissão
- **obs_ped_ger**: gerencial restrito
- **img_ped**: pode conter dados sensíveis (contrato), acesso restrito
- **valor_liquido / taxa_financeira**: visibilidade controlada

---

## 9. Métricas e Dashboards

### 9.1 Dashboard de Evolução de Vendas

**T3 — Concretização da Venda**: usa `dt_ped` do último pedido válido do lead.

**View**: `v_pedidos_validos` filtra pedidos válidos para análise (exclui cancelados, etc.).

### 9.2 PedidosGridFinan

Métricas financeiras agregadas:
- Total bruto / líquido por período
- Pendências de pagamento
- Comissões devidas
- Taxa de cancelamento

### 9.3 Indicadores por Status

- Tempo médio em cada status
- Conversão NOVO → ENTREGUE
- Análise por canal/representante

---

## 10. Glossário Técnico

| Termo | Significado |
|-------|-------------|
| **PedidosListModel** | Model não-TRecord com queries customizadas |
| **PedStatusHelper** | API de categorias (similar a LeadStatusHelper) |
| **v_pedidos_validos** | View SQL de pedidos válidos para relatórios |
| **v_evolucao_vendas** | View de lead time (T0 a T3) |
| **img_ped** | BLOB com imagem do pedido |
| **obs_ped_ger** | Observação gerencial restrita |
| **lancamento financeiro** | Registro automático em FinLancamento |

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
