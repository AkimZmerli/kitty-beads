import { useState, useRef, useCallback } from "react";
import type { HierarchyPointNode } from "d3";
import { useNavigate } from "react-router-dom";
import type { RoadmapIssue } from "../../../types/api";
import { useTreeLayout } from "../hooks/useTreeLayout";
import { TreeNodeComponent } from "./TreeNode";
import { TreeEdge } from "./TreeEdge";
import { BeadPopup } from "./BeadPopup";
import type { TreeNode, PopupPosition } from "../types";

interface TreeGraphProps {
  rootIssue: RoadmapIssue | null;
}

export function TreeGraph({ rootIssue }: TreeGraphProps) {
  const navigate = useNavigate();
  const svgRef = useRef<SVGSVGElement>(null);
  const [selectedNode, setSelectedNode] =
    useState<HierarchyPointNode<TreeNode> | null>(null);
  const [popupPosition, setPopupPosition] = useState<PopupPosition | null>(null);
  const [viewingIssue, setViewingIssue] = useState<RoadmapIssue | null>(null);

  const layout = useTreeLayout(rootIssue);

  const handleNodeClick = useCallback(
    (node: HierarchyPointNode<TreeNode>, event: React.MouseEvent) => {
      event.stopPropagation();

      // Get click position relative to viewport
      const rect = svgRef.current?.getBoundingClientRect();
      if (rect) {
        setPopupPosition({
          x: event.clientX - rect.left,
          y: event.clientY - rect.top,
        });
      }

      setSelectedNode(node);
      setViewingIssue(node.data.issue);
    },
    []
  );

  const handleClosePopup = useCallback(() => {
    setSelectedNode(null);
    setPopupPosition(null);
    setViewingIssue(null);
  }, []);

  const handleEditPlan = useCallback(
    (issueId: string) => {
      navigate(`/ideation/${issueId}`);
    },
    [navigate]
  );

  const handleViewFull = useCallback((issueId: string) => {
    // For now, navigate to ideation pad. Could open a modal instead.
    navigate(`/ideation/${issueId}`);
  }, [navigate]);

  if (!layout || !rootIssue) {
    return (
      <div className="flex items-center justify-center h-96 text-text-muted">
        Select an epic to view its tree
      </div>
    );
  }

  const { nodes, links, width, height } = layout;

  // Add padding
  const padding = 50;
  const viewBoxWidth = width + padding * 2;
  const viewBoxHeight = height + padding * 2;

  return (
    <div className="relative overflow-auto bg-night-bg rounded-lg border border-night-border">
      <svg
        ref={svgRef}
        width="100%"
        height={Math.max(400, viewBoxHeight)}
        viewBox={`${-padding} ${-padding} ${viewBoxWidth} ${viewBoxHeight}`}
        className="min-w-full"
        onClick={handleClosePopup}
      >
        {/* Edges */}
        <g className="edges">
          {links.map((link, i) => (
            <TreeEdge key={i} link={link} />
          ))}
        </g>

        {/* Nodes */}
        <g className="nodes">
          {nodes.map((node) => (
            <TreeNodeComponent
              key={node.data.id}
              node={node}
              onClick={handleNodeClick}
              isSelected={selectedNode?.data.id === node.data.id}
            />
          ))}
        </g>
      </svg>

      {/* Popup */}
      {viewingIssue && popupPosition && (
        <BeadPopup
          issue={viewingIssue}
          position={popupPosition}
          onClose={handleClosePopup}
          onEditPlan={handleEditPlan}
          onViewFull={handleViewFull}
        />
      )}
    </div>
  );
}
