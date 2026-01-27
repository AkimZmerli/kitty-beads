import Markdown from "react-markdown";

interface MarkdownViewerProps {
  content: string;
  className?: string;
}

export function MarkdownViewer({ content, className = "" }: MarkdownViewerProps) {
  return (
    <div className={`markdown-content ${className}`}>
      <Markdown>{content}</Markdown>
    </div>
  );
}
