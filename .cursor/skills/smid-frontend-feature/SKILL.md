---
name: smid-frontend-feature
description: Builds SMID 10 frontend screens with Vite, TanStack Router, shadcn/ui, TanStack Query, and the real REST API (apiClient). Use for leads, visitas, pedidos UI, login integration, SidePanel forms, or replacing shadcn-admin mocks.
---

# SMID — Feature frontend

## Pre-leitura

- `docs/specs/SPEC_UX_UI.md`
- SPEC do dominio (ex. `SPEC_LEADS.md`)
- `AGENTS.md` secao 7
- `frontend/src/lib/api/client.ts`

## Estrutura

```
frontend/src/features/<dominio>/
  components/
  hooks/
  api.ts              # wrappers TanStack Query sobre apiClient
frontend/src/routes/_authenticated/<dominio>/
frontend/src/components/smid/   # SidePanel, DataTable, etc.
```

Atualizar `frontend/src/components/layout/data/sidebar-data.ts` com rotas SMID (remover mocks de template quando substituir).

## Auth

1. `POST /api/v1/auth/login` com `{ "login", "password" }` (nao email do template).
2. Persistir `token` no `auth-store`; chamar `apiClient.setAuthToken(token)`.
3. Guard em rota autenticada: redirecionar para `/sign-in` sem token.

## Data fetching

- TanStack Query: `queryKey` por dominio + filtros.
- Erros 401: reset auth + redirect (padrao em `main.tsx`).
- Tipos TS espelhando JSON da API em `features/<dominio>/types.ts`.

## UX obrigatoria

- Mobile-first (320px minimo).
- Lista + cortina lateral (`Sheet`) para criar/editar.
- pt-BR em labels e toasts.
- react-hook-form + zod.

## Verificacao

```bash
cd frontend && pnpm test
pnpm build
```

Testar mentalmente 375px e 768px.

Nao commitar sem pedido do usuario.
