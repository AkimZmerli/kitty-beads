import type { LaneConfig } from "./types";

// Tokyo Night themed lane colors with subtle tints
export const LANES: LaneConfig[] = [
  {
    key: "planned",
    title: "Planned",
    bgColor: "rgba(125, 207, 255, 0.08)",
    borderColor: "#7dcfff",
  },
  {
    key: "doing",
    title: "In Progress",
    bgColor: "rgba(224, 175, 104, 0.08)",
    borderColor: "#e0af68",
  },
  {
    key: "for_review",
    title: "Review",
    bgColor: "rgba(122, 162, 247, 0.08)",
    borderColor: "#7aa2f7",
  },
  {
    key: "done",
    title: "Done",
    bgColor: "rgba(158, 206, 106, 0.08)",
    borderColor: "#9ece6a",
  },
];

export const PRIORITY_COLORS = [
  "#f7768e",
  "#ff9e64",
  "#e0af68",
  "#9ece6a",
  "#565f89",
];
