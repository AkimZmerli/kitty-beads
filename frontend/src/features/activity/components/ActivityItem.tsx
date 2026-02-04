import { Link } from "react-router-dom";
import {
  Plus,
  Edit3,
  RefreshCw,
  FileText,
  MessageCircle,
  UserPlus,
} from "lucide-react";
import type { ActivityEvent, ActivityType } from "../types";

interface ActivityItemProps {
  event: ActivityEvent;
}

const TYPE_CONFIG: Record<
  ActivityType,
  { icon: typeof Plus; color: string; label: string }
> = {
  issue_created: {
    icon: Plus,
    color: "text-neon-green",
    label: "created",
  },
  issue_updated: {
    icon: Edit3,
    color: "text-neon-cyan",
    label: "updated",
  },
  status_changed: {
    icon: RefreshCw,
    color: "text-neon-orange",
    label: "changed status",
  },
  plan_edited: {
    icon: FileText,
    color: "text-neon-blue",
    label: "edited plan",
  },
  comment_added: {
    icon: MessageCircle,
    color: "text-neon-magenta",
    label: "commented",
  },
  assignee_changed: {
    icon: UserPlus,
    color: "text-neon-yellow",
    label: "assigned",
  },
};

function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return "just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString();
}

export function ActivityItem({ event }: ActivityItemProps) {
  const config = TYPE_CONFIG[event.type];
  const Icon = config.icon;

  return (
    <div className="flex items-start gap-3 px-4 py-3 hover:bg-night-surface-bright transition-colors">
      {/* Icon */}
      <div
        className={`p-2 rounded-lg bg-night-surface ${config.color}`}
      >
        <Icon className="w-4 h-4" />
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          {event.actor && (
            <span className="font-medium text-text-bright">{event.actor}</span>
          )}
          <span className="text-text-muted">{config.label}</span>
          <Link
            to={`/ideation/${event.issueId}`}
            className="font-medium text-neon-cyan hover:underline truncate max-w-[200px]"
          >
            {event.issueTitle}
          </Link>
        </div>

        {/* Details */}
        {event.details && (
          <div className="mt-1 text-sm text-text-muted">
            {event.type === "status_changed" && event.details.newValue && (
              <span className="px-2 py-0.5 bg-night-surface rounded text-text-secondary">
                → {event.details.newValue.replace("_", " ")}
              </span>
            )}
            {event.type === "assignee_changed" && event.details.newValue && (
              <span>→ {event.details.newValue}</span>
            )}
            {event.details.comment && (
              <p className="mt-1 text-text-secondary line-clamp-2">
                "{event.details.comment}"
              </p>
            )}
          </div>
        )}
      </div>

      {/* Timestamp */}
      <span className="text-xs text-text-muted whitespace-nowrap">
        {formatTimestamp(event.timestamp)}
      </span>
    </div>
  );
}
