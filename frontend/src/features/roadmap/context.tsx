import {
  createContext,
  useContext,
  useState,
  useCallback,
  useMemo,
  type ReactNode,
} from "react";
import type { RoadmapIssue } from "../../types/api";
import type { FilterType, GroupByOption } from "./types";

interface ContextMenuState {
  x: number;
  y: number;
  issueId: string;
}

interface RoadmapState {
  selectedIssueId: string | null;
  activeFilter: FilterType;
  searchQuery: string;
  groupBy: GroupByOption;
  expandedIds: Set<string>;
  contextMenu: ContextMenuState | null;
  editingTitleId: string | null;
  modalIssueId: string | null;
}

interface RoadmapContextValue extends RoadmapState {
  // Selection
  selectIssue: (id: string | null) => void;
  // Expand/collapse
  toggleExpand: (id: string) => void;
  expandIssue: (id: string) => void;
  collapseIssue: (id: string) => void;
  isExpanded: (id: string) => boolean;
  // Filters
  setFilter: (filter: FilterType) => void;
  setSearchQuery: (query: string) => void;
  setGroupBy: (groupBy: GroupByOption) => void;
  // Context menu
  showContextMenu: (x: number, y: number, issueId: string) => void;
  hideContextMenu: () => void;
  // Inline edit
  startEditingTitle: (id: string) => void;
  stopEditingTitle: () => void;
  // Modal
  openModal: (issueId: string) => void;
  closeModal: () => void;
  // Flattened issues for keyboard nav
  flattenedIssues: RoadmapIssue[];
  setFlattenedIssues: (issues: RoadmapIssue[]) => void;
  // Navigation helpers
  focusedIndex: number;
  focusNext: () => void;
  focusPrevious: () => void;
}

const RoadmapContext = createContext<RoadmapContextValue | null>(null);

export function RoadmapProvider({ children }: { children: ReactNode }) {
  const [selectedIssueId, setSelectedIssueId] = useState<string | null>(null);
  const [activeFilter, setActiveFilter] = useState<FilterType>("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [groupBy, setGroupBy] = useState<GroupByOption>("none");
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [contextMenu, setContextMenu] = useState<ContextMenuState | null>(null);
  const [editingTitleId, setEditingTitleId] = useState<string | null>(null);
  const [modalIssueId, setModalIssueId] = useState<string | null>(null);
  const [flattenedIssues, setFlattenedIssues] = useState<RoadmapIssue[]>([]);

  // Selection
  const selectIssue = useCallback((id: string | null) => {
    setSelectedIssueId(id);
  }, []);

  // Expand/collapse
  const toggleExpand = useCallback((id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  const expandIssue = useCallback((id: string) => {
    setExpandedIds((prev) => new Set(prev).add(id));
  }, []);

  const collapseIssue = useCallback((id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      next.delete(id);
      return next;
    });
  }, []);

  const isExpanded = useCallback(
    (id: string) => expandedIds.has(id),
    [expandedIds],
  );

  // Filters
  const setFilter = useCallback((filter: FilterType) => {
    setActiveFilter(filter);
  }, []);

  // Context menu
  const showContextMenu = useCallback(
    (x: number, y: number, issueId: string) => {
      setContextMenu({ x, y, issueId });
    },
    [],
  );

  const hideContextMenu = useCallback(() => {
    setContextMenu(null);
  }, []);

  // Inline edit
  const startEditingTitle = useCallback((id: string) => {
    setEditingTitleId(id);
  }, []);

  const stopEditingTitle = useCallback(() => {
    setEditingTitleId(null);
  }, []);

  // Modal
  const openModal = useCallback((issueId: string) => {
    setModalIssueId(issueId);
  }, []);

  const closeModal = useCallback(() => {
    setModalIssueId(null);
  }, []);

  // Navigation helpers
  const focusedIndex = useMemo(() => {
    if (!selectedIssueId) return -1;
    return flattenedIssues.findIndex((i) => i.id === selectedIssueId);
  }, [selectedIssueId, flattenedIssues]);

  const focusNext = useCallback(() => {
    if (flattenedIssues.length === 0) return;
    const nextIndex =
      focusedIndex < 0
        ? 0
        : Math.min(focusedIndex + 1, flattenedIssues.length - 1);
    setSelectedIssueId(flattenedIssues[nextIndex]?.id ?? null);
  }, [flattenedIssues, focusedIndex]);

  const focusPrevious = useCallback(() => {
    if (flattenedIssues.length === 0) return;
    const prevIndex = focusedIndex < 0 ? 0 : Math.max(focusedIndex - 1, 0);
    setSelectedIssueId(flattenedIssues[prevIndex]?.id ?? null);
  }, [flattenedIssues, focusedIndex]);

  const value: RoadmapContextValue = useMemo(
    () => ({
      // State
      selectedIssueId,
      activeFilter,
      searchQuery,
      groupBy,
      expandedIds,
      contextMenu,
      editingTitleId,
      modalIssueId,
      flattenedIssues,
      focusedIndex,
      // Actions
      selectIssue,
      toggleExpand,
      expandIssue,
      collapseIssue,
      isExpanded,
      setFilter,
      setSearchQuery,
      setGroupBy,
      showContextMenu,
      hideContextMenu,
      startEditingTitle,
      stopEditingTitle,
      openModal,
      closeModal,
      setFlattenedIssues,
      focusNext,
      focusPrevious,
    }),
    [
      selectedIssueId,
      activeFilter,
      searchQuery,
      groupBy,
      expandedIds,
      contextMenu,
      editingTitleId,
      modalIssueId,
      flattenedIssues,
      focusedIndex,
      selectIssue,
      toggleExpand,
      expandIssue,
      collapseIssue,
      isExpanded,
      setFilter,
      setSearchQuery,
      setGroupBy,
      showContextMenu,
      hideContextMenu,
      startEditingTitle,
      stopEditingTitle,
      openModal,
      closeModal,
      setFlattenedIssues,
      focusNext,
      focusPrevious,
    ],
  );

  return (
    <RoadmapContext.Provider value={value}>{children}</RoadmapContext.Provider>
  );
}

export function useRoadmapContext() {
  const context = useContext(RoadmapContext);
  if (!context) {
    throw new Error("useRoadmapContext must be used within a RoadmapProvider");
  }
  return context;
}
