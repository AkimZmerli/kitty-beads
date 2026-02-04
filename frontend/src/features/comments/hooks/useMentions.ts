import { useState, useCallback, useMemo } from "react";

// Sample team members for mention suggestions
const TEAM_MEMBERS = [
  "alice",
  "bob",
  "charlie",
  "diana",
  "eve",
  "frank",
  "system",
  "claude",
];

interface MentionState {
  isOpen: boolean;
  query: string;
  position: { top: number; left: number } | null;
}

export function useMentions() {
  const [state, setState] = useState<MentionState>({
    isOpen: false,
    query: "",
    position: null,
  });

  const suggestions = useMemo(() => {
    if (!state.query) return TEAM_MEMBERS.slice(0, 5);
    const query = state.query.toLowerCase();
    return TEAM_MEMBERS.filter((member) =>
      member.toLowerCase().includes(query)
    ).slice(0, 5);
  }, [state.query]);

  const openMentions = useCallback(
    (query: string, position: { top: number; left: number }) => {
      setState({ isOpen: true, query, position });
    },
    []
  );

  const closeMentions = useCallback(() => {
    setState({ isOpen: false, query: "", position: null });
  }, []);

  const updateQuery = useCallback((query: string) => {
    setState((prev) => ({ ...prev, query }));
  }, []);

  return {
    isOpen: state.isOpen,
    query: state.query,
    position: state.position,
    suggestions,
    openMentions,
    closeMentions,
    updateQuery,
  };
}
