import { useMemo } from "react";
import * as d3 from "d3";
import type { RoadmapIssue } from "../../../types/api";
import type { TreeNode } from "../types";
import { TREE_DIMENSIONS } from "../constants";

interface TreeLayoutResult {
  nodes: d3.HierarchyPointNode<TreeNode>[];
  links: d3.HierarchyPointLink<TreeNode>[];
  width: number;
  height: number;
}

function issueToTreeNode(issue: RoadmapIssue): TreeNode {
  return {
    id: issue.id,
    title: issue.title,
    status: issue.status,
    priority: issue.priority,
    issue_type: issue.issue_type,
    issue,
    children: (issue.children || []).map(issueToTreeNode),
  };
}

export function useTreeLayout(
  rootIssue: RoadmapIssue | null
): TreeLayoutResult | null {
  return useMemo(() => {
    if (!rootIssue) return null;

    const rootNode = issueToTreeNode(rootIssue);

    // Create D3 hierarchy
    const hierarchy = d3.hierarchy<TreeNode>(rootNode);

    // Calculate tree dimensions based on node count
    const maxDepth = hierarchy.height;
    const leafCount = hierarchy.leaves().length;

    // Dynamic sizing
    const width = Math.max(
      TREE_DIMENSIONS.width,
      (maxDepth + 1) * TREE_DIMENSIONS.horizontalSpacing + 100
    );
    const height = Math.max(
      TREE_DIMENSIONS.height,
      leafCount * TREE_DIMENSIONS.verticalSpacing
    );

    // Create tree layout
    const treeLayout = d3
      .tree<TreeNode>()
      .size([height - 100, width - 200])
      .separation((a, b) => (a.parent === b.parent ? 1 : 1.2));

    // Apply layout
    const root = treeLayout(hierarchy);

    return {
      nodes: root.descendants(),
      links: root.links(),
      width,
      height,
    };
  }, [rootIssue]);
}
