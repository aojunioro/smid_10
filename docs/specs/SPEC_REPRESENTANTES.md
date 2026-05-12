# SPEC_REPRESENTANTES.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Representantes** define a experiência operacional do usuário que executa visitas externas, registra históricos, altera status de visita, cria pedidos em campo, acompanha comissões, lança despesas extras e participa do cálculo de KM rodado. Ele é uma visão especializada sobre Visitas, Históricos, Pedidos, Comissões, Financeiro e KM.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **VisitaRepre** | Visão/modelo de visita para representante | 0..N por representante |
| **HistoricoRepre** | Histórico registrado em campo | 0..N por visita/lead |
| **PedidosRepre** | Pedido criado por representante | 0..N por representante |
| **Comissoes** | Comissão por pedido | 0..N |
| **ComissItem** | Pagamentos de comissão | 0..N por comissão |
| **RepreDespesaExtra** | Despesa extra lançada pelo representante | 0..N |
| **RepreDespesaCateg** | Categoria de despesa extra | catálogo |
| **RepreDespesaComprovante** | Comprovante de despesa | 0..N |
| **VisitaCheckinEvent** | Evento GPS de presença | 0..N por visita |
| **KmReembolsoLote** | Reembolso de KM | 0..N por representante |

### 1.3 Relacionamentos

```
Representante (SystemUser.login) 1 ──── 0..N VisitaRepre.login_repre
Representante 1 ──── 0..N PedidosRepre.login_repre
Representante 1 ──── 0..N RepreDespesaExtra.login_repre
Representante 1 ──── 0..N KmReembolsoLote.login_repre
VisitaRepre 1 ──── 0..N HistoricoRepre
Pedido 1 ──── 0..1 Comissoes
Comissoes 1 ──── 0..N ComissItem
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| VISITAS | Entrada/Saída | Lista e execução das visitas atribuídas |
| HISTORICOS | Saída | Registro de ocorridos e fotos |
| PEDIDOS | Saída | Criação de pedidos em campo e prospecção |
| COMISSOES | Leitura | Representante acompanha valores e pagamentos |
| KM | Saída/Leitura | GPS e visitas geram reembolso |
| FINANCEIRO | Saída | Despesas extras e comissões geram contas a pagar |
| LEADS | Leitura | Dados do cliente e histórico completo do lead |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Representante** | Usuário do grupo 4 que executa visitas e vendas externas |
| **Visita Atribuída** | Visita cujo `login_repre` é o login do representante |
| **Pedido em Campo** | Pedido criado pelo representante a partir de uma visita/lead |
| **Pedido Prospec** | Pedido de prospecção, criado em fluxo especializado |
| **Histórico em Campo** | Registro do ocorrido pelo representante durante/após visita |
| **Status da Visita** | Atualização operacional feita pelo representante |
| **Despesa Extra** | Gasto reembolsável lançado pelo representante |
| **Comissão** | Valor devido por venda vinculada ao representante |
| **KM Rodado** | Reembolso de deslocamento validado por GPS |

---

## 3. Fluxos Principais

### 3.1 Lista de Visitas do Representante

```
1. Representante acessa VisitaRepreList
2. Sistema filtra visitas por login_repre = usuário logado
3. Aplica regras de visibilidade:
   - visita ativa
   - status/ocorrido permitidos
   - regra opcional para ocultar motivo_id=1 "Apresentei e não vendi"
4. Lista exibe dados completos do lead, endereço, data/hora e status
5. Representante abre VisitaRepreForm em cortina
```

**Regras**:
- RN-001: Representante vê somente suas visitas, salvo perfil elevado
- RN-002: Regra de `motivo_id=1` depende de `application.php -> representantes.mostrar_motivo_apresentei_nao_vendi`

### 3.2 Visualização de Visita e Navegação em Cortina

```
1. VisitaRepreList abre no conteúdo principal
2. VisitaRepreForm abre na cortina com onView&key=<visita_id>
3. Ações internas mantêm target `adianti_right_panel`
4. Ao abrir histórico/status, a cortina empilha novos forms
5. Retorno fecha apenas o topo e recarrega VisitaRepreForm com key válido
```

**Regras**:
- RN-003: Nunca retornar via `window.parent.__adianti_load_page` sem target explícito
- RN-004: Histórico deve listar todo histórico do lead, não apenas da visita atual

### 3.3 Alteração de Status da Visita

```
1. Representante acessa VisitaRepreStatus
2. Seleciona novo status operacional
3. Sistema captura GPS via VisitaRepreGPS/GpsPermissionGate se exigido
4. Cria VisitaCheckinEvent
5. Atualiza visita.status_id/stts_lead conforme regra
6. LeadVisitStatusSyncService sincroniza lead quando aplicável
7. Retorna para VisitaRepreForm mantendo cortina
```

**Regras**:
- RN-005: GPS pode ser obrigatório para status que confirmam presença
- RN-006: Status deve respeitar catálogo de visita (tipo V)

### 3.4 Registro de Histórico pelo Representante

```
1. Representante clica Histórico na visita
2. HistoricoRepreList carrega visita atual e resolve lead_id
3. Lista exibe históricos da visita atual e visitas anteriores do mesmo lead
4. HistoricoRepreForm registra:
   - vis_id
   - lead_id
   - motivo_id
   - ocorr_id
   - hist
   - foto_hist opcional
5. Salva e volta para VisitaRepreForm
```

**Regras**:
- RN-007: `historicos` deve persistir `vis_id` e `lead_id`
- RN-008: Histórico é contexto de lead completo, não apenas visita isolada

### 3.5 Criação de Pedido pelo Representante

```
1. Representante abre fluxo PedidosRepreStep1Form → Step5Form
2. Sistema pré-carrega lead e visita quando houver contexto
3. Fluxo coleta:
   - dados cliente
   - dados comerciais
   - pagamento
   - produtos (apenas produtos televendas='N')
   - confirmação final
4. Cria PedidosRepre
5. Pode atualizar visita/lead para vendido
6. Comissão é gerada conforme regra do pedido
```

**Variação**: `pedido_prospec/PedidosRepreStep*Prospec` permite fluxo de prospecção.

### 3.6 Dashboard do Representante

```
1. Representante acessa DashboardRepre
2. Visualiza indicadores:
   - visitas do período
   - pedidos
   - faturamento
   - comissões
   - pendências
```

### 3.7 Comissões do Representante

```
1. Representante acessa ComissoesList
2. Sistema filtra por pedidos do login_repre
3. Exibe comissão, total pago, saldo e previsão
4. Detalhes mostram ComissItem (pagamentos)
```

### 3.8 Despesas Extras

```
1. Representante acessa DespesaExtraForm
2. Informa:
   - data da despesa
   - categoria
   - valor
   - descrição
   - comprovantes
3. Status inicial: PENDENTE
4. Gestor aprova/rejeita no Financeiro
5. Financeiro paga e vincula FinContaPagar
```

---

## 4. Estados Operacionais

### 4.1 Visita do Representante

```
AGENDADO → ANDAMENTO → HISTÓRICO REGISTRADO → VENDIDO/REAGENDADO/PERDIDO
```

### 4.2 Despesa Extra

```
PENDENTE → APROVADA → PAGA
     │
     └──→ REJEITADA
```

### 4.3 Comissão

```
PENDENTE → PARCIAL → PAGO
     │
     └──→ CANCELADO
```

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `VisitaRepreList` | Lista de visitas do representante |
| `VisitaRepreForm` | Detalhe da visita em campo |
| `VisitaRepreStatus` | Alteração de status |
| `VisitaRepreGPS` | Captura/controle GPS |
| `HistoricoRepreList/Form` | Histórico de visitas/leads |
| `PedidosRepreStep1Form`..`Step5Form` | Wizard de pedido |
| `PedidosRepreStep*Prospec` | Wizard de pedido prospecção |
| `PedidosRepreList/Form` | Lista/form de pedidos do representante |
| `DashboardRepre` | Dashboard operacional |
| `ComissoesList` | Comissões próprias |
| `DespesaExtraList/Form` | Despesas extras |

---

## 6. Decisões Arquiteturais

### ADR-001: Representantes como Visão Especializada

**Decisão**: Não criar domínio físico separado; usar views/controllers especializados sobre visitas, pedidos, históricos e financeiro.

**Consequências**:
- Reuso de tabelas principais
- Regras de permissão centralizadas por login_repre
- Fluxos mobile/cortina específicos

### ADR-002: Navegação Protegida em Cortina

**Decisão**: Fluxos internos sempre usam `adianti_target_container=adianti_right_panel`.

**Consequências**:
- Evita fechamento da lista base
- Mantém contexto da visita
- Reduz telas vazias após retorno

### ADR-003: Histórico por Lead, não só por Visita

**Decisão**: HistoricoRepreList usa a visita atual para resolver lead_id e exibe histórico completo do lead.

**Consequências**:
- Representante vê contexto completo do cliente
- Exige `historicos.lead_id` preenchido

### ADR-004: GPS como Evidência de Campo

**Decisão**: Ações relevantes de visita capturam GPS.

**Consequências**:
- Base para KM rodado
- Auditoria de presença
- Tratamento de permissão negada

### ADR-005: Produtos de Representante Separados de Televendas

**Decisão**: Fluxo do representante usa produtos com `televendas='N'`.

**Consequências**:
- Evita mistura de tabelas/preços do canal Televendas
- Validação no Step4 do pedido

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `VisitaRepreList/Form/Status/GPS` | `app/control/representantes/visitas/` |
| `HistoricoRepreList/Form` | `app/control/representantes/historicos/` |
| `GpsPermissionGate` | `app/control/representantes/historicos/` |
| `PedidosRepreList/Form` | `app/control/representantes/pedidos/` |
| `PedidosRepreStep1Form`..`Step5Form` | `app/control/representantes/pedidos/` |
| `PedidosRepreStep*Prospec` | `app/control/representantes/pedidos/pedido_prospec/` |
| `DashboardRepre` | `app/control/representantes/pedidos/` |
| `ComissoesList` | `app/control/representantes/comissoes/` |
| `DespesaExtraList/Form` | `app/control/representantes/despesas/` |

### 7.2 Models/Helpers

| Artefato | Função |
|----------|--------|
| `VisitaRepre` | Visão/modelo de visitas para representante |
| `HistoricoRepre` | Histórico de visita/lead |
| `PedidosRepre` | Pedido do representante |
| `Comissoes` / `ComissItem` | Comissão e pagamentos |
| `RepreDespesaExtra` | Despesas extras |
| `RepreDespesaCateg` | Categorias de despesa |
| `RepreDespesaComprovante` | Comprovantes |
| `PedidosRepreNavigationHelper` | Navegação contextual |

### 7.3 Docs

| Documento | Descrição |
|-----------|-----------|
| `docs/representantes/FLUXO_ENTRE_FORMS_LIST.md` | Regras de navegação entre lista, visita, histórico e status |
| `docs/categorizacao/leads/PLANO_VISAO_REPRESENTANTE_POR_STATUS.md` | Plano de visão por status |

---

## 8. Segurança

- Representante vê apenas dados vinculados ao próprio `login_repre`
- Edição de pedido limitada ao próprio pedido e status permitido
- Comissões em modo consulta para representante
- Despesas só podem ser editadas antes de aprovação
- GPS e histórico são auditáveis

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
