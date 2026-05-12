# ADR 0004 — Postgres como Destino Pós-Cutover

**Status**: Aceito (intenção)
**Data**: 2026-05-12
**Decisores**: aojunioro
**Supersede parcialmente**: nenhum
**Relaciona-se com**: ADR 0002 (Reuso da Base Legada), ADR 0003 (Stack Detalhada)

---

## Contexto

O SMID 10 foi bootstrapped sobre os quatro bancos MySQL do SMID 8.x (`smid`, `permission`, `log`, `communication`) por exigência da coexistência (ADR 0002). Durante a coexistência (Fases 0 a 7) trocar de SGBD é inviável: o SMID 8.x continua escrevendo em MySQL e qualquer outro destino precisaria de CDC bidirecional sobre quatro schemas com encodings mistos — risco operacional desproporcional ao ganho técnico nesta fase.

No entanto, ao chegar no **cutover (Fase 8)** o SMID 8.x é desligado e a restrição cai. Nesse ponto, Postgres oferece vantagens concretas que o roadmap pós-cutover poderá explorar.

---

## Decisão

1. **Durante as Fases 0 a 7**: o SMID 10 permanece em MySQL, conforme ADR 0002.
2. **Na Fase 8 (cutover)**: o destino preferencial passa a ser **PostgreSQL 16+**, em **schema único** (`smid10`) substituindo os quatro bancos legados, salvo justificativa contrária documentada em novo ADR.
3. **Preparação contínua a partir da Fase 0.2**: o código é escrito de forma a tornar a troca de driver mecânica, sem refactor de domínio.

### Critérios de Gatilho para Executar a Migração

A migração para Postgres só é iniciada quando **todos** forem verdadeiros:

- SMID 8.x está fora de produção (sem escritas há ≥ 7 dias).
- 100% das telas/SPECs migradas para SMID 10 e validadas em produção.
- Backup completo dos 4 bancos MySQL congelado e arquivado.
- Plano de rollback documentado (retorno ao snapshot MySQL em ≤ 4 h).

### Estratégia de Migração (esboço)

- **Ferramenta**: `pgloader` para a carga inicial.
- **Janela**: única, sem CDC, aproveitando que o legado já está desligado.
- **Mapeamentos críticos** a validar antes:
  - `TINYINT(1)` → `boolean`
  - `DATETIME` → `timestamptz` (com revisão de timezone)
  - Encodings `latin1`/`utf8mb4` mistos → `UTF-8` com sanitização prévia
  - `AUTO_INCREMENT` → `IDENTITY`
  - Collations case-insensitive → `citext` ou índices funcionais com `lower()`
- **Triggers e views legados**: reavaliar caso a caso; muitos podem virar regras de domínio em Go.
- **`SystemChangeLog`** e demais tabelas de auditoria: migradas como estão; histórico preservado.

---

## Medidas Preparatórias Adotadas Agora (a partir da Fase 0.2)

Para que o cutover seja uma troca de driver e não uma reescrita:

1. **Repositórios atrás de interfaces**: nenhum `handler` ou `service` de domínio importa `database/sql` direto. Acesso a dados sempre via `XxxRepository` em `internal/domain/<dominio>`.
2. **SQL portável por convenção**: evitar dialectismos MySQL quando uma forma padrão equivalente existir.
   - Evitar: `GROUP_CONCAT`, `ON DUPLICATE KEY UPDATE`, `LIMIT x, y`, `IFNULL`, backticks como delimitador de identificador, `UNSIGNED`, `ENUM` inline.
   - Preferir: `STRING_AGG`-equivalente em código quando necessário, `INSERT ... RETURNING`-style isolado em queries `sqlc`, `LIMIT n OFFSET m`, `COALESCE`, aspas duplas em identificadores quando inevitável.
3. **Queries não portáveis ficam confinadas** a arquivos `.sql` do `sqlc`, nunca espalhadas em strings inline pelo código Go. Assim a substituição é localizada.
4. **Tipos de domínio independentes do driver**: `time.Time` (não `mysql.NullTime`), `bool` (não `int8`), `*string` para nuláveis em vez de `sql.NullString` na borda externa do domínio.
5. **Testes de repositório usam `testcontainers`** (futuro, Fase 1.x): permitem trocar a imagem de `mysql` para `postgres` no dia da virada com um único parâmetro.

Essas medidas **não custam tempo extra** no MySQL atual — são apenas higiene de código.

---

## Consequências

### Positivas

- Decisão estratégica registrada; não vira surpresa em 2026/2027.
- Equipe tem norte ao escolher entre construções SQL equivalentes.
- Pós-cutover, ganhos imediatos esperados: JSONB, CTEs recursivas, índices parciais/funcionais, `FOR UPDATE SKIP LOCKED` para filas, `pg_trgm` para busca, `tstzrange` para janelas de visita.
- Schema único elimina a complexidade de quatro pools, transações cruzadas e DSNs separadas (mantemos a abstração `Alias` por compatibilidade, podendo virar `Schema` no Postgres).

### Negativas

- Pequeno custo cognitivo: desenvolvedores precisam evitar dialectismos MySQL conscientemente.
- Risco de migração mal feita se executada sem auditoria de tipos.
- Dependência adicional no toolchain (`pgloader`) no dia do cutover.

### Neutras

- `sqlc`, `migrate`, `database/sql`, `slog`, Echo: todos suportam ambos os SGBDs sem alteração estrutural.

---

## Limites

| O que **PODE** fazer agora | O que **NÃO PODE** fazer agora |
|---|---|
| Escolher construções SQL portáveis quando equivalentes | Migrar qualquer banco para Postgres |
| Documentar mapeamentos de tipo MySQL→Postgres em SPECs | Introduzir Postgres em paralelo durante a coexistência |
| Manter `internal/db/` com abstração genérica | Adicionar driver `pgx` ao `go.mod` antes da Fase 7 |
| Escrever migrations `golang-migrate` em dialeto neutro quando possível | Reescrever ADR 0002 |

---

## Alternativas Consideradas

- **CDC bidirecional MySQL↔Postgres durante coexistência** (Debezium + Kafka): descartado por complexidade operacional, latência e risco de divergência em quatro schemas heterogêneos.
- **Postgres como réplica somente-leitura via CDC**: descartado porque o SMID 10 precisa escrever, não só ler.
- **Permanecer em MySQL após o cutover**: viável; reavaliar na Fase 7 com dados reais de operação. Caso vença, este ADR é rebaixado a "Rejeitado" via novo ADR.

---

## Referências

- ADR 0002 — Reuso da Base Legada Durante Coexistência
- ADR 0003 — Stack Detalhada
- `docs/specs/SPEC_DATABASE.md`
- `docs/handoff/HANDOFF_BOOTSTRAP.md` §9 (riscos de encoding e compatibilidade)
