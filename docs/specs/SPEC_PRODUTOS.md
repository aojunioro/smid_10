# SPEC_PRODUTOS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Produtos** define o catálogo comercial do SMID. Ele sustenta pedidos de representantes, orçamentos de televendas, compras e análises de margem. Um produto possui categoria, medida, fornecedor, preços de compra/venda, parâmetros de estoque e classificação por canal comercial.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Produto** | Produto base vendido ou orçado | catálogo |
| **ProdCateg** | Categoria do produto | catálogo |
| **ProdMedidas** | Unidade/medida do produto | catálogo |
| **ProdModelos** | Modelo/variante do produto | catálogo |
| **PedProdItem** | Item de produto em pedido | 0..N por pedido |
| **OrcamProdItem** | Item de produto em orçamento | 0..N por orçamento |
| **ProdOrcamento** | Produto comercial específico para Televendas | catálogo/ponte |

### 1.3 Relacionamentos

```
Produto N ──── 1 ProdCateg
Produto N ──── 1 ProdMedidas
Produto N ──── 0..1 Fornecedor
Produto 1 ──── 0..N PedProdItem
Produto 1 ──── 0..N ProdOrcamento
ProdOrcamento 1 ──── 0..N OrcamProdItem
Pedido 1 ──── 0..N PedProdItem
Orçamento 1 ──── 0..N OrcamProdItem
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Saída | Produtos compõem itens de pedido |
| TELEVENDAS | Saída | Produtos com `televendas='S'` compõem orçamentos/pedidos tele |
| COMPRAS | Entrada/Saída | Fornecedor, preço de compra e reposição |
| FINANCEIRO | Leitura | Margem, custo, valor líquido |
| ESTOQUE | Futuro | Estoque mínimo/máximo |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Produto Base** | Cadastro principal em `produtos` |
| **Produto Televendas** | Produto habilitado para canal televendas (`televendas='S'`) |
| **Produto Domicílio** | Produto vendido por representantes (`televendas='N'`) |
| **Categoria** | Agrupamento comercial/visual do produto |
| **Medida** | Unidade de medida (un, cx, kit, etc.) |
| **Modelo** | Variação/modelo do produto |
| **Preço de Compra** | Custo de aquisição (`vlr_prod_compra`) |
| **Preço de Venda** | Preço base comercial (`vlr_prod_venda`) |
| **Preço Tabela Tele** | Preço do produto em orçamento (`vlr_tabela`) |
| **Canal Exclusivo** | Regra de separação Televendas S/N, sem modo "ambos" |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Produto Base

```
1. Usuário acessa ProdutoForm
2. Preenche:
   - nome_prod
   - categ_id
   - fornec_id
   - med_id
   - vlr_prod_compra
   - vlr_prod_venda
   - estoq_min / estoq_max
   - ativo (S/N)
   - televendas (S/N)
3. Sistema valida campos obrigatórios e valores monetários
4. Salva Produto
5. ProdutoList exibe categoria, fornecedor, preços, ativo e canal Televendas
```

**Regras**:
- RN-001: Produto inativo não aparece em combos operacionais
- RN-002: `vlr_prod_venda` deve ser maior ou igual a zero
- RN-003: `vlr_prod_compra` deve ser maior ou igual a zero
- RN-004: `televendas` aceita somente `S` ou `N`

### 3.2 Separação por Canal Comercial

**Regra Canônica**: O campo `produtos.televendas` é exclusivo por canal.

```
- Televendas: combos exibem somente produtos com televendas='S'
- Representantes/domicílio: combos exibem somente produtos com televendas='N'
- Não existe modo "ambos"
```

**Impacto**:
- `PedidosTeleStep4Form` bloqueia avanço se produto base não for `S`
- Fluxos de representantes filtram apenas `N`
- `TeleProdutosForm` força vínculo com produto base `S`

### 3.3 Inclusão Rápida pelo Televendas

```
1. Usuário está em TeleProdutosForm
2. Clica '+' ao lado da combo Produto Base
3. ProdutoForm abre em cortina contextual
4. Campo Televendas vem pré-definido como S e bloqueado
5. Backend força persistência `televendas='S'`
6. Ao salvar, retorna para TeleProdutosForm com novo produto selecionado
```

### 3.4 Uso em Pedidos

```
1. PedProdDetalhes carrega produtos ativos conforme canal
2. Usuário seleciona produto
3. Sistema preenche medida e preço base
4. Usuário informa quantidade e descontos, quando aplicável
5. Calcula total do item
6. Pedido recalcula total_ped
```

### 3.5 Uso em Orçamentos Televendas

```
1. OrcamForm / PedidosTeleStep4Form carrega ProdOrcamento
2. Produto precisa ter vínculo com Produto base (`prod_id`)
3. Produto base deve estar marcado como `televendas='S'`
4. OrcamProdItem armazena quantidade, valor item e valor total
```

---

## 4. Cadastros Auxiliares

| Cadastro | Controller | Campos principais |
|----------|------------|-------------------|
| Categoria | `ProdCategForm/List` | categoria, cor |
| Medida | `ProdMedidasForm/List` | medida/sigla |
| Modelo | `ProdModelosForm/List` | modelo |
| Produto Base | `ProdutoForm/List` | nome, categoria, fornecedor, preços, canal |
| Produto Televendas | `TeleProdutosForm/List` | produto base, preço tabela, medida |

---

## 5. Decisões Arquiteturais

### ADR-001: Canal Exclusivo por Produto

**Decisão**: `produtos.televendas` é `S/N`, sem terceira opção.

**Consequências**:
- Simplicidade nos filtros
- Evita venda de produto errado por canal
- Necessidade de duplicar produto se for vendido nos dois canais com regras distintas

### ADR-002: Produto Televendas como Camada Comercial

**Decisão**: Televendas usa `ProdOrcamento`/`orcam_prod` como camada de preço comercial, vinculada ao Produto base.

**Consequências**:
- Preço de tabela separado do preço base
- Suporte a orçamento antes de pedido
- Produto base continua como fonte de categoria/medida/canal

### ADR-003: Estoque como Campo Preparatório

**Decisão**: Produto mantém `estoq_min` e `estoq_max`, mesmo sem módulo de estoque completo.

**Consequências**:
- Base pronta para estoque futuro
- Hoje serve como informação gerencial

### ADR-004: Categoria com Cor

**Decisão**: Categoria de produto possui `cor` para UI.

**Consequências**:
- Badges em listas e cards
- Visual consistente com status/equipes

---

## 6. Referência à Implementação Atual

### 6.1 Controllers

| Controller | Localização |
|------------|-------------|
| `ProdutoForm` / `ProdutoList` | `app/control/produtos/` |
| `ProdCategForm` / `ProdCategList` | `app/control/produtos/categoria/` |
| `ProdMedidasForm` / `ProdMedidasList` | `app/control/produtos/medida/` |
| `ProdModelosForm` / `ProdModelosList` | `app/control/produtos/modelo/` |
| `TeleProdutosForm` / `TeleProdutosList` | `app/control/televendas/produtos/` |
| `PedProdDetalhes` | `app/control/pedidos/PedProdDetalhes.php` |

### 6.2 Models

| Model | Tabela |
|-------|--------|
| `Produto` | `produtos` |
| `ProdCateg` | `prod_categ` |
| `ProdMedidas` | medidas de produto |
| `ProdModelos` | modelos de produto |
| `PedProdItem` | itens de pedido |
| `OrcamProdItem` | `orcam_prod_item` |
| `ProdOrcamento` | produtos de orçamento/televendas |

### 6.3 Scripts e Docs

| Referência | Descrição |
|------------|-----------|
| `docs/televendas/PRODUTOS_CANAL_TELEVENDAS.md` | Regra de canal Televendas S/N |
| `app/database/scripts/13 - televendas/004_televendas_canal_produtos.sql` | Adiciona `produtos.televendas` |

---

## 7. Segurança e Validação

- Somente usuários autorizados editam produto base
- Produto inativo não deve ser usado em novos pedidos/orçamentos
- Validação backend obrigatória para canal em Televendas
- Preços devem ser numéricos e não negativos

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
