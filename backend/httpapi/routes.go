package httpapi

import (
	"example.com/carbon-capture-monitor-service/health"
	"example.com/carbon-capture-monitor-service/ops"
	"example.com/carbon-capture-monitor-service/readings"
	"example.com/carbon-capture-monitor-service/store"
	"io/fs"
	"net/http"
	"strings"
)

func NewHandler(st *store.Store, alertSvc *ops.OpsService, readingStore *readings.Store, staticFS fs.FS) http.Handler {
	s := &server{store: st, alerts: alertSvc, readings: readingStore}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler("carbon-capture-monitor-service"))
	mux.HandleFunc("/api/capture-units", s.collection)
	mux.HandleFunc("/api/capture-units/status", s.status)
	mux.HandleFunc("/api/alerts", s.alertsDispatch)
	mux.HandleFunc("/api/alerts/", s.alertSub)
	mux.HandleFunc("/api/alerts/status/", s.alertByStatus)
	mux.HandleFunc("/api/readings", s.readingsDispatch)
	mux.HandleFunc("/api/readings/summary", s.readingsSummary)
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	return mux
}

// knownAlertStatuses mirrors the alert state machine's valid statuses for the
// status-scoped listing route.
var knownAlertStatuses = map[string]bool{"active": true, "paused": true, "closed": true}

func (s *server) alertByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	status := strings.TrimPrefix(r.URL.Path, "/api/alerts/status/")
	if !knownAlertStatuses[status] && status != "" {
		writeError(w, http.StatusBadRequest, "unknown alert status")
		return
	}
	page, err := s.alerts.Search(r.Context(), ops.OpsQuery{Status: ops.OpsStatus(status), Page: 1, PageSize: 100})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "alert search failed")
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *server) alertsDispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.alertsCollection(w, r)
	case http.MethodPost:
		s.alertsCreate(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) readingsDispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.readingsList(w, r)
	case http.MethodPost:
		s.readingsPost(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
