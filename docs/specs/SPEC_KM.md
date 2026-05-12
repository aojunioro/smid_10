# SPEC_KM.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **KM Rodado** calcula e controla reembolsos de deslocamento dos representantes comerciais. Ele reconstrói a rota diária com base nas visitas confirmadas por GPS, calcula distância real via provedor de mapas, aplica valor do KM por vigência e gera lotes de reembolso integráveis ao financeiro.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **KmConfig** | Configurações globais de KM/GPS | 1 |
| **KmValorKmVigencia** | Valor do KM por período de vigência | 0..N |
| **KmReembolsoLote** | Lote de reembolso por representante/período | 0..N |
| **KmReembolsoDia** | Cálculo diário dentro do lote | 0..N por lote |
| **KmReembolsoTrecho** | Trecho individual da rota diária | 0..N por dia |
| **VisitaCheckinEvent** | Evento GPS que valida presença | 0..N por visita |
| **VisStatusGpsEventMap** | Mapeia status de visita para evento GPS esperado | catálogo |
| **FinContaPagar** | Pagamento financeiro do lote | 0..1 por lote |

### 1.3 Relacionamentos

```
Representante 1 ──── 0..N KmReembolsoLote
KmReembolsoLote 1 ──── 0..N KmReembolsoDia
KmReembolsoDia 1 ──── 0..N KmReembolsoTrecho
Visita 1 ──── 0..N VisitaCheckinEvent
KmReembolsoTrecho N ──── 0..1 VisitaCheckinEvent
KmReembolsoLote 0..1 ──── FinContaPagar
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| VISITAS | Entrada | Visitas do representante ordenadas por data/hora |
| HISTORICOS | Leitura | Evidência textual, mas não valida presença |
| GPS | Entrada | Check-in/out obrigatório para reembolso |
| FINANCEIRO | Saída | Lote pago gera conta a pagar |
| PERMISSION | Leitura | Endereço residencial do representante |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Lote de Reembolso** | Apuração de KM por representante e período |
| **Dia de Reembolso** | Cálculo de um dia dentro do lote |
| **Trecho** | Segmento da rota (origem → destino) |
| **Ponto Inicial** | Residência do representante ou endereço customizado |
| **Retorno ao Início** | Último trecho obrigatório de volta ao ponto inicial |
| **Presença Confirmada** | Visita com evento GPS válido |
| **Accuracy** | Precisão em metros retornada pelo GPS |
| **Valor KM Vigente** | Valor monetário aplicado na data de referência |

---

## 3. Fluxos Principais

### 3.1 Configuração de KM

```
1. Gestor acessa KmConfigForm
2. Define parâmetros:
   - gps_accuracy_max_m (default 200m)
   - gps_distancia_max_lead_m (default 200m)
   - provider de mapas
   - flags de cache/validação
3. Configuração é usada por KmConfigService
```

### 3.2 Cadastro de Valor por Vigência

```
1. Gestor acessa KmValorKmVigenciaForm
2. Informa:
   - data início
   - data fim opcional
   - valor_km
3. Sistema valida sobreposição de vigências
4. Cálculo usa valor vigente na data da visita/dia
```

### 3.3 Captura de Presença GPS

```
1. Representante altera status/realiza check-in em visita
2. Browser solicita geolocalização
3. VisitaCheckinEvent registra:
   - vis_id
   - login_repre
   - tipo CHECK_IN/CHECK_OUT
   - lat/lng
   - accuracy_m
   - permission S/N
   - erro_codigo / erro_msg se falha
   - capturado_em
4. KmGpsValidationService valida presença:
   - permission = S
   - accuracy_m <= config
   - distância GPS → endereço do lead <= config
```

**Regra Canônica**: Sem GPS válido, não há reembolso. Histórico textual não substitui GPS.

### 3.4 Geração de Lote

```
1. Gestor acessa KmReembolsoLoteList
2. Cria lote para login_repre + dt_inicio/dt_fim
3. KmCalculationService busca visitas do período
4. Filtra somente visitas com GPS válido
5. Agrupa por dia
6. Para cada dia:
   - resolve ponto inicial (residência ou custom)
   - ordena visitas por hr_visita e id
   - monta rota: início → visita1 → visita2 → ... → retorno
   - calcula trechos via KmGeoService
   - aplica valor_km vigente
7. Persiste KmReembolsoDia e KmReembolsoTrecho
8. Atualiza totais do lote:
   - km_total
   - valor_km_total
   - valor_total
```

### 3.5 Revisão e Pagamento

```
1. Gestor abre KmReembolsoLoteDetalhe
2. Revisa dias e trechos
3. Status permanece EM_APURACAO até aprovação
4. Financeiro abre KmReembolsoPagarForm
5. Ao pagar:
   - status_pagamento = PAGO
   - pago_por / pago_em preenchidos
   - cria/vincula FinContaPagar
```

---

## 4. Estados

### 4.1 Lote de Reembolso

```
EM_APURACAO → PAGO
      │
      └──→ CANCELADO
```

### 4.2 Dia de Reembolso

```
PENDENTE → OK
    │       │
    │       └── PARCIAL (algum trecho com erro)
    └── ERRO
```

### 4.3 Critérios GPS

| Critério | Regra |
|----------|-------|
| Permissão | `permission = 'S'` |
| Precisão | `accuracy_m <= gps_accuracy_max_m` |
| Distância do lead | `gps_dist_to_lead_m <= gps_distancia_max_lead_m` |
| Evento | Deve corresponder ao status mapeado em `VisStatusGpsEventMap` |

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `KmConfigForm` | Configurações gerais |
| `KmValorKmVigenciaList/Form` | Valor do KM por período |
| `KmReembolsoLoteList` | Lotes por representante/período |
| `KmReembolsoLoteDetalhe` | Dias e trechos do lote |
| `KmReembolsoPagarForm` | Pagamento do lote |

---

## 6. Decisões Arquiteturais

### ADR-001: GPS como Única Fonte de Presença

**Decisão**: Somente GPS válido confirma deslocamento.

**Consequências**:
- Anti-fraude robusto
- Sem fallback manual/histórico
- Requer captura mobile confiável

### ADR-002: Rota Sequencial por Dia

**Decisão**: Reconstruir rota diária por `dt_visita + hr_visita`.

**Consequências**:
- Distância calculada pelo trajeto real esperado
- Empates resolvidos por id
- Retorno ao início obrigatório

### ADR-003: Valor KM por Vigência

**Decisão**: Valor global por período, não por representante.

**Consequências**:
- Simplicidade operacional
- Histórico preservado quando valor muda

### ADR-004: Trechos Persistidos

**Decisão**: Persistir cada trecho calculado.

**Consequências**:
- Auditoria detalhada
- Reprocessamento controlado
- Cache e provider_status visíveis

### ADR-005: Integração Financeira por Lote

**Decisão**: FinContaPagar referencia `km_reembolso_lote_id`.

**Consequências**:
- Pagamento rastreável
- Evita duplicidade financeira

---

## 7. Referência à Implementação Atual

### 7.1 Controllers

| Controller | Localização |
|------------|-------------|
| `KmConfigForm` | `app/control/km/` |
| `KmValorKmVigenciaList/Form` | `app/control/km/` |
| `KmReembolsoLoteList` | `app/control/km/` |
| `KmReembolsoLoteDetalhe` | `app/control/km/` |
| `KmReembolsoPagarForm` | `app/control/km/` |

### 7.2 Models

| Model | Tabela |
|-------|--------|
| `KmConfig` | `km_config` |
| `KmValorKmVigencia` | `km_valor_km_vigencia` |
| `KmReembolsoLote` | `km_reembolso_lote` |
| `KmReembolsoDia` | `km_reembolso_dia` |
| `KmReembolsoTrecho` | `km_reembolso_trecho` |
| `VisStatusGpsEventMap` | mapeamento status/GPS |
| `VisitaCheckinEvent` | `visita_checkin_event` |

### 7.3 Services

| Service | Responsabilidade |
|---------|------------------|
| `KmCalculationService` | Geração de lote/dias/trechos |
| `KmConfigService` | Leitura de config |
| `KmGeoService` | Geocoding/distância/provider/cache |
| `KmGpsValidationService` | Validação anti-fraude GPS |
| `KmReembolsoDataService` | Dados para tela/listagem |

---

## 8. Segurança

- Eventos GPS são auditáveis e não devem ser alterados manualmente
- Pagamento restrito ao financeiro/gestão
- Reprocessamento de lote deve preservar histórico ou registrar auditoria
- Endereço residencial do representante é dado sensível
- Provider de mapas não deve expor chaves no frontend sem controle

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
