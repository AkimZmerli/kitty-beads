import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Layout } from "./components/Layout";
import { Overview } from "./pages/Overview";
import { Specify } from "./pages/Specify";
import { Plan } from "./pages/Plan";
import { Tasks } from "./pages/Tasks";
import { Kanban } from "./pages/Kanban";
import { Research } from "./pages/Research";
import { Quickstart } from "./pages/Quickstart";
import { DataModel } from "./pages/DataModel";
import { Constitution } from "./pages/Constitution";
import { Diagnostics } from "./pages/Diagnostics";
import { TerminalPage } from "./pages/TerminalPage";
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
        <Route index element={<Overview feature={currentFeature} />} />
        <Route path="/specify" element={<Specify featureId={featureId} />} />
        <Route path="/plan" element={<Plan featureId={featureId} />} />
        <Route path="/tasks" element={<Tasks featureId={featureId} />} />
        <Route path="/kanban" element={<Kanban featureId={featureId} />} />
        <Route path="/research" element={<Research featureId={featureId} />} />
        <Route
          path="/quickstart"
          element={<Quickstart featureId={featureId} />}
        />
        <Route
          path="/data-model"
          element={<DataModel featureId={featureId} />}
        />
        <Route path="/constitution" element={<Constitution />} />
        <Route path="/diagnostics" element={<Diagnostics />} />
        <Route path="/terminal" element={<TerminalPage />} />
      </Route>
    </Routes>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AppContent />
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
