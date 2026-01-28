import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Layout } from "./layouts";
import { Kanban } from "./pages/Kanban";
import { Diagnostics } from "./pages/Diagnostics";
import { Roadmap } from "./pages/Roadmap";
import { IdeationPad } from "./pages/IdeationPad";
import { TerminalProvider } from "./features/terminal";
import { useFeatures } from "./hooks/useFeatures";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000,
      retry: 1,
    },
  },
});

// Placeholder for Whiteboard (tldraw integration - future)
function Whiteboard() {
  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <h2 className="text-neon-cyan text-2xl font-bold mb-4">Whiteboard</h2>
      <p className="text-text-muted">
        Freeform brainstorming canvas coming soon.
      </p>
      <p className="text-text-secondary mt-4 text-sm">
        Will integrate tldraw for visual thinking.
      </p>
    </div>
  );
}

// Placeholder for Activity Feed (future)
function Activity() {
  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <h2 className="text-neon-cyan text-2xl font-bold mb-4">Activity Feed</h2>
      <p className="text-text-muted">Collaboration activity log coming soon.</p>
      <p className="text-text-secondary mt-4 text-sm">
        Plan edits, status changes, comments, assignments.
      </p>
    </div>
  );
}

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
        {/* Default: Tree Graph (currently Roadmap, will become D3 tree) */}
        <Route index element={<Navigate to="/tree" replace />} />
        <Route path="/tree" element={<Roadmap />} />

        {/* Core views */}
        <Route path="/kanban" element={<Kanban featureId={featureId} />} />
        <Route path="/whiteboard" element={<Whiteboard />} />

        {/* Collaboration */}
        <Route path="/activity" element={<Activity />} />

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
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <TerminalProvider>
          <AppContent />
        </TerminalProvider>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
