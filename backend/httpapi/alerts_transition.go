package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/carbon-capture-monitor-service/ops"
)

func (s *server) alertSub(w http.ResponseWriter, r *http.Request) {
	remainder := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/alerts/"), "/")
	switch remainder {
	case "snapshot":
		s.alertSnapshot(w, r)
		return
	case "rules":
		s.alertRules(w, r)
		return
	}
	parts := strings.Split(remainder, "/")
	if len(parts) != 2 {
		writeError(w, http.StatusNotFound, "alert not found")
		return
	}
	id, action := parts[0], parts[1]
	switch action {
	case "transition":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.alertTransition(w, r, id)
	case "audit":
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.alertAudit(w, r, id)
	default:
		writeError(w, http.StatusNotFound, "alert not found")
	}
}

func (s *server) alertTransition(w http.ResponseWriter, r *http.Request, id string) {
	var body struct {
		ExpectedRevision int           `json:"expected_revision"`
		TargetStatus     ops.OpsStatus `json:"target_status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2048)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid transition payload")
		return
	}
	actor := strings.TrimSpace(r.Header.Get("X-Operator"))
	if actor == "" {
		actor = "web"
	}
	updated, err := s.alerts.Transition(r.Context(), id, body.ExpectedRevision, body.TargetStatus, actor)
	if err != nil {
		status := http.StatusBadRequest
		switch ops.ErrorCode(err) {
		case "not_found":
			status = http.StatusNotFound
		case "conflict":
			status = http.StatusConflict
		case "transition":
			status = http.StatusUnprocessableEntity
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *server) alertAudit(w http.ResponseWriter, r *http.Request, id string) {
	events := s.alerts.Audit(id)
	if len(events) == 0 {
		if _, err := s.alerts.Get(r.Context(), id); err != nil {
			writeError(w, http.StatusNotFound, "alert not found")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string][]ops.OpsEvent{"events": events})
}

func (s *server) alertSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, s.alerts.Snapshot())
}

func (s *server) alertRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rules := ops.Rules()
	severity := strings.ToLower(r.URL.Query().Get("severity"))
	if severity != "" {
		filtered := rules[:0]
		for _, rule := range rules {
			if strings.ToLower(string(rule.Severity)) == severity {
				filtered = append(filtered, rule)
			}
		}
		rules = filtered
	}
	writeJSON(w, http.StatusOK, map[string][]ops.OpsRule{"rules": rules})
}
