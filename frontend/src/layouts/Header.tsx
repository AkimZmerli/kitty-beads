import { useTerminal } from "../features/terminal";
import type { Feature } from "../types/api";
import { NeonSelect } from "../components/ui/NeonSelect";
import { CatLogo } from "../components/ui/CatLogo";

interface HeaderProps {
  features: Feature[];
  currentFeature: string | null;
  projectPath: string;
  onFeatureChange: (featureId: string) => void;
}

export function Header({
  features,
  currentFeature,
  projectPath,
  onFeatureChange,
}: HeaderProps) {
  const { togglePanel, isOpen } = useTerminal();
  const lastUpdate = new Date().toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });

  // Get workspace name (last folder in path, excluding kitty-beads)
  const workspaceName =
    projectPath?.split("/").filter(Boolean).pop() || "workspace";

  // Convert features to options
  const featureOptions = features.map((f) => ({
    value: f.id,
    label: `${f.id}: ${f.name}`,
  }));

  return (
    <header className="fixed top-0 left-0 right-0 z-30 flex items-center pl-4 pr-8 h-18 bg-night-bg border-b border-night-border">
      {/* Logo Section - never shrinks */}
      <div className="flex items-center gap-4 shrink-0">
        <CatLogo />
        <div>
          <h1 className="text-2xl font-bold text-neon-cyan whitespace-nowrap drop-shadow-[0_0_10px_rgba(125,207,255,0.3)]">
            Kitty-Beads
          </h1>
          <div className="flex items-center gap-1.5 text-sm font-mono whitespace-nowrap">
            <span className="text-neon-green">$</span>
            <span className="text-text-muted">pwd</span>
            <span className="text-neon-magenta">&gt;</span>
            <span className="text-text-secondary">{workspaceName}</span>
          </div>
        </div>
      </div>

      {/* Feature Selector - absolutely positioned, locked location */}
      <div className="absolute left-80 flex items-center gap-3">
        <label className="font-medium text-neon-magenta whitespace-nowrap">
          クエスト:
        </label>
        <NeonSelect
          options={featureOptions}
          value={currentFeature}
          onChange={onFeatureChange}
          placeholder="Select feature..."
        />
      </div>

      {/* Right Section - pushed to end */}
      <div className="flex items-center gap-5 ml-auto shrink-0">
        {/* Terminal Button */}
        <button
          onClick={togglePanel}
          className={`flex items-center gap-2 px-4 py-2 bg-terminal-bg text-terminal-text border rounded-md transition-all ${
            isOpen
              ? "border-neon-cyan shadow-[0_0_10px_rgba(125,207,255,0.3)]"
              : "border-terminal-border hover:border-neon-cyan hover:shadow-[0_0_10px_rgba(125,207,255,0.2)]"
          }`}
          title="Toggle Terminal (Ctrl+`)"
        >
          <span className="text-sm font-mono text-neon-cyan">&gt;_</span>
          <span>Terminal</span>
        </button>

        {/* Last Update */}
        <div className="text-sm text-text-muted">{lastUpdate} | v0.0.1</div>
      </div>
    </header>
  );
}
