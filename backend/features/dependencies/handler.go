// Package dependencies provides the dependency management vertical slice.
package dependencies

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandler provides HTTP handlers for dependency operations.
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

// HandleAdd handles POST /api/dependencies - add a dependency.
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

// HandleRemove handles DELETE /api/dependencies - remove a dependency.
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

// HandleGet handles GET /api/dependencies/{issueID} - get dependencies for an issue.
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

	deps, err := h.service.GetDependenciesWithMetadata(r.Context(), issueID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, deps)
}

// HandleGetDependents handles GET /api/dependents/{issueID} - get dependents for an issue.
func (h *HTTPHandler) HandleGetDependents(w http.ResponseWriter, r *http.Request) {
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

	deps, err := h.service.GetDependentsWithMetadata(r.Context(), issueID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, deps)
}

// HandleTree handles GET /api/dependencies/tree/{issueID} - get dependency tree.
func (h *HTTPHandler) HandleTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract issue ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Issue ID required")
		return
	}
	issueID := strings.Join(parts[3:], "/")

	args := TreeArgs{
		ID:           issueID,
		ShowAllPaths: r.URL.Query().Get("show_all_paths") == "true",
		Reverse:      r.URL.Query().Get("reverse") == "true",
	}

	if maxDepth := r.URL.Query().Get("max_depth"); maxDepth != "" {
		if d, err := strconv.Atoi(maxDepth); err == nil {
			args.MaxDepth = d
		}
	}

	result, err := h.service.GetTree(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleCycles handles GET /api/dependencies/cycles - detect cycles.
func (h *HTTPHandler) HandleCycles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	result, err := h.service.DetectCycles(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
