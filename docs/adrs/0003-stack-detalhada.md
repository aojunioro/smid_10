# ADR 0003 — Stack Detalhada do SMID 10

**Status**: Aceito
**Data**: 2026-05-12
**Decisores**: aojunioro

---

## Contexto

O ADR 0001 fixou o cenário (Go + Next.js). Este ADR detalha as bibliotecas concretas para evitar discussões repetidas e garantir consistência entre os módulos.

---

## Decisões

### Backend

| Camada | Escolha | Alternativa descartada | Motivo |
|--------|---------|------------------------|--------|
| Framework HTTP | **Echo v4** | Gin, Fiber, Chi | Maturidade, middlewares ricos, comunidade |
| Driver MySQL | **go-sql-driver/mysql** | — | Padrão de mercado |
| Query layer | **sqlc** | GORM, ent | Type-safe, sem overhead de ORM, controle de SQL |
| Migrations | **golang-migrate** | atlas, goose | Mais simples, sintaxe SQL pura |
| JWT | **golang-jwt/jwt v5** | paseto | Compatibilidade com firebase/php-jwt do legado |
| Validação | **go-playground/validator v10** | ozzo-validation | Padrão de mercado, tag-based |
| Logging | **log/slog** (stdlib) | zerolog, zap | Stdlib, structured, sem dependência |
| Hashing | **bcrypt** (`golang.org/x/crypto/bcrypt`) | argon2 | Compatibilidade com legado |
| Testes | **testing** + **testify/assert** | gocheck | Padrão de mercado |
| HTTP client | `net/http` + **resty** quando precisar de retries | — | Stdlib primeiro |
| Config | `.env` via **godotenv** + `os.Getenv` | viper | Mínimo viável |
| Cron/Jobs | **robfig/cron v3** | — | Padrão de mercado para schedules |

### Frontend

| Camada | Escolha | Alternativa descartada | Motivo |
|--------|---------|------------------------|--------|
| Framework | **Next.js 14+ (App Router)** | Remix, SvelteKit | Ecossistema React, SSR opcional |
| Linguagem | **TypeScript** | JavaScript puro | Type safety obrigatório |
| UI base | **shadcn/ui** + Radix | MUI, Mantine, Chakra | Copy-paste, controle total, premium |
| Estilo | **TailwindCSS v3+** | CSS modules, styled-components | Mobile-first, performance, padrão shadcn |
| Forms | **react-hook-form + zod** | Formik | DX superior, validação tipada |
| Tabelas | **TanStack Table v8** | AG Grid | Headless, customizável, gratuito |
| Drag-and-drop | **dnd-kit** | react-beautiful-dnd | Touch nativo, ativamente mantido |
| HTTP | **TanStack Query v5** + fetch | SWR, Axios | Cache poderoso, padrão atual |
| Datas | **date-fns** + **react-day-picker** | moment, dayjs | Leve, tree-shakeable |
| Tema | **next-themes** | — | Padrão Next.js para dark/light |
| Ícones | **lucide-react** | heroicons, react-icons | Padrão shadcn |
| Notifications/Toast | **sonner** | react-hot-toast | Padrão shadcn atual |
| Template base | **shadcn-admin** (kiranism / satnaing) | TailAdmin | Open source MIT, alinhado com shadcn/ui |
| Package manager | **pnpm** | npm, yarn | Performance, disk efficiency |

---

## Versões Mínimas

| Ferramenta | Versão |
|-----------|--------|
| Go | 1.22+ |
| Node.js | 20 LTS |
| pnpm | 9+ |
| MySQL | 5.7+ (compatibilidade com legado) |

---

## Consequências

### Positivas

- Stack decidida, sem retrabalho de discussão em cada módulo
- Bibliotecas todas open source, sem custos de licença
- Toda escolha tem comunidade ativa e documentação rica
- Compatibilidade total com legado (bcrypt, JWT HS256)

### Negativas

- Mudanças de stack exigem novo ADR e revisão dos módulos já implementados
- Curva de aprendizado para quem não conhece TypeScript/React + Go simultaneamente

---

## Revisão

Revisar este ADR ao final de cada fase (ver `README.md` seção 1.2) e antes de adicionar nova dependência transversal.

---

## Referências

- ADR 0001 — Cenário 2 (Go + Next.js)
- ADR 0002 — Reuso da base legada
- `SPEC_REST_API.md` — contrato de API
- `SPEC_UX_UI.md` — padrões de UX
