import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateIssue } from "../../../lib/api";
import type { IssueCard, KanbanLanes, KanbanResponse } from "../../../types/api";

// Map column keys to status values
const COLUMN_TO_STATUS: Record<keyof KanbanLanes, string> = {
  planned: "open",
  doing: "in_progress",
  for_review: "review",
  done: "closed",
};

export function useKanbanMutation(featureId: string | null) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      card,
      toColumn,
    }: {
      card: IssueCard;
      fromColumn: keyof KanbanLanes;
      toColumn: keyof KanbanLanes;
    }) => {
      const newStatus = COLUMN_TO_STATUS[toColumn];
      return updateIssue(card.id, { status: newStatus });
    },

    // Optimistic update
    onMutate: async ({ card, fromColumn, toColumn }) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: ["kanban", featureId] });

      // Snapshot the previous value
      const previousData = queryClient.getQueryData<KanbanResponse>([
        "kanban",
        featureId,
      ]);

      // Optimistically update
      queryClient.setQueryData<KanbanResponse>(
        ["kanban", featureId],
        (old) => {
          if (!old) return old;

          const newLanes = { ...old.lanes };

          // Remove from source column
          newLanes[fromColumn] = newLanes[fromColumn].filter(
            (c) => c.id !== card.id
          );

          // Add to target column with updated status
          const updatedCard = { ...card, status: COLUMN_TO_STATUS[toColumn] };
          newLanes[toColumn] = [...newLanes[toColumn], updatedCard];

          return { ...old, lanes: newLanes };
        }
      );

      return { previousData };
    },

    // If error, roll back
    onError: (_err, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(["kanban", featureId], context.previousData);
      }
    },

    // Always refetch after error or success
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["kanban", featureId] });
      // Also invalidate roadmap since status changed
      queryClient.invalidateQueries({ queryKey: ["roadmap"] });
    },
  });
}
