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
├── molecules/   # mol create, mol pour, wisp 🔶 PARTIAL (cmd group done)
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
| molecules    | 🔶 Partial | Command group migrated; subcommands remain in main (~4200 LOC) |
| sync         | ⏳ Pending | Large, complex dependencies |
| issues       | ⏳ Pending | Core commands, may stay in main |

**Molecules Migration - Partial (Phase 1 Complete):**

Phase 1 (Complete):
- ✅ Created `commands/shared/template/` package with types and functions
- ✅ Types: `Subgraph`, `CloneOptions`, `InstantiateResult`, `IssueDetailsFromShow`
- ✅ Functions: `LoadSubgraph`, `CloneSubgraph`, `ExtractAllVariables`, etc.
- ✅ `template.go` now uses type aliases and delegates to shared package
- ✅ Created `commands/molecules/` package with command group
- ✅ `mol.go` updated to use molecules package; subcommands use `molCmd.GetMolCmd()`
- ✅ Backward compatibility maintained - all mol_*.go files work without changes

Phase 2 (In Progress):
- ✅ `mol_show.go` → `commands/molecules/show.go` (migrated with exported parallel analysis)
- ✅ `cook.go` → `commands/molecules/cook.go` (formula cooking infrastructure)

Phase 3 (Remaining - ~2700 LOC):
Subcommand files now unblocked:
- `pour.go`, `wisp.go` - now can be migrated (use molcmd.ResolveAndCookFormulaWithVars)
- `mol_bond.go` - core bonding logic, uses molcmd
- `mol_burn.go`, `mol_current.go`, `mol_distill.go`, `mol_progress.go`
- `mol_ready_gated.go`, `mol_seed.go`, `mol_squash.go`, `mol_stale.go`

Migration order (suggested):
1. ~~cook.go (formula cooking infrastructure)~~ ✅ DONE
2. mol_bond.go (bonding operations)
3. pour.go, wisp.go (spawn commands)
4. Remaining mol_*.go files

Files in main already updated to use `molcmd.` imports for cook functions.

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

**What Was Done:**
1. Created `commands/shared/template/` package (~300 LOC) - types and functions for template operations
2. Created `commands/molecules/` package with:
   - `molecules.go` - command group, helpers, Register()
   - `show.go` - mol show with parallel analysis (~500 LOC)
   - `cook.go` - formula cooking infrastructure (~1060 LOC)
3. Updated `template.go` in main to use type aliases and delegation
4. Updated dependent files in main (`pour.go`, `wisp.go`, `mol_bond.go`, `mol_seed.go`) to use `molcmd.` imports
5. All tests passing, build clean

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
├── molecules/    🔶 Partial
│   ├── molecules.go  ✅
│   ├── show.go       ✅
│   └── cook.go       ✅
└── shared/
    ├── context.go
    ├── errors.go
    └── template/     ✅
        ├── types.go
        ├── operations.go
        └── variables.go
```

**What Needs To Be Done (Remaining ~2700 LOC):**

Migration order:
1. `mol_bond.go` (~600 LOC) - core bonding logic, uses molcmd
2. `pour.go` + `wisp.go` (~500 LOC) - spawn commands, already using molcmd imports
3. Remaining mol_*.go files:
   - `mol_burn.go`, `mol_current.go`, `mol_distill.go`
   - `mol_progress.go`, `mol_ready_gated.go`, `mol_seed.go`
   - `mol_squash.go`, `mol_stale.go`

**Migration Pattern for Each File:**
1. Create new file in `commands/molecules/` with `package molecules`
2. Import `cli` and use `cliCtx := cli.Get()` for globals
3. Import `template` from `commands/shared/template` for template functions
4. Export functions that other files in main depend on (prefix with capital letter)
5. Update `molecules.go` Register() to call the new registerXCmd() function
6. Update any callers in main to use `molcmd.ExportedFunction`
7. Delete old file from main
8. Run `go build ./cmd/bd/...` and `go test ./cmd/bd/... -short` to verify

**Key Exports Already Available in molcmd:**
- `ResolveAndCookFormulaWithVars` - formula loading and cooking
- `CookFormulaToSubgraph`, `CookFormulaToSubgraphWithVars` - in-memory subgraph creation
- `AnalyzeMoleculeParallel`, `ParallelInfo` - parallel step analysis
- `MoleculeSubgraph`, `MoleculeLabel` - type aliases
- `GetMolCmd()` - returns mol command for subcommand registration

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
