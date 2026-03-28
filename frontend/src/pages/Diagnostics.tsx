import { useQuery } from "@tanstack/react-query";
import { getDiagnostics } from "../lib/api";
import { RefreshCw } from "lucide-react";

export function Diagnostics() {
  const { data, isLoading, isFetching, error, refetch } = useQuery({
    queryKey: ["diagnostics"],
    queryFn: getDiagnostics,
  });

  return (
    <div className="bg-card-bg rounded-xl p-8 border border-neon-magenta">
      <h2 className="text-neon-cyan text-2xl font-bold mb-6 drop-shadow-[0_0_8px_rgba(125,207,255,0.3)]">
        Diagnostics
      </h2>

      {isLoading && (
        <p className="text-text-muted text-center py-8">Loading...</p>
      )}

      {error && (
        <p className="text-neon-pink text-center py-8">
          Error loading diagnostics
        </p>
      )}

      {data && (
        <>
          <div className="bg-night-surface p-5 rounded-lg mb-6 border border-night-border">
            <h3 className="text-neon-cyan font-semibold mb-4">
              Current Status
            </h3>
            <div className="grid grid-cols-2 gap-4 text-text-normal">
              <div>
                <span className="font-medium text-text-primary">
                  Total Issues:
                </span>{" "}
                {data.status?.total_issues || 0}
              </div>
              <div>
                <span className="font-medium text-text-primary">
                  Open Issues:
                </span>{" "}
                {data.status?.open_issues || 0}
              </div>
              <div className="col-span-2">
                <span className="font-medium text-text-primary">
                  Project Path:
                </span>{" "}
                <code className="bg-night-bg px-2 py-1 rounded text-sm text-neon-cyan border border-night-border">
                  {data.status?.project_path || ""}
                </code>
              </div>
            </div>
          </div>

          <button
            onClick={() => refetch()}
            disabled={isFetching}
            className="flex items-center gap-2 px-5 py-2.5 bg-neon-cyan hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] text-night-bg rounded-md font-semibold transition-all cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <RefreshCw
              className={`w-4 h-4 ${isFetching ? "animate-spin" : ""}`}
            />
            {isFetching ? "Refreshing..." : "Refresh Diagnostics"}
          </button>
        </>
      )}
    </div>
  );
}
