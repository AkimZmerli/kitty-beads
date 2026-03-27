# Bubble Tea Styling & Performance

## Performance: Bubble Tea vs Current Setup

### Current Stack
```
React → xterm.js → WebSocket → PTY → bash → bd list (prints text)
```

### With Bubble Tea
```
React → xterm.js → WebSocket → PTY → bd list --tui (Bubble Tea renders)
```

### Speed Comparison

| Aspect | Current | Bubble Tea |
|--------|---------|------------|
| Initial render | Instant (text dump) | Instant (same ANSI output) |
| Updates | Re-run command | Diff-based (only changed chars) |
| Scrolling | Shell handles it | 60fps viewport with batching |
| Input latency | Shell → command | Direct keystroke handling |
| Filtering | Re-run with `--title` flag | Live fuzzy filter, no re-query |

### Why Bubble Tea is Faster for Interactive Use
- No shell parsing overhead
- No re-executing commands for each action
- Framerate-based rendering prevents flicker
- Only redraws changed portions

---

## Styling Power: Lip Gloss

Lip Gloss provides CSS-like control over terminal styling.

### Capabilities

| Feature | Support |
|---------|---------|
| Foreground/background colors | True color (16M colors) |
| Bold, italic, underline | Yes |
| Strikethrough, blink | Yes |
| Padding & margins | Character-based |
| Borders | Rounded, thick, double, hidden, custom |
| Border colors | Per-side control |
| Width/height constraints | Yes |
| Text alignment | Left, center, right |
| Inline vs block | Yes |

---

## Tokyo Night Theme for Bubble Tea

Matches `frontend/src/index.css` color variables.

### Color Palette

```go
package styles

import "github.com/charmbracelet/lipgloss"

// Tokyo Night - Backgrounds
var (
    NightBg          = lipgloss.Color("#1a1b26")
    NightBgDark      = lipgloss.Color("#16161e")
    NightBgHighlight = lipgloss.Color("#24283b")
    NightSurface     = lipgloss.Color("#1f2335")
    NightSurfaceBright = lipgloss.Color("#292e42")
)

// Tokyo Night - Borders
var (
    NightBorder   = lipgloss.Color("#3b4261")
    NightBorderHi = lipgloss.Color("#565f89")
)

// Tokyo Night - Neon Accents
var (
    NeonCyan    = lipgloss.Color("#7dcfff")
    NeonMagenta = lipgloss.Color("#bb9af7")
    NeonPink    = lipgloss.Color("#f7768e")
    NeonGreen   = lipgloss.Color("#9ece6a")
    NeonOrange  = lipgloss.Color("#ff9e64")
    NeonYellow  = lipgloss.Color("#e0af68")
    NeonBlue    = lipgloss.Color("#7aa2f7")
)

// Tokyo Night - Text
var (
    TextBright = lipgloss.Color("#c0caf5")
    TextNormal = lipgloss.Color("#a9b1d6")
    TextMuted  = lipgloss.Color("#565f89")
    TextDark   = lipgloss.Color("#414868")
)
```

### Component Styles

```go
// Status icons
var (
    StatusOpen       = lipgloss.NewStyle().Foreground(NeonGreen)
    StatusInProgress = lipgloss.NewStyle().Foreground(NeonYellow)
    StatusBlocked    = lipgloss.NewStyle().Foreground(NeonPink)
    StatusClosed     = lipgloss.NewStyle().Foreground(TextMuted)
    StatusDeferred   = lipgloss.NewStyle().Foreground(NeonCyan)
)

// Priority badges
var (
    PriorityP0 = lipgloss.NewStyle().Foreground(NeonPink).Bold(true)
    PriorityP1 = lipgloss.NewStyle().Foreground(NeonOrange)
    PriorityP2 = lipgloss.NewStyle().Foreground(TextNormal)
    PriorityP3 = lipgloss.NewStyle().Foreground(TextMuted)
    PriorityP4 = lipgloss.NewStyle().Foreground(TextDark)
)

// Type badges
var (
    TypeEpic    = lipgloss.NewStyle().Foreground(NeonMagenta)
    TypeBug     = lipgloss.NewStyle().Foreground(NeonPink)
    TypeFeature = lipgloss.NewStyle().Foreground(NeonGreen)
    TypeTask    = lipgloss.NewStyle().Foreground(TextNormal)
)

// Panels
var (
    Panel = lipgloss.NewStyle().
        Background(NightSurface).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(NightBorder).
        Padding(1, 2)

    PanelHeader = lipgloss.NewStyle().
        Background(NightBgDark).
        Foreground(TextBright).
        Padding(0, 2).
        Bold(true)
)

// Selection
var (
    Selected = lipgloss.NewStyle().
        Background(NightBgHighlight).
        Foreground(NeonCyan)

    Cursor = lipgloss.NewStyle().
        Foreground(NeonMagenta).
        Bold(true)
)

// Help bar
var (
    HelpKey  = lipgloss.NewStyle().Foreground(NeonMagenta)
    HelpDesc = lipgloss.NewStyle().Foreground(TextMuted)
    HelpSep  = lipgloss.NewStyle().Foreground(NightBorder)
)

// Filter input
var (
    FilterPrompt = lipgloss.NewStyle().Foreground(NeonCyan)
    FilterText   = lipgloss.NewStyle().Foreground(TextBright)
    FilterCursor = lipgloss.NewStyle().Foreground(NeonMagenta)
)
```

---

## Visual Reference

```
┌─────────────────────────────────────────────────────────────┐  ← #3b4261 border
│ Issues (42)                            Filter: █            │  ← #1f2335 surface
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  > ● P0  bd-c4rq     Fix stale database        bug         │  ← #24283b selected
│    ○ P1  bd-qnz.4    Phase 4: Command...       task        │    #7dcfff cursor
│    ◐ P2  bd-2ds      Developer Hub             epic        │
│    ● P1  bd-auth     Auth timeout              bug         │
│    ❄ P3  bd-later    Deferred feature          task        │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ j/k nav  /filter  Enter view  c close  q quit              │  ← #565f89 muted
└─────────────────────────────────────────────────────────────┘
```

### Color Key

| Element | Color | Hex |
|---------|-------|-----|
| Selected row bg | NightBgHighlight | `#24283b` |
| Cursor (>) | NeonMagenta | `#bb9af7` |
| Status ○ open | NeonGreen | `#9ece6a` |
| Status ◐ progress | NeonYellow | `#e0af68` |
| Status ● blocked | NeonPink | `#f7768e` |
| Status ✓ closed | TextMuted | `#565f89` |
| Status ❄ deferred | NeonCyan | `#7dcfff` |
| Priority P0 | NeonPink | `#f7768e` |
| Priority P1 | NeonOrange | `#ff9e64` |
| Type epic | NeonMagenta | `#bb9af7` |
| Type bug | NeonPink | `#f7768e` |
| Border | NightBorder | `#3b4261` |
| Help text | TextMuted | `#565f89` |
| Help keys | NeonMagenta | `#bb9af7` |

---

## Leverage Comparison

| Capability | Current CLI | Bubble Tea |
|------------|-------------|------------|
| Color depth | 256 (xterm) | True color (16M) |
| Styling granularity | Per-line ANSI | Per-character |
| Layout | Manual spacing | Flexbox-like Join() |
| Borders | Manual box-drawing | 6+ built-in styles |
| Responsive | Fixed width | Adapts to terminal |
| Animation | None | Tick-based |
| Mouse support | Basic | Full click/scroll/drag |
| Themes | Hardcoded | Swappable structs |

---

## xterm.js Theme Sync

Update `TerminalInstance.tsx` to match Tokyo Night:

```typescript
theme: {
  background: '#1a1b26',      // NightBg
  foreground: '#c0caf5',      // TextBright
  cursor: '#bb9af7',          // NeonMagenta
  cursorAccent: '#1a1b26',
  selectionBackground: '#24283b',  // NightBgHighlight
  black: '#414868',
  red: '#f7768e',             // NeonPink
  green: '#9ece6a',           // NeonGreen
  yellow: '#e0af68',          // NeonYellow
  blue: '#7aa2f7',            // NeonBlue
  magenta: '#bb9af7',         // NeonMagenta
  cyan: '#7dcfff',            // NeonCyan
  white: '#c0caf5',
  brightBlack: '#565f89',
  brightRed: '#f7768e',
  brightGreen: '#9ece6a',
  brightYellow: '#e0af68',
  brightBlue: '#7aa2f7',
  brightMagenta: '#bb9af7',
  brightCyan: '#7dcfff',
  brightWhite: '#c0caf5',
}
```
