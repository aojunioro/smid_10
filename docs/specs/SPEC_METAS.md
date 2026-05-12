# SPEC_METAS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Metas** permite que gestores definam objetivos comerciais flexíveis para faturamento, vendas e agendamentos em diferentes níveis organizacionais: individual, unidade ou equipe. O sistema calcula realizado, percentual e projeção com base em pedidos e visitas dentro de um período customizado.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Meta** | Configuração da meta | 0..N |
| **MetaService** | Serviço de cálculo de realizado/percentual/projeção | service |
| **SystemUser** | Usuário alvo (atendente/representante) | permission |
| **SystemUnit** | Unidade alvo | permission |
| **Equipe** | Equipe alvo | smid |
| **PedidosRepre** | Fonte de faturamento/vendas | leitura |
| **Visita** | Fonte de agendamentos | leitura |

### 1.3 Relacionamentos

```
Meta ──── entidade_tipo = atendente/representante → SystemUser.login
Meta ──── entidade_tipo = unidade → SystemUnit.id
Meta ──── entidade_tipo = equipe → Equipe.id
MetaService ──── lê PedidosRepre para faturamento/vendas
MetaService ──── lê Visita para agendamentos
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| PEDIDOS | Leitura | Cálculo de faturamento e vendas |
| VISITAS | Leitura | Cálculo de agendamentos |
| LEADS | Leitura indireta | Filtro de unidade via leads |
| EQUIPES | Leitura | Metas por equipe |
| PERMISSION | Leitura | Usuários e unidades |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Meta** | Objetivo comercial em um período |
| **Tipo de Meta** | Faturamento, vendas ou agendamentos |
| **Entidade Alvo** | Quem/qual organização deve cumprir a meta |
| **Período Flexível** | Intervalo com início e fim livres |
| **Realizado** | Valor apurado até o momento |
| **Percentual** | Realizado / alvo * 100 |
| **Projeção** | Estimativa baseada no ritmo atual |
| **Semáforo** | Cor visual conforme progresso |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Meta

```
1. Gestor acessa MetaForm
2. Define:
   - tipo_meta: faturamento, vendas ou agendamentos
   - entidade_tipo: atendente, representante, unidade ou equipe
   - entidade_login ou entidade_id conforme tipo
   - valor_meta
   - dt_inicio
   - dt_fim
   - descrição/observação
3. Form alterna campos via TScript conforme entidade_tipo
4. Salva Meta com timestamps
```

**Regras**:
- RN-001: `dt_fim` deve ser maior ou igual a `dt_inicio`
- RN-002: `valor_meta` deve ser maior que zero
- RN-003: entidade_login é obrigatória para metas individuais
- RN-004: entidade_id é obrigatória para unidade/equipe

### 3.2 Cálculo de Realizado

`MetaService::getRealizado($meta)`:

```
1. Identifica período (dt_inicio até dt_fim)
2. Monta filtro de entidade:
   - atendente/representante: login_recep ou login_repre
   - unidade: lead_id IN (SELECT id FROM leads WHERE unidd_id = :id)
   - equipe: busca logins no permission/system_users por equipe_id
3. Identifica tipo_meta:
   - faturamento: SUM(pedidos.total_ped)
   - vendas: COUNT(pedidos.id)
   - agendamentos: COUNT(visitas.id)
4. Retorna número realizado
```

### 3.3 Percentual e Semáforo

```
percentual = min(100, (realizado / valor_meta) * 100)

Cores:
- Vermelho: < 50%
- Amarelo: 50% a 89%
- Verde: >= 90%
```

### 3.4 Projeção

```
projecao = (realizado / dias_passados) * dias_totais

Se período ainda não começou: projeção = 0
Se dias_passados <= 0: projeção = realizado
```

### 3.5 Listagem Administrativa

```
1. MetaList carrega metas
2. Para cada meta:
   - calcula realizado
   - calcula percentual
   - calcula projeção
   - exibe barra de progresso
3. Busca rápida com debounce 400ms
4. Filtros persistidos no formulário
```

### 3.6 Dashboard de Metas

```
1. MetaDashboard carrega metas filtradas
2. Exibe gráfico Gantt (TGantt)
3. Visualização inicial: mês anterior + mês atual
4. Usuário navega de 2 em 2 meses
5. Estado de zoom e tamanho persistido em sessão
6. ReflectionObject limpa rows do TGantt antes de recarregar
```

---

## 4. Tipos de Meta

| Tipo | Fonte | Cálculo |
|------|-------|---------|
| **Faturamento** | `pedidos` | `SUM(total_ped)` |
| **Vendas** | `pedidos` | `COUNT(id)` |
| **Agendamentos** | `visitas` | `COUNT(id)` |

---

## 5. Entidades Alvo

| Entidade | Campo | Fonte |
|----------|-------|-------|
| Atendente | `entidade_login` | `system_users.login` / `login_recep` |
| Representante | `entidade_login` | `system_users.login` / `login_repre` |
| Unidade | `entidade_id` | `system_unit.id` / `leads.unidd_id` |
| Equipe | `entidade_id` | `equipes.id` / logins da equipe |

---

## 6. Visualizações

| Tela | Função |
|------|--------|
| `MetaForm` | Cadastro/edição de metas |
| `MetaList` | Listagem com barras de progresso |
| `MetaDashboard` | Gantt com períodos e progresso |

---

## 7. Decisões Arquiteturais

### ADR-001: Período Flexível

**Decisão**: Usar `dt_inicio` e `dt_fim` em vez de periodicidade fixa.

**Consequências**:
- Campanhas de qualquer duração
- Metas trimestrais, mensais ou promocionais no mesmo modelo

### ADR-002: MetaService Multidatabase

**Decisão**: Cálculo usa `smid` para pedidos/visitas e `permission` para usuários/unidades.

**Consequências**:
- Transações devem alternar bases com cuidado
- Equipes exigem busca de logins antes do cálculo

### ADR-003: MetaList como TPage Customizada

**Decisão**: Não usar TStandardList para permitir controle manual e AJAX estático.

**Consequências**:
- Busca rápida com debounce 400ms
- Evita loop usando `$form->setData()` em vez de `TForm::sendData`

### ADR-004: TGantt com ReflectionObject

**Contexto**: TGantt não possui `clearRows` nativo.

**Decisão**: Usar ReflectionObject para limpar rows antes de recarregar.

**Consequências**:
- Evita duplicação de metas na navegação
- Dependência de estrutura interna do componente

### ADR-005: Semáforo Visual Padronizado

**Decisão**: Cores por percentual (<50, 50-89, >=90).

**Consequências**:
- Leitura executiva rápida
- Consistência entre lista e dashboard

---

## 8. Referência à Implementação Atual

### 8.1 Controllers

| Controller | Localização |
|------------|-------------|
| `MetaForm` | `app/control/metas/MetaForm.php` |
| `MetaList` | `app/control/metas/MetaList.php` |
| `MetaDashboard` | `app/control/metas/MetaDashboard.php` |

### 8.2 Model/Service

| Artefato | Localização |
|----------|-------------|
| `Meta` | `app/model/Meta.php` |
| `MetaService` | `app/lib/smid/MetaService.php` |

### 8.3 Docs/Scripts

| Referência | Descrição |
|------------|-----------|
| `docs/METAS.md` | Documento técnico do módulo |
| `app/database/scripts/9 - metas` | Scripts versionados do domínio |

---

## 9. Segurança

- Gestores criam/alteram metas
- Usuários comuns visualizam apenas metas aplicáveis ao seu perfil
- Filtros por unidade/equipe devem respeitar permissões
- Queries por equipe devem evitar interpolação insegura de logins
- Cálculos devem ignorar registros excluídos (`excluido_em IS NULL`)

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
