# ADR 0001 — Adoção de Backend Go + Frontend Next.js (Cenário 2)

**Status**: Aceito
**Data**: 2026-05-12
**Decisores**: aojunioro

---

## Contexto

O SMID 8.x atual é um monolito PHP construído sobre Adianti Framework 8.4. Para evoluir o sistema sem amarras de framework e abrir caminho para apps mobile e integrações, foi avaliada a melhor combinação de stack para reimplementação.

Duas alternativas foram consideradas:

- **Cenário 1**: Go + templ + HTMX + Tailwind (frontend embutido, SSR)
- **Cenário 2**: Go (API REST) + Next.js + shadcn/ui (SPA separada)

---

## Decisão

Adotar **Cenário 2**: Go como API REST e Next.js + shadcn/ui como aplicação SPA separada.

---

## Consequências

### Positivas

- Casa exatamente com o contrato já documentado em `SPEC_REST_API.md`
- Ecossistema React/Next/shadcn é o mais maduro do mercado para painéis admin modernos
- Cortina lateral, kanban, tabelas avançadas e date pickers premium prontos via shadcn/ui + Radix
- Frontend e backend podem evoluir em paralelo por times distintos
- Mesma API serve futuro app mobile nativo (React Native ou Flutter)
- Mobile premium (responsivo, touch-first) é nativo do stack Tailwind + Radix
- Template `shadcn-admin` (open source, MIT) acelera Fase 0 em 2-3 meses

### Negativas

- Dois repositórios/deploys (backend e frontend) — mitigado com monorepo `smid_10/`
- Build chain JavaScript (Node, npm) exigida no frontend
- Equipe precisa conhecer TypeScript/React, não só Go

### Neutras

- SEO não é prioridade (sistema admin), mas Next.js entrega SSR caso necessário

---

## Alternativas Descartadas

### Cenário 1 (templ + HTMX)

Descartado porque:
- Ecossistema de admin UI ainda imaturo em Go
- Cortina lateral, kanban drag-and-drop e gestures mobile exigiriam código manual extenso
- Padrão de mobile premium (breakpoints, transições) é trivial em Tailwind/shadcn e custoso em HTMX puro

### Frameworks alternativos (Rails, Laravel, NestJS)

Descartados porque:
- Objetivo é sair de monolitos full-stack PHP
- Go entrega performance, deploy de binário único e concorrência nativa superiores

---

## Referências

- `SPEC_REST_API.md` — contrato da API a ser implementada
- `SPEC_UX_UI.md` — padrões de UX que o frontend deve preservar
- [shadcn/ui](https://ui.shadcn.com/)
- [shadcn-admin](https://github.com/satnaing/shadcn-admin)
