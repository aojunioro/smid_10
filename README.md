# SMID 10 — Reimplementação em Go + Next.js

> Sucessor moderno do SMID 8.x (Adianti/PHP), construído com **Go** no backend e **Next.js + shadcn/ui** no frontend, **reaproveitando os bancos de dados legados** para integração transparente com o sistema atual durante o período de coexistência.

---

## 1. Visão e Modo de Operação

### 1.1 Modo

**Reuso da base de dados legada**: o SMID 10 conecta-se aos mesmos quatro bancos MySQL do SMID 8.x (`smid`, `permission`, `log`, `communication`). Isso permite:

- Operação simultânea SMID 8.x (PHP) ↔ SMID 10 (Go/Next) durante a migração
- Validação progressiva de funcionalidade por domínio
- Rollback imediato a qualquer momento
- Reaproveitamento integral de dados históricos

### 1.2 Estratégia

| Fase | Conteúdo |
|------|----------|
| **0. Bootstrap** | Estrutura, ADRs, dependências, conexões funcionando |
| **1. Plataforma** | Autenticação, autorização, multi-unidade, sessão, JWT |
| **2. Núcleo comercial** | Leads → Visitas → Históricos → Pedidos |
| **3. Canais** | Representantes, Televendas, Produtos |
| **4. Financeiro/incentivos** | Financeiro, Comissões, KM, Metas |
| **5. Pós-venda** | Compras, Suporte |
| **6. Análise** | Relatórios, dashboards |
| **7. Integrações** | Evolution, 3C+, GoTo Connect, jobs |
| **8. Cutover** | Migração definitiva, desativação do SMID 8.x |

### 1.3 Princípios Não Negociáveis

1. O **SPEC é a fonte de verdade**, não o código PHP atual
2. O schema legado **não pode ser quebrado** durante o período de coexistência
3. Toda nova feature deve respeitar os **invariantes** documentados em `docs/specs/SPEC_INDEX.md`
4. **Mobile premium** é requisito de aceitação, não opcional
5. **REST + JWT** para toda comunicação backend ↔ frontend
6. **Idioma**: documentação e interface em pt-BR; código em inglês

---

## 2. Estrutura do Projeto

```
smid_10/
├── README.md                    ← este arquivo
├── AGENTS.md                    ← regras de operação para agentes de IA
├── docs/
│   ├── specs/                   ← 24 SPECs canônicos (contrato)
│   ├── legacy-schema/           ← DDL do legado (referência de schema)
│   │   ├── permission.sql
│   │   ├── communication.sql
│   │   ├── log.sql
│   │   └── smid-padrao/
│   ├── legacy-reference/        ← Padrões UX/UI a preservar
│   │   ├── BUSCA_RAPIDA.md
│   │   ├── TMultiCombo.md
│   │   ├── PERMISSAO_DONO.md
│   │   ├── SMID_PERSONALIZACOES.md
│   │   └── mobile/
│   └── adrs/                    ← Architecture Decision Records
├── backend/                     ← API Go (Echo + sqlc + JWT)
│   ├── cmd/server/
│   └── internal/
│       ├── config/
│       ├── db/
│       ├── auth/
│       ├── domain/
│       └── http/
└── frontend/                    ← App Next.js + shadcn/ui
```

---

## 3. Stack

### 3.1 Backend (Go)

| Camada | Tecnologia | Motivo |
|--------|-----------|--------|
| Linguagem | Go 1.22+ | Performance, deploy simples, concorrência |
| HTTP | [Echo](https://echo.labstack.com/) | Maduro, middlewares ricos, performance |
| DB | `database/sql` + `go-sql-driver/mysql` | Reuso do MySQL legado |
| Queries tipadas | [sqlc](https://sqlc.dev/) | Type-safe sem ORM pesado |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) | Para deltas futuros, não para schema base |
| Validação | [go-playground/validator](https://github.com/go-playground/validator) | Padrão de mercado |
| JWT | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | Equivalente ao firebase/php-jwt |
| Logging | `log/slog` (stdlib) | Estruturado, sem dependências |
| Config | `os.Getenv` + arquivo `.env` | Simples, sem viper |
| Testes | `testing` + `testify` | Stdlib + assertions |

### 3.2 Frontend (Next.js)

| Camada | Tecnologia | Motivo |
|--------|-----------|--------|
| Framework | Next.js 14+ (App Router) | SSR/SSG, ecossistema React maduro |
| Linguagem | TypeScript | Type safety |
| UI | [shadcn/ui](https://ui.shadcn.com/) + Radix | Componentes premium, acessíveis, copy-paste |
| Estilo | TailwindCSS | Mobile-first, performance |
| Forms | react-hook-form + zod | DX excelente, validação |
| Tabelas | TanStack Table v8 | Headless, customizável |
| Drag-and-drop | dnd-kit | Touch-friendly para mobile |
| HTTP | TanStack Query + fetch | Cache, refetch, optimistic UI |
| Tema | next-themes | Dark/light |
| Datas | date-fns + react-day-picker | Leve e precisa |
| Template base | [shadcn-admin](https://github.com/satnaing/shadcn-admin) | Layout, sidebar, theming prontos |

---

## 4. Como Começar

### 4.1 Pré-requisitos

- Go 1.22+
- Node.js 20+ e npm/pnpm
- MySQL 5.7+ (ou conexão para o legado)
- Git

### 4.2 Backend

Ver `backend/README.md` (a ser preenchido na Fase 0).

### 4.3 Frontend

Ver `frontend/README.md` (a ser preenchido na Fase 0).

---

## 5. Documentação Crítica

| O que ler primeiro | Onde |
|--------------------|------|
| Visão geral e índice dos 24 SPECs | `docs/specs/SPEC_INDEX.md` |
| Arquitetura de dados | `docs/specs/SPEC_DATABASE.md` |
| Contrato REST/JWT | `docs/specs/SPEC_REST_API.md` |
| Padrões de UX/UI | `docs/specs/SPEC_UX_UI.md` |
| Regras de operação para IA | `AGENTS.md` |
| Harness Cursor (skills, rules, MCP) | `.cursor/README.md` |
| Continuidade entre sessões | `docs/handoff/CURRENT.md` (serie: `docs/handoff/README.md`) |
| Decisões arquiteturais | `docs/adrs/` |

---

## 6. Status

**Fase atual**: 0 — Bootstrap

**Próximos passos**:
1. Inicializar `go.mod` em `backend/`
2. Implementar conexão MySQL com pool e health check
3. Scaffold `shadcn-admin` em `frontend/`
4. Implementar endpoint `POST /auth/login` e tela de login conforme `SPEC_ADMIN`
5. Validar autenticação end-to-end contra o banco `permission` do legado

---

**Versão**: 0.1.0
**Status**: bootstrap
**Origem**: derivado de `smid_8` (Adianti/PHP) em maio/2026
