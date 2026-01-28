# Kitty-Beads: Rating & Improvements

## Architecture Ratings

### Frontend: 7.5/10

| Category | Score | Assessment |
|----------|-------|------------|
| **Structure** | 8/10 | Feature-based organization, proper separation |
| **TypeScript** | 7/10 | Good typing, some areas could be stricter |
| **State Management** | 8/10 | React Query + Context appropriate for scale |
| **Code Quality** | 7/10 | Clean components, some lint issues |
| **Reusability** | 8/10 | Good barrel exports, shared UI components |
| **Testing** | 4/10 | No test files observed |

### Backend: 8/10

| Category | Score | Assessment |
|----------|-------|------------|
| **Structure** | 8/10 | Vertical slices in features/, solid internal/ |
| **Test Coverage** | 9/10 | 355 test files / 866 total (41%) |
| **Code Quality** | 7/10 | Some god files need splitting |
| **API Design** | 8/10 | RESTful, proper status codes |
| **Separation of Concerns** | 8/10 | Handler → Service → Repository |
| **CLI Tooling** | 9/10 | Comprehensive 322-command CLI |

### Overall: 7.5/10

| Category | Score | Assessment |
|----------|-------|------------|
| **Product Vision** | 9/10 | Clear value prop, solves real problem |
| **Architecture Coherence** | 8/10 | FE + BE aligned on vertical slices |
| **Developer Experience** | 7/10 | Good CLI, missing E2E tests |
| **Production Readiness** | 6/10 | MVP-ready, needs error handling polish |

---

## Identified Issues

### Critical (P1)

#### 1. God Files in Backend RPC
**File:** `internal/rpc/server_issues_epics.go` (2,367 lines)

**Problem:** Single file handling all issue and epic operations violates single responsibility.

**Impact:**
- Hard to navigate
- Merge conflicts likely
- Testing becomes complex

**Recommendation:**
```
internal/rpc/
├── server_issues.go         # Issue CRUD
├── server_issues_filters.go # List/query operations
├── server_epics.go          # Epic-specific handlers
├── server_hierarchy.go      # Parent-child operations
└── server_bulk.go           # Bulk operations
```

#### 2. React Lint Errors
**Files:** `App.tsx`, `features/kanban/components/IssueModal.tsx`

**Problem:** `setState` called directly in `useEffect` causing render cascades.

**Current:**
```tsx
useEffect(() => {
  if (!currentFeatureId && features.length > 0) {
    setCurrentFeatureId(features[0].id);  // Lint error
  }
}, [features, currentFeatureId]);
```

**Fix:**
```tsx
// Option 1: Initialize in useState
const [currentFeatureId, setCurrentFeatureId] = useState<string | null>(
  () => features[0]?.id ?? null
);

// Option 2: useMemo for derived state
const effectiveFeatureId = useMemo(
  () => currentFeatureId ?? features[0]?.id ?? null,
  [currentFeatureId, features]
);
```

---

### High (P2)

#### 3. No Error Boundaries
**Problem:** React errors crash entire app.

**Recommendation:**
```tsx
// src/components/ErrorBoundary.tsx
export class ErrorBoundary extends Component<Props, State> {
  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} />;
    }
    return this.props.children;
  }
}

// Wrap in App.tsx
<ErrorBoundary>
  <AppContent />
</ErrorBoundary>
```

#### 4. Duplicate Utility Packages
**Problem:** Both `internal/util/` and `internal/utils/` exist.

**Recommendation:** Consolidate into single `internal/util/` package.

#### 5. Missing Frontend Tests
**Problem:** No test files in frontend.

**Recommendation:** Add Vitest + React Testing Library:
```bash
npm install -D vitest @testing-library/react @testing-library/jest-dom
```

Priority test targets:
- `features/terminal/context.tsx` (state logic)
- `lib/api.ts` (API functions)
- `lib/planParser.ts` (parsing logic)

---

### Medium (P3)

#### 6. CLI Commands in Flat Directory
**Problem:** 322 `.go` files in `cmd/bd/` - hard to navigate.

**Recommendation:**
```
cmd/bd/
├── main.go
├── core/           # create, list, show, close
├── sync/           # sync_git, sync_branch, linear_sync
├── migrate/        # migrate_*, compact_*
├── daemon/         # daemon_*, activity_*
└── admin/          # doctor, audit, formula
```

#### 7. No Route-Level Code Splitting
**Problem:** All pages loaded upfront.

**Recommendation:**
```tsx
// App.tsx
const Kanban = lazy(() => import('./pages/Kanban'));
const Roadmap = lazy(() => import('./pages/Roadmap'));

<Suspense fallback={<PageLoader />}>
  <Route path="/kanban" element={<Kanban />} />
</Suspense>
```

#### 8. Favicon Still Vite Default
**File:** `public/vite.svg`, referenced in `index.html`

**Recommendation:** Replace with branded Kitty-Beads icon.

---

### Low (P4)

#### 9. Hardcoded Layout Dimensions
**Files:** `layouts/Layout.tsx`, `layouts/Sidebar.tsx`

**Problem:**
```tsx
const HEADER_HEIGHT = 73; // Duplicated
const SIDEBAR_WIDTH = 224;
```

**Recommendation:** Create `src/constants/layout.ts`:
```ts
export const LAYOUT = {
  HEADER_HEIGHT: 73,
  SIDEBAR_WIDTH: 224,
  TERMINAL_MIN_HEIGHT: 150,
} as const;
```

#### 10. API Error Handling Inconsistent
**Problem:** Some handlers check `strings.Contains(err.Error(), "not found")`.

**Recommendation:** Use typed errors:
```go
var ErrNotFound = errors.New("not found")

// In handler
if errors.Is(err, ErrNotFound) {
    writeError(w, http.StatusNotFound, err.Error())
}
```

---

## Improvement Roadmap

### Phase 1: Quick Wins (1-2 days)
- [ ] Fix React lint errors in `App.tsx` and `IssueModal.tsx`
- [ ] Replace Vite favicon
- [ ] Add error boundary component
- [ ] Consolidate `util/` + `utils/`
- [ ] Extract layout constants

### Phase 2: Code Quality (3-5 days)
- [ ] Split `server_issues_epics.go` into focused files
- [ ] Add Vitest setup + first tests for `planParser.ts`
- [ ] Group CLI commands into subdirectories
- [ ] Add typed errors in backend handlers

### Phase 3: Production Hardening (1-2 weeks)
- [ ] Add React lazy loading for pages
- [ ] Add E2E tests with Playwright
- [ ] Add API request retry/timeout handling
- [ ] Add health check endpoint
- [ ] Add structured logging (backend)

### Phase 4: DX Improvements (ongoing)
- [ ] Add Storybook for UI components
- [ ] Add OpenAPI/Swagger docs
- [ ] Add commit hooks (lint-staged, husky)
- [ ] Add CI/CD pipeline

---

## Metrics to Track

| Metric | Current | Target |
|--------|---------|--------|
| Frontend test coverage | 0% | 60% |
| Backend test coverage | ~41% (by file) | 70% |
| Largest file (backend) | 2,367 lines | <500 lines |
| Lint errors | 2 | 0 |
| Bundle size (gzipped) | Unknown | <200KB |
| Lighthouse score | Unknown | >90 |

---

## Summary

**Strengths to Preserve:**
- Vertical slice architecture (both FE and BE)
- Excellent backend test coverage
- Comprehensive CLI tooling
- Clean barrel exports in frontend
- Local-first SQLite approach

**Top 3 Actions:**
1. Split `server_issues_epics.go` - biggest maintenance risk
2. Fix React lint errors - prevents runtime bugs
3. Add error boundaries - prevents full-app crashes

The codebase is solid for an MVP. Focus on the P1/P2 items before adding new features.
