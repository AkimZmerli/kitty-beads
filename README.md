# Kitty-Beads

**Lightning-fast dev planning dashboard with Git-native issue tracking and terminal integration**

Kitty-Beads combines the **sub-100ms query performance** of [Beads](https://github.com/steveyegge/beads) with a beautiful React dashboard for team planning. Everything lives in Git—no external databases, no merge conflicts.

Perfect for:
- **AI-assisted teams** - CLI + REST API + visual dashboard in one binary
- **Distributed teams** - Issues stored in Git, sync like code
- **High-velocity development** - Fast queries, instant feedback, zero context switching

## Core Features

### 🚀 Performance
- **Sub-100ms queries** - SQLite + intelligent caching
- **Single binary deployment** - Everything embedded (Go server + React frontend)
- **No external services** - Entirely self-contained

### 📋 Development Planning
- **Visual Kanban** - Drag-drop lanes (Planned / In Progress / Review / Done)
- **Feature tracking** - Epics with child tasks, blockers, and dependencies
- **Rich artifacts** - Design docs, specifications, research notes, checklists
- **Priority levels** - 5-tier priority system with visual indicators

### 🔗 Git-Native Design
- **Issues as code** - Stored in `.beads/issues.jsonl`, versioned with Git
- **No merge conflicts** - Hash-based IDs (`bd-a1b2` format) + intelligent conflict resolution
- **Offline-first** - Works without network, syncs when you push

### 💻 Terminal Integration
- **In-browser shell** - Full PTY bash/zsh via WebSocket
- **Run Beads CLI** - `bd create`, `bd list`, `bd sync` from dashboard
- **Real-time sync** - Dashboard updates reflect terminal changes instantly

### 🎨 Modern UI
- **React 19 + Vite** - Fast builds, hot reload in dev
- **Tailwind CSS** - Responsive design that works on any screen
- **Markdown rendering** - Rich text with syntax highlighting
- **Real-time updates** - Auto-refreshing data with 1-second cache

## Quick Start

### Install (One-liner)
```bash
curl -sSL https://raw.githubusercontent.com/AkimZmerli/kitty-beads/main/install.sh | bash
```

Then:
```bash
cd your-project && bd init  # if first time
kitty-beads                 # opens http://localhost:8080
```

### Or Manual Build
```bash
git clone https://github.com/AkimZmerli/kitty-beads
cd kitty-beads
make build && ./bin/kitty-beads
```

## Requirements

- **Go 1.24+** (if building locally)
- **Existing `.beads` directory** - Create one with: `bd init`
- **Node.js 18+** (only needed for frontend development)

## Architecture

```
kitty-beads (single binary)
│
├── HTTP Server (Go, port 8080)
│   ├── /api/features        - List epics with stats
│   ├── /api/kanban/{id}     - Kanban lanes for feature
│   ├── /api/issues          - CRUD for issues
│   ├── /api/ready           - Ready-to-work issues
│   ├── /api/artifact/*      - Design docs, specs, etc.
│   ├── /api/diagnostics     - System stats
│   ├── /api/terminal        - WebSocket PTY shell
│   └── /                    - React SPA dashboard
│
├── Storage Backend
│   ├── SQLite (.beads/beads.db)     - Fast queries, WAL mode
│   ├── JSONL (.beads/issues.jsonl)  - Git version control
│   └── Git hooks                     - Auto-sync on commit
│
└── Frontend (React 19 + Tailwind)
    ├── Overview         - Feature dashboard & stats
    ├── Kanban           - Visual board with lanes
    ├── Tasks            - List view of work
    ├── Plan/Spec/Research - Artifact editors
    ├── Terminal         - Full-screen bash shell
    └── Diagnostics      - Health & statistics
```

### Data Flow
1. **Dashboard loads** → Fetches issues from `/api/issues`
2. **Update issue** → PUT `/api/issues/{id}` → Syncs to JSONL
3. **Terminal command** → WebSocket to PTY → Executes `bd sync`
4. **Git push** → Hook triggers → SQLite cache refreshed
5. **Refresh dashboard** → 1-second cache miss → Latest data

## Usage

### Command-Line Options

```bash
kitty-beads [OPTIONS]

Options:
  -port PORT         HTTP server port (default: 8080)
  -beads-dir PATH    Path to .beads directory (default: .beads)
  -open              Automatically open dashboard in browser
```

### Examples

```bash
# Start on default port (8080)
kitty-beads

# Use custom port
kitty-beads -port 3000

# Point to different .beads directory
kitty-beads -beads-dir /path/to/.beads

# Auto-open browser on start
kitty-beads -open

# Combine options
kitty-beads -port 3000 -open
```

### Accessing from Remote

```bash
# Listen on all interfaces (careful with security)
kitty-beads -port 3000
# Then access from other machines: http://your-ip:3000
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
