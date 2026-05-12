# SMID_PERSONALIZACOES.md

Guia rápido de onde personalizar o template adminbs5 sem comprometer upgrades futuros.

## 1. Arquivos permitidos para customização

| Finalidade                    | Arquivo                                               | Observações |
|-------------------------------|--------------------------------------------------------|-------------|
| Estilos SMID gerais           | `app/templates/adminbs5/custom.css`                    | Estilos globais, ajustes de componentes, regras específicas do projeto. |
| Estilos de formulários/listas | `app/templates/adminbs5/theme.css`                     | Manter o padrão dark/light com `[data-bs-theme]`. |
| Regras mobile premium         | `app/templates/css/mobile-premium.css`                 | Sempre incluir via `TPage::include_css/js`. |
| Hook JS (menu, busca, painel) | `app/templates/adminbs5/js/theme.js`                   | Utilizar bloco SMID (`Template.applySmidOverrides`). |
| Pontos de injeção no layout   | `app/templates/adminbs5/layout.html`                   | Alterações mínimas e pontuais (ex.: `#phone-search-box` e include do `mobile-premium.css`). |
| TMessage/TQuestion (dark)     | `app/templates/adminbs5/js/theme.js`, `app/templates/adminbs5/custom.css` | Tema do SweetAlert2 por CSS vars (dark/light) + ajustes visuais do popup. |
| Popups (Perfil/Notificações)  | `app/templates/adminbs5/js/theme.js`                   | Não instanciar Tooltip em elementos `.dropdown-toggle` (evita conflito com Dropdown no Bootstrap). |
| Menu/atalhos no topo          | `app/lib/menu/AdiantiNavBarParser.php`                 | Mantém tradução `_t{}` e adiciona regras SMID (ex.: `quick_lead=1`, links em `target_container=adianti_right_panel`). |
| Aba “abrir em nova guia”      | `app/lib/menu/AdiantiMenuBuilder.php`                  | Ajuste SMID para permitir abrir itens do menu lateral em abas (exceto logout). |
| Widgets/Helpers específicos   | `app/control/PhoneSearchBox.php`, `app/lib/widget/*`   | Manter headers padrão SMID. |
| Entrypoints (exceção)         | `engine.php`, `init.php`                               | Preservar ajustes SMID (ex.: `upload_config.php`, permissões padrão e mensagem/redirect de acesso). |
| Controllers de log (upstream) | `app/control/log/*`                                    | Manter alinhado ao template limpo (preferir sobrescrever no upgrade; não adicionar customizações aqui). |

## 2. Arquivos que **não** devem ser modificados

- `app/templates/adminbs5/css/*.css` (exceto `custom.css` e `theme.css`).
- `app/templates/adminbs5/js/*.js` (exceto `theme.js`).
- `app/templates/adminbs5/*.html` (exceto `layout.html` para os pontos SMID documentados acima).
- `.htaccess` em `app/*` (manter alinhado ao template limpo: `Require all denied`).
- Qualquer arquivo em `template_exemplo/` (referência limpa).
- `lib/adianti/*` (core do framework).

## 3. Boas práticas

1. **Nunca sobrescreva** arquivos padrão: crie regras mais específicas em `custom.css` ou `theme.css`.
2. **Isolar ganchos JS** em funções como `Template.applySmidOverrides()` para reaplicar após upgrades.
3. **Guardar assets mobile** (`mobile-premium.css`, `mobile-premium.js`) e incluí-los via `TPage::include_*`.
4. **Documentar mudanças** relevantes neste arquivo e em `UPGRADES_FUTUROS.md`.
5. **Usar commits granulares**: identifique o que é update de core vs. customização SMID.
6. **TMultiCombo:** padronizar placeholder via hook em `app/templates/adminbs5/js/theme.js` e estilos em `custom.css`/`theme.css` (ver `docs/TMultiCombo.md`). Nunca alterar `lib/adianti/*`.

## 4. Checklist antes de commitar

- [ ] Apenas `custom.css`, `theme.css`, `mobile-premium.css/js`, `theme.js` foram alterados?
- [ ] Não há ajustes diretos em `css/layout.css`, `css/sidebar.css`, etc.?
- [ ] O Busca Fone, monitor de tarefas e menu continuam funcionando?
- [ ] `app/control/log/*` continua igual ao `template_exemplo/`?
- [ ] Em `Framework → SystemFilesDiff`, a lista de “Modificado” está restrita às exceções SMID (ex.: `engine.php`, `init.php`, `app/lib/menu/*`, `app/templates/adminbs5/{custom.css,theme.css,js/theme.js,layout.html}`)?
- [ ] Testes mobile realizados nos breakpoints 320, 375, 414 e 428 px?

Seguindo estas diretrizes, os upgrades futuros ficam restritos a substituir os arquivos padrão e reaplicar as customizações centralizadas, evitando retrabalho.
