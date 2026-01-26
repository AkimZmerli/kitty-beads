# Vertical Slice Architecture Migration Plan

## Goals

- Developer Experience: Easier navigation and understanding
- Maintainability: Reduced coupling, clearer feature boundaries
- Performance/Deployment: Enable independent feature work

---

## Migration Status Summary (Updated 2026-01-26)

### Overall Progress

| Phase | Description | Status | Notes |
|-------|-------------|--------|-------|
| Phase 0 | Cleanup | ✅ Done | `beads-upstream/`, `spec-kitty-upstream/` deleted |
| Phase 1 | Infrastructure Foundation | ✅ Done | `shared/middleware/` created (172 LOC) |
| Phase 2 | Storage Interface Split | 🔶 Partial | `shared/storage/adapter.go` (323 LOC) |
| Phase 3-4 | Backend Vertical Slices | 🔶 Scaffolded | `features/` structure exists (5,489 LOC) |
| Phase 5 | CLI Reorganization | 🔶 80% Done | 8/10 feature packages migrated |
| Phase 6 | Frontend Vertical Slices | ⏳ Not started | |
| Phase 7 | Final Cleanup | ⏳ Not started | |

### Phase 5 CLI Status (Primary Focus)

| Feature | Status | LOC |
|---------|--------|-----|
| labels | ✅ Done | Full migration |
| comments | ✅ Done | Full migration |
| epics | ✅ Done | Full migration |
| gates | ✅ Done | Full migration |
| dependencies | ✅ Done | Full migration |
| config | ✅ Done | Full migration |
| admin | ✅ Done | Command group only |
| molecules | ✅ Done | ~6,000 LOC migrated |
| **sync** | ⏳ Pending | Large, complex deps |
| **issues** | ⏳ Pending | Core commands |

### Key Metrics

- **322 .go files** remain in `cmd/bd/` root (main package)
- **~6,000 LOC** migrated to `commands/molecules/`
- **~5,500 LOC** scaffolded in `features/` (backend vertical slices)
- **172 LOC** in `shared/middleware/` (logger, recovery, request_id)
- **323 LOC** in `shared/storage/adapter.go`

### Current Directory Structure

```
src/
├── cmd/bd/
│   ├── commands/           # Migrated CLI packages
│   │   ├── admin/          ✅
│   │   ├── comments/       ✅
│   │   ├── config/         ✅
│   │   ├── dependencies/   ✅
│   │   ├── epics/          ✅
│   │   ├── gates/          ✅
│   │   ├── labels/         ✅
│   │   ├── molecules/      ✅ (~6,000 LOC)
│   │   └── shared/template/
│   └── *.go                # 322 files still in main package
├── features/               # Backend vertical slices (scaffolded)
│   ├── comments/           handler.go, repository.go, service.go, types.go
│   ├── compaction/
│   ├── dependencies/
│   ├── epics/
│   ├── export/
│   ├── gates/
│   ├── issues/             + rpc.go
│   ├── kanban/
│   ├── labels/
│   └── statistics/
└── shared/
    ├── middleware/         logger.go, recovery.go, request_id.go
    └── storage/            adapter.go
```

### Next Steps

1. **Phase 5 completion**: Migrate `sync` and `issues` CLI commands
2. **Phase 3-4 completion**: Wire up scaffolded `features/` with actual business logic
3. **Phase 6**: Frontend vertical slices
4. **Phase 7**: Final cleanup and dead code removal

---

## Current Pain Points

**Project Structure:**

- ~~158MB bloat from `beads-upstream/` and `spec-kitty-upstream/`~~ ✅ DELETED
- 322 files still in `src/cmd/bd/` root - needs further migration (sync, issues)
- 11 levels of nesting in Go packages

**Backend:**

- Monolithic `main.go` (659 lines) with manual path parsing
- Storage interface with 45+ methods (SRP violation) - adapter.go started
- RPC layer: 14,866 LOC in horizontal layer - features/ scaffolded
- ~~No middleware~~ ✅ DONE (shared/middleware: logger, recovery, request_id)

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

| Order | Feature      | LOC   | Notes                 |
| ----- | ------------ | ----- | --------------------- |
| 1     | labels       | ~200  | Simple CRUD           |
| 2     | comments     | ~300  | Simple CRUD           |
| 3     | statistics   | ~200  | Read-only             |
| 4     | kanban       | ~200  | HTTP-only aggregation |
| 5     | dependencies | ~500  | Graph operations      |
| 6     | gates        | ~500  | Async coordination    |
| 7     | epics        | ~300  | Aggregation           |
| 8     | molecules    | ~1000 | Workflow templates    |
| 9     | compaction   | ~800  | Background jobs       |
| 10    | export       | ~1000 | JSONL sync            |

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
├── dependencies/# dep add, dep remove, dep tree ✅ DONE
├── labels/      # label add, label remove ✅ DONE
├── comments/    # comment add, comment list ✅ DONE
├── gates/       # gate create, gate show, gate wait ✅ DONE
├── molecules/   # mol create, mol pour, wisp ✅ DONE
├── epics/       # epic status ✅ DONE
├── compaction/  # compact
├── export/      # flush, import
├── sync/        # sync, daemon
├── config/      # config get, config set ✅ DONE
└── admin/       # doctor, audit, cleanup ✅ DONE (grouping only)
```

**Phase 5 Progress:**

| Feature      | Status | Notes |
|--------------|--------|-------|
| labels       | ✅ Done | Full migration |
| comments     | ✅ Done | Full migration |
| epics        | ✅ Done | Full migration |
| gates        | ✅ Done | Full migration |
| dependencies | ✅ Done | Full migration, includes relate/unrelate |
| config       | ✅ Done | Full migration |
| admin        | ✅ Done | Command group only; cleanup/compact/reset stay in main due to complex deps |
| molecules    | ✅ Done | Full migration - all subcommands in commands/molecules/ |
| sync         | ⏳ Pending | Large, complex dependencies |
| issues       | ⏳ Pending | Core commands, may stay in main |

**Molecules Migration - Complete:**

Phase 1 (Complete):
- ✅ Created `commands/shared/template/` package with types and functions
- ✅ Types: `Subgraph`, `CloneOptions`, `InstantiateResult`, `IssueDetailsFromShow`
- ✅ Functions: `LoadSubgraph`, `CloneSubgraph`, `ExtractAllVariables`, etc.
- ✅ `template.go` now uses type aliases and delegates to shared package
- ✅ Created `commands/molecules/` package with command group
- ✅ `mol.go` updated to use molecules package; subcommands use `molCmd.GetMolCmd()`
- ✅ Backward compatibility maintained - all mol_*.go files work without changes

Phase 2 (Complete):
- ✅ `mol_show.go` → `commands/molecules/show.go` (migrated with exported parallel analysis)
- ✅ `cook.go` → `commands/molecules/cook.go` (formula cooking infrastructure)
- ✅ `mol_bond.go` → `commands/molecules/bond.go` (polymorphic bonding operations)
- ✅ `pour.go` → `commands/molecules/pour.go` (persistent mol spawning)
- ✅ `wisp.go` → `commands/molecules/wisp.go` (ephemeral wisp management)

Phase 3 (Complete):
- ✅ `mol_burn.go` → `commands/molecules/burn.go`
- ✅ `mol_current.go` → `commands/molecules/current.go`
- ✅ `mol_distill.go` → `commands/molecules/distill.go`
- ✅ `mol_progress.go` → `commands/molecules/progress.go`
- ✅ `mol_ready_gated.go` → `commands/molecules/ready_gated.go`
- ✅ `mol_seed.go` → `commands/molecules/seed.go`
- ✅ `mol_squash.go` → `commands/molecules/squash.go`
- ✅ `mol_stale.go` → `commands/molecules/stale.go`

All mol_*.go files from main have been migrated to commands/molecules/.

**Migration Pattern:**
```go
// 1. Change package main → package <feature>
package gates

// 2. Import cli
import "github.com/steveyegge/beads/cmd/bd/cli"

// 3. Replace globals with cli.Get() methods
cliCtx := cli.Get()
cliCtx.GetDaemonClient()  // was: daemonClient
cliCtx.GetStore()         // was: store
cliCtx.GetRootCtx()       // was: rootCtx
cliCtx.IsJSONOutput()     // was: jsonOutput
cliCtx.CheckReadonly()    // was: CheckReadonly()
cliCtx.MarkDirty()        // was: markDirtyAndScheduleFlush()

// 4. Replace init() with Register()
func Register(root *cobra.Command) {
    root.AddCommand(featureCmd)
}
```

**Verification:** All CLI commands work identically via cobra aliases

---

### Handoff Summary (Updated 2026-01-26)

**Molecules Migration - COMPLETE**

All mol_*.go files have been migrated from `cmd/bd/` to `cmd/bd/commands/molecules/`.

**Current Structure:**
```
src/cmd/bd/commands/
├── admin/        ✅ Done
├── comments/     ✅ Done
├── config/       ✅ Done
├── dependencies/ ✅ Done
├── epics/        ✅ Done
├── gates/        ✅ Done
├── labels/       ✅ Done
├── molecules/    ✅ Done
│   ├── molecules.go  - command group, helpers, Register()
│   ├── show.go       - mol show with parallel analysis
│   ├── cook.go       - formula cooking infrastructure
│   ├── bond.go       - polymorphic bonding operations
│   ├── pour.go       - persistent mol spawning
│   ├── wisp.go       - ephemeral wisp management
│   ├── burn.go       - molecule deletion
│   ├── current.go    - current molecule state
│   ├── distill.go    - molecule extraction
│   ├── progress.go   - progress tracking
│   ├── ready_gated.go - gated ready checks
│   ├── seed.go       - molecule seeding
│   ├── squash.go     - molecule squashing
│   └── stale.go      - stale molecule detection
└── shared/
    ├── context.go
    ├── errors.go
    └── template/     ✅
        ├── types.go
        ├── operations.go
        └── variables.go
```

**Key Exports Available in molcmd:**
- `ResolveAndCookFormulaWithVars` - formula loading and cooking
- `CookFormulaToSubgraph`, `CookFormulaToSubgraphWithVars` - in-memory subgraph creation
- `AnalyzeMoleculeParallel`, `ParallelInfo` - parallel step analysis
- `MoleculeSubgraph`, `MoleculeLabel` - type aliases
- `GetMolCmd()` - returns mol command for subcommand registration
- `BondProtoMol`, `BondMolMol`, `BondProtoProto` - bonding operations
- `IsProtoIssue`, `BondResult` - proto checking and result type
- `SpawnMolecule`, `SpawnMoleculeWithOptions` - molecule spawning
- `FormatTimeAgo` - human-readable time formatting
- `WispListItem`, `WispListResult`, `WispGCResult` - wisp result types

**Verification:** Build passes, all tests pass

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

| Risk                     | Mitigation                                            |
| ------------------------ | ----------------------------------------------------- |
| Breaking storage queries | Keep old interface as deprecated alias                |
| Breaking RPC protocol    | Keep protocol.go unchanged, only restructure handlers |
| Breaking CLI scripts     | Use cobra aliases to maintain old command paths       |
| Frontend import breaks   | Use index.ts barrel exports                           |

**Rollback:** Tag before each phase (`v0.x.x-pre-phase-N`), each phase is independently reversible.

Bloat/Junk Analysis

1. CLI Flat Structure (Being Fixed)  


- 347 files in src/cmd/bd/ root - extreme width
- After migration: should be ~20-30 files in root, rest in commands/\*/  


2. RPC Layer - 14,866 LOC  


src/internal/rpc/  
 ├── server_issues_epics.go # 2,367 lines - horizontal bloat  
 ├── server_labels_deps_comments.go  
 ├── server_routing_validation_diagnostics.go  
 └── ... (38 files total)  
 After migration: RPC handlers become thin dispatchers calling features/\*/service.go. Can delete most of this once features  
 have their own services.

3. Storage Interface - 83 methods  


- src/internal/storage/storage.go - God interface (SRP violation)
- After migration: Split into feature-specific repositories (~10 methods each in features/\*/repository.go)  


4. Dolt Backend - Possibly Dead?  


src/cmd/bd/dolt\_\*.go # 12 files  
 src/internal/storage/dolt/ # Full backend implementation

- If SQLite is primary and Dolt is unused → can delete ~2,000+ LOC
- Confirm with: "Is Dolt backend still actively used?"  


5. node_modules in src (214MB)  


- src/cmd/server/frontend/node_modules/
- Not in git (good), but creates 929 directories locally
- Add to .gitignore if not already  


6. Migration Scripts - Keep as Admin Tools  


src/cmd/bd/migrate\_\*.go # 10 files

- Per earlier decision: move to commands/admin/ not delete  


What Gets Deleted After Full Migration  
 ┌────────────────────────────┬─────────┬──────────────────────────┐  
 │ What │ LOC │ Status │  
 ├────────────────────────────┼─────────┼──────────────────────────┤  
 │ RPC server**.go handlers │ ~10,000 │ Becomes thin dispatchers │  
 ├────────────────────────────┼─────────┼──────────────────────────┤  
 │ Storage interface bloat │ ~500 │ Split into feature repos │  
 ├────────────────────────────┼─────────┼──────────────────────────┤  
 │ Dolt backend (if unused) │ ~2,000 │ Confirm first │  
 ├────────────────────────────┼─────────┼──────────────────────────┤  
 │ Old horizontal cmd/bd/*.go │ ~15,000 │ Replaced by commands/\*/ │  
 └────────────────────────────┴─────────┴──────────────────────────┘  
 What Stays but Gets Reorganized  
 ┌─────────────────────────────┬──────────────────────────────┐  
 │ From │ To │  
 ├─────────────────────────────┼──────────────────────────────┤  
 │ internal/rpc/server*_.go │ features/_/rpc.go (thin) │  
 ├─────────────────────────────┼──────────────────────────────┤  
 │ internal/storage/storage.go │ features/_/repository.go │  
 ├─────────────────────────────┼──────────────────────────────┤  
 │ cmd/bd/_.go (347 files) │ cmd/bd/commands/\*/ (10 dirs) │  
 └─────────────────────────────┴──────────────────────────────┘  
 TL;DR: The main junk is architectural bloat (horizontal layers) not leftover files. The beads-upstream/ and  
 spec-kitty-upstream/ directories mentioned in the migration doc appear to have been deleted already.
