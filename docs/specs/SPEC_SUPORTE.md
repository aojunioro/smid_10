# SPEC_SUPORTE.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Suporte** gerencia atendimentos pós-venda vinculados a pedidos. Um suporte registra solicitações do cliente, priorização, departamento responsável, responsável atribuído, SLA, relato do cliente, relato técnico, anexos e resolução. Ele conecta a etapa de pedido/entrega com atendimento operacional e eventual necessidade de compra/troca.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Suporte** | Chamado/atendimento de suporte | 0..N por pedido |
| **Pedido** | Pedido relacionado ao chamado | 1 por suporte |
| **Status Suporte** | Status do chamado (via status_id) | catálogo |
| **Solicitante** | Origem/tipo de solicitante (`solicit_id`) | catálogo |
| **Departamento** | Departamento responsável (`depart_id`) | catálogo |
| **Usuário Responsável** | Login atribuído (`atrib_login`) | 0..1 |

### 1.3 Relacionamentos

```
Pedido 1 ──── 0..N Suporte
Suporte N ──── 1 StatusSuporte
Suporte N ──── 1 Solicitante
Suporte N ──── 1 Departamento
Suporte N ──── 1 SystemUser (login - criador)
Suporte N ──── 0..1 SystemUser (atrib_login - responsável)
Pedido N ──── 1 Lead (via pedidos.lead_id)
Lead N ──── 1 Unidade (filtro por unidade)
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Entrada | Suporte sempre relacionado a pedido (`ped_id`) |
| LEADS | Leitura | Busca dados do cliente via pedido.lead_id |
| COMPRAS | Indireta | Suporte pode demandar compra/troca/reposição |
| FINANCEIRO | Leitura | Situação de pagamento pode afetar atendimento |
| AUDITORIA | Saída | SystemChangeLogTrait registra alterações |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Suporte** | Chamado pós-venda ou atendimento técnico |
| **SLA** | Prazo limite (`dt_limit`) para resolução |
| **Prioridade** | Grau de urgência do suporte |
| **Relato Cliente** | Descrição inicial feita pelo cliente |
| **Relato Técnico** | Retorno técnico/resolução |
| **Imagem Ordem** | Anexo/foto da ordem ou evidência (`img_ordem`) |
| **Atribuição** | Designação de responsável (`atrib_login`) |
| **Departamento** | Área que deve resolver o suporte |

---

## 3. Fluxos Principais

### 3.1 Abertura de Suporte

```
1. Usuário acessa SuporteForm
2. Seleciona pedido (ped_id) ou recebe contextual do pedido
3. Sistema carrega dados do cliente via pedido → lead
4. Preenche:
   - status_id inicial
   - fone_sup
   - solicit_id
   - depart_id
   - prioridade
   - dt_sup
   - dt_limit
   - relato_cli
   - img_ordem (opcional)
5. Sistema define login = usuário logado
6. Salva suporte
7. SystemChangeLogTrait registra criação
```

**Regras**:
- RN-001: `ped_id` é obrigatório
- RN-002: `relato_cli` é obrigatório na abertura
- RN-003: `dt_limit` deve respeitar SLA por prioridade/departamento
- RN-004: `fone_sup` pode divergir do telefone do lead (contato alternativo)

### 3.2 Triagem e Atribuição

```
1. Suporte aparece em SuporteList
2. Supervisor/gestor analisa prioridade e departamento
3. Define `atrib_login` (responsável)
4. Atualiza status para Em Atendimento
5. Responsável recebe chamado em sua fila
```

### 3.3 Resolução

```
1. Responsável abre SuporteForm
2. Preenche relato_tec com diagnóstico/solução
3. Opcionalmente atualiza img_ordem
4. Preenche dt_resol
5. Atualiza status para Resolvido/Concluído
6. Sistema registra alteração no audit log
```

**Regras**:
- RN-005: Para status resolvido, `relato_tec` e `dt_resol` são obrigatórios
- RN-006: `dt_resol` não pode ser anterior a `dt_sup`

### 3.4 Reabertura

```
1. Cliente retorna com problema não resolvido
2. Usuário reabre suporte ou cria novo vinculado ao mesmo pedido
3. Status volta para Aberto/Em Atendimento
4. dt_resol é preservada como histórico ou limpa conforme política
```

### 3.5 Listagem e Busca

`Suporte::getListData()` consolida dados de suporte + pedido + lead:
- filtro por telefone (`fone_sup`)
- busca rápida por nome, CPF ou ID suporte
- filtro por status
- restrição por unidades permitidas do usuário

---

## 4. Estados

### 4.1 Máquina de Estados Sugerida

```
ABERTO → TRIAGEM → EM_ATENDIMENTO → AGUARDANDO_CLIENTE → RESOLVIDO
   │          │            │                  │              │
   └──────────┴────────────┴──────────────────┴──────────────┘
                         REABERTO

ABERTO / EM_ATENDIMENTO → CANCELADO
```

### 4.2 Prioridades

| Prioridade | SLA sugerido | Descrição |
|------------|--------------|-----------|
| Baixa | 5 dias úteis | Dúvida ou solicitação simples |
| Média | 3 dias úteis | Problema operacional sem bloqueio |
| Alta | 1 dia útil | Cliente impactado |
| Crítica | mesmo dia | Risco comercial ou entrega travada |

---

## 5. Visualizações

### 5.1 SuporteList

- Listagem principal com paginação server-side
- Filtro por telefone
- Busca rápida por cliente/CPF/ID
- Filtro por status
- Restrição por unidade do lead
- Badges de prioridade e status

### 5.2 SuporteForm

- Cortina lateral para criação/edição
- Dados do pedido/cliente em contexto
- Campos de abertura e resolução
- Upload/visualização de `img_ordem`

### 5.3 SuportesGrid

- Grid embutido em Pedido ou Lead
- Lista suportes vinculados ao pedido
- Ações rápidas: abrir, editar, visualizar resolução

---

## 6. Decisões Arquiteturais

### ADR-001: Suporte Vinculado ao Pedido

**Decisão**: Chamado usa `ped_id` como vínculo obrigatório.

**Consequências**:
- Todo suporte tem contexto comercial/pós-venda
- Dados do cliente vêm via pedido → lead
- Permite múltiplos chamados por pedido

### ADR-002: Query Consolidada no Model

**Decisão**: `Suporte::getListData()` monta joins e filtros.

**Consequências**:
- Controller sem SQL
- Performance adequada na listagem
- Restrição por unidade aplicada no model

### ADR-003: Atribuição por Login

**Decisão**: Campo `atrib_login` guarda responsável.

**Consequências**:
- Simples e compatível com login legado
- Filtro por responsável direto
- Não depende de ID numérico do usuário

### ADR-004: Suporte com Dois Relatos

**Decisão**: Separar `relato_cli` e `relato_tec`.

**Consequências**:
- Clareza entre problema reportado e solução técnica
- `relato_cli` obrigatório na abertura
- `relato_tec` obrigatório na resolução

### ADR-005: Auditoria Automática

**Decisão**: `Suporte` usa `SystemChangeLogTrait`.

**Consequências**:
- Rastreia mudanças de status, responsável e relatos
- Suporte a compliance de atendimento

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `SuporteForm` | `app/control/suportes/SuporteForm.php` |
| `SuporteList` | `app/control/suportes/SuporteList.php` |
| `SuportesGrid` | `app/control/suportes/SuportesGrid.php` |

### 7.2 Model

| Model | Tabela | Auditoria |
|-------|--------|-----------|
| `Suporte` | `suportes` | SystemChangeLogTrait |

### 7.3 Campos Principais

| Campo | Descrição |
|-------|-----------|
| `ped_id` | Pedido vinculado |
| `status_id` | Status do chamado |
| `fone_sup` | Telefone de contato |
| `solicit_id` | Solicitante/origem |
| `depart_id` | Departamento responsável |
| `login` | Usuário criador |
| `atrib_login` | Usuário responsável |
| `prioridade` | Prioridade/SLA |
| `dt_sup` | Data de abertura |
| `dt_limit` | Prazo limite |
| `relato_cli` | Relato do cliente |
| `dt_resol` | Data de resolução |
| `relato_tec` | Relato técnico |
| `img_ordem` | Anexo/imagem |

---

## 8. Segurança

- Usuário só vê suportes de unidades permitidas
- Resolução restrita ao responsável, gestor ou admin
- Campos técnicos editáveis somente após atribuição
- `img_ordem` pode conter dado sensível; acesso restrito
- Auditoria automática em todas as alterações

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
