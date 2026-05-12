# TMultiCombo

Objetivo
- Padronizar o uso do `TMultiCombo` no SMID (Adianti 8.3/8.4) com foco em UX, tema dark/light e placeholder consistente.

Contexto
- O placeholder padrao do `TMultiCombo` vem do core (`AdiantiCoreTranslator`) e e iniciado no `tmulticombo_start` do framework.
- Nao devemos alterar arquivos em `lib/adianti/*` (core).

Onde aplicar
- Formularios com filtro multi-selecao.
- Preferir padrao quando o campo precisa selecionar mais de um item (ex.: Status no LeadsList).

Padrao de implementacao (Controller)
- Exemplo (usar `BootstrapFormBuilder`):

```php
use Adianti\Widget\Form\TMultiCombo;

$status_id = new TMultiCombo('status_id');
$status_id->setSize('100%');
$status_id->addItems($statusItems); // array [id => nome]

$row = $this->form->addFields(
    [new TLabel('Status'), $status_id]
);
$row->layout = ['col-sm-6', 'col-sm-6'];
```

Placeholder padrao (sem mexer no core)
- O framework usa "Click to search" / "Clique para buscar" no placeholder do `TMultiCombo`.
- Para padronizar como "Selecionar", use override em `app/templates/adminbs5/js/theme.js` (hook permitido).

Padrao recomendado (theme.js)
- Interceptar o `tmulticombo_start` e trocar os labels globalmente para todo `TMultiCombo`.
- Exemplo (global):

```js
Template.overrideMultiComboLabels = function () {
    if (typeof window.tmulticombo_start !== 'function') {
        return;
    }

    if (window._smid_tmulticombo_start) {
        return;
    }

    window._smid_tmulticombo_start = window.tmulticombo_start;
    window.tmulticombo_start = function (element_id, size, labels) {
        labels = labels || {};
        labels.placeholder = 'Selecionar';
        labels.noneSelected = 'Selecionar';
        return window._smid_tmulticombo_start(element_id, size, labels);
    };
};
```

Estilos recomendados
- Ajustes de alinhamento, altura e cores devem ficar em:
- `app/templates/adminbs5/custom.css` (estilos gerais)
- `app/templates/adminbs5/theme.css` (dark/light)

Hooks utilizados (CSS/JS)
- CSS (custom.css)
- `.ms-options-wrap` e `.ms-options-wrap > button` para tamanho, padding, fonte e alinhamento visual do controle.
- `select[widget="tmulticombo"]` para ocultar o `<select>` original (renderizado pelo widget).
- CSS (theme.css)
- `[data-bs-theme="dark"] .ms-options-wrap > button` e `.ms-options-wrap.ms-has-selections > button` para cores do placeholder e do valor selecionado.
- `[data-bs-theme="dark"] .ms-options-wrap > .ms-options` e `.ms-options-wrap > .ms-options > .ms-search input` para fundo/borda do dropdown.
- Equivalentes em `[data-bs-theme="light"]` para o tema claro.
- JS (theme.js)
- `Template.overrideMultiComboLabels()` faz o hook do `window.tmulticombo_start` e ajusta `labels.placeholder`/`labels.noneSelected` globalmente para todo `TMultiCombo`.
- O hook é chamado uma única vez no load do template (evita duplicar override).

Onde esta no codigo atualmente
- `app/templates/adminbs5/custom.css` (bloco "TMultiCombo - alinhamento visual com inputs padrão")
- `app/templates/adminbs5/theme.css` (blocos `[data-bs-theme="dark"] .ms-options-wrap...` e `[data-bs-theme="light"] .ms-options-wrap...`)
- `app/templates/adminbs5/js/theme.js` (função `Template.overrideMultiComboLabels()` + chamada no carregamento do template)

Padroes e referencias
- `AGENTS.md` (padroes gerais e arquivos permitidos)
- `docs/upgrade/SMID_PERSONALIZACOES.md` (customizacoes permitidas/proibidas)
- `docs/upgrade/UPGRADES_FUTUROS.md` (reaplicar hook e estilos em upgrades)

Checklist
1. Nao alterar `lib/adianti/*`.
2. Usar override em `app/templates/adminbs5/js/theme.js`.
3. Garantir placeholder padrao "Selecionar".
4. Validar dark e light mode.
5. Testar em 320, 375, 414 e 428 px.
