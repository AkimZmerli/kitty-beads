import { useState, useCallback } from "react";
import type { IssueCard, KanbanLanes } from "../../../types/api";

export interface DragState {
  isDragging: boolean;
  draggedCard: IssueCard | null;
  sourceColumn: keyof KanbanLanes | null;
  targetColumn: keyof KanbanLanes | null;
}

export interface DragDropHandlers {
  onDragStart: (card: IssueCard, sourceColumn: keyof KanbanLanes) => void;
  onDragOver: (targetColumn: keyof KanbanLanes) => void;
  onDragLeave: () => void;
  onDrop: (targetColumn: keyof KanbanLanes) => void;
  onDragEnd: () => void;
}

export function useDragDrop(
  onMoveCard: (
    card: IssueCard,
    fromColumn: keyof KanbanLanes,
    toColumn: keyof KanbanLanes
  ) => void
): [DragState, DragDropHandlers] {
  const [dragState, setDragState] = useState<DragState>({
    isDragging: false,
    draggedCard: null,
    sourceColumn: null,
    targetColumn: null,
  });

  const onDragStart = useCallback(
    (card: IssueCard, sourceColumn: keyof KanbanLanes) => {
      setDragState({
        isDragging: true,
        draggedCard: card,
        sourceColumn,
        targetColumn: null,
      });
    },
    []
  );

  const onDragOver = useCallback((targetColumn: keyof KanbanLanes) => {
    setDragState((prev) => ({
      ...prev,
      targetColumn,
    }));
  }, []);

  const onDragLeave = useCallback(() => {
    setDragState((prev) => ({
      ...prev,
      targetColumn: null,
    }));
  }, []);

  const onDrop = useCallback(
    (targetColumn: keyof KanbanLanes) => {
      const { draggedCard, sourceColumn } = dragState;
      if (draggedCard && sourceColumn && sourceColumn !== targetColumn) {
        onMoveCard(draggedCard, sourceColumn, targetColumn);
      }
      setDragState({
        isDragging: false,
        draggedCard: null,
        sourceColumn: null,
        targetColumn: null,
      });
    },
    [dragState, onMoveCard]
  );

  const onDragEnd = useCallback(() => {
    setDragState({
      isDragging: false,
      draggedCard: null,
      sourceColumn: null,
      targetColumn: null,
    });
  }, []);

  return [
    dragState,
    { onDragStart, onDragOver, onDragLeave, onDrop, onDragEnd },
  ];
}
