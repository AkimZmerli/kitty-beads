// Package kanban provides the kanban/ready work vertical slice.
package kanban

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandler provides HTTP handlers for kanban operations.
type HTTPHandler struct {
	service *Service
}

// NewHTTPHandler creates a new HTTP handler.
func NewHTTPHandler(service *Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// HandleReady handles GET /api/ready - get ready work items.
func (h *HTTPHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filter := parseReadyFilter(r)

	issues, err := h.service.GetReady(r.Context(), filter)
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

	filter := BlockedFilter{
		ParentID: r.URL.Query().Get("parent_id"),
	}

	issues, err := h.service.GetBlocked(r.Context(), filter)
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

	issues, err := h.service.GetStale(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}

// HandleEpicStatus handles GET /api/epics/status - get epic statuses.
func (h *HTTPHandler) HandleEpicStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filter := EpicStatusFilter{
		EligibleOnly: r.URL.Query().Get("eligible_only") == "true",
	}

	epics, err := h.service.GetEpicsEligibleForClosure(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, epics)
}

// HandleKanban handles GET /api/kanban/{featureID} - get kanban board.
func (h *HTTPHandler) HandleKanban(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract feature ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Feature ID required")
		return
	}
	featureID := parts[3]

	board, err := h.service.GetKanbanBoard(r.Context(), featureID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, board)
}

func parseReadyFilter(r *http.Request) ReadyFilter {
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
