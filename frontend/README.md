# Frontend SMID 10 (Next.js + shadcn/ui)

Aplicação SPA do SMID 10 escrita em Next.js (App Router) com TypeScript e shadcn/ui.

---

## 1. Stack

Ver `docs/adrs/0003-stack-detalhada.md`.

Resumo:
- Next.js 14+ (App Router)
- TypeScript
- shadcn/ui + Radix
- TailwindCSS
- react-hook-form + zod
- TanStack Query + TanStack Table
- dnd-kit (kanban)
- lucide-react
- sonner (toast)
- date-fns + react-day-picker
- next-themes (dark/light)

---

## 2. Bootstrap (Fase 0)

### 2.1 Opção A — Partir do template `shadcn-admin`

Acelera Fase 0 em 2–3 meses (layout, sidebar, theming prontos).

```bash
cd frontend
git clone https://github.com/satnaing/shadcn-admin.git .
rm -rf .git
pnpm install
pnpm dev
```

Adaptar:
- Branding (logo, cores)
- Configurar `lib/api/` apontando para `http://localhost:8080`
- Substituir mocks por chamadas reais

### 2.2 Opção B — Scaffold do zero

```bash
cd frontend
pnpm create next-app@latest . --typescript --tailwind --app --eslint --src-dir=false --import-alias="@/*"
pnpm dlx shadcn@latest init
pnpm add @tanstack/react-query @tanstack/react-table
pnpm add react-hook-form @hookform/resolvers zod
pnpm add date-fns react-day-picker
pnpm add @dnd-kit/core @dnd-kit/sortable
pnpm add lucide-react sonner next-themes
```

Depois adicionar componentes shadcn conforme precisar:

```bash
pnpm dlx shadcn@latest add button input form sheet table dialog dropdown-menu
```

### 2.3 Variáveis de ambiente

`.env.local`:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### 2.4 Rodar

```bash
pnpm dev
```

App em `http://localhost:3000`.

---

## 3. Estrutura

```
frontend/
├── app/
│   ├── (auth)/
│   │   └── login/page.tsx
│   ├── (app)/                   ← rotas autenticadas
│   │   ├── layout.tsx           ← layout com sidebar + topbar
│   │   ├── leads/
│   │   ├── visitas/
│   │   ├── pedidos/
│   │   └── ...
│   ├── api/                     ← route handlers Next (apenas para edge cases)
│   ├── layout.tsx
│   └── globals.css
├── components/
│   ├── ui/                      ← shadcn/ui (gerado)
│   └── smid/                    ← componentes próprios
│       ├── side-panel.tsx       ← cortina lateral (wrapper de Sheet)
│       ├── data-table.tsx       ← TanStack Table + shadcn
│       ├── date-range-field.tsx ← date picker com presets
│       ├── multi-combo.tsx      ← multi-select com busca
│       └── kanban-board.tsx     ← dnd-kit + cards
├── lib/
│   ├── api/                     ← clients REST por domínio
│   │   ├── client.ts            ← fetch wrapper com auth
│   │   ├── leads.ts
│   │   ├── visitas.ts
│   │   └── ...
│   ├── auth/
│   │   ├── auth-context.tsx
│   │   └── use-auth.ts
│   ├── utils.ts                 ← clsx, twMerge, helpers
│   └── validators/              ← schemas zod por domínio
├── hooks/
├── types/                       ← types TS espelhando JSON da API
└── public/
```

---

## 4. Padrões Importantes

### 4.1 Cortina Lateral (Side Panel)

Substitui `TSidePanel` do legado. Implementação:

```tsx
// components/smid/side-panel.tsx
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet"

export function SidePanel({ open, onOpenChange, title, children }) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-2xl overflow-y-auto">
        <SheetHeader>
          <SheetTitle>{title}</SheetTitle>
        </SheetHeader>
        {children}
      </SheetContent>
    </Sheet>
  )
}
```

Regras (de `SPEC_UX_UI.md`):
- Abre à direita em desktop, full em mobile
- Fechamento por backdrop, ESC ou botão
- Profundidade = 1 (cortina não abre outra cortina)

### 4.2 Busca Rápida

Server-side por padrão. Ver `docs/legacy-reference/BUSCA_RAPIDA.md` para árvore de decisão.

### 4.3 TDateRangeField

Componente próprio com presets (today, yesterday, last7days, etc.) — cálculo no backend, frontend só envia o preset ou intervalo.

### 4.4 Mobile-first

Classes Tailwind começam pelo mobile e crescem:

```tsx
<div className="px-4 md:px-8 lg:px-12">
```

Validar nos breakpoints 320, 375, 414, 428, 768 e desktop.

### 4.5 Tema dark/light

Via `next-themes`. Cores e variáveis em `tailwind.config.ts` + `globals.css` conforme padrão shadcn.

---

## 5. Convenções

Ver `AGENTS.md` na raiz, seção 5.

---

## 6. Roadmap do Frontend

| Fase | Entregável |
|------|-----------|
| 0.1 | Scaffold + tema dark/light + layout app/auth |
| 0.2 | Cliente HTTP com interceptor de auth |
| 0.3 | Componentes próprios em `components/smid/` |
| 1.1 | Tela de login conectada à API |
| 1.2 | Layout autenticado com sidebar |
| 1.3 | Seletor de unidade |
| 2.1 | Listagem de leads |
| 2.2 | Formulário de lead (cortina) |
| ... | (ver SPECs por domínio) |

---

## 7. Build e Deploy

```bash
pnpm build
pnpm start
```

Para deploy em Vercel, Netlify ou container Docker:

```bash
docker build -t smid10-frontend .
```

(Dockerfile a ser criado na Fase 0.)
