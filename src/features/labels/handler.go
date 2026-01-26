// Package labels provides the label management vertical slice.
package labels

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HTTPHandler provides HTTP handlers for label operations.
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

// HandleAdd handles POST /api/labels - add a label to an issue.
func (h *HTTPHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var args AddArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	result, err := h.service.Add(r.Context(), args, actor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleRemove handles DELETE /api/labels - remove a label from an issue.
func (h *HTTPHandler) HandleRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var args RemoveArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	result, err := h.service.Remove(r.Context(), args, actor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleGet handles GET /api/labels/{issueID} - get labels for an issue.
func (h *HTTPHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract issue ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}
	issueID := strings.Join(parts[2:], "/")

	labels, err := h.service.Get(r.Context(), issueID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, labels)
}

// HandleGetByLabel handles GET /api/labels/search?label=X - get issues by label.
func (h *HTTPHandler) HandleGetByLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	label := r.URL.Query().Get("label")
	if label == "" {
		writeError(w, http.StatusBadRequest, "label query parameter required")
		return
	}

	issues, err := h.service.GetIssuesByLabel(r.Context(), label)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, issues)
}
