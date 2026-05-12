# SPEC_COMISSOES.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Comissões** controla a apuração, acompanhamento e pagamento de comissões vinculadas a pedidos. Ele permite calcular o valor devido ao representante/operador, registrar pagamentos parciais, controlar saldo e integrar o pagamento ao módulo financeiro.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Comissoes** | Registro agregado da comissão de um pedido | 0..1 por pedido |
| **ComissItem** | Pagamentos/parciais de comissão | 0..N por comissão |
| **PedidosRepre** | Pedido que origina a comissão | 1 |
| **FinContaPagar** | Conta financeira vinculada à comissão | 0..N |

### 1.3 Relacionamentos

```
Pedido 1 ──── 0..1 Comissoes
Comissoes 1 ──── 0..N ComissItem
Comissoes 1 ──── 0..N FinContaPagar
Pedido N ──── 1 SystemUser (login_repre)
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Entrada | Pedido vendido gera comissão |
| TELEVENDAS | Entrada | Pedidos Tele podem gerar comissão própria do canal |
| FINANCEIRO | Saída | Pagamentos de comissão viram contas a pagar |
| METAS | Leitura | Vendas/faturamento podem influenciar metas |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Comissão** | Valor devido por uma venda/pedido |
| **Valor Comissão** | Valor total calculado (`vlr_comissao`) |
| **Total Pago** | Soma de pagamentos já realizados (`total_pago`) |
| **Saldo** | `vlr_comissao - total_pago` |
| **Data Prevista** | Previsão de pagamento (`dt_prevista`) |
| **Status Comissão** | Situação da comissão (`stt_comis`) |
| **Item de Comissão** | Pagamento parcial ou total registrado em `comis_ped_item` |

---

## 3. Fluxos Principais

### 3.1 Geração de Comissão por Pedido

```
1. Pedido é criado ou atinge status comissionável
2. Sistema identifica representante/operador associado
3. Calcula valor da comissão conforme regra vigente
4. Cria Comissoes:
   - ped_id
   - vlr_comissao
   - total_pago = 0
   - vlr_saldo = vlr_comissao
   - dt_prevista
   - stt_comis = PENDENTE/EM_ABERTO
5. Registra auditoria via SystemChangeLogTrait
```

**Regras**:
- RN-001: Um pedido deve ter no máximo uma comissão agregada ativa
- RN-002: Comissão não deve ser gerada para pedido cancelado
- RN-003: Comissão cancelada não entra em saldo a pagar

### 3.2 Pagamento Parcial/Total

```
1. Usuário abre ComisPedDetalhes
2. Informa pagamento:
   - vlr_pago
   - dt_pgto
   - obs_pgto
3. Cria ComissItem
4. Atualiza Comissoes:
   - total_pago = soma(ComissItem.vlr_pago)
   - vlr_saldo = vlr_comissao - total_pago
   - stt_comis = PAGO se saldo <= 0; PARCIAL se saldo > 0
5. Opcionalmente cria/atualiza FinContaPagar
```

**Regras**:
- RN-004: `vlr_pago` deve ser maior que zero
- RN-005: Total pago não deve exceder valor da comissão (exceto ajuste autorizado)
- RN-006: Data de pagamento é obrigatória para item pago

### 3.3 Visão Geral de Comissões

```
1. Gestor acessa ComissoesGeralList
2. Filtros:
   - Período
   - Representante
   - Status comissão
   - Pedido
3. Listagem exibe:
   - Pedido
   - Cliente
   - Representante
   - Valor comissão
   - Total pago
   - Saldo
   - Data prevista
   - Status
4. Ação abre ComisPedDetalhes
```

### 3.4 Visão do Representante

```
1. Representante acessa ComissoesList
2. Sistema filtra por login_repre do usuário logado
3. Exibe apenas comissões próprias
4. Permite visualizar detalhes e histórico de pagamentos
```

### 3.5 Integração com Financeiro

```
1. Comissão apurada gera ou agenda FinContaPagar
2. Conta possui comissao_id
3. Ao pagar financeiramente:
   - dt_pagamento preenchida
   - status = PAGO
   - ComissItem pode ser criado ou conciliado
4. Evita duplicidade por comissao_id
```

---

## 4. Estados

### 4.1 Máquina de Estados

```
PENDENTE → PARCIAL → PAGO
    │          │
    └──────────┴──→ CANCELADO
```

### 4.2 Cálculo de Saldo

```
vlr_saldo = vlr_comissao - total_pago

total_pago = SUM(comis_ped_item.vlr_pago WHERE comis_id = comissão.id)
```

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `ComissoesGeralList` | Gestão geral de comissões |
| `ComisPedDetalhes` | Detalhes e pagamentos da comissão |
| `ComissoesList` | Visão do representante |

---

## 6. Decisões Arquiteturais

### ADR-001: Comissão Agregada + Itens de Pagamento

**Decisão**: `comissoes` guarda total/saldo; `comis_ped_item` guarda pagamentos.

**Consequências**:
- Pagamentos parciais suportados
- Saldo calculado de forma simples
- Auditoria por item de pagamento

### ADR-002: Integração Financeira por comissao_id

**Decisão**: `FinContaPagar` possui `comissao_id`.

**Consequências**:
- Rastreabilidade entre comissão e pagamento financeiro
- Evita duplicidade de lançamento
- Permite conciliação

### ADR-003: Visões Separadas para Gestão e Representante

**Decisão**: `ComissoesGeralList` para gestão; `representantes/comissoes/ComissoesList` para representante.

**Consequências**:
- Permissões simples
- Interface do representante somente leitura/consulta

### ADR-004: Auditoria Automática

**Decisão**: `Comissoes` e `ComissItem` usam `SystemChangeLogTrait`.

**Consequências**:
- Rastreia alterações de valores e pagamentos
- Compliance em pagamentos de comissão

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `ComissoesGeralList` | `app/control/comissoes/ComissoesGeralList.php` |
| `ComisPedDetalhes` | `app/control/comissoes/ComisPedDetalhes.php` |
| `ComissoesList` | `app/control/representantes/comissoes/ComissoesList.php` |

### 7.2 Models

| Model | Tabela |
|-------|--------|
| `Comissoes` | `comissoes` |
| `ComissItem` | `comis_ped_item` |
| `FinContaPagar` | `fin_contas_pagar` |

### 7.3 Campos Principais

| Campo | Descrição |
|-------|-----------|
| `ped_id` | Pedido associado |
| `vlr_comissao` | Valor total devido |
| `total_pago` | Soma paga |
| `vlr_saldo` | Saldo pendente |
| `dt_prevista` | Previsão de pagamento |
| `stt_comis` | Status da comissão |
| `vlr_pago` | Valor de item pago |
| `dt_pgto` | Data do pagamento |
| `obs_pgto` | Observação do pagamento |

---

## 8. Segurança

- Representante vê apenas suas próprias comissões
- Gestão/financeiro controla pagamentos
- Valores e saldos auditados
- Cancelamentos exigem permissão elevada

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
