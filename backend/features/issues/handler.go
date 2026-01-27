// Package issues provides the issue management vertical slice.
package issues

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandler provides HTTP handlers for issue operations.
type HTTPHandler struct {
	service *Service
}

// NewHTTPHandler creates a new HTTP handler.
func NewHTTPHandler(service *Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// HandleList handles GET /api/issues - list all issues.
func (h *HTTPHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	filter := parseListFilterFromQuery(r)

	issues, err := h.service.List(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}

// HandleCreate handles POST /api/issues - create a new issue.
func (h *HTTPHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	var args CreateArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	issue, err := h.service.Create(ctx, args, actor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, issue)
}

// HandleGet handles GET /api/issues/{id} - get a single issue.
func (h *HTTPHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	issueID := extractIssueID(r.URL.Path)
	if issueID == "" {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}

	details, err := h.service.Show(ctx, issueID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, details)
}

// HandleUpdate handles PUT /api/issues/{id} - update an issue.
func (h *HTTPHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	issueID := extractIssueID(r.URL.Path)
	if issueID == "" {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}

	var args UpdateArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	args.ID = issueID

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	issue, err := h.service.Update(ctx, args, actor)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, issue)
}

// HandleDelete handles DELETE /api/issues/{id} - delete an issue.
func (h *HTTPHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	issueID := extractIssueID(r.URL.Path)
	if issueID == "" {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	args := DeleteArgs{
		IDs:    []string{issueID},
		Force:  r.URL.Query().Get("force") == "true",
		DryRun: r.URL.Query().Get("dry_run") == "true",
	}

	result, err := h.service.Delete(ctx, args, actor)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleClose handles POST /api/issues/{id}/close - close an issue.
func (h *HTTPHandler) HandleClose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()

	// Extract issue ID from path: /api/issues/{id}/close
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}
	issueID := strings.Join(pathParts[2:len(pathParts)-1], "/")

	var args CloseArgs
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	}
	args.ID = issueID

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	result, err := h.service.Close(ctx, args, actor)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleReady handles GET /api/ready - get ready work items.
func (h *HTTPHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	filter := parseReadyFilterFromQuery(r)

	issues, err := h.service.GetReady(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}

// HandleBlocked handles GET /api/blocked - get blocked issues.
func (h *HTTPHandler) HandleBlocked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	filter := BlockedFilter{
		ParentID: r.URL.Query().Get("parent_id"),
	}

	issues, err := h.service.GetBlocked(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}

// HandleStale handles GET /api/stale - get stale issues.
func (h *HTTPHandler) HandleStale(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	filter := StaleFilter{
		Status: r.URL.Query().Get("status"),
	}
	if days := r.URL.Query().Get("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil {
			filter.Days = d
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	issues, err := h.service.GetStale(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}

// HandleStats handles GET /api/stats - get statistics.
func (h *HTTPHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	stats, err := h.service.GetStats(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// HandleCount handles GET /api/issues/count - count issues.
func (h *HTTPHandler) HandleCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	filter := parseCountFilterFromQuery(r)

	result, err := h.service.Count(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Helper functions

func extractIssueID(path string) string {
	// Extract issue ID from path: /api/issues/{id} or /api/issues/{id}/...
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return ""
	}
	// Rejoin everything after "issues" to handle IDs with dots
	return strings.Join(parts[2:], "/")
}

func parseListFilterFromQuery(r *http.Request) ListFilter {
	q := r.URL.Query()
	filter := ListFilter{
		Query:               q.Get("query"),
		Status:              q.Get("status"),
		IssueType:           q.Get("type"),
		Assignee:            q.Get("assignee"),
		Label:               q.Get("label"),
		ParentID:            q.Get("parent_id"),
		MolType:             q.Get("mol_type"),
		TitleContains:       q.Get("title_contains"),
		DescriptionContains: q.Get("description_contains"),
		NotesContains:       q.Get("notes_contains"),
		CreatedAfter:        q.Get("created_after"),
		CreatedBefore:       q.Get("created_before"),
		UpdatedAfter:        q.Get("updated_after"),
		UpdatedBefore:       q.Get("updated_before"),
		ClosedAfter:         q.Get("closed_after"),
		ClosedBefore:        q.Get("closed_before"),
		DeferAfter:          q.Get("defer_after"),
		DeferBefore:         q.Get("defer_before"),
		DueAfter:            q.Get("due_after"),
		DueBefore:           q.Get("due_before"),
	}

	if limit := q.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	if priority := q.Get("priority"); priority != "" {
		if p, err := strconv.Atoi(priority); err == nil {
			filter.Priority = &p
		}
	}
	if priorityMin := q.Get("priority_min"); priorityMin != "" {
		if p, err := strconv.Atoi(priorityMin); err == nil {
			filter.PriorityMin = &p
		}
	}
	if priorityMax := q.Get("priority_max"); priorityMax != "" {
		if p, err := strconv.Atoi(priorityMax); err == nil {
			filter.PriorityMax = &p
		}
	}

	filter.EmptyDescription = q.Get("empty_description") == "true"
	filter.NoAssignee = q.Get("no_assignee") == "true"
	filter.NoLabels = q.Get("no_labels") == "true"
	filter.IncludeTemplates = q.Get("include_templates") == "true"
	filter.Deferred = q.Get("deferred") == "true"
	filter.Overdue = q.Get("overdue") == "true"

	if pinned := q.Get("pinned"); pinned != "" {
		p := pinned == "true"
		filter.Pinned = &p
	}
	if ephemeral := q.Get("ephemeral"); ephemeral != "" {
		e := ephemeral == "true"
		filter.Ephemeral = &e
	}

	if labels := q["labels"]; len(labels) > 0 {
		filter.Labels = labels
	}
	if labelsAny := q["labels_any"]; len(labelsAny) > 0 {
		filter.LabelsAny = labelsAny
	}
	if ids := q["ids"]; len(ids) > 0 {
		filter.IDs = ids
	}
	if excludeStatus := q["exclude_status"]; len(excludeStatus) > 0 {
		filter.ExcludeStatus = excludeStatus
	}
	if excludeTypes := q["exclude_types"]; len(excludeTypes) > 0 {
		filter.ExcludeTypes = excludeTypes
	}

	return filter
}

func parseReadyFilterFromQuery(r *http.Request) ReadyFilter {
	q := r.URL.Query()
	filter := ReadyFilter{
		Assignee:   q.Get("assignee"),
		Type:       q.Get("type"),
		SortPolicy: q.Get("sort_policy"),
		ParentID:   q.Get("parent_id"),
		MolType:    q.Get("mol_type"),
	}

	filter.Unassigned = q.Get("unassigned") == "true"
	filter.IncludeDeferred = q.Get("include_deferred") == "true"

	if limit := q.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	if priority := q.Get("priority"); priority != "" {
		if p, err := strconv.Atoi(priority); err == nil {
			filter.Priority = &p
		}
	}

	if labels := q["labels"]; len(labels) > 0 {
		filter.Labels = labels
	}
	if labelsAny := q["labels_any"]; len(labelsAny) > 0 {
		filter.LabelsAny = labelsAny
	}

	return filter
}

func parseCountFilterFromQuery(r *http.Request) CountFilter {
	q := r.URL.Query()
	filter := CountFilter{
		Query:               q.Get("query"),
		Status:              q.Get("status"),
		IssueType:           q.Get("type"),
		Assignee:            q.Get("assignee"),
		GroupBy:             q.Get("group_by"),
		TitleContains:       q.Get("title_contains"),
		DescriptionContains: q.Get("description_contains"),
		NotesContains:       q.Get("notes_contains"),
		CreatedAfter:        q.Get("created_after"),
		CreatedBefore:       q.Get("created_before"),
		UpdatedAfter:        q.Get("updated_after"),
		UpdatedBefore:       q.Get("updated_before"),
		ClosedAfter:         q.Get("closed_after"),
		ClosedBefore:        q.Get("closed_before"),
	}

	if priority := q.Get("priority"); priority != "" {
		if p, err := strconv.Atoi(priority); err == nil {
			filter.Priority = &p
		}
	}
	if priorityMin := q.Get("priority_min"); priorityMin != "" {
		if p, err := strconv.Atoi(priorityMin); err == nil {
			filter.PriorityMin = &p
		}
	}
	if priorityMax := q.Get("priority_max"); priorityMax != "" {
		if p, err := strconv.Atoi(priorityMax); err == nil {
			filter.PriorityMax = &p
		}
	}

	filter.EmptyDescription = q.Get("empty_description") == "true"
	filter.NoAssignee = q.Get("no_assignee") == "true"
	filter.NoLabels = q.Get("no_labels") == "true"

	if labels := q["labels"]; len(labels) > 0 {
		filter.Labels = labels
	}
	if labelsAny := q["labels_any"]; len(labelsAny) > 0 {
		filter.LabelsAny = labelsAny
	}
	if ids := q["ids"]; len(ids) > 0 {
		filter.IDs = ids
	}

	return filter
}
