# Documentação Mobile - Pedidos (SMID 8.0)

Este documento detalha a implementação da visualização e ações mobile para o formulário de Pedidos (`PedidosForm`), abrangendo as abas **Produtos**, **Comissões** e **Produção**.

## Visão Geral

A adaptação mobile foca em oferecer uma experiência otimizada para telas pequenas, substituindo as grids de dados (`TFieldList`) por listas de cartões responsivos e utilizando modais para edição e inserção de dados. O layout foi ajustado para remover espaçamentos desnecessários e integrar-se ao tema Dark/Light do sistema.

## Arquiteturas e Componentes

### 1. Controladores e Wrappers
Cada aba que possui uma lista de itens (`PedProdDetalhes`, `ComisPedDetalhes`, `ComprasDetalhes`) foi modificada para incluir dois wrappers de visualização:
- **Desktop Wrapper (`d-none d-md-block`)**: Mantém a visualização original em `TFieldList` para telas maiores.
- **Mobile Wrapper (`d-block d-md-none`)**: Renderiza os itens como cartões (Cards) otimizados para toque.

### 2. Ações Mobile Centralizadas (_Pattern_ MobileActions)
Um novo controlador, `app/control/pedidos/PedidosMobileActions.php`, foi criado para centralizar as ações exclusivas do mobile (modais de add/edit/delete). Isso evita a complexidade de adaptar os métodos existentes que dependem fortemente de `TFieldList` JavaScript.

#### Métodos Principais:
- `onAddProduto`, `onEditProduto`, `onSaveProduto`
- `onAddComissao`, `onEditComissao`, `onSaveComissao`
- `onAddProducao`, `onEditProducao`, `onSaveProducao`
- `onAddPagamento`, `onEditPagamento`, `onDeletePagamento`, `onSavePagamento` (Gestão de pagamentos internos de comissões)

### 3. Layout e Estilo (CSS)
- Foram utilizadas classes utilitárias do Bootstrap 5 (`d-block`, `d-none`, `p-0`) para controle de responsividade.
- As abas do `PedidosForm` tiveram o padding fixo (`20px`) substituído por classes responsivas (`p-0 p-md-4`) para garantir que os cartões mobile ocupem toda a largura da tela (Flush Layout).
- **Theming**: Botões de ação nos cartões utilizam a classe `.btn-default` (do `theme.css`) para adaptação automática aos modos Claro e Escuro, com ícones herdando cores apropriadas.
- **Correção Dark Mode (Pagamentos)**: Os cartões de pagamento internos (aba Comissões) utilizam classes padrão (`card`, `shadow-sm`) **sem** `bg-light` ou `text-dark` fixos, garantindo compatibilidade total com o modo escuro.
- **Remoção de Espaçamento Extra**: Uma regra CSS específica em `app/templates/css/mobile-premium.css` utiliza o seletor modern `:has()` para remover o padding do elemento `.card-body` pai quando este contém qualquer um dos novos wrappers mobile (`.mobile-produtos-wrapper`, `.mobile-comissoes-wrapper`, etc.). Isso resolve o problema de "falta de espaço" e permite que os cartões utilizem toda a largura disponível.

```css
/* mobile-premium.css */
@media (max-width: 768px) {
    .card-body:has(.mobile-visits-wrapper),
    .card-body:has(.mobile-produtos-wrapper),
    .card-body:has(.mobile-comissoes-wrapper),
    .card-body:has(.mobile-pagamentos-wrapper),
    .card-body:has(.mobile-producao-wrapper) {
        padding: 0 !important;
    }
}
```

## Implementação por Aba

### Produtos (`PedProdDetalhes`)
- **Lista**: Cartões mostrando Produto, Quantidade e Medida.
- **Ações**: Editar (Modal) e Excluir. Botões configurados com `type="button"` para evitar submissões acidentais do formulário principal.
- **Formulário Mobile**: Modal simplificado com busca de produtos e medidas.

### Comissões (`ComisPedDetalhes`)
- **Lista**: Cartões mostrando Valor, Data Prevista e Saldo.
- **Ações**: 
    - **Pagamentos**: Botão dedicado que expande/lista os pagamentos relacionados à comissão.
    - **Editar/Excluir**: Manutenção da comissão.
    - **Obs**: Métodos padronizados (`openMobileForm`) para consistência.
- **Sub-lista (Pagamentos)**: Renderizada abaixo da lista principal ou via interação, mostrando pagamentos vinculados.

### Produção (`ComprasDetalhes`)
- **Lista**: Cartões mostrando Fornecedor, Status ("Bling") e datas do processo (Produção/Envio/Chegada).
- **Ações**: Editar e Excluir lançamentos de compras/produção.

## Considerações Técnicas
- **Prevenção de Submissão**: Todos os botões de ação mobile (`Add`, `Edit`, `Delete`, `Pagamentos`) têm o atributo `type="button"` explicitamente definido. Isso impede que o navegador interprete o clique como um `submit` do formulário `PedidosForm` pai.
- **Validações**: As validações no mobile são simplificadas em relação ao desktop (focadas em campos obrigatórios essenciais), mas mantêm a integridade dos dados via Model.
- **Sincronização**: Ao salvar um item via mobile, a listagem correspondente é recarregada via AJAX, mantendo a interface fluida sem recarregar a página inteira erroneamente.
- **Manutenção**: Alterações de lógica de negócios devem ser refletidas tanto nos métodos de desktop (`onUpdateItem` etc.) quanto nos métodos mobile em `PedidosMobileActions`.

## Arquivos Modificados/Criados
- `app/control/pedidos/PedidosForm.php` (Ajustes de layout das abas)
- `app/control/pedidos/PedProdDetalhes.php` (Dual View)
- `app/control/comissoes/ComisPedDetalhes.php` (Dual View + Pagamentos)
- `app/control/compras/ComprasDetalhes.php` (Dual View)
- `app/control/pedidos/PedidosMobileActions.php` (Novo controlador de ações)
- `docs/PEDIDOS_MOBILE.md` (Esta documentação)

---
**Autor**: Agente AI (Antigravity)
**Data**: Dezembro/2025
