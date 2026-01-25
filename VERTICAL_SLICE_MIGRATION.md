# Vertical Slice Architecture Migration Plan

## Goals
- Developer Experience: Easier navigation and understanding
- Maintainability: Reduced coupling, clearer feature boundaries
- Performance/Deployment: Enable independent feature work

## Current Pain Points

**Project Structure:**
- 158MB bloat from `beads-upstream/` and `spec-kitty-upstream/` (to be deleted)
- 360 subdirectories in `src/cmd/bd/` - extreme width for CLI
- 11 levels of nesting in Go packages

**Backend:**
- Monolithic `main.go` (659 lines) with manual path parsing
- Storage interface with 45+ methods (SRP violation)
- RPC layer: 14,866 LOC in horizontal layer
- No middleware (logging, auth, validation)

**Frontend:**
- Horizontal layering: `components/`, `hooks/`, `pages/`, `lib/`
- Only Terminal has vertical structure
- Large page components (Kanban 270 LOC, Roadmap 220 LOC)

---

## Target Structure

### Backend Vertical Slices
```
src/
├── cmd/
│   ├── server/main.go        # Simplified: router + middleware setup
│   └── bd/commands/          # CLI grouped by feature
├── features/                 # Vertical slices
│   ├── issues/
│   │   ├── handler.go        # HTTP handlers
│   │   ├── rpc.go            # RPC handlers
│   │   ├── service.go        # Business logic
│   │   ├── repository.go     # Data access interface
│   │   └── types.go          # DTOs
│   ├── dependencies/
│   ├── kanban/
│   ├── labels/
│   ├── comments/
│   ├── gates/
│   ├── epics/
│   ├── compaction/
│   ├── export/
│   └── statistics/
└── shared/
    ├── storage/              # Core interface (~10 methods)
    │   └── sqlite/           # Implementation
    ├── middleware/           # Logging, recovery, request ID
    ├── rpc/                  # Infrastructure only
    └── types/                # Slimmed core types
```

### Frontend Vertical Slices
```
frontend/src/
├── features/
│   ├── kanban/
│   │   ├── components/       # KanbanBoard, KanbanLane, KanbanCard
│   │   ├── hooks/            # useKanban.ts
│   │   └── index.ts
│   ├── roadmap/
│   ├── ideation/
│   ├── terminal/             # Already partially sliced
│   ├── overview/
│   ├── specify/
│   └── plan/
└── shared/
    ├── components/           # Layout, Header, Sidebar, MarkdownViewer
    ├── lib/                  # api.ts
    └── types/                # Common types only
```

---

## Implementation Phases

### Phase 0: Cleanup (Day 1)
**Delete bloat directories:**
1. `rm -rf beads-upstream/`
2. `rm -rf spec-kitty-upstream/`
3. Update any import references
4. Run `go build ./...` and `go test ./...`

**Verification:** Build passes, tests pass, web UI works

---

### Phase 1: Infrastructure Foundation (Week 1)
**Add middleware without changing handlers:**

1. Create `src/shared/middleware/`:
   - `logger.go` - Request timing + path logging
   - `recovery.go` - Panic recovery with stack trace
   - `request_id.go` - UUID injection

2. Wrap existing mux in `main.go`:
   ```go
   handler := middleware.Chain(
       middleware.Logger(),
       middleware.Recovery(),
   )(mux)
   ```

3. Create empty `src/features/` directory structure

**Critical files:**
- `src/cmd/server/main.go` - Add middleware wrapping

**Verification:** Logs show request timing, panic recovery works

---

### Phase 2: Storage Interface Split (Week 2)
**Create feature-specific repository interfaces:**

1. Define interfaces in `src/features/*/repository.go`:
   ```go
   // features/issues/repository.go
   type Repository interface {
       Create(ctx, issue, actor) error
       Get(ctx, id) (*types.Issue, error)
       Update(ctx, id, updates, actor) error
       Close(ctx, id, reason, actor, session) error
       Delete(ctx, id) error
       Search(ctx, query, filter) ([]*types.Issue, error)
   }
   ```

2. Create adapter in `src/shared/storage/adapter.go`:
   ```go
   func (a *StorageAdapter) Issues() issues.Repository {
       return &issueRepoAdapter{a.sqlite}
   }
   ```

3. Keep existing `Storage` interface for backward compatibility

**Critical files:**
- `src/internal/storage/storage.go` (45+ methods to split)
- New: `src/shared/storage/adapter.go`
- New: `src/features/issues/repository.go`

**Verification:** All existing tests pass, new interface tests pass

---

### Phase 3: First Vertical Slice - Issues (Week 3-4)
**Complete the issues feature slice:**

1. Create `src/features/issues/`:
   - `types.go` - CreateArgs, UpdateArgs, ListFilter DTOs
   - `repository.go` - Interface (from Phase 2)
   - `service.go` - Extract business logic from RPC handlers
   - `rpc.go` - Thin RPC handlers calling service
   - `handler.go` - HTTP handlers

2. Extract from `src/internal/rpc/server_issues_epics.go` (2368 lines):
   - Move `handleCreate`, `handleUpdate`, `handleClose`, `handleDelete`, `handleList`, `handleShow` logic to `service.go`
   - Keep thin dispatcher in RPC

3. Update HTTP server to use new handlers

**Critical files:**
- `src/internal/rpc/server_issues_epics.go` (extract from)
- `src/cmd/server/main.go` (wire new handlers)
- New: `src/features/issues/*.go`

**Verification:** All issue CRUD works via HTTP and RPC

---

### Phase 4: Remaining Backend Slices (Week 5-8)
**Migrate by complexity (simplest first):**

| Order | Feature | LOC | Notes |
|-------|---------|-----|-------|
| 1 | labels | ~200 | Simple CRUD |
| 2 | comments | ~300 | Simple CRUD |
| 3 | statistics | ~200 | Read-only |
| 4 | kanban | ~200 | HTTP-only aggregation |
| 5 | dependencies | ~500 | Graph operations |
| 6 | gates | ~500 | Async coordination |
| 7 | epics | ~300 | Aggregation |
| 8 | molecules | ~1000 | Workflow templates |
| 9 | compaction | ~800 | Background jobs |
| 10 | export | ~1000 | JSONL sync |

**Per-slice process:**
1. Create feature directory with interface
2. Extract service logic from RPC handlers
3. Update RPC to call service
4. Update HTTP to call service
5. Write feature-specific tests

---

### Phase 5: CLI Reorganization (Week 9)
**Group 350 CLI files by feature:**

```
cmd/bd/commands/
├── issues/      # create, update, close, delete, list, show, ready, blocked
├── dependencies/# dep add, dep remove, dep tree
├── labels/      # label add, label remove
├── comments/    # comment add, comment list
├── gates/       # gate create, gate show, gate wait
├── molecules/   # mol create, mol pour, wisp
├── epics/       # epic status
├── compaction/  # compact
├── export/      # flush, import
├── sync/        # sync, daemon
├── config/      # config get, config set
└── admin/       # doctor, audit, cleanup
```

**Verification:** All CLI commands work identically via cobra aliases

---

### Phase 6: Frontend Vertical Slices (Week 10)
**Reorganize by feature:**

1. Create `features/` directory
2. Use Terminal as template (already sliced)
3. Migrate Kanban:
   - `pages/Kanban.tsx` -> `features/kanban/components/KanbanBoard.tsx`
   - Extract: `KanbanLane.tsx`, `KanbanCard.tsx`, `IssueModal.tsx`
   - Move `hooks/useKanban.ts` -> `features/kanban/hooks/`
4. Migrate Roadmap, IdeationPad similarly
5. Keep shared components in `shared/`

**Critical files:**
- `src/cmd/server/frontend/src/pages/Kanban.tsx` (270 lines - extract)
- `src/cmd/server/frontend/src/pages/Roadmap.tsx` (220 lines - extract)

**Verification:** All pages render, no visual regressions

---

### Phase 7: Cleanup (Week 11)
1. Remove empty `internal/rpc/server_*.go` files
2. Remove old horizontal directories from frontend
3. Update import paths
4. Run `staticcheck` for dead code

---

## Verification Strategy

**Per-phase checklist:**
```bash
# Build
go build ./...
go vet ./...

# Tests
go test ./... -race
go test ./... -count=10 -short  # Flaky detection

# Frontend
npm run build
npm run test
```

**Integration verification:**
- Start daemon, run CLI commands
- Start web server, interact with UI
- Run full export/import cycle

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Breaking storage queries | Keep old interface as deprecated alias |
| Breaking RPC protocol | Keep protocol.go unchanged, only restructure handlers |
| Breaking CLI scripts | Use cobra aliases to maintain old command paths |
| Frontend import breaks | Use index.ts barrel exports |

**Rollback:** Tag before each phase (`v0.x.x-pre-phase-N`), each phase is independently reversible.
