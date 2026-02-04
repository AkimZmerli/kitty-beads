import { useRef, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import type { RoadmapIssue } from "../../../types/api";
import { parsePlan, formatCriteriaProgress } from "../../../lib/planParser";
import { CheckCircle2, HelpCircle, ListTodo, ChevronRight } from "lucide-react";
import { useRoadmapContext } from "../context";
import { InlineEdit } from "./InlineEdit";
import { StatusDropdown } from "./StatusDropdown";

interface BeadCardProps {
  issue: RoadmapIssue;
  depth?: number;
  isLastChild?: boolean;
}

// Depth-based background opacity
const DEPTH_OPACITY = [1, 0.85, 0.7, 0.55, 0.4];

// Depth-based visual priority (for hierarchy display)
// depth 0 = P1 (red), depth 1 = P2 (orange), depth 2+ = P3 (yellow)
const DEPTH_PRIORITY_COLORS = ["#f7768e", "#ff9e64", "#e0af68", "#9ece6a"];
const DEPTH_PRIORITY_LABELS = ["P1", "P2", "P3", "P4"];

// Indentation per depth level (pixels)
const INDENT_PER_DEPTH = 40;

export function BeadCard({ issue, depth = 0, isLastChild = false }: BeadCardProps) {
  const navigate = useNavigate();
  const cardRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const [contentHeight, setContentHeight] = useState(0);
  const [showStatusDropdown, setShowStatusDropdown] = useState(false);

  const {
    selectedIssueId,
    selectIssue,
    isExpanded,
    toggleExpand,
    showContextMenu,
    editingTitleId,
    openModal,
  } = useRoadmapContext();

  const expanded = isExpanded(issue.id);
  const isSelected = selectedIssueId === issue.id;
  const isEditing = editingTitleId === issue.id;

  const parsedPlan = parsePlan(issue.design || issue.description);

  // Use depth-based priority for visual hierarchy
  const visualPriorityIndex = Math.min(depth, DEPTH_PRIORITY_COLORS.length - 1);
  const priorityColor = DEPTH_PRIORITY_COLORS[visualPriorityIndex];
  const priorityLabel = DEPTH_PRIORITY_LABELS[visualPriorityIndex];
  const priorityIndicator = depth === 0 ? "●" : "○";

  const hasChildren = issue.children && issue.children.length > 0;
  const isEpic = issue.issue_type === "epic" || hasChildren;
  const hasExpandableContent = hasChildren || parsedPlan.hasContent;

  const statusBadge = getStatusBadge(issue.status);

  // Calculate content height for animation
  useEffect(() => {
    if (contentRef.current) {
      setContentHeight(contentRef.current.scrollHeight);
    }
  }, [expanded, issue]);

  // Scroll into view when selected via keyboard
  useEffect(() => {
    if (isSelected && cardRef.current) {
      cardRef.current.scrollIntoView({ behavior: "smooth", block: "nearest" });
    }
  }, [isSelected]);

  const handleClick = () => {
    selectIssue(issue.id);
    if (hasExpandableContent) {
      toggleExpand(issue.id);
    }
  };

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    showContextMenu(e.clientX, e.clientY, issue.id);
  };

  const handleEditPlan = (e: React.MouseEvent) => {
    e.stopPropagation();
    navigate(`/ideation/${issue.id}`);
  };

  const handleViewFull = (e: React.MouseEvent) => {
    e.stopPropagation();
    openModal(issue.id);
  };

  // Depth-based styling
  const depthOpacity = DEPTH_OPACITY[Math.min(depth, DEPTH_OPACITY.length - 1)];
  const marginLeft = depth * INDENT_PER_DEPTH;

  return (
    <div className="relative" style={{ marginLeft }}>
      {/* Tree connector lines */}
      {depth > 0 && (
        <div
          className="absolute left-0 top-0 bottom-0 pointer-events-none"
          style={{ marginLeft: -INDENT_PER_DEPTH }}
        >
          {/* Vertical line */}
          <div
            className={`absolute w-px ${isLastChild ? "h-[22px]" : "h-full"}`}
            style={{
              left: 12,
              backgroundColor: DEPTH_PRIORITY_COLORS[Math.min(depth - 1, DEPTH_PRIORITY_COLORS.length - 1)],
              opacity: 0.4,
            }}
          />
          {/* Horizontal connector */}
          <div
            className="absolute top-[22px] h-px"
            style={{
              left: 12,
              width: INDENT_PER_DEPTH - 12,
              backgroundColor: DEPTH_PRIORITY_COLORS[Math.min(depth - 1, DEPTH_PRIORITY_COLORS.length - 1)],
              opacity: 0.4,
            }}
          />
        </div>
      )}

      {/* Main card row */}
      <div
        ref={cardRef}
        data-issue-id={issue.id}
        className={`
          flex items-center gap-3 px-4 py-3 rounded-lg cursor-pointer
          transition-all duration-200 ease-out
          ${isSelected ? "ring-2 ring-neon-cyan ring-offset-1 ring-offset-night-bg" : ""}
          ${expanded ? "bg-night-bg-highlight" : "hover:bg-night-bg-highlight"}
          border-l-2
        `}
        style={{
          backgroundColor: expanded
            ? undefined
            : `rgba(31, 35, 53, ${depthOpacity})`,
          borderLeftColor: priorityColor,
        }}
        onClick={handleClick}
        onContextMenu={handleContextMenu}
        tabIndex={isSelected ? 0 : -1}
      >
        {/* Expand/collapse chevron */}
        <span
          className={`
            text-text-muted
            ${expanded ? "rotate-90" : ""}
            ${hasExpandableContent ? "" : "opacity-0"}
          `}
          style={{
            transition: "transform 400ms cubic-bezier(0.4, 0, 0.2, 1)",
          }}
        >
          <ChevronRight className="w-4 h-4" />
        </span>

        {/* Priority indicator - shows depth-based hierarchy */}
        <span
          className="font-bold"
          style={{ color: priorityColor }}
        >
          {priorityIndicator} {priorityLabel}
        </span>

        {/* Issue ID */}
        <span className="text-text-muted font-mono text-sm">{issue.id}</span>

        {/* Title - editable on double click */}
        {isEditing ? (
          <InlineEdit issue={issue} />
        ) : (
          <span
            className="font-medium text-text-bright flex-1 truncate"
            onDoubleClick={(e) => {
              e.stopPropagation();
              // startEditingTitle handled by context menu or keyboard
            }}
          >
            {issue.title}
          </span>
        )}

        {/* Type badge */}
        {isEpic && (
          <span className="px-2 py-0.5 bg-purple-900/30 text-neon-magenta text-xs rounded-full font-medium">
            EPIC
          </span>
        )}

        {/* Status badge - clickable */}
        <button
          onClick={(e) => {
            e.stopPropagation();
            setShowStatusDropdown(true);
          }}
          className={`px-2 py-0.5 text-xs rounded-full font-medium transition-all hover:ring-1 hover:ring-white/20 ${statusBadge.className}`}
        >
          {statusBadge.label}
        </button>

        {/* Blocked indicator */}
        {issue.isBlocked && (
          <span className="px-2 py-0.5 bg-red-900/30 text-neon-pink text-xs rounded-full font-medium animate-pulse">
            BLOCKED
          </span>
        )}
      </div>

      {/* Expanded content with smooth animation */}
      <div
        className="overflow-hidden"
        style={{
          maxHeight: expanded ? `${contentHeight + 500}px` : "0px",
          opacity: expanded ? 1 : 0,
          transition: expanded
            ? "max-height 600ms cubic-bezier(0.4, 0, 0.2, 1), opacity 500ms ease-out"
            : "max-height 350ms cubic-bezier(0.4, 0, 0.2, 1), opacity 250ms ease-out",
          transitionDelay: expanded ? "0ms, 100ms" : "50ms, 0ms",
        }}
      >
        <div ref={contentRef}>
          {/* Expanded details panel */}
          <div
            className="mt-2 px-4 py-4 rounded-lg border border-night-border"
            style={{
              backgroundColor: `rgba(41, 46, 66, ${depthOpacity})`,
            }}
          >
            {/* Summary - reveal 1 */}
            {parsedPlan.summary && (
              <p
                className="text-text-normal mb-4 leading-relaxed"
                style={{
                  opacity: expanded ? 1 : 0,
                  transform: expanded ? "translateY(0)" : "translateY(-8px)",
                  transition: "opacity 400ms ease-out, transform 400ms ease-out",
                  transitionDelay: expanded ? "150ms" : "0ms",
                }}
              >
                {parsedPlan.summary}
              </p>
            )}

            {/* Stats row - reveal 2 */}
            <div
              className="flex items-center gap-4 text-sm text-text-muted mb-4"
              style={{
                opacity: expanded ? 1 : 0,
                transform: expanded ? "translateY(0)" : "translateY(-8px)",
                transition: "opacity 400ms ease-out, transform 400ms ease-out",
                transitionDelay: expanded ? "250ms" : "0ms",
              }}
            >
              {/* Acceptance criteria */}
              {parsedPlan.acceptanceCriteriaCount > 0 && (
                <div className="flex items-center gap-1.5">
                  <CheckCircle2 className="w-4 h-4 text-neon-green" />
                  <span>{formatCriteriaProgress(parsedPlan)}</span>
                </div>
              )}

              {/* Open questions */}
              {parsedPlan.openQuestions.length > 0 && (
                <div className="flex items-center gap-1.5">
                  <HelpCircle className="w-4 h-4 text-neon-orange" />
                  <span>{parsedPlan.openQuestions.length} questions</span>
                </div>
              )}

              {/* Children count */}
              {hasChildren && (
                <div className="flex items-center gap-1.5">
                  <ListTodo className="w-4 h-4 text-neon-blue" />
                  <span>{issue.children!.length} subtasks</span>
                </div>
              )}
            </div>

            {/* Open questions preview - reveal 3 */}
            {parsedPlan.openQuestions.length > 0 && (
              <div
                className="mb-4"
                style={{
                  opacity: expanded ? 1 : 0,
                  transform: expanded ? "translateY(0)" : "translateY(-8px)",
                  transition: "opacity 400ms ease-out, transform 400ms ease-out",
                  transitionDelay: expanded ? "350ms" : "0ms",
                }}
              >
                <div className="text-xs uppercase tracking-wider text-text-muted mb-2">
                  Open Questions
                </div>
                <ul className="text-sm text-text-normal space-y-1">
                  {parsedPlan.openQuestions.slice(0, 3).map((q, i) => (
                    <li key={i} className="flex items-start gap-2">
                      <span className="text-neon-orange">?</span>
                      <span className="truncate">{q}</span>
                    </li>
                  ))}
                  {parsedPlan.openQuestions.length > 3 && (
                    <li className="text-text-muted italic">
                      +{parsedPlan.openQuestions.length - 3} more...
                    </li>
                  )}
                </ul>
              </div>
            )}

            {/* Action buttons - reveal 4 */}
            <div
              className="flex gap-2"
              style={{
                opacity: expanded ? 1 : 0,
                transform: expanded ? "translateY(0)" : "translateY(-8px)",
                transition: "opacity 400ms ease-out, transform 400ms ease-out",
                transitionDelay: expanded ? "450ms" : "0ms",
              }}
            >
              <button
                onClick={handleEditPlan}
                className="px-3 py-1.5 bg-neon-cyan text-night-bg text-sm rounded-md
                  hover:shadow-[0_0_12px_rgba(125,207,255,0.4)] transition-all font-medium"
              >
                Edit Plan
              </button>
              <button
                onClick={handleViewFull}
                className="px-3 py-1.5 bg-night-surface text-text-normal text-sm rounded-md
                  hover:bg-night-bg-highlight transition-colors"
              >
                View Full
              </button>
            </div>
          </div>

          {/* Render children with stagger reveal */}
          {hasChildren && (
            <div className="mt-3 space-y-1">
              {issue.children!.map((child, index) => (
                <div
                  key={child.id}
                  style={{
                    opacity: expanded ? 1 : 0,
                    transform: expanded ? "translateX(0)" : "translateX(-12px)",
                    transition: "opacity 450ms ease-out, transform 450ms ease-out",
                    transitionDelay: expanded ? `${400 + index * 80}ms` : "0ms",
                  }}
                >
                  <BeadCard
                    issue={child}
                    depth={depth + 1}
                    isLastChild={index === issue.children!.length - 1}
                  />
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Dropdowns */}
      {showStatusDropdown && (
        <StatusDropdown
          issue={issue}
          onClose={() => setShowStatusDropdown(false)}
        />
      )}
    </div>
  );
}

function getStatusBadge(status: string): { label: string; className: string } {
  switch (status) {
    case "open":
      return { label: "OPEN", className: "bg-blue-900/30 text-neon-blue" };
    case "in_progress":
      return {
        label: "IN PROGRESS",
        className: "bg-yellow-900/30 text-neon-yellow",
      };
    case "closed":
      return { label: "DONE", className: "bg-green-900/30 text-neon-green" };
    case "blocked":
      return { label: "BLOCKED", className: "bg-red-900/30 text-neon-pink" };
    default:
      return {
        label: status.toUpperCase(),
        className: "bg-night-surface-bright text-text-normal",
      };
  }
}
