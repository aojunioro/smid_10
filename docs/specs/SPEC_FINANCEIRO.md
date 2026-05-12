# SPEC_FINANCEIRO.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Financeiro** centraliza contas a pagar, contas a receber, categorias financeiras, dashboards, notificações, extrato e lançamentos automáticos originados por pedidos, comissões, KM rodado e despesas extras de representantes.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **FinContaPagar** | Título/lançamento a pagar | 0..N |
| **FinContaReceber** | Título/lançamento a receber | 0..N |
| **FinCategoria** | Categoria financeira de receita/despesa | catálogo |
| **FinNotificacaoConfig** | Configuração de alertas financeiros | catálogo/config |
| **FinNotificacaoLog** | Histórico de notificações | 0..N |
| **RepreDespesaExtra** | Despesa extra de representante | 0..N |
| **RepreDespesaCateg** | Categoria de despesa extra | catálogo |
| **RepreDespesaComprovante** | Comprovante anexado à despesa | 0..N |

### 1.3 Relacionamentos

```
FinCategoria 1 ──── 0..N FinContaPagar
FinCategoria 1 ──── 0..N FinContaReceber
Pedido 1 ──── 0..N FinContaReceber
Pedido 1 ──── 0..N FinContaPagar
Comissao 1 ──── 0..N FinContaPagar
KmReembolsoLote 1 ──── 0..1 FinContaPagar
RepreDespesaExtra 1 ──── 0..1 FinContaPagar
SystemGroup 1 ──── 0..N FinContaPagar/Receber
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Entrada | Lançamentos automáticos de contas a receber |
| COMISSOES | Entrada | Contas a pagar de comissão |
| KM | Entrada | Contas a pagar de reembolso por lote |
| COMPRAS | Entrada futura | Contas a pagar de compras |
| SUPORTE | Leitura | Pendências financeiras podem impactar atendimento |
| NOTIFICAÇÕES | Saída | Alertas de vencimento, atraso e histórico |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Conta a Pagar** | Obrigação financeira da empresa |
| **Conta a Receber** | Valor que a empresa tem a receber |
| **Categoria Financeira** | Classificação contábil/gerencial (`R` receita, `D` despesa) |
| **Vencimento Original** | Data original antes de remarcações |
| **Lançamento Automático** | Criado por outro módulo (pedido, comissão, KM etc.) |
| **Recorrência** | Conjunto de contas geradas periodicamente |
| **Despesa Fixa / Receita Fixa** | Lançamento recorrente com término opcional |
| **Extrato** | Visão unificada de entradas e saídas |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Conta a Pagar

```
1. Usuário acessa FinContasPagarForm
2. Preenche:
   - categoria_id (tipo D)
   - descricao
   - valor
   - dt_vencimento
   - observacao
   - system_group_id
   - flags de recorrência/despesa fixa se aplicável
3. Sistema define status inicial PENDENTE
4. Se recorrência habilitada, gera parcelas vinculadas por recorrencia_id
5. Salva com SystemChangeLogTrait
```

**Regras**:
- RN-001: Categoria de conta a pagar deve ser tipo `D`
- RN-002: Valor deve ser maior que zero
- RN-003: `dt_venc_orig` preserva vencimento inicial

### 3.2 Pagamento de Conta a Pagar

```
1. Usuário abre FinContasPagarForm ou ação rápida
2. Informa dt_pagamento
3. Status passa para PAGO
4. Dashboard e extrato passam a considerar como realizado
```

### 3.3 Cadastro de Conta a Receber

```
1. Usuário acessa FinContasReceberForm
2. Preenche:
   - categoria_id (tipo R)
   - pedido_id opcional
   - cliente_nome / cliente_cpf
   - descricao
   - valor
   - dt_vencimento
3. Sistema define status inicial PENDENTE
4. Se receita fixa, gera recorrência
```

### 3.4 Recebimento de Conta

```
1. Usuário informa dt_recebimento
2. Status passa para RECEBIDO
3. Extrato exibe entrada realizada
```

### 3.5 Lançamento Automático por Pedido

```
1. Pedido é criado/aprovado
2. Regras financeiras definem categoria e valor
3. Cria FinContaReceber com:
   - pedido_id
   - cliente_nome / cpf
   - valor
   - dt_vencimento
   - lancamento_automatico = 'S'
4. Alterações posteriores podem reconciliar valor/status
```

### 3.6 Lançamento Automático por Comissão

```
1. Comissão é apurada
2. Cria FinContaPagar vinculada a comissao_id
3. Pagamento financeiro atualiza conta e pode refletir na comissão
```

### 3.7 Lançamento por KM/Reembolso

```
1. Lote de KM aprovado/pago
2. Cria FinContaPagar com km_reembolso_lote_id
3. Conta reflete valor_total do lote
```

### 3.8 Despesa Extra de Representante

```
1. Representante lança despesa extra com comprovantes
2. Gestor aprova via DespesaExtraAprovarForm
3. Financeiro paga via DespesaExtraPagarForm
4. FinContaPagar é criada/vinculada se aplicável
```

---

## 4. Estados

### 4.1 Contas a Pagar

```
PENDENTE → PAGO
    │
    ├── vencimento < hoje → VENCIDO
    └── cancelamento → CANCELADO
```

### 4.2 Contas a Receber

```
PENDENTE → RECEBIDO
    │
    ├── vencimento < hoje → VENCIDO
    └── cancelamento → CANCELADO
```

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `FinDashboard` | Visão executiva financeira |
| `FinDashboardIndicators` | Indicadores parciais/assíncronos |
| `FinContasPagarList/Form` | CRUD e controle de pagamentos |
| `FinContasReceberList/Form` | CRUD e controle de recebimentos |
| `FinExtratoList` | Extrato unificado |
| `FinCategoriasList/Form` | Categorias financeiras |
| `FinConfigNotificacoes/List` | Configuração de alertas |
| `FinNotificacaoHistorico` | Histórico de notificações |
| `DespesaExtra*` | Fluxo de despesas extras |

---

## 6. Decisões Arquiteturais

### ADR-001: Contas Separadas por Natureza

**Decisão**: Usar tabelas distintas para contas a pagar e receber.

**Consequências**:
- Campos especializados por natureza
- Extrato precisa unir fontes
- Categorias filtradas por tipo (`R`/`D`)

### ADR-002: Vínculos Explícitos com Módulos de Origem

**Decisão**: `FinContaPagar` possui `comissao_id`, `pedido_id`, `km_reembolso_lote_id`, `repre_despesa_extra_id`.

**Consequências**:
- Rastreabilidade de origem
- Evita duplicidade de lançamento automático
- Facilita auditoria cruzada

### ADR-003: Recorrência por recorrencia_id

**Decisão**: Contas recorrentes compartilham `recorrencia_id`.

**Consequências**:
- Exclusão/alteração em lote possível
- `FinContasPagarRecurrenceDelete` remove recorrência controlada

### ADR-004: Auditoria Automática

**Decisão**: Models financeiros usam `SystemChangeLogTrait`.

**Consequências**:
- Rastreia mudanças de valor, vencimento e status
- Compliance financeiro

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `FinDashboard` / `FinDashboardIndicators` | `app/control/financeiro/` |
| `FinContasPagarForm/List` | `app/control/financeiro/` |
| `FinContasReceberForm/List` | `app/control/financeiro/` |
| `FinExtratoList` | `app/control/financeiro/` |
| `FinCategoriasForm/List` | `app/control/financeiro/` |
| `FinConfigNotificacoes/List` | `app/control/financeiro/` |
| `FinNotificacaoHistorico` | `app/control/financeiro/` |
| `DespesaExtra*` | `app/control/financeiro/despesaExtra/` |

### 7.2 Models

| Model | Tabela |
|-------|--------|
| `FinContaPagar` | `fin_contas_pagar` |
| `FinContaReceber` | `fin_contas_receber` |
| `FinCategoria` | `fin_categorias` |
| `FinNotificacaoConfig` | config notificações |
| `FinNotificacaoLog` | log notificações |
| `RepreDespesaExtra` | despesas extras |
| `RepreDespesaCateg` | categorias despesa extra |
| `RepreDespesaComprovante` | comprovantes |

---

## 8. Segurança

- Alterações de valor/status restritas a perfis financeiros
- Lançamentos automáticos devem ser idempotentes
- Cancelamentos exigem observação
- Notificações não podem expor dados sensíveis indevidamente
- Auditoria obrigatória em valores e datas

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
