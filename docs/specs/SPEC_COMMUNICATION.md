# SPEC_COMMUNICATION.md

## 1. Visão Geral

### 1.1 Propósito

O domínio **Communication** provê recursos colaborativos e de comunicação interna: mensagens, notificações, feed/postagens, comentários, wiki, agenda e documentos/drive. Ele funciona como camada transversal para interação entre usuários e grupos.

### 1.2 Subdomínios

| Subdomínio | Descrição |
|------------|-----------|
| **Messages** | Mensagens internas, tags e dropdown |
| **Notifications** | Notificações do sistema |
| **Posts/Feed** | Publicações, comentários, likes e compartilhamento |
| **Wiki/Pages** | Páginas internas pesquisáveis e compartilháveis |
| **Documents/Drive** | Pastas, documentos, bookmarks e compartilhamento |
| **Schedule** | Agenda/compromissos internos |

### 1.3 Schema Canônico (banco `communication`)

| Tabela | Entidade | Campos principais |
|--------|----------|-------------------|
| `adianti_messages` | `SystemMessage` | `sender_login`, `subject`, `message`, `dt_message`, `status` |
| `adianti_message_recipients` | (pivot) | `message_id`, `recipient_login` |
| `adianti_message_tags` | `SystemMessageTag` | `tag`, `message_id` |
| `adianti_notifications` | `SystemNotification` | `login`, `subject`, `message`, `action_url`, `dt_notification`, `dt_read` |
| `adianti_posts` | `SystemPost` | `login`, `title`, `content`, `dt_post`, `pinned` |
| `adianti_post_comments` | `SystemPostComment` | `post_id`, `login`, `content`, `dt_comment` |
| `adianti_post_likes` | `SystemPostLike` | `post_id`, `login` |
| `adianti_post_share_groups` | `SystemPostShareGroup` | `post_id`, `system_group_id` |
| `adianti_wiki_pages` | `SystemWikiPage` | `login`, `title`, `content`, `tags`, `dt_page` |
| `adianti_wiki_share_groups` | `SystemWikiShareGroup` | `page_id`, `system_group_id` |
| `adianti_documents` | `SystemDocument` | `folder_id`, `login`, `filename`, `description`, `dt_document` |
| `adianti_folders` | `SystemFolder` | `login`, `name`, `parent_id` |
| `adianti_document_users` | `SystemDocumentUser` | `document_id`, `login` |
| `adianti_document_groups` | `SystemDocumentGroup` | `document_id`, `system_group_id` |
| `adianti_schedules` | `SystemSchedule` | `title`, `class_name`, `method`, `schedule_type`, `hour`, `minute`, `monthday`, `weekday`, `active` |

---

## 2. Entidades Principais

| Entidade | Descrição |
|----------|-----------|
| `SystemNotification` | Notificação para usuário |
| `SystemMessage` | Mensagem interna |
| `SystemMessageTag` | Tag de mensagem |
| `SystemPost` | Postagem/feed |
| `SystemPostComment` | Comentário em postagem |
| `SystemPostLike` | Curtida |
| `SystemPostShareGroup` | Compartilhamento com grupos |
| `SystemWikiPage` | Página wiki |
| `SystemWikiShareGroup` | Compartilhamento wiki |
| `SystemDocument` | Documento/arquivo |
| `SystemFolder` | Pasta |
| `SystemDocumentUser/Group` | Permissões por usuário/grupo |
| `SystemSchedule` | Agenda |

---

## 3. Fluxos Principais

### 3.1 Mensagens Internas

```
1. Usuário acessa SystemMessageForm/List
2. Cria mensagem para destinatários/grupos
3. Pode classificar por tags
4. Dropdown exibe mensagens recentes/não lidas
5. Usuário abre SystemMessageFormView para leitura
```

### 3.2 Notificações

```
1. Módulo de negócio cria SystemNotification
2. Notificação é associada a usuário/login
3. Interface global exibe pendências
4. Ao abrir/confirmar, notificação muda estado de leitura
```

### 3.3 Feed e Posts

```
1. Usuário cria SystemPost
2. Compartilha com grupos permitidos
3. Usuários comentam via SystemPostCommentForm/List
4. Likes são registrados em SystemPostLike
5. FeedView monta linha do tempo filtrada por permissão
```

### 3.4 Wiki

```
1. Usuário autorizado cria página wiki
2. Define conteúdo, tags e grupos compartilhados
3. WikiSearchList permite busca textual
4. WikiView renderiza página
5. WikiPagePicker permite seleção contextual
```

### 3.5 Drive/Documentos

```
1. Usuário cria pasta em SystemFolderForm
2. Faz upload em SystemDriveDocumentUploadForm
3. Documento é vinculado à pasta
4. Compartilhamento por usuário/grupo controla acesso
5. Bookmark permite favoritos
6. SystemTextDocumentEditor edita documentos de texto
```

### 3.6 Agenda

```
1. Usuário cria compromisso em SystemScheduleForm
2. SystemScheduleList exibe agenda/listagem
3. Eventos podem alimentar notificações internas
```

### 3.7 Jobs Agendados (SystemSchedule)

```
SystemSchedule armazena jobs recorrentes com tipo M/W/D/F:

1. CLI/cron externo chama SystemScheduleService::run()
2. Service filtra schedules ativos conforme data/hora atual
3. Para cada schedule:
   a. Executa classe/método configurado
   b. Registra SystemScheduleLog com status Y (ok) ou N (erro)
4. Configuração de novos jobs via SyncProgramForm ou SQL direto

Tipos:
- M (monthly): dia do mês + hora + minuto
- W (weekly): dia da semana + hora + minuto
- D (daily): hora + minuto
- F (fixed): sempre executado (sem filtro de data)
```

**Invariante**: SystemScheduleService deve ser acionado por um único scheduler externo (ex: cron do OS) para evitar execuções concorrentes.

---

## 4. Permissões e Compartilhamento

| Recurso | Estratégia |
|---------|------------|
| Mensagens | destinatário/login/grupo |
| Posts | grupos compartilhados |
| Wiki | grupos compartilhados |
| Documentos | usuário e grupo |
| Pastas | usuário e grupo |
| Agenda | proprietário/participantes |

---

## 5. Decisões Arquiteturais

### ADR-001: Communication como Base Separada

**Decisão**: Artefatos de comunicação usam estrutura `app/model/communication` e banco/alias `communication` quando configurado.

**Consequências**:
- Separação dos dados de negócio
- Permite evolução independente

### ADR-002: Compartilhamento por Usuário e Grupo

**Decisão**: Documentos/pastas/wiki/posts usam tabelas de vínculo com usuários/grupos.

**Consequências**:
- Controle granular de visibilidade
- Consultas precisam aplicar permissão sempre

### ADR-003: Feed e Wiki como Colaboração Interna

**Decisão**: Manter comunicação interna no SMID, sem depender de ferramenta externa.

**Consequências**:
- Conhecimento operacional centralizado
- Precisa de governança de conteúdo

### ADR-004: Upload/Drive com Permissões Próprias

**Decisão**: Drive tem pastas, documentos, bookmarks e ACL própria.

**Consequências**:
- Reuso para anexos internos
- Atenção a armazenamento e segurança de arquivos

---

## 6. Referência à Implementação Atual

### 6.1 Controllers

| Área | Localização |
|------|-------------|
| Documentos | `app/control/communication/documents/` |
| Mensagens | `app/control/communication/messages/` |
| Wiki | `app/control/communication/pages/` |
| Posts | `app/control/communication/posts/` |
| Agenda | `app/control/communication/schedule/` |

### 6.2 Models

| Área | Localização |
|------|-------------|
| Notificações | `app/model/communication/SystemNotification.php` |
| Documentos | `app/model/communication/documents/` |
| Mensagens | `app/model/communication/messages/` |
| Wiki | `app/model/communication/pages/` |
| Posts | `app/model/communication/posts/` |
| Agenda | `app/model/communication/schedule/` |

---

## 7. Segurança

- Toda leitura deve validar compartilhamento/permissão
- Upload deve validar extensão, tamanho e caminho
- Documentos não devem expor path físico sem controle
- Conteúdo HTML/wiki deve tratar XSS
- Notificações não devem vazar dados entre unidades/grupos

---

**Versão**: 1.0
**Data**: 2026-05-12
**Status**: Draft para revisão
