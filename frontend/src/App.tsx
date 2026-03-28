import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Layout } from "./layouts";
import { Kanban } from "./pages/Kanban";
import { Diagnostics } from "./pages/Diagnostics";
import { Roadmap } from "./pages/Roadmap";
import { TreeView } from "./pages/TreeView";
import { Editor } from "./pages/Editor";
import { IdeationPad } from "./pages/IdeationPad";
import { TerminalProvider } from "./features/terminal";
import {
  CommandPaletteProvider,
  CommandPalette,
  ShortcutsOverlay,
} from "./features/command-palette";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { WhiteboardCanvas } from "./features/whiteboard";
import { ActivityFeed } from "./features/activity";
import { useFeatures } from "./hooks/useFeatures";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000,
      retry: 1,
    },
  },
});

function AppContent() {
  const { data, isLoading } = useFeatures();
  const [currentFeatureId, setCurrentFeatureId] = useState<string | null>(null);

  const features = data?.features || [];
  const projectPath = data?.project_path || "";
  const currentFeature =
    features.find((f) => f.id === currentFeatureId) || null;

  // Auto-select first feature
  useEffect(() => {
    if (!currentFeatureId && features.length > 0) {
      setCurrentFeatureId(features[0].id);
    }
  }, [features, currentFeatureId]);

  if (isLoading && features.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-content-bg">
        <div className="text-text-muted">Loading...</div>
      </div>
    );
  }

  const featureId = currentFeature?.id || null;

  return (
    <Routes>
      <Route
        element={
          <Layout
            features={features}
            currentFeature={currentFeature}
            projectPath={projectPath}
            onFeatureChange={setCurrentFeatureId}
          />
        }
      >
        {/* Default: Roadmap (all beads view) */}
        <Route index element={<Navigate to="/roadmap" replace />} />
        <Route path="/roadmap" element={<Roadmap />} />
        <Route path="/tree" element={<TreeView />} />

        {/* Core views */}
        <Route path="/kanban" element={<Kanban featureId={featureId} />} />
        <Route path="/whiteboard" element={<WhiteboardCanvas />} />
        <Route path="/editor" element={<Editor />} />

        {/* Collaboration */}
        <Route path="/activity" element={<ActivityFeed />} />

        {/* System */}
        <Route path="/diagnostics" element={<Diagnostics />} />

        {/* Bead editing */}
        <Route path="/ideation/:issueId" element={<IdeationPad />} />
      </Route>
    </Routes>
  );
}

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <TerminalProvider>
            <CommandPaletteProvider>
              <CommandPalette />
              <ShortcutsOverlay />
              <AppContent />
            </CommandPaletteProvider>
          </TerminalProvider>
        </BrowserRouter>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}

export default App;
