# Kitty-Beads

**Beads performance + Spec Kitty aesthetics**

A hybrid issue tracker combining [Beads](https://github.com/steveyegge/beads)' lightning-fast Go backend with [Spec Kitty](https://github.com/priivacy-ai/spec-kitty)'s beautiful dashboard UI.

## Features

- **Fast**: Sub-100ms queries via SQLite cache
- **Git-native**: Issues stored as JSONL, version-controlled like code
- **Hash-based IDs**: No merge conflicts (`bd-a1b2` format)
- **Visual Kanban**: Drag-drop lanes (Planned / Doing / Review / Done)
- **Markdown rendering**: Rich issue descriptions with syntax highlighting
- **Real-time updates**: Dashboard auto-refreshes every 5 seconds
- **Dependency graph**: Track blockers and relationships

## Quick Start

```bash
# Build the server
make build

# Start the dashboard
make run

# Open http://localhost:8080
```

## Requirements

- Go 1.24+
- Existing `.beads` directory (run `bd init` first)

## Architecture

```
kitty-beads/
├── src/
│   ├── cmd/
│   │   ├── bd/           # Original Beads CLI
│   │   └── server/       # REST API + Dashboard server
│   │       ├── main.go
│   │       ├── templates/
│   │       │   └── index.html
│   │       └── static/
│   │           └── dashboard/
│   │               ├── dashboard.css
│   │               └── dashboard.js
│   └── internal/         # Beads core (storage, types, etc.)
└── Makefile
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Dashboard UI |
| `/api/features` | GET | List all epics |
| `/api/kanban/{id}` | GET | Kanban lanes for epic |
| `/api/issues` | GET/POST | List/create issues |
| `/api/issues/{id}` | GET/PUT/DELETE | CRUD single issue |
| `/api/ready` | GET | Issues without blockers |
| `/api/health` | GET | Server health check |
| `/api/constitution` | GET | Project README/CLAUDE.md |
| `/api/diagnostics` | GET | Issue statistics |

## Data Model Mapping

| Beads | Dashboard |
|-------|-----------|
| Issue (type=epic) | Feature/Epic |
| Issue (child) | Work Package / Task |
| `status: open` | Planned lane |
| `status: in_progress` | Doing lane |
| `status: closed` | Done lane |
| Dependencies | Blockers |

## License

MIT (inherited from Beads and Spec Kitty)
