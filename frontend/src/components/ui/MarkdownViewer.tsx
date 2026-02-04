import Markdown from "react-markdown";
import { MermaidDiagram } from "./MermaidDiagram";

interface MarkdownViewerProps {
  content: string;
  className?: string;
}

export function MarkdownViewer({
  content,
  className = "",
}: MarkdownViewerProps) {
  return (
    <div className={`markdown-content ${className}`}>
      <Markdown
        components={{
          code({ className, children }) {
            const match = /language-(\w+)/.exec(className || "");
            if (match?.[1] === "mermaid") {
              return <MermaidDiagram chart={String(children)} />;
            }
            return <code className={className}>{children}</code>;
          },
        }}
      >
        {content}
      </Markdown>
    </div>
  );
}
