import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getIssue, updateIssue } from "../lib/api";
import type { Issue } from "../types/api";

export interface UseIdeationOptions {
  issueId: string;
  onSaveSuccess?: (issue: Issue) => void;
  onSaveError?: (error: Error) => void;
}

/**
 * Hook for the Ideation Pad - fetches issue and provides save functionality
 */
export function useIdeation({ issueId, onSaveSuccess, onSaveError }: UseIdeationOptions) {
  const queryClient = useQueryClient();

  const query = useQuery<Issue>({
    queryKey: ["issue", issueId],
    queryFn: () => getIssue(issueId),
    enabled: !!issueId,
  });

  const mutation = useMutation({
    mutationFn: (updates: Partial<Issue>) => updateIssue(issueId, updates),
    onSuccess: (updatedIssue) => {
      // Update the cache
      queryClient.setQueryData(["issue", issueId], updatedIssue);
      // Invalidate roadmap to reflect changes
      queryClient.invalidateQueries({ queryKey: ["roadmap"] });
      onSaveSuccess?.(updatedIssue);
    },
    onError: (error: Error) => {
      onSaveError?.(error);
    },
  });

  return {
    issue: query.data,
    isLoading: query.isLoading,
    error: query.error,
    save: mutation.mutate,
    isSaving: mutation.isPending,
    saveError: mutation.error,
  };
}
