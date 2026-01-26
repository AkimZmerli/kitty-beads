// Package comments provides the comment management vertical slice.
package comments

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandler provides HTTP handlers for comment operations.
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

// HandleAdd handles POST /api/comments - add a comment to an issue.
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

	if args.Author == "" {
		args.Author = r.Header.Get("X-Actor")
		if args.Author == "" {
			args.Author = "http-api"
		}
	}

	comment, err := h.service.Add(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

// HandleList handles GET /api/comments/{issueID} - list comments for an issue.
func (h *HTTPHandler) HandleList(w http.ResponseWriter, r *http.Request) {
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

	comments, err := h.service.List(r.Context(), issueID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

// HandleListEvents handles GET /api/events/{issueID} - list events for an issue.
func (h *HTTPHandler) HandleListEvents(w http.ResponseWriter, r *http.Request) {
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

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	events, err := h.service.ListEvents(r.Context(), issueID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, events)
}
