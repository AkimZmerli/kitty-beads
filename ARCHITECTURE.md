# Kitty-Beads Architecture

> Local-first issue tracker with integrated terminal and CLI tooling.

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         kitty-beads/                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────┐          ┌──────────────────────────────┐│
│  │     Frontend         │   HTTP   │         Backend (Go)         ││
│  │     (React/TS)       │─────────▶│                              ││
│  │                      │          │  ┌────────────────────────┐  ││
│  │  localhost:5173      │    WS    │  │     features/          │  ││
│  │                      │─────────▶│  │  ├── issues/           │  ││
│  │                      │          │  │  ├── kanban/           │  ││
│  └──────────────────────┘          │  │  ├── epics/            │  ││
│                                    │  │  ├── comments/         │  ││
│  ┌──────────────────────┐          │  │  ├── dependencies/     │  ││
│  │     CLI (bd)         │   RPC    │  │  ├── labels/           │  ││
│  │                      │─────────▶│  │  └── statistics/       │  ││
│  │  bd create, bd list  │  (Unix   │  └────────────────────────┘  ││
│  │  bd show, bd close   │  Socket) │                              ││
│  └──────────────────────┘          │  ┌────────────────────────┐  ││
│                                    │  │     internal/          │  ││
│                                    │  │  ├── rpc/       (API)  │  ││
│                                    │  │  ├── storage/   (DB)   │  ││
│                                    │  │  ├── beads/     (Core) │  ││
│                                    │  │  ├── daemon/           │  ││
│                                    │  │  └── ...30 packages    │  ││
│                                    │  └────────────────────────┘  ││
│                                    │                              ││
│                                    │  localhost:8080              ││
│                                    └──────────────────────────────┘│
│                                              │                     │
│                                              ▼                     │
│                                    ┌──────────────────┐            │
│                                    │   .beads/        │            │
│                                    │   ├── beads.db   │ (SQLite)   │
│                                    │   └── config     │            │
│                                    └──────────────────┘            │
└─────────────────────────────────────────────────────────────────────┘
```

## Frontend Architecture

```
frontend/src/
├── main.tsx                 # App entry point
├── App.tsx                  # Router + providers setup
├── index.css                # Global styles (Tailwind)
│
├── layouts/                 # Page scaffolding
│   ├── Layout.tsx           # Main layout (header + sidebar + content)
│   ├── Header.tsx           # Top nav with feature selector
│   ├── Sidebar.tsx          # Left nav with route links
│   └── index.ts
│
├── pages/                   # Route-level components
│   ├── Kanban.tsx           # Kanban board view
│   ├── Roadmap.tsx          # Tree/list view of issues
│   ├── Diagnostics.tsx      # System health page
│   └── IdeationPad.tsx      # Markdown editor for issue plans
│
├── features/                # Feature modules (vertical slices)
│   ├── terminal/            # Integrated terminal
│   │   ├── components/
│   │   │   ├── TerminalPanel.tsx
│   │   │   ├── TerminalTabs.tsx
│   │   │   ├── TerminalInstance.tsx
│   │   │   └── ResizeHandle.tsx
│   │   ├── context.tsx      # Terminal state provider
│   │   ├── types.ts         # TerminalTab, TerminalState
│   │   ├── constants.ts     # Height limits, storage keys
│   │   ├── utils.ts         # ID generation, localStorage
│   │   └── index.ts         # Barrel export
│   │
│   ├── kanban/              # Kanban board feature
│   │   ├── components/
│   │   │   ├── KanbanBoard.tsx
│   │   │   ├── KanbanLane.tsx
│   │   │   ├── KanbanCard.tsx
│   │   │   └── IssueModal.tsx
│   │   ├── hooks/
│   │   ├── types.ts
│   │   ├── constants.ts
│   │   └── index.ts
│   │
│   └── roadmap/             # Roadmap/tree feature
│       ├── components/
│       │   └── BeadCard.tsx
│       └── index.ts
│
├── components/
│   └── ui/                  # Shared UI primitives
│       ├── NeonSelect.tsx   # Custom dropdown
│       ├── MarkdownViewer.tsx
│       └── index.ts
│
├── hooks/                   # App-wide hooks
│   ├── useFeatures.ts       # Fetch features list
│   ├── useRoadmap.ts        # Fetch roadmap data
│   ├── useIdeation.ts       # Issue editing
│   └── useTerminal.ts       # Re-export from feature
│
├── lib/                     # Utilities
│   ├── api.ts               # HTTP client functions
│   └── planParser.ts        # Markdown parsing helpers
│
└── types/
    └── api.ts               # API response types
```

### Key Patterns

| Pattern | Usage |
|---------|-------|
| Feature folders | Each feature is self-contained with components, hooks, types |
| Barrel exports | `index.ts` files for clean imports |
| React Query | Server state management, caching, invalidation |
| Context | Terminal state (global), feature-specific state |
| Tailwind | Utility-first CSS with Tokyo Night theme |

### Data Flow

```
User Action
    │
    ▼
Page Component (pages/)
    │
    ├── useQuery() ──────▶ lib/api.ts ──────▶ Backend HTTP
    │                           │
    │                           ▼
    │                      React Query Cache
    │
    ▼
Feature Components (features/)
    │
    └── useTerminal() ───▶ WebSocket ──────▶ Backend /api/terminal/:id
```

## Backend Architecture

```
backend/
├── beads.go                 # Package root
├── server                   # Compiled binary
│
├── cmd/
│   └── bd/                  # CLI commands (322 files)
│       ├── create.go        # bd create
│       ├── list.go          # bd list
│       ├── close.go         # bd close
│       ├── sync_*.go        # Git sync commands
│       ├── migrate_*.go     # Migration commands
│       └── ...
│
├── features/                # Vertical slices (HTTP handlers)
│   ├── issues/
│   │   ├── handler.go       # HTTP handlers
│   │   ├── service.go       # Business logic
│   │   ├── repository.go    # Data access
│   │   └── types.go         # DTOs
│   ├── kanban/
│   ├── epics/
│   ├── comments/
│   ├── dependencies/
│   ├── labels/
│   ├── gates/
│   ├── statistics/
│   ├── compaction/
│   └── export/
│
├── internal/                # Core packages (31 packages)
│   ├── rpc/                 # RPC server + handlers
│   │   ├── server.go
│   │   ├── server_core.go
│   │   ├── server_issues_epics.go  # Issue/epic handlers
│   │   ├── client.go
│   │   └── protocol.go
│   │
│   ├── storage/             # Persistence layer
│   │   ├── storage.go       # Interface
│   │   └── sqlite/          # SQLite implementation
│   │
│   ├── beads/               # Core domain logic
│   ├── daemon/              # Background processes
│   ├── config/              # Configuration
│   ├── formula/             # Issue formulas/templates
│   ├── git/                 # Git integration
│   ├── hooks/               # Git hooks
│   ├── syncbranch/          # Branch synchronization
│   ├── compact/             # Database compaction
│   ├── export/              # Export functionality
│   ├── importer/            # Import functionality
│   ├── validation/          # Input validation
│   ├── types/               # Shared types
│   └── ...
│
└── shared/                  # Cross-cutting concerns
```

### Vertical Slice Pattern

Each feature follows:

```
features/{feature}/
├── handler.go      # HTTP handlers (thin, delegates to service)
├── service.go      # Business logic (orchestrates repos)
├── repository.go   # Data access (talks to storage)
└── types.go        # Request/response DTOs
```

Example flow:
```
HTTP Request
    │
    ▼
handler.go ─────▶ Parse request, validate
    │
    ▼
service.go ─────▶ Business logic, orchestration
    │
    ▼
repository.go ──▶ SQL queries via storage interface
    │
    ▼
internal/storage/ ─▶ SQLite execution
```

### API Endpoints

| Endpoint | Method | Handler |
|----------|--------|---------|
| `/api/issues` | GET | issues.HandleList |
| `/api/issues` | POST | issues.HandleCreate |
| `/api/issues/{id}` | GET | issues.HandleGet |
| `/api/issues/{id}` | PUT | issues.HandleUpdate |
| `/api/issues/{id}` | DELETE | issues.HandleDelete |
| `/api/issues/{id}/close` | POST | issues.HandleClose |
| `/api/ready` | GET | issues.HandleReady |
| `/api/blocked` | GET | issues.HandleBlocked |
| `/api/kanban` | GET | kanban.HandleBoard |
| `/api/terminal/{id}` | WS | Terminal WebSocket |

### CLI Commands

The `bd` CLI communicates with the daemon via Unix socket RPC:

```bash
bd create "Task title"       # Create issue
bd list                      # List issues
bd show <id>                 # Show issue details
bd close <id>                # Close issue
bd sync                      # Sync with git
bd daemon start              # Start background daemon
```

## Data Model

```
┌─────────────────┐     ┌─────────────────┐
│     Issue       │     │     Label       │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │◀───▶│ issue_id (FK)   │
│ title           │     │ name            │
│ description     │     └─────────────────┘
│ design          │
│ status          │     ┌─────────────────┐
│ priority        │     │   Dependency    │
│ issue_type      │     ├─────────────────┤
│ parent_id (FK)  │◀───▶│ issue_id (FK)   │
│ assignee        │     │ depends_on (FK) │
│ created_at      │     │ dep_type        │
│ updated_at      │     └─────────────────┘
│ closed_at       │
│ due_at          │     ┌─────────────────┐
│ defer_until     │     │    Comment      │
└─────────────────┘     ├─────────────────┤
        │               │ issue_id (FK)   │
        │               │ content         │
        ▼               │ author          │
┌─────────────────┐     │ created_at      │
│  Child Issues   │     └─────────────────┘
└─────────────────┘
```

## Development

### Running Locally

```bash
# Terminal 1: Backend
cd backend
go build -o server ./cmd/server
./server

# Terminal 2: Frontend
cd frontend
pnpm install
pnpm dev

# Terminal 3: CLI (optional)
cd backend
go build -o bd ./cmd/bd
./bd list
```

### Environment

| Component | Port | Protocol |
|-----------|------|----------|
| Frontend (Vite) | 5173 | HTTP |
| Backend API | 8080 | HTTP |
| Terminal WS | 8080 | WebSocket |
| CLI ↔ Daemon | Unix socket | RPC |

### Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | React 18, TypeScript, Vite, TailwindCSS |
| State | TanStack Query (React Query v5), React Context |
| Terminal | xterm.js, WebSocket |
| Backend | Go 1.21+, Chi router |
| Database | SQLite |
| CLI | Cobra |
| Package Manager | pnpm |
