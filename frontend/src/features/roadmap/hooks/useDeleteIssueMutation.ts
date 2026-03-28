import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteIssue } from "../../../lib/api";
import type { RoadmapIssue } from "../../../types/api";
import type { RoadmapData } from "../../../hooks/useRoadmap";

function removeIssueFromTree(
  issues: RoadmapIssue[],
  id: string
): RoadmapIssue[] {
  return issues
    .filter((issue) => issue.id !== id)
    .map((issue) =>
      issue.children?.length
        ? { ...issue, children: removeIssueFromTree(issue.children, id) }
        : issue
    );
}

export function useDeleteIssueMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ issueId, force }: { issueId: string; force?: boolean }) =>
      deleteIssue(issueId, force),

    onMutate: async ({ issueId }) => {
      await queryClient.cancelQueries({ queryKey: ["roadmap"] });
      const previousData = queryClient.getQueryData<RoadmapData>(["roadmap"]);
      queryClient.setQueryData<RoadmapData>(["roadmap"], (old) => {
        if (!old) return old;
        return { ...old, issues: removeIssueFromTree(old.issues, issueId) };
      });
      return { previousData };
    },

    onError: (_err, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(["roadmap"], context.previousData);
      }
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["roadmap"] });
    },
  });
}
