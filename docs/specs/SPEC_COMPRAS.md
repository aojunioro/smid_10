# SPEC_COMPRAS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Compras** gerencia a aquisição, coleta, transporte e chegada dos produtos comprados para atender pedidos. Ele conecta pedido, fornecedor, transportadora, nota fiscal, frete, parcelas e pagamento da compra, permitindo acompanhar o ciclo de reposição/compra vinculado a vendas.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Compras** | Registro principal da compra | 0..N por pedido |
| **ComprStatus** | Status da compra | catálogo |
| **ComprFornec** | Fornecedores | catálogo |
| **ComprTransport** | Transportadoras | catálogo |
| **PedidosRepre** | Pedido que originou/justifica a compra | 0..1 |

### 1.3 Relacionamentos

```
Pedido 1 ──── 0..N Compras
Compras N ──── 1 ComprStatus
Compras N ──── 1 ComprFornec
Compras N ──── 0..1 ComprTransport
Compras N ──── 1 SystemUser (login)
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Entrada | Compra vinculada a `ped_id` |
| PRODUTOS | Indireta | Produtos do pedido indicam necessidade de compra |
| FINANCEIRO | Saída | dt_pgto, parcelas, valor da compra e frete |
| SUPORTE | Indireta | Problemas pós-venda podem gerar nova compra/troca |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Compra** | Aquisição de produto/serviço junto a fornecedor |
| **Fornecedor** | Empresa que fornece o produto |
| **Transportadora** | Empresa responsável pela entrega/coleta |
| **dt_compr** | Data em que a compra foi realizada |
| **dt_coleta** | Data de coleta pela transportadora |
| **dt_chegada** | Data de chegada/recebimento |
| **dt_pgto** | Data de pagamento da compra |
| **frete** | Valor do frete |
| **vlr_compr** | Valor da compra |
| **n_nf** | Número da nota fiscal |
| **n_parcelas** | Quantidade de parcelas |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Compra Vinculada a Pedido

```
1. Usuário acessa ComprasForm
2. Seleciona ou recebe ped_id contextual
3. Sistema carrega dados do pedido:
   - n_ped
   - cliente
   - dt_ped
   - status do pedido
4. Preenche:
   - fornecedor
   - status da compra
   - transportadora
   - dt_compr
   - frete
   - vlr_compr
   - n_nf
   - n_parcelas
   - dt_pgto
5. Salva compra
6. SystemChangeLogTrait registra criação
```

**Regras**:
- RN-001: Compra pode existir sem transportadora inicialmente
- RN-002: Compra vinculada a pedido permite rastreabilidade total
- RN-003: Valor da compra e frete devem ser não negativos

### 3.2 Acompanhamento Logístico

```
1. Compra em status inicial (ex: Solicitada)
2. Quando fornecedor confirma coleta:
   - Preenche dt_coleta
   - Seleciona transportadora
   - Atualiza status
3. Quando produto chega:
   - Preenche dt_chegada
   - Atualiza status para recebido/concluído
4. ComprasList exibe datas e status coloridos
```

### 3.3 Pagamento da Compra

```
1. Financeiro/gestor preenche n_parcelas e dt_pgto
2. Sistema registra valor total e frete
3. Dados podem alimentar lançamento financeiro futuro
```

### 3.4 Listagem com Detalhes

`Compras::fetchListWithDetails()` consolida:
- Compra
- Pedido (`n_ped`, cliente, dt_ped, status pedido)
- Status da compra
- Fornecedor
- Transportadora
- Valores/datas

Filtros esperados:
- busca rápida
- status compra
- fornecedor
- transportadora
- período

---

## 4. Estados

### 4.1 Máquina de Estados Sugerida

```
SOLICITADA → COMPRADA → COLETADA → EM_TRANSPORTE → RECEBIDA
      │            │            │             │
      └────────────┴────────────┴─────────────┴──→ CANCELADA
```

### 4.2 Marcos por Data

| Campo | Marco Operacional |
|-------|-------------------|
| `dt_compr` | Compra realizada |
| `dt_coleta` | Coleta iniciada |
| `dt_chegada` | Produto recebido |
| `dt_pgto` | Pagamento realizado/agendado |

---

## 5. Cadastros Auxiliares

| Cadastro | Controller |
|----------|------------|
| Status da Compra | `CompraStatusForm/List` |
| Fornecedores | `ComprasFornecedorForm/List` |
| Transportadoras | `ComprasTransportadoraForm/List` |

---

## 6. Decisões Arquiteturais

### ADR-001: Compra como Entidade Vinculada ao Pedido

**Decisão**: Compra possui `ped_id`, mantendo rastreabilidade com venda.

**Consequências**:
- Permite ver compra por pedido
- Suporta múltiplas compras para um pedido
- Relatórios podem cruzar margem venda × custo real

### ADR-002: Listagem via Query Consolidada no Model

**Decisão**: `Compras::fetchListWithDetails()` monta query com joins.

**Consequências**:
- Controller sem SQL
- Performance e dados ricos na listagem
- Centralização de filtros no model

### ADR-003: Transportadora Opcional

**Decisão**: `transp_id` pode estar vazio até etapa logística.

**Consequências**:
- Compra pode ser registrada antes de coleta
- Validação de transportadora depende do status

### ADR-004: Auditoria Automática

**Decisão**: `Compras` usa `SystemChangeLogTrait`.

**Consequências**:
- Rastreabilidade de alterações em status/datas/valores
- Base para painel de auditoria futuro

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `ComprasForm` | `app/control/compras/ComprasForm.php` |
| `ComprasList` | `app/control/compras/ComprasList.php` |
| `ComprasDetalhes` | `app/control/compras/ComprasDetalhes.php` |
| `CompraStatusForm/List` | `app/control/compras/status/` |
| `ComprasFornecedorForm/List` | `app/control/compras/fornecedores/` |
| `ComprasTransportadoraForm/List` | `app/control/compras/transportadoras/` |

### 7.2 Models

| Model | Tabela |
|-------|--------|
| `Compras` | `compras` |
| `ComprStatus` | status de compras |
| `ComprFornec` | fornecedores |
| `ComprTransport` | transportadoras |

**Observação**: `Compras.php` referencia `ComprStatus`, `ComprFornec` e `ComprTransport` como relações do domínio.

### 7.3 Campos Principais

| Campo | Descrição |
|-------|-----------|
| `ped_id` | Pedido vinculado |
| `fornec_id` | Fornecedor |
| `status_id` | Status da compra |
| `transp_id` | Transportadora |
| `login` | Usuário que registrou |
| `dt_compr` | Data da compra |
| `dt_coleta` | Data de coleta |
| `dt_chegada` | Data de chegada |
| `frete` | Valor do frete |
| `vlr_compr` | Valor da compra |
| `n_nf` | Nota fiscal |
| `n_parcelas` | Parcelamento |
| `dt_pgto` | Data pagamento |

---

## 8. Segurança

- Alterações financeiras (valor, frete, pagamento) restritas a perfis autorizados
- Compra usa audit log automático
- Não excluir fisicamente registros com vínculo financeiro
- Campos de NF/pagamento devem ser validados no backend

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
