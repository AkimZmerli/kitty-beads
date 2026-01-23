import { useQuery } from "@tanstack/react-query";
import { getKanban } from "../lib/api";

export function useKanban(featureId: string | null) {
  return useQuery({
    queryKey: ["kanban", featureId],
    queryFn: () => getKanban(featureId!),
    enabled: !!featureId,
    refetchInterval: 5000,
  });
}
