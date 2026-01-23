import { Terminal } from '../components/Terminal';

export function TerminalPage() {
  return (
    <div className="h-[calc(100vh-4rem)] flex flex-col">
      <div className="flex items-center justify-between px-6 py-3 bg-gray-100 border-b border-border">
        <h1 className="text-lg font-semibold text-text-primary">Terminal</h1>
        <div className="text-sm text-text-muted">
          Press <kbd className="px-1.5 py-0.5 bg-gray-200 rounded text-xs font-mono">Ctrl+C</kbd> to interrupt
        </div>
      </div>
      <div className="flex-1 min-h-0">
        <Terminal />
      </div>
    </div>
  );
}
