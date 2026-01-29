import mermaid from "mermaid";

// Tokyo Night theme configuration for Mermaid diagrams
export const mermaidConfig = {
  startOnLoad: false,
  securityLevel: "strict" as const,
  theme: "base" as const,
  themeVariables: {
    // Primary colors - Tokyo Night palette
    primaryColor: "#7dcfff", // neon-cyan
    primaryTextColor: "#1f2335", // night-surface (dark text on light nodes)
    primaryBorderColor: "#3b4261", // night-border

    // Secondary colors
    secondaryColor: "#bb9af7", // neon-magenta
    secondaryTextColor: "#1f2335",
    secondaryBorderColor: "#3b4261",

    // Tertiary colors
    tertiaryColor: "#f7768e", // neon-pink
    tertiaryTextColor: "#1f2335",
    tertiaryBorderColor: "#3b4261",

    // Background
    background: "#1f2335", // night-surface
    mainBkg: "#1f2335",

    // Text colors
    textColor: "#c0caf5", // text-bright
    lineColor: "#7dcfff", // neon-cyan

    // Node styling
    nodeBorder: "#3b4261",
    nodeTextColor: "#1f2335",

    // Flowchart specific
    edgeLabelBackground: "#1f2335",

    // Sequence diagram
    actorBkg: "#7dcfff",
    actorTextColor: "#1f2335",
    actorBorder: "#3b4261",
    actorLineColor: "#565f89",
    signalColor: "#c0caf5",
    signalTextColor: "#c0caf5",
    labelBoxBkgColor: "#1f2335",
    labelBoxBorderColor: "#3b4261",
    labelTextColor: "#c0caf5",
    loopTextColor: "#c0caf5",
    noteBkgColor: "#bb9af7",
    noteTextColor: "#1f2335",
    noteBorderColor: "#3b4261",
    activationBkgColor: "#3b4261",
    activationBorderColor: "#7dcfff",

    // State diagram
    labelColor: "#c0caf5",

    // Class diagram
    classText: "#c0caf5",

    // Git graph
    git0: "#7dcfff",
    git1: "#bb9af7",
    git2: "#f7768e",
    git3: "#9ece6a", // neon-green
    git4: "#ff9e64", // neon-orange
    git5: "#7aa2f7", // neon-blue
    git6: "#e0af68", // neon-yellow
    git7: "#ff007c", // bright pink
    gitBranchLabel0: "#1f2335",
    gitBranchLabel1: "#1f2335",
    gitBranchLabel2: "#1f2335",
    gitBranchLabel3: "#1f2335",
    commitLabelColor: "#c0caf5",
    commitLabelBackground: "#1f2335",

    // Gantt chart
    sectionBkgColor: "#24283b",
    altSectionBkgColor: "#1f2335",
    gridColor: "#3b4261",
    todayLineColor: "#f7768e",

    // Pie chart
    pie1: "#7dcfff",
    pie2: "#bb9af7",
    pie3: "#f7768e",
    pie4: "#9ece6a",
    pie5: "#ff9e64",
    pie6: "#7aa2f7",
    pie7: "#e0af68",
    pieStrokeColor: "#3b4261",
    pieOuterStrokeColor: "#3b4261",

    // Error styling
    errorBkgColor: "#f7768e",
    errorTextColor: "#1f2335",
  },
  flowchart: {
    curve: "basis" as const,
    padding: 20,
  },
  sequence: {
    diagramMarginX: 50,
    diagramMarginY: 10,
    actorMargin: 50,
    boxMargin: 10,
    boxTextMargin: 5,
    noteMargin: 10,
    messageMargin: 35,
  },
};

// Initialize mermaid with Tokyo Night theme
export function initializeMermaid() {
  mermaid.initialize(mermaidConfig);
}
