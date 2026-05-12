# SPEC_RELATORIOS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Relatórios** consolida indicadores operacionais, comerciais e gerenciais do SMID. Ele transforma dados de Leads, Visitas, Históricos, Pedidos, Televendas, Atendentes, Representantes, Mídias, Unidades, Equipes e integrações telefônicas em visões analíticas e dashboards.

### 1.2 Famílias de Relatórios

| Família | Descrição |
|---------|-----------|
| **Aproveitamento** | Conversão/eficiência por dimensão |
| **Comparativo** | Comparação entre atendentes, representantes, mídias, equipes e unidades |
| **Ligações do Dia** | Métricas de chamadas por atendente/mídia |
| **Evolução de Vendas** | Tempo entre lead, visita e venda |
| **Visitas do Dia Seguinte** | Dashboard operacional de agenda |
| **Dashboards de domínio** | Financeiro, Metas, Logs, Representante, Admin |

---

## 2. Entidades e Fontes

| Fonte | Uso |
|-------|-----|
| `leads` | Origem, mídia, atendente, unidade, status |
| `visitas` | Agendamentos, representantes, status, datas |
| `historicos` | Motivos, ocorridos, resultados |
| `pedidos` | Vendas, faturamento, datas, status |
| `lead_status_categoria` | Categorias funcionais de leads |
| `ped_status_categoria` | Categorias funcionais de pedidos |
| Integrações 3C+/GoTo | Ligações e status de chamadas |
| `system_users` / `system_unit` | Usuários e unidades |

---

## 3. Relatórios Principais

### 3.1 Aproveitamento

Relatórios:
- `AproveitamentoAtendentesReport`
- `AproveitamentoEquipesReport`
- `AproveitamentoMidiasReport`
- `AproveitamentoRepreReport`
- `AproveitamentoUnidadesReport`

Métricas típicas:
- novos leads;
- atendidos;
- agendados;
- vendidos;
- perdidos;
- taxa de conversão;
- faturamento quando aplicável.

### 3.2 Comparativos

Relatórios:
- `ComparativoAtendentesReport`
- `ComparativoEquipesReport`
- `ComparativoMidiasReport`
- `ComparativoRepreReport`
- `ComparativoUnidadesReport`

Objetivo:
```
Comparar desempenho entre entidades no mesmo período e com os mesmos critérios.
```

### 3.3 Ligações do Dia

Relatórios:
- `LigacoesDiaAtendentesReport`
- `LigacoesDiaMidiasReport`

Métricas:
- novas ligações;
- atendidas;
- abandonadas;
- agendadas;
- duplicadas;
- percentuais por status.

Services:
- `LigacoesDiaAtendentesReportService`
- `LigacoesDiaMidiasReportService`

### 3.4 Evolução de Vendas

Dashboard:
- `DashboardEvolucao`

Service:
- `EvolucaoVendasReportService`

Mede tempos entre:
```
T0: criação do lead
T1: criação/agendamento da visita
T2: data efetiva da visita
T3: criação/concretização do pedido
```

Categorias válidas de pedido usam `PedStatusCategoria`.

### 3.5 Visitas do Dia Seguinte

Dashboard:
- `LeadVisitNextDayDashboard`

Service:
- `LeadVisitNextDayDashboardService`

Regras destacadas:
- auto refresh a cada 60 segundos;
- horário comercial 8h–21h;
- mínimo de visitas por representante: 3;
- análise por atendente, representante e unidade.

---

## 4. Fluxo Geral de Relatório

```
1. Usuário abre relatório
2. Sistema carrega filtros padrão
3. Usuário ajusta período/dimensões
4. Controller valida filtros
5. Service monta consulta/agregação
6. Dados são normalizados em arrays/DTOs
7. Tela renderiza cards, tabelas ou gráficos
8. Exportação/impressão pode ser disponibilizada conforme tela
```

---

## 5. Filtros e Critérios

| Filtro | Uso |
|--------|-----|
| Período | Base temporal do relatório |
| Unidade | Restrição operacional |
| Atendente | Login de recepção/callcenter |
| Representante | Login comercial externo |
| Mídia | Origem do lead |
| Equipe | Agrupamento de usuários |
| Status/Categoria | Normalização por categoria funcional |

**Regra**: relatórios devem preferir categorias funcionais (`LeadStatusHelper`, `PedStatusHelper`) em vez de IDs fixos.

---

## 6. Decisões Arquiteturais

### ADR-001: Services para Agregação

**Decisão**: Cálculos ficam em `app/model/report/*Service.php`.

**Consequências**:
- Controllers focam UI/filtros
- Regras reutilizáveis
- Melhor portabilidade para outra linguagem

### ADR-002: Categorias Funcionais de Status

**Decisão**: Relatórios não devem depender de IDs fixos de status.

**Consequências**:
- Suporta customização por cliente
- Exige helpers canônicos de status

### ADR-003: SQL Otimizado em Services/Models

**Decisão**: Relatórios podem usar SQL agregado em camada de model/service.

**Consequências**:
- Performance superior para dashboards
- Controllers permanecem sem SQL

### ADR-004: Dashboards por Domínio

**Decisão**: Dashboards específicos vivem próximos ao domínio, mas pertencem conceitualmente a Relatórios.

**Consequências**:
- Financeiro, Metas, Logs e Representantes têm dashboards próprios
- SPEC_RELATORIOS documenta padrões comuns

---

## 7. Referência à Implementação Atual

### 7.1 Controllers de Relatórios

| Área | Localização |
|------|-------------|
| Atendentes | `app/control/relatorios/atendentes/` |
| Equipes | `app/control/relatorios/equipes/` |
| Mídias | `app/control/relatorios/midias/` |
| Representantes | `app/control/relatorios/representantes/` |
| Unidades | `app/control/relatorios/unidades/` |
| Dashboards | `app/control/*/*Dashboard*.php`, `app/control/relatorios/dashboards/` |

### 7.2 Services

| Service | Função |
|---------|--------|
| `AproveitamentoReportService` | Cálculo de aproveitamento com categorias |
| `EvolucaoVendasReportService` | Tempos da jornada comercial |
| `LeadVisitNextDayDashboardService` | Agenda operacional do dia seguinte |
| `LigacoesDiaAtendentesReportService` | Chamadas por atendente |
| `LigacoesDiaMidiasReportService` | Chamadas por mídia |
| `ThreeCPlusCallReportService` | Relatórios de chamadas 3C+ |
| `GotoConnectCallReportService` | Relatórios de chamadas GoTo |

---

## 8. Segurança e Governança

- Respeitar unidades permitidas do usuário
- Não expor dados pessoais além do necessário
- Filtros devem ser validados no backend
- Consultas agregadas devem evitar SQL injection
- Exportações devem obedecer ao mesmo escopo visualizado

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
