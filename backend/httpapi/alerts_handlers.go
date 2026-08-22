package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/carbon-capture-monitor-service/ops"
)

func (s *server) alertsCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := ops.OpsQuery{
		Subject:  r.URL.Query().Get("subject"),
		Status:   ops.OpsStatus(r.URL.Query().Get("status")),
		Priority: ops.OpsPriority(r.URL.Query().Get("priority")),
		Owner:    r.URL.Query().Get("owner"),
		Page:     queryInt(r, "page", 1),
		PageSize: queryInt(r, "page_size", 25),
	}
	// Always query fresh: the dashboard must reflect alerts created since the
	// last refresh, including ones written by concurrent create requests.
	page, err := s.alerts.Search(r.Context(), q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "alert search failed")
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *server) alertsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		Subject  string            `json:"subject"`
		Owner    string            `json:"owner"`
		Priority ops.OpsPriority   `json:"priority"`
		Status   ops.OpsStatus     `json:"status"`
		Labels   map[string]string `json:"labels"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid alert payload")
		return
	}
	if strings.TrimSpace(body.Subject) == "" {
		writeError(w, http.StatusBadRequest, "subject is required")
		return
	}
	record := ops.OpsRecord{
		Subject:  body.Subject,
		Owner:    body.Owner,
		Priority: body.Priority,
		Status:   body.Status,
		Labels:   body.Labels,
	}
	created, err := s.alerts.Create(r.Context(), record)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ops.ErrOpsConflict) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return fallback
	}
	return value
}
