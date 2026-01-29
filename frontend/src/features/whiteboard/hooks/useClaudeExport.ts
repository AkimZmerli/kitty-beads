import { useCallback } from "react";
import type { Editor, TLShape } from "tldraw";

interface ExportResult {
  markdown: string;
  shapes: TLShape[];
}

// Helper to safely get text from shape props
function getShapeText(shape: TLShape): string | null {
  const props = shape.props as Record<string, unknown>;
  if ("text" in props && typeof props.text === "string") {
    return props.text;
  }
  if ("richText" in props && props.richText) {
    // Handle rich text format
    const richText = props.richText as { content?: Array<{ text?: string }> };
    if (richText.content) {
      return richText.content.map((c) => c.text || "").join("");
    }
  }
  return null;
}

// Helper to get geo type
function getGeoType(shape: TLShape): string | null {
  const props = shape.props as Record<string, unknown>;
  if ("geo" in props && typeof props.geo === "string") {
    return props.geo;
  }
  return null;
}

export function useClaudeExport(editor: Editor | null) {
  const exportToMarkdown = useCallback((): ExportResult | null => {
    if (!editor) return null;

    const shapes = editor.getCurrentPageShapes();
    const sections: string[] = [];

    // Group shapes by type
    const notes: TLShape[] = [];
    const texts: TLShape[] = [];
    const geos: TLShape[] = [];
    const arrows: TLShape[] = [];
    const others: TLShape[] = [];

    for (const shape of shapes) {
      switch (shape.type) {
        case "note":
          notes.push(shape);
          break;
        case "text":
          texts.push(shape);
          break;
        case "geo":
          geos.push(shape);
          break;
        case "arrow":
          arrows.push(shape);
          break;
        default:
          others.push(shape);
      }
    }

    // Build markdown
    sections.push("# Whiteboard Export\n");

    if (notes.length > 0) {
      sections.push("## Notes\n");
      for (const note of notes) {
        const text = getShapeText(note) || "(empty note)";
        sections.push(`- ${text}`);
      }
      sections.push("");
    }

    if (texts.length > 0) {
      sections.push("## Text Elements\n");
      for (const textShape of texts) {
        const text = getShapeText(textShape) || "(empty text)";
        sections.push(`- ${text}`);
      }
      sections.push("");
    }

    if (geos.length > 0) {
      sections.push("## Shapes\n");
      for (const geo of geos) {
        const geoType = getGeoType(geo) || "shape";
        const label = getShapeText(geo) || "(no label)";
        sections.push(`- **${geoType}**: ${label}`);
      }
      sections.push("");
    }

    if (arrows.length > 0) {
      sections.push("## Connections\n");
      sections.push(`- ${arrows.length} connection(s) drawn`);
      sections.push("");
    }

    if (others.length > 0) {
      sections.push("## Other Elements\n");
      sections.push(`- ${others.length} other shape(s)`);
      sections.push("");
    }

    // Summary
    sections.push("---\n");
    sections.push(
      `*Exported from Kitty-Beads Whiteboard: ${shapes.length} total elements*`
    );

    return {
      markdown: sections.join("\n"),
      shapes,
    };
  }, [editor]);

  const copyToClipboard = useCallback(async () => {
    const result = exportToMarkdown();
    if (result) {
      await navigator.clipboard.writeText(result.markdown);
      return true;
    }
    return false;
  }, [exportToMarkdown]);

  return {
    exportToMarkdown,
    copyToClipboard,
  };
}
