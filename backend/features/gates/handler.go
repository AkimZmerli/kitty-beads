// Package gates provides the async coordination (gates) vertical slice.
package gates

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HTTPHandler provides HTTP handlers for gate operations.
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

// HandleCreate handles POST /api/gates - create a gate.
func (h *HTTPHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var args CreateArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	result, err := h.service.Create(r.Context(), args, actor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// HandleList handles GET /api/gates - list gates.
func (h *HTTPHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	args := ListArgs{
		All: r.URL.Query().Get("all") == "true",
	}

	gates, err := h.service.List(r.Context(), args)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, gates)
}

// HandleShow handles GET /api/gates/{id} - show a gate.
func (h *HTTPHandler) HandleShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract gate ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "Gate ID required")
		return
	}
	gateID := strings.Join(parts[2:], "/")

	gate, err := h.service.Show(r.Context(), ShowArgs{ID: gateID})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, gate)
}

// HandleClose handles POST /api/gates/{id}/close - close a gate.
func (h *HTTPHandler) HandleClose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract gate ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Gate ID required")
		return
	}
	gateID := strings.Join(parts[2:len(parts)-1], "/")

	var args CloseArgs
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}
	}
	args.ID = gateID

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	if err := h.service.Close(r.Context(), args, actor); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "closed"})
}

// HandleWait handles POST /api/gates/{id}/wait - add waiters to a gate.
func (h *HTTPHandler) HandleWait(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract gate ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Gate ID required")
		return
	}
	gateID := strings.Join(parts[2:len(parts)-1], "/")

	var args WaitArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	args.ID = gateID

	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "http-api"
	}

	result, err := h.service.Wait(r.Context(), args, actor)
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
