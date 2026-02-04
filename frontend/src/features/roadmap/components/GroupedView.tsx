import type { RoadmapIssue } from "../../../types/api";
import type { GroupByOption } from "../types";
import { groupIssues, getGroupInfo, getSortedGroupKeys } from "../utils/grouping";
import { BeadCard } from "./BeadCard";
import { getPriorityColor } from "../../../lib/planParser";

interface GroupedViewProps {
  issues: RoadmapIssue[];
  groupBy: GroupByOption;
}

export function GroupedView({ issues, groupBy }: GroupedViewProps) {
  const groups = groupIssues(issues, groupBy);
  const sortedKeys = getSortedGroupKeys(groups, groupBy);

  if (groupBy === "none") {
    return (
      <div className="space-y-1">
        {issues.map((issue, index) => (
          <BeadCard
            key={issue.id}
            issue={issue}
            isLastChild={index === issues.length - 1}
          />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {sortedKeys.map((key) => {
        const groupIssues = groups.get(key) || [];
        const groupInfo = getGroupInfo(key, groupBy);

        // Get color for the group header
        let headerColor = "#7dcfff"; // default neon-cyan
        if (groupBy === "priority") {
          const priority = parseInt(key.replace("p", ""));
          headerColor = getPriorityColor(priority);
        } else if (groupBy === "status") {
          const colorMap: Record<string, string> = {
            open: "#7aa2f7",
            in_progress: "#e0af68",
            blocked: "#f7768e",
            closed: "#9ece6a",
          };
          headerColor = colorMap[key] || headerColor;
        }

        return (
          <div key={key}>
            {/* Group header */}
            <div className="flex items-center gap-3 mb-3 px-2">
              <div
                className="w-3 h-3 rounded-full"
                style={{ backgroundColor: headerColor }}
              />
              <h3 className="font-semibold text-text-bright">
                {groupInfo.label}
              </h3>
              <span className="text-sm text-text-muted">
                {groupIssues.length} {groupIssues.length === 1 ? "item" : "items"}
              </span>
            </div>

            {/* Group issues */}
            <div className="space-y-1 pl-2 border-l-2 border-night-border ml-1">
              {groupIssues.map((issue, index) => (
                <BeadCard
                  key={issue.id}
                  issue={issue}
                  isLastChild={index === groupIssues.length - 1}
                />
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}
