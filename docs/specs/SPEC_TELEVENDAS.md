# SPEC_TELEVENDAS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Televendas** organiza a fila de leads elegíveis para abordagem remota, registra contatos, gera orçamentos com preços especiais e converte orçamentos ou vendas diretas em pedidos próprios do canal. Ele é diretamente alinhado aos domínios Leads, Visitas, Históricos, Produtos e Pedidos.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **TelevendasFila** | View consolidada da fila (`v_televendas_base`) | view |
| **TelevendasContato** | Registro de contato realizado pelo operador | 0..N por lead |
| **TelevendasStatus** | Status do contato Televendas | catálogo |
| **TeleCoeficiente** | Coeficientes comerciais/financeiros do canal | catálogo |
| **ProdOrcamento** | Produto comercial do orçamento Televendas | catálogo |
| **OrcamProdItem** | Itens do orçamento | 0..N por orçamento |
| **Orçamento** | Proposta comercial Televendas | 0..N por lead |
| **PedidosTele** | Pedido gerado pelo canal Televendas | 0..N |
| **ThreeCPlusConfig/Logs** | Integração 3C+ | configuração/auditoria |
| **GotoConnectConfig/Logs** | Integração GoTo Connect | configuração/auditoria |

### 1.3 Relacionamentos

```
Lead 1 ──── 0..N TelevendasContato
Lead 1 ──── 0..N Orçamento
Orçamento 1 ──── 0..N OrcamProdItem
Orçamento 0..1 ──── 1 PedidoTele (conversão)
TelevendasContato 0..1 ──── Orçamento
TelevendasContato 0..1 ──── Pedido
TelevendasContato 0..1 ──── Visita
TelevendasContato 0..1 ──── Histórico
Produto base 1 ──── 0..N ProdOrcamento
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| LEADS | Entrada | Fila consolidada de leads elegíveis |
| VISITAS | Leitura | Usa dados de última visita/status visita |
| HISTORICOS | Leitura | Usa último histórico/motivo |
| PRODUTOS | Leitura | Produtos Televendas (`televendas='S'`) |
| PEDIDOS | Saída | Converte orçamento ou cria pedido direto |
| 3C+ | Entrada | Chamadas, relatórios, leads receptivos |
| GoTo Connect | Entrada | Chamadas e sincronização de leads |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Fila Televendas** | Conjunto de leads elegíveis para contato remoto |
| **Contato** | Registro de interação do operador com o lead |
| **Orçamento** | Proposta comercial antes do pedido |
| **Pedido Tele** | Pedido criado pelo fluxo Televendas, separado de PedidosForm padrão |
| **Temperatura** | Grau de interesse/probabilidade do orçamento |
| **Coeficiente** | Parâmetro financeiro aplicado ao orçamento |
| **Produto Tele** | Produto habilitado para televendas (`produtos.televendas='S'`) |
| **Conversão** | Transformação de orçamento em pedido |

---

## 3. Fluxos Principais

### 3.1 Fila de Televendas

```
1. Operador abre TelevendasList
2. Sistema carrega `v_televendas_base` via TelevendasFila
3. Regra base aplicada automaticamente:
   - unidade elegível OU motivo_id=1
4. Colunas consolidadas:
   - Hora
   - Detalhe
   - Endereço
   - Lead
   - Histórico
   - Resultado
   - Observação
   - login_recep / login_repre
5. Filtros opcionais:
   - Meus contatos
   - Datas de visita (TDateRangeField)
   - Status orçamento
   - Datas orçamento
   - Temperatura
   - Mídia, unidade, status lead, atendente, representante, motivo
```

**Regras**:
- RN-001: Lista abre sem exigir filtros manuais
- RN-002: View consolidada é fonte de performance
- RN-003: Filtro "Meus contatos" usa `televendas_contatos.login`

### 3.2 Registro de Contato

```
1. Operador clica em um item da fila
2. Abre TelevendasForm
3. Form grava exclusivamente em `televendas_contatos`
4. Campos registrados:
   - lead_id
   - login
   - status_id
   - orcam_id
   - ped_id
   - visita_id
   - historico_id
   - midia_id / unidd_id
   - login_recep / login_repre
   - dt_visita / hr_visita
   - cidade / bairro
   - status_visita_id
   - motivo_id
   - observacao
   - criado_em
5. Após salvar, fila atualiza status visual
```

**Regras**:
- RN-004: TelevendasForm não altera Lead diretamente
- RN-005: Todo contato deve ter login do operador

### 3.3 Criação de Orçamento

```
1. Operador abre OrcamForm a partir do contato/lead
2. Seleciona produtos de Televendas (ProdOrcamento)
3. Produtos devem estar vinculados a Produto base com `televendas='S'`
4. Adiciona OrcamProdItem:
   - prod_orcam_id
   - med_id
   - qtdd_item
   - vlr_item
   - vlr_total_item
5. Aplica coeficientes (TeleCoeficiente)
6. Define temperatura e status orçamento
7. Salva orçamento
8. Contato registra orcam_id
```

### 3.4 Conversão de Orçamento em Pedido

```
1. Operador clica em "Converter em Pedido"
2. Fluxo obrigatório: PedidosTeleStep1Form → Step5Form
3. PedidosTeleOrcamentoPreloadHelper carrega dados do orçamento
4. Usuário confirma:
   - lead_id obrigatório
   - dados do cliente
   - forma/condição de pagamento
   - produtos/valores
   - datas comerciais
5. Sistema cria pedido via fluxo Televendas
6. Registra ped_id no orçamento/contato
7. Atualiza status orçamento para convertido
```

**Regras**:
- RN-006: Proibido usar `PedidosForm.php` para inclusão/conversão do Televendas
- RN-007: Conversão sempre via `PedidosTeleStep*`
- RN-008: `lead_id` obrigatório

### 3.5 Pedido Direto sem Orçamento

```
1. Operador decide criar pedido sem orçamento
2. Abre PedidosTeleStep1Form diretamente
3. Sistema exige lead_id
4. Fluxo Step1→Step5 executa validações do canal
5. Cria pedido e registra ped_id no contato
```

### 3.6 Integração 3C+

```
1. Configuração em ThreeCPlusConfig
2. ThreeCPlusAuthService autentica
3. ThreeCPlusCallReportService busca chamadas por período
4. ThreeCPlusLeadSyncService sincroniza leads/contatos
5. ThreeCPlusProcessedCall evita duplicidade
6. ThreeCPlusSyncLog registra execução
7. ThreeCPlusSyncPeriodForm permite sincronização manual
```

### 3.7 Integração GoTo Connect

```
1. GotoConnectOAuthService realiza OAuth
2. GotoConnectCallReportService coleta chamadas
3. GotoLeadSyncService cria/atualiza leads conforme regras
4. GotoConnectSyncLog registra execução
```

---

## 4. Estados

### 4.1 Status Televendas

`TelevendasStatus` possui:
- `stts_tele`
- `ordem`
- `cor`

Exemplos esperados:
- Novo contato
- Em negociação
- Orçamento enviado
- Convertido
- Sem interesse
- Retornar depois

### 4.2 Status Orçamento

`OrcamStatusForm/List` controla status de orçamento.

Fluxo típico:
```
NOVO → ENVIADO → NEGOCIAÇÃO → CONVERTIDO
                  └────────→ PERDIDO / CANCELADO
```

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `DashTelevendas` | Dashboard do canal |
| `TelevendasList` | Fila operacional |
| `TelevendasForm` | Registro de contato |
| `OrcamList` / `OrcamForm` / `OrcamGrid` | Orçamentos |
| `PedidosTeleList` / `PedidosTeleForm` | Pedidos do canal |
| `PedidosTeleStep1Form`..`Step5Form` | Wizard de pedido |
| `TeleProdutosList/Form` | Produtos Televendas |
| `TeleCoeficientesList/Form` | Coeficientes comerciais |
| `TeleStatusList/Form` | Status de contato |

---

## 6. Decisões Arquiteturais

### ADR-001: Fila via View Consolidada

**Decisão**: `TelevendasFila` lê `v_televendas_base`.

**Consequências**:
- Performance adequada para múltiplos joins
- Lista abre sem filtros manuais
- Requer índices nas colunas de filtro

### ADR-002: Contato Não Altera Lead

**Decisão**: `TelevendasForm` grava exclusivamente em `televendas_contatos`.

**Consequências**:
- Histórico de abordagem separado do cadastro do lead
- Evita efeitos colaterais em LeadsForm

### ADR-003: Fluxo de Pedido Próprio

**Decisão**: Televendas não usa `PedidosForm.php`; usa `PedidosTeleStep*`.

**Consequências**:
- Regras do canal isoladas
- Lead obrigatório sempre
- Conversão orçamento→pedido controlada

### ADR-004: Produtos por Canal

**Decisão**: Televendas só usa produtos `televendas='S'`.

**Consequências**:
- Evita mistura com produto de domicílio
- Validação backend obrigatória em Step4

### ADR-005: Integrações como Services Dedicados

**Decisão**: 3C+ e GoTo possuem services/config/log próprios.

**Consequências**:
- Credenciais isoladas e criptografadas
- Sincronização auditável
- Possibilidade de execução manual por período

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `DashTelevendas` | `app/control/televendas/` |
| `TelevendasList/Form` | `app/control/televendas/` |
| `OrcamList/Form/Grid` | `app/control/televendas/orcamentos/` |
| `OrcamStatusForm/List` | `app/control/televendas/orcamentos/` |
| `PedidosTeleList/Form` | `app/control/televendas/pedidos/` |
| `PedidosTeleStep1Form`..`Step5Form` | `app/control/televendas/pedidos/` |
| `TeleProdutosForm/List` | `app/control/televendas/produtos/` |
| `TeleCoeficientesForm/List` | `app/control/televendas/coeficientes/` |
| `TeleStatusForm/List` | `app/control/televendas/status/` |

### 7.2 Models e Services

| Artefato | Função |
|----------|--------|
| `TelevendasFila` | View da fila |
| `TelevendasContato` | Contatos realizados |
| `TelevendasStatus` | Status do contato |
| `TeleCoeficiente` | Coeficientes |
| `ThreeCPlus*` | Integração 3C+ |
| `GotoConnect*` | Integração GoTo |
| `PedidosTeleOrcamentoPreloadHelper` | Pré-carrega orçamento no pedido |
| `PedidosTeleNavigationHelper` | Navegação contextual |

---

## 8. Segurança

- Acesso restrito ao departamento Televendas
- Credenciais 3C+/GoTo criptografadas por services de crypto
- Logs de sincronização obrigatórios
- Nenhum fluxo Televendas cria pedido sem `lead_id`
- Produtos validados no backend por canal

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
