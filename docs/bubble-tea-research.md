# Bubble Tea TUI Enhancement for kitty-beads

## Overview

Bubble Tea is a Go TUI framework. It outputs ANSI to terminal. Our xterm.js already renders ANSI. Zero frontend changes needed.

**We already use Charm ecosystem**: `huh` in `create-form.go`

---

## Architecture Fit

```
Current:   React (xterm.js) ←── WebSocket ←── PTY ←── /bin/bash
Enhanced:  React (xterm.js) ←── WebSocket ←── PTY ←── bd list --tui
                                                       ↑
                                               Bubble Tea renders here
```

---

## Enhanced Commands

### `bd list --tui` → Interactive Issue Browser

```
┌──────────────────────────────────────────────────────────────────────┐
│ Issues (42)                                    Filter: bug█          │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  > ● P1  kitty-beads-qnz     Unified Developer-AI Collabora...  epic │
│    ○ P1  kitty-beads-qnz.4   Phase 4: Linear Center Command...  task │
│    ◐ P2  kitty-beads-2ds     Developer Collaboration Hub       epic │
│    ● P0  bd-c4rq             Fix stale database detection       bug  │
│    ○ P2  bd-7zka             Gate issues not filtering          task │
│    ✓ P3  bd-old              Completed feature                  task │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ j/k navigate  /filter  Enter view  c close  e edit  p priority  ?   │
└──────────────────────────────────────────────────────────────────────┘
```

**Detail View (on Enter)**:
```
┌──────────────────────────────────────────────────────────────────────┐
│ bd-c4rq                                                    ← Back    │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Fix stale database detection                                        │
│  ════════════════════════════                                        │
│                                                                      │
│  Status: ● blocked    Priority: P0    Type: bug                      │
│  Assignee: @akimzmerli    Labels: backend, urgent                    │
│                                                                      │
│  ──────────────────────────────────────────────────────────────────  │
│  Description:                                                        │
│                                                                      │
│  Database freshness check fails silently when daemon is running.     │
│  Users see stale data without warning.                               │
│                                                                      │
│  Steps to reproduce:                                                 │
│  1. Start daemon                                                     │
│  2. Modify issues.jsonl externally                                   │
│  3. Run bd list - shows old data                                     │
│                                                                      │
│  ──────────────────────────────────────────────────────────────────  │
│  Blocked by:                                                         │
│    ○ bd-daemon-1  Daemon auto-import on staleness                    │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ c close  e edit  a assign  l labels  d deps  Esc back               │
└──────────────────────────────────────────────────────────────────────┘
```

---

### `bd graph --tui` → Interactive Dependency Explorer

```
┌─────────────────────────────────────────────────────────────────────┐
│ Dependency Graph: kitty-beads-qnz                      [box/tree]   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ▼ ● kitty-beads-qnz [epic] Unified Developer-AI Collaboration     │
│    │                                                                │
│    ├── ✓ kitty-beads-qnz.1 [task] Phase 1: Base Layout             │
│    │                                                                │
│    ├── ✓ kitty-beads-qnz.2 [task] Phase 2: Issue Panel             │
│    │                                                                │
│    ├── ◐ kitty-beads-qnz.3 [task] Phase 3: Kanban View             │
│    │   │                                                            │
│    │   └── ● bd-c4rq [bug] Fix stale database  ◀── blocker         │
│    │                                                                │
│    └── ▶ kitty-beads-qnz.4 [task] Phase 4: Command Palette...      │
│        (3 children hidden)                                          │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ ▶/▼ expand  j/k nav  Enter details  b blockers  f focus subtree    │
└─────────────────────────────────────────────────────────────────────┘
```

**Focused Subtree View**:
```
┌─────────────────────────────────────────────────────────────────────┐
│ Focus: kitty-beads-qnz.3                               Esc unfocus  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ◐ kitty-beads-qnz.3 [task] Phase 3: Kanban View                   │
│    │                                                                │
│    ├── ● bd-c4rq [bug] Fix stale database                          │
│    │   └── ○ bd-daemon-1 [task] Daemon auto-import                 │
│    │                                                                │
│    ├── ○ bd-kanban-1 [task] Drag-drop reordering                   │
│    │                                                                │
│    └── ○ bd-kanban-2 [task] Swimlane view                          │
│                                                                     │
│  Legend: ○ open  ◐ in_progress  ● blocked  ✓ closed  ❄ deferred   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ ▶/▼ expand  j/k nav  Enter details  p show parent  Esc unfocus     │
└─────────────────────────────────────────────────────────────────────┘
```

---

### `bd status --tui` → Live Dashboard

```
┌─────────────────────────────────────────────────────────────────────┐
│ 📊 kitty-beads Status                              ↻ 5s auto       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Summary                              Activity (24h)                │
│  ───────────────────────────────      ──────────────────────────   │
│  Total:     127                       ▁▂▄▆█▇▅▃▂▁ commits           │
│                                                                     │
│  Open:      23  ████████████░░░░░░    Created:  +12                │
│  Progress:   8  █████░░░░░░░░░░░░░    Closed:    -8                │
│  Blocked:    3  ██░░░░░░░░░░░░░░░░    Net:       +4                │
│  Closed:    93  ██████████████████                                 │
│                                                                     │
│  Ready Work: 18                       Avg Lead Time: 4.2h          │
│                                                                     │
│  ──────────────────────────────────────────────────────────────    │
│  ⚠ Alerts                                                          │
│    • 2 epics ready to close                                        │
│    • 3 issues overdue                                              │
│    • 1 issue blocked >48h                                          │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ r refresh  l list ready  o overdue  e epics  q quit                │
└─────────────────────────────────────────────────────────────────────┘
```

---

### `bd create-form` → Enhanced with Dependency Picker

**Step 3: Add Dependencies**:
```
┌─────────────────────────────────────────────────────────────────────┐
│ Create Issue                                           Step 3/4    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Add Dependencies                         Search: auth█            │
│  ═══════════════                                                   │
│                                                                     │
│  Type: [blocks ▼]                                                  │
│                                                                     │
│  > ○ bd-auth-1    Implement OAuth2 flow                 P2  task   │
│    ○ bd-auth-2    Add refresh token logic               P1  task   │
│    ● bd-auth-3    Session management refactor           P1  task   │
│    ○ bd-login-1   Login page redesign                   P3  task   │
│                                                                     │
│  ──────────────────────────────────────────────────────────────    │
│  Selected Dependencies:                                            │
│    • blocks: bd-auth-3                                             │
│    • discovered-from: bd-login-1                                   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ Space select  / filter  Tab type  Enter next step  Esc back        │
└─────────────────────────────────────────────────────────────────────┘
```

**Step 4: Preview**:
```
┌─────────────────────────────────────────────────────────────────────┐
│ Create Issue                                           Step 4/4    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Preview                                                           │
│  ═══════                                                           │
│                                                                     │
│  Title:       Fix authentication timeout in login handler          │
│  Type:        bug                                                  │
│  Priority:    P1                                                   │
│  Assignee:    @akimzmerli                                          │
│  Labels:      backend, auth                                        │
│                                                                     │
│  Description:                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Users report being logged out after 5 minutes even with     │   │
│  │ "remember me" checked. Need to investigate session token    │   │
│  │ refresh logic.                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  Dependencies:                                                     │
│    • blocks: bd-auth-3 (Session management refactor)               │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ Enter create  e edit  Esc cancel                                   │
└─────────────────────────────────────────────────────────────────────┘
```

---

### `bd tui` → Full Linear-Style Interface

**Main View**:
```
┌─────────────────────────────────────────────────────────────────────┐
│ kitty-beads                    [Issues] [Graph] [Status]    ⌘K     │
├────────────────────────────────┬────────────────────────────────────┤
│                                │                                    │
│  Filter: █                     │  bd-c4rq                          │
│                                │  ═════════════════════════════    │
│  > ● P0  bd-c4rq         bug   │                                    │
│    ○ P1  bd-qnz.4        task  │  Fix stale database detection     │
│    ◐ P2  bd-2ds          epic  │                                    │
│    ○ P2  bd-7zka         task  │  Status: ● blocked                │
│    ○ P2  bd-kanban-1     task  │  Priority: P0                     │
│    ○ P3  bd-kanban-2     task  │  Assignee: @akimzmerli            │
│    ✓ P1  bd-qnz.1        task  │                                    │
│    ✓ P1  bd-qnz.2        task  │  ────────────────────────────     │
│                                │  Database freshness check fails   │
│                                │  silently when daemon is running. │
│                                │                                    │
│                                │  Blocked by:                      │
│                                │    ○ bd-daemon-1                  │
│                                │                                    │
├────────────────────────────────┴────────────────────────────────────┤
│ j/k nav  Enter focus  c close  n new  / filter  Tab pane  ⌘K cmd   │
└─────────────────────────────────────────────────────────────────────┘
```

**Command Palette (⌘K)**:
```
┌─────────────────────────────────────────────────────────────────────┐
│                              ⌘K                                     │
├─────────────────────────────────────────────────────────────────────┤
│  █                                                                  │
│                                                                     │
│  > Create new issue                                        ⌘N      │
│    Close issue                                             ⌘W      │
│    Assign to me                                            ⌘M      │
│    ────────────────────────────────────────────────────           │
│    View ready work                                         ⌘R      │
│    Show my issues                                          ⌘I      │
│    ────────────────────────────────────────────────────           │
│    Switch to Graph view                                    ⌘G      │
│    Switch to Status view                                   ⌘S      │
│    ────────────────────────────────────────────────────           │
│    Run sync                                                ⌘⇧S     │
│    Open in browser                                         ⌘O      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

### `bd daemon --monitor` → Process Monitor

```
┌─────────────────────────────────────────────────────────────────────┐
│ 🔧 Beads Daemon                                        ● Running    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Uptime: 2h 34m 12s       │  Memory: 24.5 MB                       │
│  Connections: 3           │  Pending: 0                            │
│  Requests/min: 12         │  Last sync: 45s ago                    │
│                                                                     │
│  ──────────────────────────────────────────────────────────────    │
│  Recent Activity                                                   │
│  ──────────────────────────────────────────────────────────────    │
│  14:32:01  ✓ Auto-import completed (12 issues)                     │
│  14:31:45  ◐ Sync started                                          │
│  14:30:12  ✓ RPC: List (filter: status=open)                       │
│  14:29:58  ✓ RPC: Create (bd-new-1)                                │
│  14:29:30  ✓ Health check passed                                   │
│  14:28:15  ✓ RPC: Close (bd-old-3)                                 │
│                                                                     │
│  ──────────────────────────────────────────────────────────────    │
│  Logs (tail)                                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ [INFO] Client connected from /var/run/bd.sock               │   │
│  │ [DEBUG] Query: SELECT * FROM issues WHERE status != closed  │   │
│  │ [INFO] Returned 42 issues in 3ms                            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ s sync  r restart  l logs  c compact  q quit                       │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Bubbles Components Mapping

| Component | Use Case |
|-----------|----------|
| `list` | Issue browser, dependency picker |
| `textinput` | Filter box, search |
| `textarea` | Description editor |
| `viewport` | Scrollable details, logs |
| `progress` | Status bars |
| `spinner` | Loading states |
| `table` | Structured data |
| `help` | Keybinding bar |
| `timer` | Auto-refresh |

---

## Production Usage

- CockroachDB tooling
- gh-dash (GitHub)
- chezmoi (dotfiles)
- Trufflehog (security)
- Glow (markdown)
