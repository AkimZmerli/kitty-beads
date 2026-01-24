import { Link } from 'react-router-dom';
import { useTerminal } from '../context/TerminalContext';
import type { Feature } from '../types/api';

interface HeaderProps {
  features: Feature[];
  currentFeature: string | null;
  projectPath: string;
  onFeatureChange: (featureId: string) => void;
}

// Sketch-style cat icon SVG
function CatLogo() {
  return (
    <svg className="w-12 h-12" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      {/* Ears */}
      <path d="M20 45 L30 15 L40 40" stroke="#374151" strokeWidth="3" fill="#f9fafb" strokeLinecap="round" strokeLinejoin="round"/>
      <path d="M60 40 L70 15 L80 45" stroke="#374151" strokeWidth="3" fill="#f9fafb" strokeLinecap="round" strokeLinejoin="round"/>
      {/* Inner ears */}
      <path d="M25 42 L30 22 L35 40" stroke="#f472b6" strokeWidth="2" fill="#fce7f3" strokeLinecap="round" strokeLinejoin="round"/>
      <path d="M65 40 L70 22 L75 42" stroke="#f472b6" strokeWidth="2" fill="#fce7f3" strokeLinecap="round" strokeLinejoin="round"/>
      {/* Head */}
      <ellipse cx="50" cy="58" rx="32" ry="28" stroke="#374151" strokeWidth="3" fill="#f9fafb"/>
      {/* Eyes */}
      <ellipse cx="38" cy="52" rx="6" ry="7" stroke="#374151" strokeWidth="2" fill="#fbbf24"/>
      <ellipse cx="62" cy="52" rx="6" ry="7" stroke="#374151" strokeWidth="2" fill="#fbbf24"/>
      <ellipse cx="38" cy="53" rx="3" ry="4" fill="#111827"/>
      <ellipse cx="62" cy="53" rx="3" ry="4" fill="#111827"/>
      <circle cx="36" cy="51" r="1.5" fill="white"/>
      <circle cx="60" cy="51" r="1.5" fill="white"/>
      {/* Nose */}
      <path d="M47 62 L50 66 L53 62 Z" stroke="#374151" strokeWidth="1.5" fill="#f472b6"/>
      {/* Mouth */}
      <path d="M50 66 L50 70" stroke="#374151" strokeWidth="2" strokeLinecap="round"/>
      <path d="M50 70 Q44 76 38 72" stroke="#374151" strokeWidth="2" fill="none" strokeLinecap="round"/>
      <path d="M50 70 Q56 76 62 72" stroke="#374151" strokeWidth="2" fill="none" strokeLinecap="round"/>
      {/* Whiskers */}
      <path d="M30 65 L12 60" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M30 68 L12 68" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M30 71 L12 76" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 65 L88 60" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 68 L88 68" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M70 71 L88 76" stroke="#374151" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
  );
}

export function Header({ features, currentFeature, projectPath, onFeatureChange }: HeaderProps) {
  const { togglePanel, isOpen } = useTerminal();
  const lastUpdate = new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });

  return (
    <header className="flex items-center justify-between px-8 py-4 bg-white border-b border-border">
      <div className="flex items-center gap-8">
        {/* Logo Section */}
        <div className="flex items-center gap-4">
          <CatLogo />
          <div>
            <h1 className="text-2xl font-bold text-grassy-green">Kitty-Beads</h1>
            <div className="text-sm font-mono bg-gray-100 px-3 py-1 rounded text-text-secondary">
              {projectPath || '/loading...'}
            </div>
          </div>
        </div>

        {/* Feature Selector */}
        <div className="flex items-center gap-3">
          <label className="font-medium text-text-muted">Feature:</label>
          <select
            value={currentFeature || ''}
            onChange={(e) => onFeatureChange(e.target.value)}
            className="px-4 py-2 border border-border-light rounded-md min-w-[250px]"
          >
            {features.length === 0 ? (
              <option value="">Loading...</option>
            ) : (
              features.map((f) => (
                <option key={f.id} value={f.id}>
                  {f.id}: {f.name}
                </option>
              ))
            )}
          </select>
        </div>
      </div>

      <div className="flex items-center gap-5">
        {/* Terminal Button */}
        <button
          onClick={togglePanel}
          className={`flex items-center gap-2 px-4 py-2 bg-terminal-bg text-terminal-text border rounded-md transition-colors ${
            isOpen ? 'border-grassy-green' : 'border-terminal-border hover:border-grassy-green'
          }`}
          title="Toggle Terminal (Ctrl+`)"
        >
          <span className="text-sm font-mono">&gt;_</span>
          <span>Terminal</span>
        </button>

        {/* Constitution Button */}
        <Link
          to="/constitution"
          className="flex items-center gap-2 px-5 py-2.5 bg-grassy-green text-white rounded-md font-semibold hover:bg-grassy-green-dark transition-colors"
        >
          <span>📜</span>
          <span>Constitution</span>
        </Link>

        {/* Last Update */}
        <div className="text-sm text-text-muted">
          {lastUpdate} | v0.0.1
        </div>
      </div>
    </header>
  );
}
