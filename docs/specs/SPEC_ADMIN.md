# SPEC_ADMIN.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Admin** concentra autenticação, autorização, gestão de usuários, grupos, papéis, programas, unidades, preferências, perfil, suporte administrativo, diagnóstico do sistema e configurações de integrações. Ele é a camada de governança operacional do SMID.

### 1.2 Entidades Principais

| Entidade | Descrição |
|----------|-----------|
| **SystemUser** | Usuário do sistema |
| **SystemGroup** | Grupo/perfil operacional |
| **SystemRole** | Papel/permissão lógica |
| **SystemProgram** | Programa/tela autorizável |
| **SystemUnit** | Unidade/filial |
| **SystemPreference** | Preferências globais |
| **SystemPermission** | Matriz de permissões |
| **SystemConditionalRule** | Regras granulares/condicionais |

### 1.3 Relacionamentos

```
SystemUser N ──── N SystemGroup     (via SystemUserGroup)
SystemUser N ──── N SystemRole      (via SystemUserRole)
SystemUser N ──── N SystemUnit      (via SystemUserUnit)
SystemUser N ──── N SystemProgram   (via SystemUserProgram — override individual)
SystemGroup N ──── N SystemProgram  (via SystemGroupProgram)
SystemGroup N ──── N SystemRole     (via SystemGroupRole)
SystemProgram N ──── N SystemRole   (via SystemProgramMethodRole — por método/ação)
SystemUser 1 ──── 1 SystemUnit      (unidade principal)
```

### 1.4 Schema Canônico (banco `permission`)

| Tabela | Chave | Campos relevantes |
|--------|-------|-------------------|
| `system_users` | `id` | `login`, `name`, `email`, `password`, `system_unit_id`, `active`, `frontpage_id`, `equipe_id`, `phone`, `address`… |
| `system_groups` | `id` | `name`, `frontpage_id` |
| `system_roles` | `id` | `name` |
| `system_programs` | `id` | `name`, `controller`, `description` |
| `system_units` | `id` | `name`, `parent_id` (hierarquia) |
| `system_user_group` | `id` | `system_user_id`, `system_group_id` |
| `system_user_role` | `id` | `system_user_id`, `system_role_id` |
| `system_user_unit` | `id` | `system_user_id`, `system_unit_id` |
| `system_group_program` | `id` | `system_group_id`, `system_program_id` |
| `system_group_role` | `id` | `system_group_id`, `system_role_id` |
| `system_program_method_role` | `id` | `system_program_id`, `method_name`, `system_role_id` |
| `system_preferences` | `id` | `attribute`, `value` |

---

## 2. Glossário de Negócio

| Termo | Definição |
|-------|-----------|
| **Usuário** | Pessoa autenticada por login/senha |
| **Grupo** | Perfil macro de acesso, ex.: Admin, Atendente, Representante |
| **Papel** | Regra lógica reutilizável de autorização |
| **Programa** | Classe/tela do sistema sujeita a permissão |
| **Unidade** | Escopo operacional dos dados |
| **Regra Granular** | Condição de UI/campo baseada em contexto |
| **Preferência** | Configuração global da aplicação |

---

## 3. Fluxos Principais

### 3.1 Login e Sessão

```
1. Usuário informa login e senha em LoginForm
2. ApplicationAuthenticationService.authenticate(login, password)
   a. Abre TTransaction('permission')
   b. SystemUser::validate(login)
   c. Se auth_service configurado (ex: LDAP), delega autenticação
   d. Caso contrário, SystemUser::authenticate(login, password)
   e. Carrega unit, frontpage e grupos antecipadamente
3. loadSessionVars(user):
   - session('login') = login
   - session('userid') = id
   - session('username') = nome
   - session('system_unit_id') = unidade principal
   - session('system_unit_ids') = todas as unidades
   - session('frontpage') = controller da página inicial
   - session('system_group_ids') = grupos do usuário
4. JWT opcional emitido para REST API
5. ApplicationSessionHeartbeatService mantém sessão ativa
6. AccessLog registra login, IP, resultado
```

### 3.2 Multi-Unidade

```
1. Usuário pertence a uma ou mais SystemUnit
2. system_unit_id na sessão = unidade corrente
3. ApplicationAuthenticationService.setUnit(unit_id) troca unidade
4. Módulos de negócio filtram por system_unit_ids da sessão
5. Dados sem unidade (unidd_id IS NULL) geralmente visíveis para todos
```

### 3.3 Regras Granulares (SystemConditionalRule)

```
1. Admin define regras condicionais em SystemRegraGranularForm
2. Regra associa: programa, método, condição e ação
3. ConditionalRuleService avalia contexto de usuário/sessão
4. Aplica bloqueio, visibilidade ou valor sugerido em campo
5. Avaliação ocorre na camada de controller, nunca no core
```

### 3.4 Gestão de Usuários

```
1. Admin acessa SystemUserList
2. Cria/edita usuário em SystemUserForm
3. Vincula grupos, papéis e unidades
4. Define status, senha, perfil e atributos operacionais
5. Permissões passam a valer no próximo carregamento/sessão
```

### 3.5 Matriz de Permissões

```
1. Admin cadastra programas em SystemProgramForm
2. Vincula programas a grupos em SystemPermissionController
3. Associa papéis e métodos quando necessário
4. A navegação/menu exibe apenas programas autorizados
5. Acesso direto também é validado no backend
```

### 3.6 Unidades e Escopo

```
1. Admin cadastra SystemUnit
2. Usuários são vinculados a uma ou mais unidades
3. Módulos de negócio filtram dados conforme unidades permitidas
```

### 3.7 Perfil e Segurança

```
1. Usuário acessa SystemProfileForm/View
2. Atualiza dados pessoais e senha
3. SystemProfile2FAForm gerencia segundo fator quando habilitado
4. Fluxos de reset/renovação tratam senha expirada ou esquecida
```

### 3.8 Diagnóstico Administrativo

Telas administrativas oferecem:
- informações de PHP e sistema;
- erro de PHP;
- diff de arquivos contra referência;
- explorer de banco/tabelas;
- painel SQL;
- dashboard de administração;
- verificação de módulos.

---

## 4. Estados e Permissões

### 4.1 Estado de Usuário

```
ATIVO → BLOQUEADO/INATIVO
ATIVO → SENHA_EXPIRADA → RENOVADA
```

### 4.2 Grupos Canônicos

| Grupo | Uso |
|-------|-----|
| Admin | Administração total |
| Atendente | Operação de leads/agenda |
| Representante | Visitas/pedidos em campo |
| Supervisor/Gestor | Gestão e relatórios |

---

## 5. Decisões Arquiteturais

### ADR-001: Permission como Base Separada

**Decisão**: Usuários, grupos, papéis e unidades vivem no banco `permission`.

**Consequências**:
- Autorização desacoplada dos dados de negócio
- Serviços precisam alternar transações entre `permission` e `smid`

### ADR-002: Programas como Recursos Autorizáveis

**Decisão**: Cada tela/classe é um `SystemProgram`.

**Consequências**:
- Menu e acesso direto podem usar a mesma matriz
- Deploy de nova tela exige cadastro/permissão

### ADR-003: Regras Granulares Fora do Core

**Decisão**: Regras condicionais ficam em models/controllers próprios, sem alterar core Adianti.

**Consequências**:
- Upgrade mais seguro
- UI dinâmica sem hardcode no framework

### ADR-004: Diagnósticos Restritos ao Admin

**Decisão**: Ferramentas como SQLPanel, DataBrowser e FilesDiff são administrativas.

**Consequências**:
- Alto poder operacional
- Acesso deve ser fortemente restrito

---

## 6. Referência à Implementação Atual

### 6.1 Controllers

| Área | Arquivos |
|------|----------|
| Login/Senha | `LoginForm`, `SystemPassword*`, `SystemRequestPasswordResetForm` |
| Usuários/Permissões | `SystemUser*`, `SystemGroup*`, `SystemRole*`, `SystemProgram*`, `SystemPermissionController` |
| Unidades | `SystemUnitForm/List` |
| Perfil | `SystemProfile*` |
| Diagnóstico | `SystemAdministrationDashboard`, `SystemDataBrowser`, `SystemDatabaseExplorer`, `SystemFilesDiff`, `SystemPHP*`, `SystemSQLPanel` |
| Integrações | `Evolution*`, `ThreeCPlus*`, `GotoConnect*`, `WhatsNotificacoesForm` |

### 6.2 Models

`app/model/admin/SystemUser.php`, `SystemGroup.php`, `SystemRole.php`, `SystemProgram.php`, `SystemUnit.php`, `SystemPreference.php`, `SystemPermission.php`, tabelas pivô `SystemUserGroup`, `SystemUserRole`, `SystemUserUnit`, `SystemGroupProgram`, `SystemGroupRole`.

---

## 7. Segurança

- Nunca confiar apenas no menu; validar acesso no backend
- Ferramentas de banco/diff restritas a Admin
- Reset de senha com token e expiração
- 2FA quando habilitado
- Auditoria de alterações sensíveis

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
