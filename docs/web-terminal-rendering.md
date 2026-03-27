# Web Terminal Rendering: Display & Performance Guide

How to make the kitty-beads Bubble Tea TUI look and perform great inside a
browser-based terminal component.

## Current Architecture

```
React (xterm.js) <── WebSocket <── PTY <── bd list --tui
                                            ↑
                                    Bubble Tea renders ANSI here
```

Bubble Tea outputs raw ANSI escape sequences. xterm.js interprets them. No
translation layer — the terminal _is_ the renderer. This means display quality
depends on two things: what ANSI we emit and how xterm.js is configured.

---

## 1. xterm.js Configuration

### GPU-Accelerated Rendering

Use the WebGL addon. This is the single biggest performance win — moves glyph
rendering off the main thread onto the GPU.

```typescript
import { Terminal } from '@xterm/xterm'
import { WebglAddon } from '@xterm/addon-webgl'

const term = new Terminal({
  fontFamily: '"JetBrains Mono", "Fira Code", "Cascadia Code", monospace',
  fontSize: 14,
  lineHeight: 1.2,
  cursorBlink: true,
  cursorStyle: 'bar',
  allowTransparency: false,  // faster when false
  scrollback: 1000,
  smoothScrollDuration: 100,
})

// Attach WebGL after open()
term.open(container)
const webgl = new WebglAddon()
webgl.onContextLoss(() => {
  webgl.dispose()
  // fall back to canvas renderer automatically
})
term.loadAddon(webgl)
```

**Key settings:**
- `allowTransparency: false` — enables fast path in WebGL renderer
- `scrollback: 1000` — keep it reasonable; large scrollback eats memory
- `smoothScrollDuration` — adds polish without perf cost

### Canvas Fallback

WebGL fails on some browsers/devices. Always have the canvas fallback:

```typescript
import { CanvasAddon } from '@xterm/addon-canvas'

try {
  term.loadAddon(new WebglAddon())
} catch {
  term.loadAddon(new CanvasAddon())
}
```

The default DOM renderer (no addon) is the slowest option. Avoid it.

### Font Loading

Monospace font must be loaded _before_ `term.open()` or glyphs measure wrong:

```typescript
await document.fonts.load('14px "JetBrains Mono"')
term.open(container)
```

If the font loads late, call `term.resize(cols, rows)` to force re-measure.

---

## 2. Tokyo Night Theme Sync

The TUI uses Tokyo Night colors via Lip Gloss. The xterm.js theme must match
exactly or you get color mismatches between styled and unstyled regions.

```typescript
const tokyoNight: ITheme = {
  background:       '#1a1b26',  // NightBg
  foreground:       '#a9b1d6',  // TextNormal
  cursor:           '#c0caf5',  // TextBright
  cursorAccent:     '#1a1b26',
  selectionBackground: '#24283b80',  // NightBgHighlight + alpha

  // ANSI 0-7 (normal)
  black:   '#16161e',  // NightBgDark
  red:     '#f7768e',  // NeonPink
  green:   '#9ece6a',  // NeonGreen
  yellow:  '#e0af68',  // NeonYellow
  blue:    '#7aa2f7',  // NeonBlue
  magenta: '#bb9af7',  // NeonMagenta
  cyan:    '#7dcfff',  // NeonCyan
  white:   '#a9b1d6',  // TextNormal

  // ANSI 8-15 (bright)
  brightBlack:   '#414868',  // TextDark
  brightRed:     '#f7768e',
  brightGreen:   '#9ece6a',
  brightYellow:  '#e0af68',
  brightBlue:    '#7aa2f7',
  brightMagenta: '#bb9af7',
  brightCyan:    '#7dcfff',
  brightWhite:   '#c0caf5',  // TextBright
}

term.options.theme = tokyoNight
```

**Important:** Lip Gloss uses true-color (24-bit) escape sequences by default,
not the 16 ANSI colors. The ANSI palette above matters for:
- The terminal's own UI (cursor, selection)
- Any non-Lip-Gloss output (log messages, shell prompts)
- Glamour markdown rendering (which uses ANSI colors)

### Forcing True Color

Ensure the PTY environment advertises true color support:

```go
cmd.Env = append(os.Environ(),
    "TERM=xterm-256color",
    "COLORTERM=truecolor",
)
```

Lip Gloss auto-detects from these env vars. Without them you get degraded
4-bit color output.

---

## 3. PTY Bridge (Server Side)

The server spawns `bd list --tui` in a PTY and pipes I/O over WebSocket.

```go
import "github.com/creack/pty"  // already a dependency

func serveTUI(ws *websocket.Conn, filter string) error {
    cmd := exec.Command("bd", "list", "--tui", "--filter", filter)
    cmd.Env = append(os.Environ(),
        "TERM=xterm-256color",
        "COLORTERM=truecolor",
    )

    ptmx, err := pty.Start(cmd)
    if err != nil {
        return err
    }
    defer ptmx.Close()

    // Set initial size from client
    pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

    // Bidirectional pipe
    go io.Copy(ptmx, ws)   // client keystrokes → PTY
    io.Copy(ws, ptmx)      // PTY output → client
    return cmd.Wait()
}
```

### Resize Handling

The client must send resize events. Use a simple binary protocol or JSON
messages on a control channel:

```typescript
// Client sends resize
ws.send(JSON.stringify({ type: 'resize', cols, rows }))
```

```go
// Server handles resize
pty.Setsize(ptmx, &pty.Winsize{
    Rows: uint16(msg.Rows),
    Cols: uint16(msg.Cols),
})
```

Bubble Tea picks up the new size via `SIGWINCH` automatically.

### Fit Addon

Auto-size the terminal to its container:

```typescript
import { FitAddon } from '@xterm/addon-fit'

const fit = new FitAddon()
term.loadAddon(fit)
fit.fit()

const resizeObserver = new ResizeObserver(() => {
  fit.fit()
  ws.send(JSON.stringify({
    type: 'resize',
    cols: term.cols,
    rows: term.rows,
  }))
})
resizeObserver.observe(container)
```

---

## 4. Rendering Performance

### The 60fps Target

Bubble Tea batches output by default — it collects `View()` calls and writes
once per frame. On the xterm.js side, the WebGL renderer also batches at
`requestAnimationFrame`. The pipeline is already frame-aligned on both ends.

**Where you lose frames:**

| Bottleneck | Cause | Fix |
|---|---|---|
| Network jitter | WebSocket latency spikes | Buffer + batch writes on client |
| Large repaints | Full-screen redraws on every keystroke | Bubble Tea already handles this via alt-screen |
| Slow View() | Complex string building in Go | Profile with `pprof`; pre-compute styles |
| Font shaping | Complex ligatures on every frame | Disable ligatures or use simpler font |
| Scrollback | Huge scrollback buffer | Cap at 1000 lines |

### Write Batching on the Client

Network writes arrive in chunks. Batch them before feeding to xterm.js:

```typescript
let pending = ''
let rafId: number | null = null

ws.onmessage = (event) => {
  pending += event.data
  if (!rafId) {
    rafId = requestAnimationFrame(() => {
      term.write(pending)
      pending = ''
      rafId = null
    })
  }
}
```

This coalesces multiple WebSocket messages into a single `term.write()` per
frame, which xterm.js can then render in one GPU pass.

### Avoid Reflow Thrashing

Never read `term.cols`/`term.rows` and then immediately resize in the same
frame. The fit addon handles this correctly, but manual resize code can trigger
layout thrashing.

---

## 5. Unicode & Glyph Rendering

### Status Icons

The TUI uses Unicode symbols for issue status:

| Status | Glyph | Codepoint |
|---|---|---|
| Open | ○ | U+25CB |
| In Progress | ◐ | U+25D0 |
| Blocked | ● | U+25CF |
| Closed | ✓ | U+2713 |
| Deferred | ❄ | U+2744 |
| Pinned | 📌 | U+1F4CC |

**Problem:** Not all monospace fonts have these glyphs. Missing glyphs render
as tofu (□) or fall back to a proportional font, breaking alignment.

**Fix:** Use the Unicode graphemes addon + a Nerd Font:

```typescript
import { Unicode11Addon } from '@xterm/addon-unicode11'

term.loadAddon(new Unicode11Addon())
term.unicode.activeVersion = '11'
```

And use a font stack that covers the icons:

```typescript
fontFamily: '"JetBrainsMono Nerd Font", "JetBrains Mono", monospace'
```

### Wide Characters

Emoji like 📌 are double-width. Lip Gloss accounts for this when computing
column widths, and xterm.js handles it with the Unicode addon. Without the
addon, cursor positioning breaks after any wide character.

---

## 6. Input Handling

### Keyboard

xterm.js captures keystrokes and sends escape sequences. Most work out of the
box, but watch for:

- **Browser shortcuts hijacking keys:** Ctrl+W (close tab), Ctrl+T (new tab),
  Ctrl+N (new window) can't be captured. The TUI avoids these — it uses `q`,
  `j/k`, `/`, `Enter`, `Esc`.
- **Meta/Alt key:** Behaves differently across OS. Keep bindings simple.
- **`/` for search:** Works fine — xterm.js sends it as a literal character.

### Mouse

Bubble Tea is initialized with `tea.WithMouseCellMotion()`. xterm.js sends
mouse events in SGR format by default, which Bubble Tea handles. No extra
config needed.

If mouse events aren't working, check that the PTY has mouse reporting enabled:

```typescript
// xterm.js sends mouse in SGR mode by default — should just work
// If not, force it:
term.write('\x1b[?1006h')  // enable SGR mouse mode
```

---

## 7. Session Management

Each browser tab = one PTY session running `bd list --tui`. Consider:

### Idle Timeout

Kill idle sessions to free server resources:

```go
timer := time.NewTimer(10 * time.Minute)
go func() {
    <-timer.C
    cmd.Process.Signal(syscall.SIGTERM)
}()

// Reset timer on any input
go func() {
    buf := make([]byte, 1024)
    for {
        n, err := ws.Read(buf)
        if err != nil { return }
        timer.Reset(10 * time.Minute)
        ptmx.Write(buf[:n])
    }
}()
```

### Concurrent Users

The current store is single-writer. Multiple web sessions writing to the same
store will conflict. Options:
- **Read-only web sessions** — simplest, viewers can browse but not close/edit
- **Daemon mode** — all web sessions connect to the daemon, which serializes writes
- **Per-session locks** — pessimistic locking on write operations

Daemon mode is the natural fit since it already handles serialized access.

---

## 8. Container Sizing & Layout

The terminal component lives inside a larger web layout. Get the sizing right:

```css
.terminal-container {
  /* Fill available space */
  width: 100%;
  height: 100%;

  /* Prevent content from pushing layout */
  overflow: hidden;

  /* Terminal needs explicit dimensions — can't be intrinsic */
  min-height: 300px;

  /* Dark background to prevent flash before xterm.js loads */
  background: #1a1b26;  /* NightBg */
}
```

### Responsive Behavior

When the container resizes (sidebar toggle, window resize, panel drag):

1. `ResizeObserver` fires
2. `FitAddon.fit()` recalculates cols/rows
3. Client sends resize message over WebSocket
4. Server calls `pty.Setsize()`
5. Bubble Tea receives `SIGWINCH`, calls `View()` with new dimensions
6. Fresh ANSI output flows back through the pipeline

This whole cycle takes 1-2 frames. It feels instant.

### Minimum Viable Size

The TUI needs minimum dimensions to render properly:
- **Width:** ~60 cols (list delegate truncates titles below this)
- **Height:** ~10 rows (list + status bar + margins)

Below these thresholds, consider showing a "terminal too small" message instead
of a broken layout.

---

## 9. Connection Lifecycle

```
Client                          Server
  │                               │
  ├── WebSocket connect ─────────►│
  │                               ├── Spawn PTY + bd list --tui
  │◄── Initial render ───────────┤
  │                               │
  ├── Keystrokes ────────────────►│── PTY stdin
  │◄── ANSI output ──────────────┤── PTY stdout
  │                               │
  ├── Resize events ─────────────►│── pty.Setsize()
  │                               │
  │  (idle timeout or tab close)  │
  ├── WebSocket close ───────────►│
  │                               ├── SIGTERM → bd process
  │                               ├── Close PTY
  │                               └── Cleanup session
```

### Reconnection

If the WebSocket drops, the PTY session is gone — Bubble Tea doesn't support
detach/reattach. Options:
- **Simple:** Show "disconnected" overlay, prompt to reload
- **Better:** Auto-reconnect, spawn a new session, restore the user's filter
  state via URL params (`?filter=open+P1`)

---

## 10. Checklist

Before shipping the web terminal component:

- [ ] WebGL addon loaded with canvas fallback
- [ ] Tokyo Night theme applied to xterm.js (matches Lip Gloss palette)
- [ ] Font loaded before `term.open()`
- [ ] Fit addon with ResizeObserver
- [ ] Write batching via requestAnimationFrame
- [ ] Unicode11 addon for status icons
- [ ] PTY env includes `TERM=xterm-256color` and `COLORTERM=truecolor`
- [ ] Resize messages sent over WebSocket
- [ ] Idle session timeout on server
- [ ] Minimum size threshold with graceful fallback
- [ ] Reconnection handling (at minimum: "disconnected" message)
- [ ] Mouse events working (SGR mode)
- [ ] No browser shortcut conflicts with TUI keybindings
