# SPEC_HISTORICOS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Históricos** representa o registro temporal de eventos relacionados às visitas e leads — o que aconteceu durante ou após o contato. É o "log de eventos comerciais" que documenta motivos, ocorridos, observações e fotos. Subsidia auditoria, análise de gargalos e validação de visitas efetivas.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **HistoricoRepre** | Registro de evento histórico de visita | 0..N por visita |
| **HistoricoMotivo** | Catálogo de motivos (com cor) | catálogo |
| **HistoricoOcorrido** | Catálogo de ocorridos (com cor + responsável) | catálogo |

### 1.3 Relacionamentos

```
Visita 1 ──── 0..N HistoricoRepre
Lead   1 ──── 0..N HistoricoRepre (vínculo direto via lead_id)
HistoricoRepre N ──── 1 HistoricoMotivo
HistoricoRepre N ──── 1 HistoricoOcorrido
HistoricoRepre N ──── 1 SystemUser (login)
```

### 1.4 Integrações com outros Módulos

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| VISITAS | Bidirecional | Histórico vinculado à visita; atualiza `visita.hist_feito` |
| LEADS | Saída | `historicos.lead_id` permite consulta direta sem JOIN |
| PEDIDOS | Saída | Ocorrido pode indicar conversão em pedido |
| RELATÓRIOS | Leitura | Análise de motivos, eficácia por representante |
| MOBILE | Entrada | Representante registra histórico em campo |

---

## 2. Glossário de Negócio

### 2.1 Termos do Domínio

| Termo | Definição |
|-------|-----------|
| **Histórico** | Registro de evento ou ocorrido relacionado à visita |
| **Motivo** | Razão/categoria do registro (ex: Cliente não atendeu, Reagendou) |
| **Ocorrido** | Detalhamento mais granular do que aconteceu |
| **Responsável** | Quem deve tratar o ocorrido: `callcenter` ou `comercial` |
| **Foto Hist** | Foto opcional anexada ao histórico (ex: comprovante, recibo) |
| **Hist Feito** | Flag em `visitas.hist_feito` indicando que histórico foi preenchido |

### 2.2 Tipos de Responsável

| Código | Descrição | Quem Trata |
|--------|-----------|------------|
| `callcenter` | CallCenter | Atendente / supervisão de callcenter |
| `comercial` | Comercial | Representantes / gestão comercial |

### 2.3 Convenções de Cor

Tanto `HistoricoMotivo` quanto `HistoricoOcorrido` possuem campo `cor` (hex) usado em badges visuais para diferenciação rápida na timeline e listagens.

---

## 3. Fluxos Principais

### 3.1 Registro de Histórico Durante Visita

**Ator**: Representante
**Pré-condições**: Visita em andamento; permissão GPS validada

```
1. Representante abre HistoricoRepreForm (mobile)
2. GpsPermissionGate verifica permissão GPS (registra evento)
3. Sistema pré-preenche:
   - vis_id (visita atual)
   - lead_id (de visitas.lead_id)
   - login (representante logado)
   - criado_em (timestamp)
4. Representante seleciona:
   - motivo_id (HistoricoMotivo via combo)
   - ocorr_id (HistoricoOcorrido via combo, filtrado por responsável)
5. Preenche hist (texto livre)
6. Opcionalmente anexa foto_hist
7. Salva
8. Sistema atualiza visita.hist_feito = 'S'
9. Se ocorrido invalidante (regra hist_motivo=2 conforme legado):
   9.1 Visita perde validade para relatórios de visita efetivada
```

**Regras de Negócio**:
- RN-001: lead_id deve ser preenchido a partir de visitas.lead_id (vide handoff)
- RN-002: Foto opcional, mas obrigatória para certos ocorridos (configurável)
- RN-003: Após registrar, visita.hist_feito é atualizada

### 3.2 Registro de Histórico em Tela (Atendente)

**Ator**: Atendente
**Pré-condições**: Acesso ao lead

```
1. Atendente abre LeadsForm
2. Aba Histórico (HistoricosController) é exibida
3. HistoricosGrid mostra timeline cronológica
4. Atendente clica "Adicionar Histórico"
5. HistoricoRepreForm em cortina
6. Preenche dados (mesma lógica)
7. Salva e atualiza grid
```

### 3.3 Listagem e Análise (HistoricoList)

```
1. Gestor acessa Histórico > Listagem
2. Filtros disponíveis:
   - Período
   - Motivo
   - Ocorrido
   - Representante
   - Lead
3. DataGrid exibe registros com:
   - Data/hora
   - Lead (nome + telefone)
   - Visita vinculada
   - Motivo (com cor)
   - Ocorrido (com cor)
   - Texto do histórico
   - Foto (preview/download)
4. Exportação Excel disponível
```

### 3.4 Análise por Representante (HistoricoRepreList)

**Tipo**: Listagem dedicada por representante

```
1. Gestor seleciona representante
2. Sistema exibe histórico completo do representante
3. Filtros: período, motivo, ocorrido
4. Análise de produtividade e padrões
```

### 3.5 Auditoria de lead_id Faltante (HAIFLEX_ORIG)

**Contexto Histórico**: Em ambientes legados (HAIFLEX_ORIG), `historicos.lead_id` ficou vazio sistemicamente.

**Estratégia de Saneamento**:
```
1. Auditoria via SELECT (sem escrita)
2. Identifica registros com lead_id NULL
3. Resolve via JOIN h.vis_id = v.id (visita.lead_id como fonte de verdade)
4. Migration controlada UPDATE historicos SET lead_id = visitas.lead_id
   WHERE historicos.vis_id = visitas.id AND historicos.lead_id IS NULL
5. Validação: zero divergência entre lead_id resolvido e visita.lead_id
```

**Referência**: `docs/handoff/HANDOFF_HISTORICO_SYNC_LEAD_VISITA.md`

---

## 4. Cadastros Auxiliares

### 4.1 Histórico Motivo

**Controllers**: `HistoricoMotivoForm`, `HistoricoMotivoList`

**Estrutura**:
- `motivo` (varchar) — descrição
- `cor` (hex) — cor do badge visual

**Exemplos**: "Cliente não atendeu", "Reagendou", "Confirmou compra", "Cancelou"

### 4.2 Histórico Ocorrido

**Controllers**: `HistoricoOcorridoForm`, `HistoricoOcorridoList`

**Estrutura**:
- `ocorrido` (varchar) — descrição
- `cor` (hex) — cor do badge
- `ordem` (int) — ordenação no combo
- `responsavel` (enum) — `callcenter` ou `comercial`

**Filtragem**: O combo no Form filtra ocorridos pelo responsável apropriado ao perfil logado.

---

## 5. Decisões Arquiteturais (ADRs)

### ADR-001: Vínculo Direto a Lead via lead_id

**Contexto**: Consultas frequentes de histórico por lead exigiam JOIN com visitas.

**Decisão**: Adicionar `historicos.lead_id` como redundância para acesso direto.

**Consequências**:
- Performance melhor em queries de timeline
- Necessidade de saneamento em bases legadas
- `visitas.lead_id` permanece fonte de verdade para reconciliação

### ADR-002: Cor no Cadastro

**Decisão**: Tanto `HistoricoMotivo` quanto `HistoricoOcorrido` têm campo `cor`.

**Consequências**:
- UX visual consistente em badges/timelines
- Customizável por cliente
- Padronização com outros módulos (LeadStatus, Equipe)

### ADR-003: Responsável como Filtro de Combo

**Decisão**: `HistoricoOcorrido.responsavel` filtra opções no Form conforme perfil.

**Consequências**:
- Atendente vê apenas ocorridos de callcenter
- Representante vê apenas ocorridos comerciais
- Reduz erro humano

### ADR-004: Foto Opcional com Storage Controlado

**Decisão**: `foto_hist` armazenada conforme storage do projeto (filesystem/blob).

**Consequências**:
- Comprovante visual
- Cuidado com tamanho (compressão recomendada)

### ADR-005: SystemChangeLog Apenas em HistoricoRepre

**Decisão**: Apenas a entidade principal `HistoricoRepre` usa `SystemChangeLogTrait`. Catálogos (Motivo/Ocorrido) não auditam.

**Consequências**:
- Audit trail focado em registros operacionais
- Catálogos mais leves

---

## 6. Referência à Implementação Atual

### 6.1 Controllers

| Controller | Localização | Propósito |
|------------|-------------|-----------|
| `HistoricoList` | `app/control/historicos/HistoricoList.php` | Listagem geral |
| `HistoricosGrid` | `app/control/historicos/HistoricosGrid.php` | Grid embutido |
| `HistoricoRepreList` | `app/control/representantes/historicos/HistoricoRepreList.php` | Listagem por representante |
| `HistoricoRepreForm` | `app/control/representantes/historicos/HistoricoRepreForm.php` | CRUD de histórico |
| `GpsPermissionGate` | `app/control/representantes/historicos/GpsPermissionGate.php` | Gate de permissão GPS |
| `HistoricoMotivoForm` / `List` | `app/control/historicos/motivo/` | Catálogo de motivos |
| `HistoricoOcorridoForm` / `List` | `app/control/historicos/ocorrido/` | Catálogo de ocorridos |

### 6.2 Models

| Model | Tabela | Auditoria |
|-------|--------|-----------|
| `HistoricoRepre` | `historicos` | SystemChangeLogTrait |
| `HistoricoMotivo` | `hist_motivo` | Não |
| `HistoricoOcorrido` | `hist_ocorrido` | Não |

### 6.3 Vínculos com Outros Módulos

| Origem | Destino | Vínculo |
|--------|---------|---------|
| `LeadsForm` | Aba Histórico | Subcontroller embutido |
| `LeadVisitaDetalhes` | Histórico de visita | Acesso via vis_id |
| `LeadVisitNextDayDashboard` | Métrica de visitas efetivadas | Filtro hist_motivo |

---

## 7. Considerações de Segurança

### 7.1 Permissões

| Role | Pode Criar | Pode Editar | Pode Excluir |
|------|-----------|-------------|--------------|
| Admin | Sempre | Sempre | Sempre |
| Gestor | Sempre | Sempre | Conforme política |
| Atendente | Sim | Próprios | Apenas recém-criados |
| Representante | Sim | Próprios | Não |

### 7.2 GPS e Auditoria

- Cada registro de histórico está vinculado a `vis_id`
- `VisitaCheckinEvent` registra GPS no momento da visita
- Permissão GPS é gateada por `GpsPermissionGate`

### 7.3 Foto Histórica

- Tamanho limitado pelo storage configurado
- Compressão recomendada antes do upload
- Acesso restrito a usuários com permissão no lead

---

## 8. Métricas e Análises

### 8.1 Visita Efetivada

**Critério histórico**: Registro em `historicos` cuja `hist_motivo` não seja invalidante (legado: `hist_motivo=2`).

**Uso**: Dashboard de Evolução de Vendas (T2: efetivação da visita).

**Migração**: A definição evoluiu para usar `v_visitas_validas.dt_visita`, mas o histórico continua como evidência operacional.

### 8.2 Análises de Padrões

- Motivos mais comuns por representante
- Ocorridos por período/unidade
- Taxa de visitas com histórico preenchido
- Tempo entre visita e registro do histórico

---

## 9. Glossário Técnico

| Termo | Significado |
|-------|-------------|
| **historicos.vis_id** | FK obrigatório para visita |
| **historicos.lead_id** | FK redundante (otimização de query) |
| **hist_motivo=2 (legado)** | Motivo que invalida visita para relatórios |
| **HistoricosGrid** | Grid embutido em outros forms |
| **GpsPermissionGate** | Componente de UI/Backend para verificar permissão GPS |
| **foto_hist** | BLOB ou path para foto anexada |

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
