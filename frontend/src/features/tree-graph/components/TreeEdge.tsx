import type { HierarchyPointLink } from "d3";
import * as d3 from "d3";
import type { TreeNode } from "../types";

interface TreeEdgeProps {
  link: HierarchyPointLink<TreeNode>;
}

export function TreeEdge({ link }: TreeEdgeProps) {
  // Create a curved path from source to target
  // Note: we swap x/y because the tree is horizontal (x is vertical, y is horizontal in d3.tree)
  const path = d3
    .linkHorizontal<HierarchyPointLink<TreeNode>, HierarchyPointLink<TreeNode>["source"]>()
    .x((d) => d.y)
    .y((d) => d.x);

  const pathD = path(link);

  if (!pathD) return null;

  return (
    <path
      d={pathD}
      fill="none"
      stroke="#3b4261"
      strokeWidth={1.5}
      strokeOpacity={0.6}
      className="transition-all duration-150"
    />
  );
}
