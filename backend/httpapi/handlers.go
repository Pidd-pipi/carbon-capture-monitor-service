package httpapi

import (
	"encoding/json"
	"errors"
	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"example.com/carbon-capture-monitor-service/validation"
	"log"
	"net/http"
	"strings"
	"time"
)

type server struct {
	store    *store.Store
	alerts   *ops.OpsService
	readings *readings.Store
}

func (s *server) collection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string][]domain.CaptureUnit{"items": s.store.List()})
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var request struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&request); err != nil || strings.TrimSpace(request.ID) == "" {
		writeError(w, http.StatusBadRequest, "id and status are required")
		return
	}
	if err := validation.Status(request.Status); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	current, err := s.store.Get(request.ID)
	if err != nil {
		// Preserve the underlying error so errors.Is resolves the sentinel
		// (store.ErrNotFound) instead of collapsing to a generic 500.
		status, message := statusResponse(err)
		logStatusError(r, "lookup", request.ID, err, status)
		writeError(w, status, message)
		return
	}
	if err := domain.ValidateStatusUpdate(current.Status, request.Status, updatedAt); err != nil {
		status, message := statusResponse(err)
		logStatusError(r, "validate", request.ID, err, status)
		writeError(w, status, message)
		return
	}
	item, err := s.store.UpdateStatus(request.ID, request.Status, updatedAt)
	if err != nil {
		status, message := statusResponse(err)
		logStatusError(r, "update", request.ID, err, status)
		writeError(w, status, message)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// statusResponse maps a status-path error to its HTTP status code and a
// client-facing message. Unknown devices yield 404, illegal transitions or
// invalid arguments yield 400, and anything else falls back to 500. The
// underlying error is preserved so errors.Is/As keeps working through the
// returned wrappers.
func statusResponse(err error) (int, string) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return http.StatusNotFound, "capture unit not found"
	case errors.Is(err, domain.ErrInvalidStatus):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "capture unit update failed"
	}
}

// statusErrorStatus keeps the previous mapping surface for any callers that
// still ask for just the code; new code should prefer statusResponse.
func statusErrorStatus(err error) int {
	status, _ := statusResponse(err)
	return status
}

// logStatusError records the failing step so an operator can tell from logs
// whether lookup, validation, or the write itself failed — not just that the
// request returned an error.
func logStatusError(r *http.Request, step, id string, err error, status int) {
	log.Printf("capture-unit status %s failed request_id=%s id=%s status=%d err=%v",
		step, r.Header.Get("X-Request-ID"), id, status, err)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
