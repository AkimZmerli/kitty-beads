import { useState, useCallback } from "react";
import type { Editor, TLStoreSnapshot } from "tldraw";

interface UseWhiteboardOptions {
  onSave?: (snapshot: TLStoreSnapshot) => void;
}

export function useWhiteboard({ onSave }: UseWhiteboardOptions = {}) {
  const [editor, setEditor] = useState<Editor | null>(null);
  const [isDirty, setIsDirty] = useState(false);

  const handleMount = useCallback((newEditor: Editor) => {
    setEditor(newEditor);

    // Listen for changes
    const unsubscribe = newEditor.store.listen(() => {
      setIsDirty(true);
    });

    return () => unsubscribe();
  }, []);

  const saveSnapshot = useCallback(() => {
    if (editor) {
      const snapshot = editor.store.getStoreSnapshot();
      onSave?.(snapshot);
      setIsDirty(false);
      return snapshot;
    }
    return null;
  }, [editor, onSave]);

  const loadSnapshot = useCallback(
    (snapshot: TLStoreSnapshot) => {
      if (editor) {
        editor.store.loadStoreSnapshot(snapshot);
        setIsDirty(false);
      }
    },
    [editor]
  );

  const clearCanvas = useCallback(() => {
    if (editor) {
      editor.selectAll();
      editor.deleteShapes(editor.getSelectedShapeIds());
      setIsDirty(false);
    }
  }, [editor]);

  return {
    editor,
    isDirty,
    handleMount,
    saveSnapshot,
    loadSnapshot,
    clearCanvas,
  };
}
