# SPEC_TAREFAS.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Tarefas** gerencia lembretes e atividades operacionais atribuídas a usuários do SMID. Uma tarefa possui descrição, data, hora, status, responsável (`login`) e vínculo opcional com Lead. O módulo oferece listagem, formulário, Kanban e notificações globais em tempo real para tarefas próximas, no horário ou atrasadas.

### 1.2 Entidades Principais

| Entidade | Descrição | Cardinalidade |
|----------|-----------|---------------|
| **Tarefa** | Registro principal da atividade | 0..N por usuário |
| **TarefaNotificationService** | Serviço de notificação e priorização | service |
| **TarefaNotificationCheck** | Endpoint/controller de checagem | service/controller |
| **Lead** | Vínculo opcional para contexto comercial | 0..1 |
| **SystemUser** | Responsável pela tarefa (`login`) | 1 |

### 1.3 Relacionamentos

```
SystemUser.login 1 ──── 0..N Tarefa.login
Lead 1 ──── 0..N Tarefa.lead_id
TarefaNotificationService ──── lê tarefas pendentes por login atual
```

### 1.4 Integrações

| Módulo | Tipo | Descrição |
|--------|------|-----------|
| LEADS | Vínculo | Tarefa pode apontar para lead |
| TEMPLATE / UI | Saída | Badge e toast global |
| NOTIFICAÇÕES | Serviço | Polling, lembrete e atraso |
| AUDITORIA | Planejado | Handoff de cortina/audit de tarefa |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Tarefa** | Atividade/lembrete com data e hora |
| **Responsável** | Usuário associado pelo campo `login` |
| **Pendente** | Tarefa ainda não resolvida (`status='N'`) |
| **Resolvida** | Tarefa finalizada (`status='S'`) |
| **Atrasada** | Tarefa pendente cujo horário já passou |
| **Lembrete 5 minutos** | Notificação antes do vencimento |
| **Notificação na hora** | Alerta no horário exato |
| **Toast persistente** | Notificação visual até ação do usuário |

---

## 3. Fluxos Principais

### 3.1 Cadastro de Tarefa

```
1. Usuário acessa TarefasForm
2. Preenche:
   - tarefa (descrição/título)
   - dt_tarefa
   - hr_tarefa
   - login (responsável)
   - lead_id opcional
   - status inicial N (Pendente)
3. Salva tarefa
4. TarefasList recarrega
```

**Regras**:
- RN-001: `tarefa`, `dt_tarefa`, `hr_tarefa` e `login` são obrigatórios
- RN-002: Status válido: `N` Pendente, `S` Resolvido
- RN-003: `lead_id` é opcional

### 3.2 Listagem de Tarefas

```
1. Usuário acessa TarefasList
2. Sistema filtra tarefas conforme permissão/login
3. Exibe:
   - tarefa
   - data/hora
   - status
   - responsável
   - lead vinculado se houver
4. Ações: editar, resolver, excluir conforme permissão
```

### 3.3 Kanban de Tarefas

```
1. Usuário acessa TarefasKanban
2. Sistema cria colunas:
   - Pendente
   - Resolvido
   - Atrasada
3. Carrega tarefas do usuário/contexto
4. Tarefa pendente com data/hora passada vai para Atrasada
5. Arrastar card entre colunas atualiza status:
   - para Resolvido → status='S'
   - para Pendente → status='N'
```

### 3.4 Notificações Globais

```
1. JS global faz polling a cada 60 segundos
2. Endpoint TarefaNotificationCheck chama TarefaNotificationService
3. Service carrega tarefas pendentes do login atual
4. Para cada tarefa:
   - resolve timestamp (dt_tarefa + hr_tarefa)
   - ignora tarefas futuras de outro dia
   - calcula diferença para agora
5. Classifica:
   - five_minutes: faltam até 5 minutos
   - on_time: horário atual ±1 minuto
   - overdue: passou mais de 1 minuto
6. Retorna a próxima notificação relevante
7. Frontend exibe toast persistente e atualiza badge
8. Service marca como tratado na sessão
```

### 3.5 Priorização de Notificações

Ordem de prioridade:
1. Tarefas na janela de 5 minutos antes
2. Tarefas no horário exato
3. Tarefas atrasadas

Dentro de cada grupo, exibe a mais urgente/mais próxima.

### 3.6 Controle Anti-duplicidade

```
1. Chave de sessão: TarefaNotificationService_handled_<tipo>_<id>
2. five_minutes e on_time: notifica apenas uma vez
3. overdue: pode notificar novamente após 15 minutos
4. total_pending alimenta badge visual
```

---

## 4. Estados

### 4.1 Máquina de Estados

```
PENDENTE (N) ──── resolver ───▶ RESOLVIDO (S)
     │                             │
     └── horário passou ─────▶ ATRASADA (visual, status continua N)
```

### 4.2 Status Persistidos

| Código | Label | Cor |
|--------|-------|-----|
| `N` | Pendente | #16a34a |
| `S` | Resolvido | #4f46e5 |

Aliases suportados pelo model:
- `P` como pendente
- `R` como resolvido

### 4.3 Estados Visuais

| Estado Visual | Condição |
|---------------|----------|
| Pendente | status N e horário futuro |
| Resolvido | status S |
| Atrasada | status N e timestamp < agora |

---

## 5. Visualizações

| Tela | Função |
|------|--------|
| `TarefasList` | Listagem padrão |
| `TarefasForm` | Cadastro/edição |
| `TarefasKanban` | Visão por status visual |
| `TarefaNotificationCheck` | Checagem de notificação |

---

## 6. Notificações

### 6.1 Tipos

| Tipo | Condição | Repetição |
|------|----------|-----------|
| `five_minutes` | 0 < diff <= 5min | Uma vez |
| `on_time` | abs(diff) <= 1min | Uma vez |
| `overdue` | diff < -1min | A cada 15min |

### 6.2 Comportamento UI

- Toast persistente
- Click abre TarefasForm/Tarefa relacionada
- Badge no menu mostra total pendente relevante
- Polling global a cada 60s

---

## 7. Decisões Arquiteturais

### ADR-001: Notificação Global via Polling

**Decisão**: JavaScript global consulta endpoint a cada 60 segundos.

**Consequências**:
- Funciona em qualquer página
- Não requer WebSocket
- Baixo custo operacional

### ADR-002: Notificações por Login Criador/Responsável

**Decisão**: Service filtra por `TSession::getValue('login')` e `tarefa.login`.

**Consequências**:
- Usuário vê apenas suas tarefas
- Simples e compatível com login legado

### ADR-003: Atrasada é Estado Visual

**Decisão**: Tarefa atrasada não muda status no banco; continua `N`.

**Consequências**:
- Sem escrita automática
- Kanban calcula coluna atrasada dinamicamente

### ADR-004: Sessão para Anti-duplicidade

**Decisão**: Controle de notificações já exibidas fica em sessão.

**Consequências**:
- Evita spam no mesmo login/sessão
- Ao reiniciar sessão, lembretes podem reaparecer

### ADR-005: Kanban como Visão Alternativa

**Decisão**: TarefasKanban usa `TKanban` com colunas Pendente/Resolvido/Atrasada.

**Consequências**:
- UX visual simples
- Drag-and-drop pode atualizar status

---

## 8. Referência à Implementação Atual

### 8.1 Controllers

| Controller | Localização |
|------------|-------------|
| `TarefasList` | `app/control/tarefas/TarefasList.php` |
| `TarefasForm` | `app/control/tarefas/TarefasForm.php` |
| `TarefasKanban` | `app/control/tarefas/TarefasKanban.php` |
| `TarefaNotificationCheck` | `app/control/tarefas/TarefaNotificationCheck.php` |

### 8.2 Model/Service

| Artefato | Localização |
|----------|-------------|
| `Tarefa` | `app/model/Tarefa.php` |
| `TarefaNotificationService` | `app/service/tarefas/TarefaNotificationService.php` |

### 8.3 Campos Principais

| Campo | Descrição |
|-------|-----------|
| `tarefa` | Texto da tarefa |
| `dt_tarefa` | Data |
| `hr_tarefa` | Hora |
| `status` | N/S |
| `login` | Responsável |
| `lead_id` | Lead opcional |

### 8.4 Docs

| Documento | Descrição |
|-----------|-----------|
| `docs/notificacoes/NOTIFICACOES_TAREFAS.md` | Sistema de notificações globais |
| `docs/handoff/HANDOFF_CORTINA_AUDIT_TAREFA.md` | Handoff sobre cortina/auditoria |

---

## 9. Segurança

- Usuário recebe notificação apenas de suas tarefas
- Listagem deve respeitar login/permissão
- Excluir tarefa deve ser restrito ao criador ou gestor
- Tarefas vinculadas a lead devem respeitar acesso ao lead
- Não deve haver escrita automática apenas por notificação

---

**Versão**: 1.0
**Data**: 2026-05-11
**Autor**: SMID Architecture Team
**Status**: Draft para revisão
