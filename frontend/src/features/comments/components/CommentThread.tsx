import { MessageCircle, Trash2 } from "lucide-react";
import { useComments } from "../hooks/useComments";
import { CommentInput } from "./CommentInput";
import type { Comment } from "../types";

interface CommentThreadProps {
  issueId: string;
  className?: string;
}

function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

  if (diffMins < 1) return "just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;

  return date.toLocaleDateString();
}

function highlightMentions(content: string): React.ReactNode {
  const parts = content.split(/(@\w+)/g);
  return parts.map((part, i) => {
    if (part.startsWith("@")) {
      return (
        <span key={i} className="text-neon-cyan font-medium">
          {part}
        </span>
      );
    }
    return part;
  });
}

function CommentItem({
  comment,
  onDelete,
}: {
  comment: Comment;
  onDelete: (id: string) => void;
}) {
  return (
    <div className="group flex gap-3 p-3 rounded-lg hover:bg-night-surface transition-colors">
      {/* Avatar */}
      <div className="w-8 h-8 rounded-full bg-night-surface-bright flex items-center justify-center text-text-muted text-sm font-medium">
        {comment.author[0].toUpperCase()}
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className="font-medium text-text-bright text-sm">
            {comment.author}
          </span>
          <span className="text-xs text-text-muted">
            {formatTimestamp(comment.createdAt)}
          </span>
          {comment.updatedAt && (
            <span className="text-xs text-text-muted">(edited)</span>
          )}
        </div>
        <p className="text-text-secondary text-sm whitespace-pre-wrap break-words">
          {highlightMentions(comment.content)}
        </p>
      </div>

      {/* Actions */}
      <button
        onClick={() => onDelete(comment.id)}
        className="opacity-0 group-hover:opacity-100 p-1 text-text-muted hover:text-neon-pink transition-all"
        title="Delete comment"
      >
        <Trash2 className="w-4 h-4" />
      </button>
    </div>
  );
}

export function CommentThread({ issueId, className = "" }: CommentThreadProps) {
  const { comments, isSubmitting, addComment, deleteComment } =
    useComments(issueId);

  const handleSubmit = (content: string) => {
    addComment(content, "You");
  };

  return (
    <div className={`flex flex-col ${className}`}>
      {/* Header */}
      <div className="flex items-center gap-2 px-4 py-3 border-b border-night-border">
        <MessageCircle className="w-4 h-4 text-neon-cyan" />
        <span className="text-sm font-medium text-text-primary">
          Comments ({comments.length})
        </span>
      </div>

      {/* Comments list */}
      <div className="flex-1 overflow-auto p-2">
        {comments.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <MessageCircle className="w-8 h-8 text-text-muted mb-2 opacity-50" />
            <p className="text-text-muted text-sm">No comments yet</p>
            <p className="text-text-secondary text-xs mt-1">
              Be the first to add a comment
            </p>
          </div>
        ) : (
          <div className="space-y-1">
            {comments.map((comment) => (
              <CommentItem
                key={comment.id}
                comment={comment}
                onDelete={deleteComment}
              />
            ))}
          </div>
        )}
      </div>

      {/* Input */}
      <div className="p-3 border-t border-night-border">
        <CommentInput onSubmit={handleSubmit} isSubmitting={isSubmitting} />
      </div>
    </div>
  );
}
