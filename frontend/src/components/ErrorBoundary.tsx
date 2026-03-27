import { Component, type ErrorInfo, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}
interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false };

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error("ErrorBoundary caught:", error, info);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex flex-col items-center justify-center h-screen bg-night-bg text-text-normal gap-4">
          <p className="text-neon-pink text-lg font-bold">
            Something went wrong
          </p>
          <p className="text-text-muted text-sm">{this.state.error?.message}</p>
          <button
            className="px-4 py-2 bg-night-bg-highlight text-neon-cyan rounded border border-night-border hover:border-neon-cyan"
            onClick={() => this.setState({ hasError: false })}
          >
            Try again
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
