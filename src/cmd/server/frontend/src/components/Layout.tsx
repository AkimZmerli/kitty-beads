import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import type { Feature } from '../types/api';

interface LayoutProps {
  features: Feature[];
  currentFeature: Feature | null;
  projectPath: string;
  onFeatureChange: (featureId: string) => void;
  onToggleTerminal?: () => void;
}

export function Layout({
  features,
  currentFeature,
  projectPath,
  onFeatureChange,
  onToggleTerminal,
}: LayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      <Header
        features={features}
        currentFeature={currentFeature?.id || null}
        projectPath={projectPath}
        onFeatureChange={onFeatureChange}
        onToggleTerminal={onToggleTerminal}
      />
      <div className="flex flex-1">
        <Sidebar currentFeature={currentFeature} />
        <main className="flex-1 p-8 bg-content-bg overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
