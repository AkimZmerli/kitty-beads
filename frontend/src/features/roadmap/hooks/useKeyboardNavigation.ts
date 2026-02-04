import { useEffect } from "react";
import { useRoadmapContext } from "../context";
import { KEYBOARD_SHORTCUTS } from "../constants";

/**
 * Hook to handle keyboard navigation in the roadmap
 * j/k: Move selection up/down
 * Enter: Expand/collapse selected issue
 * o: Open modal for selected issue
 * Space: Open context menu for selected issue
 * Escape: Close modal or deselect
 */
export function useKeyboardNavigation() {
  const {
    selectedIssueId,
    modalIssueId,
    contextMenu,
    editingTitleId,
    focusNext,
    focusPrevious,
    toggleExpand,
    openModal,
    closeModal,
    hideContextMenu,
    showContextMenu,
    selectIssue,
  } = useRoadmapContext();

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Skip if typing in input/textarea
      const target = e.target as HTMLElement;
      if (
        target.tagName === "INPUT" ||
        target.tagName === "TEXTAREA" ||
        target.isContentEditable
      ) {
        return;
      }

      // Skip if editing title
      if (editingTitleId) {
        return;
      }

      // Handle Escape first (works even with modal open)
      if (KEYBOARD_SHORTCUTS.CLOSE.includes(e.key)) {
        e.preventDefault();
        if (modalIssueId) {
          closeModal();
        } else if (contextMenu) {
          hideContextMenu();
        } else if (selectedIssueId) {
          selectIssue(null);
        }
        return;
      }

      // Skip remaining shortcuts if modal is open
      if (modalIssueId) {
        return;
      }

      // Navigation
      if (KEYBOARD_SHORTCUTS.MOVE_DOWN.includes(e.key)) {
        e.preventDefault();
        focusNext();
        return;
      }

      if (KEYBOARD_SHORTCUTS.MOVE_UP.includes(e.key)) {
        e.preventDefault();
        focusPrevious();
        return;
      }

      // Expand/collapse
      if (KEYBOARD_SHORTCUTS.EXPAND_COLLAPSE.includes(e.key)) {
        e.preventDefault();
        if (selectedIssueId) {
          toggleExpand(selectedIssueId);
        }
        return;
      }

      // Open modal
      if (KEYBOARD_SHORTCUTS.OPEN_MODAL.includes(e.key)) {
        e.preventDefault();
        if (selectedIssueId) {
          openModal(selectedIssueId);
        }
        return;
      }

      // Open context menu
      if (KEYBOARD_SHORTCUTS.CONTEXT_MENU.includes(e.key)) {
        e.preventDefault();
        if (selectedIssueId) {
          // Find the selected bead element and open context menu at its position
          const beadElement = document.querySelector(
            `[data-issue-id="${selectedIssueId}"]`
          );
          if (beadElement) {
            const rect = beadElement.getBoundingClientRect();
            // Position context menu to the right of the bead, vertically centered
            showContextMenu(rect.right - 100, rect.top + rect.height / 2, selectedIssueId);
          }
        }
        return;
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [
    selectedIssueId,
    modalIssueId,
    contextMenu,
    editingTitleId,
    focusNext,
    focusPrevious,
    toggleExpand,
    openModal,
    closeModal,
    hideContextMenu,
    showContextMenu,
    selectIssue,
  ]);
}
