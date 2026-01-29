import type { HierarchyPointNode } from "d3";
import type { TreeNode as TreeNodeType } from "../types";
import { NODE_COLORS, PRIORITY_COLORS, TREE_DIMENSIONS } from "../constants";

interface TreeNodeProps {
  node: HierarchyPointNode<TreeNodeType>;
  onClick: (node: HierarchyPointNode<TreeNodeType>, event: React.MouseEvent) => void;
  isSelected?: boolean;
}

function getNodeColors(node: TreeNodeType) {
  const status = node.status;
  const isEpic = node.issue_type === "epic" || (node.children?.length ?? 0) > 0;

  if (status === "done" || status === "closed" || status === "completed") {
    return NODE_COLORS.done;
  }
  if (status === "in_progress") {
    return NODE_COLORS.inProgress;
  }
  if (isEpic) {
    return NODE_COLORS.epic;
  }
  return NODE_COLORS.task;
}

export function TreeNodeComponent({ node, onClick, isSelected }: TreeNodeProps) {
  const colors = getNodeColors(node.data);
  const priorityColor = PRIORITY_COLORS[node.data.priority] || PRIORITY_COLORS[3];
  const isEpic =
    node.data.issue_type === "epic" || (node.data.children?.length ?? 0) > 0;

  const nodeWidth = TREE_DIMENSIONS.nodeWidth;
  const nodeHeight = TREE_DIMENSIONS.nodeHeight;

  // Truncate title if too long
  const maxTitleLength = 20;
  const displayTitle =
    node.data.title.length > maxTitleLength
      ? node.data.title.substring(0, maxTitleLength) + "..."
      : node.data.title;

  return (
    <g
      transform={`translate(${node.y},${node.x})`}
      onClick={(e) => onClick(node, e)}
      className="cursor-pointer"
      style={{ transition: "transform 0.2s ease" }}
    >
      {/* Node background */}
      <rect
        x={-nodeWidth / 2}
        y={-nodeHeight / 2}
        width={nodeWidth}
        height={nodeHeight}
        rx={8}
        ry={8}
        fill={colors.fill}
        stroke={isSelected ? "#fff" : colors.stroke}
        strokeWidth={isSelected ? 2 : 1.5}
        className="transition-all duration-150 hover:brightness-125"
      />

      {/* Priority indicator bar */}
      <rect
        x={-nodeWidth / 2}
        y={-nodeHeight / 2}
        width={4}
        height={nodeHeight}
        rx={2}
        fill={priorityColor}
      />

      {/* Type badge */}
      {isEpic && (
        <g transform={`translate(${nodeWidth / 2 - 30}, ${-nodeHeight / 2 + 8})`}>
          <rect
            width={24}
            height={14}
            rx={3}
            fill="rgba(187, 154, 247, 0.3)"
            stroke="#bb9af7"
            strokeWidth={0.5}
          />
          <text
            x={12}
            y={10}
            textAnchor="middle"
            fontSize={8}
            fill="#bb9af7"
            fontWeight="bold"
          >
            EPIC
          </text>
        </g>
      )}

      {/* Issue ID */}
      <text
        x={-nodeWidth / 2 + 12}
        y={-nodeHeight / 2 + 16}
        fontSize={10}
        fill="#565f89"
        fontFamily="monospace"
      >
        {node.data.id}
      </text>

      {/* Title */}
      <text
        x={0}
        y={5}
        textAnchor="middle"
        fontSize={12}
        fill={colors.text}
        fontWeight="500"
      >
        {displayTitle}
      </text>

      {/* Status badge */}
      <text
        x={0}
        y={nodeHeight / 2 - 8}
        textAnchor="middle"
        fontSize={9}
        fill="#565f89"
        style={{ textTransform: "uppercase" }}
      >
        {node.data.status}
      </text>
    </g>
  );
}
