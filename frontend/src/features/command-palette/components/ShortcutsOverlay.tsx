import { useCommandPalette } from "../context";

interface ShortcutRow {
  keys: string[];
  description: string;
}

const shortcuts: ShortcutRow[] = [
  { keys: ["Cmd+K", "Ctrl+K"], description: "Command palette" },
  { keys: ["Ctrl+`"], description: "Toggle terminal" },
  { keys: ["?"], description: "Show keyboard shortcuts" },
  { keys: ["1–9"], description: "Navigate tabs (future)" },
];

export function ShortcutsOverlay() {
  const { showShortcuts, closeShortcuts } = useCommandPalette();

  if (!showShortcuts) return null;

  return (
    <div className="fixed inset-0 z-[110]">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={closeShortcuts}
      />

      {/* Modal */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md">
        <div
          className="bg-card-bg border border-night-border rounded-xl shadow-2xl overflow-hidden"
          onKeyDown={(e) => {
            if (e.key === "Escape") {
              closeShortcuts();
            }
          }}
        >
          {/* Header */}
          <div className="flex items-center justify-between px-5 py-4 border-b border-night-border">
            <h2 className="text-text-primary font-semibold text-sm tracking-wide uppercase">
              Keyboard Shortcuts
            </h2>
            <button
              className="text-text-muted hover:text-text-primary transition-colors"
              onClick={closeShortcuts}
              aria-label="Close shortcuts overlay"
            >
              <span className="text-lg leading-none">&times;</span>
            </button>
          </div>

          {/* Shortcuts table */}
          <div className="p-5">
            <table className="w-full">
              <tbody className="divide-y divide-night-border/50">
                {shortcuts.map((row) => (
                  <tr key={row.description} className="group">
                    <td className="py-3 pr-6">
                      <div className="flex items-center gap-1.5 flex-wrap">
                        {row.keys.map((key, i) => (
                          <span key={key} className="flex items-center gap-1.5">
                            <kbd className="inline-block px-2 py-0.5 rounded bg-night-bg border border-night-border text-neon-cyan text-xs font-mono">
                              {key}
                            </kbd>
                            {i < row.keys.length - 1 && (
                              <span className="text-text-muted text-xs">/</span>
                            )}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="py-3 text-text-secondary text-sm">
                      {row.description}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Footer */}
          <div className="px-5 py-3 border-t border-night-border text-xs text-text-muted">
            Press <kbd className="bg-night-bg px-1.5 py-0.5 rounded border border-night-border mx-1">Esc</kbd> or click outside to close
          </div>
        </div>
      </div>
    </div>
  );
}
