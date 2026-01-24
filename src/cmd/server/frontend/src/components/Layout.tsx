import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { TerminalPanel } from './terminal/TerminalPanel';
import { useTerminal } from '../context/TerminalContext';
import type { Feature } from '../types/api';

interface LayoutProps {
  features: Feature[];
  currentFeature: Feature | null;
  projectPath: string;
  onFeatureChange: (featureId: string) => void;
}

export function Layout({
  features,
  currentFeature,
  projectPath,
  onFeatureChange,
}: LayoutProps) {
  const { isOpen, panelHeight, isFullscreen } = useTerminal();

  // Calculate main content padding to account for terminal panel
  const mainPaddingBottom = isOpen && !isFullscreen ? panelHeight : 0;

  return (
    <div className="min-h-screen flex flex-col">
      <Header
        features={features}
        currentFeature={currentFeature?.id || null}
        projectPath={projectPath}
        onFeatureChange={onFeatureChange}
      />
      <div className="flex flex-1">
        <Sidebar currentFeature={currentFeature} />
        <main
          className="flex-1 p-8 bg-content-bg overflow-auto transition-[padding-bottom] duration-200"
          style={{ paddingBottom: mainPaddingBottom ? `${mainPaddingBottom + 32}px` : undefined }}
        >
          <Outlet />
        </main>
      </div>
      <TerminalPanel />
    </div>
  );
}
