import type { RoadmapIssue } from "../../types/api";

export interface TreeNode {
  id: string;
  title: string;
  status: string;
  priority: number;
  issue_type?: string;
  issue: RoadmapIssue;
  children: TreeNode[];
  x?: number;
  y?: number;
}

export interface TreeDimensions {
  width: number;
  height: number;
  nodeWidth: number;
  nodeHeight: number;
  horizontalSpacing: number;
  verticalSpacing: number;
}

export interface PopupPosition {
  x: number;
  y: number;
}

export interface BeadPopupProps {
  issue: RoadmapIssue;
  position: PopupPosition;
  onClose: () => void;
  onEditPlan: (issueId: string) => void;
  onViewFull: (issueId: string) => void;
}
