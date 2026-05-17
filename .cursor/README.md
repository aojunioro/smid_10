# Cursor harness — SMID 10

Artefatos para agentes e desenvolvedores no Cursor IDE.

| Pasta / arquivo | Conteudo |
|-----------------|----------|
| `../AGENTS.md` | Guia canonico do projeto |
| `skills/` | Workflows (SPEC, frontend, deploy, handoff, invariantes) |
| `rules/` | Regras `.mdc` por glob (Go, TS, SPECs) |
| `mcp.json` | MCP: `smid-tools` + `smid-docs` |
| `mcp/` | Servidor Node `smid-tools` + README |

## Primeiro uso

```bash
cd .cursor/mcp/smid-tools && npm install
```

Reabrir o workspace no Cursor para carregar skills e MCP.

## Windsurf

Workflows equivalentes em `.windsurf/workflows/`. No Cursor, preferir `.cursor/skills/`.
