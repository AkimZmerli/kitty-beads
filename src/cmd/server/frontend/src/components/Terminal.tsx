import { useEffect, useRef, useCallback } from 'react';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';

interface TerminalMessage {
  type: 'output' | 'error' | 'exit' | 'pong';
  data?: string;
}

export function Terminal() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<XTerm | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/api/terminal`);
    wsRef.current = ws;

    ws.onopen = () => {
      xtermRef.current?.write('\x1b[32mConnected to terminal\x1b[0m\r\n');
      // Send initial size
      if (fitAddonRef.current && xtermRef.current) {
        const dims = fitAddonRef.current.proposeDimensions();
        if (dims) {
          ws.send(JSON.stringify({
            type: 'resize',
            data: { cols: dims.cols, rows: dims.rows }
          }));
        }
      }
    };

    ws.onmessage = (event) => {
      try {
        const msg: TerminalMessage = JSON.parse(event.data);
        switch (msg.type) {
          case 'output':
            if (msg.data) {
              xtermRef.current?.write(msg.data);
            }
            break;
          case 'error':
            xtermRef.current?.write(`\x1b[31mError: ${msg.data}\x1b[0m\r\n`);
            break;
          case 'exit':
            xtermRef.current?.write('\r\n\x1b[33mSession ended\x1b[0m\r\n');
            break;
        }
      } catch {
        // Handle plain text output
        xtermRef.current?.write(event.data);
      }
    };

    ws.onclose = () => {
      xtermRef.current?.write('\r\n\x1b[31mDisconnected\x1b[0m\r\n');
      // Attempt reconnect after 3 seconds
      reconnectTimeoutRef.current = window.setTimeout(() => {
        xtermRef.current?.write('\x1b[33mReconnecting...\x1b[0m\r\n');
        connect();
      }, 3000);
    };

    ws.onerror = () => {
      xtermRef.current?.write('\x1b[31mConnection error\x1b[0m\r\n');
    };
  }, []);

  useEffect(() => {
    if (!terminalRef.current) return;

    // Create terminal
    const term = new XTerm({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1e1e1e',
        foreground: '#d4d4d4',
        cursor: '#ffffff',
        cursorAccent: '#000000',
        selectionBackground: '#264f78',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#e5e5e5',
      },
    });

    xtermRef.current = term;

    // Add fit addon
    const fitAddon = new FitAddon();
    fitAddonRef.current = fitAddon;
    term.loadAddon(fitAddon);

    // Add web links addon
    const webLinksAddon = new WebLinksAddon();
    term.loadAddon(webLinksAddon);

    // Open terminal
    term.open(terminalRef.current);
    fitAddon.fit();

    // Handle input
    term.onData((data) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'input', data }));
      }
    });

    // Handle resize
    const handleResize = () => {
      fitAddon.fit();
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        const dims = fitAddon.proposeDimensions();
        if (dims) {
          wsRef.current.send(JSON.stringify({
            type: 'resize',
            data: { cols: dims.cols, rows: dims.rows }
          }));
        }
      }
    };

    window.addEventListener('resize', handleResize);

    // Connect to WebSocket
    connect();

    // Cleanup
    return () => {
      window.removeEventListener('resize', handleResize);
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      wsRef.current?.close();
      term.dispose();
    };
  }, [connect]);

  return (
    <div className="h-full w-full bg-[#1e1e1e] p-2">
      <div ref={terminalRef} className="h-full w-full" />
    </div>
  );
}
