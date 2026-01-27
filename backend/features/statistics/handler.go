// Package statistics provides the statistics vertical slice.
package statistics

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HTTPHandler provides HTTP handlers for statistics operations.
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

// HandleGet handles GET /api/stats - get aggregate statistics.
func (h *HTTPHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := h.service.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// HandleGetMoleculeProgress handles GET /api/stats/molecule/{id} - get molecule progress.
func (h *HTTPHandler) HandleGetMoleculeProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract molecule ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Molecule ID required")
		return
	}
	moleculeID := strings.Join(parts[3:], "/")

	progress, err := h.service.GetMoleculeProgress(r.Context(), moleculeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, progress)
}
