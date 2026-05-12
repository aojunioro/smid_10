# Documentação Técnica: Otimização Mobile da Aba de Visitas

**Data:** 17/12/2025  
**Contexto:** Otimização da experiência mobile na gestão de visitas dentro do `LeadsForm` (SMID Haiflex).

## 1. Visão Geral
A implementação padrão utilizando `TFieldList` (tabela com edição inline) mostra-se inviável para dispositivos móveis devido à necessidade de rolagem horizontal excessiva e dificuldade de interação com pequenos inputs.

A solução adotada foi criar uma **estratégia de dupla visualização**:
- **Desktop:** Mantém o `TFieldList` tradicional (tabela).
- **Mobile:** Renderiza uma lista de **Cartões (Cards)** verticais, com botões de ação otimizados e edição via **Modal (Pop-up)**.

## 2. Arquitetura da Solução

### 2.1. Arquivos Envolvidos
1.  `app/control/leads/LeadVisitaDetalhes.php`: Classe principal (CRUD inline).
2.  `app/control/leads/LeadVisitaMobileActions.php`: **Novo controlador** dedicado para ações mobile.
3.  `app/templates/css/mobile-premium.css`: Customizações de estilo específicas para mobile.

### 2.2. Detalhes de Implementação

#### A. Separação de Views (LeadVisitaDetalhes.php)
Utilizamos classes utilitárias do Bootstrap para alternar a visibilidade:
```php
// Desktop (TFieldList)
$this->fieldlistWrapper->class = 'd-none d-md-block';

// Mobile (Cards)
$this->mobileListWrapper->class = 'd-block d-md-none mobile-visits-wrapper';
```
O método `buildMobileList()` foi criado para iterar sobre as visitas e gerar o HTML dos cartões, utilizando a classe `btn-default` para suportar nativamente os temas Dark e Light.

#### B. Edição via Modal e Controlador Dedicado (LeadVisitaMobileActions.php)
Criamos um controlador estático `LeadVisitaMobileActions` para encapsular o ciclo de vida da transação e da instância, resolvendo problemas de "Crash" no construtor e "Tela Branca" ao salvar.

```php
class LeadVisitaMobileActions extends TPage
{
    public static function onEdit($param)
    {
        try {
            TTransaction::open('smid'); 
            $page = new LeadVisitaDetalhes(['attachToPage' => false]);
            $page->setLeadId($param['lead_id']);
            $page->loadOptions([]);
            $page->openMobileForm($param, 'Editar Visita');
            TTransaction::close();
        } catch (Exception $e) { /* ... */ }
    }
}
```

#### C. Ajustes de Layout (mobile-premium.css)
Para otimizar o espaço em tela nos dispositivos móveis, foi necessário remover o padding padrão aplicado pelo template ao container principal (`card-body`) da aba.

Utilizamos o seletor `:has()` para focar especificamente no wrapper das visitas mobile:
```css
/* Remove padding do container pai APENAS quando contiver nossa lista mobile */
@media (max-width: 768px) {
    .card-body:has(.mobile-visits-wrapper) {
        padding: 0 !important;
    }
}
```

## 3. Desafios e Correções Realizadas

### 3.1. Tela Branca ao Salvar
**Solução:** Adicionado parâmetro `static=1` na `TAction` do botão Salvar para forçar requisição AJAX.

### 3.2. Erro "Sem transação ativa"
**Solução:** Instanciação controlada via `LeadVisitaMobileActions` com gestão explícita de `TTransaction`.

### 3.3. Layout com Excesso de Espaço e Padding
**Solução:** Implementada regra CSS em `mobile-premium.css` removendo o padding do elemento pai `.card-body` ao detectar a classe `.mobile-visits-wrapper`.

### 3.4. Suporte a Dark/Light Mode
**Solução:** Substituídas classes fixas como `btn-light`, `bg-light`, e `text-dark` pela classe semântica `btn-default`. Esta classe, configurada no `theme.css`, ajusta automaticamente as cores de fundo, borda e texto conforme o tema ativo, garantindo consistência visual entre desktop e mobile.

## 4. Como Manter

1.  **Novos Campos:** Atualizar `buildFieldList` (Desktop) e `openMobileForm` (Mobile).
2.  **Estilos:** Ajustes visuais específicos devem ser feitos em `app/templates/css/mobile-premium.css` utilizando a classe `.mobile-visits-wrapper` como escopo.
3.  **Lógica:** Manter lógica de negócio (`onCreateVisit`, `onUpdateVisit`) centralizada na classe principal `LeadVisitaDetalhes`.

---
**Arquivos Criados/Modificados:**
- `app/control/leads/LeadVisitaMobileActions.php` (Novo)
- `app/control/leads/LeadVisitaDetalhes.php` (Modificado)
- `app/templates/css/mobile-premium.css` (Modificado)
