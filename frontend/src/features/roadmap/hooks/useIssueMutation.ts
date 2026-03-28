import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateIssue } from "../../../lib/api";
import type { Issue, RoadmapIssue } from "../../../types/api";
import type { RoadmapData } from "../../../hooks/useRoadmap";

interface MutationVariables {
  issueId: string;
  updates: Partial<Issue>;
}

/**
 * Recursively update an issue in a tree structure
 */
function updateIssueInTree(
  issues: RoadmapIssue[],
  id: string,
  updates: Partial<Issue>
): RoadmapIssue[] {
  return issues.map((issue) => {
    if (issue.id === id) {
      return { ...issue, ...updates };
    }
    if (issue.children?.length) {
      return {
        ...issue,
        children: updateIssueInTree(issue.children, id, updates),
      };
    }
    return issue;
  });
}

/**
 * Hook for updating issues with optimistic updates
 */
export function useIssueMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ issueId, updates }: MutationVariables) =>
      updateIssue(issueId, updates),

    // Optimistic update
    onMutate: async ({ issueId, updates }) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ["roadmap"] });

      // Snapshot the previous value
      const previousData = queryClient.getQueryData<RoadmapData>(["roadmap"]);

      // Optimistically update the cache
      queryClient.setQueryData<RoadmapData>(["roadmap"], (old) => {
        if (!old) return old;
        return {
          ...old,
          issues: updateIssueInTree(old.issues, issueId, updates),
        };
      });

      // Return context with the snapshot
      return { previousData };
    },

    // If mutation fails, rollback to the previous value
    onError: (_err, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(["roadmap"], context.previousData);
      }
    },

    // Always refetch after error or success
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["roadmap"] });
    },
  });
}
