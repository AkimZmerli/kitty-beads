# CLI Reorganization - Phase 5 Analysis

## Current State

**Problem:** 350 `.go` files flat in `src/cmd/bd/` - not developer friendly.

```
src/cmd/bd/
├── main.go
├── create.go
├── close.go
├── daemon.go
├── ... (347 more files)
```

## File Categorization

| Feature | Files | Examples |
|---------|-------|----------|
| **sync/** | ~67 | `daemon*.go`, `sync*.go` |
| **issues/** | ~25 | `create.go`, `close.go`, `delete.go`, `list.go`, `show.go`, `update.go`, `ready.go` |
| **export/** | ~20 | `export*.go`, `import*.go`, `flush*.go` |
| **admin/** | ~20 | `doctor*.go`, `audit.go`, `cleanup.go`, `repair*.go`, `migrate*.go` |
| **molecules/** | ~18 | `mol*.go`, `pour.go`, `wisp.go`, `cook.go`, `formula.go` |
| **dependencies/** | ~6 | `dep.go`, `graph.go`, `relate.go` |
| **config/** | ~4 | `config.go`, `setup.go` |
| **gates/** | ~4 | `gate*.go` |
| **labels/** | ~3 | `label*.go` |
| **comments/** | ~3 | `comments*.go` |
| **epics/** | ~3 | `epic*.go` |
| **compaction/** | ~5 | `compact*.go` |
| **core/** | ~170+ | `main.go`, `version.go`, `types.go`, helpers, tests |

## Target Structure

```
cmd/bd/
├── main.go                    # Root cmd + global setup
├── commands/
│   ├── shared/               # ✅ Created - shared state & helpers
│   │   ├── context.go        # Runtime state (store, actor, ctx)
│   │   └── errors.go         # FatalError, CheckReadonly, etc.
│   ├── issues/
│   │   ├── create.go
│   │   ├── close.go
│   │   ├── list.go
│   │   └── register.go       # func Register(root *cobra.Command)
│   ├── sync/
│   │   ├── daemon.go
│   │   ├── sync.go
│   │   └── register.go
│   ├── molecules/
│   ├── admin/
│   └── ...
```

## Challenges Discovered

### 1. Heavy Global Coupling

All 350 files are `package main` and share these globals:

```go
var (
    store        storage.Storage   // Used by ~200 files
    actor        string            // Used by ~150 files
    rootCtx      context.Context   // Used by ~100 files
    daemonClient *rpc.Client       // Used by ~80 files
    jsonOutput   bool              // Used by ~60 files
    hookRunner   *hooks.Runner     // Used by ~40 files
)
```

### 2. Shared Helper Functions

Many helpers are used across feature boundaries:

| Helper | Used By |
|--------|---------|
| `markDirtyAndScheduleFlush()` | 35 files |
| `FatalError()` | 50+ files |
| `outputJSON()` | 40+ files |
| `GetLastTouchedID()` | 15 files |
| `needsRouting()` | 20 files |
| `issueIDCompletion` | 25 files |

### 3. Go Package Constraints

Go doesn't allow the same package across multiple directories. Options:

1. **Subpackages** - Each feature is its own package, imports shared
2. **Single package** - Keep all in `main`, use file naming convention
3. **Internal package** - Move commands to `internal/cli/*` packages

## Implementation Options

### Option A: Full Vertical Slice (2-3 days)

**Pros:**
- Clean architecture
- Each feature isolated
- Easier testing per feature

**Cons:**
- Large refactor
- Must move 50+ helper functions to shared
- Risk of breaking changes

**Steps:**
1. Move all globals to `commands/shared/context.go`
2. Move all helpers to `commands/shared/`
3. Create feature packages one by one
4. Each package exports `Register(root *cobra.Command)`
5. Update `main.go` to call all Register functions

### Option B: Pragmatic Naming (1 hour)

**Pros:**
- Quick win
- No code changes
- IDE navigation works

**Cons:**
- Still 350 files in one dir
- No compile-time isolation

**Steps:**
1. Add `_categories.go` documenting file groupings
2. Use consistent prefixes (already partially done: `daemon_*.go`, `sync_*.go`)
3. Add editor config for grouping

### Option C: Hybrid (1 day)

**Pros:**
- Proves the pattern
- Incremental progress
- Lower risk

**Cons:**
- Inconsistent structure temporarily

**Steps:**
1. Fully extract ONE feature (e.g., `labels` - simplest)
2. Document the pattern
3. Gradually migrate others

## Progress So Far

### Created

- [x] `commands/` directory structure (12 feature dirs)
- [x] `commands/shared/context.go` - Runtime state management
- [x] `commands/shared/errors.go` - Error helpers

### Verified

- [x] Shared package compiles: `go build ./cmd/bd/commands/shared/...`

## Recommendation

Start with **Option C (Hybrid)** - extract `labels` as proof of concept since it's the simplest (~3 files, minimal dependencies). This validates the architecture before committing to the full migration.

## Next Steps

1. Extract `labels` package as POC
2. If successful, extract `comments` (similar complexity)
3. Then tackle `issues` (moderate complexity)
4. Leave `sync` for last (highest complexity with 67 files)

---

*Document created: 2026-01-25*
*Epic: kitty-beads-kcw (Vertical Slice Architecture Migration)*
*Phase: 5 - CLI Reorganization*
