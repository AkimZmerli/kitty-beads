import type { TreeDimensions } from "./types";

export const TREE_DIMENSIONS: TreeDimensions = {
  width: 1200,
  height: 600,
  nodeWidth: 180,
  nodeHeight: 60,
  horizontalSpacing: 220,
  verticalSpacing: 100,
};

export const NODE_COLORS = {
  epic: {
    fill: "rgba(187, 154, 247, 0.2)",
    stroke: "#bb9af7",
    text: "#bb9af7",
  },
  task: {
    fill: "rgba(125, 207, 255, 0.2)",
    stroke: "#7dcfff",
    text: "#7dcfff",
  },
  done: {
    fill: "rgba(158, 206, 106, 0.2)",
    stroke: "#9ece6a",
    text: "#9ece6a",
  },
  inProgress: {
    fill: "rgba(224, 175, 104, 0.2)",
    stroke: "#e0af68",
    text: "#e0af68",
  },
};

export const PRIORITY_COLORS: Record<number, string> = {
  1: "#f7768e",
  2: "#ff9e64",
  3: "#e0af68",
  4: "#9ece6a",
  5: "#565f89",
};
