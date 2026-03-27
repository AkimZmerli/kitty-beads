// Package main provides terminal WebSocket handling for kitty-beads.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/steveyegge/beads/internal/terminal"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local dev
	},
}

// TerminalMessage represents messages between client and server
type TerminalMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// ResizeMessage represents a terminal resize request
type ResizeMessage struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

// TerminalSession manages a single PTY session
type TerminalSession struct {
	ptmx   *os.File
	cmd    *exec.Cmd
	ws     *websocket.Conn
	writer *terminal.BufferedTerminalWriter
	mu     sync.Mutex
	closed bool
}

// handleTerminal handles WebSocket connections for the terminal
// Supports session IDs in URL: /api/terminal/{sessionId}
func (s *Server) handleTerminal(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL path (optional, for multi-tab support)
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/terminal"), "/")
	sessionID := ""
	if len(parts) > 1 && parts[1] != "" {
		sessionID = parts[1]
	}
	_ = sessionID // Session ID available for future session management

	// Upgrade HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// Build command based on ?mode= query param
	// mode=list  → bd list --tui
	// mode=graph → bd graph --tui --all
	// mode=tui   → bd tui
	// mode=shell → interactive shell (default)
	mode := r.URL.Query().Get("mode")
	var cmdArgs []string
	switch mode {
	case "list":
		cmdArgs = []string{"bd", "list", "--tui"}
	case "graph":
		cmdArgs = []string{"bd", "graph", "--tui", "--all"}
	case "tui":
		cmdArgs = []string{"bd", "tui"}
	default:
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		cmdArgs = []string{shell}
	}

	// Create command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) // #nosec G204 -- args are from a controlled allow-list
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
	)
	cmd.Dir = s.rootDir

	// Start PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("PTY start error: %v", err)
		ws.WriteJSON(TerminalMessage{
			Type: "error",
			Data: json.RawMessage(`"Failed to start terminal"`),
		})
		return
	}

	bw := terminal.NewBufferedTerminalWriter(ws)
	defer bw.Close()

	session := &TerminalSession{
		ptmx:   ptmx,
		cmd:    cmd,
		ws:     ws,
		writer: bw,
	}
	defer session.Close()

	// Set initial size (80x24 default)
	setWinsize(ptmx, 80, 24)

	// Handle PTY output -> WebSocket
	go session.readPTY()

	// Handle WebSocket input -> PTY
	session.readWebSocket()
}

// readPTY reads from PTY and sends to WebSocket via the buffered writer.
func (ts *TerminalSession) readPTY() {
	buf := make([]byte, 4096)
	for {
		n, err := ts.ptmx.Read(buf)
		if err != nil {
			ts.Close()
			return
		}

		ts.mu.Lock()
		if ts.closed {
			ts.mu.Unlock()
			return
		}
		ts.mu.Unlock()

		if _, err := ts.writer.Write(buf[:n]); err != nil {
			ts.Close()
			return
		}
	}
}

// readWebSocket reads from WebSocket and sends to PTY
func (ts *TerminalSession) readWebSocket() {
	for {
		var msg TerminalMessage
		err := ts.ws.ReadJSON(&msg)
		if err != nil {
			ts.Close()
			return
		}

		switch msg.Type {
		case "input":
			var input string
			if err := json.Unmarshal(msg.Data, &input); err != nil {
				continue
			}
			ts.ptmx.WriteString(input)

		case "resize":
			var size ResizeMessage
			if err := json.Unmarshal(msg.Data, &size); err != nil {
				continue
			}
			setWinsize(ts.ptmx, size.Cols, size.Rows)

		case "ping":
			ts.mu.Lock()
			ts.ws.WriteJSON(TerminalMessage{Type: "pong"})
			ts.mu.Unlock()
		}
	}
}

// Close terminates the PTY session
func (ts *TerminalSession) Close() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.closed {
		return
	}
	ts.closed = true

	// Close PTY
	if ts.ptmx != nil {
		ts.ptmx.Close()
	}

	// Kill process
	if ts.cmd != nil && ts.cmd.Process != nil {
		ts.cmd.Process.Kill()
		ts.cmd.Wait()
	}

	// Send close message to client
	ts.ws.WriteJSON(TerminalMessage{Type: "exit"})
}

// setWinsize sets the terminal window size
func setWinsize(f *os.File, cols, rows uint16) {
	ws := struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}{
		Row: rows,
		Col: cols,
	}
	syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
}
