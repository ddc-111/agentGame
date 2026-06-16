# AGENTS.md

Monorepo with three independent apps. No shared build system.

## Structure

```
server/     Go 1.21+ backend (Gin + GORM + WebSocket)
client/     Phaser 3 game client (Vite)
gm/         Vue 3 GM editor (Vite + Element Plus + Pinia)
```

## Quick start

```bash
# From repo root, copy config on first run:
cp server/config.example.yaml server/config.yaml

# Server (auto-migrates DB and seeds data on first run)
cd server && go run cmd/gameserver/main.go
# -> http://localhost:8080

# Client
cd client && npm install && npm run dev
# -> http://localhost:5173

# GM editor
cd gm && npm install && npm run dev
# -> http://localhost:5174
```

Windows shortcut: `start.bat` launches all three; `stop.bat` kills them.

## Server details

- Entry: `server/cmd/gameserver/main.go`
- Config: `server/config.yaml` (gitignored; copy from `config.example.yaml`)
- Database: SQLite by default (`game.db`, pure Go driver, no CGO). MySQL also supported.
- Auto-migration + seed on startup if DB is empty (`server/internal/database/seed.go`)
- Seed data: 4 scenes, 3 NPCs, 3 agents, 2 shops, 8 items, 2 tasks, 1 flow, 2 prompt templates
- Models use soft delete (GORM `DeletedAt`)
- JSON fields stored as strings in DB (Behaviors, Schedule, Effect, etc.)

### Key server packages

| Package | Purpose |
|---------|---------|
| `internal/config` | YAML config loading, defaults |
| `internal/database` | GORM init, migrations, seeding |
| `internal/database/models` | All GORM models |
| `internal/database/repository` | Data access layer |
| `internal/network` | Gin router, HTTP handlers, WebSocket |
| `internal/agent` | AI agent stub (placeholder Chat method) |
| `internal/generator` | OpenAI-based config generator |
| `internal/mcp` | MCP (Model Context Protocol) server |

### MCP endpoint

`POST http://localhost:8080/mcp` — JSON-RPC 2.0

Also available as REST wrappers:
- `GET /api/mcp/tools` — list tools
- `POST /api/mcp/call` — call a tool

Full tool docs: see `AGENT.md` in repo root.

### API routes

All under `/api/`: scenes, npcs, agents, llm/providers, prompts, shops, items, tasks, flows, players, conversations, export, import, generator/*, mcp/*, ws.

## Frontend details

- **client/**: Vite + Phaser 3. Entry: `client/src/main.js`. Scenes in `client/src/game/scenes/`.
- **gm/**: Vite + Vue 3 + Vue Router + Pinia + Element Plus. Stores in `gm/src/stores/`, views in `gm/src/views/`.

## Build

```bash
cd server && go build -o bin/gameserver.exe cmd/gameserver/main.go
cd client && npm run build
cd gm && npm run build
```

## Config notes

- `server/config.yaml` has two AI sections: `ai` (for in-game NPC agents) and `generator` (for the config generation AI, can point to a different model/endpoint)
- Both default to OpenAI gpt-4/gpt-4-turbo
- Generator is enabled by default; set `generator.enabled: false` to disable

## Conventions

- Entity codes use snake_case with prefix: `scene_*`, `npc_*`, `agent_*`, `shop_*`, `item_*`, `task_*`, `flow_*`, `template_*`
- All models have `BaseModel` (ID, CreatedAt, UpdatedAt, DeletedAt)
- Foreign keys: NPC→Agent, NPC→Shop, ShopItem→Shop+Item, SceneNPC→Scene+NPC
- No tests, linter, or formatter configured yet
- Database file (`game.db`) and `server/config.yaml` are gitignored

## Gotchas

- SQLite driver is pure Go (`glebarez/sqlite`), not the CGO one — no gcc needed
- Seed only runs if scenes table is empty; deleting all scenes triggers re-seed on next start
- Agent.Chat() is a stub returning placeholder text — real LLM integration is in `generator`
- MCP tool IDs are numeric (float64 from JSON), cast to uint for DB queries
