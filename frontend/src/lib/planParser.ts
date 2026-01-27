// Plan Parsing Utilities for Roadmap View
// Extracts summary, acceptance criteria count, and open questions from markdown plans

export interface ParsedPlan {
  /** First paragraph of the plan as summary */
  summary: string;
  /** Count of acceptance criteria (- [ ] items) */
  acceptanceCriteriaCount: number;
  /** Count of completed acceptance criteria (- [x] items) */
  completedCriteriaCount: number;
  /** List of open questions found in the plan */
  openQuestions: string[];
  /** Whether the plan has any content */
  hasContent: boolean;
}

/**
 * Parse a markdown plan and extract key information for the roadmap view
 */
export function parsePlan(markdown: string | undefined | null): ParsedPlan {
  if (!markdown || markdown.trim() === "") {
    return {
      summary: "",
      acceptanceCriteriaCount: 0,
      completedCriteriaCount: 0,
      openQuestions: [],
      hasContent: false,
    };
  }

  const content = markdown.trim();

  return {
    summary: extractSummary(content),
    acceptanceCriteriaCount: countAcceptanceCriteria(content),
    completedCriteriaCount: countCompletedCriteria(content),
    openQuestions: extractOpenQuestions(content),
    hasContent: true,
  };
}

/**
 * Extract the first paragraph as a summary
 * Stops at first blank line, heading, or list
 */
export function extractSummary(content: string): string {
  const lines = content.split("\n");
  const summaryLines: string[] = [];

  for (const line of lines) {
    const trimmed = line.trim();

    // Skip initial headings
    if (trimmed.startsWith("#") && summaryLines.length === 0) {
      continue;
    }

    // Stop at blank line, heading, or list after we have content
    if (summaryLines.length > 0) {
      if (
        trimmed === "" ||
        trimmed.startsWith("#") ||
        trimmed.startsWith("- ") ||
        trimmed.startsWith("* ") ||
        /^\d+\./.test(trimmed)
      ) {
        break;
      }
    }

    // Skip blank lines at the beginning
    if (trimmed === "" && summaryLines.length === 0) {
      continue;
    }

    summaryLines.push(trimmed);
  }

  const summary = summaryLines.join(" ").trim();

  // Truncate to reasonable length for display
  const maxLength = 200;
  if (summary.length > maxLength) {
    return summary.substring(0, maxLength).trim() + "...";
  }

  return summary;
}

/**
 * Count acceptance criteria items (- [ ] patterns)
 */
function countAcceptanceCriteria(content: string): number {
  // Match both unchecked and checked items
  const checkboxPattern = /^[\s]*[-*]\s*\[[ x]\]/gim;
  const matches = content.match(checkboxPattern);
  return matches ? matches.length : 0;
}

/**
 * Count completed acceptance criteria (- [x] patterns)
 */
function countCompletedCriteria(content: string): number {
  const checkedPattern = /^[\s]*[-*]\s*\[x\]/gim;
  const matches = content.match(checkedPattern);
  return matches ? matches.length : 0;
}

/**
 * Extract open questions from the plan
 * Looks for:
 * - Lines starting with "?" or "Q:"
 * - Lines under "## Questions" or "## Open Questions" headers
 * - Lines containing "TODO:" or "TBD"
 */
function extractOpenQuestions(content: string): string[] {
  const questions: string[] = [];
  const lines = content.split("\n");

  let inQuestionsSection = false;

  for (const line of lines) {
    const trimmed = line.trim();

    // Check for questions section header
    if (/^#{1,3}\s*(open\s*)?questions?/i.test(trimmed)) {
      inQuestionsSection = true;
      continue;
    }

    // Exit questions section on next heading
    if (inQuestionsSection && /^#{1,3}\s+/.test(trimmed)) {
      inQuestionsSection = false;
      continue;
    }

    // Collect questions from questions section
    if (inQuestionsSection && trimmed.startsWith("-")) {
      const question = trimmed.replace(/^[-*]\s*/, "").trim();
      if (question) {
        questions.push(question);
      }
      continue;
    }

    // Look for inline questions
    if (trimmed.startsWith("?") || trimmed.toLowerCase().startsWith("q:")) {
      const question = trimmed.replace(/^[?Q]:\s*/i, "").trim();
      if (question) {
        questions.push(question);
      }
      continue;
    }

    // Look for TODO/TBD items
    if (/\bTODO\b|\bTBD\b/i.test(trimmed)) {
      questions.push(trimmed);
    }
  }

  return questions;
}

/**
 * Format acceptance criteria progress as a string
 */
export function formatCriteriaProgress(parsed: ParsedPlan): string {
  if (parsed.acceptanceCriteriaCount === 0) {
    return "No criteria";
  }
  return `${parsed.completedCriteriaCount}/${parsed.acceptanceCriteriaCount} criteria`;
}

/**
 * Get priority color based on priority level
 */
export function getPriorityColor(priority: number): string {
  switch (priority) {
    case 1:
      return "#ef4444"; // red-500
    case 2:
      return "#f97316"; // orange-500
    case 3:
      return "#6b7280"; // gray-500
    default:
      return "#6b7280";
  }
}

/**
 * Get priority indicator (filled or empty circle)
 */
export function getPriorityIndicator(priority: number): string {
  return priority <= 2 ? "●" : "○";
}
