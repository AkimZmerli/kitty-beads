import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { TerminalPanel } from './terminal/TerminalPanel';
import { useTerminal } from '../context/TerminalContext';
import type { Feature } from '../types/api';

// Layout constants
const HEADER_HEIGHT = 73; // px - must match Header
const SIDEBAR_WIDTH = 224; // px - w-56 = 14rem = 224px

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
    <>
      {/* Fixed Header */}
      <Header
        features={features}
        currentFeature={currentFeature?.id || null}
        projectPath={projectPath}
        onFeatureChange={onFeatureChange}
      />
      {/* Fixed Sidebar */}
      <Sidebar currentFeature={currentFeature} />
      {/* Main content area - flows normally with margin offsets */}
      <main
        className="min-h-screen bg-content-bg p-8"
        style={{
          marginTop: `${HEADER_HEIGHT}px`,
          marginLeft: `${SIDEBAR_WIDTH}px`,
          paddingBottom: mainPaddingBottom ? `${mainPaddingBottom + 32}px` : undefined,
        }}
      >
        <Outlet />
      </main>
      <TerminalPanel />
    </>
  );
}
