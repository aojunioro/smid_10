# MCP — SMID 10

Configuracao em `.cursor/mcp.json` (versionada). Servidores opcionais em `mcp.json.example`.

## Setup (uma vez)

```bash
cd .cursor/mcp/smid-tools
npm install
```

Reinicie o Cursor apos editar `mcp.json`.

## Servidores incluidos

| Servidor | Tipo | Funcao |
|----------|------|--------|
| `smid-tools` | local (Node) | Snapshot do projeto, SPECs, dominios Go, rotas, handoff |
| `smid-docs` | npx | Leitura/escrita em `docs/` via filesystem MCP |

## Ferramentas `smid-tools`

- `smid_project_snapshot` — estado, URLs, lacunas conhecidas
- `smid_list_specs` — lista `docs/specs/SPEC_*.md`
- `smid_list_domain_packages` — pacotes em `backend/internal/domain/`
- `smid_api_routes_summary` — resumo de `routes.go`
- `smid_read_handoff_current` — handoff operacional (`CURRENT.md`)
- `smid_read_handoff_archive_section` — secao do archive (historico; evitar se possivel)
- `smid_read_spec_excerpt` — primeiras linhas de um SPEC

## Opcionais (copiar de `mcp.json.example`)

- **github** — PRs/issues; exige `GITHUB_TOKEN` no ambiente
- **mysql-smid** — consultas read-only ao banco local; **nunca** commitar senhas

## Seguranca

- Nao colocar credenciais de producao no `mcp.json`.
- MySQL MCP apenas em ambiente de dev com dados nao sensiveis ou read-only.

## Troubleshooting

1. Cursor: Settings -> MCP -> verificar `smid-tools` verde
2. Output -> MCP Logs
3. Testar manualmente: `SMID_ROOT=/caminho/smid_10 node .cursor/mcp/smid-tools/index.mjs`
