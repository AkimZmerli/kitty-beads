import { useState, useRef, useCallback } from "react";
import { Send } from "lucide-react";
import { useMentions } from "../hooks/useMentions";

interface CommentInputProps {
  onSubmit: (content: string) => void;
  isSubmitting?: boolean;
  placeholder?: string;
}

export function CommentInput({
  onSubmit,
  isSubmitting = false,
  placeholder = "Write a comment... Use @ to mention someone",
}: CommentInputProps) {
  const [content, setContent] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const mentions = useMentions();

  const handleSubmit = useCallback(() => {
    if (content.trim() && !isSubmitting) {
      onSubmit(content.trim());
      setContent("");
    }
  }, [content, isSubmitting, onSubmit]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      // Submit on Ctrl/Cmd + Enter
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        e.preventDefault();
        handleSubmit();
      }
      // Close mentions on Escape
      if (e.key === "Escape" && mentions.isOpen) {
        mentions.closeMentions();
      }
    },
    [handleSubmit, mentions]
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const value = e.target.value;
      setContent(value);

      // Check for @mention trigger
      const cursorPos = e.target.selectionStart;
      const textBeforeCursor = value.slice(0, cursorPos);
      const mentionMatch = textBeforeCursor.match(/@(\w*)$/);

      if (mentionMatch) {
        const rect = textareaRef.current?.getBoundingClientRect();
        if (rect) {
          mentions.openMentions(mentionMatch[1], {
            top: rect.bottom + 4,
            left: rect.left,
          });
        }
      } else if (mentions.isOpen) {
        mentions.closeMentions();
      }
    },
    [mentions]
  );

  const insertMention = useCallback(
    (username: string) => {
      const textarea = textareaRef.current;
      if (!textarea) return;

      const cursorPos = textarea.selectionStart;
      const textBeforeCursor = content.slice(0, cursorPos);
      const textAfterCursor = content.slice(cursorPos);

      // Find the start of the @mention
      const mentionStart = textBeforeCursor.lastIndexOf("@");
      const newContent =
        textBeforeCursor.slice(0, mentionStart) +
        `@${username} ` +
        textAfterCursor;

      setContent(newContent);
      mentions.closeMentions();

      // Focus back on textarea
      setTimeout(() => {
        textarea.focus();
        const newCursorPos = mentionStart + username.length + 2;
        textarea.setSelectionRange(newCursorPos, newCursorPos);
      }, 0);
    },
    [content, mentions]
  );

  return (
    <div className="relative">
      <div className="flex gap-2">
        <textarea
          ref={textareaRef}
          value={content}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          disabled={isSubmitting}
          className="flex-1 min-h-[80px] p-3 bg-night-surface border border-night-border rounded-lg text-text-primary placeholder:text-text-muted resize-none focus:outline-none focus:border-neon-cyan transition-colors"
          rows={2}
        />
        <button
          onClick={handleSubmit}
          disabled={!content.trim() || isSubmitting}
          className={`px-4 self-end mb-1 py-2 rounded-lg transition-all ${
            content.trim() && !isSubmitting
              ? "bg-neon-cyan text-night-bg hover:shadow-[0_0_12px_rgba(125,207,255,0.4)]"
              : "bg-night-surface text-text-muted cursor-not-allowed"
          }`}
        >
          <Send className="w-4 h-4" />
        </button>
      </div>

      {/* Mention suggestions */}
      {mentions.isOpen && mentions.suggestions.length > 0 && (
        <div className="absolute z-50 mt-1 w-48 bg-night-surface border border-night-border rounded-lg shadow-xl overflow-hidden">
          {mentions.suggestions.map((username) => (
            <button
              key={username}
              onClick={() => insertMention(username)}
              className="w-full px-3 py-2 text-left text-sm text-text-primary hover:bg-night-surface-bright transition-colors"
            >
              @{username}
            </button>
          ))}
        </div>
      )}

      <p className="mt-2 text-xs text-text-muted">
        Press Ctrl+Enter to send
      </p>
    </div>
  );
}
