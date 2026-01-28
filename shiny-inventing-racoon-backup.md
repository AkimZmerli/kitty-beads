# Kitty-Beads: Developer Collaboration Hub

## Vision

A collaboration hub for friends working together on projects. Combines visual project mapping, focused planning, and AI-assisted iteration in one place.

**Context**: Personal friend collab project. Whiteboard + full features are worth shipping.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     KITTY-BEADS COLLAB HUB                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐               │
│  │ TREE GRAPH  │   │   KANBAN    │   │ WHITEBOARD  │               │
│  │  (default)  │   │   BOARD     │   │ (brainstorm)│               │
│  └──────┬──────┘   └─────────────┘   └─────────────┘               │
│         │                                                           │
│         ▼ click node                                                │
│  ┌─────────────────────────────────────────────────────────┐       │
│  │              BEAD CARD (popup)                          │       │
│  │  ┌─────────────────────────────────────────────────┐    │       │
│  │  │ bd-zpw.2: Implement user registration API  ●P1  │    │       │
│  │  ├─────────────────────────────────────────────────┤    │       │
│  │  │ Summary of plan...                              │    │       │
│  │  │ ✓ 3 acceptance criteria  •  2 open questions    │    │       │
│  │  │                                                 │    │       │
│  │  │              [Edit Plan]  [View Full]           │    │       │
│  │  └─────────────────────────────────────────────────┘    │       │
│  └─────────────────────────────────────────────────────────┘       │
│         │                        │                                  │
│         ▼                        ▼                                  │
│  ┌─────────────┐          ┌─────────────┐                          │
│  │  IDEATION   │          │ EPIC PLAN   │                          │
│  │    PAD      │          │   (modal)   │                          │
│  └─────────────┘          └─────────────┘                          │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  TERMINAL (docked)                                          │   │
│  │  $ bd plan get bd-zpw.2 | claude "improve acceptance..."    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  [ Cmd+K → Linear Center ]  [ Activity Feed ]  [ Assignments ]     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. Tree Graph Map (Default Landing)

Visual overview of project hierarchy - epics → tasks → subtasks

```
                    ┌─────────────┐
                    │  bd-zpw     │  ← Epic
                    │  Auth Epic  │
                    └──────┬──────┘
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
    ┌───────────┐   ┌───────────┐   ┌───────────┐
    │ bd-zpw.1  │   │ bd-zpw.2  │   │ bd-zpw.3  │  ← Tasks
    │ Login UI  │   │ Reg API   │   │ Pwd Reset │
    └─────┬─────┘   └───────────┘   └───────────┘
    ┌─────┼─────┐
    ▼     ▼     ▼
  ┌───┐ ┌───┐ ┌───┐
  │.1 │ │.2 │ │.3 │  ← Subtasks
  └───┘ └───┘ └───┘
```

**Interactions**:

- Click node → Bead Card popup appears
- Read-only (no drag/drop reparenting)
- Color-coded by priority or status

**Bead Card** (on node click):

- ID, title, priority, plan summary
- Acceptance criteria count, open questions count
- [Edit Plan] → Opens Ideation Pad
- [View Full] → Shows parent epic's full plan

---

### 2. Kanban Board

Phase-based overview - what's where

```
┌─────────────┬─────────────┬─────────────┬─────────────┐
│  IDEATION   │  PLANNING   │ IN PROGRESS │    DONE     │
├─────────────┼─────────────┼─────────────┼─────────────┤
│  bd-abc.1   │  bd-zpw.3   │  bd-zpw.2   │  bd-xyz.1   │
│  bd-def.2   │             │  bd-zpw.1   │  bd-xyz.2   │
└─────────────┴─────────────┴─────────────┴─────────────┘
```

- Click card → Bead Card popup
- Drag between columns to change phase
- Filter by epic, assignee, priority

---

### 3. Whiteboard (Brainstorming)

Freeform visual thinking space (tldraw-based)

- Basic shapes, arrows, text, freehand
- Stores as JSON blob per bead or epic
- NOT direct bead creation from shapes
- **Auto-send to Claude**: [Export to Claude] button → `bd whiteboard export <bead-id> | claude` in terminal
- Bead attachment (storing whiteboard as artifact) → second pass, validate in practice first
- Export as image

---

### 4. Ideation Pad (Plan Editor)

Focused markdown editor for bead plans

```
┌─────────────────────────────────────────────────────────────────┐
│  IDEATION PAD                          Editing: bd-zpw.2        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  # User Registration API                                        │
│                                                                 │
│  Mobile-first registration flow with OAuth...                   │
│                                                                 │
│  ## Acceptance Criteria                                         │
│  - [ ] OAuth login with Google                                  │
│  - [ ] OAuth login with GitHub                                  │
│                                                                 │
│  ## Open Questions                                              │
│  - Rate limit: 5/min or 10/min?                                 │
│                                                                 │
│                                              [Save] [Cancel]    │
└─────────────────────────────────────────────────────────────────┘
```

- Markdown editor (CodeMirror)
- Auto-extracts first paragraph as summary
- Parses acceptance criteria and open questions for counts
- Auto-refresh when updated externally (via CLI)

---

### 5. Terminal + Claude Bridge

AI reads/writes plans directly from terminal

**CLI commands**:

```bash
bd plan get <bead-id>      # Output plan markdown
bd plan set <bead-id> --file=<path>   # Update plan
```

**Workflow**:

```bash
# One-liner: Claude improves plan
bd plan get bd-zpw.2 | claude "add error handling criteria" | bd plan set bd-zpw.2
```

Ideation Pad auto-refreshes when backend changes.

---

### 6. Activity Feed

Async collaboration - see what happened

- Plan edits, status changes, comments, assignments
- Filter by bead, user, type
- Click entry → jump to bead

---

### 7. Comments

Discussion threads per bead

- Accessed via bead card
- Markdown support
- @mentions

---

### 8. Assignments

Who owns what

- Assign beads to team members
- Filter views by assignee
- Shows in bead card and kanban

---

### 9. Linear Center (Cmd+K)

Quick navigation and actions

- Search beads fuzzy
- Jump to views (Tree, Kanban, Whiteboard, Activity)
- Actions: New Bead, Edit Plan, Toggle Terminal

---

## Two Workflows

### A: Top-Down (Epic-Driven)

```
Epic defined → Tree Graph shows structure → Pick bead →
Edit Plan → Iterate with Claude → Implement
```

### B: Bottom-Up (Brainstorming)

```
Whiteboard brainstorm → Screenshot to Claude → Create beads →
Organize into epic → Refine plans → Implement
```

---

## Implementation Priority

1. **Tree Graph + Bead Card** - The visual map (core navigation)
2. **Kanban Enhancement** - Drag-drop between phase columns
3. **Ideation Pad + Claude Bridge CLI** - `bd plan get/set` with auto-refresh
4. **Whiteboard** - tldraw integration with auto-send to Claude button
5. **Activity Feed** - Read-only collaboration log
6. **Comments UI** - Discussion threads + @mentions

---

## Migration: Current → New

### KEEP AS-IS (Tokyo Night Design)

- **Header** - Logo, project selector, terminal toggle
- **Sidebar** - Structure and styling (just swap nav items)
- **Layout** - Container, spacing, terminal dock
- **Terminal** - Multi-tab xterm.js, resize, fullscreen
- **NeonSelect** - Dropdown component
- **MarkdownViewer** - Markdown rendering

### KEEP & REBUILD

- **Roadmap** → **Tree Graph** (visual node graph instead of list)
- **BeadCard** → **Node Popup** (reuse for tree graph clicks)
- **Ideation Pad** → Add auto-refresh for Claude bridge

### KEEP & ENHANCE

- **Kanban** → Add drag-drop between columns

### DISCARD (replaced by bead-centric model)

- Overview, Specify, Plan, Tasks, Research, Quickstart, Data Model

### Sidebar Nav Mapping

**Old Sidebar:**

```
WORKFLOW
  Overview      ← discard
  Specify       ← discard
  Plan          ← discard
  Tasks         ← discard
  Kanban        ← keep

ARTIFACTS
  Research      ← discard
  Quickstart    ← discard
  Data Model    ← discard

SYSTEM
  Diagnostics   ← keep (dev tool)
```

**New Sidebar:**

```
VIEWS
  Tree Graph    ← NEW (default, replaces Overview)
  Kanban        ← keep
  Whiteboard    ← NEW

COLLABORATE
  Activity      ← NEW
  Assignments   ← NEW

SYSTEM
  Diagnostics   ← keep
```

Same Tokyo Night styling, same section structure, new items.

---

## Design Decisions

- **Multi-user sync**: Last-write-wins (simple, sufficient for small teams)
- **Default view**: Tree Graph
- **Plan is the artifact** - not chat history
- **UI/Layout**: Keep Tokyo Night design exactly as built

---

## Technical Notes

### Data Model Additions

```
Bead:
  + assignee: string
  + phase: enum (ideation, planning, in_progress, done)

Comment:
  + id, bead_id, author, content, created_at

Activity:
  + id, type, bead_id, user, data, created_at

Whiteboard:
  + id, bead_id (optional), data (tldraw JSON)
```

### Frontend Packages

- **D3 + SVG**: Tree graph visualization (lightweight, full Tokyo Night design control)
- **tldraw**: Whiteboard with programmatic export for Claude bridge
- **@uiw/react-codemirror**: Markdown editor (already in use)
- **React Aria**: Headless accessible primitives (free)
- **cmdk**: Command palette (later phase, nice-to-have)

### Storage Decisions

- **Whiteboard storage**: DB as tldraw JSON blob (simpler, consistent with beads)
- **Activity retention**: 90 days rolling window
- **Component library**: React Aria + existing Tokyo Night styles (no Untitled UI)

---

## Deployment & Authentication

### Hosting: Fly.io

- **Why**: Free tier, WebSocket support, auto-HTTPS, fast redeploys
- **Setup**: One-time `fly launch`, then `fly deploy` after changes
- **Public URL**: `kitty-beads.fly.dev` initially, then custom domain after MVP

### Access Model: Public Read + Authenticated Write

**Visibility:**

- Anonymous users: Can view Tree Graph, Kanban, plans (read-only)
- Authenticated users (invited friends): Can edit beads, plans, create activities
- Show "View Only" badge for anonymous users
- Hide edit buttons for unauthenticated visitors

**Data Architecture:**

```
Single SQLite DB (.beads/beads.db)
├── beads, issues, epics, tasks (existing)
├── users (id, github_id, username, email, created_at)
└── sessions (user_id, token, expires_at)
```

### GitHub OAuth

**Setup (one-time):**

1. Create OAuth App in GitHub → Developer Settings
2. Set callback URL: `https://kitty-beads.app/api/auth/callback`
3. Get Client ID + Client Secret
4. Add to Fly.io: `fly secrets set GITHUB_CLIENT_ID=xxx GITHUB_CLIENT_SECRET=yyy`

**Auth Flow:**

```
Anonymous user views kitty-beads.app
    ↓ (read-only)
Wants to edit → [Sign in with GitHub]
    ↓
Redirected to GitHub → user approves
    ↓
GitHub redirects back with code
    ↓
Server exchanges code for user profile
    ↓
Create user in DB + issue session token
    ↓
User can now edit
```

**Backend Implementation:**

- `POST /api/auth/github` → Redirect to GitHub
- `GET /api/auth/callback?code=xxx` → Exchange code, create user/session
- Middleware on `/api/*` routes: Check Authorization header for write operations
- `GET` endpoints: No auth required
- `POST/PUT/DELETE` endpoints: Auth required

**Frontend Implementation:**

- Store JWT in localStorage after auth
- Read-only mode for anonymous (hide edit buttons)
- Include token in headers for mutations
- [Sign in with GitHub] button → opens OAuth flow

### Custom Domain

**Post-MVP:**

- Buy domain (namecheap, vercel, ~$12/year)
- Point nameservers to Fly.io
- `fly certs add kitty-beads.app`
- Auto HTTPS + professional appearance for showcasing

### Deployment Checklist

- [ ] Add GitHub OAuth Client ID/Secret
- [ ] Create `users` and `sessions` tables in SQLite
- [ ] Implement auth middleware
- [ ] Add `Authorization` header to API mutations (frontend)
- [ ] Build frontend with read-only mode for anonymous users
- [ ] Test public read + auth write flow locally
- [ ] Deploy to Fly.io
- [ ] Buy + configure custom domain (post-MVP)
