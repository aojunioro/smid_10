# SPEC_DATABASE.md — Arquitetura de Dados e Migrações

> Este SPEC é o contrato canônico da camada de persistência do SMID. Qualquer reimplementação (MySQL, PostgreSQL, SQL Server, etc.) deve preservar os quatro bancos lógicos, seus aliases e a ordem de migração documentada.

> ## ⚑ Regra de Ouro do Schema
>
> - **Para projeto novo (greenfield)**: schema vem dos SPECs + `app/database/*.sql` + `app/database/scripts/smid padrao/`.
> - **Para manutenção do SMID atual** ou **importação de dados de cliente existente**: use o restante de `app/database/scripts/` conforme `README_SCRIPTS.md`.
>
> Aplicação detalhada nas seções 1.3, 2 e 9 deste SPEC.

---

## 1. Visão Geral

### 1.1 Multi-Banco Obrigatório

O SMID particiona dados em **4 bancos lógicos** com aliases canônicos. Qualquer serviço, model ou migração deve declarar explicitamente o alias em uso.

| Alias | Banco físico (exemplo) | Conteúdo |
|-------|------------------------|----------|
| `smid` | `<cliente>_smid` | Dados de negócio (leads, visitas, pedidos, produtos, financeiro, KM, metas, tarefas, televendas) |
| `permission` | `<cliente>_permission` | Autenticação, usuários, grupos, papéis, programas, unidades, regras granulares |
| `log` | `<cliente>_log` | Logs técnicos (acesso, request, SQL, change-log, schedule, webhooks) |
| `communication` | `<cliente>_communication` | Mensagens, notificações, posts, wiki, drive, agenda, jobs |

### 1.2 Aliases Auxiliares

| Alias | Aponta para | Uso |
|-------|-------------|-----|
| `main` | `smid` | Compatibilidade com código legado |
| `principal` | `smid` | Compatibilidade com código legado |

### 1.3 Arquivos SQL de Referência

Localizados em `app/database/` no projeto atual. **Para um projeto novo (greenfield), use-os apenas como referência para derivar o schema final junto com os SPECs**:

| Arquivo | Finalidade | Uso em projeto novo |
|---------|-----------|---------------------|
| `permission.sql` | Schema base da autenticação/autorização | Referência direta |
| `permission-update.sql` | Deltas históricos de `permission` | Consolidar com `permission.sql` |
| `communication.sql` | Schema base da comunicação | Referência direta |
| `communication-update.sql` | Deltas históricos de `communication` | Consolidar com `communication.sql` |
| `log.sql` | Schema base dos logs | Referência direta |
| `scripts/smid padrao/indigita_adianti_teste.sql` | Schema base de negócio | Referência direta |
| `scripts/smid padrao/hist_ocorrido.sql` | Schema de histórico | Referência direta |
| `scripts/smid padrao/v_calendario_visitas.sql` | View canônica do calendário | Referência direta |

> **Importante**: a tabela do tipo "single source of truth" do schema para um projeto novo são os próprios SPECs (entidades, atributos, relacionamentos). Os `.sql` acima ajudam a confirmar tipos, índices e collations da implementação atual.

---

## 2. Diretório de Scripts (`app/database/scripts/`)

> **Este diretório é REGISTRO HISTÓRICO da migração do legado para o SMID 8.x.** Não é receita para greenfield. Contem scripts `canonical`, `operacional-producao` (clientes específicos) e `legado` misturados.

### 2.1 Quando usar

| Cenário | Usar scripts? |
|---------|----------------|
| Novo projeto greenfield em qualquer linguagem | **Não**. Derive schema dos SPECs + arquivos base do item 1.3 |
| Reimplementação que precisa importar dados de um cliente SMID existente | **Sim, como referência** para entender shape final do dado |
| Manutenção do SMID 8.x atual (PHP) | **Sim**, conforme `README_SCRIPTS.md` |

### 2.2 Organização (13 categorias sequenciais)

| Ordem | Categoria | Conteúdo |
|-------|-----------|----------|
| 1 | `1 - unidades` | Schema e seed de unidades/filiais |
| 2 | `2 - permissoes` | Usuários, grupos, programas, papéis, scopes, sync de legado |
| 3 | `3 - leads` | Evoluções da tabela `leads` |
| 4 | `4 - financeiro` | Contas a pagar/receber e categorias |
| 5 | `5 - historico` | Histórico, motivos, ocorridos, fotos |
| 6 | `6 - despesa extra` | Despesas extras de representantes |
| 7 | `7 - km_rodado` | Reembolso de KM, lotes, trechos, config |
| 8 | `8 - planilha_agendamento` | Horários padrão de agendamento |
| 9 | `9 - metas` | Tabelas de metas |
| 10 | `10 - regras_granulares` | Regras condicionais e ações |
| 11 | `11 - categorizacao de leads` | `lead_status_categoria` + views |
| 12 | `12 - categorizacao de pedidos` | `ped_status_categoria` + views |
| 13 | `13 - televendas` | Contatos, fila e status do Televendas |

### 2.3 Política de Classificação Interna

| Classe | Semântica |
|--------|-----------|
| `canonical` | Estrutura estável; útil para confirmar shape em projeto novo |
| `operacional-producao` | Ajuste pontual de cliente específico; **não** replicar em greenfield |
| `legado` | Histórico/rollback; **não** usar |

### 2.4 Regras (somente para manutenção do SMID atual)

1. Executar categorias em ordem numérica
2. Dentro de cada categoria, executar em ordem lexicográfica
3. Scripts destrutivos têm **rollback dedicado** (ex.: `2.17-rollback-*.sql`)
4. Scripts `canonical` são **idempotentes**

---

## 3. Convenções de Schema

### 3.1 Timestamps de Domínio (`smid`)

| Coluna | Semântica |
|--------|-----------|
| `criado_em` | Criação do registro |
| `alterado_em` | Última alteração |
| `excluido_em` | Soft delete (NULL = ativo) |

### 3.2 Escopo de Unidade

| Coluna | Semântica |
|--------|-----------|
| `unidd_id` | Unidade proprietária do registro (NULL = visível a todos) |

### 3.3 Soft Delete

- Regra geral: dados de negócio usam `excluido_em IS NULL` para representar registros ativos
- Listagens devem filtrar `excluido_em IS NULL` por padrão
- Exceções: tabelas de auditoria, logs, pivôs técnicos

### 3.4 Chaves

| Convenção | Uso |
|-----------|-----|
| `id` | PK serial/auto em todas as tabelas |
| `<entidade>_id` | FK seguindo nome singular (ex.: `lead_id`, `vis_id`, `ped_id`) |
| `IDPOLICY = 'serial'` | Padrão em novas tabelas |
| `IDPOLICY = 'max'` | Tabelas legadas que preservam ID manual (ex.: `midias`) |

### 3.5 Categorização Funcional

Campos de status críticos (`lead_status`, `ped_status`) têm tabela de **categoria funcional** paralela:
- `lead_status_categoria` mapeia status físicos para categorias funcionais (NOVO, QUALIFICADO, AGENDADO, VENDIDO, PERDIDO, …)
- `ped_status_categoria` mapeia status de pedidos (PROCESSANDO, APROVADO, ENTREGA, ENTREGUE, CANCELADO)

**Invariante**: relatórios e regras de negócio devem consumir categorias funcionais via helpers canônicos (`LeadStatusHelper`, `PedStatusHelper`), nunca IDs fixos.

---

## 4. Views Canônicas

| View | Banco | Propósito |
|------|-------|-----------|
| `v_calendario_visitas` | `smid` | Calendário operacional de visitas |
| `v_televendas_fila` | `smid` | Fila consolidada do Televendas |
| Views de visitas válidas/vendidos | `smid` | Suporte a relatórios de evolução |

Views são recriadas por scripts nas categorias 11, 12 e 13.

---

## 5. Transações

### 5.1 Regras Obrigatórias

1. Todo acesso a banco deve abrir transação explícita: `TTransaction::open('<alias>')`
2. Todo controlador deve fechar em `try/catch/finally` equivalente
3. Serviços que cruzam alias devem abrir/fechar por alias sequencialmente
4. Nunca manter duas transações abertas simultaneamente no mesmo alias

### 5.2 Padrão Cross-Base

```
Exemplo: autenticar e carregar frontpage

1. TTransaction::open('permission')
2. Carrega SystemUser, SystemGroups, SystemUnits
3. TTransaction::close()

4. TTransaction::open('smid')
5. Carrega dashboards e permissões de unidade
6. TTransaction::close()

7. TTransaction::open('log')
8. Registra SystemAccessLog
9. TTransaction::close()
```

---

## 6. Portabilidade

### 6.1 Adaptação para Outro SGBD

Para PostgreSQL/SQL Server:
- Substituir `ENGINE=InnoDB` por equivalente
- Converter `TINYINT` em `SMALLINT` ou `BOOLEAN`
- Converter `DATETIME` em `TIMESTAMP`
- Rever índices fulltext (MySQL-específico)
- Ajustar `CONVERT()` e collations específicas (ver services de relatórios)

### 6.2 Aliases Obrigatórios

A aplicação consumidora deve **sempre** resolver 4 conexões lógicas. Mesmo que fisicamente sejam o mesmo servidor, o pool de conexões deve manter a separação para preservar a arquitetura.

---

## 7. Segurança de Dados

| Risco | Mitigação |
|-------|-----------|
| SQL injection | Parametrização obrigatória (`TCriteria`, `TFilter`, prepared statements) |
| Exposição de credenciais | Credenciais apenas em `app/config/app_database.php`, nunca versionadas para remotos |
| Alteração descontrolada | Todas as migrações passam por `app/database/scripts/` numeradas |
| Perda de rastro | Modelos críticos ativam `SystemChangeLogTrait` para auditoria |
| Vazamento cross-unit | Filtrar `unidd_id` na camada de consulta conforme sessão |
| Exposição de PII | Logs não devem persistir senha/token/CVV em claro |

---

## 8. Referências à Implementação Atual

| Artefato | Localização | Natureza |
|----------|-------------|----------|
| Configuração de conexões | `app/config/app_database.php` | Referência |
| Função canônica | `getDatabaseConfig($database)` | Referência |
| Schema base por banco | `app/database/*.sql` | **Útil para greenfield** |
| Schema base de negócio | `app/database/scripts/smid padrao/` | **Útil para greenfield** |
| Migrações históricas | `app/database/scripts/1..13` | Registro histórico |
| Guia oficial de migração | `app/database/scripts/README_SCRIPTS.md` | Para manutenção do SMID atual |
| Anotações históricas | `app/database/scripts/anotacoes_migracoes.md` | Histórico |
| Helper read-only | `scripts/dev/mcp_mysql_readonly.sh` | Diagnóstico |

---

## 9. Caminho Recomendado para Projeto Novo (Greenfield)

```
1. Ler SPEC_INDEX.md (visão geral)
2. Ler SPEC_DATABASE.md (este arquivo)
3. Para cada banco lógico (smid, permission, log, communication):
   a. Inspecionar o `.sql` base correspondente em app/database/
   b. Ler os SPECs de domínio que persistem nesse banco
   c. Derivar schema final juntando: SPEC.entidades + tipos do .sql base
4. Implementar convenções (timestamps, soft delete, unidade)
5. Implementar tabelas de categorização funcional (lead_status_categoria, ped_status_categoria)
6. Criar views canônicas equivalentes
7. NUNCA replicar scripts operacional-producao do legado
```

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
