# SPEC_INDEX.md — Índice Mestre do Ecossistema SMID

> Este documento é o ponto de entrada obrigatório para qualquer desenvolvedor ou agente de IA que queira compreender, reimplementar ou evoluir o SMID em qualquer linguagem ou plataforma.

> ## ⚑ Regra de Ouro do Schema
>
> - **Para projeto novo (greenfield)**: schema vem dos SPECs + `app/database/*.sql` + `app/database/scripts/smid padrao/`.
> - **Para manutenção do SMID atual** ou **importação de dados de cliente existente**: use o restante de `app/database/scripts/` conforme `app/database/scripts/README_SCRIPTS.md`.
>
> Detalhes em `SPEC_DATABASE.md` (seções 1.3, 2 e 9).

---

## 1. O que é o SMID

**SMID** é um CRM de vendas e gestão comercial com:
- Captação e qualificação de leads
- Agendamento e execução de visitas em campo
- Pedidos por representantes e televendas
- Financeiro, comissões, reembolso de KM e suporte pós-venda
- Relatórios analíticos e dashboards operacionais
- Comunicação interna, tarefas, metas e administração

---

## 2. Mapa de Domínios

```
┌────────────────────────────────────────────────────────────────────────┐
│                         CAMADAS FUNDACIONAIS                           │
│  SPEC_DATABASE  ·  SPEC_UX_UI  ·  SPEC_REST_API                        │
└────────────────────────────┬───────────────────────────────────────────┘
                             │
┌────────────────────────────▼───────────────────────────────────────────┐
│                          CAMADA ADMIN / PLATAFORMA                     │
│  SPEC_ADMIN  ·  SPEC_LOG  ·  SPEC_COMMUNICATION  ·  SPEC_TAREFAS       │
└────────────────────────────┬───────────────────────────────────────────┘
                             │
┌────────────────────────────▼───────────────────────────────────────────┐
│                         JORNADA COMERCIAL CORE                         │
│                                                                         │
│  LEADS ──► VISITAS ──► HISTÓRICOS ──► PEDIDOS                          │
│  SPEC_LEADS · SPEC_VISITAS · SPEC_HISTORICOS · SPEC_PEDIDOS            │
└──────────────────┬────────────────────────────┬────────────────────────┘
                   │                            │
     ┌─────────────▼──────────┐    ┌────────────▼────────────────┐
     │    CANAL EXTERNO       │    │   CANAL REMOTO              │
     │  SPEC_REPRESENTANTES   │    │   SPEC_TELEVENDAS           │
     └────────────────────────┘    └─────────────────────────────┘
                   │                            │
     ┌─────────────▼────────────────────────────▼──────────────────────┐
     │                      CADEIA PÓS-VENDA                           │
     │  SPEC_PRODUTOS · SPEC_COMPRAS · SPEC_SUPORTE                    │
     └──────────────────────────┬──────────────────────────────────────┘
                                │
     ┌──────────────────────────▼──────────────────────────────────────┐
     │                    FINANCEIRO / INCENTIVOS                      │
     │  SPEC_FINANCEIRO · SPEC_COMISSOES · SPEC_KM · SPEC_METAS       │
     └──────────────────────────┬──────────────────────────────────────┘
                                │
     ┌──────────────────────────▼──────────────────────────────────────┐
     │                       ANÁLISE E DADOS                           │
     │  SPEC_RELATORIOS                                                 │
     └─────────────────────────────────────────────────────────────────┘
                                │
     ┌──────────────────────────▼──────────────────────────────────────┐
     │              INTEGRAÇÕES EXTERNAS E JOBS                        │
     │  SPEC_INTEGRACOES_JOBS                                          │
     └─────────────────────────────────────────────────────────────────┘
```

---

## 3. Catálogo de SPECs

### 3.1 Camadas Fundacionais

| SPEC | Domínio | Dependências-chave |
|------|---------|-------------------|
| `SPEC_DATABASE.md` | Multi-banco, migrações, convenções de schema | base de tudo |
| `SPEC_UX_UI.md` | Padrões de UI, cortina, mobile-premium, busca rápida | base de todas as telas |
| `SPEC_REST_API.md` | Contrato REST, JWT, webhooks | camada de integração |

### 3.2 Domínios

| SPEC | Domínio | Dependências-chave |
|------|---------|-------------------|
| `SPEC_ADMIN.md` | Autenticação, usuários, grupos, papéis, programas, unidades | `permission` (banco) |
| `SPEC_LOG.md` | Auditoria, logs técnicos e de negócio | `log` (banco), todos |
| `SPEC_COMMUNICATION.md` | Mensagens, notificações, feed, wiki, drive, agenda | `communication` (banco) |
| `SPEC_TAREFAS.md` | Lembretes operacionais com notificação em tempo real | `smid`, `Leads` |
| `SPEC_LEADS.md` | Captação, qualificação, gestão e espelhamento de status | `smid`, `permission` |
| `SPEC_VISITAS.md` | Agendamento, execução e GPS de visitas | `smid`, `Leads` |
| `SPEC_HISTORICOS.md` | Registro de ocorridos e fotos por visita/lead | `smid`, `Visitas` |
| `SPEC_PEDIDOS.md` | Pedidos de venda com fluxo wizard | `smid`, `Leads`, `Produtos` |
| `SPEC_REPRESENTANTES.md` | Visão operacional do representante | `smid`, `Visitas`, `Pedidos`, `KM` |
| `SPEC_TELEVENDAS.md` | Canal remoto, fila, orçamentos, pedidos tele | `smid`, `Leads`, `Produtos` |
| `SPEC_PRODUTOS.md` | Catálogo de produtos por canal | `smid` |
| `SPEC_COMPRAS.md` | Aquisição vinculada a pedidos | `smid`, `Pedidos` |
| `SPEC_SUPORTE.md` | Atendimento pós-venda | `smid`, `Pedidos` |
| `SPEC_FINANCEIRO.md` | Contas a pagar/receber, categorias, extrato | `smid`, todos |
| `SPEC_COMISSOES.md` | Comissões por pedido, pagamentos parciais | `smid`, `Pedidos` |
| `SPEC_KM.md` | Reembolso de KM por GPS | `smid`, `Visitas`, `Financeiro` |
| `SPEC_METAS.md` | Objetivos comerciais flexíveis | `smid`, `permission` |
| `SPEC_RELATORIOS.md` | Dashboards e relatórios analíticos | `smid`, todos |
| `SPEC_INTEGRACOES_JOBS.md` | Integrações externas e jobs agendados | todos |
| `SPEC_PHONE_SEARCH_BOX.md` | Componente de busca rápida por telefone | `smid`, `Leads`, `SPEC_REST_API` |

---

## 4. Modelo de Dados Multi-Banco

O SMID usa **quatro bases de dados** com aliases canônicos:

| Alias | Conteúdo |
|-------|----------|
| `smid` | Dados de negócio (leads, visitas, pedidos, financeiro…) |
| `permission` | Usuários, grupos, papéis, programas, unidades |
| `log` | Logs de acesso, requisição, SQL, auditoria |
| `communication` | Mensagens, posts, wiki, documentos, agenda |

**Regra**: serviços que cruzam bases devem abrir/fechar transações por alias explicitamente. Nunca assumir conexão corrente.

---

## 5. Identidades dos Atores

| Ator | Grupo | Papel operacional |
|------|-------|-------------------|
| Admin | 1 | Configuração, permissões, diagnóstico |
| Atendente | 2 | Recepção de leads, agendamento de visitas |
| Representante | 4 | Visitas em campo, pedidos, GPS |
| Supervisor/Gestor | 5 | Relatórios, metas, aprovação de despesas |
| Operador Televendas | grupo tele | Fila, orçamentos, pedidos remotos |

---

## 6. Fluxo Comercial Canônico

```
Mídia externa / WhatsApp / 3C+ / GoTo
        │
        ▼
    LEAD criado
        │
        ▼
  Atendente qualifica e agenda
        │
        ▼
  VISITA criada (status L ↔ V sincronizados)
        │
        ▼
  Representante executa + GPS + Histórico
        │
      ┌─┴──────────────────────────────┐
      ▼                                ▼
  PEDIDO (representante)         PEDIDO (televendas via orçamento)
      │                                │
      ▼                                ▼
  COMISSÃO gerada           COMISSÃO gerada
  COMPRA vinculada          COMPRA vinculada
  SUPORTE pós-venda
      │
      ▼
  FINANCEIRO (contas a receber / pagar / KM)
      │
      ▼
  RELATÓRIOS / METAS
```

---

## 7. Contratos de Interface Entre Domínios

| Produtor | Consumidor | Dado produzido |
|----------|------------|----------------|
| Lead | Visita | `lead_id` obrigatório |
| Lead | Tarefa | `lead_id` opcional |
| Visita | Histórico | `vis_id` + `lead_id` obrigatórios |
| Visita | KM | `vis_id`, `dt_visita`, `hr_visita`, `login_repre` |
| Visita | `LeadVisitStatusSyncService` | `status_id` → sincroniza `leads.status_id` |
| Pedido | Comissão | `ped_id`, `login_repre`, `total_ped` |
| Pedido | Financeiro | `ped_id`, `valor`, `cliente_nome`, `cliente_cpf` |
| Pedido | Compras | `ped_id` |
| Pedido | Suporte | `ped_id` |
| Comissão | Financeiro | `comissao_id` → `fin_contas_pagar` |
| KmLote | Financeiro | `km_reembolso_lote_id` → `fin_contas_pagar` |
| DespesaExtra | Financeiro | `repre_despesa_extra_id` → `fin_contas_pagar` |
| Evolution webhook | Lead | Ingresso assíncrono via `EvolutionLeadIngressService` |
| 3C+/GoTo | Lead/Televendas | Sincronização via sync services |
| Agendamento (SystemSchedule) | Jobs | Execução via `SystemScheduleService` |

---

## 8. Invariantes de Negócio Críticos (não quebrar)

| ID | Invariante |
|----|------------|
| INV-001 | Lead nunca fica sem status |
| INV-002 | Status do lead espelha a última visita válida para sincronismo |
| INV-003 | Ao excluir visita, lead retorna ao último status do tipo L |
| INV-004 | GPS válido é pré-requisito obrigatório para reembolso de KM |
| INV-005 | Canal Televendas/Representante é exclusivo por produto (`S/N`) |
| INV-006 | Todo pedido Televendas exige `lead_id` |
| INV-007 | Comissão cancelada não entra em saldo a pagar |
| INV-008 | Suporte exige `ped_id` vinculado |
| INV-009 | Lançamento automático financeiro deve ser idempotente |
| INV-010 | Histórico deve persistir `vis_id` e `lead_id` |
| INV-011 | Relatórios devem usar categorias funcionais, nunca IDs fixos de status |
| INV-012 | Duplicidade de lead: bloqueio técnico em 20min, negócio em 6h |
| INV-013 | Notificação de tarefa só vai para o criador/responsável |
| INV-014 | Permissões validadas no backend; menu é só conveniência |
| INV-015 | Alterações de dado auditável são registradas por `SystemChangeLog` |

---

## 9. Janelas Temporais e Thresholds

| Contexto | Valor | Significado |
|----------|-------|-------------|
| Duplicidade técnica | 20 minutos | Mesmo telefone no período → bloqueia criação |
| Duplicidade negócio | 6 horas | Mesmo telefone → registra como duplicado |
| Urgência lead | 3 horas | Indicador visual de lead quente |
| Urgência crítica | 5 horas | Indicador visual de lead crítico |
| Kanban leads | 400 cards | Limite MAX_CARDS para performance |
| KM GPS accuracy | 200m | Precisão máxima aceita |
| KM distância lead | 200m | Distância máxima GPS → endereço |
| Notificação tarefa | 5 min antes | Lembrete preventivo |
| Notificação tarefa | 1 min | Janela "na hora" |
| Renotificação atrasada | 15 min | Intervalo entre alertas de atraso |
| Dashboard visitas | 60s | Auto-refresh do dia seguinte |
| Polling notif. tarefas | 60s | Intervalo de polling frontend |

---

## 10. Estrutura de Pastas de Referência

```
app/
  config/
    app_database.php        ← aliases smid/permission/log/communication
  database/                 ← SCHEMA BASE (útil para greenfield)
    permission.sql          ← schema base permission
    permission-update.sql   ← deltas históricos de permission
    communication.sql       ← schema base communication
    communication-update.sql← deltas históricos de communication
    log.sql                 ← schema base log
    scripts/                ← REGISTRO HISTÓRICO de migração do legado
      smid padrao/          ← schema base smid (útil para greenfield)
      1 - unidades..13 - televendas ← 13 categorias históricas;
                              úteis apenas para manutenção do SMID
                              atual ou importação de dados de clientes
      README_SCRIPTS.md     ← guia para manutenção do SMID atual
  control/                  ← controllers (sem SQL)
    admin/                  ← SPEC_ADMIN
    communication/          ← SPEC_COMMUNICATION
    compras/                ← SPEC_COMPRAS
    financeiro/             ← SPEC_FINANCEIRO
    km/                     ← SPEC_KM
    leads/                  ← SPEC_LEADS
    log/                    ← SPEC_LOG
    metas/                  ← SPEC_METAS
    pedidos/                ← SPEC_PEDIDOS
    produtos/               ← SPEC_PRODUTOS
    relatorios/             ← SPEC_RELATORIOS
    representantes/         ← SPEC_REPRESENTANTES
    suportes/               ← SPEC_SUPORTE
    tarefas/                ← SPEC_TAREFAS
    televendas/             ← SPEC_TELEVENDAS
    PhoneSearchBox.php      ← SPEC_PHONE_SEARCH_BOX
  model/
    admin/                  ← SystemUser, Group, Role, Program, Unit
    communication/          ← mensagens, posts, wiki, drive
    log/                    ← SystemChangeLog, AccessLog, RequestLog
    report/                 ← services de relatório
  service/
    auth/                   ← autenticação, sessão, LDAP
    integration/            ← Evolution, 3C+, GoTo, WhatsApp
    km/                     ← cálculo, GPS, geo, reembolso
    log/                    ← serviços de log
    tarefas/                ← notificação de tarefas
    visitas/                ← sync, permissão, status, validação
    financeiro/             ← contas, notificações
  lib/smid/                 ← helpers canônicos: LeadStatusHelper,
                               PedStatusHelper, MetaService…
docs/
  specs/                    ← todos os SPECs (este diretório)
  handoff/                  ← handoffs históricos (ADRs implicitías)
  upgrade/SMID_PERSONALIZACOES.md ← arquivos permitidos/proibidos
  mobile/                   ← análise e requisitos mobile
  BUSCA_RAPIDA.md           ← padrão oficial de busca em listas
  TMultiCombo.md            ← padrão de MultiCombo
  METAS.md                  ← padrão do módulo de metas
  PERMISSAO_DONO.md         ← regra do usuário dono
  televendas/               ← docs específicos de Televendas
  representantes/           ← docs específicos de Representantes
  km/                       ← docs específicos de KM
  leads/                    ← docs específicos de Leads
  notificacoes/             ← docs de notificações (tarefas, financeiras)
  financeiro/               ← docs específicos de Financeiro
  integracao/               ← docs específicos de integrações
  regras_granulares/        ← docs de regras condicionais
  categorizacao/            ← docs de categorias funcionais de status
  deploy/                   ← instruções de deploy
```

---

## 11. Guia de Reimplementação em Outra Linguagem

Para reimplementar o SMID em Go, Java, Node.js, Python ou qualquer stack:

1. **Fundações**: `SPEC_DATABASE` → `SPEC_UX_UI` → `SPEC_REST_API`
2. **Plataforma**: `SPEC_ADMIN` → `SPEC_LOG` → `SPEC_COMMUNICATION`
3. **Jornada core**: `SPEC_LEADS` → `SPEC_VISITAS` → `SPEC_HISTORICOS` → `SPEC_PEDIDOS`
4. **Canais**: `SPEC_REPRESENTANTES` + `SPEC_TELEVENDAS` + `SPEC_PRODUTOS`
5. **Cadeia pós-venda**: `SPEC_COMPRAS` + `SPEC_SUPORTE`
6. **Financeiro/Incentivos**: `SPEC_FINANCEIRO` + `SPEC_COMISSOES` + `SPEC_KM` + `SPEC_METAS`
7. **Análise**: `SPEC_RELATORIOS`
8. **Integrações**: `SPEC_INTEGRACOES_JOBS`
9. **Widgets/UX**: `SPEC_PHONE_SEARCH_BOX` + `SPEC_TAREFAS`

**Para cada domínio**:
- Implemente as entidades da seção 1.2
- Respeite os relacionamentos da seção 1.3
- Implemente os fluxos da seção 3.x em ordem
- Implemente a máquina de estados da seção 4.x
- Siga as ADRs como contratos arquiteturais
- Nunca quebre os invariantes do item 8 deste índice
- Valide responsividade nos breakpoints de `SPEC_UX_UI`
- Exponha o recurso via `SPEC_REST_API` quando aplicável

---

## 12. Fontes Primárias por Categoria

| Categoria | Localização | Como usar |
|-----------|-------------|-----------|
| SPECs (este diretório) | `docs/specs/SPEC_*.md` | Fonte primária de reimplementação |
| Schema base | `app/database/*.sql` + `scripts/smid padrao/` | DDL de referência para greenfield |
| Migrações históricas | `app/database/scripts/1..13` | Registro histórico; não replicar em greenfield |
| Models | `app/model/` | Contratos de entidade e relacionamentos |
| Services | `app/service/` | Lógica reutilizável e integrações |
| Helpers SMID | `app/lib/smid/` | Padrões transversais (status, navegação, meta) |
| Controllers | `app/control/` | Exemplos de fluxo/UI — não reimplementar verbatim |
| Handoffs | `docs/handoff/` | Decisões e planos históricos (ADRs implícitas) |
| Docs temáticos | `docs/<tema>/` | Contexto de negócio por módulo |
| Restrições de arquivos | `docs/upgrade/SMID_PERSONALIZACOES.md` | Arquivos permitidos/proibidos |

---

## 13. Checklist de Reimplementação

Antes de considerar uma reimplementação do SMID concluída, validar:

### 13.1 Dados
- [ ] Quatro conexões lógicas (`smid`, `permission`, `log`, `communication`) funcionando
- [ ] Schema derivado dos SPECs + arquivos base de `app/database/`
- [ ] Tabelas de categorização funcional (`lead_status_categoria`, `ped_status_categoria`) povoadas
- [ ] Views críticas criadas (`v_calendario_visitas`, fila televendas)
- [ ] Convenções transversais aplicadas (timestamps `criado_em`/`alterado_em`/`excluido_em`, `unidd_id`)

### 13.2 Domínios
- [ ] Todos os 21 SPECs de domínio implementados
- [ ] Invariantes INV-001 a INV-015 validados por testes
- [ ] Sincronização Lead ↔ Visita funcionando em todas as transições

### 13.3 Plataforma
- [ ] Autenticação + multi-unidade + regras granulares
- [ ] Auditoria via equivalente a `SystemChangeLogTrait`
- [ ] Notificações internas e jobs agendados
- [ ] Logs de acesso, requisição e SQL

### 13.4 Integrações
- [ ] Webhook Evolution com ingresso de leads + duplicidade
- [ ] Sync 3C+ e GoTo Connect
- [ ] Notificações WhatsApp para equipe
- [ ] LDAP (opcional conforme cliente)

### 13.5 UI/UX
- [ ] Todos os breakpoints mobile validados (320, 375, 414, 428, 768, desktop)
- [ ] Cortina lateral funcionando em forms invocados por listas
- [ ] Busca rápida seguindo árvore de decisão
- [ ] `TDateRangeField` com presets backend-only
- [ ] Tema dark/light funcionando

### 13.6 API
- [ ] JWT emitido/validado com claims obrigatórios
- [ ] Endpoints principais por domínio expostos
- [ ] Webhooks idempotentes por chave externa
- [ ] Rate limit e auditoria em endpoints sensíveis

### 13.7 Relatórios
- [ ] Aproveitamento, comparativo, ligações do dia, evolução de vendas, visitas do dia seguinte
- [ ] Uso de categorias funcionais de status (nunca IDs fixos)
- [ ] Filtros respeitando escopo de unidade

---

**Versão**: 2.0
**Data**: 2026-05-12
**Domínios cobertos**: 23 SPECs (3 fundacionais + 20 de domínio)
**Status**: Documento vivo — atualizar sempre que um SPEC for criado ou modificado
