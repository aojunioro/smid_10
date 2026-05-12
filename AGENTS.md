# AGENTS.md — SMID 10

Guia para agentes de IA e desenvolvedores trabalhando no SMID 10 (Go + Next.js). Este é o documento canônico de operação. Em conflito com qualquer outro, este prevalece.

---

## 1. Contexto Rápido

- **Sucessor** do SMID 8.x (Adianti/PHP), reescrito em **Go (backend) + Next.js/shadcn (frontend)**
- **Modo**: reuso da base MySQL do legado (`smid`, `permission`, `log`, `communication`)
- **Coexistência**: SMID 8.x e SMID 10 operam simultaneamente até o cutover
- **Fonte de verdade**: os 24 SPECs em `docs/specs/`

---

## 2. Princípios Não Negociáveis

1. **SPECs primeiro**: leia `docs/specs/SPEC_INDEX.md` antes de qualquer implementação
2. **Schema legado intocável** durante a coexistência: novas tabelas, sim; alterar tabelas existentes, **não**
3. **Multi-banco obrigatório**: respeitar os 4 aliases (`smid`, `permission`, `log`, `communication`)
4. **Invariantes**: nunca quebrar os invariantes listados em `SPEC_INDEX.md` seção 8
5. **Mobile premium**: validar nos breakpoints 320, 375, 414, 428, 768 e desktop
6. **REST + JWT**: comunicação backend ↔ frontend sempre via API documentada
7. **pt-BR**: documentação, interface, mensagens de erro
8. **Inglês**: código (classes, funções, variáveis, comentários técnicos)
9. **Sem emojis** em código, docs ou commits
10. **Idempotência**: webhooks e jobs devem ser idempotentes por chave externa

---

## 3. Workflow Padrão

```
ANALISAR → PLANEJAR → IMPLEMENTAR → TESTAR → VALIDAR → DOCUMENTAR
```

Para cada demanda:

1. **ANALISAR**: ler SPECs do domínio afetado + ADRs relevantes
2. **PLANEJAR**: definir endpoint REST (se aplicável), telas, contratos
3. **IMPLEMENTAR**: backend antes do frontend; SPECs como referência
4. **TESTAR**: testes unitários + integração; mobile responsivo
5. **VALIDAR**: contra invariantes do SPEC_INDEX
6. **DOCUMENTAR**: atualizar SPEC se houver mudança de regra; ADR se houver decisão arquitetural

---

## 4. Backend (Go)

### 4.1 Estrutura

```
backend/
  cmd/server/main.go         ← entrypoint
  internal/
    config/                  ← carregamento de .env, struct Config
    db/                      ← pools de conexão por alias (smid/permission/log/communication)
    auth/                    ← JWT, hashing, sessão
    domain/                  ← entidades, regras de negócio (por SPEC)
      leads/
      visitas/
      pedidos/
      ...
    http/
      handlers/              ← controllers REST
      middleware/            ← auth, logging, recovery, cors
      routes.go              ← registro de rotas
  pkg/                       ← código reusável fora do internal
```

### 4.2 Regras

- **Cada domínio é um pacote**: `internal/domain/leads`, `internal/domain/visitas`, etc.
- **Handlers finos**: parsing/validação de input + chamada para serviço de domínio + serialização
- **Serviços de domínio**: regras de negócio puras, sem dependência de HTTP
- **Repositórios**: acesso a banco via sqlc (queries tipadas) ou `database/sql`
- **Transações**: abrir/fechar explicitamente por alias, nunca cruzar transações entre aliases
- **Erros**: tipados com sentinel errors ou `errors.Is/As`; nunca expor `error.Error()` ao cliente sem sanitizar
- **Logging**: `slog` estruturado; nunca logar senha/token/PII em claro
- **Contexto**: propagar `context.Context` em todas as camadas

### 4.3 Convenções

| Item | Convenção |
|------|-----------|
| Nome de pacote | `lowercase` (`leads`, `visitas`) |
| Nome de tipo | `PascalCase` (`Lead`, `LeadService`) |
| Nome de função | `PascalCase` se exportada, `camelCase` se não |
| Nome de arquivo | `snake_case.go` ou `nome_servico.go` |
| Erros | `ErrXxx` para sentinel; wrappear com `fmt.Errorf("%w", err)` |
| Testes | `*_test.go` no mesmo pacote; usar `testify/assert` |

---

## 5. Frontend (Next.js + shadcn/ui)

### 5.1 Estrutura

```
frontend/
  app/                       ← App Router (Next.js 14+)
    (auth)/login/page.tsx
    (app)/                   ← rotas autenticadas
      leads/
      visitas/
      pedidos/
      ...
    layout.tsx
    globals.css
  components/
    ui/                      ← shadcn/ui (gerados)
    smid/                    ← componentes próprios reutilizáveis
      side-panel.tsx         ← cortina lateral
      data-table.tsx
      date-range-field.tsx
      multi-combo.tsx
  lib/
    api/                     ← clients REST por domínio
    auth/                    ← guarda de rota, hooks
    utils.ts
  hooks/
  types/                     ← types TypeScript espelhando JSON da API
```

### 5.2 Regras

- **shadcn/ui primeiro**: usar componentes prontos; criar custom apenas em `components/smid/`
- **TanStack Query** para fetch: cache, refetch, optimistic updates
- **react-hook-form + zod**: todos os formulários
- **Mobile-first**: classes Tailwind começam por base (mobile) e crescem para `md:`, `lg:`
- **Cortina lateral**: usar `<Sheet>` do shadcn como base; envolver em wrapper `<SidePanel>` próprio
- **Tema**: `next-themes` com classe `dark` no `<html>`
- **Acessibilidade**: Radix (via shadcn) já entrega; respeitar `aria-*` em customizações
- **Sem emojis** em UI; usar ícones `lucide-react`

### 5.3 Convenções

| Item | Convenção |
|------|-----------|
| Componente | `PascalCase.tsx` |
| Hook | `useXxx.ts` |
| Tipo TS | `PascalCase` |
| Rota | `kebab-case` |
| Variável | `camelCase` |

---

## 6. Banco de Dados

### 6.1 Reuso da Base Legada

- **Schema existente é imutável** durante coexistência
- **Conexões por alias** explicitamente nomeadas no pool
- **Novas features** que precisem de tabelas: criar com prefixo `s10_` ou em schema dedicado
- **Migrations**: `golang-migrate` apenas para deltas do SMID 10

### 6.2 Queries

- Preferir **sqlc** (queries SQL em `.sql`, código tipado gerado em Go)
- Filtrar **sempre** por `unidd_id` quando aplicável
- Filtrar **sempre** por `excluido_em IS NULL` em listagens (soft delete)
- Usar `lead_status_categoria` / `ped_status_categoria` para categorização funcional, **nunca** IDs fixos

---

## 7. Segurança

| Item | Regra |
|------|-------|
| Senhas | Hash + salt; nunca em claro; comparar com `bcrypt` ou compatível com legado |
| JWT | HS256, expiração 8h, refresh dedicado |
| Validação | Backend valida 100% do input, mesmo com validação no frontend |
| SQL injection | Apenas queries parametrizadas; sqlc força isso |
| CORS | Configurado por ambiente; produção restrita |
| Rate limit | Endpoints de autenticação e busca |
| Auditoria | Toda escrita relevante gera entrada equivalente a `SystemChangeLog` |
| Secrets | Apenas em `.env` (gitignored); produção em variáveis de ambiente do orquestrador |

---

## 8. Workflow Git

- Branch `main` é integração; nunca commitar diretamente
- Trabalho em `feature/<modulo>-<descricao>` ou `fix/<descricao>`
- Commits convencionais: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`
- Mensagens de commit em inglês, sem emoji
- Merge para `main` apenas após review e testes

---

## 9. Checklist Antes de Concluir uma Demanda

- [ ] SPECs do domínio lidos e respeitados
- [ ] Invariantes do `SPEC_INDEX.md` validados
- [ ] Testes unitários passando
- [ ] Testes de integração (se aplicável) passando
- [ ] Mobile validado nos breakpoints
- [ ] Sem credenciais ou PII em código/logs
- [ ] Documentação atualizada (SPEC ou ADR) se houve mudança de regra
- [ ] Sem código morto ou comentários TODO sem responsável

---

## 10. Quando Atualizar Este Documento

Sempre que:
- Surgir nova convenção que se aplique a múltiplos domínios
- Mudar um princípio não negociável (raro; requer ADR)
- Adicionar dependência transversal (ex.: nova lib de auth, nova lib de UI)

Atualizar acompanhado de um ADR em `docs/adrs/`.

---

**Versão**: 0.1.0
**Status**: documento vivo
