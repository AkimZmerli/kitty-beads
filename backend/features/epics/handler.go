// Package epics provides the epic management vertical slice.
package epics

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HTTPHandler provides HTTP handlers for epic operations.
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

// HandleStatus handles GET /api/epics/status - get epic statuses.
func (h *HTTPHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	args := StatusArgs{
		EligibleOnly: r.URL.Query().Get("eligible_only") == "true",
	}

	epics, err := h.service.GetStatus(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, epics)
}

// HandleGetSummary handles GET /api/epics/{id} - get epic summary.
func (h *HTTPHandler) HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract epic ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "Epic ID required")
		return
	}
	epicID := strings.Join(parts[2:], "/")

	summary, err := h.service.GetEpicSummary(r.Context(), epicID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// HandleFeatures handles GET /api/features - list all features/epics.
func (h *HTTPHandler) HandleFeatures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	features, err := h.service.ListFeatures(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, features)
}
