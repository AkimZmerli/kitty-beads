// Package compaction provides the issue compaction vertical slice.
package compaction

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// HTTPHandler provides HTTP handlers for compaction operations.
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

// HandleCompact handles POST /api/compact - run compaction.
func (h *HTTPHandler) HandleCompact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var args CompactArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	result, err := h.service.Compact(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// HandleStats handles GET /api/compact/stats - get compaction stats.
func (h *HTTPHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	args := CompactStatsArgs{}
	if tier := r.URL.Query().Get("tier"); tier != "" {
		if t, err := strconv.Atoi(tier); err == nil {
			args.Tier = t
		}
	}

	stats, err := h.service.GetStats(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
