# BUSCA_RAPIDA.md

Padrão oficial de Busca Rápida para listagens (`*List.php`) no SMID 8.0 (Adianti 8.4).

## Objetivo
- Padronizar UX, estado e implementação.
- Evitar busca dinâmica indevida em listas paginadas.
- Evitar conflito entre Busca Rápida, paginação, Filtro Avançado e exportações.

## Regra de decisão (obrigatória)
1. A lista é longa, paginada, usa Filtro Avançado, exporta dados ou precisa buscar no conjunto completo?
   - Use busca **server-side manual**.
2. A lista é curta e a busca local na página atual é suficiente?
   - Pode usar busca **client-side** com `enableSearch`.
3. O usuário pediu explicitamente que a busca "não seja dinâmica"?
   - **Não** use `oninput`, `onkeyup`, debounce, `setExitAction()` nem `enableSearch()` implícito.
   - Use busca **server-side manual por lupa**.

## Padrão consolidado do projeto
Hoje, o padrão preferencial para listas operacionais é:
- campo `input_search` no header;
- botão `Buscar` com lupa;
- botão `Limpar`;
- termo salvo em sessão;
- reidratação do campo no construtor;
- aplicação da busca apenas ao clicar na lupa;
- limpeza da mesma chave de sessão em `onClearSearch()`.

## Server-side manual (padrão oficial)
### Quando usar
- Lista paginada.
- Busca global entre páginas.
- Tela com Filtro Avançado.
- Tela com exportação.
- Quando a busca não pode ser dinâmica.

### Requisitos obrigatórios
- `input_search` no header.
- `onQuickSearch()` salvando o termo em `SEARCH_SESSION`.
- `onClearSearch()` limpando a **mesma** `SEARCH_SESSION`.
- Construtor reidratando o campo com o valor da mesma sessão.
- `onReload()` lendo:
  - `input_search` do request, quando vier;
  - senão, o valor salvo em sessão.
- `onkeydown` no input para impedir submit por Enter.
- Botão `Buscar` disparando `onQuickSearch()`.
- Botão `Limpar` disparando `onClearSearch()`.
- Busca e Filtro Avançado com sessões separadas.

### Regras obrigatórias de UX
- A busca só aplica no clique da lupa.
- `Limpar` deve zerar o campo visualmente e limpar a sessão da busca.
- Paginação deve preservar a busca ativa.
- Exportação deve usar o mesmo critério da busca ativa.

### O que não usar neste padrão
- `oninput`
- `onkeyup`
- debounce JavaScript
- `setExitAction(new TAction([__CLASS__, 'onQuickSearch']))`
- `datagrid->enableSearch(...)`

### Snippet base
```php
private const SEARCH_SESSION = __CLASS__ . '_search_filter';

$input_search = new TEntry('input_search');
$input_search->placeholder = 'Busca Rápida';
$input_search->setSize('100%');
$input_search->setProperty('onkeydown', "if (event.key === 'Enter') { event.preventDefault(); return false; }");

$storedSearch = TSession::getValue(self::SEARCH_SESSION);
if (is_string($storedSearch)) {
    $input_search->setValue($storedSearch);
}

$searchForm = new TForm('form_quick_search');

$searchWrapper = new TElement('div');
$searchWrapper->{'class'} = 'list-header-quick flex-grow-1 d-flex gap-2 align-items-center';
$searchWrapper->add($input_search);

$searchActions = new TElement('div');
$searchActions->{'class'} = 'smid-quick-actions d-flex align-items-center';

$btn_search = new TButton('btn_quick_search');
$btn_search->setAction(new TAction([__CLASS__, 'onQuickSearch']), '');
$btn_search->setFormName($searchForm->getName());
$btn_search->addStyleClass('btn btn-default');
$btn_search->setImage('fa:search blue');
$btn_search->setProperty('title', 'Buscar');
$btn_search->setProperty('type', 'button');
$searchActions->add($btn_search);

$btn_clear = new TButton('btn_clear_search');
$btn_clear->setAction(new TAction([__CLASS__, 'onClearSearch']), '');
$btn_clear->setFormName($searchForm->getName());
$btn_clear->addStyleClass('btn btn-default');
$btn_clear->setImage('fa:broom orange');
$btn_clear->setProperty('title', 'Limpar busca');
$btn_clear->setProperty('type', 'button');
$btn_clear->setProperty('onmousedown', "var input = document.querySelector('input[name=\"input_search\"]'); if (input) { input.value = ''; }");
$searchActions->add($btn_clear);

$searchWrapper->add($searchActions);

$headerWrapper = new TElement('div');
$headerWrapper->{'class'} = 'list-header-bar d-flex align-items-center gap-2 w-100';
$headerWrapper->add($searchWrapper);

$searchForm->add($headerWrapper);
$searchForm->setFields([$input_search, $btn_search, $btn_clear]);
```

### Métodos base
```php
public static function onQuickSearch($param = null)
{
    $param = (array) ($param ?? []);
    $search = isset($param['input_search']) ? trim((string) $param['input_search']) : '';
    TSession::setValue(self::SEARCH_SESSION, $search !== '' ? $search : null);

    $args = ['page' => 1, 'focus_search' => 1];
    if ($search !== '') {
        $args['input_search'] = $search;
    }

    AdiantiCoreApplication::loadPage(__CLASS__, 'onReload', $args);
}

public static function onClearSearch($param = null)
{
    TSession::setValue(self::SEARCH_SESSION, null);

    AdiantiCoreApplication::loadPage(__CLASS__, 'onReload', [
        'page' => 1,
        'focus_search' => 1,
        'clear_search' => 1,
        'input_search' => '',
    ]);
}
```

### Leitura correta no `onReload()`
```php
if (array_key_exists('input_search', $param)) {
    $searchValue = trim((string) $param['input_search']);
    TSession::setValue(self::SEARCH_SESSION, $searchValue !== '' ? $searchValue : null);
} else {
    $storedSearch = TSession::getValue(self::SEARCH_SESSION);
    $searchValue = is_string($storedSearch) ? trim($storedSearch) : '';
}

if ($searchValue !== '') {
    // aplicar filtros do repositório/criteria/model
    $param['input_search'] = $searchValue;
} else {
    unset($param['input_search']);
}
```

## Client-side (`enableSearch`)
### Quando usar
- Lista curta.
- Busca local na página atual.
- Sem necessidade de paginação global.
- Busca instantânea por digitação é aceitável.

### Requisitos
- `input_search` no header.
- `datagrid->enableSearch($input_search, 'campos')`.
- Sem sessão obrigatória.
- Botão `Limpar` é opcional.

### Snippet base
```php
$input_search = new TEntry('input_search');
$input_search->placeholder = 'Busca Rápida';
$input_search->setSize('100%');

$this->datagrid->enableSearch($input_search, 'campo1,campo2');
$panel->addHeaderWidget($input_search);
```

### Limitação importante
`enableSearch()` é **dinâmico por natureza**. Ele liga a busca ao `keyup`.

Se a tela precisar:
- lupa manual;
- sessão;
- paginação coerente;
- ou o usuário pedir que "não seja dinâmica";

então **não basta adicionar botões**. É obrigatório remover `enableSearch()` e migrar a tela para o padrão server-side manual.

## Layout e spacing dos botões
### Padrão geral
- Wrapper da busca:
  - `list-header-quick flex-grow-1 d-flex gap-2 align-items-center`
- Grupo dos botões:
  - `smid-quick-actions d-flex align-items-center`

### Regra de CSS consolidada
Existe regra global `.btn+.btn { margin-left: 8px; }`.

Para colar `Buscar` e `Limpar` no padrão atual, usar em `app/templates/adminbs5/custom.css`:
```css
div[page_name] .smid-quick-actions {
    gap: 0;
}

div[page_name] .smid-quick-actions .btn+.btn {
    margin-left: 0;
}
```

### Caso legado do `TelevendasList`
`TelevendasList` usa `televendas-quick-actions`. O mesmo princípio vale:
```css
div[page_name="TelevendasList"] .televendas-quick-actions {
    gap: 0;
}

div[page_name="TelevendasList"] .televendas-quick-actions .btn+.btn {
    margin-left: 0;
}
```

## Busca Rápida + Filtro Avançado
### Regras
- Use chaves separadas:
  - `SEARCH_SESSION` para busca rápida;
  - `FILTER_SESSION` para filtro avançado.
- `onClearSearch()` limpa apenas a busca rápida.
- `onClearFilters()` limpa apenas os filtros avançados, salvo regra explícita da tela.
- `onReload()` deve compor os dois critérios sem misturar as chaves.

### Exportação
Se a lista exporta dados:
- o export deve ler o `input_search` do request quando vier;
- senão, deve ler a mesma `SEARCH_SESSION`;
- os filtros aplicados na grid e na exportação devem ser espelhados.

## Armadilhas conhecidas
- Adicionar a lupa, mas deixar `enableSearch()` ativo:
  - a busca continua dinâmica.
- Limpar a chave errada de sessão:
  - o campo parece limpo, mas a busca continua aplicada.
- Limpar só o campo visual:
  - a sessão continua filtrando a grid.
- Ler uma chave no construtor e outra no `onReload()`:
  - a tela fica inconsistente.
- Registrar `TElement` como campo do `TForm`:
  - wrappers de layout entram com `add($wrapper)`;
  - só widgets reais entram em `setFields([...])`.

## Testes manuais obrigatórios
### Server-side manual
1. Digitar texto e confirmar que nada acontece antes da lupa.
2. Clicar na lupa e validar o filtro.
3. Paginar com busca ativa e validar persistência.
4. Clicar `Limpar` e validar:
   - campo vazio;
   - sessão limpa;
   - grid completa.
5. Aplicar Filtro Avançado junto e validar combinação dos critérios.
6. Se houver exportação, validar que exporta o mesmo recorte.

### Client-side
1. Digitar texto e validar filtro instantâneo.
2. Confirmar com o time que busca local atende ao caso.
3. Não usar este padrão quando o comportamento esperado for manual por botão.

## Referências atuais de implementação
### Server-side manual
- `app/control/televendas/TelevendasList.php`
- `app/control/leads/LeadsList.php`
- `app/control/historicos/HistoricoList.php`
- `app/control/visitas/pos-visita/PosVisitaList.php`
- `app/control/visitas/VisitaList.php`
- `app/control/leads/LeadDuplicadoList.php`
- `app/control/pedidos/PedidosList.php`
- `app/control/financeiro/FinExtratoList.php`
- `app/control/financeiro/FinContasPagarList.php`
- `app/control/financeiro/FinContasReceberList.php`
- `app/control/produtos/ProdutoList.php`
- `app/control/compras/ComprasList.php`
- `app/control/suportes/SuporteList.php`
- `app/control/televendas/pedidos/PedidosTeleList.php`
- `app/control/televendas/orcamentos/OrcamList.php`
- `app/control/metas/MetaList.php`

### Client-side
- `app/control/leads/equipes/EquipeList.php`
