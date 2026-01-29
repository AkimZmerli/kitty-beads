import { useState, useCallback } from "react";
import type { Comment } from "../types";

// In-memory storage for comments (simulates backend)
// In production, this would be replaced with API calls
const commentStore = new Map<string, Comment[]>();

let commentIdCounter = 1;

function generateCommentId(): string {
  return `comment-${commentIdCounter++}`;
}

export function useComments(issueId: string) {
  const [comments, setComments] = useState<Comment[]>(() => {
    return commentStore.get(issueId) || [];
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const addComment = useCallback(
    async (content: string, author: string = "Anonymous") => {
      setIsSubmitting(true);

      // Simulate API delay
      await new Promise((resolve) => setTimeout(resolve, 300));

      const newComment: Comment = {
        id: generateCommentId(),
        issueId,
        content,
        author,
        createdAt: new Date().toISOString(),
        mentions: extractMentions(content),
      };

      const updatedComments = [...(commentStore.get(issueId) || []), newComment];
      commentStore.set(issueId, updatedComments);
      setComments(updatedComments);
      setIsSubmitting(false);

      return newComment;
    },
    [issueId]
  );

  const deleteComment = useCallback(
    async (commentId: string) => {
      const updatedComments = comments.filter((c) => c.id !== commentId);
      commentStore.set(issueId, updatedComments);
      setComments(updatedComments);
    },
    [issueId, comments]
  );

  const updateComment = useCallback(
    async (commentId: string, content: string) => {
      const updatedComments = comments.map((c) =>
        c.id === commentId
          ? {
              ...c,
              content,
              updatedAt: new Date().toISOString(),
              mentions: extractMentions(content),
            }
          : c
      );
      commentStore.set(issueId, updatedComments);
      setComments(updatedComments);
    },
    [issueId, comments]
  );

  return {
    comments,
    isSubmitting,
    addComment,
    deleteComment,
    updateComment,
  };
}

// Extract @mentions from content
function extractMentions(content: string): string[] {
  const mentionRegex = /@(\w+)/g;
  const matches = content.match(mentionRegex);
  if (!matches) return [];
  return [...new Set(matches.map((m) => m.slice(1)))];
}
