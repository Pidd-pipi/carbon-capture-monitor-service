package httpapi

import (
	"encoding/json"
	"errors"
	"example.com/carbon-capture-monitor-service/domain"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"example.com/carbon-capture-monitor-service/validation"
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
		writeError(w, statusErrorStatus(err), "capture unit update failed")
		return
	}
	if err := domain.ValidateStatusUpdate(current.Status, request.Status, updatedAt); err != nil {
		writeError(w, statusErrorStatus(err), err.Error())
		return
	}
	item, err := s.store.UpdateStatus(request.ID, request.Status, updatedAt)
	if err != nil {
		writeError(w, statusErrorStatus(err), "capture unit update failed")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func statusErrorStatus(err error) int {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidStatus):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
