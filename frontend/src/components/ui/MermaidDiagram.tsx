import { useEffect, useState, useId } from "react";
import mermaid from "mermaid";
import { AlertTriangle, ChevronDown, ChevronRight } from "lucide-react";
import { initializeMermaid } from "./mermaidConfig";

// Initialize mermaid once at module load
initializeMermaid();

interface MermaidDiagramProps {
  chart: string;
  className?: string;
}

type RenderState =
  | { status: "loading" }
  | { status: "success"; svg: string }
  | { status: "error"; message: string };

export function MermaidDiagram({ chart, className = "" }: MermaidDiagramProps) {
  const id = useId();
  const [state, setState] = useState<RenderState>({ status: "loading" });
  const [showSource, setShowSource] = useState(false);

  useEffect(() => {
    const renderDiagram = async () => {
      const trimmedChart = chart.trim();

      if (!trimmedChart) {
        setState({ status: "error", message: "Empty diagram content" });
        return;
      }

      setState({ status: "loading" });

      try {
        // Validate syntax first
        await mermaid.parse(trimmedChart);

        // Generate unique ID for this render
        const elementId = `mermaid-${id.replace(/:/g, "-")}-${Date.now()}`;

        // Render the diagram
        const { svg } = await mermaid.render(elementId, trimmedChart);
        setState({ status: "success", svg });
      } catch (error) {
        const message =
          error instanceof Error ? error.message : "Failed to render diagram";
        setState({ status: "error", message });
      }
    };

    renderDiagram();
  }, [chart, id]);

  if (state.status === "loading") {
    return (
      <div className={`mermaid-container ${className}`}>
        <div className="flex items-center justify-center p-8 text-[var(--color-text-muted)]">
          <div className="animate-spin h-5 w-5 border-2 border-current border-t-transparent rounded-full mr-2" />
          Rendering diagram...
        </div>
      </div>
    );
  }

  if (state.status === "error") {
    return (
      <div
        className={`mermaid-error rounded-lg border border-[var(--color-neon-pink)] bg-[var(--color-night-surface)] p-4 my-4 ${className}`}
      >
        <div className="flex items-start gap-3">
          <AlertTriangle className="h-5 w-5 text-[var(--color-neon-pink)] flex-shrink-0 mt-0.5" />
          <div className="flex-1 min-w-0">
            <p className="text-[var(--color-neon-pink)] font-medium mb-1">
              Diagram Error
            </p>
            <p className="text-[var(--color-text-muted)] text-sm break-words">
              {state.message}
            </p>

            <button
              onClick={() => setShowSource(!showSource)}
              className="flex items-center gap-1 mt-3 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text-normal)] transition-colors"
            >
              {showSource ? (
                <ChevronDown className="h-4 w-4" />
              ) : (
                <ChevronRight className="h-4 w-4" />
              )}
              {showSource ? "Hide" : "Show"} source
            </button>

            {showSource && (
              <pre className="mt-2 p-3 rounded bg-[var(--color-night-bg)] text-[var(--color-text-normal)] text-xs overflow-x-auto">
                <code>{chart}</code>
              </pre>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`mermaid-container my-4 ${className}`}
      dangerouslySetInnerHTML={{ __html: state.svg }}
    />
  );
}
