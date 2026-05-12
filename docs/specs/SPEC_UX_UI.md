# SPEC_UX_UI.md — Padrões de UX e UI do SMID

> Este SPEC define o contrato de experiência do SMID. É **agnóstico de framework** — o que importa é preservar os padrões de interação, responsividade e composição de tela.

---

## 1. Princípios

1. **Desktop first, mobile premium**: telas devem funcionar em 320, 375, 414, 428, 768 e desktop
2. **Uma tela, uma intenção**: formulários e listas com propósito claro
3. **Contexto preservado**: navegação não deve perder filtros/estado
4. **Resposta previsível**: busca dinâmica só quando a UX pede
5. **Acessibilidade via cortina**: painéis laterais em vez de modais empilhados
6. **Tema dual**: suporte claro e escuro via variáveis CSS

---

## 2. Componentes Canônicos

### 2.1 SMIDStandardForm

Classe base para formulários do SMID. Deve oferecer:

- Cabeçalho com título e breadcrumb
- Corpo com `TForm` responsivo
- Rodapé com ações (Salvar, Limpar, Voltar, Excluir quando aplicável)
- Integração automática com cortina lateral quando chamada em modo painel
- Suporte a `mobile-premium` ativado por padrão
- Hooks: `onSave`, `onEdit`, `onClear`, `onDelete`, `onReload`

### 2.2 SMIDStandardList

Classe base para listagens. Deve oferecer:

- Cabeçalho com título, botão Novo e Busca Rápida (conforme árvore de decisão)
- `TDataGrid` com ordenação, paginação e ações por linha
- Filtro Avançado opcional com sessão isolada
- Rodapé com total de registros e paginação
- Hooks: `onReload`, `onQuickSearch`, `onClearSearch`, `onApplyAdvancedFilters`

### 2.3 SMIDPage + TSidePanel (Cortina Lateral)

Cortina lateral é o mecanismo padrão para formulários invocados a partir de listas/dashboards.

**Regras**:
- Abre à direita em desktop, sobrepondo parcialmente a tela
- Ocupa 100% em mobile
- Fechamento por backdrop, tecla ESC ou botão explícito
- Ao salvar com sucesso, deve recarregar a listagem pai
- Cortina não deve abrir outra cortina (profundidade = 1)

### 2.4 Kanban (TKanban)

Usado para tarefas e, quando aplicável, pedidos/visitas:

- Colunas representam estados
- Cards exibem resumo da entidade
- Drag-and-drop dispara ação de mudança de estado
- Limite de cards por coluna (ex.: Kanban de Leads usa `MAX_CARDS = 400` por coluna)

### 2.5 Dashboard (TDashboard)

Composto por cards, gráficos e listas de drill-down. Filtros de período usam `TDateRangeField`.

---

## 3. Busca Rápida

### 3.1 Árvore de Decisão Obrigatória

```
A lista é longa/paginada, usa Filtro Avançado ou precisa busca global?
    ├── SIM → busca server-side manual (padrão oficial)
    │         campo `input_search` + botão lupa + botão Limpar
    │         termo salvo em sessão `<Class>_search_filter`
    │
    └── NÃO → lista curta? Pode usar `enableSearch` client-side
              (sem botão Limpar, apenas filtro visual)
```

### 3.2 Server-side (padrão oficial)

Obrigatório:
- Campo `input_search` no header
- Método `onQuickSearch()` salva termo em sessão dedicada
- Método `onClearSearch()` limpa a mesma sessão
- Construtor reidrata o campo a partir da sessão
- `onReload()` aplica o filtro persistido
- `onkeydown` impede submit por Enter
- Botão lupa dispara a busca; botão Limpar a reseta
- Paginação e exportação preservam a busca ativa
- Filtro Avançado usa sessão separada

Proibido nesse padrão:
- `oninput` com debounce
- `onkeyup` automático
- `setExitAction` sobre o input
- `datagrid->enableSearch()`

### 3.3 Client-side (enableSearch)

Apenas para listas curtas (< 50 registros) que não precisam de filtro global.

---

## 4. TDateRangeField

Componente canônico para filtros por período.

### 4.1 Presets Obrigatórios

| Preset | Semântica |
|--------|-----------|
| `equal` | Data exata |
| `less` | Antes de |
| `greater` | Depois de |
| `interval` | Intervalo entre duas datas |
| `today` | Hoje |
| `yesterday` | Ontem |
| `last7days` | Últimos 7 dias |
| `last30days` | Últimos 30 dias |
| `last3months` | Últimos 3 meses |
| `thisMonth` | Mês atual |
| `lastMonth` | Mês anterior |

### 4.2 Regras Obrigatórias

- Cálculo de presets **apenas no backend**, nunca via JavaScript
- Presets normalizados em `onSearch` / `onApplyAdvancedFilters`
- Filtros persistidos em sessão dedicada por tela
- Listar todos os campos-base em `$dateFields` para reidratação correta

---

## 5. TMultiCombo

Componente canônico para seleção múltipla.

### 5.1 Regras

- Placeholder padrão hooked via `app/templates/adminbs5/js/theme.js`
- Nunca alterar core (`lib/adianti/*`)
- Preencher via `addItems(array $items)` com `id => label`
- Respeitar ordenação do negócio (alfabética ou hierárquica)

---

## 6. TFieldList

Para inclusão dinâmica de itens (ex.: produtos em pedido).

### 6.1 Regra

Preferir **reconstrução das linhas pelo servidor** (`addDetail`) em vez de sincronização frágil via `TForm::sendData` + timeout. Isso evita perda de dados em telas mobile.

---

## 7. Mobile Premium

### 7.1 Assets Obrigatórios

| Arquivo | Propósito |
|---------|-----------|
| `app/templates/css/mobile-premium.css` | Estilos responsivos |
| `app/templates/js/mobile-premium.js` | Comportamentos responsivos |

Incluir em **todas** as páginas, forms e listagens.

### 7.2 Breakpoints

| Faixa | Dispositivo |
|-------|-------------|
| 320–480px | Mobile pequeno |
| 481–767px | Mobile padrão |
| 768–991px | Tablet |
| 992+px | Desktop |

### 7.3 Validação Obrigatória

Antes de concluir qualquer tela, validar em:
- 320px
- 375px
- 414px
- 428px
- 768px
- Desktop padrão

### 7.4 Padrões Mobile

- Campos de formulário empilham em mobile
- Tabelas viram cards ou rolam horizontal com indicador
- Cortina lateral ocupa 100% em mobile
- Kanban vira carrossel horizontal
- Botões críticos ficam fixos no rodapé

---

## 8. Tema Visual

### 8.1 Arquivos Permitidos

| Arquivo | Uso |
|---------|-----|
| `app/templates/adminbs5/custom.css` | Customizações gerais |
| `app/templates/adminbs5/theme.css` | Tema claro/escuro |
| `app/templates/adminbs5/js/theme.js` | Hooks de UI (ex.: TMultiCombo) |

### 8.2 Arquivos Proibidos

| Caminho | Motivo |
|---------|--------|
| `lib/adianti/*` | Core do framework |
| `tutor/*` | Referência |
| `template_exemplo/*` | Referência |
| `app/templates/adminbs5/css/*` (exceto permitidos) | Estrutura do template |
| `app/templates/adminbs5/js/*` (exceto `theme.js`) | Estrutura do template |

### 8.3 Tema Dark/Light

Alterações de esquema de cor **apenas** em `theme.css`. Nunca hardcodar cores em componentes.

---

## 9. Notificações e Mensagens

| Tipo | Função |
|------|--------|
| `__adianti_message('success', msg)` | Confirmação |
| `__adianti_message('error', msg)` | Falha |
| `__adianti_message('warning', msg)` | Alerta |
| `__adianti_message('info', msg)` | Informação |

Mensagens devem ser curtas, em português, sem emojis.

---

## 10. Navegação Protegida

Para fluxos de leads/visitas/pedidos:

### 10.1 Padrão Stack Fallback

```
1. Tentar recuperar contexto da sessão corrente
2. Se falhar, consultar histórico de navegação (SMIDVoltaPagina*)
3. Se falhar, redirecionar para listagem raiz do módulo
4. Nunca deixar usuário em tela sem contexto
```

### 10.2 Implementações

- `LeadNavigationHelper`
- `VisitaNavigationHelper`
- `PedidosRepreNavigationHelper`
- `PedidosTeleNavigationHelper`
- `SMIDVoltaPaginaForm/List/Calendario`

---

## 11. Ações por Linha (cores canônicas)

| Cor | Semântica |
|-----|-----------|
| Azul | Edição/visualização |
| Verde | Ação positiva (aprovar, salvar) |
| Amarelo | Alerta/atenção |
| Vermelho | Destrutivo (excluir, cancelar) |
| Cinza | Neutro/info |

---

## 12. Identidade Textual

- Idioma da interface: **português do Brasil**
- Código (classes, métodos, variáveis): **inglês**
- Sem emojis em qualquer lugar
- Rótulos em formulários: frase capitalizada (ex.: "Nome completo")
- Cabeçalhos de tabela: frase capitalizada curta

---

## 13. Referências à Implementação Atual

| Artefato | Localização |
|----------|-------------|
| Bases SMID | `app/lib/smid/SMIDStandardForm.php`, `SMIDStandardList.php`, `SMIDPage.php` |
| Cortina | `app/lib/widget/TSidePanel.php` |
| Date Range | `app/lib/widget/TDateRangeField.php` |
| Helpers navegação | `app/lib/smid/*NavigationHelper.php` |
| Mobile | `app/templates/css/mobile-premium.css`, `app/templates/js/mobile-premium.js` |
| Docs | `docs/BUSCA_RAPIDA.md`, `docs/mobile/*`, `docs/TMultiCombo.md`, `docs/upgrade/SMID_PERSONALIZACOES.md` |

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
