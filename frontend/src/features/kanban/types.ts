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
}

export interface KanbanLaneProps {
  config: LaneConfig;
  issues: IssueCard[];
  onCardClick: (issue: IssueCard) => void;
}

export interface IssueModalProps {
  issueId: string;
  onClose: () => void;
}
