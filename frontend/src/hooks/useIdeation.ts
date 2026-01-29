import { useState, useEffect, useRef, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getIssue, updateIssue } from "../lib/api";
import type { Issue } from "../types/api";

export interface UseIdeationOptions {
  issueId: string;
  onSaveSuccess?: (issue: Issue) => void;
  onSaveError?: (error: Error) => void;
  /** Polling interval for conflict detection (ms). Default: 5000 */
  pollInterval?: number;
}

/**
 * Hook for the Ideation Pad - fetches issue and provides save functionality
 * with conflict detection via polling
 */
export function useIdeation({
  issueId,
  onSaveSuccess,
  onSaveError,
  pollInterval = 5000,
}: UseIdeationOptions) {
  const queryClient = useQueryClient();

  // Track the last known version of content when editing began
  const [lastKnownContent, setLastKnownContent] = useState<string | null>(null);
  const [hasExternalChanges, setHasExternalChanges] = useState(false);
  const isEditing = useRef(false);

  const query = useQuery<Issue>({
    queryKey: ["issue", issueId],
    queryFn: () => getIssue(issueId),
    enabled: !!issueId,
    refetchInterval: pollInterval,
    refetchIntervalInBackground: false,
  });

  // Detect external changes
  useEffect(() => {
    if (query.data && lastKnownContent !== null && isEditing.current) {
      const serverContent = query.data.design || query.data.description || "";
      if (serverContent !== lastKnownContent) {
        setHasExternalChanges(true);
      }
    }
  }, [query.data, lastKnownContent]);

  // Mark editing as started when content is first captured
  const startEditing = useCallback((currentContent: string) => {
    if (!isEditing.current) {
      isEditing.current = true;
      setLastKnownContent(currentContent);
      setHasExternalChanges(false);
    }
  }, []);

  // Refresh to get latest changes
  const refreshContent = useCallback(() => {
    setHasExternalChanges(false);
    setLastKnownContent(null);
    isEditing.current = false;
    queryClient.invalidateQueries({ queryKey: ["issue", issueId] });
  }, [queryClient, issueId]);

  // Dismiss the conflict warning without refreshing
  const dismissConflict = useCallback(() => {
    if (query.data) {
      const serverContent = query.data.design || query.data.description || "";
      setLastKnownContent(serverContent);
      setHasExternalChanges(false);
    }
  }, [query.data]);

  const mutation = useMutation({
    mutationFn: (updates: Partial<Issue>) => updateIssue(issueId, updates),
    onSuccess: (updatedIssue) => {
      // Update the cache
      queryClient.setQueryData(["issue", issueId], updatedIssue);
      // Reset editing state
      const newContent = updatedIssue.design || updatedIssue.description || "";
      setLastKnownContent(newContent);
      setHasExternalChanges(false);
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
    // Conflict detection
    hasExternalChanges,
    startEditing,
    refreshContent,
    dismissConflict,
  };
}
