# ADR 0002 — Reuso da Base de Dados Legada Durante Coexistência

**Status**: Aceito
**Data**: 2026-05-12
**Decisores**: aojunioro

---

## Contexto

O SMID 10 substituirá o SMID 8.x (Adianti/PHP) de forma progressiva. Durante o período de transição, ambos precisam operar simultaneamente sobre os mesmos dados, sem que o legado quebre nem o novo perca consistência.

---

## Decisão

O SMID 10 conecta-se **diretamente aos quatro bancos MySQL do SMID 8.x** (`smid`, `permission`, `log`, `communication`), preservando schema, dados e integridade referencial.

### Regras de Coexistência

1. **Schema legado é imutável** durante a coexistência
2. **Novas tabelas exclusivas do SMID 10** usam prefixo `s10_` ou schema dedicado (a definir)
3. **Tabelas de categorização** (`lead_status_categoria`, `ped_status_categoria`) já existem no legado e são reutilizadas
4. **Triggers e views** do legado são preservados
5. **SystemChangeLog**: SMID 10 grava na mesma tabela do legado para auditoria unificada
6. **Sessões e JWT**: SMID 10 emite JWT próprio; sessões PHP do legado permanecem independentes
7. **Cutover**: ao migrar, o SMID 8.x é desligado e o SMID 10 assume; nenhuma migração de dados é necessária

---

## Consequências

### Positivas

- Zero migração de dados
- Rollback instantâneo (basta voltar ao SMID 8.x)
- Validação progressiva por domínio em produção real
- Dados históricos preservados sem transformação
- Usuários, permissões, unidades e configurações compartilhadas

### Negativas

- SMID 10 herda particularidades do schema legado (collations, encodings mistos, IDs sem FK explícita em alguns casos)
- Mudanças de schema do legado (raras, conforme `SMID_PERSONALIZACOES.md`) precisam ser sincronizadas com SMID 10
- Necessidade de manter compatibilidade de tipos (ex.: `TINYINT(1)` legado ↔ `bool` em Go)

### Neutras

- Performance: ambos competem pelo MySQL; monitorar e otimizar com índices se necessário

---

## Limites

| O que **PODE** fazer | O que **NÃO PODE** fazer |
|----------------------|--------------------------|
| Ler qualquer tabela do legado | Alterar schema de tabela legada |
| Escrever em tabelas do legado | Renomear tabelas/colunas legadas |
| Criar tabelas com prefixo `s10_` | Remover triggers/views do legado |
| Criar índices novos | Quebrar contratos de FK existentes |
| Adicionar colunas a tabelas novas | Adicionar colunas em tabelas legadas sem ADR específico |

Toda exceção exige novo ADR documentando o impacto no SMID 8.x.

---

## Referências

- `SPEC_DATABASE.md` — arquitetura multi-banco e convenções
- `docs/legacy-schema/` — DDL de referência
- `docs/legacy-reference/SMID_PERSONALIZACOES.md` — restrições herdadas
