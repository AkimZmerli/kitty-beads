import { useQuery } from "@tanstack/react-query";
import { getDiagnostics } from "../lib/api";

export function Diagnostics() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["diagnostics"],
    queryFn: getDiagnostics,
  });

  return (
    <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
      <h2 className="text-grassy-green text-2xl font-bold mb-6">Diagnostics</h2>

      {isLoading && (
        <p className="text-text-muted text-center py-8">Loading...</p>
      )}

      {error && (
        <p className="text-red-500 text-center py-8">
          Error loading diagnostics
        </p>
      )}

      {data && (
        <>
          <div className="bg-sidebar-bg p-5 rounded-lg mb-6">
            <h3 className="text-grassy-green font-semibold mb-4">
              Current Status
            </h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <span className="font-medium">Total Issues:</span>{" "}
                {data.status?.total_issues || 0}
              </div>
              <div>
                <span className="font-medium">Open Issues:</span>{" "}
                {data.status?.open_issues || 0}
              </div>
              <div className="col-span-2">
                <span className="font-medium">Project Path:</span>{" "}
                <code className="bg-white px-2 py-1 rounded text-sm">
                  {data.status?.project_path || ""}
                </code>
              </div>
            </div>
          </div>

          <button
            onClick={() => refetch()}
            className="px-5 py-2.5 bg-grassy-green hover:bg-grassy-green-dark text-white rounded-md font-semibold transition-colors"
          >
            Refresh Diagnostics
          </button>
        </>
      )}
    </div>
  );
}
