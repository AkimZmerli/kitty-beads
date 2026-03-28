// Package terminal provides utilities for terminal WebSocket sessions.
package terminal

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	flushInterval = 16 * time.Millisecond // ~60fps
	flushSize     = 4 * 1024              // flush at 4KB
)

// wsMessage mirrors the server's TerminalMessage for encoding.
type wsMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// BufferedTerminalWriter batches PTY output and flushes to the WebSocket
// at ~60fps or when the buffer exceeds 4KB, whichever comes first.
// This eliminates per-read-syscall WebSocket frames and reduces flicker.
type BufferedTerminalWriter struct {
	ws     *websocket.Conn
	mu     sync.Mutex
	buf    []byte
	ticker *time.Ticker
	done   chan struct{}
}

// NewBufferedTerminalWriter creates a writer and starts the flush goroutine.
// Call Close() to stop the background goroutine.
func NewBufferedTerminalWriter(ws *websocket.Conn) *BufferedTerminalWriter {
	w := &BufferedTerminalWriter{
		ws:     ws,
		buf:    make([]byte, 0, flushSize),
		ticker: time.NewTicker(flushInterval),
		done:   make(chan struct{}),
	}
	go w.flusher()
	return w
}

// Write appends data to the buffer. Flushes immediately if threshold reached.
func (w *BufferedTerminalWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	w.buf = append(w.buf, p...)
	shouldFlush := len(w.buf) >= flushSize
	w.mu.Unlock()

	if shouldFlush {
		return len(p), w.flush()
	}
	return len(p), nil
}

// Close flushes any remaining data and stops the background goroutine.
func (w *BufferedTerminalWriter) Close() {
	w.ticker.Stop()
	close(w.done)
	_ = w.flush() // drain remaining bytes
}

// flusher is the background goroutine that fires on the 60fps ticker.
func (w *BufferedTerminalWriter) flusher() {
	for {
		select {
		case <-w.ticker.C:
			_ = w.flush()
		case <-w.done:
			return
		}
	}
}

// flush drains the buffer and sends one WebSocket text frame.
func (w *BufferedTerminalWriter) flush() error {
	w.mu.Lock()
	if len(w.buf) == 0 {
		w.mu.Unlock()
		return nil
	}
	data := w.buf
	w.buf = make([]byte, 0, flushSize)
	w.mu.Unlock()

	encoded, err := json.Marshal(string(data))
	if err != nil {
		return err
	}
	msg := wsMessage{Type: "output", Data: json.RawMessage(encoded)}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return w.ws.WriteMessage(websocket.TextMessage, payload)
}
