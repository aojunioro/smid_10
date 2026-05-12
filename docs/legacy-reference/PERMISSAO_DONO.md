# Permissao do Proprietario

Objetivo: proteger o proprietario fixo do sistema, que e exclusivamente o usuario `system_users.id = 1`, impedindo que outros usuarios/clientes visualizem ou manipulem esse login, acessem programas exclusivos do proprietario ou herdem o grupo reservado do proprietario.

## Convencoes oficiais
- `OWNER_USER_ID = 1`: regra canonica para definir quem e proprietario.
- `OWNER_GROUP_ID = 1`: grupo reservado do proprietario.
- Apenas o usuario `ID = 1` pode ser considerado proprietario.
- O grupo `1` nao pode permanecer vinculado a nenhum usuario diferente do `ID = 1`.

## Modelo atual de protecao

### 1. Proprietario fixo por usuario
- Arquivo: `app/model/admin/SystemPermission.php`
- Metodo: `isOwner()`
- Regra atual: compara o `userid` em sessao com `OWNER_USER_ID`.
- O valor de sessao `is_owner` e recalculado em `ApplicationAuthenticationService::loadSessionVars()`, inclusive em login, recarga de permissao e personificacao.

### 2. Programas exclusivos por metadata
- Arquivo: `app/model/admin/SystemProgram.php`
- Coluna: `system_program.access_scope`
- Valores oficiais:
  - `standard`
  - `owner_only`
- A protecao deixou de depender do hardcode em runtime. Agora o controller consulta o escopo cadastrado do programa.

### 3. Bloqueio central no backend
- Arquivo: `app/model/admin/SystemPermission.php`
- Metodo: `checkPermission()`
- Fluxo:
  - identifica o `access_scope` do controller via `SystemProgram::getAccessScopeByController()`
  - se o escopo for `owner_only` e o usuario atual nao for o proprietario `ID = 1`, o acesso e negado
  - logado: mostra `Acesso restrito` e redireciona para `EmptyPage`
  - nao logado: mostra `Acesso restrito` e redireciona para `LoginForm`

### 4. Menu, busca e sessao tambem respeitam o escopo
- Arquivo: `app/model/admin/SystemUser.php`
- Metodo: `getAllPrograms()`
- Todos os programas carregados para a sessao passam por `SystemProgram::filterProgramsByUserId()`.
- Efeito pratico:
  - `TSession['programs']` nao inclui programas `owner_only` para usuarios comuns
  - `SearchBox` deixa de listar programas exclusivos do proprietario
  - frontpages exclusivas do proprietario nao sao escolhidas para usuarios comuns

## Frontpage (tela inicial)

### Regra central
- Arquivo: `app/model/admin/SystemPermission.php`
- Metodo: `resolveFrontpageForUser()`
- Hierarquia:
  - frontpage individual do usuario
  - frontpage do primeiro grupo aplicavel
  - `EmptyPage`
- A selecao ignora qualquer `SystemProgram` com `access_scope = owner_only` quando o usuario nao for o proprietario.

### Pontos que usam essa regra
- `app/control/admin/LoginForm.php`
- `app/control/admin/SystemUserList.php` em `onImpersonation`
- `app/service/system/SystemPermissionService.php`

## Administracao do access_scope

### 1. Cadastro de programas
- Arquivo: `app/control/admin/SystemProgramForm.php`
- Novo campo: `Escopo de acesso`
- Opcoes atuais:
  - `Publico`
  - `Somente proprietario`
- Regras:
  - somente o proprietario deve administrar esse cadastro
  - ao salvar como `owner_only`, o grupo `1` e garantido automaticamente entre os grupos do programa
  - se o programa continuar vinculado a grupos diferentes do `1`, o sistema exibe um aviso de confirmacao
  - essa confirmacao nao limpa vinculos automaticamente; ela apenas pergunta se o proprietario deseja manter esses vinculos por enquanto

### 2. Listagem de programas
- Arquivo: `app/control/admin/SystemProgramList.php`
- A listagem agora exibe:
  - coluna `Escopo`
  - filtro `Todos / Publico / Somente proprietario`

### 3. Formularios de grupo e usuario
- `app/control/admin/SystemGroupForm.php`
- `app/control/admin/SystemUserForm.php`
- Comportamento para nao proprietario:
  - programas `owner_only` nao aparecem como opcoes de tela inicial
  - programas `owner_only` nao aparecem na checklist de programas do grupo
  - validacao backend recusa tentativa manual de atribuir uma frontpage `owner_only`
  - vinculos ocultos existentes sao preservados no save, evitando remocao silenciosa

### 4. Listagens de grupo e usuario
- `app/control/admin/SystemGroupList.php`
- `app/control/admin/SystemUserList.php`
- Se a frontpage configurada for `owner_only` e o usuario atual nao for o proprietario, a listagem exibe `-` em vez do nome do programa.

## Seed inicial dos programas exclusivos do proprietario

O script de migracao inicial marca como `owner_only` os programas historicamente protegidos:

- `SystemAdministrationDashboard`
- `SystemProgramList`
- `SystemProgramForm`
- `SyncProgramForm`
- `SystemRoleList`
- `SystemRoleForm`
- `SystemPreferenceForm`
- `SystemWikiList`
- `SystemPostList`
- `SystemScheduleList`
- `SystemLogDashboard`
- `SystemAccessLogList`
- `SystemChangeLogView`
- `SystemSqlLogList`
- `SystemRequestLogList`
- `SystemRequestLogView`
- `SystemScheduleLogList`
- `SystemPHPErrorLogView`
- `SystemSessionVarsView`
- `SystemDatabaseExplorer`
- `SystemTableList`
- `SystemDataBrowser`
- `SystemSQLPanel`
- `SystemModulesCheckView`
- `SystemFilesDiff`
- `SystemInformationView`
- `SystemPHPInfoView`

Observacao:
- essa lista e usada apenas para a migracao inicial do dado.
- a protecao ativa em runtime passa a depender do `access_scope` salvo em `system_program`.

## Banco de dados e saneamento

### 1. Exclusividade do usuario proprietario
- Script: `app/database/scripts/2 - permissoes/2.11-enforce_owner_user_id_1.sql`
- Objetivo:
  - garantir o vinculo `system_user_id = 1` com `system_group_id = 1`
  - remover vinculos do grupo `1` para qualquer outro usuario

### 2. Migracao do escopo de acesso
- Script: `app/database/scripts/2 - permissoes/2.12-add_system_program_access_scope.sql`
- Objetivo:
  - adicionar a coluna `access_scope` em `system_program`
  - normalizar valores invalidos para `standard`
  - marcar os programas historicamente protegidos como `owner_only`

## Resumo operacional
- O proprietario do SMID e apenas o usuario `ID = 1`.
- O grupo `1` e reservado e nao pode ser atribuido a outros usuarios.
- Programas exclusivos do proprietario sao controlados por `system_program.access_scope`.
- A protecao nao depende apenas do menu: existe bloqueio central, filtragem da sessao, filtros de listagem, validacao de formulario e saneamento SQL.
- A personificacao recalcula `is_owner`, evitando heranca indevida da sessao do proprietario.
- Ao marcar um programa como `owner_only`, o sistema avisa sobre vinculos de grupos comuns, mas nao apaga esses vinculos automaticamente.
