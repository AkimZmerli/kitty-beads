import { useQuery } from "@tanstack/react-query";
import { getConstitution } from "../lib/api";
import { MarkdownViewer } from "../components/MarkdownViewer";

export function Constitution() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["constitution"],
    queryFn: getConstitution,
  });

  return (
    <div className="bg-white rounded-xl p-8 border-2 border-sunny-yellow-border">
      <h2 className="text-grassy-green text-2xl font-bold mb-6">
        Project Constitution
      </h2>

      {isLoading && (
        <p className="text-text-muted text-center py-8">Loading...</p>
      )}

      {error && (
        <p className="text-text-muted text-center py-8">
          No constitution file found (CLAUDE.md or README.md).
        </p>
      )}

      {data && <MarkdownViewer content={data} />}
    </div>
  );
}
