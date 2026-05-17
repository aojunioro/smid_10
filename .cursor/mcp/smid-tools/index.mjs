#!/usr/bin/env node
/**
 * SMID 10 — MCP tools (stdio).
 * Read-only helpers for SPECs, routes, domain packages, handoff.
 */
import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js'
import { z } from 'zod'
import fs from 'node:fs/promises'
import fsSync from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const defaultRoot = path.resolve(__dirname, '../../..')

function normalizeRootCandidate(raw) {
  if (!raw || typeof raw !== 'string') return null
  let value = raw.trim()
  if (value.startsWith('~/')) {
    value = path.join(process.env.HOME || '', value.slice(2))
  }
  const tildeDup = value.indexOf('/~/')
  if (tildeDup !== -1) {
    value = value.slice(0, tildeDup)
  }
  return path.resolve(value)
}

function resolveProjectRoot() {
  const candidates = [process.env.SMID_ROOT, defaultRoot]
  for (const raw of candidates) {
    const root = normalizeRootCandidate(raw)
    if (root && fsSync.existsSync(path.join(root, 'AGENTS.md'))) {
      return root
    }
  }
  return defaultRoot
}

async function readText(filePath) {
  return fs.readFile(filePath, 'utf8')
}

async function fileExists(filePath) {
  try {
    await fs.access(filePath)
    return true
  } catch {
    return false
  }
}

async function listSpecs(root) {
  const dir = path.join(root, 'docs/specs')
  const entries = await fs.readdir(dir)
  return entries
    .filter((f) => f.startsWith('SPEC_') && f.endsWith('.md'))
    .sort()
}

async function listDomainPackages(root) {
  const dir = path.join(root, 'backend/internal/domain')
  if (!(await fileExists(dir))) return []
  const entries = await fs.readdir(dir, { withFileTypes: true })
  return entries.filter((e) => e.isDirectory()).map((e) => e.name).sort()
}

function summarizeRoutes(routesSource) {
  const groups = []
  const groupRe = /v1\.Group\("([^"]+)"\)/g
  let m
  while ((m = groupRe.exec(routesSource)) !== null) {
    groups.push(m[1])
  }
  const methodRe = /(GET|POST|PUT|PATCH|DELETE)\("([^"]*)"/g
  const routes = []
  while ((m = methodRe.exec(routesSource)) !== null) {
    routes.push(`${m[1]} ${m[2] || '/'}`)
  }
  return { groups: [...new Set(groups)], routeCount: routes.length, sample: routes.slice(0, 40) }
}

function extractHandoffSection(markdown, sectionNum) {
  const re = new RegExp(
    `## ${sectionNum}\\.[^\\n]*\\n([\\s\\S]*?)(?=\\n## \\d+\\.|\\n---\\n|$)`,
    'm'
  )
  const match = markdown.match(re)
  return match ? match[1].trim().slice(0, 12000) : null
}

const server = new McpServer({
  name: 'smid-tools',
  version: '1.0.0',
})

server.registerTool(
  'smid_project_snapshot',
  {
    description:
      'Returns SMID 10 project paths, API URLs, known gaps, and harness locations (AGENTS.md, skills, MCP).',
    inputSchema: {},
  },
  async () => {
    const root = resolveProjectRoot()
    const text = [
      `project_root: ${root}`,
      'api_staging: https://api.s10.smydi.com.br',
      'app_staging: https://s10.smydi.com.br',
      'health: GET /healthz',
      'harness:',
      '  - AGENTS.md',
      '  - docs/handoff/CURRENT.md',
      '  - docs/handoff/README.md (indice da serie)',
      '  - .cursor/skills/',
      '  - .cursor/rules/',
      'known_gaps:',
      '  - JWT middleware not applied on protected routes (routes.go)',
      '  - Frontend login still mock; apiClient unused in features',
      '  - Backend domains missing: metas, relatorios, integracoes_jobs',
      '  - Business rules mostly CRUD-only vs full SPECs',
    ].join('\n')
    return { content: [{ type: 'text', text }] }
  }
)

server.registerTool(
  'smid_list_specs',
  {
    description: 'Lists all SPEC_*.md files in docs/specs/.',
    inputSchema: {},
  },
  async () => {
    const root = resolveProjectRoot()
    const specs = await listSpecs(root)
    return {
      content: [
        {
          type: 'text',
          text: specs.map((s) => `docs/specs/${s}`).join('\n') || 'No SPECs found.',
        },
      ],
    }
  }
)

server.registerTool(
  'smid_list_domain_packages',
  {
    description: 'Lists Go domain packages under backend/internal/domain/.',
    inputSchema: {},
  },
  async () => {
    const root = resolveProjectRoot()
    const domains = await listDomainPackages(root)
    return {
      content: [
        {
          type: 'text',
          text: domains.length
            ? domains.map((d) => `backend/internal/domain/${d}`).join('\n')
            : 'No domain packages found.',
        },
      ],
    }
  }
)

server.registerTool(
  'smid_api_routes_summary',
  {
    description: 'Summarizes Echo route groups and HTTP methods from backend/internal/http/routes.go.',
    inputSchema: {},
  },
  async () => {
    const root = resolveProjectRoot()
    const routesPath = path.join(root, 'backend/internal/http/routes.go')
    if (!(await fileExists(routesPath))) {
      return { content: [{ type: 'text', text: `Missing: ${routesPath}` }], isError: true }
    }
    const source = await readText(routesPath)
    const summary = summarizeRoutes(source)
    const text = [
      `file: ${routesPath}`,
      `route_groups: ${summary.groups.join(', ')}`,
      `http_bindings_found: ${summary.routeCount}`,
      'sample:',
      ...summary.sample.map((r) => `  ${r}`),
    ].join('\n')
    return { content: [{ type: 'text', text }] }
  }
)

server.registerTool(
  'smid_read_handoff_current',
  {
    description:
      'Reads docs/handoff/CURRENT.md (operational handoff). Use instead of the archived monolith.',
    inputSchema: {},
  },
  async () => {
    const root = resolveProjectRoot()
    const currentPath = path.join(root, 'docs/handoff/CURRENT.md')
    if (!(await fileExists(currentPath))) {
      return { content: [{ type: 'text', text: `Missing: ${currentPath}` }], isError: true }
    }
    const md = await readText(currentPath)
    return { content: [{ type: 'text', text: md }] }
  }
)

server.registerTool(
  'smid_read_handoff_archive_section',
  {
    description:
      'Reads a numbered section from docs/handoff/archive/BOOTSTRAP_HISTORY.md (historical only; large file).',
    inputSchema: {
      section: z.number().int().min(1).max(11).describe('Section number from archived bootstrap handoff'),
    },
  },
  async ({ section }) => {
    const root = resolveProjectRoot()
    const archivePath = path.join(root, 'docs/handoff/archive/BOOTSTRAP_HISTORY.md')
    if (!(await fileExists(archivePath))) {
      return { content: [{ type: 'text', text: `Missing: ${archivePath}` }], isError: true }
    }
    const md = await readText(archivePath)
    const body = extractHandoffSection(md, section)
    if (!body) {
      return {
        content: [{ type: 'text', text: `Section ${section} not found in archive.` }],
        isError: true,
      }
    }
    return {
      content: [
        {
          type: 'text',
          text: `# ARCHIVE section ${section} (historical)\n\n${body.slice(0, 15000)}`,
        },
      ],
    }
  }
)

server.registerTool(
  'smid_read_spec_excerpt',
  {
    description: 'Reads the first N lines of a SPEC file by name (e.g. LEADS or SPEC_LEADS).',
    inputSchema: {
      spec: z.string().describe('SPEC name: LEADS, VISITAS, or SPEC_LEADS.md'),
      max_lines: z.number().int().min(20).max(400).optional(),
    },
  },
  async ({ spec, max_lines: maxLines = 120 }) => {
    const root = resolveProjectRoot()
    let fileName = spec.trim()
    if (!fileName.endsWith('.md')) {
      fileName = fileName.startsWith('SPEC_') ? `${fileName}.md` : `SPEC_${fileName.toUpperCase()}.md`
    }
    const specPath = path.join(root, 'docs/specs', fileName)
    if (!(await fileExists(specPath))) {
      return { content: [{ type: 'text', text: `Missing: ${specPath}` }], isError: true }
    }
    const full = await readText(specPath)
    const lines = full.split('\n').slice(0, maxLines)
    return {
      content: [
        {
          type: 'text',
          text: `# ${fileName} (first ${maxLines} lines)\n\n${lines.join('\n')}`,
        },
      ],
    }
  }
)

async function main() {
  const transport = new StdioServerTransport()
  await server.connect(transport)
}

main().catch((err) => {
  console.error('smid-tools MCP failed:', err)
  process.exit(1)
})
