import { useTerminal } from '../context/TerminalContext';
import type { Feature } from '../types/api';
import { NeonSelect } from './NeonSelect';

interface HeaderProps {
  features: Feature[];
  currentFeature: string | null;
  projectPath: string;
  onFeatureChange: (featureId: string) => void;
}

// Tokyo Night themed cat logo SVG
function CatLogo() {
  return (
    <svg className="w-12 h-12" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      {/* Ears */}
      <path d="M20 45 L30 15 L40 40" stroke="#565f89" strokeWidth="3" fill="#24283b" strokeLinecap="round" strokeLinejoin="round"/>
      <path d="M60 40 L70 15 L80 45" stroke="#565f89" strokeWidth="3" fill="#24283b" strokeLinecap="round" strokeLinejoin="round"/>
      {/* Inner ears - neon pink */}
      <path d="M25 42 L30 22 L35 40" stroke="#f7768e" strokeWidth="2" fill="#f7768e33" strokeLinecap="round" strokeLinejoin="round"/>
      <path d="M65 40 L70 22 L75 42" stroke="#f7768e" strokeWidth="2" fill="#f7768e33" strokeLinecap="round" strokeLinejoin="round"/>
      {/* Head */}
      <ellipse cx="50" cy="58" rx="32" ry="28" stroke="#565f89" strokeWidth="3" fill="#24283b"/>
      {/* Eyes - neon yellow */}
      <ellipse cx="38" cy="52" rx="6" ry="7" stroke="#565f89" strokeWidth="2" fill="#e0af68"/>
      <ellipse cx="62" cy="52" rx="6" ry="7" stroke="#565f89" strokeWidth="2" fill="#e0af68"/>
      <ellipse cx="38" cy="53" rx="3" ry="4" fill="#1a1b26"/>
      <ellipse cx="62" cy="53" rx="3" ry="4" fill="#1a1b26"/>
      <circle cx="36" cy="51" r="1.5" fill="#7dcfff"/>
      <circle cx="60" cy="51" r="1.5" fill="#7dcfff"/>
      {/* Nose - neon pink */}
      <path d="M47 62 L50 66 L53 62 Z" stroke="#565f89" strokeWidth="1.5" fill="#f7768e"/>
      {/* Mouth */}
      <path d="M50 66 L50 70" stroke="#565f89" strokeWidth="2" strokeLinecap="round"/>
      <path d="M50 70 Q44 76 38 72" stroke="#565f89" strokeWidth="2" fill="none" strokeLinecap="round"/>
      <path d="M50 70 Q56 76 62 72" stroke="#565f89" strokeWidth="2" fill="none" strokeLinecap="round"/>
      {/* Whiskers */}
      <path d="M30 65 L12 60" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M30 68 L12 68" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M30 71 L12 76" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 65 L88 60" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 68 L88 68" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 71 L88 76" stroke="#565f89" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
  );
}

export function Header({ features, currentFeature, projectPath, onFeatureChange }: HeaderProps) {
  const { togglePanel, isOpen } = useTerminal();
  const lastUpdate = new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });

  // Get workspace name (last folder in path, excluding kitty-beads)
  const workspaceName = projectPath?.split('/').filter(Boolean).pop() || 'workspace';

  // Convert features to options
  const featureOptions = features.map(f => ({
    value: f.id,
    label: `${f.id}: ${f.name}`
  }));

  return (
    <header className="fixed top-0 left-0 right-0 z-30 flex items-center px-8 h-18 bg-night-bg border-b border-night-border">
      {/* Logo Section - never shrinks */}
      <div className="flex items-center gap-4 shrink-0">
        <CatLogo />
        <div>
          <h1 className="text-2xl font-bold text-neon-cyan whitespace-nowrap drop-shadow-[0_0_10px_rgba(125,207,255,0.3)]">Kitty-Beads</h1>
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
        <label className="font-medium text-neon-magenta whitespace-nowrap">クエスト:</label>
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
              ? 'border-neon-cyan shadow-[0_0_10px_rgba(125,207,255,0.3)]'
              : 'border-terminal-border hover:border-neon-cyan hover:shadow-[0_0_10px_rgba(125,207,255,0.2)]'
          }`}
          title="Toggle Terminal (Ctrl+`)"
        >
          <span className="text-sm font-mono text-neon-cyan">&gt;_</span>
          <span>Terminal</span>
        </button>

        {/* Last Update */}
        <div className="text-sm text-text-muted">
          {lastUpdate} | v0.0.1
        </div>
      </div>
    </header>
  );
}
