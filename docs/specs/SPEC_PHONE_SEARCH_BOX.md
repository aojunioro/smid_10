# SPEC_PHONE_SEARCH_BOX.md

> **Este SPEC descreve um componente de UI/UX e seu contrato de backend. É independente de linguagem ou framework — implemente o comportamento descrito em qualquer stack.**

## 1. Visão Geral

### 1.1 Propósito

O **PhoneSearchBox** é um componente de busca rápida de leads por telefone no menu superior do sistema. Ele localiza um lead pelo número de telefone normalizado e apresenta o resultado em contexto isolado (aba, modal ou rota), respeitando o escopo de unidade do usuário autenticado.

### 1.2 Escopo

| Item | Descrição |
|------|----------|
| Tipo | Componente de UI + endpoint de busca |
| Local | Menu/topbar do layout |
| Público | Usuários autorizados a consultar leads |
| Resultado | Lista de leads matcháveis exibida em contexto isolado |
| Mobile | Ocultável ou adaptável conforme design do layout |

### 1.3 Contrato de Backend (endpoint canônico)

**Endpoint**: `GET /leads/search-by-phone`

| Parâmetro | Tipo | Obrigatório | Descrição |
|----------|------|------------|----------|
| `phone` | string | sim | Telefone com máscara ou só dígitos, mín. 10 dígitos |
| `token` | string | opcional | Token de requisição para cache-bust e auditoria |

**Resposta de sucesso**:
```json
{
  "status": "found",
  "lead_id": 123,
  "nome": "João Silva",
  "fone1": "(11) 99999-1234",
  "unidade_id": 2
}
```

**Resposta sem resultado**:
```json
{ "status": "not_found", "phone": "(11) 99999-1234" }
```

**Resposta restrita por unidade**:
```json
{ "status": "restricted", "message": "Lead pertence a outra unidade" }
```

**Normalização de telefone** (backend):
```
1. Remove todos os não-dígitos
2. Normaliza para 10 ou 11 dígitos
3. Busca em fone1 e fone2 da tabela leads
4. Filtra apenas leads visíveis pela unidade do usuário
```

**Regra de permissão**:
- Lead na unidade do usuário: retorna `found`
- Lead sem unidade (unidade_id IS NULL): retorna `found`
- Lead em outra unidade: retorna `restricted`
- Nenhum lead: retorna `not_found`

---

## 2. Elementos de Interface

| Elemento | Descrição |
|----------|-----------|
| Container | `div.d-flex.align-items-center.hide-mobile` |
| Input | `#phone-search-input` |
| Placeholder | `Busca Fone` |
| Máscara | `(99) 99999-9999` |
| Inputmode | `tel` |
| Largura | 180px |

---

## 3. Fluxo Principal

### 3.1 Busca por Telefone

```
1. Usuário digita telefone no campo Busca Fone
2. Máscara inline formata enquanto digita
3. Usuário pressiona Enter
4. Componente cancela submit/navegação padrão
5. Remove caracteres não numéricos
6. Se telefone tiver menos de 10 dígitos:
   - exibe warning "Digite um telefone válido com DDD"
   - não abre aba
7. Gera token temporal
8. Monta URL:
   index.php?class=LeadListResultSearch
     &method=onPhoneSearch
     &input_phone_search=<telefone>
     &phone_search_token=<timestamp>
     &template=iframe
9. Fecha aba anterior "Busca Fone:" se existir
10. Abre nova aba com label "Busca Fone: <telefone>"
11. Ativa iframe da aba
12. Limpa input após abertura
```

### 3.2 Resultado sem Lead

```
1. LeadListResultSearch detecta ausência de resultado
2. Chama window.phoneSearchEmptyResult(phone)
3. Widget exibe mensagem info: "Fone não encontrado: <phone>"
4. Fecha aba da busca após confirmação
```

### 3.3 Resultado Restrito por Unidade

```
1. Busca encontra lead em unidade não autorizada
2. Chama window.phoneSearchRestrictedResult(phone)
3. Widget exibe mensagem info: "Lead de outra Unidade"
4. Fecha aba da busca após confirmação
```

---

## 4. Regras Funcionais

| Regra | Descrição |
|-------|-----------|
| RN-001 | Campo só aparece em desktop |
| RN-002 | Enter dispara a busca |
| RN-003 | Telefone deve conter ao menos 10 dígitos |
| RN-004 | Nova busca substitui aba anterior de Busca Fone |
| RN-005 | Resultado abre em iframe para preservar navegação principal |
| RN-006 | Callback global fecha aba quando não houver resultado útil |
| RN-007 | Busca final deve respeitar permissões/unidades do usuário |

---

## 5. Decisões Arquiteturais

### ADR-001: Widget Isolado em `app/control/PhoneSearchBox.php`

**Decisão**: Criar componente próprio, não embutir lógica diretamente no template.

**Consequências**:
- Reuso e manutenção mais simples
- Arquivo consta como ponto permitido de customização SMID

### ADR-002: Desktop Only

**Decisão**: Usar classe `hide-mobile`.

**Consequências**:
- Evita poluição da topbar mobile
- Mobile deve buscar telefone por telas responsivas de leads

### ADR-003: Aba Iframe em vez de Cortina

**Decisão**: Abrir `LeadListResultSearch` em aba `template=iframe`.

**Consequências**:
- Preserva página atual
- Permite repetir busca sem perder contexto
- Exige controle de fechamento/ativação de abas

### ADR-004: Callbacks Globais

**Decisão**: Resultado chama `window.phoneSearchEmptyResult` e `window.phoneSearchRestrictedResult`.

**Consequências**:
- Baixo acoplamento com a tela de resultado
- Dependência de nomes globais documentados

---

## 6. Referência à Implementação Atual

| Artefato | Localização |
|----------|-------------|
| Widget | `app/control/PhoneSearchBox.php` |
| Resultado | `LeadListResultSearch::onPhoneSearch` |
| Ponto de template | `app/templates/adminbs5/layout.html` |
| JS/CSS permitido | `app/templates/adminbs5/js/theme.js`, `custom.css`, `theme.css` |

### APIs de Template Usadas

- `Template.createPageTab(label, url)`
- `Template.openPageTab(url)`
- `Template.closePageTab(href)`
- `Template.updatePageTabsLocalStorage()`
- `__adianti_message(type, message, callback)`

---

## 7. Segurança e Privacidade

- Busca deve validar permissões de unidade no backend
- Telefone é dado pessoal; evitar expor em logs desnecessários
- Nunca retornar dados de lead de unidade não autorizada
- Token temporal evita reaproveitamento visual/cache simples, mas não substitui autorização

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
