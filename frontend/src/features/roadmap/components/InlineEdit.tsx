import { useState, useRef, useEffect } from "react";
import type { RoadmapIssue } from "../../../types/api";
import { useRoadmapContext } from "../context";
import { useIssueMutation } from "../hooks/useIssueMutation";

interface InlineEditProps {
  issue: RoadmapIssue;
}

export function InlineEdit({ issue }: InlineEditProps) {
  const [value, setValue] = useState(issue.title);
  const inputRef = useRef<HTMLInputElement>(null);
  const { stopEditingTitle } = useRoadmapContext();
  const { mutate } = useIssueMutation();

  useEffect(() => {
    inputRef.current?.focus();
    inputRef.current?.select();
  }, []);

  const handleSave = () => {
    const trimmed = value.trim();
    if (trimmed && trimmed !== issue.title) {
      mutate({ issueId: issue.id, updates: { title: trimmed } });
    }
    stopEditingTitle();
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSave();
    } else if (e.key === "Escape") {
      e.preventDefault();
      stopEditingTitle();
    }
  };

  return (
    <input
      ref={inputRef}
      value={value}
      onChange={(e) => setValue(e.target.value)}
      onBlur={handleSave}
      onKeyDown={handleKeyDown}
      onClick={(e) => e.stopPropagation()}
      className="flex-1 bg-transparent border-b-2 border-neon-cyan text-text-bright font-medium outline-none px-1 py-0.5"
      placeholder="Enter title..."
    />
  );
}
