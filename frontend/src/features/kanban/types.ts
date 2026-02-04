import type { KanbanLanes, IssueCard } from "../../types/api";

export interface KanbanProps {
  featureId: string | null;
}

export interface LaneConfig {
  key: keyof KanbanLanes;
  title: string;
  bgColor: string;
  borderColor: string;
}

export interface KanbanCardProps {
  issue: IssueCard;
  onClick: () => void;
  onDragStart?: (e: React.DragEvent) => void;
  onDragEnd?: (e: React.DragEvent) => void;
  isDragging?: boolean;
}

export interface KanbanLaneProps {
  config: LaneConfig;
  issues: IssueCard[];
  onCardClick: (issue: IssueCard) => void;
  onDragStart?: (card: IssueCard) => void;
  onDragEnd?: () => void;
  onDragOver?: (e: React.DragEvent) => void;
  onDragLeave?: (e: React.DragEvent) => void;
  onDrop?: (e: React.DragEvent) => void;
  isDropTarget?: boolean;
  draggedCardId?: string | null;
}

export interface IssueModalProps {
  issueId: string;
  onClose: () => void;
}
